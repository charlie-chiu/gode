package gode

import (
	"encoding/json"
)

type FakePhpGame struct{}

type response struct {
	Action string `json:"action"`
	Result result `json:"result"`
}

type result struct {
	Event bool        `json:"event"`
	Data  interface{} `json:"data"`
}

func (g FakePhpGame) OnReady() string {
	r := &response{
		Action: "onReady",
		Result: result{
			Event: true,
			Data:  nil,
		},
	}
	bytes, _ := json.Marshal(r)

	return string(bytes)
}

func (g FakePhpGame) OnLogin() string {
	return "onLogin"
}

func (g FakePhpGame) OnTakeMachine() string {
	return "OnTakeMachine"
}

func (g FakePhpGame) OnLoadInfo() string {
	return "OnLoadInfo"
}

func (g FakePhpGame) OnGetMachineDetail() string {
	return "OnGetMachineDetail"
}
