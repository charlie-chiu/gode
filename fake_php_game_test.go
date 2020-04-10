package gode

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestFakePhpGame_OnReady(t *testing.T) {
	t.Run("should return valid JSON", func(t *testing.T) {
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
}

func assertJSONEqual(t *testing.T, want *response, got *response) {
	t.Helper()
	if !reflect.DeepEqual(want, got) {
		t.Errorf("want %v got %v", want, got)
	}
}

func unmarshalJSON(t *testing.T, str string, got *response) {
	t.Helper()
	err := json.Unmarshal([]byte(str), got)
	if err != nil {
		t.Fatalf("problem unmarshal json %v", err)
	}
}
