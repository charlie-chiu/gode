package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/charlie-chiu/gode/example"
)

func main() {
	server := example.NewServer()

	fmt.Println("start listen...")
	err := http.ListenAndServe(":80", server)
	if err != nil {
		log.Fatalf("could not start server %v", err)
	}
}
