//
// Upgraded viper package
//
package viper

import (
	"os"
	"errors"
	"sync"
	"io/ioutil"
	"github.com/tidwall/gjson"
	"github.com/hashicorp/consul/api"
)

// ErrNilReadFromConsul is for when determing when consul has read an unset key
var ErrNilReadFromConsul = errors.New("nil was read")

// Global internal variables
var mtx sync.Mutex
var cfgType string
var cfgFilePath string
var cfgContents string

var consulAddr string
var consulKey string

// Setup mutex
func init() {
	mtx = sync.Mutex{}
}

// SetConfigType only accepts JSON for the time being, this is mostly a placeholder method
func SetConfigType(t string) {
	mtx.Lock()
	cfgType = t
	mtx.Unlock()
}

func SetConfigFile(f string) {
	mtx.Lock()
	cfgFilePath = f
	mtx.Unlock()
}

func ReadInConfig() error {
        if _, err := os.Stat(cfgFilePath); err != nil {
                return errors.New("config path not valid")
        }

	mtx.Lock()
	byteVal, err := ioutil.ReadFile(cfgFilePath)
	cfgContents = string(byteVal)
	mtx.Unlock()
	return err
}

// SetRemoteProvider t is type, currently unused, addr is consul address, keyPref is our configuration key
// (typically servicename/config, a la encrypt/config, guid/config, etc.)
func SetRemoteProvider(t string, addr string, keyPref string) {
	mtx.Lock()
	consulAddr = addr
	consulKey = keyPref
	mtx.Unlock()
}


// ReadRemoteConfig reads our remote config
func ReadRemoteConfig() error {
	var err error

        client, err := api.NewClient(&api.Config{Address: consulAddr})
	if err != nil {
		return err
	}

	mtx.Lock()
	cfgContents, err = consulGet( client, consulKey )
	mtx.Unlock()
	return err
}

func GetStringSlice(k string) []string {
	mtx.Lock()
	f := gjson.Get( cfgContents, k )
	mtx.Unlock()

	fArray := f.Array()
	strSlice := make([]string, len ( fArray ) )
	for i, v := range fArray {
		strSlice[i] = v.String()
	}

	return strSlice
}

func GetString(k string) string {
	mtx.Lock()
	s := gjson.Get( cfgContents, k )
	mtx.Unlock()
	return s.String()
}

func GetBool(k string) bool {
	mtx.Lock()
	b := gjson.Get( cfgContents, k )
	mtx.Unlock()
	return b.Bool()
}

func GetInt(k string) int {
	mtx.Lock()
	b := gjson.Get( cfgContents, k )
	mtx.Unlock()
	return int(b.Int())
}

func GetInt64(k string) int64 {
	mtx.Lock()
	b := gjson.Get( cfgContents, k )
	mtx.Unlock()
	return b.Int()
}

