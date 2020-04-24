package gode

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

// this implementation fetch user data from flash2db
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

func (c *Flash2dbClient) Fetch() {
	url := c.host
	rawMsg := c.call(url)

	loginCheck := &LoginCheck{}
	err := json.Unmarshal(rawMsg, loginCheck)
	if err != nil {
		log.Panicf("JSON unmarshal error, %v", err)
		return
	}

	c.hid = toHallID(loginCheck.Data.HallID)
	c.uid = UserID(loginCheck.Data.UserID)
	c.sid = SessionID(loginCheck.Data.Sid)
}

func (c *Flash2dbClient) call(url string) json.RawMessage {
	response, err := http.Get(url)
	if err != nil {
		log.Fatal("http Get Error", err)
	}
	//todo: should we do anything?
	//noinspection GoUnhandledErrorResult
	defer response.Body.Close()

	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Print("ioutil ReadAll error : ", err)
	}
	return bytes
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
