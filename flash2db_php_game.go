package gode

import (
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
)

type Flash2dbPhpGame struct {
	url string
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

func (g *Flash2dbPhpGame) OnTakeMachine(id UserID) []byte {
	url := g.generateURL(MachineOccupyAuto, id)

	return g.call(url)
}

func (g *Flash2dbPhpGame) OnLoadInfo(id UserID, gc GameCode) []byte {
	url := g.generateURL(OnLoadInfo, id, gc)

	return g.call(url)
}

func (Flash2dbPhpGame) OnGetMachineDetail() []byte {
	panic("implement me")
}

func (Flash2dbPhpGame) BeginGame() []byte {
	panic("implement me")
}

func (Flash2dbPhpGame) OnCreditExchange() []byte {
	panic("implement me")
}

func (Flash2dbPhpGame) OnBalanceExchange() []byte {
	panic("implement me")
}

func (g *Flash2dbPhpGame) call(url string) []byte {
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
		b.WriteString(fmt.Sprintf("/%d", p))
	}

	return b.String()
}
