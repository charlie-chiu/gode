package main

import (
	"fmt"
	"log"
	"net/http"

	go_ws "github.com/charlie-chiu/go-ws"
)

func main() {
	server := go_ws.NewWSServer()

	fmt.Println("start listen...")
	err := http.ListenAndServe(":80", server)
	if err != nil {
		log.Fatalf("could not start server %v", err)
	}
}
