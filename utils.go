package main

import (
	"encoding/json"
	"io/ioutil"
)

func readFile(file string) []byte {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil
	}
	return data
}

func writeFile(file string, data []byte) error {
	return ioutil.WriteFile(file, data, 0644)
}

func serialize(data interface{}) []byte {
	res, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return nil
	}
	return res
}

func deserialize(data []byte, res interface{}) error {
	return json.Unmarshal(data, res)
}
