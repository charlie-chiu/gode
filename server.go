package gode

import (
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Server struct {
	http.Handler
	g Game
}

type Game interface {
	OnReady() string
	OnLogin() string
	OnTakeMachine() string
	OnLoadInfo() string
	OnGetMachineDetail() string
}

func NewServer(g Game) *Server {
	server := new(Server)
	server.g = g

	router := http.NewServeMux()
	router.Handle("/game", http.HandlerFunc(server.demoPageHandler))
	router.Handle("/ws/game", http.HandlerFunc(server.gameHandler))

	server.Handler = router

	return server
}

func (s *Server) demoPageHandler(w http.ResponseWriter, r *http.Request) {
	const demoTemplatePath = "demo.html"
	tmpl, err := template.ParseFiles(demoTemplatePath)
	if err != nil {
		log.Fatalf("problem opening %s %v", demoTemplatePath, err)
	}

	const welcomeMsg = "a simple API demo page"
	tmpl.Execute(w, struct{ WelcomeMsg string }{welcomeMsg})
}

func (s *Server) gameHandler(w http.ResponseWriter, r *http.Request) {
	ws := newWSServer(w, r)
	ws.WriteMessage(websocket.BinaryMessage, []byte(s.g.OnReady()))
	ws.WriteMessage(websocket.BinaryMessage, []byte(s.g.OnLogin()))
	ws.WriteMessage(websocket.BinaryMessage, []byte(s.g.OnTakeMachine()))
	ws.WriteMessage(websocket.BinaryMessage, []byte(s.g.OnLoadInfo()))
	ws.WriteMessage(websocket.BinaryMessage, []byte(s.g.OnGetMachineDetail()))
}
