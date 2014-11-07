package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/fabioberger/coinbase-go"
	"github.com/go-martini/martini"
)

var o *coinbase.OAuth

func main() {
	m := martini.New()
	m.Use(martini.Logger())
	m.Use(martini.Recovery())
	m.Use(martini.Static("public"))
	r := martini.NewRouter()
	m.MapTo(r, (*martini.Routes)(nil))
	m.Action(r.Handle)

	// Instantiate OAuthService with the OAuth App Client Id & Secret from the Environment Variables
	o, err := coinbase.OAuthService(os.Getenv("COINBASE_CLIENT_ID"), os.Getenv("COINBASE_CLIENT_SECRET"), "https://localhost:8443/tokens")
	if err != nil {
		panic(err)
		return
	}

	// At https://localhost:8443/ we will display an "authorize" link
	r.Get("/", func() string {
		authorizeUrl := o.CreateAuthorizeUrl([]string{
			"all",
		})
		link := "<a href='" + authorizeUrl + "'>authorize</a>"
		return link
	})

	// AuthorizeUrl redirects to https://localhost:8443/tokens with 'code' in its
	// query params. If you dont have SSL enabled, replace 'https' with 'http'
	// and reload the page. If successful, the user's balance will show
	r.Get("/tokens", func(res http.ResponseWriter, req *http.Request) string {
		// Get the tokens given the 'code' query param
		tokens, err := o.NewTokensFromRequest(req) // Will use 'code' query param from req
		if err != nil {
			return err.Error()
		}
		// instantiate the OAuthClient
		c := coinbase.OAuthClient(tokens)
		amount, err := c.GetBalance()
		if err != nil {
			return err.Error()
		}
		return strconv.FormatFloat(amount, 'f', 6, 64)
	})

	// HTTP
	go func() {
		if err := http.ListenAndServe(":8080", m); err != nil {
			log.Fatal(err)
		}
	}()

	// HTTPS
	// To generate a development cert and key, run the following from your *nix terminal:
	// go run $(go env GOROOT)/src/pkg/crypto/tls/generate_cert.go --host="localhost"
	if err := http.ListenAndServeTLS(":8443", "cert.pem", "key.pem", m); err != nil {
		log.Fatal(err)
	}
}
