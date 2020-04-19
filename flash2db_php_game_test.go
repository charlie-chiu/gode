package gode

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFlash2dbPhpGame(t *testing.T) {
	t.Run("constructor should return an error when game path not exist", func(t *testing.T) {
		dummyURL := "127.0.0.1"
		_, err := NewFlash2dbPhpGame(dummyURL, 99888)
		assertError(t, err)
	})

	t.Run("OnTakeMachine get correct url and return result", func(t *testing.T) {
		var userID UserID = 362907402
		gamePath := "/casino.slot.line243.BuBuGaoSheng.machineOccupyAuto/"
		expectedURL := fmt.Sprintf("/amfphp/json.php%s%d", gamePath, userID)

		srv := NewTestingServer(t, expectedURL, `help`)
		defer srv.Close()

		g, err := NewFlash2dbPhpGame(srv.URL, 5145)

		assertNoError(t, err)

		want := []byte(`help`)
		got := g.OnTakeMachine(userID)
		assertByteEqual(t, got, want)
	})

	t.Run("onLoadInfo get correct url and return result", func(t *testing.T) {
		var userID UserID = 362907402
		var gameCode GameCode = 1
		gamePath := "/casino.slot.line243.BuBuGaoSheng.onLoadInfo/"
		expectedURL := fmt.Sprintf("/amfphp/json.php%s%d/%d", gamePath, userID, gameCode)
		fmt.Println(expectedURL)
		srv := NewTestingServer(t, expectedURL, `onLoadInfo`)
		defer srv.Close()

		g, err := NewFlash2dbPhpGame(srv.URL, 5145)

		assertNoError(t, err)

		want := []byte(`onLoadInfo`)
		got := g.OnLoadInfo(userID, gameCode)
		assertByteEqual(t, got, want)
	})

}

func NewTestingServer(t *testing.T, expectedURL string, response string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertURLEqual(t, r, expectedURL)
		fmt.Fprint(w, response)
	}))
}

func assertError(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Errorf("expect an error but not got one")
	}
}

func assertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal("didn't expecting an error but got one", err)
	}
}

func assertURLEqual(t *testing.T, r *http.Request, want string) {
	t.Helper()
	if r.URL.Path != want {
		t.Errorf("wanted URL %q, got %q", want, r.URL)
	}
}
