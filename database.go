package main

import (
	"bytes"
	"crypto/rsa"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	gp "github.com/number571/gopeer"
	"strings"
)

const (
	PASWDIFF = 25
)

func NewDB(filename string) *DB {
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		return nil
	}
	_, err = db.Exec(`
CREATE TABLE IF NOT EXISTS users (
	id   INTEGER,
	hashn VARCHAR(255) UNIQUE,
	hashp VARCHAR(255),
	salt VARCHAR(255),
	priv TEXT,
	PRIMARY KEY(id)
);
`)
	if err != nil {
		return nil
	}
	return &DB{
		ptr: db,
	}
}

func (db *DB) SetUser(name, pasw string, priv *rsa.PrivateKey) error {
	db.mtx.Lock()
	defer db.mtx.Unlock()
	if priv == nil {
		return fmt.Errorf("private key is null")
	}
	name = strings.TrimSpace(name)
	if len(name) < 6 || len(name) > 64 {
		return fmt.Errorf("need len username >= 6 and <= 64")
	}
	if len(pasw) < 8 {
		return fmt.Errorf("need len password >= 8")
	}
	if db.userExist(name) {
		return fmt.Errorf("user already exist")
	}
	salt := gp.GenerateBytes(32)
	bpasw := gp.RaiseEntropy([]byte(pasw), salt, PASWDIFF)
	hpasw := gp.HashSum(bytes.Join(
		[][]byte{
			bpasw,
			[]byte(name),
		},
		[]byte{},
	))
	_, err := db.ptr.Exec(
		"INSERT INTO users (hashn, hashp, salt, priv) VALUES ($1, $2, $3, $4)",
		gp.Base64Encode(gp.HashSum([]byte(name))),
		gp.Base64Encode(hpasw),
		gp.Base64Encode(salt),
		gp.Base64Encode(gp.EncryptAES(bpasw, gp.PrivateKeyToBytes(priv))),
	)
	return err
}

func (db *DB) GetUser(name, pasw string) *User {
	db.mtx.Lock()
	defer db.mtx.Unlock()
	var (
		id    int
		hpasw string
		ssalt string
		spriv string
	)
	name = strings.TrimSpace(name)
	row := db.ptr.QueryRow(
		"SELECT id, hashp, salt, priv FROM users WHERE hashn=$1",
		gp.Base64Encode(gp.HashSum([]byte(name))),
	)
	row.Scan(&id, &hpasw, &ssalt, &spriv)
	if spriv == "" {
		return nil
	}
	salt := gp.Base64Decode(ssalt)
	bpasw := gp.RaiseEntropy([]byte(pasw), salt, PASWDIFF)
	chpasw := gp.HashSum(bytes.Join(
		[][]byte{
			bpasw,
			[]byte(name),
		},
		[]byte{},
	))
	if !bytes.Equal(chpasw, gp.Base64Decode(hpasw)) {
		return nil
	}
	priv := gp.BytesToPrivateKey(gp.DecryptAES(bpasw, gp.Base64Decode(spriv)))
	if priv == nil {
		return nil
	}
	return &User{
		Id:   id,
		Name: name,
		Pasw: bpasw,
		Priv: priv,
	}
}

func (db *DB) userExist(name string) bool {
	var (
		namee string
	)
	row := db.ptr.QueryRow(
		"SELECT name FROM users WHERE hashn=$1",
		gp.Base64Encode(gp.HashSum([]byte(name))),
	)
	row.Scan(&namee)
	return namee != ""
}
