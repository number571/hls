package main

import (
	"fmt"
	"flag"
	"os"
	"crypto/rsa"
	"strings"
	gp "github.com/number571/gopeer"
)

func init() {
	gp.Set(gp.SettingsType{
		"POWS_DIFF": uint(20),
		"PACK_SIZE": uint(8 << 20),
		"AKEY_SIZE": uint(2 << 10),
		"SKEY_SIZE": uint(1 << 5),
		"RAND_SIZE": uint(1 << 4),
	})
	signupPtr := flag.String("signup", "", "create private key by nickname@password")
	signinPtr := flag.String("signin", "", "unlock private key by nickname@password")
	prvkeyPtr := flag.String("prvkey", "", "upload private key with signup")
	viewPtr   := flag.String("view", "", 
		"print user information: \n\t1. pubkey;\n\t2. prvkey;\n\t3. all;")
	flag.Parse()
	if *signupPtr != "" {
		slice := strings.Split(*signupPtr, "@")
		if len(slice) < 2 {
			panic("signup failed")
		}
		name := slice[0]
		pasw := strings.Join(slice[1:], "@")
		priv := new(rsa.PrivateKey)
		if *prvkeyPtr == "" {
			priv = gp.GenerateKey(gp.Get("AKEY_SIZE").(uint))
		} else {
			priv = gp.StringToPrivateKey(*prvkeyPtr)
		}
		err := DATABASE.SetUser(name, pasw, priv)
		if err != nil {
			panic(err.Error())
		}
	}
	if *signinPtr == "" {
		if *signupPtr != "" {
			os.Exit(0)
		}
		panic("signin undefined")
	}
	slice := strings.Split(*signinPtr, "@")
	if len(slice) < 2 {
		panic("signin failed")
	}
	name := slice[0]
	pasw := strings.Join(slice[1:], "@")
	USER = DATABASE.GetUser(name, pasw)
	if USER == nil {
		panic("user is null")
	}
	switch *viewPtr {
	case "":
	case "pubkey":
		fmt.Printf("Public key:\n%s\n\n", gp.PublicKeyToString(&(USER.Priv).PublicKey))
	case "prvkey":
		fmt.Printf("Private key:\n%s\n\n", gp.PrivateKeyToString(USER.Priv))
	case "all":
		fmt.Printf("Public key:\n%s\n\n", gp.PublicKeyToString(&(USER.Priv).PublicKey))
		fmt.Printf("Private key:\n%s\n\n", gp.PrivateKeyToString(USER.Priv))
	default:
		fmt.Println("view: undefined arg\n")
	}
}
