package coinbase

// Params includes all the struct parameters that are required for specific API requests
// By defining a specific param struct, a developer can know which parameters are allowed
// for a given request. Also included here are the return object structs returned by
// specific API calls

// Parameter Struct for GET /api/v1/addresses Requests
type AddressesParams struct {
	Page      int64  `json:"page,omitempty"`
	Limit     int64  `json:"limit,omitempty"`
	AccountId string `json:"account_id,omitempty"`
	Query     string `json:"query,omitempty"`
}

// Parameter Struct for POST /api/v1/account/generate_receive_address Requests
type AddressParams struct {
	Label       string `json:"label,omitempty"`
	CallbackUrl string `json:"callback_url,omitempty"`
}

// Parameter Struct for POST /api/v1/transactions/(request_money,send_money) Requests
type TransactionParams struct {
	To                string `json:"to,omitempty"`
	From              string `json:"from,omitempty"`
	Amount            string `json:"amount,omitempty"`
	AmountString      string `json:"amount_string,omitempty"`
	AmountCurrencyIso string `json:"amount_currency_iso,omitempty"`
	Notes             string `json:"notes,omitempty"`
	UserFee           string `json:"user_fee,omitempty"`
	ReferrerId        string `json:"refferer_id,omitempty"`
	Idem              string `json:"idem,omitempty"`
	InstantBuy        bool   `json:"instant_buy,omitempty"`
	OrderId           string `json:"order_id,omitempty"`
}

// Parameter Struct for GET /api/v1/contacts Requests
type ContactsParams struct {
	Page  int64  `json:"page,omitempty"`
	Limit int64  `json:"limit,omitempty"`
	Query string `json:"query,omitempty"`
}

// The OAuth Tokens Struct returned from OAuth Authentication
type oauthTokens struct {
	AccessToken  string
	RefreshToken string
	ExpireTime   int64
}

// The return response from SendMoney, RequestMoney, CompleteRequest
type transactionConfirmation struct {
	Transaction transaction
	Transfer    transfer
}

// The return response from GetAllAddresses
type addresses struct {
	paginationStats
	Addresses []address
}

// The structure for one returned address from GetAllAddresses
type address struct {
	Address     string `json:"address,omitempty"`
	CallbackUrl string `json:"callback_url,omitempty"`
	Label       string `json:"label,omitempty"`
	CreatedAt   string `json:"created_at,omitempty"`
}

// The sub-structure of a response denominating a currency
type currency struct {
	Name string `json:"name,omitempty"`
	Iso  string `json:"iso,omitempty"`
}

// The sub-structure of a response denominating a contact
type contact struct {
	Contact struct {
		Email string `json:"email,omitempty"`
	} `json:"contact,omitempty"`
}

// The return response from CreateButton
type Button struct {
	Name                string `json:"name,omitempty"`
	PriceString         string `json:"price_string,omitempty"`
	PriceCurrencyIso    string `json:"price_currency_iso,omitempty"`
	Type                string `json:"type,omitempty"`
	Subscription        bool   `json:"subscription,omitempty"`
	Repeat              string `json:"repeat,omitempty"`
	Style               string `json:"style,omitempty"`
	Text                string `json:"text,omitempty"`
	Description         string `json:"description,omitempty"`
	Custom              string `json:"custom,omitempty"`
	CustomSecure        bool   `json:"custom_secure,omitempty"`
	CallbackUrl         string `json:"callback_url,omitempty"`
	SuccessUrl          string `json:"success_url,omitempty"`
	CancelUrl           string `json:"cancel_url,omitempty"`
	InfoUrl             string `json:"info_url,omitempty"`
	AutoRedirect        bool   `json:"auto_redirect,omitempty"`
	AutoRedirectSuccess bool   `json:"auto_redirect_success,omitempty"`
	AutoRedirectCancel  bool   `json:"auto_redirect_cancel,omitempty"`
	VariablePrice       bool   `json:"variable_price,omitempty"`
	ChoosePrice         bool   `json:"choose_price,omitempty"`
	IncludeAddress      bool   `json:"include_address,omitempty"`
	IncludeEmail        bool   `json:"include_email,omitempty"`
	Price1              string `json:"price1,omitempty"`
	Price2              string `json:"price2,omitempty"`
	Price3              string `json:"price3,omitempty"`
	Price4              string `json:"price4,omitempty"`
	Price5              string `json:"price5,omitempty"`
	Code                string `json:"code,omitempty"`
	Price               fee    `json:"price,omitempty"`
	Id                  string `json:"id,omitempty"`
	EmbedHtml           string `json:"embed_html"` //Added embed_html for convenience
}

// The return response from GetUser and CreateUser
type user struct {
	Id             string   `json:"id,omitempty"`
	Name           string   `json:"name,omitempty"`
	Email          string   `json:"email,omitempty"`
	ReceiveAddress string   `json:"receive_address,omitempty"`
	TimeZone       string   `json:"timezone,omitempty"`
	NativeCurrency string   `json:"native_currency,omitempty"`
	Balance        amount   `json:"balance,omitempty"`
	Merchant       merchant `json:"merchant,omitempty"`
	BuyLevel       int64    `json:"buy_level,omitempty"`
	SellLevel      int64    `json:"sell_level,omitempty"`
	BuyLimit       amount   `json:"buy_limit,omitempty"`
	SellLimit      amount   `json:"sell_limit,omitempty"`
}

// The sub-structure of a response denominating a merchant
type merchant struct {
	CompanyName string `json:"company_name,omitempty"`
	Logo        struct {
		Small  string `json:"small,omitempty"`
		Medium string `json:"medium,omitempty"`
		Url    string `json:"url,omitempty"`
	} `json:"logo,omitempty"`
}

// The sub-structure of a response denominating the oauth data
type oauth struct {
	AccessToken  string `json:"access_token,omitempty"`
	TokenType    string `json:"token_type,omitempty"`
	ExpiresIn    int64  `json:"expires_in,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Scope        string `json:"scope,omitempty"`
}

// The sub-structure of a response denominating pagination stats
type paginationStats struct {
	TotalCount  int64 `json:"total_count,omitempty"`
	NumPages    int64 `json:"num_pages,omitempty"`
	CurrentPage int64 `json:"current_page,omitempty"`
}

// The return response from GetTransfers
type transfers struct {
	paginationStats
	Transfers []transfer
}

// The sub-structure of a response denominating a transfer
type transfer struct {
	Id            string `json:"id,omitempty"`
	Type          string `json:"type,omitempty"`
	Code          string `json:"code,omitempty"`
	CreatedAt     string `json:"created_at,omitempty"`
	Fees          fees   `json:"fees,omitempty"`
	Status        string `json:"status,omitempty"`
	PayoutDate    string `json:"payout_date,omitempty"`
	Btc           amount `json:"btc,omitempty"`
	Subtotal      amount `json:"subtotal,omitempty"`
	Total         amount `json:"total,omitempty"`
	Description   string `json:"description,omitempty"`
	TransactionId string `json:"transaction_id,omitempty"`
}

// The sub-structure of a response denominating an amount
type amount struct {
	Amount   string `json:"amount,omitempty"`
	Currency string `json:"currency,omitempty"`
}

// The sub-structure of a response denominating a fee
type fee struct {
	Cents       float64 `json:"cents,omitempty"`
	CurrencyIso string  `json:"currency_iso,omitempty"`
}

// The sub-structure of a response denominating fees
type fees struct {
	Coinbase fee `json:"coinbase,omitempty"`
	Bank     fee `json:"bank,omitempty"`
}

// The sub-structure of a response denominating a transaction actor
type transactionActor struct {
	Id    string `json:"id,omitempty"`
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
}

// The return response from GetTransactions
type transactions struct {
	paginationStats
	Transactions []transaction
}

// The sub-structure of a response denominating a transaction
type transaction struct {
	Id                 string           `json:"id,omitempty"`
	CreateAt           string           `json:"create_at,omitempty"`
	Hsh                string           `json:"hsh,omitempty"`
	Notes              string           `json:"notes,omitempty"`
	Idem               string           `json:"idem,omitempty"`
	Amount             amount           `json:"amount,omitempty"`
	Request            bool             `json:"request,omitempty"`
	Status             string           `json:"status,omitempty"`
	Sender             transactionActor `json:"sender,omitempty"`
	Recipient          transactionActor `json:"recipient,omitempty"`
	RecipientAddress   string           `json:"recipient_address,omitempty"`
	Type               string           `json:"type,omitempty"`
	Signed             bool             `json:"signed,omitempty"`
	SignaturesRequired int64            `json:"signature_required,omitempty"`
	SignaturesPresent  int64            `json:"signatures_present,omitempty"`
	SignaturesNeeded   int64            `json:"signatures_needed,omitempty"`
	Hash               string           `json:"hash,omitempty"`
	Confirmations      int64            `json:"confirmations,omitempty"`
}

// The return response from GetOrders
type orders struct {
	paginationStats
	Orders []order
}

// The sub-structure of a response denominating an order
type order struct {
	Id             string      `json:"id,omitempty"`
	CreatedAt      string      `json:"created_at,omitempty"`
	Status         string      `json:"status,omitempty"`
	TotalBtc       fee         `json:"total_btc,omitempty"`
	TotalNative    fee         `json:"total_native,omitempty"`
	Custom         string      `json:"custom,omitempty"`
	ReceiveAddress string      `json:"receive_address,omitempty"`
	Button         Button      `json:"button,omitempty"`
	Transaction    transaction `json:"transaction,omitempty"`
}
