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
	ReadyMessage            string
	LoginMessage            string
	LoadInfoMessage         string
	TakeMachineMessage      string
	GetMachineDetailMessage string
	BeginGameMessage        string
}

func (s StubPhpGame) BeginGame() string {
	return s.BeginGameMessage
}

func (s StubPhpGame) OnReady() string {
	return s.ReadyMessage
}

func (s StubPhpGame) OnTakeMachine() string {
	return s.TakeMachineMessage
}

func (s StubPhpGame) OnLoadInfo() string {
	return s.LoadInfoMessage
}

func (s StubPhpGame) OnGetMachineDetail() string {
	return s.GetMachineDetailMessage
}

func (s StubPhpGame) OnLogin() string {
	return s.LoginMessage
}

func TestWebSocketGame(t *testing.T) {
	const timeOut = time.Second
	t.Run("/ws/game receive and return binary", func(t *testing.T) {
		stubGame := StubPhpGame{
			ReadyMessage:            "OnReady",
			LoginMessage:            "OnLogin",
			LoadInfoMessage:         "OnLoadInfo",
			TakeMachineMessage:      "OnTakeMachine",
			GetMachineDetailMessage: "OnGetMachineDetail",
			BeginGameMessage:        "OnBeginGame",
		}
		server := httptest.NewServer(gode.NewServer(stubGame))
		url := makeWebSocketURL(server, "/ws/game")
		dialer := mustDialWS(t, url)
		defer server.Close()

		within(t, timeOut, func() {
			mType := websocket.BinaryMessage
			assertWSReceiveMessage(t, dialer, mType, "OnReady")
			assertWSReceiveMessage(t, dialer, mType, "OnLogin")
			assertWSReceiveMessage(t, dialer, mType, "OnTakeMachine")
			assertWSReceiveMessage(t, dialer, mType, "OnLoadInfo")
			assertWSReceiveMessage(t, dialer, mType, "OnGetMachineDetail")
			assertWSReceiveMessage(t, dialer, mType, "OnBeginGame")
		})

		err := dialer.Close()
		if err != nil {
			t.Errorf("problem closing dialer %v", err)
		}
	})
}

func assertWSReceiveMessage(t *testing.T, dialer *websocket.Conn, expectedType int, want string) {
	t.Helper()

	mt, p, err := dialer.ReadMessage()
	if err != nil {
		t.Fatal("ReadMessageError", err)
	}
	if mt != expectedType {
		t.Errorf("expect got message type %d, got %d", expectedType, mt)
	}
	got := string(p)
	if got != want {
		t.Errorf("expected message %q from web socket, got %q", want, got)
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
