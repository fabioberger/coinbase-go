package coinbase

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

// Initialize the client without mock mode enabled on rpc.
// All calls hit the coinbase API and tests focus on checking
// format of the response and validity of sent requests
func initClient() Client {
	return ApiKeyClient(os.Getenv("COINBASE_KEY"), os.Getenv("COINBASE_SECRET"))
}

func TestGetBalanceEndpoint(t *testing.T) {
	c := initClient()
	amount, err := c.GetBalance()
	if err != nil {
		log.Fatal(err)
	}
	assert.IsType(t, "string", amount)
}

func TestGetReceiveAddressEndpoint(t *testing.T) {
	t.Skip("Skipping GetReceiveAddressEndpoint")
	fmt.Println("Skipped to avoid generating many addresses")
	c := initClient()
	params := &ReceiveAddressParams{
		Address: &AddressParams{
			Callback_url: "http://www.wealthlift.com",
			Label:        "My Test Address",
		},
	}
	address, err := c.GenerateReceiveAddress(params)
	if err != nil {
		log.Fatal(err)
	}
	assert.IsType(t, "string", address)
}

func TestGetAllAddressesEndpoint(t *testing.T) {
	c := initClient()
	params := &AddressesParams{
		Page:  1,
		Limit: 5,
	}
	addresses, err := c.GetAllAddresses(params)
	if err != nil {
		log.Fatal(err)
	}
	assert.IsType(t, "string", addresses.Addresses[0].Address.Created_at)
	assert.IsType(t, "string", addresses.Addresses[1].Address.Address)
}

func TestGetButton(t *testing.T) {
	c := initClient()
	params := &ButtonParams{
		Button: &button{
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
		},
	}
	data, err := c.GetButton(params)
	if err != nil {
		log.Fatal(err)
	}
	assert.IsType(t, true, data.Success)
	assert.IsType(t, "string", data.Button.Type)
}

func TestGetExchangeRates(t *testing.T) {
	c := initClient()
	data, err := c.GetExchangeRates()
	if err != nil {
		log.Fatal(err)
	}
	assert.IsType(t, "string", data["btc_to_usd"])
}

func TestGetExchangeRate(t *testing.T) {
	c := initClient()
	data, err := c.GetExchangeRate("btc", "usd")
	if err != nil {
		log.Fatal(err)
	}
	assert.IsType(t, "string", data)
}

func TestGetTransactions(t *testing.T) {
	c := initClient()
	data, err := c.GetTransactions(1)
	if err != nil {
		log.Fatal(err)
	}
	assert.IsType(t, "string", data.Current_user.Id)
	assert.IsType(t, "string", data.Native_balance.Amount)
}

func TestGetBuyPrice(t *testing.T) {
	c := initClient()
	data, err := c.GetBuyPrice(1)
	if err != nil {
		log.Fatal(err)
	}
	assert.IsType(t, "string", data.Subtotal.Currency)
	assert.IsType(t, "string", data.Total.Amount)
}

func TestGetSellPrice(t *testing.T) {
	c := initClient()
	data, err := c.GetSellPrice(1)
	if err != nil {
		log.Fatal(err)
	}
	assert.IsType(t, "string", data.Subtotal.Currency)
	assert.IsType(t, "string", data.Total.Amount)
}
