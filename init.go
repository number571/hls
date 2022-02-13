package main

import (
	"flag"
	"os"

	cr "github.com/number571/go-peer/crypto"
	st "github.com/number571/go-peer/settings"
)

func hlsDefaultInit() {
	var initOnly bool

	flag.BoolVar(&initOnly, "init-only", false, "run initialization only")
	flag.Parse()

	Settings = st.NewSettings()
	PrivKey = newPrivKey("priv.key")
	Config = newConfig("hls.cfg")

	if initOnly {
		os.Exit(0)
	}
}

func newPrivKey(filepath string) cr.PrivKey {
	var priv cr.PrivKey

	if !fileIsExist(filepath) {
		priv = cr.NewPrivKey(AkeySize)
		writeFile(filepath, []byte(priv.String()))
		writeFile("pub.key", []byte(priv.PubKey().String()))
	} else {
		spriv := string(readFile(filepath))
		priv = cr.LoadPrivKeyByString(spriv)
	}

	return priv
}

func newConfig(filepath string) *CFG {
	var cfg = new(CFG)

	if !fileIsExist(filepath) {
		cfg = &CFG{
			Address:  "localhost:9571",           // create local hls
			Connects: []string{"127.0.0.2:9571"}, // connect to another hls's
			Services: map[string]string{
				// crypto-address -> network-address
				"hidden-default-service": "http://localhost:8080",
			},
		}
		writeFile(filepath, serialize(cfg))
	} else {
		deserialize(readFile(filepath), cfg)
		if cfg == nil {
			panic("error: config is nil")
		}
	}

	return cfg
}

func fileIsExist(file string) bool {
	_, err := os.Stat(file)
	return !os.IsNotExist(err)
}
