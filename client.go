package gode

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
	Fetch()
}

type FakeClient struct {
	UID UserID
	HID HallID
	SID SessionID
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
