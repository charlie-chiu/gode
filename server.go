package gode

import (
	"encoding/json"
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
	BeginGame() string
	OnCreditExchange() string
	OnBalanceExchange() string
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

type beginGameSchema struct {
	Action  string                 `json:"action"`
	SID     string                 `json:"sid"`
	BetInfo map[string]interface{} `json:"betInfo"`
}

func (s *Server) gameHandler(w http.ResponseWriter, r *http.Request) {
	ws := newWSServer(w, r)
	ws.WriteMessage(websocket.BinaryMessage, []byte(s.g.OnReady()))
	ws.WriteMessage(websocket.BinaryMessage, []byte(s.g.OnLogin()))
	ws.WriteMessage(websocket.BinaryMessage, []byte(s.g.OnTakeMachine()))
	ws.WriteMessage(websocket.BinaryMessage, []byte(s.g.OnLoadInfo()))
	ws.WriteMessage(websocket.BinaryMessage, []byte(s.g.OnGetMachineDetail()))
	ws.WriteMessage(websocket.BinaryMessage, []byte(s.g.BeginGame()))

	for {
		messageType, bytes, err := ws.ReadMessage()
		if err != nil {
			log.Println("ReadMessage Error: ", err)
			break
		}

		if !json.Valid(bytes) {
			continue
		}

		b := &beginGameSchema{}
		err = json.Unmarshal(bytes, b)
		if err != nil {
			log.Println("Json Unmarshal Error: ", err)
			continue
		}

		//  {"action":"beginGame4"}
		if b.Action == "beginGame4" {
			err = ws.WriteMessage(messageType, []byte(s.g.BeginGame()))
			if err != nil {
				log.Println("Write Error: ", err)
				break
			}
		}
	}
}
