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

const (
	login            = "loginBySid"
	onLoadInfo       = "onLoadInfo2"
	getMachineDetail = "getMachineDetail"
	beginGame        = "beginGame4"
	creditExchange   = "creditExchange"
	balanceExchange  = "balanceExchange"
)

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
	//tell client we are ready to handle msg
	const msgOnReady = `{"action":"ready","result":{"event":true,"data":null}}`
	s.writeBinaryMsg(ws, []byte(msgOnReady))

	wsMsg := make(chan []byte)
	go s.readMessage(ws, wsMsg)

	for {
		select {
		case msg := <-wsMsg:
			s.handleMessage(ws, msg)
		}
	}
}

func (s *Server) readMessage(ws *wsServer, wsMsg chan []byte) {
	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			log.Println("ReadMessage Error: ", err)
			break
		}

		//maybe shouldn't valid JSON here
		if !json.Valid(msg) {
			log.Println("not Valid JSON", string(msg))
			continue
		}

		wsMsg <- msg
	}
}

func (s *Server) handleMessage(ws *wsServer, msg []byte) {
	// todo: 應該寫成獨立的handler 之類的，之後再視需求修改
	const msgOnLogin = `{"action":"onLogin","result":{"event":true,"data":{"COID":2688,"ExchangeRate":1,"GameID":0,"HallID":6,"Sid":"","Test":1,"UserID":0}}}`
	data := &wsData{}
	err := json.Unmarshal(msg, data)
	if err != nil {
		log.Println("Json Unmarshal Error: ", err)
	}

	var (
		sid      SessionID = ""
		uid      UserID    = 455648515
		hid      HallID    = 6
		gameCode GameCode  = 0
		bet      string    = `{"BetLevel":1}`
		betBase  string    = "1:1"
		credit   int       = 1000
	)

	switch data.Action {
	case login:
		s.writeBinaryMsg(ws, []byte(msgOnLogin))
		s.writeBinaryMsg(ws, s.g.OnTakeMachine(uid))
	case onLoadInfo:
		s.writeBinaryMsg(ws, s.g.OnLoadInfo(uid, gameCode))
	case getMachineDetail:
		s.writeBinaryMsg(ws, s.g.OnGetMachineDetail(uid, gameCode))
	case beginGame:
		s.writeBinaryMsg(ws, s.g.BeginGame(sid, gameCode, bet))
	case creditExchange:
		s.writeBinaryMsg(ws, s.g.OnCreditExchange(sid, gameCode, betBase, credit))
	case balanceExchange:
		s.writeBinaryMsg(ws, s.g.OnBalanceExchange(uid, hid, gameCode))
	}
}

func (s *Server) writeBinaryMsg(ws *wsServer, msg []byte) {
	err := ws.WriteMessage(websocket.BinaryMessage, msg)
	if err != nil {
		log.Println("Write Error: ", err)
	}
}
