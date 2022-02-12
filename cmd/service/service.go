package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	cr "github.com/number571/go-peer/crypto"
	lc "github.com/number571/go-peer/local"
	nt "github.com/number571/go-peer/network"

	st "github.com/number571/hls/settings"
)

type CFG struct {
	Address  string            `json:"address"`
	Services map[string]string `json:"services"`
	Connects []string          `json:"connects"`
}

var (
	PrivKey cr.PrivKey
	Config  *CFG
)

const (
	FileWithPrivKey = "priv.key"
	FileWithConfig  = "hls.cfg"
)

func init() {
	var initOnly bool

	flag.BoolVar(&initOnly, "init-only", false, "run initialization only")
	flag.Parse()

	if !st.FileIsExist(FileWithPrivKey) {
		priv := cr.NewPrivKey(st.AKEYSIZE)
		st.WriteFile(FileWithPrivKey, []byte(priv.String()))
		st.WriteFile(st.FileWithPubKey, []byte(priv.PubKey().String()))
	}
	spriv := string(st.ReadFile(FileWithPrivKey))
	PrivKey = cr.LoadPrivKeyByString(spriv)

	if !st.FileIsExist(FileWithConfig) {
		config := &CFG{
			Address:  "localhost:9571",
			Connects: []string{"127.0.0.2:9571"},
			Services: map[string]string{
				st.ServerAddressInHLS: "http://localhost:8080",
			},
		}
		st.WriteFile(FileWithConfig, st.Serialize(config))
	}
	st.Deserialize(st.ReadFile(FileWithConfig), &Config)

	if initOnly {
		os.Exit(0)
	}
}

func main() {
	fmt.Println("Service is listening...")

	client := lc.NewClient(PrivKey, st.SETTINGS)
	node := nt.NewNode(client).
		Handle([]byte(st.HLS), hlservice)

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

func hlservice(client lc.Client, msg lc.Message) []byte {
	request := new(st.Request)

	st.Deserialize(msg.Body.Data, request)
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
