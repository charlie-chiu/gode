package gode

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
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
	Action    string  `json:"action"`
	SessionID string  `json:"sid"`
	BetBase   string  `json:"rate"`
	Credit    string  `json:"credit"`
	BetInfo   betInfo `json:"betInfo"`
}

type betInfo struct {
	BetLevel int `json:"BetLevel"`
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
		//extract interval and jp msg to Jackpot interface
		case <-time.Tick(300 * time.Millisecond):
			ws.writeBinaryMsg(s.makeSendJSON("updateJP", []byte(`[4,3,2,1]`)))
		}
		if closed {
			break
		}
	}
}

func (s *Server) handleDisconnect() {
	uid := s.client.UserID()
	hid := s.client.HallID()
	s.game.BalanceExchange(uid, hid)
	s.game.LeaveMachine(uid, hid)
}

func (s *Server) handleMessage(ws *wsServer, msg []byte) {
	data := &wsDataReceive{}
	err := json.Unmarshal(msg, data)
	if err != nil {
		log.Println("Json Unmarshal Error: ", err)
	}

	//fmt.Printf("%#v\n", data)
	sid := SessionID(data.SessionID)
	switch data.Action {
	case ClientLogin:
		ws.writeBinaryMsg(s.makeSendJSON(ServerLogin, s.client.Login(sid, "127.0.0.1")))
		ws.writeBinaryMsg(s.makeSendJSON("onTakeMachine", s.game.TakeMachine(s.client.UserID())))
	case ClientOnLoadInfo:
		msg := s.makeSendJSON("onOnLoadInfo2", s.game.OnLoadInfo(s.client.UserID()))
		ws.writeBinaryMsg(msg)
	case ClientGetMachineDetail:
		msg := s.makeSendJSON("onGetMachineDetail", s.game.GetMachineDetail(s.client.UserID()))
		ws.writeBinaryMsg(msg)
	case ClientBeginGame:
		//todo: handle error
		betInfo, _ := parseBetInfo(data)
		msg := s.makeSendJSON("onBeginGame", s.game.BeginGame(sid, betInfo))
		ws.writeBinaryMsg(msg)
	case ClientExchangeCredit:
		//todo: handle error
		credit, _ := parseExchangeCredit(data)
		msg := s.makeSendJSON("onCreditExchange", s.game.CreditExchange(sid, BetBase(data.BetBase), credit))
		ws.writeBinaryMsg(msg)
	case ClientExchangeBalance:
		msg := s.makeSendJSON("onBalanceExchange", s.game.BalanceExchange(s.client.UserID(), s.client.HallID()))
		ws.writeBinaryMsg(msg)
	}
}

func parseExchangeCredit(data *wsDataReceive) (credit int, err error) {
	if data.Credit != "" {
		credit, err = strconv.Atoi(data.Credit)
		if err != nil {
			return 0, fmt.Errorf("parse credit error, %v", err)
		}
	}

	return credit, err
}

func parseBetInfo(data *wsDataReceive) (BetInfo, error) {
	betInfo, err := json.Marshal(data.BetInfo)
	if err != nil {
		return BetInfo(``), fmt.Errorf("parse bet info error %v", err)
	}

	return betInfo, nil
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
