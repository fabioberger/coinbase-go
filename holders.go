package coinbase

// Holders includes all the structs used for marshaling JSON responses from the coinbase API

type tokensHolder struct {
	Access_token  string `json:"access_token,omitempty"`
	Token_type    string `json:"token_type,omitempty"`
	Expires_in    int64  `json:"expires_in,omitempty"`
	Refresh_token string `json:"refresh_token,omitempty"`
	Scope         string `json:"scope,omitempty"`
}

type addressesHolder struct {
	PaginationStats
	Addresses []struct {
		Address address `json:"address,omitempty"`
	} `json:"addresses,omitempty"`
}

type orderHolder struct {
	Response
	Order order `json:"order,omitempty"`
}

type ordersHolder struct {
	PaginationStats
	Orders []struct {
		Order order `json:"order,omitempty"`
	} `json:"orders,omitempty"`
}

type buttonHolder struct {
	Response
	Button button `json:"button,omitempty"`
}

type transfersHolder struct {
	PaginationStats
	Transfers []struct {
		Transfer transfer `json:"transfer,omitempty"`
	} `json:"transfers,omitempty"`
}

type transferHolder struct {
	Response
	Transfer transfer `json:"transfer,omitempty"`
}

type pricesHolder struct {
	Subtotal amount `json:"subtotal,omitempty"`
	Fees     []struct {
		Coinbase amount `json:"coinbase,omitempty"`
		Bank     amount `json:"bank,omitempty"`
	} `json:"fees,omitempty"`
	Total amount `json:"total,omitempty"`
}

type usersHolder struct {
	Response
	Users []struct {
		User user `json:"user,omitempty"`
	} `json:"users,omitempty"`
}

type userHolder struct {
	Response
	User  user  `json:"user,omitempty"`
	Oauth oauth `json:"oauth,omitempty"`
}

type Response struct {
	Success bool     `json:"success"`
	Errors  []string `json:"errors"`
	Error   string   `json:"error"`
}

type contactsHolder struct {
	PaginationStats
	Contacts []contact `json:"contacts,omitempty"`
	Emails   []string  `json:"emails,omitempty"` // Add for convenience
}

type transactionHolder struct {
	Response
	Transaction transaction `json:"transaction"`
	Transfer    transfer    `json:"transfer"`
}

type transactionsHolder struct {
	PaginationStats
	Current_user   user   `json:"current_user,omitempty"`
	Balance        amount `json:"balance,omitempty"`
	Native_balance amount `json:"native_balance,omitempty"`
	Transactions   []struct {
		Transaction transaction `json:"transaction,omitempty"`
	} `json:"transactions,omitempty"`
}
