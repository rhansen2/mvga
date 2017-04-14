//
// Upgraded viper package
//
package viper

import (
	"errors"
	"github.com/hashicorp/consul/api"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"os"

	"sync/atomic"
	"unsafe"
)

// ErrNilReadFromConsul is for when determing when consul has read an unset key
var ErrNilReadFromConsul = errors.New("nil was read")

// Global internal variables
//var mtx sync.Mutex
var cfgType string
var cfgFilePath string
var cfgContents *string

var consulAddr string
var consulKey string

// Setup mutex
func init() {
	//mtx = sync.Mutex{}
}

// SetConfigType only accepts JSON for the time being, this is mostly a placeholder method
func SetConfigType(t string) {
	//mtx.Lock()
	cfgType = t
	//mtx.Unlock()
}

func SetConfigFile(f string) {
	cfgFilePath = f
}

func ReadInConfig() error {
	if _, err := os.Stat(cfgFilePath); err != nil {
		return errors.New("config path not valid")
	}

	byteVal, err := ioutil.ReadFile(cfgFilePath)
	if err != nil {
		return err
	}
	tmpByteVal := string(byteVal)
	UpdateConfig(&tmpByteVal)
	return nil
}

var config unsafe.Pointer // actual type is *Config
// CurrentConfig atomically returns the current configuration
func CurrentConfig() *string { return (*string)(atomic.LoadPointer(&config)) }

// UpdateConfig atomically swaps the current configuration
func UpdateConfig(cfg *string) { atomic.StorePointer(&config, unsafe.Pointer(cfg)) }

// SetRemoteProvider t is type, currently unused, addr is consul address, keyPref is our configuration key
func SetRemoteProvider(t string, addr string, keyPref string) {
	consulAddr = addr
	consulKey = keyPref
}

// ReadRemoteConfig reads our remote config
func ReadRemoteConfig() error {
	client, err := api.NewClient(&api.Config{Address: consulAddr})
	if err != nil {
		return err
	}

	data, err := consulGet(client, consulKey)
	if err != nil {
		return err
	}
	UpdateConfig(&data)
	return nil
}

func GetStringSlice(k string) []string {
	f := gjson.Get(*CurrentConfig(), k)

	fArray := f.Array()
	strSlice := make([]string, len(fArray))
	for i, v := range fArray {
		strSlice[i] = v.String()
	}

	return strSlice
}

func GetString(k string) string {
	s := gjson.Get(*CurrentConfig(), k)
	return s.String()
}

func GetBool(k string) bool {
	b := gjson.Get(*CurrentConfig(), k)
	return b.Bool()
}

func GetInt(k string) int {
	b := gjson.Get(*CurrentConfig(), k)
	return int(b.Int())
}

func GetInt64(k string) int64 {
	b := gjson.Get(*CurrentConfig(), k)
	return b.Int()
}
