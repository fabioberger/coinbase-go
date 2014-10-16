# Coinbase Go Client Library

An easy way to buy, send, and accept [bitcoin](http://en.wikipedia.org/wiki/Bitcoin) through the [Coinbase API](https://coinbase.com/docs/api/overview).

This library supports both the [API key authentication method](https://coinbase.com/docs/api/overview) and OAuth. The below examples use an API key - for instructions on how to use OAuth, see [OAuth Authentication](#oauth-authentication).

## Installation

Obtain the latest version of the Coinbase Go library with:

    git clone https://github.com/fabioberger/coinbase-go

Then, add the following to your Go project:

```go
import (
	"github.com/fabioberger/coinbase-go"
)
```

## Usage

Start by [enabling an API Key on your account](https://coinbase.com/settings/api).

Next, create an instance of the client using the `ApiKeyClient` method:

```go
c := coinbase.ApiKeyClient(os.Getenv("COINBASE_KEY"), os.Getenv("COINBASE_SECRET"))
```

Now you can call methods on `c` similar to the ones described in the [API reference](https://coinbase.com/api/doc).  For example:

```go
amount, err := c.GetBalance()
if err != nil {
	log.Fatal(err)
}
fmt.Println(amount)
```
