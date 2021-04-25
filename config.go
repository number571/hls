package main 

import (
	"os"
	"io/ioutil"
	"encoding/json"
)

type CFG struct {
	Openaddr   string            `json:"openaddr"`
	Connects   []string          `json:"connects"`
	Services   map[string]string `json:"services"`
}

func NewCFG(filename string) *CFG {
	var config = new(CFG)
	if !fileIsExist(filename) {
		config.Openaddr = ""
		config.Connects = []string{
			"addr",
		}
		config.Services = map[string]string{
			"host": "addr",
		}
		err := ioutil.WriteFile(filename, serialize(config), 0644)
		if err != nil {
			return nil
		}
	}
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil
	}
	err = json.Unmarshal(content, config)
	if err != nil {
		return nil
	}
	return config
}

func fileIsExist(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}
	return true
}
