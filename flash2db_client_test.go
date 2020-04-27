package gode_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/charlie-chiu/gode"
)

func TestFlash2dbClientBeforeFetchInformation(t *testing.T) {
	client := gode.NewFlash2dbClient(gode.NewFlash2dbConnector(""))

	t.Run("UserID() returns zero value", func(t *testing.T) {
		want := gode.UserID(0)
		got := client.UserID()

		assertUserIDEqual(t, got, want)
	})

	t.Run("HallID returns zero value", func(t *testing.T) {
		want := gode.HallID(0)
		got := client.HallID()

		assertHallIDEqual(t, got, want)
	})

	t.Run("SessionID returns zero value", func(t *testing.T) {
		want := gode.SessionID("")
		got := client.SessionID()

		assertSessionIDEqual(t, got, want)
	})
}

func TestFlash2dbClient_Login(t *testing.T) {
	t.Run("get correct url", func(t *testing.T) {
		// arrange
		spyHandler := &SpyHandler{}
		svr := httptest.NewServer(http.HandlerFunc(spyHandler.spy))
		defer svr.Close()

		client := gode.NewFlash2dbClient(gode.NewFlash2dbConnector(svr.URL))

		sid := gode.SessionID("19870604xi")
		ip := "127.0.0.1"
		expectedURL := fmt.Sprintf(`/amfphp/json.php/Client.loginCheck/%s/%s`, sid, ip)

		// act
		client.Login(sid)

		// assert
		assertURLEqual(t, spyHandler.requestedURL[0], expectedURL)
	})

	t.Run("store updated sid, uid and hid after successful login", func(t *testing.T) {
		sid := gode.SessionID("19870604xi")
		uid := gode.UserID(362907402)
		hid := gode.HallID(32)
		svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
			//client.loginCheck return HallID as string.
			_, _ = fmt.Fprintf(w, `{"data":{"UserID":%d,"Sid":"%s","HallID":"%d"},"event":true}`, uid, sid, hid)
		}))
		defer svr.Close()

		connector := gode.NewFlash2dbConnector(svr.URL)
		client := gode.NewFlash2dbClient(connector)
		client.Login("")

		assertUserIDEqual(t, client.UserID(), uid)
		assertHallIDEqual(t, client.HallID(), hid)
		assertSessionIDEqual(t, client.SessionID(), sid)
	})

	t.Run("login return msg got from flash2db", func(t *testing.T) {
		uid := gode.UserID(9527)
		msg := fmt.Sprintf(`{"data":{"UserID":%d},"event":true}`, uid)
		svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
			_, _ = fmt.Fprintf(w, msg)
		}))
		defer svr.Close()

		client := gode.NewFlash2dbClient(gode.NewFlash2dbConnector(svr.URL))
		got := client.Login("")

		assertRawJSONEqual(t, got, json.RawMessage(msg))
	})
}
