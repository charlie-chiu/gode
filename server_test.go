package gode_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/charlie-chiu/gode"
	"github.com/gorilla/websocket"
)

type StubPhpGame struct {
	LoadInfoResult         string
	TakeMachineResult      string
	GetMachineDetailResult string
	BeginGameResult        string
	BalanceExchangeResult  string
	CreditExchangeResult   string
	LeaveMachineResult     string

	BalanceExchangeCalled bool
	LeaveMachineCalled    bool
}

func (s *StubPhpGame) OnTakeMachine(uid gode.UserID) json.RawMessage {
	return json.RawMessage(s.TakeMachineResult)
}
func (s *StubPhpGame) OnLoadInfo(uid gode.UserID) json.RawMessage {
	return json.RawMessage(s.LoadInfoResult)
}
func (s *StubPhpGame) OnGetMachineDetail(uid gode.UserID) json.RawMessage {
	return json.RawMessage(s.GetMachineDetailResult)
}
func (s *StubPhpGame) OnCreditExchange(sid gode.SessionID, bb string, credit int) json.RawMessage {
	return json.RawMessage(s.CreditExchangeResult)
}
func (s *StubPhpGame) OnBalanceExchange(uid gode.UserID, hid gode.HallID) json.RawMessage {
	s.BalanceExchangeCalled = true
	return json.RawMessage(s.BalanceExchangeResult)
}
func (s *StubPhpGame) BeginGame(sid gode.SessionID, betInfo string) json.RawMessage {
	return json.RawMessage(s.BeginGameResult)
}
func (s *StubPhpGame) OnLeaveMachine(uid gode.UserID, hid gode.HallID) json.RawMessage {
	s.LeaveMachineCalled = true
	return json.RawMessage(s.LeaveMachineResult)
}

type StubClient struct {
	UID gode.UserID
	HID gode.HallID
	SID gode.SessionID
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

func TestWebSocketGame(t *testing.T) {
	const timeOut = time.Second

	t.Run("/ws/game can process game", func(t *testing.T) {
		stubClient := &StubClient{}
		stubGame := &StubPhpGame{
			LoadInfoResult:         `{"event":"LoadInfo"}`,
			TakeMachineResult:      `{"event":"TakeMachine"}`,
			GetMachineDetailResult: `{"event":"MachineDetail"}`,
			BeginGameResult:        `{"event":"BeginGame"}`,
			BalanceExchangeResult:  `{"event":"BalanceExchange"}`,
			CreditExchangeResult:   `{"event":"CreditExchange"}`,
			LeaveMachineResult:     `{"event":"LeaveMachine"}`,
		}
		server := httptest.NewServer(gode.NewServer(stubClient, stubGame))
		url := makeWebSocketURL(server, "/ws/game")
		wsClient := mustDialWS(t, url)
		defer server.Close()

		within(t, timeOut, func() {
			//ready
			assertWSReceiveBinaryMsg(t, wsClient, `{"action":"ready","result":{"event":true,"data":null}}`)

			//ClientLogin
			writeBinaryMsg(t, wsClient, `{"action":"loginBySid","sid":"21d9b36e42c8275a4359f6815b859df05ec2bb0a"}`)
			assertWSReceiveBinaryMsg(t, wsClient, `{"action":"onLogin","result":{"event":true,"data":{"COID":2688,"ExchangeRate":1,"GameID":0,"HallID":6,"Sid":"","Test":1,"UserID":0}}}`)
			assertWSReceiveBinaryMsg(t, wsClient, `{"action":"onTakeMachine","result":{"event":"TakeMachine"}}`)

			//ClientOnLoadInfo
			writeBinaryMsg(t, wsClient, `{"action":"onLoadInfo2","sid":"21d9b36e42c8275a4359f6815b859df05ec2bb0a"}`)
			assertWSReceiveBinaryMsg(t, wsClient, `{"action":"onOnLoadInfo2","result":{"event":"LoadInfo"}}`)

			//ClientGetMachineDetail
			writeBinaryMsg(t, wsClient, `{"action":"getMachineDetail","sid":"21d9b36e42c8275a4359f6815b859df05ec2bb0a"}`)
			assertWSReceiveBinaryMsg(t, wsClient, `{"action":"onGetMachineDetail","result":{"event":"MachineDetail"}}`)

			//開分
			writeBinaryMsg(t, wsClient, `{"action":"creditExchange"}`)
			assertWSReceiveBinaryMsg(t, wsClient, `{"action":"onCreditExchange","result":{"event":"CreditExchange"}}`)

			//begin game
			writeBinaryMsg(t, wsClient, `{"action":"beginGame4"}`)
			assertWSReceiveBinaryMsg(t, wsClient, `{"action":"onBeginGame","result":{"event":"BeginGame"}}`)

			//洗分
			writeBinaryMsg(t, wsClient, `{"action":"balanceExchange"}`)
			assertWSReceiveBinaryMsg(t, wsClient, `{"action":"onBalanceExchange","result":{"event":"BalanceExchange"}}`)
		})

		err := wsClient.Close()
		if err != nil {
			t.Errorf("problem closing dialer %v", err)
		}
	})

	t.Run("should call leaveMachine when ws disconnect", func(t *testing.T) {
		stubClient := &StubClient{}
		stubGame := &StubPhpGame{
			LoadInfoResult:     `{"event":"LoadInfo"}`,
			TakeMachineResult:  `{"event":"TakeMachine"}`,
			LeaveMachineCalled: false,
		}
		svr := httptest.NewServer(gode.NewServer(stubClient, stubGame))
		url := makeWebSocketURL(svr, "/ws/game")
		wsClient := mustDialWS(t, url)
		defer svr.Close()

		//ready
		assertWSReceiveBinaryMsg(t, wsClient, `{"action":"ready","result":{"event":true,"data":null}}`)

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
}

func writeBinaryMsg(t *testing.T, wsClient *websocket.Conn, msg string) {
	err := wsClient.WriteMessage(websocket.BinaryMessage, []byte(msg))
	if err != nil {
		t.Error("ws WriteMessage Error", err)
	}
}

func assertWSReceiveBinaryMsg(t *testing.T, dialer *websocket.Conn, want string) {
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

func TestGet(t *testing.T) {
	t.Run("/ returns 404", func(t *testing.T) {
		stubClient := &StubClient{}
		stubGame := &StubPhpGame{}
		server := gode.NewServer(stubClient, stubGame)

		request, _ := http.NewRequest(http.MethodGet, "/", nil)
		responseRecorder := httptest.NewRecorder()

		server.ServeHTTP(responseRecorder, request)

		assertResponseCode(t, responseRecorder.Code, http.StatusNotFound)
	})
}

func assertResponseCode(t *testing.T, got, expected int) {
	t.Helper()
	if got != expected {
		t.Errorf("expect response status code %d, got %d", expected, got)
	}
}
