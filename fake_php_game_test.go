package gode

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestFakePhpGame(t *testing.T) {
	t.Run("onReady", func(t *testing.T) {
		game := &FakePhpGame{}
		want := &response{
			Action: "onReady",
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
		t.Errorf("want %#v got %#v", want, got)
	}
}

func unmarshalJSON(t *testing.T, str string, got *response) {
	t.Helper()
	err := json.Unmarshal([]byte(str), got)
	if err != nil {
		t.Fatalf("problem unmarshal json %v", err)
	}
}
