package gode_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/charlie-chiu/gode"
	"github.com/gorilla/websocket"
)

type SpyPhpGame struct {
	LoadInfoResult         string
	TakeMachineResult      string
	GetMachineDetailResult string
	BeginGameResult        string
	BalanceExchangeResult  string
	CreditExchangeResult   string
	LeaveMachineResult     string

	BalanceExchangeCalled bool
	LeaveMachineCalled    bool

	ReceivedArgs struct {
		SID            gode.SessionID
		BetInfo        gode.BetInfo
		BetBase        gode.BetBase
		exchangeCredit int
	}
}

func (s *SpyPhpGame) TakeMachine(gode.UserID) json.RawMessage {
	return json.RawMessage(s.TakeMachineResult)
}
func (s *SpyPhpGame) OnLoadInfo(gode.UserID) json.RawMessage {
	return json.RawMessage(s.LoadInfoResult)
}
func (s *SpyPhpGame) GetMachineDetail(gode.UserID) json.RawMessage {
	return json.RawMessage(s.GetMachineDetailResult)
}
func (s *SpyPhpGame) CreditExchange(sid gode.SessionID, betBase gode.BetBase, credit int) json.RawMessage {
	s.ReceivedArgs.SID = sid
	s.ReceivedArgs.BetBase = betBase
	s.ReceivedArgs.exchangeCredit = credit
	return json.RawMessage(s.CreditExchangeResult)
}
func (s *SpyPhpGame) BalanceExchange(gode.UserID, gode.HallID) json.RawMessage {
	s.BalanceExchangeCalled = true
	return json.RawMessage(s.BalanceExchangeResult)
}
func (s *SpyPhpGame) BeginGame(sid gode.SessionID, betInfo gode.BetInfo) json.RawMessage {
	s.ReceivedArgs.SID = sid
	s.ReceivedArgs.BetInfo = betInfo
	return json.RawMessage(s.BeginGameResult)
}
func (s *SpyPhpGame) LeaveMachine(gode.UserID, gode.HallID) json.RawMessage {
	s.LeaveMachineCalled = true
	return json.RawMessage(s.LeaveMachineResult)
}

type StubClient struct {
	UID gode.UserID
	HID gode.HallID
	SID gode.SessionID
}

func (c StubClient) Login(gode.SessionID, string) json.RawMessage {
	return json.RawMessage(`{"event":"login"}`)
}
func (c StubClient) UserID() gode.UserID {
	return c.UID
}
func (c StubClient) HallID() gode.HallID {
	return c.HID
}
func (c StubClient) SessionID() gode.SessionID {
	return c.SID
}

type dummyJackpot struct{}

func (j *dummyJackpot) Interval() time.Duration {
	return time.Second
}
func (j *dummyJackpot) Fetch() json.RawMessage {
	return json.RawMessage(`{"event":"dada"}`)
}

func TestGet(t *testing.T) {
	t.Run("/ returns 404", func(t *testing.T) {
		stubClient := &StubClient{}
		stubGame := &SpyPhpGame{}
		server := gode.NewServer(stubClient, stubGame, &dummyJackpot{})

		request, _ := http.NewRequest(http.MethodGet, "/", nil)
		responseRecorder := httptest.NewRecorder()

		server.ServeHTTP(responseRecorder, request)

		assertResponseCode(t, responseRecorder.Code, http.StatusNotFound)
	})
}

func TestWebSocketGame(t *testing.T) {
	const timeOut = time.Second

	t.Run("/ws/game can process game", func(t *testing.T) {
		stubClient := &StubClient{}
		stubGame := &SpyPhpGame{
			LoadInfoResult:         `{"event":"LoadInfo"}`,
			TakeMachineResult:      `{"event":"TakeMachine"}`,
			GetMachineDetailResult: `{"event":"MachineDetail"}`,
			BeginGameResult:        `{"event":"BeginGame"}`,
			BalanceExchangeResult:  `{"event":"BalanceExchange"}`,
			CreditExchangeResult:   `{"event":"CreditExchange"}`,
			LeaveMachineResult:     `{"event":"LeaveMachine"}`,
		}
		server := httptest.NewServer(gode.NewServer(stubClient, stubGame, &dummyJackpot{}))
		wsClient := mustDialWS(t, makeWebSocketURL(server, "/ws/game"))
		defer server.Close()
		defer wsClient.Close()

		within(t, timeOut, func() {
			//ready
			assertReceiveBinaryMsg(t, wsClient, `{"action":"ready","result":{"event":true,"data":null}}`)

			//ClientLogin
			writeBinaryMsg(t, wsClient, `{"action":"loginBySid","sid":"21d9b36e42c8275a4359f6815b859df05ec2bb0a"}`)
			assertReceiveBinaryMsg(t, wsClient, `{"action":"onLogin","result":{"event":"login"}}`)
			assertReceiveBinaryMsg(t, wsClient, `{"action":"onTakeMachine","result":{"event":"TakeMachine"}}`)

			//ClientOnLoadInfo
			writeBinaryMsg(t, wsClient, `{"action":"onLoadInfo2","sid":"21d9b36e42c8275a4359f6815b859df05ec2bb0a"}`)
			assertReceiveBinaryMsg(t, wsClient, `{"action":"onOnLoadInfo2","result":{"event":"LoadInfo"}}`)

			//ClientGetMachineDetail
			writeBinaryMsg(t, wsClient, `{"action":"getMachineDetail","sid":"21d9b36e42c8275a4359f6815b859df05ec2bb0a"}`)
			assertReceiveBinaryMsg(t, wsClient, `{"action":"onGetMachineDetail","result":{"event":"MachineDetail"}}`)

			//開分
			writeBinaryMsg(t, wsClient, `{"action":"creditExchange"}`)
			assertReceiveBinaryMsg(t, wsClient, `{"action":"onCreditExchange","result":{"event":"CreditExchange"}}`)

			//begin game
			writeBinaryMsg(t, wsClient, `{"action":"beginGame4"}`)
			assertReceiveBinaryMsg(t, wsClient, `{"action":"onBeginGame","result":{"event":"BeginGame"}}`)

			//洗分
			writeBinaryMsg(t, wsClient, `{"action":"balanceExchange"}`)
			assertReceiveBinaryMsg(t, wsClient, `{"action":"onBalanceExchange","result":{"event":"BalanceExchange"}}`)
		})
	})

	t.Run("should call leaveMachine when ws disconnect", func(t *testing.T) {
		stubClient := &StubClient{}
		stubGame := &SpyPhpGame{
			LoadInfoResult:     `{"event":"LoadInfo"}`,
			TakeMachineResult:  `{"event":"TakeMachine"}`,
			LeaveMachineCalled: false,
		}
		svr := httptest.NewServer(gode.NewServer(stubClient, stubGame, &dummyJackpot{}))
		url := makeWebSocketURL(svr, "/ws/game")
		wsClient := mustDialWS(t, url)
		defer svr.Close()

		//ready
		assertReceiveBinaryMsg(t, wsClient, `{"action":"ready","result":{"event":true,"data":null}}`)

		//ClientLogin
		writeBinaryMsg(t, wsClient, `{"action":"loginBySid","sid":"21d9b36e42c8275a4359f6815b859df05ec2bb0a"}`)

		_ = wsClient.Close()
		time.Sleep(1 * time.Millisecond)
		if !stubGame.LeaveMachineCalled {
			t.Error("expected game.LeaveMachine called but not")
		}
		if !stubGame.BalanceExchangeCalled {
			t.Error("expected game.BalanceExchange called but not")
		}

	})

	t.Run("forward param from client to php beginGame", func(t *testing.T) {
		stubClient := &StubClient{}
		stubGame := &SpyPhpGame{
			BeginGameResult: `{"action":"beginGame"}`,
		}
		svr := httptest.NewServer(gode.NewServer(stubClient, stubGame, &dummyJackpot{}))
		url := makeWebSocketURL(svr, "/ws/game")
		wsClient := mustDialWS(t, url)
		defer svr.Close()

		//ready
		assertReceiveBinaryMsg(t, wsClient, `{"action":"ready","result":{"event":true,"data":null}}`)

		//beginGame
		wantedSID := gode.SessionID("21d9")
		betInfo := gode.BetInfo(`{"BetLevel":1}`)
		msg := fmt.Sprintf(`{"action":"beginGame4","sid":"%s","betInfo":%s}`, wantedSID, betInfo)
		writeBinaryMsg(t, wsClient, msg)

		time.Sleep(1 * time.Millisecond)
		assertSessionIDEqual(t, stubGame.ReceivedArgs.SID, wantedSID)
		if bytes.Compare(stubGame.ReceivedArgs.BetInfo, betInfo) != 0 {
			t.Errorf("expected stubGame receive BetInfo %#q, got %#q", betInfo, stubGame.ReceivedArgs.BetInfo)
		}
	})

	t.Run("forward param from client to php creditExchange", func(t *testing.T) {
		stubClient := &StubClient{}
		stubGame := &SpyPhpGame{
			CreditExchangeResult: `{"action":"creditExchange"}`,
		}
		svr := httptest.NewServer(gode.NewServer(stubClient, stubGame, &dummyJackpot{}))
		url := makeWebSocketURL(svr, "/ws/game")
		wsClient := mustDialWS(t, url)
		defer svr.Close()

		//ready
		assertReceiveBinaryMsg(t, wsClient, `{"action":"ready","result":{"event":true,"data":null}}`)

		wantedSID := gode.SessionID("21d9")
		betBase := gode.BetBase("1:1")
		// client passing credit in string
		credit := 788
		msg := fmt.Sprintf(`{"action":"creditExchange","sid":"%s", "rate":"%s","credit":"%v"}`, wantedSID, betBase, credit)
		writeBinaryMsg(t, wsClient, msg)

		time.Sleep(1 * time.Millisecond)
		assertSessionIDEqual(t, stubGame.ReceivedArgs.SID, wantedSID)
		if stubGame.ReceivedArgs.BetBase != betBase {
			t.Errorf("expected stubGame receive BetBase %q, got %q", betBase, stubGame.ReceivedArgs.BetBase)
		}
		if stubGame.ReceivedArgs.exchangeCredit != credit {
			t.Errorf("expected stubGame receive credit %d, got %d", credit, stubGame.ReceivedArgs.exchangeCredit)
		}
	})
}

type StubJackpot struct {
	BroadcastInterval time.Duration
	FetchResult       json.RawMessage
}

func (j *StubJackpot) Interval() time.Duration {
	return j.BroadcastInterval
}

func (j *StubJackpot) Fetch() json.RawMessage {
	return j.FetchResult
}

func TestWebSocketJackpot(t *testing.T) {
	t.Run("should send JP info after take machine", func(t *testing.T) {
		stubClient := &StubClient{}
		stubGame := &SpyPhpGame{
			TakeMachineResult: `{"event":"TakeMachine"}`,
		}
		stubJackpot := &StubJackpot{
			BroadcastInterval: 30 * time.Millisecond,
			FetchResult:       json.RawMessage(`[4,3,2,1]`),
		}
		svr := httptest.NewServer(gode.NewServer(stubClient, stubGame, stubJackpot))
		wsClient := mustDialWS(t, makeWebSocketURL(svr, "/ws/game"))
		defer svr.Close()
		defer wsClient.Close()

		//initialize
		assertReceiveBinaryMsg(t, wsClient, `{"action":"ready","result":{"event":true,"data":null}}`)
		writeBinaryMsg(t, wsClient, `{"action":"loginBySid","sid":""}`)
		assertReceiveBinaryMsg(t, wsClient, `{"action":"onLogin","result":{"event":"login"}}`)
		assertReceiveBinaryMsg(t, wsClient, `{"action":"onTakeMachine","result":{"event":"TakeMachine"}}`)

		//should receive 3 Jackpot update with 100 millisecond.
		timeOut := 100 * time.Millisecond
		const UpdateJackpot = `{"action":"updateJP","result":[4,3,2,1]}`
		within(t, timeOut, func() {
			assertReceiveBinaryMsg(t, wsClient, UpdateJackpot)
			assertReceiveBinaryMsg(t, wsClient, UpdateJackpot)
			assertReceiveBinaryMsg(t, wsClient, UpdateJackpot)
		})
	})
}

func writeBinaryMsg(t *testing.T, wsClient *websocket.Conn, msg string) {
	err := wsClient.WriteMessage(websocket.BinaryMessage, []byte(msg))
	if err != nil {
		t.Error("ws WriteMessage Error", err)
	}
}

func assertReceiveBinaryMsg(t *testing.T, dialer *websocket.Conn, want string) {
	t.Helper()
	const expectedType = websocket.BinaryMessage
	mt, p, err := dialer.ReadMessage()
	if err != nil {
		t.Fatal("ReadMessageError", err)
	}
	if mt != expectedType {
		t.Errorf("expect got message type %d, got %d", expectedType, mt)
	}
	got := string(p)
	if got != want {
		t.Errorf("message from web socket not matched\nwant %q\n got %q", want, got)
	}
}

func within(t *testing.T, d time.Duration, assert func()) {
	t.Helper()

	done := make(chan struct{}, 1)

	go func() {
		assert()
		done <- struct{}{}
	}()

	select {
	case <-time.After(d):
		t.Error("timed out")
	case <-done:
	}
}

func mustDialWS(t *testing.T, url string) *websocket.Conn {
	t.Helper()
	dialer, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Fatalf("could not open a ws connection on %s %v", url, err)
	}
	return dialer
}

func makeWebSocketURL(server *httptest.Server, path string) string {
	url := "ws" + strings.TrimPrefix(server.URL, "http") + path
	return url
}

func assertResponseCode(t *testing.T, got, expected int) {
	t.Helper()
	if got != expected {
		t.Errorf("expect response status code %d, got %d", expected, got)
	}
}
