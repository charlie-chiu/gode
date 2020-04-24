package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/charlie-chiu/gode"
)

func main() {
	host := "http://103.241.238.141/"
	client := gode.NewFlash2dbClient(host)
	game, err := gode.NewFlash2dbGame(host, 5145)
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
