package gode_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/charlie-chiu/gode"
)

func assertRawJSONEqual(t *testing.T, got, want json.RawMessage) {
	t.Helper()
	if bytes.Compare(got, want) != 0 {
		fmt.Printf("want: %s \n", want)
		fmt.Printf(" got: %s \n", got)
		t.Errorf("not matched")
	}
}

func assertError(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Errorf("expect an error but not got one")
	}
}

func assertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal("didn't expecting an error but got one", err)
	}
}

func assertURLEqual(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("URL not matched\n want %q\n, got %q", want, got)
	}
}

func assertSessionIDEqual(t *testing.T, got, want gode.SessionID) {
	t.Helper()
	if got != want {
		t.Errorf("wanted SessionID %s, got %s", want, got)
	}
}

func assertHallIDEqual(t *testing.T, got, want gode.HallID) {
	t.Helper()
	if got != want {
		t.Errorf("wanted HallID %d, got %d", want, got)
	}
}

func assertUserIDEqual(t *testing.T, got, want gode.UserID) {
	t.Helper()
	if got != want {
		t.Errorf("wanted UserID %d, got %d", want, got)
	}
}
