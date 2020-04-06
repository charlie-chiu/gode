package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/charlie-chiu/gode"
)

func main() {
	server := gode.NewWSServer()

	fmt.Println("start listen...")
	err := http.ListenAndServe(":80", server)
	if err != nil {
		log.Fatalf("could not start server %v", err)
	}
}
