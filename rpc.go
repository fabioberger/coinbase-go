package coinbase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

// basePath needed for reading mock JSON files in simulateRequest and for
// referencing ca-coinbase.crt in OAuthService. Must be fixed path because
// relative path changes depending on where library called
var basePath string = os.Getenv("GOPATH") + "/src/github.com/fabioberger/coinbase-go"

// Rpc handles the remote procedure call requests
type rpc struct {
	auth authenticator
	mock bool
}

// Request sends a request with params marshaled into a JSON payload in the body
// The response value is marshaled from JSON into the specified holder struct
func (r rpc) Request(method string, endpoint string, params interface{}, holder interface{}) error {

	jsonParams, err := json.Marshal(params)
	if err != nil {
		return err
	}

	request, err := r.createRequest(method, endpoint, jsonParams)
	if err != nil {
		return err
	}

	var data []byte
	if r.mock == true { // Mock mode: Replace actual request with expected JSON from file
		data, err = r.simulateRequest(endpoint, method)
	} else {
		data, err = r.executeRequest(request)
	}
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, &holder); err != nil {
		return err
	}

	return nil
}

// CreateRequest formats a request with all the necessary headers
func (r rpc) createRequest(method string, endpoint string, params []byte) (*http.Request, error) {

	endpoint = r.auth.getBaseUrl() + endpoint //BaseUrl depends on Auth type used

	req, err := http.NewRequest(method, endpoint, bytes.NewBuffer(params))
	if err != nil {
		return nil, err
	}

	// Authenticate the request
	r.auth.authenticate(req, endpoint, params)

	req.Header.Set("User-Agent", "CoinbaseGo/v1")
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

// executeRequest takes a prepared http.Request and returns the body of the response
// If the response is not of HTTP Code 200, an error is returned
func (r rpc) executeRequest(req *http.Request) ([]byte, error) {
	resp, err := r.auth.getClient().Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bytes := buf.Bytes()
	if resp.StatusCode != 200 {
		if len(bytes) == 0 { // Log response body for debugging purposes
			log.Printf("Response body was empty")
		} else {
			log.Printf("Response body:\n\t%s\n", bytes)
		}
		return nil, fmt.Errorf("%s %s failed. Response code was %s", req.Method, req.URL, resp.Status)
	}
	return bytes, nil
}

// simulateRequest simulates a request by returning a sample JSON from file
func (r rpc) simulateRequest(endpoint string, method string) ([]byte, error) {
	// Test files conform to replacing '/' in endpoint with '_'
	fileName := strings.Replace(endpoint, "/", "_", -1)
	// file names also have method type prepended to ensure uniqueness
	filePath := basePath + "/test_data/" + method + "_" + fileName + ".json"
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return data, nil
}
