package coinbase

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strconv"
	"time"
)

// ApiKey Authentication Struct implements the Authentication interface and takes
// care of authenticating RPC requests for clients with a Key & Secret pair
type ApiKeyAuthentication struct {
	Key     string
	Secret  string
	BaseUrl string
	Client  http.Client
}

// Instantiate ApiKeyAuthentication with the API key & secret
func ApiKeyAuth(key string, secret string) *ApiKeyAuthentication {
	a := ApiKeyAuthentication{
		Key:     key,
		Secret:  secret,
		BaseUrl: "https://api.coinbase.com/v1/",
		Client: http.Client{
			Transport: &http.Transport{
				Dial: dialTimeout,
			},
		},
	}
	return &a
}

// API Key + Secret authentication requires a request header of the HMAC SHA-256
// signature of the "message" as well as an incrementing nonce and the API key
func (a ApiKeyAuthentication) Authenticate(req *http.Request, endpoint string, params []byte) error {

	nonce := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
	message := nonce + endpoint + string(params) //As per Coinbase Documentation

	req.Header.Set("ACCESS_KEY", a.Key)

	h := hmac.New(sha256.New, []byte(a.Secret))
	h.Write([]byte(message))

	signature := hex.EncodeToString(h.Sum(nil))

	req.Header.Set("ACCESS_SIGNATURE", signature)
	req.Header.Set("ACCESS_NONCE", nonce)

	return nil
}

func (a ApiKeyAuthentication) GetBaseUrl() string {
	return a.BaseUrl
}

func (a ApiKeyAuthentication) GetClient() *http.Client {
	return &a.Client
}
