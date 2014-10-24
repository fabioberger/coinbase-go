package coinbase

import (
	"net/url"
	"strings"
	"time"
)

// OAuth handles all service oauth related functionality (i.e GetTokens(), RefreshTokens()
type OAuth struct {
	ClientId     string
	ClientSecret string
	RedirectUri  string
	Rpc          Rpc
}

// OAuthService Instantiates OAuth Struct in order to send service related OAuth requests
func OAuthService(clientId string, clientSecret string, redirectUri string) (*OAuth, error) {
	certFilePath := BASE_PATH + "/ca-coinbase.crt"
	serviceAuth, err := ServiceOAuth(certFilePath)
	if err != nil {
		return nil, err
	}
	o := OAuth{
		ClientId:     clientId,
		ClientSecret: clientSecret,
		RedirectUri:  redirectUri,
		Rpc: Rpc{
			auth: serviceAuth,
			mock: false,
		},
	}
	return &o, nil
}

// CreateAuthorizeUrl create the Authorize Url used to redirect users for
// coinbase app authorization
func (o OAuth) CreateAuthorizeUrl(scope []string) string {
	Url, _ := url.Parse("https://coinbase.com")
	Url.Path += "/oauth/authorize"

	parameters := url.Values{}
	parameters.Add("response_type", "code")
	parameters.Add("client_id", o.ClientId)
	parameters.Add("redirect_uri", o.RedirectUri)
	parameters.Add("scope", strings.Join(scope, "+"))
	Url.RawQuery = parameters.Encode()

	return Url.String()
}

// RefreshTokens refreshes a users existing OAuth tokens
func (o OAuth) RefreshTokens(oldTokens map[string]interface{}) (*oauthTokens, error) {
	refresh_token := oldTokens["refresh_token"].(string)
	return o.GetTokens(refresh_token, "refresh_token")
}

// NewTokens generates new tokens for an OAuth user
func (o OAuth) NewTokens(code string) (*oauthTokens, error) {
	return o.GetTokens(code, "authorization_code")
}

// GetTokens gets tokens for an OAuth user specifying a grantType (i.e authorization_code)
func (o OAuth) GetTokens(code string, grantType string) (*oauthTokens, error) {

	postVars := map[string]string{
		"grant_type":    grantType,
		"redirect_uri":  o.RedirectUri,
		"client_id":     o.ClientId,
		"client_secret": o.ClientSecret,
	}

	if grantType == "refresh_token" {
		postVars["refresh_token"] = code
	} else {
		postVars["code"] = code
	}
	holder := tokensHolder{}
	err := o.Rpc.Request("POST", "oauth/token", postVars, &holder)
	if err != nil {
		return nil, err
	}

	tokens := oauthTokens{
		Access_token:  holder.Access_token,
		Refresh_token: holder.Refresh_token,
		Expire_time:   time.Now().UTC().Unix() + holder.Expires_in,
	}

	return &tokens, nil
}
