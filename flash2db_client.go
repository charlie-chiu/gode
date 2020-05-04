package gode

import (
	"encoding/json"
	"fmt"
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

func (c *Flash2dbClient) Login(sid SessionID, ip string) json.RawMessage {
	const function = "Client.loginCheck"
	rawMsg, _ := c.conn.Connect(function, sid, ip)

	_ = c.storeLoginCheckResult(rawMsg)

	return rawMsg
}

// result of flash2db client.loginCheck
type LoginCheck struct {
	Data struct {
		UserID int    `json:"UserID"`
		Sid    string `json:"Sid"`
		//HallID       string `json:"HallID"`
		GameID       string `json:"GameID"`
		COID         string `json:"COID"`
		Test         string `json:"Test"`
		ExchangeRate string `json:"ExchangeRate"`
		IP           string `json:"IP"`
	} `json:"data"`
	Event bool `json:"event"`
}

func (c *Flash2dbClient) storeLoginCheckResult(rawMsg json.RawMessage) error {
	// store
	okData := &LoginCheck{}
	err := json.Unmarshal(rawMsg, okData)
	if err != nil {
		return fmt.Errorf("login check parsing error, %v", err)
	}
	c.uid = UserID(okData.Data.UserID)
	c.sid = SessionID(okData.Data.Sid)

	//store HallID : number or string

	//fmt.Printf("(%T)%+v\n", unknownType.Data["HallID"], unknownType.Data["HallID"])
	unknownType := struct {
		Data  map[string]interface{} `json:"data"`
		Event bool                   `json:"event"`
	}{}
	err = json.Unmarshal(rawMsg, &unknownType)
	if err != nil {
		return fmt.Errorf("login check parsing error, %v", err)
	}
	c.hid = toHallID(unknownType.Data["HallID"])

	return nil
}

func toHallID(input interface{}) (hid HallID) {
	switch input.(type) {
	case HallID:
		return
	case int:
		return HallID(input.(int))
	case float64:
		return HallID(int(input.(float64)))
	case string:
		i, _ := strconv.Atoi(input.(string))
		return HallID(i)
	}

	return 0
}
