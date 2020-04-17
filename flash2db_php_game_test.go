package gode

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFlash2dbPhpGame(t *testing.T) {
	t.Run("OnTakeMachine return []byte", func(t *testing.T) {
		const wantedURL = "/amfphp/json.php/casino.slot.line243.BuBuGaoSheng.machineOccupyAuto/362907402"
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assertURLEqual(t, r, wantedURL)
			fmt.Fprint(w, `help`)
		})
		srv := httptest.NewServer(handler)
		defer srv.Close()

		g := NewFlash2dbPhpGame(srv.URL, 5146)

		want := []byte(`help`)
		got := g.OnTakeMachine()

		assertByteEqual(t, got, want)
	})
}

func assertURLEqual(t *testing.T, r *http.Request, want string) {
	t.Helper()
	if r.URL.Path != want {
		t.Errorf("wanted URL %q, got %q", want, r.URL)
	}
}
