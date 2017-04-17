//
// This is a thin thread-safe wrapper over the standard spf13/viper package.
// It supports simple remote polling with consul for config changes over YAML.
// We do not use much of the core viper code for doing this because none of it is thread-
// safe or supports 'dynamic reloading' the way I'd like to.
//
package viper

import (
	"errors"
	"time"
	//"fmt"

	"github.com/hashicorp/consul/api"
)

// StartWatcher encapsulates a basic polling watcher on a given consul key.
// Calling this routine spawns a goroutine which polls a given consul address for a specific consul key.
// It tracks internal state based on the last value read and only sends nil to the channel when the config
//
// It returns an error channel, or nil and an error if it is unable to create a new consul client.
func StartWatcher(consulAddr string, consulKey string, delayInSeconds int) (chan error, error) {
	client, err := api.NewClient(&api.Config{Address: consulAddr})
	if err != nil {
		return nil, err
	}

	watcherCh := make(chan error)

	go func(c *api.Client, caddr string, ckey string, delay int) {
		lastVal := ""
		for {
			time.Sleep(time.Second * time.Duration(delay)) // delay after each request

			// TODO: avoid reading remote config twice?
			val, err := consulGet(c, ckey)
			// val changed somehow, reload
			if err != nil {
				// TODO: error, couldn't read remote config, should we do our own backoff?
				watcherCh <- err
				continue
			}

			if val != lastVal && val != "" {
				err = ReadRemoteConfig()
				if err != nil {
					// TODO: error, couldn't read remote config, should we do our own backoff?
					watcherCh <- err
				} else {
					lastVal = val
					//fmt.Println("just set lastVal to: " + lastVal)
					watcherCh <- nil
				}
			}
		}
	}(client, consulAddr, consulKey, delayInSeconds)
	return watcherCh, nil
}

func StartWatcherListener(watcherCh chan error) {
	go func(wch chan error) {
		for {
			select {
			case err := <-wch:
				if err != nil {
					//fmt.Println("error:")
					//fmt.Println(err)
				} else {
					//fmt.Println("config should be reloaded now?")
				}
			}
		}
	}(watcherCh)
}

func consulGet(c *api.Client, k string) (string, error) {
	// Get a handle to the KV API
	kv := c.KV()

	v, _, err := kv.Get(k, &api.QueryOptions{})
	if err != nil {
		return "", err
	}
	// This can happen if consul goes away and then comes back with an
	// non-existing key
	if v == nil {
		return "", errors.New("consul value was non-existent")
	}
	return string(v.Value), nil
}
