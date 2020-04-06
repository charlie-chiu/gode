package gode

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type wsServer struct {
	*websocket.Conn
}

func newWSServer(w http.ResponseWriter, r *http.Request) *wsServer {
	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Panicf("problem upgrading connection to web socket %v\n", err)
	}

	return &wsServer{conn}
}

func (w *wsServer) waitForMessage() string {
	_, p, err := w.ReadMessage()
	if err != nil {
		log.Panicf("error reading from web socket %v\n", err)
	}
	msg := string(p)
	//log.Printf("WS got messages type: %d / content : %v", messageType, msg)

	return msg
}

func (w *wsServer) write(p []byte) (n int, err error) {
	err = w.WriteMessage(1, p)
	if err != nil {
		return 0, err
	}

	return len(p), nil
}
