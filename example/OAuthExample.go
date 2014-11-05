package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/codegangsta/negroni"
	"github.com/fabioberger/coinbase-go"
)

var o *coinbase.OAuth

func main() {
	mux := http.NewServeMux()
	o, err := coinbase.OAuthService(os.Getenv("COINBASE_CLIENT_ID"), os.Getenv("COINBASE_CLIENT_SECRET"), "https://localhost:3000/tokens")
	if err != nil {
		panic(err)
		return
	}
	// At http://localhost:3000/ we will display an "authorize" link
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		authorizeUrl := o.CreateAuthorizeUrl([]string{
			"all",
		})
		link := "<a href='" + authorizeUrl + "'>authorize</a>"
		fmt.Fprintln(w, link)
	})

	// AuthorizeUrl redirects to https://localhost:3000/tokens with 'code' in its
	// query params. If you dont have SSL enabled, replace 'https' with 'http'
	// and reload the page. If successful, the user's balance will show
	mux.HandleFunc("/tokens", func(w http.ResponseWriter, req *http.Request) {
		// Get the tokens given the 'code' query param
		tokens, err := o.NewTokensFromRequest(req) // Will use 'code' query param from req
		if err != nil {
			fmt.Fprintln(w, err)
			return
		}
		// instantiate the OAuthClient
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
