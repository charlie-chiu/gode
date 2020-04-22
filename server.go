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
	ServerLogin = "onLogin"
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

	ws.writeBinaryMsg(s.makeSendJSON(ServerReady, []byte(`{"event":true,"data":null}`)))

	// keep listen new message and handle it
	wsMsg := make(chan []byte)
	go ws.listenJSON(wsMsg)
	for {
		closed := false
		select {
		case msg, ok := <-wsMsg:
			if ok {
				s.handleMessage(ws, msg)
			} else {
				s.handleDisconnect()
				closed = true
			}
		}

		if closed {
			break
		}
	}
}

func (s *Server) handleDisconnect() {
	uid := s.client.UserID()
	hid := s.client.HallID()
	s.game.OnBalanceExchange(uid, hid)
	s.game.OnLeaveMachine(uid, hid)
}

func (s *Server) handleMessage(ws *wsServer, msg []byte) {
	data := &wsDataReceive{}
	err := json.Unmarshal(msg, data)
	if err != nil {
		log.Println("Json Unmarshal Error: ", err)
	}

	var (
		sid     = s.client.SessionID()
		uid     = s.client.UserID()
		hid     = s.client.HallID()
		bet     = `{"BetLevel":1}`
		betBase = "1:1"
		credit  = 1000
	)

	switch data.Action {
	case ClientLogin:
		const loginResult = `{"event":true,"data":{"COID":2688,"ExchangeRate":1,"GameID":0,"HallID":6,"Sid":"","Test":1,"UserID":0}}`
		ws.writeBinaryMsg(s.makeSendJSON(ServerLogin, []byte(loginResult)))
		ws.writeBinaryMsg(s.makeSendJSON("onTakeMachine", s.game.OnTakeMachine(uid)))
	case ClientOnLoadInfo:
		msg := s.makeSendJSON("onOnLoadInfo2", s.game.OnLoadInfo(uid))
		ws.writeBinaryMsg(msg)
	case ClientGetMachineDetail:
		msg := s.makeSendJSON("onGetMachineDetail", s.game.OnGetMachineDetail(uid))
		ws.writeBinaryMsg(msg)
	case ClientBeginGame:
		msg := s.makeSendJSON("onBeginGame", s.game.BeginGame(sid, bet))
		ws.writeBinaryMsg(msg)
	case ClientExchangeCredit:
		msg := s.makeSendJSON("onCreditExchange", s.game.OnCreditExchange(sid, betBase, credit))
		ws.writeBinaryMsg(msg)
	case ClientExchangeBalance:
		msg := s.makeSendJSON("onBalanceExchange", s.game.OnBalanceExchange(uid, hid))
		ws.writeBinaryMsg(msg)
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
