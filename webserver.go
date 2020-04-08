package gode

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type WSServer struct {
	http.Handler
}

func NewWSServer() *WSServer {
	server := new(WSServer)

	router := http.NewServeMux()
	router.Handle("/", http.HandlerFunc(server.demoPageHandler))
	router.Handle("/ws/echo", http.HandlerFunc(server.wsEchoHandler))
	router.Handle("/ws/time", http.HandlerFunc(server.wsTimeHandler))

	server.Handler = router

	return server
}

func (s *WSServer) demoPageHandler(w http.ResponseWriter, r *http.Request) {
	const demoTemplatePath = "demo.html"
	tmpl, err := template.ParseFiles(demoTemplatePath)
	if err != nil {
		log.Fatalf("problem opening %s %v", demoTemplatePath, err)
	}

	const welcomeMsg = "a simple API demo page"
	tmpl.Execute(w, struct{ WelcomeMsg string }{welcomeMsg})
}

const (
	echoPrefix = "your message : "
	goodbyeMsg = "goodbye."
)

func (s *WSServer) wsEchoHandler(w http.ResponseWriter, r *http.Request) {
	ws := newWSServer(w, r)

	msg := ws.waitForMessage()

	ws.write([]byte(echoPrefix + msg))

	ws.write([]byte(goodbyeMsg))

	//this well generate close code 1006 at client
	//ws.Close()
	ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "server closed"))
}

func (s *WSServer) wsTimeHandler(w http.ResponseWriter, r *http.Request) {
	const timeFormat = "15:04:05"

	ws := newWSServer(w, r)
	ws.write([]byte(time.Now().Format(timeFormat)))

	ws.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "server closed"))
}
