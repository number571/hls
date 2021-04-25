package main

import (
	"crypto/rsa"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"sync"
)

type DB struct {
	ptr *sql.DB
	mtx sync.Mutex
}

type User struct {
	Id   int
	Name string
	Pasw []byte
	Priv *rsa.PrivateKey
}
