package v1

import (
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/kolide/kolide/shared/hub"
)

const (
	pingPeriod = (pongWait * 9) / 10
	pongWait   = 60 * time.Second
)

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func websocketUpgrade(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	conn, err := wsupgrader.Upgrade(w, r, nil)

	if err != nil {
		return nil, err
	}

	return conn, nil
}

// Websocket route
func Websocket(c *gin.Context) {
	client, err := websocketUpgrade(c.Writer, c.Request)

	if err != nil {
		log.Warn(err)
		return
	}

	defer client.Close()

	conn := &hub.Connection{
		Send: make(chan []byte, 256),
		WS:   client,
	}

	hub.Websocket.Register <- conn

	go conn.WritePump()

	conn.ReadPump()
}
