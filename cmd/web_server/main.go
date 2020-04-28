package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/charlie-chiu/gode"
)

func main() {
	host := "http://103.241.238.141/"
	flash2dbConnector := gode.NewFlash2dbConnector(host)
	client := gode.NewFlash2dbClient(flash2dbConnector)
	game, err := gode.NewFlash2dbGame(flash2dbConnector, 5145)
	if err != nil {
		log.Fatal("error when NewFlash2dbGame", err)
	}
	stubJackpot := &StubJackpot{
		BroadcastInterval: 1 * time.Minute,
		FetchResult:       json.RawMessage(`[4,3,2,1]`),
	}
	server := gode.NewServer(client, game, stubJackpot)

	fmt.Println("start listen...")
	err = http.ListenAndServe(":80", server)
	if err != nil {
		log.Fatalf("could not start server %v", err)
	}
}

type StubJackpot struct {
	BroadcastInterval time.Duration
	FetchResult       json.RawMessage
}

func (j *StubJackpot) Interval() time.Duration {
	return j.BroadcastInterval
}

func (j *StubJackpot) Fetch() json.RawMessage {
	return j.FetchResult
}
