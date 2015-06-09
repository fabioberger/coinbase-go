[![GoDoc](https://godoc.org/github.com/fabioberger/coinbase-go?status.svg)](https://godoc.org/github.com/fabioberger/coinbase-go)

# Coinbase Go Client Library

An easy way to buy, send, and accept [bitcoin](http://en.wikipedia.org/wiki/Bitcoin) through the [Coinbase API](https://coinbase.com/docs/api/overview).

This library supports both the [API key authentication method](https://coinbase.com/docs/api/overview) and OAuth. The below examples use an API key - for instructions on how to use OAuth, see [OAuth Authentication](#oauth-authentication).

A detailed step-by-step tutorial on how to use this library can be found at [this blog](http://fabioberger.com/blog/2014/11/06/building-a-coinbase-app-in-go/).

## Installation

Make sure you have set the environment variable $GOPATH

```bash
export GOPATH="path/to/your/go/folder"
```

Obtain the latest version of the Coinbase Go library with:

```bash
go get github.com/fabioberger/coinbase-go
```

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

Notice here that we did not hard code the API key into our codebase, but set it in an environment variable instead.  This is just one example, but keeping your credentials separate from your code base is a good [security practice](https://coinbase.com/docs/api/overview#security). Here is a [step-by-step guide](http://fabioberger.com/blog/2014/11/06/building-a-coinbase-app-in-go/#env) on how to add these environment variables to your shell config file.

Now you can call methods on `c` similar to the ones described in the [API reference](https://coinbase.com/api/doc).  For example:

```go
balance, err := c.GetBalance()
if err != nil {
	log.Fatal(err)
}
fmt.Printf("Balance is %f BTC", balance)
```

A working API key example is available in example/ApiKeyExample.go. To run it, execute:

`go run ./example/ApiKeyExample.go`

## Error Handling

All errors generated at runtime will be returned to the calling client method. Any API request for which Coinbase returns an error encoded in a JSON response will be parsed and returned by the client method as a Golang error struct. Lastly, it is important to note that for HTTP requests, if the response code returned is not '200 OK', an error will be returned to the client method detailing the response code that was received.

## Examples

### Get user information

```go
user, err := c.GetUser()
if err != nil {
	log.Fatal(err)
}
fmt.Println(user.Name)
// 'User One'
fmt.Println(user.Email)
// 'user1@example.com'
```

### Check your balance

```go
amount, err := c.GetBalance()
if err != nil {
	log.Fatal(err)
}
fmt.Printf("Balance is %f BTC", amount)
// 'Balance is 24.229801 BTC'
```

### Send bitcoin

`func (c Client) SendMoney(params *TransactionParams) (*transactionConfirmation, error) `

```go
params := &coinbase.TransactionParams{
		To:     "1HHNtsSVWuJXzTZrAmq71busSKLHzgm4Wb",
		Amount: "0.0026",
		Notes:  "Thanks for the coffee!",
	}
confirmation, err := c.SendMoney(params)
if err != nil {
	log.Fatal(err)
}
fmt.Println(confirmation.Transaction.Status)
// 'pending'
fmt.Println(confirmation.Transaction.Id)
// '518d8567ed3ddcd4fd000034'
```

The "To" parameter can also be a bitcoin address and the "Notes" parameter can be a note or description of the transaction.  Descriptions are only visible on Coinbase (not on the general bitcoin network).

You can also send money in a number of currencies (see `GetCurrencies()`).  The amount will be automatically converted to the correct BTC amount using the current exchange rate.

All possible transaction parameters are detailed below:

```go
type TransactionParams struct {
	To                string
	From              string
	Amount            string
	AmountString      string
	AmountCurrencyIso string
	Notes             string
	UserFee           string
	ReferrerId        string
	Idem              string
	InstantBuy        bool
	OrderId           string
}
```
Note that parameters are equivalent to those of the coinbase API except in camelcase rather then with underscores between words (Golang standard). This can also be assumed for accessing return values. For detailed information on each parameter, check out the ['send_money' documentation](https://www.coinbase.com/api/doc/1.0/transactions/send_money.html)

### Request bitcoin

This will send an email to the recipient, requesting payment, and give them an easy way to pay.

```go
params := &coinbase.TransactionParams{
	From:   "client@example.com", //Who are you requesting Bitcoins from
	Amount: "2.5",
	Notes:  "contractor hours in January (website redesign for 50 BTC)",
}
confirmation, err := c.RequestMoney(params)
if err != nil {
	log.Fatal(err)
}
fmt.Println(confirmation.Transaction.Request)
// 'true'
fmt.Println(confirmation.Transaction.Id)
// '518d8567ed3ddcd4fd000034'


success, err := c.ResendRequest("501a3554f8182b2754000003")
if err != nil {
	log.Fatal(err)
}
fmt.Println(success)
// 'true'


success, err := c.CancelRequest("501a3554f8182b2754000003")
if err != nil {
	log.Fatal(err)
}
fmt.Println(success)
// 'true'


// From the other account:
success, err := c.CompleteRequest("501a3554f8182b2754000003")
if err != nil {
	log.Fatal(err)
}
fmt.Println(success)
// 'true'
```

### List your current transactions

Sorted in descending order by timestamp, 30 per page.  You can pass an integer as the first param to page through results, for example `c.GetTransactions(2)` would grab the second page.

```go
response, err := c.GetTransactions(1)
if err != nil {
	log.Fatal(err)
}
fmt.Println(response.CurrentPage)
// '1'
fmt.Println(response.NumPages)
// '2'
fmt.Println(response.Transactions[0].Id)
// '5018f833f8182b129c00002f'
```

Transactions will always have an `id` attribute which is the primary way to identity them through the Coinbase api.  They will also have a `hsh` (bitcoin hash) attribute once they've been broadcast to the network (usually within a few seconds).

### Check bitcoin prices

Check the buy or sell price by passing a `quantity` of bitcoin that you'd like to buy or sell.  The price can be given with or without the Coinbase's fee of 1% and the bank transfer fee of $0.15.

```go
price, err := c.GetBuyPrice(1)
if err != nil {
	log.Fatal(err)
}
fmt.Println(price.Subtotal.Amount) // Subtotal does not include fees
// '303.00'
fmt.Println(price.Total.Amount) // Total includes coinbase & bank fee
// '306.18'

price, err = c.GetSellPrice(1)
if err != nil {
	log.Fatal(err)
}
fmt.Println(price.Subtotal.Amount) // Subtotal is current market price
// '9.90'
fmt.Println(price.Total.Amount) // Total is amount you will receive (after fees)
// '9.65'
```

### Buy or sell bitcoin

Buying and selling bitcoin requires you to [link and verify a bank account](https://coinbase.com/payment_methods) through the web interface first.

Then you can call `buy` or `sell` and pass a `quantity` of bitcoin you want to buy.

On a buy, coinbase will debit your bank account and the bitcoin will arrive in your Coinbase account four business days later (this is shown as the `payoutDate` below).  This is how long it takes for the bank transfer to complete and verify, although they are working on shortening this window. In some cases, they may not be able to guarantee a price, and buy requests will fail. In that case, set the second parameter (`agreeBtcAmountVaries`) to true in order to purchase bitcoin at the future market price when your money arrives.

On a sell they will credit your bank account in a similar way and it will arrive within two business days.

	func (c Client) Buy(amount float64, agreeBtcAmountVaries bool) (*transfer, error)

```go
transfer, err := c.Buy(1.0, true)
if err != nil {
	log.Fatal(err)
}
fmt.Println(transfer.Code)
// '6H7GYLXZ'
fmt.Println(transfer.Btc.Amount)
// '1.00000000'
fmt.Println(transfer.Total.Amount)
// '$361.55'
fmt.Println(transfer.PayoutDate)
// '2013-02-01T18:00:00-08:00' (ISO 8601 format - can be parsed with time.Parse(transfer.PayoutDate, "2013-06-05T14:10:43.678Z"))
```

```go
transfer, err := c.Sell(1.0, true)
if err != nil {
	log.Fatal(err)
}
fmt.Println(transfer.Code)
// '6H7GYLXZ'
fmt.Println(transfer.Btc.Amount)
// '1.00000000'
fmt.Println(transfer.Total.Amount)
// '$361.55'
fmt.Println(transfer.PayoutDate)
// '2013-02-01T18:00:00-08:00' (ISO 8601 format - can be parsed with time.Parse(transfer.PayoutDate, "2013-06-05T14:10:43.678Z"))
```

### Create a payment button

This will create the code for a payment button (and modal window) that you can use to accept bitcoin on your website.  You can read [more about payment buttons here and try a demo](https://coinbase.com/docs/merchant_tools/payment_buttons).

The allowed ButtonParams are:

```go
type Button struct {
	Name                string
	PriceString         string
	PriceCurrencyIso    string
	Type                string
	Subscription        bool
	Repeat              string
	Style               string
	Text                string
	Description         string
	Custom              string
	CustomSecure        bool
	CallbackUrl         string
	SuccessUrl          string
	CancelUrl           string
	InfoUrl             string
	AutoRedirect        bool
	AutoRedirectSuccess bool
	AutoRedirectCancel  bool
	VariablePrice       bool
	ChoosePrice         bool
	IncludeAddress      bool
	IncludeEmail        bool
	Price1              string
	Price2              string
	Price3              string
	Price4              string
	Price5              string
}
```
The `custom` param will get passed through in [callbacks](https://coinbase.com/docs/merchant_tools/callbacks) to your site.  The list of valid `options` [are described here](https://coinbase.com/api/doc/1.0/buttons/create.html).

For detailed information on each parameter, check out the ['buttons' documentation](https://www.coinbase.com/api/doc/1.0/buttons/create.html)

```go
params := &coinbase.Button{
		Name:             "test",
		Type:             "buy_now",
		Subscription:     false,
		PriceString:      "1.23",
		PriceCurrencyIso: "USD",
		Custom:           "Order123",
		CallbackUrl:      "http://www.example.com/my_custom_button_callback",
		Description:      "Sample Description",
		Style:            "custom_large",
		IncludeEmail:     true,
	}
button, err := c.CreateButton(params)
if err != nil {
	log.Fatal(err)
}
fmt.Println(button.Code)
// '93865b9cae83706ae59220c013bc0afd'
fmt.Println(button.EmbedHtml)
// '<div class=\"coinbase-button\" data-code=\"93865b9cae83706ae59220c013bc0afd\"></div><script src=\"https://coinbase.com/assets/button.js\" type=\"text/javascript\"></script>'
```

### Exchange rates and currency utilities

You can fetch a list of all supported currencies and ISO codes with the `GetCurrencies()` method.

```go
currencies, err := c.GetCurrencies()
if err != nil {
	log.Fatal()
}
fmt.Println(currencies[0].Name)
// 'Afghan Afghani (AFN)'
```

`GetExchangeRates()` will return a list of exchange rates.

```go
exchanges, err := c.GetExchangeRates()
if err != nil {
	log.Fatal(err)
}
fmt.Println(exchanges["btc_to_cad"])
// '117.13892'
```

`GetExchangeRate(from string, to string)` will return a single exchange rate

```go
exchange, err := c.GetExchangeRate("btc", "usd")
if err != nil {
	log.Fatal(err)
}
fmt.Println(exchange)
// 117.13892
```

### Create a new user

```go
user, err := c.CreateUser("test@email.com", "password")
if err != nil {
	log.Fatal(err)
}
fmt.Println(user.Email)
// 'newuser@example.com'
fmt.Println(user.ReceiveAddress)
// 'mpJKwdmJKYjiyfNo26eRp4j6qGwuUUnw9x'
```

A receive address is returned also in case you need to send the new user a payment right away.

### Get autocomplete contacts

This will return a list of contacts the user has previously sent to or received from. Useful for auto completion. By default, 30 contacts are returned at a time; use the `$page` and `$limit` parameters to adjust how pagination works.

The allowed ContactsParams are:

```go
type ContactsParams struct {
	Page  int
	Limit int
	Query string
}
```

```go
params := &coinbase.ContactsParams{
	Page:  1,
	Limit: 5,
	Query: "user",
}
contacts, err := c.GetContacts(params)
if err != nil {
	log.Fatal(err)
}
fmt.Println(strings.Join(contacts.Emails, ","))
// 'user1@example.com, user2@example.com'
```

## Adding new methods

You can see a [list of method calls here](https://github.com/fabioberger/coinbase-go/blob/master/coinbase.go) and how they are implemented.  They are all wrappers around the [Coinbase JSON API](https://coinbase.com/api/doc).

If there are any methods listed in the [API Reference](https://coinbase.com/api/doc) that don't have an explicit function name in the library, you can also call `Get`, `Post`, `Put`, or `Delete` with a `path`, `params` and holder struct for a quick implementation. Holder should be a pointer to some data structure that correctly reflects the structure of the returned JSON response. The library will attempt to unmarshal the response from the server into holder. For example:

```go
balance := map[string]string{} // Holder struct depends on JSON format returned from API
if err := c.Get("account/balance", nil, &balance); err != nil {
	log.Fatal(err)
}
fmt.Println(balance)
// map[amount:36.62800000 currency:BTC]
```

Or feel free to add a new wrapper method and submit a pull request.

# OAuth Authentication

For an indepth tutorial on how to implement OAuth Authentication, visit this [step-by-step  tutorial](http://fabioberger.com/blog/2014/11/06/building-a-coinbase-app-in-go/#oauth).

## Higher Level Overview

To authenticate with OAuth, first create an OAuth application at [https://coinbase.com/oauth/applications](https://coinbase.com/oauth/applications).
When a user wishes to connect their Coinbase account, redirect them to a URL created with `func (o OAuth) CreateAuthorizeUrl(scope []string) string`:

```go
o, err := coinbase.OAuthService(YOUR_CLIENT_ID, YOUR_CLIENT_SECRET, YOUR_REDIRECT_URL)
if err != nil {
	log.Fatal(err)
}
scope := []string{"all",}
header("Location: " . o.CreateAuthorizeUrl(scope));
```

After the user has authorized your application, they will be redirected back to the redirect URL specified above. A `code` parameter will be included - pass this into `GetTokens` to receive a set of tokens:

```go
query := req.URL.Query()
code := query.Get("code")
tokens, err := o.GetTokens(code, "authorization_code")
if err != nil {
	log.Fatal(err)
}
```

Store these tokens safely, and use them to make Coinbase API requests in the future. For example:

```go
c := coinbase.OAuthClient(tokens)
amount, err := c.GetBalance()
if err != nil {
	log.Fatal(err)
}
```

A full example implementation is available in the `example` directory. In order to run this example implementation, you will need to install the following dependency:

```bash
go get github.com/go-martini/martini
```

You will also need to set your coinbase application client_id and client_secret as environment variables by adding these environment variables to your bash config file (i.e ~/.bashrc, ~/.bash_profile, etc...) and reload them:

```bash
export COINBASE_CLIENT_ID="YOUR_CLIENT_ID"
```
```bash
export COINBASE_CLIENT_SECRET="YOUR_CLIENT_SECRET"
```

```bash
source ~/.bash_profile
```

The last step we need to take is to generate a cert and key pair in order to run our OAuth server over SSL. To do this, run the following command from within the example directory:

```bash
go run $(go env GOROOT)/src/pkg/crypto/tls/generate_cert.go --host="localhost"
```

Once you have done this, run the example:

```bash
go run OAuthExample.go
```


## Security notes

If someone gains access to your API Key they will have complete control of your Coinbase account.  This includes the abillity to send all of your bitcoins elsewhere.

For this reason, API access is disabled on all Coinbase accounts by default.  If you decide to enable API key access you should take precautions to store your API key securely in your application.  How to do this is application specific, but it's something you should [research](http://programmers.stackexchange.com/questions/65601/is-it-smart-to-store-application-keys-ids-etc-directly-inside-an-application) if you have never done this before.

## Testing

In order to run the tests for this library, you will first need to install the Testify/Assert dependency with the following command:

 ```bash
 go get github.com/stretchr/testify/assert
 ```

Then run all tests by executing the following in your command line:

 	go test . -v


To run either only the endpoint or mock tests, use the below commands:

Endpoint(Live) :

	go test . -v -test.run=TestEndpoint

Mock :

	go test . -v -test.run=TestMock

If you would like to use the sandbox testnet instead of the live API endpoint, edit the "sandbox" variable in the config package to "true":

```
Sandbox = true
```


