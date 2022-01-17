package main

import (
	"fmt"
	"os"

	cr "github.com/number571/go-peer/crypto"
	lc "github.com/number571/go-peer/local"
	nt "github.com/number571/go-peer/network"
	gp "github.com/number571/go-peer/settings"
)

func main() {
	priv := cr.NewPrivKey(gp.Get("AKEY_SIZE").(uint))

	client := lc.NewClient(priv)
	node := nt.NewNode(client).
		Handle([]byte(HLS), nil)

	err := node.Connect("localhost:9571")
	if err != nil {
		fmt.Println("error: connection")
		os.Exit(1)
	}

	msg := lc.NewMessage(
		[]byte(HLS),
		serialize(&Request{
			Host:   ServerAddressInHLS,
			Path:   "/echo",
			Method: "GET",
			Head: map[string]string{
				"Content-Type": "application/json",
			},
			Body: []byte(`{"message": "hello, world!"}`),
		}),
		gp.Get("POWS_DIFF").(uint),
	)

	spub := string(readFile(FileWithPubKey))
	route := lc.NewRoute(cr.LoadPubKeyByString(spub))

	res, err := node.Send(msg, route)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(res))
}
