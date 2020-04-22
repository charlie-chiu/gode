package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/charlie-chiu/gode"
)

func main() {
	client := gode.FakeClient{
		UID: 362907402, //dev angel888
		HID: 32,
		SID: gode.SessionID("197af9c6341e4f846d6defe4da1aaf0489dc15d5"),
	}
	game, err := gode.NewFlash2dbPhpGame("http://103.241.238.141/", 5145)
	if err != nil {
		log.Fatal("error when NewFlash2dbPhpGame", err)
	}
	server := gode.NewServer(client, game)

	fmt.Println("start listen...")
	err = http.ListenAndServe(":80", server)
	if err != nil {
		log.Fatalf("could not start server %v", err)
	}
}
