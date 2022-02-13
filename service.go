package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	lc "github.com/number571/go-peer/local"
	nt "github.com/number571/go-peer/network"
)

func main() {
	hlsDefaultInit()
	fmt.Println("Service is listening...")

	client := lc.NewClient(PrivKey, Settings)
	node := nt.NewNode(client).
		Handle([]byte(PatternHLS), routeHLS)

	for _, conn := range Config.Connects {
		err := node.Connect(conn)
		if err != nil {
			fmt.Println(err)
		}
	}

	if Config.Address == "" {
		select {}
	}

	err := node.Listen(Config.Address)
	if err != nil {
		fmt.Println(err)
	}
}

func routeHLS(client lc.Client, msg lc.Message) []byte {
	request := new(Request)

	deserialize(msg.Body.Data, request)
	if request == nil {
		return nil
	}

	addr, ok := Config.Services[request.Host]
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
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil
	}

	return data
}
