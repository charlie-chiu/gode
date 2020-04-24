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
	client := gode.NewFlash2dbClient("")

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
	t.Run("store updated sid, uid and hid after successful login", func(t *testing.T) {
		sid := gode.SessionID("19870604xi")
		uid := gode.UserID(362907402)
		hid := gode.HallID(32)
		svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
			//client.loginCheck return HallID as string.
			_, _ = fmt.Fprintf(w, `{"data":{"UserID":%d,"Sid":"%s","HallID":"%d","GameID":"0","COID":"216310","Test":"0","ExchangeRate":"1","IP":"127.0.0.1"},"event":true}`, uid, sid, hid)
		}))
		defer svr.Close()

		client := gode.NewFlash2dbClient(svr.URL)
		client.Login()

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

		client := gode.NewFlash2dbClient(svr.URL)
		got := client.Login()

		assertRawJSONEqual(t, got, json.RawMessage(msg))
	})
}
