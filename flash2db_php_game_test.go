package gode

import (
	"testing"
)

func TestFlash2dbPhpGame(t *testing.T) {
	game := NewFlash2dbPhpGame(5145)

	t.Run("OnLogin", func(t *testing.T) {
		var sid SessionID = "21d9b36e42c8275a4359f6815b859df05ec2bb0a"
		want := &response{
			Action: "onLogin",
			Result: result{
				Event: true,
				Data: map[string]interface{}{
					"COID":         2688,
					"ExchangeRate": 1,
					"GameID":       0,
					"HallID":       6,
					"Sid":          "21d9b36e42c8275a4359f6815b859df05ec2bb0a",
					"Test":         1,
					"UserID":       0,
				},
			},
		}
		got := &response{}

		unmarshalJSON(t, game.OnLogin(sid), got)
		assertJSONEqual(t, want, got)
	})
}
