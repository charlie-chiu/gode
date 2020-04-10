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
	Event bool                   `json:"event"`
	Data  map[string]interface{} `json:"data"`
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
	r := &response{
		Action: "onLogin",
		Result: result{
			Event: true,
			Data: map[string]interface{}{
				"COID":         2688,
				"ExchangeRate": 1,
				"GameID":       0,
				"HallID":       6,
				"Sid":          "",
				"Test":         1,
				"UserID":       0,
			},
		},
	}

	bytes, _ := json.Marshal(r)
	return string(bytes)
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
