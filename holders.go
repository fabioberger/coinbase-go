package coinbase

// Holders includes all the structs used for marshaling JSON responses from the coinbase API

// The holder used to marshal the JSON request returned in GetTokens
type tokensHolder struct {
	AccessToken  string `json:"access_token,omitempty"`
	TokenType    string `json:"token_type,omitempty"`
	ExpiresIn    int64  `json:"expires_in,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Scope        string `json:"scope,omitempty"`
}

// addressesHolder used to marshal the JSON request returned in GetAllAddresses
type addressesHolder struct {
	paginationStats
	Addresses []struct {
		Address address `json:"address,omitempty"`
	} `json:"addresses,omitempty"`
}

// orderHolder used to marshal the JSON request returned in CreateOrderFromButtonCode and GetOrder
type orderHolder struct {
	response
	Order order `json:"order,omitempty"`
}

// orderHolders used to marshal the JSON request returned in GetOrders
type ordersHolder struct {
	paginationStats
	Orders []struct {
		Order order `json:"order,omitempty"`
	} `json:"orders,omitempty"`
}

// buttonHolder used to marshal the JSON request returned in CreateButton
type buttonHolder struct {
	response
	Button Button `json:"button,omitempty"`
}

// transfersHolder used to marshal the JSON request returned in GetTransfers
type transfersHolder struct {
	paginationStats
	Transfers []struct {
		Transfer transfer `json:"transfer,omitempty"`
	} `json:"transfers,omitempty"`
}

// transferHolder used to marshal the JSON request returned in Buy & Sell
type transferHolder struct {
	response
	Transfer transfer `json:"transfer,omitempty"`
}

// pricesHolder used to marshal the JSON request returned in GetBuyPrice & GetSellPrice
type pricesHolder struct {
	Subtotal amount `json:"subtotal,omitempty"`
	Fees     []struct {
		Coinbase amount `json:"coinbase,omitempty"`
		Bank     amount `json:"bank,omitempty"`
	} `json:"fees,omitempty"`
	Total amount `json:"total,omitempty"`
}

// usersHolder used to marshal the JSON request returned in GetUser
type usersHolder struct {
	response
	Users []struct {
		User user `json:"user,omitempty"`
	} `json:"users,omitempty"`
}

// userHolder used to marshal the JSON request returned in CreateUser
type userHolder struct {
	response
	User  user  `json:"user,omitempty"`
	Oauth oauth `json:"oauth,omitempty"`
}

// The sub-structure of a response denominating its success and/or errors
type response struct {
	Success bool     `json:"success"`
	Errors  []string `json:"errors"`
	Error   string   `json:"error"`
}

// contactsHolder used to marshal the JSON request returned in GetContacts
type contactsHolder struct {
	paginationStats
	Contacts []contact `json:"contacts,omitempty"`
	Emails   []string  `json:"emails,omitempty"` // Add for convenience
}

// transactionHolder used to marshal the JSON request returned in SendMoney, RequestMoney,
// GetTransaction
type transactionHolder struct {
	response
	Transaction transaction `json:"transaction"`
	Transfer    transfer    `json:"transfer"`
}

// transactionsHolder used to marshal the JSON request returned in GetTransactions
type transactionsHolder struct {
	paginationStats
	CurrentUser   user   `json:"current_user,omitempty"`
	Balance       amount `json:"balance,omitempty"`
	NativeBalance amount `json:"native_balance,omitempty"`
	Transactions  []struct {
		Transaction transaction `json:"transaction,omitempty"`
	} `json:"transactions,omitempty"`
}
