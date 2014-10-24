# Coinbase Go Client Library

An easy way to buy, send, and accept [bitcoin](http://en.wikipedia.org/wiki/Bitcoin) through the [Coinbase API](https://coinbase.com/docs/api/overview).

This library supports both the [API key authentication method](https://coinbase.com/docs/api/overview) and OAuth. The below examples use an API key - for instructions on how to use OAuth, see [OAuth Authentication](#oauth-authentication).

## Installation

Make sure you have set the environment variable $GOPATH

	export GOPATH=path/to/your/go/folder

Obtain the latest version of the Coinbase Go library with:

    go get github.com/fabioberger/coinbase-go

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
fmt.Printf("Balance is %f BTC", balance)
```

A working API key example is available in example/ApiKeyExample.go.

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
fmt.Printf("Balance is %f BTC", balance)
// 'Balance is 24.229801 BTC'
```

### Send bitcoin

`func (c Client) SendMoney(params *TransactionRequestParams) (*transactionConfirmation, error) `

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

You can also send money in a number of currencies (see `getCurrencies()`).  The amount will be automatically converted to the correct BTC amount using the current exchange rate.

All possible transaction parameters are detailed below:

```go
type TransactionParams struct {
	To                  string 
	From                string 
	Amount              string 
	Amount_string       string 
	Amount_currency_iso string 
	Notes               string 
	User_fee            string 
	Referrer_id         string 
	Idem                string 
	Instant_buy         bool   
	Order_id            string 
}
```
For detailed information on each parameter, check out the ['send_money' documentation](https://www.coinbase.com/api/doc/1.0/transactions/send_money.html)

### Request bitcoin

This will send an email to the recipient, requesting payment, and give them an easy way to pay.

```go
params := &coinbase.TransactionParams{
	To:     "client@example.com",
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
```

### List your current transactions

Sorted in descending order by timestamp, 30 per page.  You can pass an integer as the first param to page through results, for example `c.GetTransactions(2)` would grab the second page.

```go
response, err := c.GetTransactions(1)
if err != nil {
	log.Fatal(err)
}
fmt.Println(response.Current_page)
// '1'
fmt.Println(response.Num_pages)
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

price, err := c.GetSellPrice(1)
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

On a buy, coinbase will debit your bank account and the bitcoin will arrive in your Coinbase account four business days later (this is shown as the `payout_date` below).  This is how long it takes for the bank transfer to complete and verify, although they are working on shortening this window. In some cases, they may not be able to guarantee a price, and buy requests will fail. In that case, set the second parameter (`agreeBtcAmountVaries`) to true in order to purchase bitcoin at the future market price when your money arrives.

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
fmt.Println(transfer.Payout_date)
// '2013-02-01T18:00:00-08:00' (ISO 8601 format - can be parsed with time.Parse(payout_date, "2013-06-05T14:10:43.678Z"))
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
fmt.Println(transfer.Payout_date)
// '2013-02-01T18:00:00-08:00' (ISO 8601 format - can be parsed with time.Parse(payout_date, "2013-06-05T14:10:43.678Z"))
```

### Create a payment button

This will create the code for a payment button (and modal window) that you can use to accept bitcoin on your website.  You can read [more about payment buttons here and try a demo](https://coinbase.com/docs/merchant_tools/payment_buttons).

The allowed ButtonParams are: 

```go
type button struct {
	Name                  string 
	Price_string          string 
	Price_currency_iso    string
	Type                  string 
	Subscription          bool   
	Repeat                string 
	Style                 string 
	Text                  string 
	Description           string 
	Custom                string 
	Custom_secure         bool   
	Callback_url          string 
	Success_url           string 
	Cancel_url            string 
	Info_url              string 
	Auto_redirect         bool   
	Auto_redirect_success bool   
	Auto_redirect_cancel  bool   
	Variable_price        bool   
	Choose_price          bool   
	Include_address       bool   
	Include_email         bool   
	Price1                string 
	Price2                string 
	Price3                string 
	Price4                string 
	Price5                string 
}
```
The `custom` param will get passed through in [callbacks](https://coinbase.com/docs/merchant_tools/callbacks) to your site.  The list of valid `options` [are described here](https://coinbase.com/api/doc/1.0/buttons/create.html).

For detailed information on each parameter, check out the ['buttons' documentation](https://www.coinbase.com/api/doc/1.0/buttons/create.html)

```go
params := &button{
	Name:               "test",
	Type:               "buy_now",
	Subscription:       false,
	Price_string:       "1.23",
	Price_currency_iso: "USD",
	Custom:             "Order123",
	Callback_url:       "http://www.example.com/my_custom_button_callback",
	Description:        "Sample Description",
	Style:              "custom_large",
	Include_email:      true,
}
button, err := c.CreateButton(params)
if err != nil {
	log.Fatal(err)
}
fmt.Println(button.Code)
// '93865b9cae83706ae59220c013bc0afd'
fmt.Println(button.Embed_html)
// '<div class=\"coinbase-button\" data-code=\"93865b9cae83706ae59220c013bc0afd\"></div><script src=\"https://coinbase.com/assets/button.js\" type=\"text/javascript\"></script>'
```

### Exchange rates and currency utilties

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
data, err := c.GetExchangeRate("btc", "usd")
if err != nil {
	log.Fatal(err)
}
// '117.13892'
```

### Create a new user

```go
$response = $coinbase->createUser("newuser@example.com", "some password");
echo $response->user->email;
// 'newuser@example.com'
echo $response->user->receive_address;
// 'mpJKwdmJKYjiyfNo26eRp4j6qGwuUUnw9x'
```

A receive address is returned also in case you need to send the new user a payment right away.

### Get autocomplete contacts

This will return a list of contacts the user has previously sent to or received from. Useful for auto completion. By default, 30 contacts are returned at a time; use the `$page` and `$limit` parameters to adjust how pagination works.

```go
$response = $coinbase->getContacts("exa");
echo implode(', ', $response->contacts);
// 'user1@example.com, user2@example.com'
```

## Adding new methods

You can see a [list of method calls here](https://github.com/coinbase/coinbase-go/blob/master/lib/Coinbase/Coinbase.go) and how they are implemented.  They are a wrapper around the [Coinbase JSON API](https://coinbase.com/api/doc).

If there are any methods listed in the [API Reference](https://coinbase.com/api/doc) that don't have an explicit function name in the library, you can also call `get`, `post`, `put`, or `delete` with a `$path` and optional `$params` array for a quick implementation.  The raw JSON object will be returned. For example:

```go
var_dump($coinbase->get('/account/balance'));
// object(stdClass)#4 (2) {
//   ["amount"]=>
//   string(10) "0.56902981"
//   ["currency"]=>
//   string(3) "BTC"
// }
```

Or feel free to add a new wrapper method and submit a pull request.
