package gode_test

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/charlie-chiu/gode"
)

func TestFlash2dbClientBeforeFetchInformation(t *testing.T) {
	client := gode.NewFlash2dbClient(gode.NewFlash2dbConnector(""))

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

type SpyConnector struct {
	funcCalled []funcCall
}

type funcCall struct {
	funcName string
	args     []interface{}
}

func (c *SpyConnector) Connect(function string, parameters ...interface{}) (json.RawMessage, error) {
	c.funcCalled = append(c.funcCalled, funcCall{
		funcName: function,
		args:     parameters,
	})

	return json.RawMessage(`{}`), nil
}

type MockConnector struct {
	returnMsg json.RawMessage
	returnErr error
}

func (c *MockConnector) Connect(string, ...interface{}) (json.RawMessage, error) {
	return c.returnMsg, c.returnErr
}

func TestFlash2dbClient_Login(t *testing.T) {
	const LoginFunction = "Client.loginCheck"

	t.Run("connect with correct args", func(t *testing.T) {
		spyConnector := &SpyConnector{}
		client := gode.NewFlash2dbClient(spyConnector)

		sid := gode.SessionID("19870604xi")
		ip := "127.0.0.1"
		client.Login(sid, ip)

		expectedCalls := []funcCall{
			{LoginFunction, []interface{}{sid, ip}},
		}

		assertFuncCalledSame(t, expectedCalls, spyConnector.funcCalled)
	})

	t.Run("store updated sid, uid and hid after successful login", func(t *testing.T) {
		sid := gode.SessionID("19870604xi")
		uid := gode.UserID(362907402)
		hid := gode.HallID(32)
		msg := fmt.Sprintf(`{"data":{"UserID":%d,"Sid":"%s","HallID":"%d"},"event":true}`, uid, sid, hid)

		client := gode.NewFlash2dbClient(&MockConnector{
			returnMsg: json.RawMessage(msg),
		})

		client.Login("dummySID", "dummyIP")

		assertUserIDEqual(t, client.UserID(), uid)
		assertHallIDEqual(t, client.HallID(), hid)
		assertSessionIDEqual(t, client.SessionID(), sid)
	})

	t.Run("login return msg got from flash2db", func(t *testing.T) {
		msg := fmt.Sprintf(`{"data":{"UserID":123,"Sid":"dummySID","HallID":"123"},"event":true}`)

		client := gode.NewFlash2dbClient(&MockConnector{
			returnMsg: json.RawMessage(msg),
		})

		got := client.Login("dummySID", "dummyIP")

		assertRawJSONEqual(t, got, json.RawMessage(msg))
	})
}

func assertFuncCalledSame(t *testing.T, expectedCalls, got []funcCall) {
	t.Helper()
	if !reflect.DeepEqual(expectedCalls, got) {
		fmt.Printf("expected:%v\n", expectedCalls)
		fmt.Printf("     got:%v\n", got)
		t.Errorf("called function not match")
	}
}
