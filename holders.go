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
	Order   order `json:"order,omitempty"`
}

type ordersHolder struct {
	Orders []struct {
		Order order `json:"order,omitempty"`
	} `json:"orders,omitempty"`
	Total_count  int `json:"total_count,omitempty"`
	Num_pages    int `json:"num_pages,omitempty"`
	Current_page int `json:"current_page,omitempty"`
}

type buttonHolder struct {
	Success    bool     `json:"success,omitempty"`
	Button     button   `json:"button,omitempty"`
	Errors     []string `json:"errors,omitempty"`
	Embed_html string   `json:"embed_html,omitempty"` //Added embed_html for convenience
}

type transfersHolder struct {
	Transfers []struct {
		Transfer transfer `json:"transfer,omitempty"`
	} `json:"transfers,omitempty"`
	Total_count  int `json:"total_count,omitempty"`
	Num_pages    int `json:"num_pages,omitempty"`
	Current_page int `json:"current_page,omitempty"`
}

type transferHolder struct {
	Success  bool     `json:"success,omitempty"`
	Errors   []string `json:"errors,omitempty"`
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
	Users []struct {
		User user `json:"user,omitempty"`
	} `json:"users,omitempty"`
}

type userHolder struct {
	Success bool     `json:"success,omitempty"`
	Errors  []string `json:"errors,omitempty"`
	User    user     `json:"user,omitempty"`
	Oauth   oauth    `json:"oauth,omitempty"`
}

type contactsHolder struct {
	Contacts     []contact `json:"contacts,omitempty"`
	Total_count  int       `json:"total_count,omitempty"`
	Num_pages    int       `json:"num_pages,omitempty"`
	Current_page int       `json:"current_page,omitempty"`
	Emails       []string  `json:"emails,omitempty"` // Add for convenience
}

type transactionHolder struct {
	Success     bool        `json:"success,omitempty"`
	Errors      []string    `json:"errors,omitempty"`
	Transaction transaction `json:"transaction,omitempty"`
	Transfer    transfer    `json:"transfer,omitempty"`
}

type transactionsHolder struct {
	Current_user   user   `json:"current_user,omitempty"`
	Balance        amount `json:"balance,omitempty"`
	Native_balance amount `json:"native_balance,omitempty"`
	Total_count    int    `json:"total_count,omitempty"`
	Num_pages      int    `json:"num_pages,omitempty"`
	Current_page   int    `json:"current_page,omitempty"`
	Transactions   []struct {
		Transaction transaction `json:"transaction,omitempty"`
	} `json:"transactions,omitempty"`
}
