package gode

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

const PathPrefix = "/amfphp/json.php"
const Path5145 = "/casino.slot.line243.BuBuGaoSheng."

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
	url      string
	gameCode GameCode
}

func NewFlash2dbGame(host string, gameType GameType) (*Flash2dbGame, error) {
	path := map[GameType]string{
		5145: Path5145,
	}

	gamePath, ok := path[gameType]
	if !ok {
		return nil, fmt.Errorf("game %d not define", gameType)
	}
	url := host + PathPrefix + gamePath

	return &Flash2dbGame{
		url: url,
	}, nil
}

func (g *Flash2dbGame) TakeMachine(id UserID) json.RawMessage {
	url := g.generateURL(MachineOccupyAuto, id)

	apiResult := g.call(url)
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
	url := g.generateURL(OnLoadInfo, uid, g.gameCode)

	return g.call(url)
}

func (g *Flash2dbGame) GetMachineDetail(uid UserID) json.RawMessage {
	url := g.generateURL(GetMachineDetail, uid, g.gameCode)

	return g.call(url)
}

func (g *Flash2dbGame) BeginGame(sid SessionID, betInfo BetInfo) json.RawMessage {
	u := g.generateURL(BeginGame, sid, g.gameCode, string(betInfo))

	return g.call(u)
}

func (g *Flash2dbGame) CreditExchange(sid SessionID, betBase BetBase, credit int) json.RawMessage {
	url := g.generateURL(CreditExchange, sid, g.gameCode, betBase, credit)

	return g.call(url)
}

func (g *Flash2dbGame) BalanceExchange(uid UserID, hid HallID) json.RawMessage {
	url := g.generateURL(BalanceExchange, uid, hid, g.gameCode)

	return g.call(url)
}

func (g *Flash2dbGame) LeaveMachine(uid UserID, hid HallID) json.RawMessage {
	url := g.generateURL(MachineLeave, uid, hid, g.gameCode)

	return g.call(url)
}

func (g *Flash2dbGame) call(url string) json.RawMessage {
	response, err := http.Get(url)
	if err != nil {
		log.Fatal("http Get Error", err)
	}
	//todo: should we do anything?
	//noinspection GoUnhandledErrorResult
	defer response.Body.Close()

	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Print("ioutil ReadAll error : ", err)
	}
	return bytes
}

func (g *Flash2dbGame) generateURL(phpFunctionName string, param ...interface{}) string {
	b := strings.Builder{}

	b.WriteString(fmt.Sprintf("%s%s", g.url, phpFunctionName))

	for _, p := range param {
		b.WriteString(fmt.Sprintf("/%v", p))
	}

	return b.String()
}
