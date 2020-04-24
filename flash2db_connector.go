package gode

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type Connector interface {
	Connect(function string, parameters ...interface{}) json.RawMessage
}

type Flash2db struct {
	url string
}

func NewFlash2db(host string) *Flash2db {
	const prefixPath = `/amfphp/json.php`
	return &Flash2db{
		url: host + prefixPath,
	}
}

func (f *Flash2db) Connect(function string, parameters ...interface{}) json.RawMessage {
	f.get(f.generateURL(function, parameters...))

	return json.RawMessage(``)
}

func (f *Flash2db) get(url string) json.RawMessage {
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

func (f *Flash2db) generateURL(phpFunctionName string, param ...interface{}) string {
	b := strings.Builder{}

	b.WriteString(fmt.Sprintf("%s/%s", f.url, phpFunctionName))

	for _, p := range param {
		b.WriteString(fmt.Sprintf("/%v", p))
	}

	return b.String()
}
