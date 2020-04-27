package gode

import (
	"encoding/json"
	"log"
	"strconv"
)

// this implementation fetch user data from flash2db
type Flash2dbClient struct {
	uid  UserID
	hid  HallID
	sid  SessionID
	conn Connector
}

func NewFlash2dbClient(connector Connector) *Flash2dbClient {
	return &Flash2dbClient{
		conn: connector,
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

// result of flash2db client.loginCheck
type LoginCheck struct {
	Data struct {
		UserID       int    `json:"UserID"`
		Sid          string `json:"Sid"`
		HallID       string `json:"HallID"`
		GameID       string `json:"GameID"`
		COID         string `json:"COID"`
		Test         string `json:"Test"`
		ExchangeRate string `json:"ExchangeRate"`
		IP           string `json:"IP"`
	} `json:"data"`
	Event bool `json:"event"`
}

func (c *Flash2dbClient) Login(sid SessionID) json.RawMessage {
	const function = "Client.loginCheck"
	ip := "127.0.0.1"
	rawMsg, _ := c.conn.Connect(function, sid, ip)

	loginCheck := &LoginCheck{}
	err := json.Unmarshal(rawMsg, loginCheck)
	if err != nil {
		log.Panicf("JSON unmarshal error, %v", err)
	}

	c.hid = toHallID(loginCheck.Data.HallID)
	c.uid = UserID(loginCheck.Data.UserID)
	c.sid = SessionID(loginCheck.Data.Sid)

	return rawMsg
}

func toHallID(input interface{}) (hid HallID) {
	switch input.(type) {
	case HallID:
		return
	case int:
		return HallID(input.(int))
	case string:
		i, _ := strconv.Atoi(input.(string))
		return HallID(i)
	}

	return 0
}
