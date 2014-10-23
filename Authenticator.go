package coinbase

import (
	"net/http"
)

// Authenticator is an interface that objects can implement in order to act as the
// authentication mechanism for RPC requests to Coinbase
type Authenticator interface {
	GetBaseUrl() string
	GetClient() *http.Client
	Authenticate(req *http.Request, endpoint string, params []byte) error
}
