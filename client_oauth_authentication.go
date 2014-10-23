package coinbase

import (
	"errors"
	"net/http"
	"time"
)

// Service OAuth Authentication Struct implements the Authentication interface
// and takes care of authenticating OAuth RPC requests on behalf of a client
// (i.e GetBalance())
type ClientOAuthAuthentication struct {
	Tokens  *oauthTokens
	BaseUrl string
	Client  http.Client
}

// Instantiate ClientOAuthAuthentication with the client OAuth tokens
func ClientOAuth(tokens *oauthTokens) *ClientOAuthAuthentication {
	a := ClientOAuthAuthentication{
		Tokens:  tokens,
		BaseUrl: "https://api.coinbase.com/v1/",
		Client: http.Client{
			Transport: &http.Transport{
				Dial: dialTimeout,
			},
		},
	}
	return &a
}

// Client OAuth authentication requires us to attach an unexpired OAuth token to
// the request header
func (a ClientOAuthAuthentication) Authenticate(req *http.Request, endpoint string, params []byte) error {
	// Ensure tokens havent expired
	if time.Now().UTC().Unix() > a.Tokens.Expire_time {
		return errors.New("The OAuth tokens are expired. Use refreshTokens to refresh them")
	}
	req.Header.Set("Authorization", "Bearer "+a.Tokens.Access_token)
	return nil
}

func (a ClientOAuthAuthentication) GetBaseUrl() string {
	return a.BaseUrl
}

func (a ClientOAuthAuthentication) GetClient() *http.Client {
	return &a.Client
}
