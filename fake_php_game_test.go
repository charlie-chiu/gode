package gode

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestFakePhpGame_OnReady(t *testing.T) {
	t.Run("should return valid json string", func(t *testing.T) {
		game := &FakePhpGame{}
		want := &response{
			Action: "onReady",
			Result: result{
				Event: true,
				//Data:  ???,
			},
		}
		got := &response{}
		err := json.Unmarshal([]byte(game.OnReady()), got)
		if err != nil {
			t.Fatalf("problem unmarshal json %v", err)
		}

		if !reflect.DeepEqual(want, got) {
			t.Errorf("want %v got %v", want, got)
		}
	})
}
