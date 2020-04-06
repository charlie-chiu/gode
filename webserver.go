package gode

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type WSServer struct {
	http.Handler
}

func NewWSServer() *WSServer {
	server := new(WSServer)

	router := http.NewServeMux()
	router.Handle("/", http.HandlerFunc(server.demoPageHandler))
	router.Handle("/ws/echo", http.HandlerFunc(server.wsEchoHandler))

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

func (s *WSServer) wsEchoHandler(w http.ResponseWriter, r *http.Request) {
	ws := newWSServer(w, r)

	msg := ws.waitForMessage()

	ws.write([]byte(fmt.Sprintf("your message : %s", msg)))

	ws.write([]byte("goodbye."))

	ws.Close()
}
