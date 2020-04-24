package gode_test

import (
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
		flash2db := gode.NewFlash2db(svr.URL)
		flash2db.Connect("Client.loginCheck", sid, ip)

		expectedURL := fmt.Sprintf(`/amfphp/json.php/Client.loginCheck/%s/%s`, sid, ip)
		assertURLEqual(t, spyHandler.requestedURL[0], expectedURL)
	})
}
