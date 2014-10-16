package coinbase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"
)

// Rpc handles the remote procedure call requests
type Rpc struct {
	auth   Authenticator
	client http.Client
}

// dialTimeout is used to enforce a timeout for all http requests.
func dialTimeout(network, addr string) (net.Conn, error) {
	var timeout = time.Duration(2 * time.Second) //how long to wait when trying to connect to the coinbase
	return net.DialTimeout(network, addr, timeout)
}

func (r Rpc) Request(method string, endpoint string, params map[string]interface{}) ([]byte, error) {

	request, err := r.createRequest(method, endpoint, params)
	if err != nil {
		return nil, err
	}
	data, err := r.executeRequest(request)

	return data, err
}

// CreateRequest formats a request with all the necessary authenticated headers
func (r Rpc) createRequest(method string, endpoint string, params map[string]interface{}) (*http.Request, error) {

	if params == nil {
		params = map[string]interface{}{}
	}

	nonce := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
	endpoint = API_BASE + endpoint

	// Convert params payload to json
	payloadJson, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	message := nonce + endpoint + string(payloadJson) //Needed for HMAC Signature

	req, err := http.NewRequest(method, endpoint, bytes.NewBuffer(payloadJson))
	if err != nil {
		return nil, err
	}

	// Authenticate the request
	r.auth.Authenticate(req, message, nonce)

	req.Header.Set("User-Agent", "CoinbasePHP/v1")
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

// executeRequest takes a request and returns the body of the response
func (r Rpc) executeRequest(req *http.Request) ([]byte, error) {
	resp, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bytes := buf.Bytes()
	if resp.StatusCode != 200 {
		if len(bytes) == 0 {
			log.Printf("Response body was empty")
		} else {
			log.Printf("Response body:\n\t%s\n", bytes)
		}
		return nil, fmt.Errorf("%s %s failed. Response code was %s", req.Method, req.URL, resp.Status)
	}
	return bytes, nil
}
