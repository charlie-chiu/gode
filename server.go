package gode

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
)

type Server struct {
	http.Handler
	g Game
}

const (
	// action from client
	ClientLogin            = "loginBySid"
	ClientOnLoadInfo       = "onLoadInfo2"
	ClientGetMachineDetail = "getMachineDetail"
	ClientBeginGame        = "beginGame4"
	ClientExchangeCredit   = "creditExchange"
	ClientExchangeBalance  = "balanceExchange"
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

func (s *Server) gameHandler(w http.ResponseWriter, r *http.Request) {
	ws := newWSServer(w, r)
	//tell client we are ready to handle msg
	const msgOnReady = `{"action":"ready","result":{"event":true,"data":null}}`
	ws.writeBinaryMsg(ws, []byte(msgOnReady))

	wsMsg := make(chan []byte)
	go s.readMessage(ws, wsMsg)

	for {
		select {
		case msg := <-wsMsg:
			s.handleMessage(ws, msg)
		}
	}
}

type wsDataReceive struct {
	Action string `json:"action"`
}

type wsDataSend struct {
	Action string          `json:"action"`
	Result json.RawMessage `json:"result"`
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
	data := &wsDataReceive{}
	err := json.Unmarshal(msg, data)
	if err != nil {
		log.Println("Json Unmarshal Error: ", err)
	}

	var (
		sid      SessionID = "b285306cc11c53d9877791427f892d87354bb8a8"
		uid      UserID    = 455648515
		hid      HallID    = 6
		gameCode GameCode  = 36
		bet      string    = `{"BetLevel":1}`
		betBase  string    = "1:1"
		credit   int       = 1000
	)

	switch data.Action {
	case ClientLogin:
		ws.writeBinaryMsg(ws, []byte(msgOnLogin))
		msg := s.makeSendJSON("onTakeMachine", s.g.OnTakeMachine(uid))
		ws.writeBinaryMsg(ws, msg)
	case ClientOnLoadInfo:
		msg := s.makeSendJSON("onOnLoadInfo2", s.g.OnLoadInfo(uid, gameCode))
		ws.writeBinaryMsg(ws, msg)
	case ClientGetMachineDetail:
		msg := s.makeSendJSON("onGetMachineDetail", s.g.OnGetMachineDetail(uid, gameCode))
		ws.writeBinaryMsg(ws, msg)
	case ClientBeginGame:
		msg := s.makeSendJSON("onBeginGame", s.g.BeginGame(sid, gameCode, bet))
		ws.writeBinaryMsg(ws, msg)
	case ClientExchangeCredit:
		msg := s.makeSendJSON("onCreditExchange", s.g.OnCreditExchange(sid, gameCode, betBase, credit))
		ws.writeBinaryMsg(ws, msg)
	case ClientExchangeBalance:
		msg := s.makeSendJSON("onBalanceExchange", s.g.OnBalanceExchange(uid, hid, gameCode))
		ws.writeBinaryMsg(ws, msg)
	}
}

func (s *Server) makeSendJSON(action string, APIResult []byte) []byte {
	msg, err := json.Marshal(&wsDataSend{
		Action: action,
		Result: APIResult,
	})
	if err != nil {
		log.Print("Problem marshal JSON", err)
	}
	return msg
}
