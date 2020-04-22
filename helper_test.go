package gode_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
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

func assertURLEqual(t *testing.T, r *http.Request, want string) {
	t.Helper()
	if r.URL.Path != want {
		t.Errorf("URL not matched\n want %q\n, got %q", want, r.URL)
	}
}
