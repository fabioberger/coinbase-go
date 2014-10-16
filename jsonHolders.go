package coinbase

type addressesHolder struct {
	Addresses []struct {
		Address struct {
			Address      string `json:"address"`
			Callback_url string `json:"callback_url"`
			Label        string `json:"label"`
			Created_at   string `json:"created_at"`
		} `json:"address"`
	} `json:"addresses"`
	Total_count  int `json:"total_count"`
	Num_pages    int `json:"num_pages"`
	Current_page int `json:"current_page"`
}
