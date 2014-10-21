package coinbase

import (
	"fmt"
	"log"
	"os"
	"testing"
)

// Initialize the client with mock mode enabled on rpc
// All calls return the corresponding json response from the test_data files
func initTestClient() Client {
	return ApiKeyClientTest(os.Getenv("COINBASE_KEY"), os.Getenv("COINBASE_SECRET"))
}

func TestGetBalanceParse(t *testing.T) {
	c := initTestClient()
	amount, err := c.GetBalance()
	if err != nil {
		log.Fatal(err)
	}
	compareString(t, "GetBalanceParse", "36.62800000", amount)
}

func TestGetReceiveAddressParse(t *testing.T) {
	c := initTestClient()
	params := &ReceiveAddressParams{
		Address: &AddressParams{},
	}
	address, err := c.GenerateReceiveAddress(params)
	if err != nil {
		log.Fatal(err)
	}
	compareString(t, "GetReceiveAddressParse", "muVu2JZo8PbewBHRp6bpqFvVD87qvqEHWA", address)
}

func TestGetAllAddressesParse(t *testing.T) {
	c := initTestClient()
	params := &AddressesParams{}
	addresses, err := c.GetAllAddresses(params)
	if err != nil {
		log.Fatal(err)
	}
	compareString(t, "GetAllAddressesParse", "2013-05-09T23:07:08-07:00", addresses.Addresses[0].Address.Created_at)
	compareString(t, "GetAllAddressesParse", "mwigfecvyG4MZjb6R5jMbmNcs7TkzhUaCj", addresses.Addresses[1].Address.Address)
	compareInt(t, "GetAllAddressesParse", 1, int64(addresses.Num_pages))
}

func TestSendMoneyParse(t *testing.T) {
	c := initTestClient()
	params := &TransactionRequestParams{
		Transaction: &TransactionParams{},
	}
	data, err := c.SendMoney(params)
	if err != nil {
		log.Fatal(err)
	}
	compareBool(t, "SendMoneyParse", true, data.Success)
	compareString(t, "SendMoneyParse", "-1.23400000", data.Transaction.Amount.Amount)
	compareString(t, "SendMoneyParse", "37muSN5ZrukVTvyVh3mT5Zc5ew9L9CBare", data.Transaction.Recipient_address)
}

func TestRequestMoneyParse(t *testing.T) {
	c := initTestClient()
	params := &TransactionRequestParams{
		Transaction: &TransactionParams{},
	}
	data, err := c.RequestMoney(params)
	if err != nil {
		log.Fatal(err)
	}
	compareBool(t, "RequestMoneyParse", true, data.Success)
	compareString(t, "RequestMoneyParse", "1.23400000", data.Transaction.Amount.Amount)
	compareString(t, "RequestMoneyParse", "5011f33df8182b142400000e", data.Transaction.Recipient.Id)
}

func TestResendRequestParse(t *testing.T) {
	c := initTestClient()
	data, err := c.ResendRequest("ID")
	if err != nil {
		log.Fatal(err)
		fmt.Println(err)
	}
	compareBool(t, "ResendRequestParse", true, data)
}

func TestCancelRequestParse(t *testing.T) {
	c := initTestClient()
	data, err := c.CancelRequest("ID")
	if err != nil {
		log.Fatal(err)
	}
	compareBool(t, "CancelRequestParse", false, data)
}

func TestCompleteRequestParse(t *testing.T) {
	c := initTestClient()
	data, err := c.CompleteRequest("ID")
	if err != nil {
		log.Fatal(err)
	}
	compareBool(t, "CancelRequestParse", false, data.Success)
}

func TestCreateOrderFromButtonCodeParse(t *testing.T) {
	c := initTestClient()
	data, err := c.CreateOrderFromButtonCode("ID")
	if err != nil {
		log.Fatal(err)
	}
	compareBool(t, "CreateOrderFromButtonCodeParse", true, data.Success)
	compareString(t, "CreateOrderFromButtonCodeParse", "new", data.Order.Status)
}

func TestCreateUserParse(t *testing.T) {
	c := initTestClient()
	data, err := c.CreateUser("test@email.com", "password")
	if err != nil {
		log.Fatal(err)
	}
	compareBool(t, "CreateUser", false, data.Success)
	compareString(t, "CreateUser", "Email is already taken", data.Errors[0])
	compareString(t, "CreateUser", "501a3d9ef8182b2754000018", data.User.Id)
}

func TestBuyParse(t *testing.T) {
	c := initTestClient()
	data, err := c.Buy("1000", true)
	if err != nil {
		log.Fatal(err)
	}
	compareBool(t, "Buys", true, data.Success)
	compareString(t, "Buys", "2013-01-28T16:08:58-08:00", data.Transfer.Created_at)
	compareString(t, "Buys", "USD", data.Transfer.Fees.Bank.Currency_iso)
	compareString(t, "Buys", "13.55", data.Transfer.Subtotal.Amount)
}

func TestSellParse(t *testing.T) {
	c := initTestClient()
	data, err := c.Sell("1000")
	if err != nil {
		log.Fatal(err)
	}
	compareBool(t, "Sells", true, data.Success)
	compareString(t, "Sells", "2013-01-28T16:32:35-08:00", data.Transfer.Created_at)
	compareString(t, "Sells", "USD", data.Transfer.Fees.Bank.Currency_iso)
	compareString(t, "Sells", "13.50", data.Transfer.Subtotal.Amount)
}

func TestGetContactsParse(t *testing.T) {
	c := initTestClient()
	params := &ContactsParams{
		Page:  1,
		Limit: 5,
	}
	data, err := c.GetContacts(params)
	if err != nil {
		log.Fatal(err)
	}
	compareString(t, "GetContacts", "user1@example.com", data.Emails[0])
}

func TestGetTransactionsParse(t *testing.T) {
	c := initTestClient()
	data, err := c.GetTransactions(1)
	if err != nil {
		log.Fatal(err)
	}
	compareString(t, "GetTransactions", "5011f33df8182b142400000e", data.Current_user.Id)
	compareString(t, "GetTransactions", "500.00", data.Native_balance.Amount)
	compareString(t, "GetTransactions", "5018f833f8182b129c00002f", data.Transactions[0].Transaction.Id)
	compareString(t, "GetTransactions", "-1.00000000", data.Transactions[1].Transaction.Amount.Amount)
}

func TestGetOrdersParse(t *testing.T) {
	c := initTestClient()
	data, err := c.GetOrders(1)
	if err != nil {
		log.Fatal(err)
	}
	compareString(t, "GetOrders", "A7C52JQT", data.Orders[0].Order.Id)
	compareString(t, "GetOrders", "mgrmKftH5CeuFBU3THLWuTNKaZoCGJU5jQ", data.Orders[0].Order.Receive_address)
	compareString(t, "GetOrders", "buy_now", data.Orders[0].Order.Button.Type)
	compareInt(t, "GetOrders", int64(0), int64(data.Orders[0].Order.Transaction.Confirmations))
}

func TestGetTransfersParse(t *testing.T) {
	c := initTestClient()
	data, err := c.GetTransfers(1)
	if err != nil {
		log.Fatal(err)
	}
	compareString(t, "GetTransfers", "Buy", data.Transfers[0].Transfer.Type)
	compareString(t, "GetTransfers", "Pending", data.Transfers[0].Transfer.Status)
	compareString(t, "GetTransfers", "BTC", data.Transfers[0].Transfer.Btc.Currency)
	compareInt(t, "GetTransfers", int64(1), int64(data.Num_pages))
}

func TestGetTransaction(t *testing.T) {
	c := initTestClient()
	data, err := c.GetTransaction("ID")
	if err != nil {
		log.Fatal(err)
	}
	compareString(t, "GetTransaction", "5018f833f8182b129c00002f", data.Id)
	compareString(t, "GetTransaction", "BTC", data.Amount.Currency)
	compareString(t, "GetTransaction", "User One", data.Recipient.Name)
}

func TestGetOrder(t *testing.T) {
	c := initTestClient()
	data, err := c.GetOrder("ID")
	if err != nil {
		log.Fatal(err)
	}
	compareString(t, "GetTransaction", "A7C52JQT", data.Id)
	compareString(t, "GetTransaction", "BTC", data.Total_btc.Currency_iso)
	compareString(t, "GetTransaction", "test", data.Button.Name)
}

func TestGetUser(t *testing.T) {
	c := initTestClient()
	data, err := c.GetUser()
	if err != nil {
		log.Fatal(err)
	}
	compareString(t, "GetTransaction", "512db383f8182bd24d000001", data.Id)
	compareString(t, "GetTransaction", "49.76000000", data.Balance.Amount)
	compareString(t, "GetTransaction", "Company Name, Inc.", data.Merchant.Company_name)
}

func compareFloat(t *testing.T, prefix string, expected float64, got float64) {
	if expected != got {
		t.Errorf(`%s
			Expected %0.32f
			but got  %0.32f`, prefix, expected, got)
	}
}

func compareInt(t *testing.T, prefix string, expected int64, got int64) {
	if expected != got {
		t.Errorf("%s Expected %d but got %d", prefix, expected, got)
	}
}

func compareString(t *testing.T, prefix string, expected string, got string) {
	if expected != got {
		t.Errorf("%s Expected '%s' but got '%s'", prefix, expected, got)
	}
}

func compareBool(t *testing.T, prefix string, expected bool, got bool) {
	if expected != got {
		t.Errorf("%s Expected '%t' but got '%t'", prefix, expected, got)
	}
}
