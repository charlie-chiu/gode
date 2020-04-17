package gode

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFlash2dbPhpGame(t *testing.T) {
	t.Run("OnTakeMachine return []byte", func(t *testing.T) {
		const path = "/amfphp/json.php/casino.slot.line243.BuBuGaoSheng.machineOccupyAuto/362907402"
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, `help`)
		})
		srv := httptest.NewServer(handler)
		httptest.
		defer srv.Close()


		host := srv.URL
		host = host + "/hello"
		fmt.Println(host)
		g := NewFlash2dbPhpGame(host, 5146)

		want := []byte(`help`)
		got := g.OnTakeMachine()

		assertByteEqual(t, got, want)
	})
}
