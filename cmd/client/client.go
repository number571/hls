package main

import (
	"fmt"
	"os"

	cr "github.com/number571/go-peer/crypto"
	lc "github.com/number571/go-peer/local"
	nt "github.com/number571/go-peer/network"

	st "github.com/number571/hls/settings"
)

func main() {
	priv := cr.NewPrivKey(st.AKEYSIZE)
	client := lc.NewClient(priv, st.SETTINGS)

	node := nt.NewNode(client).
		Handle([]byte(st.HLS), nil)

	err := node.Connect("localhost:9571")
	if err != nil {
		fmt.Println("error: connection")
		os.Exit(1)
	}

	msg := lc.NewMessage(
		[]byte(st.HLS),
		st.Serialize(&st.Request{
			Host:   st.ServerAddressInHLS,
			Path:   "/echo",
			Method: "GET",
			Head: map[string]string{
				"Content-Type": "application/json",
			},
			Body: []byte(`{"message": "hello, world!"}`),
		}),
	)

	spub := string(st.ReadFile(fmt.Sprintf("../service/%s", st.FileWithPubKey)))
	route := lc.NewRoute(cr.LoadPubKeyByString(spub), nil, nil)

	res, err := node.Broadcast(route, msg)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(res))
}
