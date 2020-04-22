package gode_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/charlie-chiu/gode"
)

type SpyHandler struct {
	requestedURL []string
}

func (h *SpyHandler) spy(w http.ResponseWriter, r *http.Request) {
	h.requestedURL = append(h.requestedURL, r.URL.Path)
	fmt.Fprint(w, `{"event":true,"data":{"event":true,"GameCode":43}}`)
}

func TestFlash2dbPhpGame(t *testing.T) {
	t.Run("constructor return an error when game path not exist", func(t *testing.T) {
		dummyURL := "127.0.0.1"
		_, err := gode.NewFlash2dbPhpGame(dummyURL, 99888)
		assertError(t, err)
	})

	const PathPrefix = "/amfphp/json.php"

	t.Run("OnTakeMachine get correct url and return result", func(t *testing.T) {
		expectedURL := PathPrefix + `/casino.slot.line243.BuBuGaoSheng.machineOccupyAuto/362907402`
		srv := NewTestingServer(t, expectedURL, `{"event":true,"data":{"event":true,"GameCode":43}}`)
		defer srv.Close()

		g, err := gode.NewFlash2dbPhpGame(srv.URL, 5145)
		assertNoError(t, err)

		want := json.RawMessage(`{"event":true,"data":{"event":true,"GameCode":43}}`)
		got := g.OnTakeMachine(362907402)
		assertRawJSONEqual(t, got, want)
	})

	t.Run("will using game code after from take machine api", func(t *testing.T) {
		spyHandler := &SpyHandler{}
		svr := httptest.NewServer(http.HandlerFunc(spyHandler.spy))
		g, _ := gode.NewFlash2dbPhpGame(svr.URL, 5145)

		UserID := gode.UserID(111)
		hid := gode.HallID(6)
		dummyGameCode := gode.GameCode(0)
		sid := gode.SessionID(`SessionID466`)
		g.OnTakeMachine(UserID)
		g.OnLoadInfo(UserID, dummyGameCode)
		g.OnGetMachineDetail(UserID, dummyGameCode)
		g.OnCreditExchange(sid, dummyGameCode, "1:1", 1000)
		g.BeginGame(sid, dummyGameCode, `{"BetLevel":1}`)
		g.OnBalanceExchange(UserID, hid, dummyGameCode)
		g.OnLeaveMachine(UserID, hid, dummyGameCode)

		const Prefix = "/amfphp/json.php/casino.slot.line243.BuBuGaoSheng."
		expectedURLs := []string{
			Prefix + `machineOccupyAuto/111`,
			Prefix + `onLoadInfo/111/43`,
			Prefix + `getMachineDetail/111/43`,
			Prefix + `creditExchange/SessionID466/43/1:1/1000`,
			Prefix + `beginGame/SessionID466/43/{"BetLevel":1}`,
			Prefix + `balanceExchange/111/6/43`,
			Prefix + `machineLeave/111/6/43`,
		}

		// assert SUT called correct URL
		if !reflect.DeepEqual(expectedURLs, spyHandler.requestedURL) {
			fmt.Printf("expected: %v\n", expectedURLs)
			fmt.Printf("     got: %v\n", spyHandler.requestedURL)
			t.Errorf("URLs not match")
		}
	})

	t.Run("getMachineDetail get correct url and return result", func(t *testing.T) {

		var userID gode.UserID = 362907402
		var gameCode gode.GameCode = 0
		gamePath := "/casino.slot.line243.BuBuGaoSheng."
		phpFunctionName := "getMachineDetail"
		expectedURL := fmt.Sprintf("%s%s%s/%d/%d", PathPrefix, gamePath, phpFunctionName, userID, gameCode)

		srv := NewTestingServer(t, expectedURL, `getMachineDetail`)
		defer srv.Close()

		g, err := gode.NewFlash2dbPhpGame(srv.URL, 5145)

		assertNoError(t, err)

		want := []byte(`getMachineDetail`)
		got := g.OnGetMachineDetail(userID, gameCode)
		assertRawJSONEqual(t, got, want)
	})

	t.Run("creditExchange get correct url and return result", func(t *testing.T) {
		var sid gode.SessionID = "sidSid123"
		var gameCode gode.GameCode = 0
		var betBase string = "1:5"
		var credit int = 1000
		expectedURL := `/amfphp/json.php/casino.slot.line243.BuBuGaoSheng.creditExchange/sidSid123/0/1:5/1000`

		srv := NewTestingServer(t, expectedURL, `credit`)
		defer srv.Close()

		g, err := gode.NewFlash2dbPhpGame(srv.URL, 5145)

		assertNoError(t, err)

		want := []byte(`credit`)
		got := g.OnCreditExchange(sid, gameCode, betBase, credit)
		assertRawJSONEqual(t, got, want)
	})

	t.Run("beginGame get correct url and return result", func(t *testing.T) {
		var sid gode.SessionID = "sidSid123"
		var gameCode gode.GameCode = 0
		var betInfo string = `{"BetLevel":1}`
		expectedURL := `/amfphp/json.php/casino.slot.line243.BuBuGaoSheng.beginGame/sidSid123/0/{"BetLevel":1}`

		srv := NewTestingServer(t, expectedURL, `begin`)
		defer srv.Close()

		g, err := gode.NewFlash2dbPhpGame(srv.URL, 5145)

		assertNoError(t, err)

		want := []byte(`begin`)
		got := g.BeginGame(sid, gameCode, betInfo)
		assertRawJSONEqual(t, got, want)
	})

	t.Run("balanceExchange get correct url and return result", func(t *testing.T) {
		var userID gode.UserID = 362907402
		var gameCode gode.GameCode = 0
		var hallID gode.HallID = 6
		expectedURL := `/amfphp/json.php/casino.slot.line243.BuBuGaoSheng.balanceExchange/362907402/6/0`

		srv := NewTestingServer(t, expectedURL, `balance`)
		defer srv.Close()

		g, err := gode.NewFlash2dbPhpGame(srv.URL, 5145)

		assertNoError(t, err)

		want := []byte(`balance`)
		got := g.OnBalanceExchange(userID, hallID, gameCode)
		assertRawJSONEqual(t, got, want)
	})

	t.Run("machineLeave get correct url and return result", func(t *testing.T) {
		var userID gode.UserID = 362907402
		var gameCode gode.GameCode = 1
		var hallID gode.HallID = 6
		expectedURL := `/amfphp/json.php/casino.slot.line243.BuBuGaoSheng.machineLeave/362907402/6/0`

		srv := NewTestingServer(t, expectedURL, `leave`)
		defer srv.Close()

		g, err := gode.NewFlash2dbPhpGame(srv.URL, 5145)

		assertNoError(t, err)

		want := []byte(`leave`)
		got := g.OnLeaveMachine(userID, hallID, gameCode)
		assertRawJSONEqual(t, got, want)
	})
}

func NewTestingServer(t *testing.T, expectedURL string, response string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertURLEqual(t, r, expectedURL)
		fmt.Fprint(w, response)
	}))
}
