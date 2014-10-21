package coinbase

type AddressesParams struct {
	Page       int    `json:"page,omitempty"`
	Limit      int    `json:"limit,omitempty"`
	Account_id string `json:"account_id,omitempty"`
	Query      string `json:"query,omitempty"`
}

type ReceiveAddressParams struct {
	Account_id string         `json:"account_id,omitempty"`
	Address    *AddressParams `json:"address,omitempty"`
}

type AddressParams struct {
	Label        string `json:"label,omitempty"`
	Callback_url string `json:"callback_url,omitempty"`
}

type TransactionRequestParams struct {
	Account_id  string             `json:"account_id,omitempty"`
	Transaction *TransactionParams `json:"transaction"`
}

type TransactionParams struct {
	To                  string `json:"to,omitempty"`
	From                string `json:"from,omitempty"`
	Amount              string `json:"amount,omitempty"`
	Amount_string       string `json:"amount_string,omitempty"`
	Amount_currency_iso string `json:"amount_currency_iso,omitempty"`
	Notes               string `json:"notes,omitempty"`
	User_fee            string `json:"user_fee,omitempty"`
	Referrer_id         string `json:"refferer_id,omitempty"`
	Idem                string `json:"idem,omitempty"`
	Instant_buy         bool   `json:"instant_buy,omitempty"`
	Order_id            string `json:"order_id,omitempty"`
}

type ButtonParams struct {
	Account_id string  `json:"account_id,omitempty"`
	Button     *Button `json:"button,omitempty"`
}

type ContactsParams struct {
	Page  int    `json:"page,omitempty"`
	Limit int    `json:"limit,omitempty"`
	Query string `json:"query,omitempty"`
}

type Currency struct {
	Name string `json:"name,omitempty"`
	Iso  string `json:"iso,omitempty"`
}

type Contact struct {
	Contact struct {
		Email string `json:"email,omitempty"`
	} `json:"contact,omitempty"`
}

type Button struct {
	Name                  string `json:"name,omitempty"`
	Price_string          string `json:"price_string,omitempty"`
	Price_currency_iso    string `json:"price_currency_iso,omitempty"`
	Type                  string `json:"type,omitempty"`
	Subscription          bool   `json:"subscription,omitempty"`
	Repeat                string `json:"repeat,omitempty"`
	Style                 string `json:"style,omitempty"`
	Text                  string `json:"text,omitempty"`
	Description           string `json:"description,omitempty"`
	Custom                string `json:"custom,omitempty"`
	Custom_secure         bool   `json:"custom_secure,omitempty"`
	Callback_url          string `json:"callback_url,omitempty"`
	Success_url           string `json:"success_url,omitempty"`
	Cancel_url            string `json:"cancel_url,omitempty"`
	Info_url              string `json:"info_url,omitempty"`
	Auto_redirect         bool   `json:"auto_redirect,omitempty"`
	Auto_redirect_success bool   `json:"auto_redirect_success,omitempty"`
	Auto_redirect_cancel  bool   `json:"auto_redirect_cancel,omitempty"`
	Variable_price        bool   `json:"variable_price,omitempty"`
	Choose_price          bool   `json:"choose_price,omitempty"`
	Include_address       bool   `json:"include_address,omitempty"`
	Include_email         bool   `json:"include_email,omitempty"`
	Price1                string `json:"price1,omitempty"`
	Price2                string `json:"price2,omitempty"`
	Price3                string `json:"price3,omitempty"`
	Price4                string `json:"price4,omitempty"`
	Price5                string `json:"price5,omitempty"`
	Code                  string `json:"code,omitempty"`
	Price                 Fee    `json:"price,omitempty"`
	Id                    string `json:"id,omitempty"`
}

type User struct {
	Id              string   `json:"id,omitempty"`
	Name            string   `json:"name,omitempty"`
	Email           string   `json:"email,omitempty"`
	Receive_address string   `json:"receive_address,omitempty"`
	Time_zone       string   `json:"timezone,omitempty"`
	Native_currency string   `json:"native_currency,omitempty"`
	Balance         Amount   `json:"balance,omitempty"`
	Merchant        Merchant `json:"merchant,omitempty"`
	Buy_level       int      `json:"buy_level,omitempty"`
	Sell_level      int      `json:"sell_level,omitempty"`
	Buy_limit       Amount   `json:"buy_limit,omitempty"`
	Sell_limit      Amount   `json:"sell_limit,omitempty"`
}

type Merchant struct {
	Company_name string `json:"company_name,omitempty"`
	Logo         struct {
		Small  string `json:"small,omitempty"`
		Medium string `json:"medium,omitempty"`
		Url    string `json:"url,omitempty"`
	} `json:"logo,omitempty"`
}

type Oauth struct {
	Access_token  string `json:"access_token,omitempty"`
	Token_type    string `json:"token_type,omitempty"`
	Expires_in    int    `json:"expires_in,omitempty"`
	Refresh_token string `json:"refresh_token,omitempty"`
	Scope         string `json:"scope,omitempty"`
}

type Transfer struct {
	Id             string `json:"id,omitempty"`
	Type           string `json:"type,omitempty"`
	Code           string `json:"code,omitempty"`
	Created_at     string `json:"created_at,omitempty"`
	Fees           Fees   `json:"fees,omitempty"`
	Status         string `json:"status,omitempty"`
	Payout_date    string `json:"payout_date,omitempty"`
	Btc            Amount `json:"btc,omitempty"`
	Subtotal       Amount `json:"subtotal,omitempty"`
	Total          Amount `json:"total,omitempty"`
	Description    string `json:"description,omitempty"`
	Transaction_id string `json:"transaction_id,omitempty"`
}

type Amount struct {
	Amount   string `json:"amount,omitempty"`
	Currency string `json:"currency,omitempty"`
}

type Fee struct {
	Cents        int    `json:"cents,omitempty"`
	Currency_iso string `json:"currency_iso,omitempty"`
}

type Fees struct {
	Coinbase Fee `json:"coinbase,omitempty"`
	Bank     Fee `json:"bank,omitempty"`
}

type TransactionActor struct {
	Id    string `json:"id,omitempty"`
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
}

type Transaction struct {
	Id                  string           `json:"id,omitempty"`
	Create_at           string           `json:"create_at,omitempty"`
	Hsh                 string           `json:"hsh,omitempty"`
	Notes               string           `json:"notes,omitempty"`
	Idem                string           `json:"idem,omitempty"`
	Amount              Amount           `json:"amount,omitempty"`
	Request             bool             `json:"request,omitempty"`
	Status              string           `json:"status,omitempty"`
	Sender              TransactionActor `json:"sender,omitempty"`
	Recipient           TransactionActor `json:"recipient,omitempty"`
	Recipient_address   string           `json:"recipient_address,omitempty"`
	Type                string           `json:"type,omitempty"`
	Signed              bool             `json:"signed,omitempty"`
	Signatures_required int              `json:"signature_required,omitempty"`
	Signatures_present  int              `json:"signatures_present,omitempty"`
	Signatures_needed   int              `json:"signatures_needed,omitempty"`
	Hash                string           `json:"hash,omitempty"`
	Confirmations       int              `json:"confirmations,omitempty"`
}

type Order struct {
	Id              string      `json:"id,omitempty"`
	Created_at      string      `json:"created_at,omitempty"`
	Status          string      `json:"status,omitempty"`
	Total_btc       Fee         `json:"total_btc,omitempty"`
	Total_native    Fee         `json:"total_native,omitempty"`
	Custom          string      `json:"custom,omitempty"`
	Receive_address string      `json:"receive_address,omitempty"`
	Button          Button      `json:"button,omitempty"`
	Transaction     Transaction `json:"transaction,omitempty"`
}
