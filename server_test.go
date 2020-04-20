package gode_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/charlie-chiu/gode"
	"github.com/gorilla/websocket"
)

type StubPhpGame struct {
	LoadInfoMessage         []byte
	TakeMachineMessage      []byte
	GetMachineDetailMessage []byte
	BeginGameMessage        []byte
	BalanceExchangeMsg      []byte
	CreditExchangeMsg       []byte
	LeaveMachineMessage     []byte
}

func (s StubPhpGame) OnTakeMachine(uid gode.UserID) []byte {
	return s.TakeMachineMessage
}

func (s StubPhpGame) OnLoadInfo(uid gode.UserID, gc gode.GameCode) []byte {
	return s.LoadInfoMessage
}

func (s StubPhpGame) OnGetMachineDetail(uid gode.UserID, gc gode.GameCode) []byte {
	return s.GetMachineDetailMessage
}

func (s StubPhpGame) OnCreditExchange(sid gode.SessionID, gc gode.GameCode, bb string, credit int) []byte {
	return s.CreditExchangeMsg
}

func (s StubPhpGame) OnBalanceExchange(uid gode.UserID, hid gode.HallID, gc gode.GameCode) []byte {
	return s.BalanceExchangeMsg
}

func (s StubPhpGame) BeginGame(sid gode.SessionID, gc gode.GameCode, betInfo string) []byte {
	return s.BeginGameMessage
}

func (s StubPhpGame) OnLeaveMachine(uid gode.UserID, hid gode.HallID, gameCode gode.GameCode) []byte {
	return s.LeaveMachineMessage
}

func TestWebSocketGame(t *testing.T) {
	const timeOut = time.Second

	t.Run("/ws/game can process game", func(t *testing.T) {
		stubGame := StubPhpGame{
			LoadInfoMessage:         []byte("OnLoadInfo"),
			TakeMachineMessage:      []byte("OnTakeMachine"),
			GetMachineDetailMessage: []byte("OnGetMachineDetail"),
			BeginGameMessage:        []byte("OnBeginGame"),
			BalanceExchangeMsg:      []byte("balance"),
			CreditExchangeMsg:       []byte("credit"),
			LeaveMachineMessage:     []byte("leave"),
		}
		server := httptest.NewServer(gode.NewServer(stubGame))
		url := makeWebSocketURL(server, "/ws/game")
		wsClient := mustDialWS(t, url)
		defer server.Close()

		within(t, timeOut, func() {
			//ready
			assertWSReceiveBinaryMsg(t, wsClient, `{"action":"ready","result":{"event":true,"data":null}}`)

			//login
			writeBinaryMsg(t, wsClient, `{"action":"loginBySid","sid":"21d9b36e42c8275a4359f6815b859df05ec2bb0a"}`)
			assertWSReceiveBinaryMsg(t, wsClient, `{"action":"onLogin","result":{"event":true,"data":{"COID":2688,"ExchangeRate":1,"GameID":0,"HallID":6,"Sid":"","Test":1,"UserID":0}}}`)
			assertWSReceiveBinaryMsg(t, wsClient, `{"action":"onTakeMachine","result":"OnTakeMachine"}`)

			//onLoadInfo
			writeBinaryMsg(t, wsClient, `{"action":"onLoadInfo2","sid":"21d9b36e42c8275a4359f6815b859df05ec2bb0a"}`)
			assertWSReceiveBinaryMsg(t, wsClient, `{"action":"onOnLoadInfo2","result":"OnLoadInfo"}`)

			//getMachineDetail
			writeBinaryMsg(t, wsClient, `{"action":"getMachineDetail","sid":"21d9b36e42c8275a4359f6815b859df05ec2bb0a"}`)
			assertWSReceiveBinaryMsg(t, wsClient, `{"action":"onGetMachineDetail","result":"OnGetMachineDetail"}`)

			//開分
			writeBinaryMsg(t, wsClient, `{"action":"creditExchange"}`)
			assertWSReceiveBinaryMsg(t, wsClient, `{"action":"onCreditExchange","result":"credit"}`)

			//begin game
			writeBinaryMsg(t, wsClient, `{"action":"beginGame4"}`)
			assertWSReceiveBinaryMsg(t, wsClient, `{"action":"onBeginGame","result":"OnBeginGame"}`)

			//洗分
			writeBinaryMsg(t, wsClient, `{"action":"balanceExchange"}`)
			assertWSReceiveBinaryMsg(t, wsClient, `{"action":"onBalanceExchange","result":"balance"}`)
		})

		err := wsClient.Close()
		if err != nil {
			t.Errorf("problem closing dialer %v", err)
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
		stubGame := StubPhpGame{}
		server := gode.NewServer(stubGame)

		request, _ := http.NewRequest(http.MethodGet, "/", nil)
		responseRecorder := httptest.NewRecorder()

		server.ServeHTTP(responseRecorder, request)

		assertResponseCode(t, responseRecorder.Code, http.StatusNotFound)
	})

	t.Run("/game returns 200", func(t *testing.T) {
		stubGame := StubPhpGame{}
		server := gode.NewServer(stubGame)

		request, _ := http.NewRequest(http.MethodGet, "/game", nil)
		responseRecorder := httptest.NewRecorder()

		server.ServeHTTP(responseRecorder, request)

		assertResponseCode(t, responseRecorder.Code, http.StatusOK)
	})
}

func assertResponseCode(t *testing.T, got, expected int) {
	t.Helper()
	if got != expected {
		t.Errorf("expect response status code %d, got %d", expected, got)
	}
}
