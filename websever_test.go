package go_ws_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	ws "github.com/charlie-chiu/go-ws"
)

func TestGet(t *testing.T) {
	t.Run("/ returns 200", func(t *testing.T) {
		server := ws.NewWSServer()

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
