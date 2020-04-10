package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/charlie-chiu/gode"
)

type StubPhpGame struct {
	ReadyMessage            string
	LoginMessage            string
	LoadInfoMessage         string
	TakeMachineMessage      string
	GetMachineDetailMessage string
}

func (s StubPhpGame) OnReady() string {
	return s.ReadyMessage
}

func (s StubPhpGame) OnTakeMachine() string {
	return s.TakeMachineMessage
}

func (s StubPhpGame) OnLoadInfo() string {
	return s.LoadInfoMessage
}

func (s StubPhpGame) OnGetMachineDetail() string {
	return s.GetMachineDetailMessage
}

func (s StubPhpGame) OnLogin() string {
	return s.LoginMessage
}

func main() {
	stubGame := StubPhpGame{
		ReadyMessage:            "OnReady",
		LoginMessage:            "OnLogin",
		LoadInfoMessage:         "OnLoadInfo",
		TakeMachineMessage:      "OnTakeMachine",
		GetMachineDetailMessage: "OnGetMachineDetail",
	}
	server := gode.NewServer(stubGame)

	fmt.Println("start listen...")
	err := http.ListenAndServe(":80", server)
	if err != nil {
		log.Fatalf("could not start server %v", err)
	}
}
