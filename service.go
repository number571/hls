package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	gp "github.com/number571/gopeer"
	"io/ioutil"
	"net/http"
	"time"
)

type Request struct {
	Host   string
	Path   string
	Method string
	Head   map[string]string
	Body   []byte
}

const (
	HLS_PROTO = "[HLS-PROTO]"
)

var (
	USER     = new(User)
	DATABASE = NewDB("service.db")
	FLCONFIG = NewCFG("service.cfg")
)

func main() {
	fmt.Println("Service is listening...")
	client := gp.NewClient(
		USER.Priv,
		handleFunc,
	)
	for _, addr := range FLCONFIG.Connects {
		client.Connect(addr)
	}
	if FLCONFIG.Openaddr == "" {
		for {
			time.Sleep(1 * time.Hour)
		}
	}
	err := client.RunNode(FLCONFIG.Openaddr)
	if err != nil {
		panic(err)
	}
}

func handleFunc(client *gp.Client, pack *gp.Package) {
	client.Handle(HLS_PROTO, pack, handleService)
}

func handleService(client *gp.Client, pack *gp.Package) []byte {
	request := deserialize(pack.Body.Data)
	if request == nil {
		return nil
	}
	addr, ok := FLCONFIG.Services[request.Host]
	if !ok {
		return nil
	}
	req, err := http.NewRequest(
		request.Method,
		addr+request.Path,
		bytes.NewReader(request.Body),
	)
	if err != nil {
		return nil
	}
	for key, val := range request.Head {
		req.Header.Add(key, val)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil
	}
	resp.Body.Close()
	return data
}

func serialize(data interface{}) []byte {
	res, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return nil
	}
	return res
}

func deserialize(data []byte) *Request {
	var req = new(Request)
	err := json.Unmarshal(data, req)
	if err != nil {
		return nil
	}
	return req
}
