package settings

import (
	"encoding/json"
	"io/ioutil"
	"os"

	gp "github.com/number571/go-peer/settings"
)

type Request struct {
	Host   string
	Path   string
	Method string
	Head   map[string]string
	Body   []byte
}

const (
	AKEYSIZE           = 2048
	HLS                = "hidden-lake-service"
	FileWithPubKey     = "pub.key"
	ServerAddressInHLS = "route-service"
)

var (
	SETTINGS = gp.NewSettings()
)

func FileIsExist(file string) bool {
	_, err := os.Stat(file)
	return !os.IsNotExist(err)
}

func ReadFile(file string) []byte {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil
	}
	return data
}

func WriteFile(file string, data []byte) error {
	return ioutil.WriteFile(file, data, 0644)
}

func Serialize(data interface{}) []byte {
	res, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return nil
	}
	return res
}

func Deserialize(data []byte, res interface{}) error {
	return json.Unmarshal(data, res)
}
