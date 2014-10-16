# Coinbase Go Client Library

An easy way to buy, send, and accept [bitcoin](http://en.wikipedia.org/wiki/Bitcoin) through the [Coinbase API](https://coinbase.com/docs/api/overview).

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
