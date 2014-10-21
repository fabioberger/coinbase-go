package coinbase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Rpc handles the remote procedure call requests
type Rpc struct {
	auth   Authenticator
	client http.Client
	mock   bool
}

// dialTimeout is used to enforce a timeout for all http requests.
func dialTimeout(network, addr string) (net.Conn, error) {
	var timeout = time.Duration(2 * time.Second) //how long to wait when trying to connect to the coinbase
	return net.DialTimeout(network, addr, timeout)
}

func (r Rpc) Request(method string, endpoint string, params interface{}, holder interface{}) error {

	jsonParams, err := json.Marshal(params)
	if err != nil {
		return err
	}

	request, err := r.createRequest(method, endpoint, jsonParams)
	if err != nil {
		return err
	}

	var data []byte
	if r.mock == true {
		data, err = r.simulateRequest(endpoint, method)
	} else {
		data, err = r.executeRequest(request)
	}

	if err := json.Unmarshal(data, &holder); err != nil {
		return err
	}

	return nil
}

// CreateRequest formats a request with all the necessary headers
func (r Rpc) createRequest(method string, endpoint string, params []byte) (*http.Request, error) {

	nonce := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
	endpoint = API_BASE + endpoint

	message := nonce + endpoint + string(params) //Needed for HMAC Signature

	req, err := http.NewRequest(method, endpoint, bytes.NewBuffer(params))
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

// Simulate a request by returning a sample JSON from file
func (r Rpc) simulateRequest(endpoint string, method string) ([]byte, error) {
	fileName := strings.Replace(endpoint, "/", "_", -1)
	filePath := "test_data/" + method + "_" + fileName + ".json"
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return data, nil
}
