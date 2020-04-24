package gode

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
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

// result of flash2db client.getUserInfo
type UserInfo struct {
	Event bool `json:"event"`
	Data  struct {
		Sid      string `json:"Sid"`
		UserID   int    `json:"UserID"`
		UserName string `json:"UserName"`
		HallID   int    `json:"HallID"`
		Currency string `json:"Currency"`
		Credit   int    `json:"Credit"`
	} `json:"data"`
}

func (c *Flash2dbClient) Fetch() {
	url := c.host
	rawMsg := c.call(url)

	//rawMsg := json.RawMessage(`{"event":true,"data":{"Sid":"197af9c6341e4f846d6defe4da1aaf0489dc15d5","UserID":362907402,"UserName":"angel888","HallID":32,"Currency":"RMB","Credit":19872311}}`)
	userInfo := &UserInfo{}
	err := json.Unmarshal(rawMsg, userInfo)
	if err != nil {
		log.Panicf("JSON unmarshal error, %v", err)
		return
	}

	c.hid = HallID(userInfo.Data.HallID)
	c.uid = UserID(userInfo.Data.UserID)
	c.sid = SessionID(userInfo.Data.Sid)
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
