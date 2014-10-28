package coinbase

import (
	"errors"
	"net"
	"strings"
	"time"
)

// checkApiErrors checks for errors returned by coinbase API JSON response
// i.e { "success": false, "errors": ["Button with code code123456 does not exist"], ...}
func checkApiErrors(resp Response, method string) error {
	if resp.Success == false { // Return errors received from API here
		err := " in " + method + "()"
		if resp.Errors != nil {
			err = strings.Join(resp.Errors, ",") + err
			return errors.New(err)
		}
		if resp.Error != "" { // Return errors received from API here
			err = resp.Error + err
			return errors.New(err)
		}
	}
	return nil
}

// dialTimeout is used to enforce a timeout for all http requests.
func dialTimeout(network, addr string) (net.Conn, error) {
	var timeout = time.Duration(2 * time.Second) //how long to wait when trying to connect to the coinbase
	return net.DialTimeout(network, addr, timeout)
}
