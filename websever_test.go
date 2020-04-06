package gode_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/charlie-chiu/gode"
	"github.com/gorilla/websocket"
)

func TestWebSocketTime(t *testing.T) {
	t.Run("/ws/echo echo user message then disconnect", func(t *testing.T) {

		server := httptest.NewServer(gode.NewWSServer())

		url := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws" + "/echo"
		ws, _, err := websocket.DefaultDialer.Dial(url, nil)
		if err != nil {
			t.Fatalf("could not open a ws connection on %s %v", url, err)
		}

		defer server.Close()
		defer ws.Close()

		want := "your message : msg from test"
		writeWSMessage(t, ws, "msg from test")

		_, bytes, _ := ws.ReadMessage()
		got := string(bytes)
		if got != want {
			t.Errorf("expected %q, got %q", want, got)
		}

		ws.CloseHandler()

		//todo: complete "then disconnect"
	})
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
