package main

import (
	"fmt"
	"net/http"
	"os"

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
	// Be sure your coinbase OAuth CLIENT_ID and CLIENT_SECRET are set as shell environment variables
	// In shell run:
	// export COINBASE_CLIENT_ID="YOUR_CLIENT_ID"
	// export COINBASE_CLIENT_SECRET="YOUR_CLIENT_SECRET"
	return coinbase.OAuthService(os.Getenv("COINBASE_CLIENT_ID"), os.Getenv("COINBASE_CLIENT_SECRET"), "https://localhost:3000/tokens")
}
