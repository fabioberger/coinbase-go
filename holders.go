package coinbase

type addressesHolder struct {
	Addresses []struct {
		Address struct {
			Address      string `json:"address,omitempty"`
			Callback_url string `json:"callback_url,omitempty"`
			Label        string `json:"label,omitempty"`
			Created_at   string `json:"created_at,omitempty"`
		} `json:"address,omitempty"`
	} `json:"addresses,omitempty"`
	Total_count  int `json:"total_count,omitempty"`
	Num_pages    int `json:"num_pages,omitempty"`
	Current_page int `json:"current_page,omitempty"`
}

type orderHolder struct {
	Success bool  `json:"success,omitempty"`
	Order   Order `json:"order,omitempty"`
}

type ordersHolder struct {
	Orders []struct {
		Order Order `json:"order,omitempty"`
	} `json:"orders,omitempty"`
	Total_count  int `json:"total_count,omitempty"`
	Num_pages    int `json:"num_pages,omitempty"`
	Current_page int `json:"current_page,omitempty"`
}

type buttonHolder struct {
	Success    bool     `json:"success,omitempty"`
	Button     Button   `json:"button,omitempty"`
	Errors     []string `json:"errors,omitempty"`
	Embed_html string   `json:"embed_html,omitempty"` //Added embed_html for convenience
}

type transfersHolder struct {
	Transfers []struct {
		Transfer Transfer `json:"transfer,omitempty"`
	} `json:"transfers,omitempty"`
	Total_count  int `json:"total_count,omitempty"`
	Num_pages    int `json:"num_pages,omitempty"`
	Current_page int `json:"current_page,omitempty"`
}

type transferHolder struct {
	Success  bool     `json:"success,omitempty"`
	Errors   []string `json:"errors,omitempty"`
	Transfer Transfer `json:"transfer,omitempty"`
}

type pricesHolder struct {
	Subtotal Amount `json:"subtotal,omitempty"`
	Fees     []struct {
		Coinbase Amount `json:"coinbase,omitempty"`
		Bank     Amount `json:"bank,omitempty"`
	} `json:"fees,omitempty"`
	Total Amount `json:"total,omitempty"`
}

type usersHolder struct {
	Users []struct {
		User User `json:"user,omitempty"`
	} `json:"users,omitempty"`
}

type userHolder struct {
	Success bool     `json:"success,omitempty"`
	Errors  []string `json:"errors,omitempty"`
	User    User     `json:"user,omitempty"`
	Oauth   Oauth    `json:"oauth,omitempty"`
}

type contactsHolder struct {
	Contacts     []Contact `json:"contacts,omitempty"`
	Total_count  int       `json:"total_count,omitempty"`
	Num_pages    int       `json:"num_pages,omitempty"`
	Current_page int       `json:"current_page,omitempty"`
	Emails       []string  `json:"emails,omitempty"` // Add for convenience
}

type transactionHolder struct {
	Success     bool        `json:"success,omitempty"`
	Errors      []string    `json:"errors,omitempty"`
	Transaction Transaction `json:"transaction,omitempty"`
	Transfer    Transfer    `json:"transfer,omitempty"`
}

type transactionsHolder struct {
	Current_user   User   `json:"current_user,omitempty"`
	Balance        Amount `json:"balance,omitempty"`
	Native_balance Amount `json:"native_balance,omitempty"`
	Total_count    int    `json:"total_count,omitempty"`
	Num_pages      int    `json:"num_pages,omitempty"`
	Current_page   int    `json:"current_page,omitempty"`
	Transactions   []struct {
		Transaction Transaction `json:"transaction,omitempty"`
	} `json:"transactions,omitempty"`
}
