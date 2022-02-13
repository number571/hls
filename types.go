package main

type CFG struct {
	Address  string            `json:"address"`
	Connects []string          `json:"connects"`
	Services map[string]string `json:"services"`
}

type Request struct {
	Host   string
	Path   string
	Method string
	Head   map[string]string
	Body   []byte
}
