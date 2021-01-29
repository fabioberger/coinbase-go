// Coinbase-go is a convenient Go wrapper for the Coinbase API
package coinbase

import (
	"errors"
	"strconv"
	"strings"
)

// Client is the struct from which all API requests are made
type Client struct {
	rpc rpc
}

// ApiKeyClient instantiates the client with ApiKey Authentication
func ApiKeyClient(key string, secret string) Client {
	c := Client{
		rpc: rpc{
			auth: apiKeyAuth(key, secret),
			mock: false,
		},
	}
	return c
}

// OAuthClient instantiates the client with OAuth Authentication
func OAuthClient(tokens *OauthTokens) Client {
	c := Client{
		rpc: rpc{
			auth: clientOAuth(tokens),
			mock: false,
		},
	}
	return c
}

// ApiKeyClientTest instantiates Testing ApiKeyClient. All client methods execute
// normally except responses are returned from a test_data/ file instead of the coinbase API
func apiKeyClientTest(key string, secret string) Client {
	c := ApiKeyClient(key, secret)
	c.rpc.mock = true
	return c
}

// Get sends a GET request and marshals response data into holder
func (c Client) Get(path string, params interface{}, holder interface{}) error {
	return c.rpc.Request("GET", path, params, &holder)
}

// Post sends a POST request and marshals response data into holder
func (c Client) Post(path string, params interface{}, holder interface{}) error {
	return c.rpc.Request("POST", path, params, &holder)
}

// Delete sends a DELETE request and marshals response data into holder
func (c Client) Delete(path string, params interface{}, holder interface{}) error {
	return c.rpc.Request("DELETE", path, params, &holder)
}

// Put sends a PUT request and marshals response data into holder
func (c Client) Put(path string, params interface{}, holder interface{}) error {
	return c.rpc.Request("PUT", path, params, &holder)
}

// GetBalance returns current balance in BTC
func (c Client) GetBalance() (float64, error) {
	balance := map[string]string{}
	if err := c.Get("account/balance", nil, &balance); err != nil {
		return 0.0, err
	}
	balanceFloat, err := strconv.ParseFloat(balance["amount"], 64)
	if err != nil {
		return 0, err
	}
	return balanceFloat, nil
}

// GetReceiveAddress returns clients current bitcoin receive address
func (c Client) GetReceiveAddress() (string, error) {
	holder := map[string]interface{}{}
	if err := c.Get("account/receive_address", nil, &holder); err != nil {
		return "", err
	}
	return holder["address"].(string), nil
}

// GetAllAddresses returns bitcoin addresses associated with client account
func (c Client) GetAllAddresses(params *AddressesParams) (*addresses, error) {
	holder := addressesHolder{}
	if err := c.Get("addresses", params, &holder); err != nil {
		return nil, err
	}
	addresses := addresses{
		paginationStats: holder.paginationStats,
	}
	// Remove one layer of nesting
	for _, addr := range holder.Addresses {
		addresses.Addresses = append(addresses.Addresses, addr.Address)
	}
	return &addresses, nil
}

// GenerateReceiveAddress generates and returns a new bitcoin receive address
func (c Client) GenerateReceiveAddress(params *AddressParams) (string, error) {
	holder := map[string]interface{}{}
	if err := c.Post("account/generate_receive_address", params, &holder); err != nil {
		return "", err
	}
	return holder["address"].(string), nil
}

// SendMoney to either a bitcoin or email address
func (c Client) SendMoney(params *TransactionParams) (*transactionConfirmation, error) {
	return c.transactionRequest("POST", "send_money", params)
}

// RequestMoney from either a bitcoin or email address
func (c Client) RequestMoney(params *TransactionParams) (*transactionConfirmation, error) {
	return c.transactionRequest("POST", "request_money", params)
}

func (c Client) transactionRequest(method string, kind string, params *TransactionParams) (*transactionConfirmation, error) {
	finalParams := &struct {
		Transaction *TransactionParams `json:"transaction"`
	}{
		Transaction: params,
	}
	holder := transactionHolder{}
	var err error
	if method == "POST" {
		err = c.Post("transactions/"+kind, finalParams, &holder)
	} else if method == "PUT" {
		err = c.Put("transactions/"+kind, finalParams, &holder)
	}
	if err != nil {
		return nil, err
	}
	if err := checkApiErrors(holder.response, kind); err != nil {
		return nil, err
	}
	confirmation := transactionConfirmation{
		Transaction: holder.Transaction,
		Transfer:    holder.Transfer,
	}
	return &confirmation, nil
}

// ResendRequest resends a transaction request referenced by id
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

// CancelRequest cancels a transaction request referenced by id
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

// CompleteRequest completes a money request referenced by id
func (c Client) CompleteRequest(id string) (*transactionConfirmation, error) {
	return c.transactionRequest("PUT", id+"/complete_request", nil)
}

// CreateButton gets a new payment button including EmbedHtml as a field on button struct
func (c Client) CreateButton(params *Button) (*Button, error) {
	finalParams := &struct {
		Button *Button `json:"button"`
	}{
		Button: params,
	}
	holder := buttonHolder{}
	if err := c.Post("buttons", finalParams, &holder); err != nil {
		return nil, err
	}
	if err := checkApiErrors(holder.response, "CreateButton"); err != nil {
		return nil, err
	}
	button := holder.Button
	button.EmbedHtml = "<div class=\"coinbase-button\" data-code=\"" + button.Code + "\"></div><script src=\"https://coinbase.com/assets/button.js\" type=\"text/javascript\"></script>"
	return &button, nil
}

// CreateOrderFromButtonCode creates an order for a given button code
func (c Client) CreateOrderFromButtonCode(buttonCode string) (*order, error) {
	holder := orderHolder{}
	if err := c.Post("buttons/"+buttonCode+"/create_order", nil, &holder); err != nil {
		return nil, err
	}
	if err := checkApiErrors(holder.response, "CreateOrderFromButtonCode"); err != nil {
		return nil, err
	}
	return &holder.Order, nil
}

// CreateUser creates a new user given an email and password
func (c Client) CreateUser(email string, password string) (*user, error) {
	params := map[string]interface{}{
		"user[email]":    email,
		"user[password]": password,
	}
	holder := userHolder{}
	if err := c.Post("users", params, &holder); err != nil {
		return nil, err
	}
	if err := checkApiErrors(holder.response, "CreateUser"); err != nil {
		return nil, err
	}
	return &holder.User, nil
}

// Buy an amount of BTC and bypass rate limits by setting agreeBtcAmountVaries to true
func (c Client) Buy(amount float64, agreeBtcAmountVaries bool) (*transfer, error) {
	params := map[string]interface{}{
		"qty": amount,
		"agree_btc_amount_varies": agreeBtcAmountVaries,
	}
	holder := transferHolder{}
	if err := c.Post("buys", params, &holder); err != nil {
		return nil, err
	}
	if err := checkApiErrors(holder.response, "Buy"); err != nil {
		return nil, err
	}
	return &holder.Transfer, nil
}

// Sell an amount of BTC
func (c Client) Sell(amount float64) (*transfer, error) {
	params := map[string]interface{}{
		"qty": amount,
	}
	holder := transferHolder{}
	if err := c.Post("sells", params, &holder); err != nil {
		return nil, err
	}
	if err := checkApiErrors(holder.response, "Sell"); err != nil {
		return nil, err
	}
	return &holder.Transfer, nil
}

// GetContacts gets a users contacts
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

// GetCurrencies gets all currency names and ISO's
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

// GetExchangeRates gets the current exchange rates
func (c Client) GetExchangeRates() (map[string]string, error) {
	holder := map[string]string{}
	if err := c.Get("currencies/exchange_rates", nil, &holder); err != nil {
		return nil, err
	}
	return holder, nil
}

// GetExchangeRate gets the exchange rate between two specified currencies
func (c Client) GetExchangeRate(from string, to string) (float64, error) {
	exchanges, err := c.GetExchangeRates()
	if err != nil {
		return 0.0, err
	}
	key := from + "_to_" + to
	if exchanges[key] == "" {
		return 0.0, errors.New("The exchange rate does not exist for this currency pair")
	}
	exchangeFloat, err := strconv.ParseFloat(exchanges[key], 64)
	if err != nil {
		return 0.0, err
	}
	return exchangeFloat, nil
}

// GetTransactions gets transactions associated with an account
func (c Client) GetTransactions(page int) (*transactions, error) {
	params := map[string]int{
		"page": page,
	}
	holder := transactionsHolder{}
	if err := c.Get("transactions", params, &holder); err != nil {
		return nil, err
	}
	transactions := transactions{
		paginationStats: holder.paginationStats,
	}
	// Remove one layer of nesting
	for _, tx := range holder.Transactions {
		transactions.Transactions = append(transactions.Transactions, tx.Transaction)
	}
	return &transactions, nil
}

// GetOrders gets orders associated with an account
func (c Client) GetOrders(page int) (*orders, error) {
	holder := ordersHolder{}
	params := map[string]int{
		"page": page,
	}
	if err := c.Get("orders", params, &holder); err != nil {
		return nil, err
	}
	orders := orders{
		paginationStats: holder.paginationStats,
	}
	// Remove one layer of nesting
	for _, o := range holder.Orders {
		orders.Orders = append(orders.Orders, o.Order)
	}
	return &orders, nil
}

// GetTransfers get transfers associated with an account
func (c Client) GetTransfers(page int) (*transfers, error) {
	params := map[string]int{
		"page": page,
	}
	holder := transfersHolder{}
	if err := c.Get("transfers", params, &holder); err != nil {
		return nil, err
	}
	transfers := transfers{
		paginationStats: holder.paginationStats,
	}
	// Remove one layer of nesting
	for _, t := range holder.Transfers {
		transfers.Transfers = append(transfers.Transfers, t.Transfer)
	}
	return &transfers, nil
}

// GetBuyPrice gets the current BTC buy price
func (c Client) GetBuyPrice(qty int) (*pricesHolder, error) {
	return c.getPrice("buy", qty)
}

// GetSellPrice gets the current BTC sell price
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

// GetTransaction gets a particular transaction referenced by id
func (c Client) GetTransaction(id string) (*transaction, error) {
	holder := transactionHolder{}
	if err := c.Get("transactions/"+id, nil, &holder); err != nil {
		return nil, err
	}
	if err := checkApiErrors(holder.response, "GetTransaction"); err != nil {
		return nil, err
	}
	return &holder.Transaction, nil
}

// GetOrder gets a particular order referenced by id
func (c Client) GetOrder(id string) (*order, error) {
	holder := orderHolder{}
	if err := c.Get("orders/"+id, nil, &holder); err != nil {
		return nil, err
	}
	if err := checkApiErrors(holder.response, "GetOrder"); err != nil {
		return nil, err
	}
	return &holder.Order, nil
}

// GetUser gets the user associated with the authentication
func (c Client) GetUser() (*user, error) {
	holder := usersHolder{}
	if err := c.Get("users", nil, &holder); err != nil {
		return nil, err
	}
	return &holder.Users[0].User, nil
}

// checkApiErrors checks for errors returned by coinbase API JSON response
// i.e { "success": false, "errors": ["Button with code code123456 does not exist"], ...}
func checkApiErrors(resp response, method string) error {
	if resp.Success == false { // Return errors received from API here
		err := " in " + method + "()"
		if resp.Errors != nil {
			err = strings.Join(resp.Errors, ",") + err
			return errors.New(err)
		}
		if resp.Error != "" { // Return errors received from API here
			err = resp.Error + err
			return errors.New(err)
		}
	}
	return nil
}
