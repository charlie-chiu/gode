package gode_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/charlie-chiu/gode"
)

type SpyHandler struct {
	requestedURL []string
}

func (h *SpyHandler) spy(w http.ResponseWriter, r *http.Request) {
	h.requestedURL = append(h.requestedURL, r.URL.Path)
	_, _ = fmt.Fprint(w, `{"event":true,"data":{"event":true,"GameCode":43}}`)
}

func TestFlash2dbGame(t *testing.T) {
	t.Run("constructor return an error when game path not exist", func(t *testing.T) {
		dummyConnector := &StubConnector{}

		_, err := gode.NewFlash2dbGame(dummyConnector, 99888)
		assertError(t, err)
	})

	const TakeMachineFunction = "casino.slot.line243.BuBuGaoSheng.machineOccupyAuto"
	dummyUID := gode.UserID(23542)
	dummyHID := gode.HallID(6)
	dummySID := gode.SessionID("19870604Xi")
	dummyCredit := 466
	dummyBetBase := gode.BetBase("bb")
	dummyBetInfo := gode.BetInfo(`{"BetLevel":1}`)

	t.Run("TakeMachine return connect result", func(t *testing.T) {
		connector := &StubConnector{
			returnMsg: json.RawMessage(`{"event":true,"data":{"event":true,"GameCode":43}}`),
		}

		g, err := gode.NewFlash2dbGame(connector, 5145)
		assertNoError(t, err)

		want := json.RawMessage(`{"event":true,"data":{"event":true,"GameCode":43}}`)
		got := g.TakeMachine(362907402)
		assertRawJSONEqual(t, got, want)
	})

	t.Run("connect with correct args", func(t *testing.T) {
		spyConnector := &SpyConnector{
			returnMsg: json.RawMessage(fmt.Sprintf(`{"event":true,"data":{"event":true,"GameCode":%d}}`, 431)),
		}
		game, err := gode.NewFlash2dbGame(spyConnector, 5145)
		assertNoError(t, err)

		game.TakeMachine(gode.UserID(362907402))

		expectedCalls := []funcCall{
			{TakeMachineFunction, []interface{}{gode.UserID(362907402)}},
		}

		assertFuncCalledSame(t, expectedCalls, spyConnector.funcCalled)
	})

	t.Run("will using game code after from take machine api", func(t *testing.T) {
		spyConnector := &SpyConnector{}
		game, err := gode.NewFlash2dbGame(spyConnector, 5145)
		assertNoError(t, err)

		uid := gode.UserID(111)
		hid := gode.HallID(6)
		sid := gode.SessionID(`SessionID466`)
		gc := gode.GameCode(43)

		spyConnector.returnMsg = json.RawMessage(fmt.Sprintf(`{"event":true,"data":{"event":true,"GameCode":%d}}`, gc))
		game.TakeMachine(uid)
		game.OnLoadInfo(uid)
		game.GetMachineDetail(uid)
		game.CreditExchange(sid, "1:1", 1000)
		game.BeginGame(sid, gode.BetInfo(`{"BetLevel":1}`))
		game.BalanceExchange(uid, hid)
		game.LeaveMachine(uid, hid)

		const Prefix = "casino.slot.line243.BuBuGaoSheng."
		expectedCalls := []funcCall{
			{Prefix + "machineOccupyAuto", []interface{}{uid}},
			{Prefix + "onLoadInfo", []interface{}{uid, gc}},
			{Prefix + "getMachineDetail", []interface{}{uid, gc}},
			{Prefix + "creditExchange", []interface{}{sid, gc, gode.BetBase("1:1"), 1000}},
			{Prefix + "beginGame", []interface{}{sid, gc, gode.BetInfo(`{"BetLevel":1}`)}},
			{Prefix + "balanceExchange", []interface{}{uid, hid}},
			{Prefix + "machineLeave", []interface{}{uid, hid}},
		}

		assertFuncCalledSame(t, expectedCalls, spyConnector.funcCalled)
	})

	t.Run("TakeMachine return connect result", func(t *testing.T) {
		wantedMsg := json.RawMessage(`{"event":true,"data":{"event":true}}`)
		connector := &StubConnector{
			returnMsg: wantedMsg,
		}

		g, err := gode.NewFlash2dbGame(connector, 5145)
		assertNoError(t, err)

		got := g.GetMachineDetail(dummyUID)
		assertRawJSONEqual(t, got, wantedMsg)
	})

	t.Run("getMachineDetail return connect result", func(t *testing.T) {
		wantedMsg := json.RawMessage(`{"event":true,"data":{"event":true}}`)
		connector := &StubConnector{
			returnMsg: wantedMsg,
		}

		g, err := gode.NewFlash2dbGame(connector, 5145)
		assertNoError(t, err)

		got := g.GetMachineDetail(dummyUID)
		assertRawJSONEqual(t, got, wantedMsg)
	})

	t.Run("creditExchange return connect result", func(t *testing.T) {
		wantedMsg := json.RawMessage(`{"event":true,"data":{"event":true}}`)
		connector := &StubConnector{
			returnMsg: wantedMsg,
		}

		g, err := gode.NewFlash2dbGame(connector, 5145)
		assertNoError(t, err)

		got := g.CreditExchange(dummySID, dummyBetBase, dummyCredit)
		assertRawJSONEqual(t, got, wantedMsg)
	})

	t.Run("beginGame return connect result", func(t *testing.T) {
		wantedMsg := json.RawMessage(`{"event":true,"data":{"event":true}}`)
		connector := &StubConnector{
			returnMsg: wantedMsg,
		}

		g, err := gode.NewFlash2dbGame(connector, 5145)
		assertNoError(t, err)

		got := g.BeginGame(dummySID, dummyBetInfo)
		assertRawJSONEqual(t, got, wantedMsg)
	})

	t.Run("balanceExchange return connect result", func(t *testing.T) {
		wantedMsg := json.RawMessage(`{"event":true,"data":{"event":true}}`)
		connector := &StubConnector{
			returnMsg: wantedMsg,
		}

		g, err := gode.NewFlash2dbGame(connector, 5145)
		assertNoError(t, err)

		got := g.BalanceExchange(dummyUID, dummyHID)
		assertRawJSONEqual(t, got, wantedMsg)
	})

	t.Run("machineLeave return connect result", func(t *testing.T) {
		wantedMsg := json.RawMessage(`{"event":true,"data":{"event":true}}`)
		connector := &StubConnector{
			returnMsg: wantedMsg,
		}

		g, err := gode.NewFlash2dbGame(connector, 5145)
		assertNoError(t, err)

		got := g.LeaveMachine(dummyUID, dummyHID)
		assertRawJSONEqual(t, got, wantedMsg)
	})
}
