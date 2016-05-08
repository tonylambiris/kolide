package v1

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
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
