package gode_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/charlie-chiu/gode"
)

func TestFlash2db_connect(t *testing.T) {
	t.Run("should get correct url", func(t *testing.T) {
		spyHandler := &SpyHandler{}
		svr := httptest.NewServer(http.HandlerFunc(spyHandler.spy))
		defer svr.Close()

		sid := gode.SessionID("19870604xi")
		ip := "127.0.0.1"
		flash2db := gode.NewFlash2dbConnector(svr.URL)
		_, _ = flash2db.Connect("Client.loginCheck", sid, ip)

		expectedURL := fmt.Sprintf(`/amfphp/json.php/Client.loginCheck/%s/%s`, sid, ip)
		assertURLEqual(t, spyHandler.requestedURL[0], expectedURL)
	})

	t.Run("should return response from flash2db", func(t *testing.T) {
		svrResponse := "msg from server"
		svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
			_, _ = fmt.Fprint(w, svrResponse)
		}))
		defer svr.Close()

		flash2db := gode.NewFlash2dbConnector(svr.URL)
		msg, _ := flash2db.Connect("dummy.function.name")

		assertRawJSONEqual(t, msg, json.RawMessage(svrResponse))
	})

	t.Run("should return an error when server return status not found", func(t *testing.T) {
		svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer svr.Close()

		flash2db := gode.NewFlash2dbConnector(svr.URL)
		_, err := flash2db.Connect("dummy.function.name")

		assertError(t, err)
	})
	t.Run("should return an error when server return status error", func(t *testing.T) {
		svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer svr.Close()

		flash2db := gode.NewFlash2dbConnector(svr.URL)
		_, err := flash2db.Connect("dummy.function.name")

		assertError(t, err)
	})
}
