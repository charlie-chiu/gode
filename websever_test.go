package gode_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/charlie-chiu/gode"
	"github.com/gorilla/websocket"
)

func TestWebSocket(t *testing.T) {
	t.Run("/ws/echo echo user message then close normally", func(t *testing.T) {

		server := httptest.NewServer(gode.NewWSServer())

		url := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws" + "/echo"
		dialer, _, err := websocket.DefaultDialer.Dial(url, nil)
		if err != nil {
			t.Fatalf("could not open a ws connection on %s %v", url, err)
		}

		defer server.Close()
		defer dialer.Close()

		want := "your message : msg from test"
		writeWSMessage(t, dialer, "msg from test")

		assertWSMessage(t, dialer, want)
		assertWSMessage(t, dialer, "goodbye.")

		_, _, err = dialer.ReadMessage()
		if !websocket.IsCloseError(err, websocket.CloseNormalClosure) {
			t.Errorf("expected CloseNormalClosure, got %v", err)
		}
	})
}

func assertWSMessage(t *testing.T, conn *websocket.Conn, want string) {
	t.Helper()
	_, bytes, _ := conn.ReadMessage()
	got := string(bytes)
	if got != want {
		t.Errorf("expected message %q from web socket, got %q", want, got)
	}
}

func writeWSMessage(t *testing.T, conn *websocket.Conn, message string) {
	t.Helper()
	err := conn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		t.Fatalf("could not send message over ws connetcion %v", err)
	}
}

func TestGet(t *testing.T) {
	t.Run("/ returns 200", func(t *testing.T) {
		server := gode.NewWSServer()

		request, _ := http.NewRequest(http.MethodGet, "/", nil)
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
