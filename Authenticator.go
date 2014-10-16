package coinbase

import (
	"net/http"
)

// Create an authentication interface to enable polymorphism between auth types
type Authenticator interface {
	Authenticate(req *http.Request, message string, nonce string) error
}
