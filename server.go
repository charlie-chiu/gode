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
	OnReady() []byte
	OnLogin() []byte
	OnTakeMachine() []byte
	OnLoadInfo() []byte
	OnGetMachineDetail() []byte
	BeginGame() []byte
	OnCreditExchange() []byte
	OnBalanceExchange() []byte
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
	err = tmpl.Execute(w, struct{ WelcomeMsg string }{welcomeMsg})
	if err != nil {
		log.Fatal("template.Execute Error", err)
	}
}

type wsData struct {
	Action string `json:"action"`
}

func (s *Server) gameHandler(w http.ResponseWriter, r *http.Request) {
	ws := newWSServer(w, r)
	writeBinaryMsg(ws, s.g.OnReady())

	for {
		_, bytes, err := ws.ReadMessage()
		if err != nil {
			log.Println("ReadMessage Error: ", err)
			break
		}

		if !json.Valid(bytes) {
			continue
		}

		data := &wsData{}
		err = json.Unmarshal(bytes, data)
		if err != nil {
			log.Println("Json Unmarshal Error: ", err)
			break
		}

		switch data.Action {
		case "loginBySid":
			writeBinaryMsg(ws, s.g.OnLogin())
			writeBinaryMsg(ws, s.g.OnTakeMachine())
		case "onLoadInfo2":
			writeBinaryMsg(ws, s.g.OnLoadInfo())
		case "getMachineDetail":
			writeBinaryMsg(ws, s.g.OnGetMachineDetail())
		case "beginGame4":
			writeBinaryMsg(ws, s.g.BeginGame())
		case "creditExchange":
			writeBinaryMsg(ws, s.g.OnCreditExchange())
		case "balanceExchange":
			writeBinaryMsg(ws, s.g.OnBalanceExchange())
		}
	}
}

func writeBinaryMsg(ws *wsServer, msg []byte) {
	err := ws.WriteMessage(websocket.BinaryMessage, msg)
	if err != nil {
		log.Println("Write Error: ", err)
	}
}
