package main

import (
	"fmt"
	"log"
	"net/http"

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
	server := gode.NewServer(client, game)

	fmt.Println("start listen...")
	err = http.ListenAndServe(":80", server)
	if err != nil {
		log.Fatalf("could not start server %v", err)
	}
}
