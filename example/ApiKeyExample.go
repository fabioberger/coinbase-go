package main

import (
	"fmt"
	"log"
	"os"

	"github.com/fabioberger/coinbase-go"
)

func main() {
	// Instantiate a ApiKeyClient with key and secret set at environment variables
	c := coinbase.ApiKeyClient(os.Getenv("COINBASE_KEY"), os.Getenv("COINBASE_SECRET"))

	balance, err := c.GetBalance()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Balance is %f BTC", balance)

}
