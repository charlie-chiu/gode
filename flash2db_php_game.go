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

type Flash2dbPhpGame struct {
	url      string
	gameCode GameCode
}

func NewFlash2dbPhpGame(host string, gameType GameType) (*Flash2dbPhpGame, error) {
	path := map[GameType]string{
		5145: Path5145,
	}

	gamePath, ok := path[gameType]
	if !ok {
		return nil, fmt.Errorf("game %d not define", gameType)
	}
	url := host + PathPrefix + gamePath

	return &Flash2dbPhpGame{
		url: url,
	}, nil
}

func (g *Flash2dbPhpGame) OnTakeMachine(id UserID) json.RawMessage {
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

func (g *Flash2dbPhpGame) OnLoadInfo(id UserID, gc GameCode) json.RawMessage {
	url := g.generateURL(OnLoadInfo, id, g.gameCode)

	return g.call(url)
}

func (g *Flash2dbPhpGame) OnGetMachineDetail(id UserID, gc GameCode) json.RawMessage {
	url := g.generateURL(GetMachineDetail, id, g.gameCode)

	return g.call(url)
}

func (g *Flash2dbPhpGame) BeginGame(sid SessionID, gameCode GameCode, betInfo string) json.RawMessage {
	u := g.generateURL(BeginGame, sid, g.gameCode, betInfo)

	return g.call(u)
}

func (g *Flash2dbPhpGame) OnCreditExchange(id SessionID, code GameCode, betBase string, credit int) json.RawMessage {
	url := g.generateURL(CreditExchange, id, g.gameCode, betBase, credit)

	return g.call(url)
}

func (g *Flash2dbPhpGame) OnBalanceExchange(uid UserID, hid HallID, code GameCode) json.RawMessage {
	url := g.generateURL(BalanceExchange, uid, hid, g.gameCode)

	return g.call(url)
}

func (g *Flash2dbPhpGame) OnLeaveMachine(uid UserID, hid HallID, gameCode GameCode) json.RawMessage {
	url := g.generateURL(MachineLeave, uid, hid, g.gameCode)

	return g.call(url)
}

func (g *Flash2dbPhpGame) call(url string) json.RawMessage {
	response, err := http.Get(url)
	if err != nil {
		log.Fatal("http Get Error", err)
	}
	defer response.Body.Close()

	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Print("ioutil ReadAll error : ", err)
	}
	return bytes
}

func (g *Flash2dbPhpGame) generateURL(phpFunctionName string, param ...interface{}) string {
	b := strings.Builder{}

	b.WriteString(fmt.Sprintf("%s%s", g.url, phpFunctionName))

	for _, p := range param {
		b.WriteString(fmt.Sprintf("/%v", p))
	}

	return b.String()
}
