package gode

import (
	"encoding/json"
	"math/rand"
	"time"
)

type FakeGame5145 struct{}

type response struct {
	Action string `json:"action"`
	Result result `json:"result"`
}

type result struct {
	Event bool                   `json:"event"`
	Data  map[string]interface{} `json:"data"`
}

func (g FakeGame5145) TakeMachine(_ UserID) json.RawMessage {
	r := &response{
		Action: "onTakeMachine",
		Result: result{
			Event: true,
			Data: map[string]interface{}{
				"GameCode": 1,
				"gameCode": 1,
				"event":    true,
			},
		},
	}

	bytes, _ := json.Marshal(r)
	return bytes
}

func (g FakeGame5145) OnLoadInfo(_ UserID) json.RawMessage {
	msg := `{"action":"onOnLoadInfo2","result":{"data":{"Balance":99999999.99,"Base":"1:10,1:5,1:2,1:1,2:1,5:1,10:1","BetBase":"","BetTotal":"1000000.0000","Credit":"0.00","Currency":"RMB","DefaultBase":"1:1","ExchangeRate":"1","HallID":6,"LastBetTime":"2019-09-27 00:36:55","LoginName":"fakeuser","PayTotal":"975000.0000","Percent":"99","PercentLow":"96","Roller":{"Normal":[[5,9,10,6,9,10,5,9,9,10,6,9,1,7,9,9,5,10,10,6,1,9,6,10,9,5,8,9,6,10,10,6,9,9,3,8,8,1,9,9,3,8,1,10,5,8,1,5,8,8,3,10,10,6,8,10,3,8,1,9,9,5,10,7,1,6,9,10,6,9,10,4,9,1,10,5,8,8,9,1,8,10,10,5,9,10,1,9,4,10,10,3,9,8],[5,10,10,5,7,7,1,8,4,1,8,8,10,2,7,10,2,7,4,1,7,7,5,8,2,4,7,1,10,10,4,8,8,4,10,9,5,8,8,4,10,7,7,1,10,10,1,6,7,3,8,2,5,10,1,4,8,5,1,7,7,5,3,10,5,2,7,8,1,10,10,4,9,1,10,10,4,8,2,6,10,10,7,5,10,1,5,7,7,4,3,10],[4,7,7,8,4,7,7,1,9,4,7,7,9,5,7,4,5,7,3,8,2,4,10,1,3,8,8,4,1,7,7,5,2,4,9,8,6,1,8,8,2,9,6,7,1,5,9,4,8,1,9,9,3,7,7,1,8,6,9,9,3,2,6,3,7,7,1,6,9,5,4,9,9,1,7,7,2,8,6,9,1,7,4,6,1,8,8,6,4,10,2,3],[4,9,10,4,9,7,1,5,7,4,9,1,4,7,10,1,7,6,1,8,10,4,3,10,6,4,7,5,10,7,4,9,3,4,9,10,7,9,2,2,2,7,9,10,6,7,1,3,8,1,6,5,9,10,4,1,5,7,1,5,10,1,9,10,6,7,9,6,1,7,6,9,5,7,6,10,1,6,7,1,3,9,2,2,4,10,1,8,5,3,10],[5,3,7,1,5,8,7,3,4,8,9,4,3,7,9,5,7,8,4,1,9,8,1,6,8,4,6,8,1,9,8,1,3,4,8,2,2,2,3,7,9,6,8,4,5,9,1,7,8,1,7,6,10,2,2,8,5,1,7,6,10,7,6,9,1,7,9,1,3,7,8,4,7,1,6,7,5,9,7,1,6,10,9,1,1,1,10,8,5,10,3,7,10]],"Free":[[5,9,10,7,9,10,6,9,9,10,5,8,9,1,10,10,5,9,10,7,1,9,6,7,10,5,7,10,6,7,10,6,7,9,3,7,9,1,7,9,3,7,1,10,6,7,1,5,7,9,3,10,9,7,10,9,3,10,1,7,7,5,10,9,1,6,10,7,6,9,10,6,9,1,10,5,8,9,10,1,9,10,7,5,9,7,1,9,4,7,10,3,7,10],[4,8,10,3,8,10,1,7,3,8,7,1,8,10,4,8,3,10,7,1,8,9,2,8,10,4,8,1,10,8,3,7,9,4,8,8,1,7,8,2,5,7,7,10,2,7,10,1,5,10,8,4,10,8,1,4,10,5,1,10,8,5,10,7,5,10,8,5,1,10,8,4,10,1,9,8,4,10,8,4,10,10,8,5,10,1,5,10,8,6,3,10],[3,7,7,10,5,7,8,1,7,9,5,7,9,6,7,4,5,7,4,8,6,4,7,8,1,10,7,9,1,8,8,2,9,4,7,8,1,7,9,2,8,7,9,8,1,8,9,4,7,1,9,8,4,9,8,1,7,8,9,6,7,9,5,2,6,9,1,5,9,9,4,7,5,1,9,8,4,9,5,8,1,9,4,3,8,1,6,8,4,5,9,8],[5,9,8,4,9,8,1,6,8,4,9,1,4,8,9,1,8,4,1,9,8,4,3,8,4,3,8,6,9,8,6,9,4,10,9,3,8,9,2,2,2,7,9,4,10,8,1,3,9,1,6,5,9,8,4,1,6,8,1,6,9,1,4,9,6,3,10,6,1,7,6,9,5,7,6,9,1,7,6,1,7,9,3,8,10,1,9,8,5,3,10],[5,3,8,1,5,7,8,6,3,7,8,4,3,7,8,5,7,3,4,1,8,7,1,6,9,4,6,10,1,9,7,1,3,10,7,2,2,2,7,8,9,6,7,4,6,10,1,4,10,1,6,10,8,4,9,10,7,1,10,6,5,8,6,10,1,7,10,1,3,8,10,4,8,1,6,10,5,9,10,1,6,8,9,1,10,8,7,5,10,7,3,10,9]]},"Status":"N","Test":true,"UserID":"000000000","WagersID":"0","event":true,"isCash":true,"userSetting":{}},"event":true}}`
	return json.RawMessage(msg)
}

func (g FakeGame5145) GetMachineDetail(UserID) json.RawMessage {
	r := &response{
		Action: "onGetMachineDetail",
		Result: result{
			Event: true,
			Data: map[string]interface{}{
				"Balance":      99999999.99,
				"Base":         "1:10,1:5,1:2,1:1,2:1,5:1,10:1",
				"BetBase":      "",
				"BetTotal":     "1000000.0000",
				"Credit":       "0.00",
				"Currency":     "RMB",
				"DefaultBase":  "1:1",
				"ExchangeRate": "1",
				"HallID":       6,
				"LastBetTime":  "2019-09-27 00:36:55",
				"LoginName":    "fakeuser",
				"PayTotal":     "975000.0000",
				"Percent":      "99",
				"PercentLow":   "96",
				"Status":       "N",
				"Test":         true,
				"UserID":       "000000000",
				"WagersID":     "0",
				"event":        true,
			},
		},
	}

	bytes, _ := json.Marshal(r)
	return bytes
}

const beginGameResult1 = `{"action":"onBeginGame","result":{"data":{"AllPayTotal":0,"AxisLocation":"4-78-57-27-45","BBJackpot":{},"BetInfo":{},"BetValue":50,"Cards":["9-10-5","2-6-10","6-9-9","5-10-7","9-1-7"],"Credit":49950,"Credit_End":49950,"EncryID":null,"FreeGameBonusTime":{"ScatterBonus":[[0,0,0],[0,0,0],[0,0,0],[0,0,0],[0,0,0]]},"GamePayTotal":{"Lines":0},"LinePayoff":0,"Lines":[{"Element":["9","2","9"],"ElementID":9,"GridNum":3,"Grids":["1","4","8"],"Payoff":10},{"Element":["9","2","9"],"ElementID":9,"GridNum":3,"Grids":["1","4","9"],"Payoff":10}],"PayTotal":0,"PayValue":0,"RollerNumber":0,"WagersID":310631818546},"event":true}}`

const beginGameResult2 = `{"action":"onBeginGame","result":{"data":{"AllPayTotal":0,"AxisLocation":"1-2-3-4-5","BBJackpot":{},"BetInfo":{},"BetValue":50,"Cards":["1-2-3","2-3-4","3-4-5","4-5-6","5-6-7"],"Credit":49950,"Credit_End":49950,"EncryID":null,"FreeGameBonusTime":{"ScatterBonus":[[0,0,0],[0,0,0],[0,0,0],[0,0,0],[0,0,0]]},"GamePayTotal":{"Lines":0},"LinePayoff":0,"Lines":[],"PayTotal":0,"PayValue":0,"RollerNumber":0,"WagersID":310631818546},"event":true}}`

func (g FakeGame5145) BeginGame(SessionID, BetInfo) json.RawMessage {
	rand.Seed(time.Now().UnixNano())
	if rand.Intn(10)%2 == 0 {
		return json.RawMessage(beginGameResult1)
	} else {
		return json.RawMessage(beginGameResult2)
	}
}

func (g FakeGame5145) CreditExchange(SessionID, BetBase, int) json.RawMessage {
	msg := `{"action":"onCreditExchange","result":{"data":{"Balance":99949999.99,"BetBase":"1:1","Credit":100,"event":true},"event":true}}`
	return json.RawMessage(msg)
}

func (g FakeGame5145) BalanceExchange(UserID, HallID) json.RawMessage {
	msg := `{"action":"onBalanceExchange","result":{"data":{"Amount":50000,"Balance":9999999.99,"BetBase":"","ErrorID":1354000000,"TransCredit":"0.00","event":true},"event":true}}`
	return json.RawMessage(msg)
}

func (g FakeGame5145) LeaveMachine(UserID, HallID) json.RawMessage {
	panic("implement me")
}
