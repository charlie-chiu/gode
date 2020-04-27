package gode

import (
	"encoding/json"
	"fmt"
	"log"
)

const Path5145 = "casino.slot.line243.BuBuGaoSheng."

const (
	MachineOccupyAuto = "machineOccupyAuto"
	OnLoadInfo        = "onLoadInfo"
	GetMachineDetail  = "getMachineDetail"
	CreditExchange    = "creditExchange"
	BeginGame         = "beginGame"
	BalanceExchange   = "balanceExchange"
	MachineLeave      = "machineLeave"
)

type Flash2dbGame struct {
	gameCode   GameCode
	conn       Connector
	funcPrefix string
}

func NewFlash2dbGame(connector Connector, gameType GameType) (*Flash2dbGame, error) {
	path := map[GameType]string{
		5145: Path5145,
	}

	gamePath, ok := path[gameType]
	if !ok {
		return nil, fmt.Errorf("game %d not define", gameType)
	}

	return &Flash2dbGame{
		conn:       connector,
		funcPrefix: gamePath,
	}, nil
}

func (g *Flash2dbGame) TakeMachine(id UserID) json.RawMessage {
	apiResult, _ := g.conn.Connect(g.funcPrefix+MachineOccupyAuto, id)
	type data struct {
		Event bool `json:"event"`
		Data  struct {
			Event    bool `json:"event"`
			GameCode int  `json:"GameCode"`
		} `json:"data"`
	}
	message := &data{}

	err := json.Unmarshal(apiResult, message)
	if err != nil {
		log.Fatalf("problem unmarshal JSON when parsing %s %v", apiResult, err)
	}
	g.gameCode = GameCode(message.Data.GameCode)

	return apiResult
}

func (g *Flash2dbGame) OnLoadInfo(uid UserID) json.RawMessage {
	apiResult, _ := g.conn.Connect(g.funcPrefix+OnLoadInfo, uid, g.gameCode)

	return apiResult
}

func (g *Flash2dbGame) GetMachineDetail(uid UserID) json.RawMessage {
	apiResult, _ := g.conn.Connect(g.funcPrefix+GetMachineDetail, uid, g.gameCode)

	return apiResult
}

func (g *Flash2dbGame) BeginGame(sid SessionID, betInfo BetInfo) json.RawMessage {
	apiResult, _ := g.conn.Connect(g.funcPrefix+BeginGame, sid, g.gameCode, betInfo)

	return apiResult
}

func (g *Flash2dbGame) CreditExchange(sid SessionID, betBase BetBase, credit int) json.RawMessage {
	apiResult, _ := g.conn.Connect(g.funcPrefix+CreditExchange, sid, g.gameCode, betBase, credit)

	return apiResult
}

func (g *Flash2dbGame) BalanceExchange(uid UserID, hid HallID) json.RawMessage {
	apiResult, _ := g.conn.Connect(g.funcPrefix+BalanceExchange, uid, hid, g.gameCode)

	return apiResult
}

func (g *Flash2dbGame) LeaveMachine(uid UserID, hid HallID) json.RawMessage {
	apiResult, _ := g.conn.Connect(g.funcPrefix+MachineLeave, uid, hid, g.gameCode)

	return apiResult
}
