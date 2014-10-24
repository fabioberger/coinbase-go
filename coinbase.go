// Coinbase-go is a convenient Go wrapper for the Coinbase API
package coinbase

import (
	"errors"
	"fmt"
	"strconv"
)

// The Client from which all API requests are made
type Client struct {
	rpc Rpc
}

// Instantiate the client with ApiKey Authentication
func ApiKeyClient(key string, secret string) Client {
	c := Client{
		rpc: Rpc{
			auth: ApiKeyAuth(key, secret),
			mock: false,
		},
	}
	return c
}

// Instantiate the client with OAuth Authentication
func OAuthClient(tokens *oauthTokens) Client {
	c := Client{
		rpc: Rpc{
			auth: ClientOAuth(tokens),
			mock: false,
		},
	}
	return c
}

// Instantiates Testing ApiKeyClient. All client methods execute normally except
// responses are returned from a test_data/ file instead of the coinbase API
func ApiKeyClientTest(key string, secret string) Client {
	c := ApiKeyClient(key, secret)
	c.rpc.mock = true
	return c
}

// Sends a GET request and marshals response data into holder
func (c Client) Get(path string, params interface{}, holder interface{}) error {
	return c.rpc.Request("GET", path, params, &holder)
}

// Sends a POST request and marshals response data into holder
func (c Client) Post(path string, params interface{}, holder interface{}) error {
	return c.rpc.Request("POST", path, params, &holder)
}

// Sends a DELETE request and marshals response data into holder
func (c Client) Delete(path string, params interface{}, holder interface{}) error {
	return c.rpc.Request("DELETE", path, params, &holder)
}

// Sends a PUT request and marshals response data into holder
func (c Client) Put(path string, params interface{}, holder interface{}) error {
	return c.rpc.Request("PUT", path, params, &holder)
}

// Returns current balance in BTC
func (c Client) GetBalance() (float64, error) {
	balance := map[string]string{}
	if err := c.Get("account/balance", nil, &balance); err != nil {
		return 0.0, err
	}
	fmt.Println(balance)
	balanceFloat, err := strconv.ParseFloat(balance["amount"], 64)
	if err != nil {
		return 0, err
	}
	return balanceFloat, nil
}

// Returns clients current bitcoin receive address
func (c Client) GetReceiveAddress() (string, error) {
	holder := map[string]interface{}{}
	if err := c.Get("account/receive_address", nil, &holder); err != nil {
		return "", err
	}
	return holder["address"].(string), nil
}

// Returns bitcoin addresses associated with client account
func (c Client) GetAllAddresses(params *AddressesParams) (*addresses, error) {
	holder := addressesHolder{}
	if err := c.Get("addresses", params, &holder); err != nil {
		return nil, err
	}
	addresses := addresses{
		Total_count:  holder.Total_count,
		Num_pages:    holder.Current_page,
		Current_page: holder.Current_page,
	}
	// Remove one layer of nesting
	for _, addr := range holder.Addresses {
		addresses.Addresses = append(addresses.Addresses, addr.Address)
	}
	return &addresses, nil
}

// Generates and returns a new bitcoin receive address
func (c Client) GenerateReceiveAddress(params *ReceiveAddressParams) (string, error) {
	holder := map[string]interface{}{}
	if err := c.Post("account/generate_receive_address", params, &holder); err != nil {
		return "", err
	}
	return holder["address"].(string), nil
}

// Sends money to either a bitcoin or email address
func (c Client) SendMoney(params *TransactionParams) (*transactionConfirmation, error) {
	return c.transactionRequest("POST", "send_money", params)
}

// Request money from either a bitcoin or email address
func (c Client) RequestMoney(params *TransactionParams) (*transactionConfirmation, error) {
	return c.transactionRequest("POST", "request_money", params)
}

// Execute a transaction request (i.e send_money, request_money)
func (c Client) transactionRequest(method string, kind string, params *TransactionParams) (*transactionConfirmation, error) {
	finalParams := &TransactionRequestParams{
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
	if err := checkApiErrors(holder.Response); err != nil {
		return nil, err
	}
	confirmation := transactionConfirmation{
		Transaction: holder.Transaction,
		Transfer:    holder.Transfer,
	}
	return &confirmation, nil
}

// Resend a transaction request referenced by id
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

// Cancel a transaction request referenced by id
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

// Complete a money request referenced by id
func (c Client) CompleteRequest(id string) (*transactionConfirmation, error) {
	return c.transactionRequest("PUT", id+"/complete_request", nil)
}

// Get a new payment button including Embed_html as a field on button struct
func (c Client) CreateButton(params *button) (*button, error) {
	finalParams := &ButtonParams{
		Button: params,
	}
	holder := buttonHolder{}
	if err := c.Post("buttons", finalParams, &holder); err != nil {
		return nil, err
	}
	if err := checkApiErrors(holder.Response); err != nil {
		return nil, err
	}
	button := holder.Button
	button.Embed_html = "<div class=\"coinbase-button\" data-code=\"" + button.Code + "\"></div><script src=\"https://coinbase.com/assets/button.js\" type=\"text/javascript\"></script>"
	return &button, nil
}

// Create an order for a given button code
func (c Client) CreateOrderFromButtonCode(buttonCode string) (*order, error) {
	holder := orderHolder{}
	if err := c.Post("buttons/"+buttonCode+"/create_order", nil, &holder); err != nil {
		return nil, err
	}
	if err := checkApiErrors(holder.Response); err != nil {
		return nil, err
	}
	return &holder.Order, nil
}

// Create a new user given an email and password
func (c Client) CreateUser(email string, password string) (*user, error) {
	params := map[string]interface{}{
		"user[email]":    email,
		"user[password]": password,
	}
	holder := userHolder{}
	if err := c.Post("users", params, &holder); err != nil {
		return nil, err
	}
	if err := checkApiErrors(holder.Response); err != nil {
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
	if err := checkApiErrors(holder.Response); err != nil {
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
	if err := checkApiErrors(holder.Response); err != nil {
		return nil, err
	}
	return &holder.Transfer, nil
}

// Get a users contacts
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

// Get all currency names and ISO's
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

// Get the current exchange rates
func (c Client) GetExchangeRates() (map[string]string, error) {
	holder := map[string]string{}
	if err := c.Get("currencies/exchange_rates", nil, &holder); err != nil {
		return nil, err
	}
	return holder, nil
}

// Get the exchange rate between two specified currencies
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

// Get transactions associated with an account
func (c Client) GetTransactions(page int) (*transactions, error) {
	params := map[string]int{
		"page": page,
	}
	holder := transactionsHolder{}
	if err := c.Get("transactions", params, &holder); err != nil {
		return nil, err
	}
	transactions := transactions{
		PaginationStats: holder.PaginationStats,
	}
	// Remove one layer of nesting
	for _, tx := range holder.Transactions {
		transactions.Transactions = append(transactions.Transactions, tx.Transaction)
	}
	return &transactions, nil
}

// Get orders associated with an account
func (c Client) GetOrders(page int) (*orders, error) {
	holder := ordersHolder{}
	params := map[string]int{
		"page": page,
	}
	if err := c.Get("orders", params, &holder); err != nil {
		return nil, err
	}
	orders := orders{
		PaginationStats: holder.PaginationStats,
	}
	// Remove one layer of nesting
	for _, o := range holder.Orders {
		orders.Orders = append(orders.Orders, o.Order)
	}
	return &orders, nil
}

// Get transfers associated with an account
func (c Client) GetTransfers(page int) (*transfers, error) {
	params := map[string]int{
		"page": page,
	}
	holder := transfersHolder{}
	if err := c.Get("transfers", params, &holder); err != nil {
		return nil, err
	}
	transfers := transfers{
		PaginationStats: holder.PaginationStats,
	}
	// Remove one layer of nesting
	for _, t := range holder.Transfers {
		transfers.Transfers = append(transfers.Transfers, t.Transfer)
	}
	return &transfers, nil
}

// Get the current BTC buy price
func (c Client) GetBuyPrice(qty int) (*pricesHolder, error) {
	return c.getPrice("buy", qty)
}

// Get the current BTC sell price
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

// Get a particular transaction referenced by id
func (c Client) GetTransaction(id string) (*transaction, error) {
	holder := transactionHolder{}
	if err := c.Get("transactions/"+id, nil, &holder); err != nil {
		return nil, err
	}
	if err := checkApiErrors(holder.Response); err != nil {
		return nil, err
	}
	return &holder.Transaction, nil
}

// Get a particular order referenced by id
func (c Client) GetOrder(id string) (*order, error) {
	holder := orderHolder{}
	if err := c.Get("orders/"+id, nil, &holder); err != nil {
		return nil, err
	}
	if err := checkApiErrors(holder.Response); err != nil {
		return nil, err
	}
	return &holder.Order, nil
}

// Get the user associated with the authentication
func (c Client) GetUser() (*user, error) {
	holder := usersHolder{}
	if err := c.Get("users", nil, &holder); err != nil {
		return nil, err
	}
	return &holder.Users[0].User, nil
}
