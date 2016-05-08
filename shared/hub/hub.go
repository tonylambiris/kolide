// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hub

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

// Message for websocket transport
type Message struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

// JSON version fo a message
func (n *Message) JSON() []byte {
	b, _ := json.Marshal(n)
	return b
}

// Websocket global
var Websocket = Hub{
	Broadcast:   make(chan []byte),
	Register:    make(chan *Connection),
	Unregister:  make(chan *Connection),
	Connections: make(map[*Connection]bool),
}

// Connection is an middleman between the websocket connection and the hub.
type Connection struct {
	// WS connection.
	WS *websocket.Conn

	// Send buffered channel of outbound messages.
	Send chan []byte
}

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Hub maintains the set of active connections and broadcasts messages to the
// connections.
type Hub struct {
	// Connections registered connections.
	Connections map[*Connection]bool

	// Broadcast inbound messages from the connections.
	Broadcast chan []byte

	// Register requests from the connections.
	Register chan *Connection

	// Unregister requests from connections.
	Unregister chan *Connection
}

// ReadPump pumps messages from the websocket hub.Connection to the hub.
func (c *Connection) ReadPump() {
	defer func() {
		Websocket.Unregister <- c
		c.WS.Close()
	}()
	c.WS.SetReadLimit(maxMessageSize)
	c.WS.SetReadDeadline(time.Now().Add(pongWait))
	c.WS.SetPongHandler(func(string) error { c.WS.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.WS.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}
		Websocket.Broadcast <- message
	}
}

// write writes a message with the given message type and payload.
func (c *Connection) write(mt int, payload []byte) error {
	c.WS.SetWriteDeadline(time.Now().Add(writeWait))
	return c.WS.WriteMessage(mt, payload)
}

// WritePump pumps messages from the hub to the websocket hub.Connection.
func (c *Connection) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.WS.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.write(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

// Run websocket loop
func (h *Hub) Run() {
	for {
		select {
		case c := <-Websocket.Register:
			Websocket.Connections[c] = true
		case c := <-Websocket.Unregister:
			if _, ok := Websocket.Connections[c]; ok {
				delete(Websocket.Connections, c)
				close(c.Send)
			}
		case m := <-Websocket.Broadcast:
			for c := range Websocket.Connections {
				select {
				case c.Send <- m:
				default:
					close(c.Send)
					delete(Websocket.Connections, c)
				}
			}
		}
	}
}
