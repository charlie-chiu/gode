package gode

import "encoding/json"

type (
	SessionID string
	UserID    int
	HallID    int
)

type Client interface {
	UserID() UserID
	HallID() HallID
	SessionID() SessionID

	// fetch user information from source(ACC)
	Login(sid SessionID) json.RawMessage
}

type FakeClient struct {
	UID UserID
	HID HallID
	SID SessionID
}

func (c FakeClient) Login(SessionID) json.RawMessage {
	return json.RawMessage(`{"event":true,"data":{"COID":2688,"ExchangeRate":1,"GameID":0,"HallID":6,"Sid":"","Test":1,"UserID":0}}`)
}

func (c FakeClient) UserID() UserID {
	return c.UID
}
func (c FakeClient) HallID() HallID {
	return c.HID
}
func (c FakeClient) SessionID() SessionID {
	return c.SID
}
