package coinbase

import (
	"net/http"
)

const API_BASE = "https://api.coinbase.com/v1/" //"https://www.bitstamp.net/"

type Client struct {
	rpc Rpc
}

// Instantiate the client with ApiKey Authentication
func ApiKeyClient(key string, secret string) Client {
	c := Client{
		rpc: Rpc{
			client: http.Client{
				Transport: &http.Transport{
					Dial: dialTimeout,
				},
			},
			auth: ApiKeyAuthentication{
				Key:    key,
				Secret: secret,
			},
		},
	}
	return c
}

func (c Client) Get(path string, params map[string]interface{}) ([]byte, error) {
	return c.rpc.Request("GET", path, params)
}

func (c Client) Post(path string, params map[string]interface{}) ([]byte, error) {
	return c.rpc.Request("POST", path, params)
}

func (c Client) Delete(path string, params map[string]interface{}) ([]byte, error) {
	return c.rpc.Request("DELETE", path, params)
}

func (c Client) Put(path string, params map[string]interface{}) ([]byte, error) {
	return c.rpc.Request("PUT", path, params)
}
