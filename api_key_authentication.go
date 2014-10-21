package coinbase

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
)

// ApiKey Authentication Struct of type Authentication
type ApiKeyAuthentication struct {
	Key    string
	Secret string
}

func (a ApiKeyAuthentication) Authenticate(req *http.Request, message string, nonce string) error {
	req.Header.Set("ACCESS_KEY", a.Key)

	h := hmac.New(sha256.New, []byte(a.Secret))
	h.Write([]byte(message))
	signature := hex.EncodeToString(h.Sum(nil))
	req.Header.Set("ACCESS_SIGNATURE", signature)
	req.Header.Set("ACCESS_NONCE", nonce)
	return nil
}
