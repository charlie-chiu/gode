package gode

// this client implement get user date from flash2db
type Flash2dbClient struct {
	host string
	uid  UserID
	hid  HallID
	sid  SessionID
}

func NewFlash2dbClient(host string) *Flash2dbClient {
	return &Flash2dbClient{
		host: host,
	}
}

func (c *Flash2dbClient) UserID() UserID {
	return c.uid
}

func (c *Flash2dbClient) HallID() HallID {
	return c.hid
}

func (c *Flash2dbClient) SessionID() SessionID {
	return c.sid
}
