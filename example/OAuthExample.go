package main

import (
	"fmt"
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/fabioberger/coinbase-go"
)

func main() {
	mux := http.NewServeMux()

	// At http://localhost:3000/ we will display an "authorize" link
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		o, err := GetOAuthService()
		if err != nil {
			fmt.Fprintln(w, err)
			return
		}
		authorizeUrl := o.CreateAuthorizeUrl([]string{
			"all",
		})
		link := "<a href='" + authorizeUrl + "'>authorize</a>"
		fmt.Fprintln(w, link)
	})

	// AuthorizeUrl redirects to https://localhost:3000/tokens with code in its
	// query params. If you dont have SSL enabled, replace 'https' with 'http'
	// and reload the page. If successful, the user's balance will show
	mux.HandleFunc("/tokens", func(w http.ResponseWriter, req *http.Request) {
		query := req.URL.Query()
		code := query.Get("code")
		o, err := GetOAuthService()
		tokens, err := o.GetTokens(code, "authorization_code")
		if err != nil {
			fmt.Fprintln(w, err)
			return
		}
		c := coinbase.OAuthClient(tokens)
		amount, err := c.GetBalance()
		if err != nil {
			fmt.Fprintln(w, err)
			return
		}
		fmt.Fprintln(w, amount)
	})

	n := negroni.Classic()
	n.UseHandler(mux)
	n.Run(":3000")
}

func GetOAuthService() (*coinbase.OAuth, error) {
	// Be sure to replace CLIENT_ID and CLIENT_SECRET with your OAuth app id and secret
	return coinbase.OAuthService(CLIENT_ID, CLIENT_SECRET, "https://localhost:3000/tokens")
}
