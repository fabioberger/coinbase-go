package coinbase

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strconv"
	"time"

	"github.com/fabioberger/coinbase-go/config"
)

// ApiKeyAuthentication Struct implements the Authentication interface and takes
// care of authenticating RPC requests for clients with a Key & Secret pair
type apiKeyAuthentication struct {
	Key     string
	Secret  string
	BaseUrl string
	Client  http.Client
}

// ApiKeyAuth instantiates ApiKeyAuthentication with the API key & secret
func apiKeyAuth(key string, secret string) *apiKeyAuthentication {
	a := apiKeyAuthentication{
		Key:     key,
		Secret:  secret,
		BaseUrl: config.BaseUrl,
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
func (a apiKeyAuthentication) authenticate(req *http.Request, endpoint string, params []byte) error {

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

func (a apiKeyAuthentication) getBaseUrl() string {
	return a.BaseUrl
}

func (a apiKeyAuthentication) getClient() *http.Client {
	return &a.Client
}
