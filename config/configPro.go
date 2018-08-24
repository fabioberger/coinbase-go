package config

var (
	BaseUrl string
	Sandbox = false // set to true if you want to use the sandbox API endpoint
)

func init() {
	BaseUrl = "https://api.pro.coinbase.com/v1/"
	if Sandbox == true {
		BaseUrl = "https://api.pro.sandbox.coinbase.com/v1/"
	}
}
