// Coinbase-go is a convenient Go wrapper for the Coinbase API
package coinbase

import (
	"errors"
	"net/http"
)

const API_BASE = "https://api.coinbase.com/v1/" //"https://www.bitstamp.net/"

// The Client from which all API requests are made
type Client struct {
	rpc Rpc
}

// Instantiate the client with ApiKey Authentication
func ApiKeyClient(key string, secret string) Client {
	c := Client{
		rpc: Rpc{
			client: http.Client{
				Transport: &http.Transport{
					Dial: dialTimeout,
				},
			},
			auth: ApiKeyAuthentication{
				Key:    key,
				Secret: secret,
			},
			mock: false,
		},
	}
	return c
}

// Only used when testing the wrapper
func ApiKeyClientTest(key string, secret string) Client {
	c := ApiKeyClient(key, secret)
	c.rpc.mock = true
	return c
}

func (c Client) Get(path string, params interface{}, holder interface{}) error {
	return c.rpc.Request("GET", path, params, &holder)
}

func (c Client) Post(path string, params interface{}, holder interface{}) error {
	return c.rpc.Request("POST", path, params, &holder)
}

func (c Client) Delete(path string, params interface{}, holder interface{}) error {
	return c.rpc.Request("DELETE", path, params, &holder)
}

func (c Client) Put(path string, params interface{}, holder interface{}) error {
	return c.rpc.Request("PUT", path, params, &holder)
}

func (c Client) getPaginatedResource(resource string, listElement string, unwrapElement string, params interface{}, holder interface{}) error {
	if err := c.Get(resource, params, &holder); err != nil {
		return err
	}
	return nil
}

func (c Client) GetBalance() (string, error) {
	balance := map[string]string{}
	if err := c.Get("account/balance", nil, &balance); err != nil {
		return "", err
	}
	return balance["amount"], nil
}

func (c Client) GetReceiveAddress() (string, error) {
	holder := map[string]interface{}{}
	if err := c.Get("account/receive_address", nil, &holder); err != nil {
		return "", err
	}
	return holder["address"].(string), nil
}

func (c Client) GetAllAddresses(params *AddressesParams) (*addressesHolder, error) {
	holder := addressesHolder{}
	if err := c.getPaginatedResource("addresses", "addresses", "address", params, &holder); err != nil {
		return nil, err
	}
	return &holder, nil
}

func (c Client) GenerateReceiveAddress(params *ReceiveAddressParams) (string, error) {
	holder := map[string]interface{}{}
	if err := c.Post("account/generate_receive_address", params, &holder); err != nil {
		return "", err
	}
	return holder["address"].(string), nil
}

func (c Client) SendMoney(params *TransactionRequestParams) (*transactionHolder, error) {
	holder := transactionHolder{}
	if err := c.Post("transactions/send_money", params, &holder); err != nil {
		return nil, err
	}
	return &holder, nil
}

func (c Client) RequestMoney(params *TransactionRequestParams) (*transactionHolder, error) {
	holder := transactionHolder{}
	if err := c.Post("transactions/request_money", params, &holder); err != nil {
		return nil, err
	}
	return &holder, nil
}

func (c Client) ResendRequest(id string) (bool, error) {
	holder := map[string]interface{}{}
	if err := c.Put("transactions/"+id+"/resend_request", nil, &holder); err != nil {
		return false, err
	}
	if holder["success"].(bool) {
		return true, nil
	}
	return false, nil
}

func (c Client) CancelRequest(id string) (bool, error) {
	holder := map[string]interface{}{}
	if err := c.Delete("transactions/"+id+"/cancel_request", nil, &holder); err != nil {
		return false, err
	}
	if holder["success"].(bool) {
		return true, nil
	}
	return false, nil
}

func (c Client) CompleteRequest(id string) (*transactionHolder, error) {
	holder := transactionHolder{}
	if err := c.Put("transactions/"+id+"/complete_request", nil, &holder); err != nil {
		return nil, err
	}
	return &holder, nil
}

func (c Client) GetButton(params *ButtonParams) (*buttonHolder, error) {
	holder := buttonHolder{}
	if err := c.Post("buttons", params, &holder); err != nil {
		return nil, err
	}
	holder.Embed_html = "<div class=\"coinbase-button\" data-code=\"" + holder.Button.Code + "\"></div><script src=\"https://coinbase.com/assets/button.js\" type=\"text/javascript\"></script>"
	return &holder, nil
}

func (c Client) CreateOrderFromButtonCode(buttonCode string) (*orderHolder, error) {
	holder := orderHolder{}
	if err := c.Post("buttons/"+buttonCode+"/create_order", nil, &holder); err != nil {
		return nil, err
	}
	return &holder, nil
}

func (c Client) CreateUser(email string, password string) (*userHolder, error) {
	params := map[string]interface{}{
		"user[email]":    email,
		"user[password]": password,
	}
	holder := userHolder{}
	if err := c.Post("users", params, &holder); err != nil {
		return nil, err
	}
	return &holder, nil
}

func (c Client) Buy(amount string, agreeBtcAmountVaries bool) (*transferHolder, error) {
	params := map[string]interface{}{
		"qty": amount,
		"agree_btc_amount_varies": agreeBtcAmountVaries,
	}
	holder := transferHolder{}
	if err := c.Post("buys", params, &holder); err != nil {
		return nil, err
	}
	return &holder, nil
}

func (c Client) Sell(amount string) (*transferHolder, error) {
	params := map[string]interface{}{
		"qty": amount,
	}
	holder := transferHolder{}
	if err := c.Post("sells", params, &holder); err != nil {
		return nil, err
	}
	return &holder, nil
}

func (c Client) GetContacts(params *ContactsParams) (*contactsHolder, error) {
	holder := contactsHolder{}
	if err := c.Get("contacts", params, &holder); err != nil {
		return nil, err
	}
	for _, contact := range holder.Contacts {
		if contact.Contact.Email != "" {
			holder.Emails = append(holder.Emails, contact.Contact.Email)
		}
	}
	return &holder, nil
}

func (c Client) GetCurrencies() ([]currency, error) {
	holder := [][]string{}
	if err := c.Get("currencies", nil, &holder); err != nil {
		return nil, err
	}
	finalData := []currency{}
	for _, curr := range holder {
		class := currency{
			Name: curr[0],
			Iso:  curr[1],
		}
		finalData = append(finalData, class)
	}
	return finalData, nil
}

func (c Client) GetExchangeRates() (map[string]string, error) {
	holder := map[string]string{}
	if err := c.Get("currencies/exchange_rates", nil, &holder); err != nil {
		return nil, err
	}
	return holder, nil
}

func (c Client) GetExchangeRate(from string, to string) (string, error) {
	exchanges, err := c.GetExchangeRates()
	if err != nil {
		return "", err
	}
	key := from + "_to_" + to
	if exchanges[key] == "" {
		return "", errors.New("The exchange rate does not exist for this currency pair")
	}
	return exchanges[key], nil
}

func (c Client) GetTransactions(page int) (*transactionsHolder, error) {
	params := map[string]int{
		"page": page,
	}
	holder := transactionsHolder{}
	err := c.getPaginatedResource("transactions", "transactions", "transaction", params, &holder)
	if err != nil {
		return nil, err
	}
	return &holder, nil
}

func (c Client) GetOrders(page int) (*ordersHolder, error) {
	holder := ordersHolder{}
	params := map[string]int{
		"page": page,
	}
	if err := c.getPaginatedResource("orders", "orders", "order", params, &holder); err != nil {
		return nil, err
	}
	return &holder, nil
}

func (c Client) GetTransfers(page int) (*transfersHolder, error) {
	params := map[string]int{
		"page": page,
	}
	holder := transfersHolder{}
	if err := c.getPaginatedResource("transfers", "transfers", "transfer", params, &holder); err != nil {
		return nil, err
	}
	return &holder, nil
}

func (c Client) GetBuyPrice(qty int) (*pricesHolder, error) {
	return c.getPrice("buy", qty)
}

func (c Client) GetSellPrice(qty int) (*pricesHolder, error) {
	return c.getPrice("sell", qty)
}

func (c Client) getPrice(kind string, qty int) (*pricesHolder, error) {
	params := map[string]int{
		"qty": qty,
	}
	holder := pricesHolder{}
	if err := c.Get("prices/"+kind, params, &holder); err != nil {
		return nil, err
	}
	return &holder, nil
}

func (c Client) GetTransaction(id string) (*transaction, error) {
	holder := transactionHolder{}
	if err := c.Get("transactions/"+id, nil, &holder); err != nil {
		return nil, err
	}
	return &holder.Transaction, nil
}

func (c Client) GetOrder(id string) (*order, error) {
	holder := orderHolder{}
	if err := c.Get("orders/"+id, nil, &holder); err != nil {
		return nil, err
	}
	return &holder.Order, nil
}

func (c Client) GetUser() (*user, error) {
	holder := usersHolder{}
	if err := c.Get("users", nil, &holder); err != nil {
		return nil, err
	}
	return &holder.Users[0].User, nil
}
