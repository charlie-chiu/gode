package gode

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestInternalValue(t *testing.T) {
	t.Run("storage GameCode after take machine", func(t *testing.T) {
		response := `{"event":true,"data":{"event":true,"GameCode":43}}`
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, response)
		})
		svr := httptest.NewServer(handler)
		game, _ := NewFlash2dbPhpGame(svr.URL, 5145)
		want := GameCode(0)
		got := game.gameCode
		if got != want {
			t.Errorf("game code zero value error, got %d, want %d", gameCode, want)
		}
		dummyUserID := UserID(0)
		game.OnTakeMachine(dummyUserID)

		want = GameCode(43)
		got = game.gameCode
		if got != want {
			t.Errorf("game code zero value error, got %d, want %d", gameCode, want)
		}
	})
}

func TestFlash2dbPhpGame(t *testing.T) {
	t.Run("constructor should return an error when game path not exist", func(t *testing.T) {
		dummyURL := "127.0.0.1"
		_, err := NewFlash2dbPhpGame(dummyURL, 99888)
		assertError(t, err)
	})

	t.Run("OnTakeMachine get correct url and return result", func(t *testing.T) {
		var userID UserID = 362907402
		gamePath := "/casino.slot.line243.BuBuGaoSheng."
		phpFunctionName := "machineOccupyAuto"
		expectedURL := fmt.Sprintf("%s%s%s/%d", PathPrefix, gamePath, phpFunctionName, userID)
		srv := NewTestingServer(t, expectedURL, `{"event":true,"data":{"event":true,"GameCode":43}}`)
		defer srv.Close()

		g, err := NewFlash2dbPhpGame(srv.URL, 5145)

		assertNoError(t, err)

		want := []byte(`{"event":true,"data":{"event":true,"GameCode":43}}`)
		got := g.OnTakeMachine(userID)
		assertByteEqual(t, got, want)
	})

	t.Run("onLoadInfo get correct url and return result", func(t *testing.T) {
		var userID UserID = 362907402
		var gameCode GameCode = 0
		gamePath := "/casino.slot.line243.BuBuGaoSheng."
		phpFunctionName := "onLoadInfo"
		expectedURL := fmt.Sprintf("%s%s%s/%d/%d", PathPrefix, gamePath, phpFunctionName, userID, gameCode)

		srv := NewTestingServer(t, expectedURL, `onLoadInfo`)
		defer srv.Close()

		g, err := NewFlash2dbPhpGame(srv.URL, 5145)

		assertNoError(t, err)

		want := []byte(`onLoadInfo`)
		got := g.OnLoadInfo(userID, gameCode)
		assertByteEqual(t, got, want)
	})

	t.Run("getMachineDetail get correct url and return result", func(t *testing.T) {

		var userID UserID = 362907402
		var gameCode GameCode = 0
		gamePath := "/casino.slot.line243.BuBuGaoSheng."
		phpFunctionName := "getMachineDetail"
		expectedURL := fmt.Sprintf("%s%s%s/%d/%d", PathPrefix, gamePath, phpFunctionName, userID, gameCode)

		srv := NewTestingServer(t, expectedURL, `getMachineDetail`)
		defer srv.Close()

		g, err := NewFlash2dbPhpGame(srv.URL, 5145)

		assertNoError(t, err)

		want := []byte(`getMachineDetail`)
		got := g.OnGetMachineDetail(userID, gameCode)
		assertByteEqual(t, got, want)
	})

	t.Run("creditExchange get correct url and return result", func(t *testing.T) {
		var sid SessionID = "sidSid123"
		var gameCode GameCode = 0
		var betBase string = "1:5"
		var credit int = 1000
		expectedURL := `/amfphp/json.php/casino.slot.line243.BuBuGaoSheng.creditExchange/sidSid123/0/1:5/1000`

		srv := NewTestingServer(t, expectedURL, `credit`)
		defer srv.Close()

		g, err := NewFlash2dbPhpGame(srv.URL, 5145)

		assertNoError(t, err)

		want := []byte(`credit`)
		got := g.OnCreditExchange(sid, gameCode, betBase, credit)
		assertByteEqual(t, got, want)
	})

	t.Run("beginGame get correct url and return result", func(t *testing.T) {
		var sid SessionID = "sidSid123"
		var gameCode GameCode = 0
		var betInfo string = `{"BetLevel":1}`
		expectedURL := `/amfphp/json.php/casino.slot.line243.BuBuGaoSheng.beginGame/sidSid123/0/{"BetLevel":1}`

		srv := NewTestingServer(t, expectedURL, `begin`)
		defer srv.Close()

		g, err := NewFlash2dbPhpGame(srv.URL, 5145)

		assertNoError(t, err)

		want := []byte(`begin`)
		got := g.BeginGame(sid, gameCode, betInfo)
		assertByteEqual(t, got, want)
	})

	t.Run("balanceExchange get correct url and return result", func(t *testing.T) {
		var userID UserID = 362907402
		var gameCode GameCode = 0
		var hallID HallID = 6
		expectedURL := `/amfphp/json.php/casino.slot.line243.BuBuGaoSheng.balanceExchange/362907402/6/0`

		srv := NewTestingServer(t, expectedURL, `balance`)
		defer srv.Close()

		g, err := NewFlash2dbPhpGame(srv.URL, 5145)

		assertNoError(t, err)

		want := []byte(`balance`)
		got := g.OnBalanceExchange(userID, hallID, gameCode)
		assertByteEqual(t, got, want)
	})

	t.Run("machineLeave get correct url and return result", func(t *testing.T) {
		var userID UserID = 362907402
		var gameCode GameCode = 1
		var hallID HallID = 6
		expectedURL := `/amfphp/json.php/casino.slot.line243.BuBuGaoSheng.machineLeave/362907402/6/0`

		srv := NewTestingServer(t, expectedURL, `leave`)
		defer srv.Close()

		g, err := NewFlash2dbPhpGame(srv.URL, 5145)

		assertNoError(t, err)

		want := []byte(`leave`)
		got := g.OnLeaveMachine(userID, hallID, gameCode)
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
		t.Errorf("URL not matched\n want %q\n, got %q", want, r.URL)
	}
}
