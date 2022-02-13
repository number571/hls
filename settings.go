package main

import (
	cr "github.com/number571/go-peer/crypto"
	st "github.com/number571/go-peer/settings"
)

const (
	AkeySize   = 2048
	PatternHLS = "hidden-lake-service"
)

var (
	Settings st.Settings
	PrivKey  cr.PrivKey
	Config   *CFG
)
