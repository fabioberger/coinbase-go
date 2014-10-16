// Coinbase-go is a convenient Go wrapper for the Coinbase API
package coinbase

import (
	"encoding/json"
	"net/http"
)

const API_BASE = "https://api.coinbase.com/v1/" //"https://www.bitstamp.net/"

// The Client from which all API requests are made
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

func (c Client) getPaginatedResource(resource string, listElement string, unwrapElement string, params map[string]interface{}, holder interface{}) error {
	data, err := c.Get(resource, params)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, &holder); err != nil {
		return err
	}
	return nil
}

func (c Client) GetBalance() (string, error) {
	data, err := c.Get("account/balance", nil)
	if err != nil {
		return "", err
	}
	balance := map[string]string{}
	if err := json.Unmarshal(data, &balance); err != nil {
		return "", err
	}
	return balance["amount"], nil
}

func (c Client) GetReceiveAddress() (string, error) {
	data, err := c.Get("account/receive_address", nil)
	if err != nil {
		return "", err
	}
	holder := map[string]interface{}{}
	if err := json.Unmarshal(data, &holder); err != nil {
		return "", err
	}
	return holder["address"].(string), nil
}

func (c Client) GetAllAddresses(query string, page int, limit int) (addressesHolder, error) {
	params := map[string]interface{}{}
	if query != "" {
		params["query"] = query
	}
	if limit != 0 {
		params["limit"] = limit
	}
	params["page"] = page

	holder := addressesHolder{}

	err := c.getPaginatedResource("addresses", "addresses", "address", params, &holder)
	if err != nil {
		return addressesHolder{}, err
	}
	return holder, nil
}

func (c Client) GenerateReceiveAddress(callback string, label string) (string, error) {
	params := map[string]interface{}{}
	if callback != "" {
		params["address[callback_url]"] = callback
	}
	if label != "" {
		params["address[label]"] = label
	}
	data, err := c.Post("account/generate_receive_address", params)
	if err != nil {
		return "", err
	}
	holder := map[string]interface{}{}
	if err := json.Unmarshal(data, &holder); err != nil {
		return "", err
	}
	return holder["address"].(string), nil
}

// func (c Client) sendMoney(to string, amount float64, notes string, userFee float64, amountCurrency string) ([]byte, error) {
// 	params := map[string]interface{}{
// 		"transaction[to]": to,
// 	}

// 	if(amountCurrency != "") {

// 	}

// }
