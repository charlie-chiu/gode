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

func TestWebSocketGame(t *testing.T) {
	const timeOut = time.Second
	t.Run("/ws/game receive and return binary", func(t *testing.T) {
		server := httptest.NewServer(gode.NewWSServer())
		url := makeWebSocketURL(server, "/ws/game")
		dialer := mustDialWS(t, url)
		defer server.Close()

		within(t, timeOut, func() {
			mType := websocket.BinaryMessage
			assertWSReceiveMessage(t, dialer, mType, "onReady")
			assertWSReceiveMessage(t, dialer, mType, "onLogin")
			assertWSReceiveMessage(t, dialer, mType, "onTakeMachine")
			assertWSReceiveMessage(t, dialer, mType, "onLoadInfo")
			assertWSReceiveMessage(t, dialer, mType, "onGetMachineDetail")
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

func TestWebSocketEcho(t *testing.T) {
	//ws server must response with 1 sec
	const timeOut = time.Second
	const echoPrefix = "ECHO: "
	t.Run("/ws/echo echo every received message with ECHO: prefix", func(t *testing.T) {
		server := httptest.NewServer(gode.NewWSServer())
		url := makeWebSocketURL(server, "/ws/echo")
		dialer := mustDialWS(t, url)
		defer server.Close()

		within(t, timeOut, func() {
			msg := "msg from test"
			writeWebSocketMessage(t, dialer, msg)
			assertWSMessage(t, dialer, echoPrefix+msg)
		})
		within(t, timeOut, func() {
			msg := "another message"
			writeWebSocketMessage(t, dialer, msg)
			assertWSMessage(t, dialer, echoPrefix+msg)
		})
		within(t, timeOut, func() {
			msg := "third message"
			writeWebSocketMessage(t, dialer, msg)
			assertWSMessage(t, dialer, echoPrefix+msg)
		})

		err := dialer.Close()
		if err != nil {
			t.Errorf("problem closing dialer %v", err)
		}
	})
}

func TestWebSocketTime(t *testing.T) {
	const timeOut = time.Second
	t.Run("must response then close normally before timeout", func(t *testing.T) {
		server := httptest.NewServer(gode.NewWSServer())
		url := makeWebSocketURL(server, "/ws/time")
		dialer := mustDialWS(t, url)
		defer server.Close()

		within(t, timeOut, func() {
			assertWSMessage(t, dialer, time.Now().Format("15:04:05"))
			assertWSCloseWithExpectError(t, dialer, websocket.CloseNormalClosure)
		})

		err := dialer.Close()
		if err != nil {
			t.Errorf("problem closing dialer %v", err)
		}
	})
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

func writeWebSocketMessage(t *testing.T, conn *websocket.Conn, message string) {
	t.Helper()
	err := conn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		t.Fatalf("could not send message over ws connetcion %v", err)
	}
}

// Close codes defined in RFC 6455, section 11.7.
func assertWSCloseWithExpectError(t *testing.T, dialer *websocket.Conn, closeCode int) {
	t.Helper()
	_, _, err := dialer.ReadMessage()
	if !websocket.IsCloseError(err, closeCode) {
		t.Errorf("expected close code %d, got %v", closeCode, err)
	}
}

func assertWSMessage(t *testing.T, conn *websocket.Conn, want string) {
	t.Helper()
	_, bytes, _ := conn.ReadMessage()
	got := string(bytes)
	if got != want {
		t.Errorf("expected message %q from web socket, got %q", want, got)
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
