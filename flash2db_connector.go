package gode

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

//By convention, one-method interfaces are named by the method name plus an -er suffix ...
//https://golang.org/doc/effective_go.html#interface-names
type Connector interface {
	Connect(function string, parameters ...interface{}) (json.RawMessage, error)
}

type Flash2dbConnector struct {
	url string
}

func NewFlash2dbConnector(host string) *Flash2dbConnector {
	const prefixPath = `/amfphp/json.php`
	return &Flash2dbConnector{
		url: host + prefixPath,
	}
}

func (f *Flash2dbConnector) Connect(function string, parameters ...interface{}) (json.RawMessage, error) {
	msg, err := f.get(f.generateURL(function, parameters...))
	if err != nil {
		return nil, err
	}

	return msg, nil
}

func (f *Flash2dbConnector) get(url string) (json.RawMessage, error) {
	response, err := http.Get(url)
	if err != nil {
		log.Fatal("http Get Error", err)
		return nil, err
	}
	//todo: should we do anything?
	//noinspection GoUnhandledErrorResult
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("flash2db response code not OK, got %d", response.StatusCode)
	}

	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Print("ioutil ReadAll error : ", err)
		return nil, err
	}
	fmt.Printf("---\nrequest:%s\nresponse:%v\n", url, string(bytes))
	return bytes, nil
}

func (f *Flash2dbConnector) generateURL(phpFunctionName string, param ...interface{}) string {
	b := strings.Builder{}

	b.WriteString(fmt.Sprintf("%s/%s", f.url, phpFunctionName))

	for _, p := range param {
		b.WriteString(fmt.Sprintf("/%v", p))
	}

	return b.String()
}
