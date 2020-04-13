package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/charlie-chiu/gode"
)

func main() {
	game := gode.FakePhpGame{}
	server := gode.NewServer(game)

	fmt.Println("start listen...")
	err := http.ListenAndServe(":80", server)
	if err != nil {
		log.Fatalf("could not start server %v", err)
	}
}
