package coinbase

import (
	"log"
	"os"
	"testing"
)

// Initialize the client with mock mode enabled on rpc
// All calls return the corresponding json response from the test_data files
func initTestClient() Client {
	return apiKeyClientTest(os.Getenv("COINBASE_KEY"), os.Getenv("COINBASE_SECRET"))
}

// About Mock Tests:
// All Mock Tests simulate requests to the coinbase API by returning the expected
// return values from a file under the test_data folder. The values received from
// the file are compared with the expected value given the marshaling of the JSON
// executed correctly.

func TestMockGetBalanceParse(t *testing.T) {
	c := initTestClient()
	amount, err := c.GetBalance()
	if err != nil {
		log.Fatal(err)
	}
	compareFloat(t, "GetBalanceParse", 36.62800000, amount)
}

func TestMockGetReceiveAddressParse(t *testing.T) {
	c := initTestClient()
	params := &AddressParams{}
	address, err := c.GenerateReceiveAddress(params)
	if err != nil {
		log.Fatal(err)
	}
	compareString(t, "GetReceiveAddressParse", "muVu2JZo8PbewBHRp6bpqFvVD87qvqEHWA", address)
}

func TestMockGetAllAddressesParse(t *testing.T) {
	c := initTestClient()
	params := &AddressesParams{}
	addresses, err := c.GetAllAddresses(params)
	if err != nil {
		log.Fatal(err)
	}
	compareString(t, "GetAllAddressesParse", "2013-05-09T23:07:08-07:00", addresses.Addresses[0].CreatedAt)
	compareString(t, "GetAllAddressesParse", "mwigfecvyG4MZjb6R5jMbmNcs7TkzhUaCj", addresses.Addresses[1].Address)
	compareInt(t, "GetAllAddressesParse", 1, int64(addresses.NumPages))
}

func TestMockCreateButtonParse(t *testing.T) {
	c := initTestClient()
	params := &Button{}
	data, err := c.CreateButton(params)
	if err != nil {
		log.Fatal(err)
	}
	compareString(t, "CreateButtonParse", "93865b9cae83706ae59220c013bc0afd", data.Code)
	compareString(t, "CreateButtonParse", "Sample description", data.Description)
}

func TestMockSendMoneyParse(t *testing.T) {
	c := initTestClient()
	params := &TransactionParams{}
	data, err := c.SendMoney(params)
	if err != nil {
		log.Fatal(err)
	}
	compareString(t, "SendMoneyParse", "-1.23400000", data.Transaction.Amount.Amount)
	compareString(t, "SendMoneyParse", "37muSN5ZrukVTvyVh3mT5Zc5ew9L9CBare", data.Transaction.RecipientAddress)
}

func TestMockRequestMoneyParse(t *testing.T) {
	c := initTestClient()
	params := &TransactionParams{}
	data, err := c.RequestMoney(params)
	if err != nil {
		log.Fatal(err)
	}
	compareString(t, "RequestMoneyParse", "1.23400000", data.Transaction.Amount.Amount)
	compareString(t, "RequestMoneyParse", "5011f33df8182b142400000e", data.Transaction.Recipient.Id)
}

func TestMockResendRequestParse(t *testing.T) {
	c := initTestClient()
	data, err := c.ResendRequest("ID")
	if err != nil {
		log.Fatal(err)
	}
	compareBool(t, "ResendRequestParse", true, data)
}

func TestMockCancelRequestParse(t *testing.T) {
	c := initTestClient()
	data, err := c.CancelRequest("ID")
	if err != nil {
		log.Fatal(err)
	}
	compareBool(t, "CancelRequestParse", false, data)
}

func TestMockCompleteRequestParse(t *testing.T) {
	c := initTestClient()
	data, err := c.CompleteRequest("ID")
	if err != nil {
		log.Fatal(err)
	}
	compareString(t, "CancelRequestParse", "503c46a3f8182b106500009b", data.Transaction.Id)
}

func TestMockCreateOrderFromButtonCodeParse(t *testing.T) {
	c := initTestClient()
	data, err := c.CreateOrderFromButtonCode("ID")
	if err != nil {
		log.Fatal(err)
	}
	compareString(t, "CreateOrderFromButtonCodeParse", "7RTTRDVP", data.Id)
	compareString(t, "CreateOrderFromButtonCodeParse", "new", data.Status)
}

func TestMockCreateUserParse(t *testing.T) {
	c := initTestClient()
	data, err := c.CreateUser("test@email.com", "password")
	if err != nil {
		log.Fatal(err)
	}
	compareString(t, "CreateUser", "newuser@example.com", data.Email)
	compareString(t, "CreateUser", "501a3d22f8182b2754000011", data.Id)
}

func TestMockBuyParse(t *testing.T) {
	c := initTestClient()
	data, err := c.Buy(1000.0, true)
	if err != nil {
		log.Fatal(err)
	}
	compareString(t, "Buys", "2013-01-28T16:08:58-08:00", data.CreatedAt)
	compareString(t, "Buys", "USD", data.Fees.Bank.CurrencyIso)
	compareString(t, "Buys", "13.55", data.Subtotal.Amount)
}

func TestMockSellParse(t *testing.T) {
	c := initTestClient()
	data, err := c.Sell(1000.0)
	if err != nil {
		log.Fatal(err)
	}
	compareString(t, "Sells", "2013-01-28T16:32:35-08:00", data.CreatedAt)
	compareString(t, "Sells", "USD", data.Fees.Bank.CurrencyIso)
	compareString(t, "Sells", "13.50", data.Subtotal.Amount)
}

func TestMockGetContactsParse(t *testing.T) {
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

func TestMockGetTransactionsParse(t *testing.T) {
	c := initTestClient()
	data, err := c.GetTransactions(1)
	if err != nil {
		log.Fatal(err)
	}
	compareInt(t, "GetTransactions", 2, data.TotalCount)
	compareString(t, "GetTransactions", "5018f833f8182b129c00002f", data.Transactions[0].Id)
	compareString(t, "GetTransactions", "-1.00000000", data.Transactions[1].Amount.Amount)
}

func TestMockGetOrdersParse(t *testing.T) {
	c := initTestClient()
	data, err := c.GetOrders(1)
	if err != nil {
		log.Fatal(err)
	}
	compareString(t, "GetOrders", "buy_now", data.Orders[0].Button.Type)
	compareInt(t, "GetOrders", int64(0), int64(data.Orders[0].Transaction.Confirmations))
}

func TestMockGetTransfersParse(t *testing.T) {
	c := initTestClient()
	data, err := c.GetTransfers(1)
	if err != nil {
		log.Fatal(err)
	}
	compareString(t, "GetTransfers", "BTC", data.Transfers[0].Btc.Currency)
	compareInt(t, "GetTransfers", int64(1), int64(data.NumPages))
}

func TestMockGetTransaction(t *testing.T) {
	c := initTestClient()
	data, err := c.GetTransaction("ID")
	if err != nil {
		log.Fatal(err)
	}
	compareString(t, "GetTransaction", "5018f833f8182b129c00002f", data.Id)
	compareString(t, "GetTransaction", "BTC", data.Amount.Currency)
	compareString(t, "GetTransaction", "User One", data.Recipient.Name)
}

func TestMockGetOrder(t *testing.T) {
	c := initTestClient()
	data, err := c.GetOrder("ID")
	if err != nil {
		log.Fatal(err)
	}
	compareString(t, "GetTransaction", "A7C52JQT", data.Id)
	compareString(t, "GetTransaction", "BTC", data.TotalBtc.CurrencyIso)
	compareString(t, "GetTransaction", "test", data.Button.Name)
}

func TestMockGetUser(t *testing.T) {
	c := initTestClient()
	data, err := c.GetUser()
	if err != nil {
		log.Fatal(err)
	}
	compareString(t, "GetTransaction", "512db383f8182bd24d000001", data.Id)
	compareString(t, "GetTransaction", "49.76000000", data.Balance.Amount)
	compareString(t, "GetTransaction", "Company Name, Inc.", data.Merchant.CompanyName)
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
