package gode

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

func TestFakePhpGame(t *testing.T) {
	t.Run("onReady", func(t *testing.T) {
		game := &FakePhpGame{}
		want := &response{
			Action: "ready",
			Result: result{
				Event: true,
				//Data:  ???,
			},
		}
		got := &response{}

		unmarshalJSON(t, game.OnReady(), got)
		assertJSONEqual(t, want, got)
	})

	t.Run("onLogin", func(t *testing.T) {
		game := &FakePhpGame{}
		want := &response{
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
		got := &response{}

		unmarshalJSON(t, game.OnLogin(), got)
		assertJSONEqual(t, want, got)
	})

	t.Run("onTakeMachine", func(t *testing.T) {
		game := &FakePhpGame{}
		want := &response{
			Action: "onTakeMachine",
			Result: result{
				Event: true,
				Data: map[string]interface{}{
					"GameCode": 1,
					"event":    true,
					"gameCode": 1,
				},
			},
		}
		got := &response{}

		unmarshalJSON(t, game.OnTakeMachine(), got)
		assertJSONEqual(t, want, got)
	})

	t.Run("onLoadInfo", func(t *testing.T) {
		game := &FakePhpGame{}
		want := &response{
			Action: "onOnLoadInfo2",
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
					"Roller": map[string]interface{}{
						"Normal": []interface{}{
							//[]interface{}{5, 9, 10, 6, 9, 10, 5, 9, 9, 10, 6, 9, 1, 7, 9, 9, 5, 10, 10, 6, 1, 9, 6, 10, 9, 5, 8, 9, 6, 10, 10, 6, 9, 9, 3, 8, 8, 1, 9, 9, 3, 8, 1, 10, 5, 8, 1, 5, 8, 8, 3, 10, 10, 6, 8, 10, 3, 8, 1, 9, 9, 5, 10, 7, 1, 6, 9, 10, 6, 9, 10, 4, 9, 1, 10, 5, 8, 8, 9, 1, 8, 10, 10, 5, 9, 10, 1, 9, 4, 10, 10, 3, 9, 8},
							//[]interface{}{5, 10, 10, 5, 7, 7, 1, 8, 4, 1, 8, 8, 10, 2, 7, 10, 2, 7, 4, 1, 7, 7, 5, 8, 2, 4, 7, 1, 10, 10, 4, 8, 8, 4, 10, 9, 5, 8, 8, 4, 10, 7, 7, 1, 10, 10, 1, 6, 7, 3, 8, 2, 5, 10, 1, 4, 8, 5, 1, 7, 7, 5, 3, 10, 5, 2, 7, 8, 1, 10, 10, 4, 9, 1, 10, 10, 4, 8, 2, 6, 10, 10, 7, 5, 10, 1, 5, 7, 7, 4, 3, 10},
							//[]interface{}{4, 7, 7, 8, 4, 7, 7, 1, 9, 4, 7, 7, 9, 5, 7, 4, 5, 7, 3, 8, 2, 4, 10, 1, 3, 8, 8, 4, 1, 7, 7, 5, 2, 4, 9, 8, 6, 1, 8, 8, 2, 9, 6, 7, 1, 5, 9, 4, 8, 1, 9, 9, 3, 7, 7, 1, 8, 6, 9, 9, 3, 2, 6, 3, 7, 7, 1, 6, 9, 5, 4, 9, 9, 1, 7, 7, 2, 8, 6, 9, 1, 7, 4, 6, 1, 8, 8, 6, 4, 10, 2, 3},
							//[]interface{}{4, 9, 10, 4, 9, 7, 1, 5, 7, 4, 9, 1, 4, 7, 10, 1, 7, 6, 1, 8, 10, 4, 3, 10, 6, 4, 7, 5, 10, 7, 4, 9, 3, 4, 9, 10, 7, 9, 2, 2, 2, 7, 9, 10, 6, 7, 1, 3, 8, 1, 6, 5, 9, 10, 4, 1, 5, 7, 1, 5, 10, 1, 9, 10, 6, 7, 9, 6, 1, 7, 6, 9, 5, 7, 6, 10, 1, 6, 7, 1, 3, 9, 2, 2, 4, 10, 1, 8, 5, 3, 10},
							//[]interface{}{5, 3, 7, 1, 5, 8, 7, 3, 4, 8, 9, 4, 3, 7, 9, 5, 7, 8, 4, 1, 9, 8, 1, 6, 8, 4, 6, 8, 1, 9, 8, 1, 3, 4, 8, 2, 2, 2, 3, 7, 9, 6, 8, 4, 5, 9, 1, 7, 8, 1, 7, 6, 10, 2, 2, 8, 5, 1, 7, 6, 10, 7, 6, 9, 1, 7, 9, 1, 3, 7, 8, 4, 7, 1, 6, 7, 5, 9, 7, 1, 6, 10, 9, 1, 1, 1, 10, 8, 5, 10, 3, 7, 10},
						},
						"Free": []interface{}{
							//[]interface{}{5, 9, 10, 7, 9, 10, 6, 9, 9, 10, 5, 8, 9, 1, 10, 10, 5, 9, 10, 7, 1, 9, 6, 7, 10, 5, 7, 10, 6, 7, 10, 6, 7, 9, 3, 7, 9, 1, 7, 9, 3, 7, 1, 10, 6, 7, 1, 5, 7, 9, 3, 10, 9, 7, 10, 9, 3, 10, 1, 7, 7, 5, 10, 9, 1, 6, 10, 7, 6, 9, 10, 6, 9, 1, 10, 5, 8, 9, 10, 1, 9, 10, 7, 5, 9, 7, 1, 9, 4, 7, 10, 3, 7, 10},
							//[]interface{}{4, 8, 10, 3, 8, 10, 1, 7, 3, 8, 7, 1, 8, 10, 4, 8, 3, 10, 7, 1, 8, 9, 2, 8, 10, 4, 8, 1, 10, 8, 3, 7, 9, 4, 8, 8, 1, 7, 8, 2, 5, 7, 7, 10, 2, 7, 10, 1, 5, 10, 8, 4, 10, 8, 1, 4, 10, 5, 1, 10, 8, 5, 10, 7, 5, 10, 8, 5, 1, 10, 8, 4, 10, 1, 9, 8, 4, 10, 8, 4, 10, 10, 8, 5, 10, 1, 5, 10, 8, 6, 3, 10},
							//[]interface{}{3, 7, 7, 10, 5, 7, 8, 1, 7, 9, 5, 7, 9, 6, 7, 4, 5, 7, 4, 8, 6, 4, 7, 8, 1, 10, 7, 9, 1, 8, 8, 2, 9, 4, 7, 8, 1, 7, 9, 2, 8, 7, 9, 8, 1, 8, 9, 4, 7, 1, 9, 8, 4, 9, 8, 1, 7, 8, 9, 6, 7, 9, 5, 2, 6, 9, 1, 5, 9, 9, 4, 7, 5, 1, 9, 8, 4, 9, 5, 8, 1, 9, 4, 3, 8, 1, 6, 8, 4, 5, 9, 8},
							//[]interface{}{5, 9, 8, 4, 9, 8, 1, 6, 8, 4, 9, 1, 4, 8, 9, 1, 8, 4, 1, 9, 8, 4, 3, 8, 4, 3, 8, 6, 9, 8, 6, 9, 4, 10, 9, 3, 8, 9, 2, 2, 2, 7, 9, 4, 10, 8, 1, 3, 9, 1, 6, 5, 9, 8, 4, 1, 6, 8, 1, 6, 9, 1, 4, 9, 6, 3, 10, 6, 1, 7, 6, 9, 5, 7, 6, 9, 1, 7, 6, 1, 7, 9, 3, 8, 10, 1, 9, 8, 5, 3, 10},
							//[]interface{}{5, 3, 8, 1, 5, 7, 8, 6, 3, 7, 8, 4, 3, 7, 8, 5, 7, 3, 4, 1, 8, 7, 1, 6, 9, 4, 6, 10, 1, 9, 7, 1, 3, 10, 7, 2, 2, 2, 7, 8, 9, 6, 7, 4, 6, 10, 1, 4, 10, 1, 6, 10, 8, 4, 9, 10, 7, 1, 10, 6, 5, 8, 6, 10, 1, 7, 10, 1, 3, 8, 10, 4, 8, 1, 6, 10, 5, 9, 10, 1, 6, 8, 9, 1, 10, 8, 7, 5, 10, 7, 3, 10, 9},
						},
					},
					"Status":      "N",
					"Test":        true,
					"UserID":      "000000000",
					"WagersID":    "0",
					"event":       true,
					"isCash":      true,
					"userSetting": map[string]interface{}{},
				},
			},
		}
		got := &response{}

		unmarshalJSON(t, game.OnLoadInfo(), got)
		assertJSONEqual(t, want, got)
	})

	t.Run("onGetMachineDetail", func(t *testing.T) {
		game := &FakePhpGame{}
		want := &response{
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
		got := &response{}

		unmarshalJSON(t, game.OnGetMachineDetail(), got)
		assertJSONEqual(t, want, got)
	})
}

func TestFakePhpGame_BeginGame(t *testing.T) {
	//todo: refactor this test ...maybe
	g := &FakePhpGame{}
	t.Run("got result 1", func(t *testing.T) {
		got := g.BeginGame()
		if got != beginGameResult1 && got != beginGameResult2 {
			t.Errorf("begin game result not matched")
		}
	})
}

func assertJSONEqual(t *testing.T, want *response, got *response) {
	// json.unmarshal will store float64 for JSON numbers
	// see https://golang.org/pkg/encoding/json/#Unmarshal for more
	for k, v := range want.Result.Data {
		if reflect.TypeOf(v).Kind() == reflect.Int {
			want.Result.Data[k] = float64(v.(int))
		}
	}

	t.Helper()
	if !reflect.DeepEqual(want, got) {
		fmt.Printf("want: %#v \n", want)
		fmt.Printf(" got: %#v \n", got)
		t.Errorf("not matched")
	}
}

func unmarshalJSON(t *testing.T, str string, got *response) {
	t.Helper()
	err := json.Unmarshal([]byte(str), got)
	if err != nil {
		t.Fatalf("problem unmarshal json %v", err)
	}
}
