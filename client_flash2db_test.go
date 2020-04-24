package gode_test

import (
	"testing"

	"github.com/charlie-chiu/gode"
)

func TestFlash2dbClientBeforeFetchInformation(t *testing.T) {
	client := gode.NewFlash2dbClient("")

	t.Run("UserID() returns zero value", func(t *testing.T) {
		want := gode.UserID(0)
		got := client.UserID()

		assertUserIDEqual(t, got, want)
	})

	t.Run("HallID returns zero value", func(t *testing.T) {
		want := gode.HallID(0)
		got := client.HallID()

		assertHallIDEqual(t, got, want)
	})

	t.Run("SessionID returns zero value", func(t *testing.T) {
		want := gode.SessionID("")
		got := client.SessionID()

		assertSessionIDEqual(t, got, want)
	})
}
