package coinbase

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strconv"
	"time"
)

// ApiKeyAuthentication Struct implements the Authentication interface and takes
// care of authenticating RPC requests for clients with a Key & Secret pair
type apiKeyAuthentication struct {
	Key     string
	Secret  string
	BaseUrl string
	Client  http.Client
}

// ApiKeyAuthWithEnv instantiates ApiKeyAuthentication with the API key & secret & environment (Live or Sandbox)
func apiKeyAuthWithEnv(key string, secret string, sandbox bool) *apiKeyAuthentication {
	baseUrl := "https://api.coinbase.com/v1/" // Live Url
	
	// Check if should use sandbox
	if sandbox {
		baseUrl = "https://api.sandbox.coinbase.com/v1/" // Sandbox Url
	}
	a := apiKeyAuthentication{
		Key:     key,
		Secret:  secret,
		BaseUrl: baseUrl,
		Client: http.Client{
			Transport: &http.Transport{
				Dial: dialTimeout,
			},
		},
	}
	return &a
}

// ApiKeyAuth instantiates ApiKeyAuthentication with the API key & secret
// TODO: Maybe remove this (not sure if it would break backwards compatability)
func apiKeyAuth(key string, secret string) *apiKeyAuthentication {
	return apiKeyAuthWithEnv(key, secret, false)
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
