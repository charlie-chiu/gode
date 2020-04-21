package gode

import (
	"encoding/json"
	"log"
	"net/http"
)

type Server struct {
	http.Handler
	game   Game
	client Client
}

const (
	// actions from client
	ClientLogin            = "loginBySid"
	ClientOnLoadInfo       = "onLoadInfo2"
	ClientGetMachineDetail = "getMachineDetail"
	ClientBeginGame        = "beginGame4"
	ClientExchangeCredit   = "creditExchange"
	ClientExchangeBalance  = "balanceExchange"
)
const (
	// actions to client
	ServerReady = "ready"
)

func NewServer(c Client, g Game) *Server {
	server := &Server{
		game:   g,
		client: c,
	}

	router := http.NewServeMux()
	router.Handle("/ws/game", http.HandlerFunc(server.gameHandler))

	server.Handler = router

	return server
}

type wsDataReceive struct {
	Action string `json:"action"`
}

type wsDataSend struct {
	Action string          `json:"action"`
	Result json.RawMessage `json:"result"`
}

func (s *Server) gameHandler(w http.ResponseWriter, r *http.Request) {
	ws := newWSServer(w, r)

	readyMsg := s.makeSendJSON(ServerReady, []byte(`{"event":true,"data":null}`))
	ws.writeBinaryMsg(ws, readyMsg)

	// keep listen new message and handle it
	wsMsg := make(chan []byte)
	go s.listenMessage(ws, wsMsg)
	for {
		select {
		case msg := <-wsMsg:
			s.handleMessage(ws, msg)
		}
	}
}

func (s *Server) listenMessage(ws *wsServer, wsMsg chan []byte) {
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
		sid           = s.client.SessionID()
		uid           = s.client.UserID()
		hid           = s.client.HallID()
		dummyGameCode = GameCode(0)
		bet           = `{"BetLevel":1}`
		betBase       = "1:1"
		credit        = 1000
	)

	switch data.Action {
	case ClientLogin:
		ws.writeBinaryMsg(ws, []byte(msgOnLogin))
		msg := s.makeSendJSON("onTakeMachine", s.game.OnTakeMachine(uid))
		ws.writeBinaryMsg(ws, msg)
	case ClientOnLoadInfo:
		msg := s.makeSendJSON("onOnLoadInfo2", s.game.OnLoadInfo(uid, dummyGameCode))
		ws.writeBinaryMsg(ws, msg)
	case ClientGetMachineDetail:
		msg := s.makeSendJSON("onGetMachineDetail", s.game.OnGetMachineDetail(uid, dummyGameCode))
		ws.writeBinaryMsg(ws, msg)
	case ClientBeginGame:
		msg := s.makeSendJSON("onBeginGame", s.game.BeginGame(sid, dummyGameCode, bet))
		ws.writeBinaryMsg(ws, msg)
	case ClientExchangeCredit:
		msg := s.makeSendJSON("onCreditExchange", s.game.OnCreditExchange(sid, dummyGameCode, betBase, credit))
		ws.writeBinaryMsg(ws, msg)
	case ClientExchangeBalance:
		msg := s.makeSendJSON("onBalanceExchange", s.game.OnBalanceExchange(uid, hid, dummyGameCode))
		ws.writeBinaryMsg(ws, msg)
	}
}

func (s *Server) makeSendJSON(action string, APIResult json.RawMessage) json.RawMessage {
	msg, err := json.Marshal(&wsDataSend{
		Action: action,
		Result: APIResult,
	})
	if err != nil {
		log.Printf("Problem marshal JSON: action %q, APIResult %v, %v", action, APIResult, err)
	}
	return msg
}
