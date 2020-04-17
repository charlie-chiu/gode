package gode

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const PathPrefix = "/amfphp/json.php"
const Path5145 = "/casino.slot.line243.BuBuGaoSheng."

type Flash2dbPhpGame struct {
	url string
}

func NewFlash2dbPhpGame(host string, gameType GameType) *Flash2dbPhpGame {
	path := map[GameType]string{
		5145: Path5145,
	}

	url := host + PathPrefix + path[gameType]

	return &Flash2dbPhpGame{
		url: url,
	}
}

func (g *Flash2dbPhpGame) OnTakeMachine() []byte {
	url := fmt.Sprintf("%smachineOccupyAuto/%d", g.url, 362907402)
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

func (Flash2dbPhpGame) OnLoadInfo() []byte {
	panic("implement me")
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
