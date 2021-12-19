package configs

import (
	"github.com/HaoxuanXu/TradingBot/configs/credentials/live"
	"github.com/HaoxuanXu/TradingBot/configs/credentials/paper"
)

type Credentials struct {
	API_KEY    string `json:"api_key"`
	API_SECRET string `json:"api_secret"`
	BASE_URL   string `json:"base_url"`
}

func GetCredentials(accountType string) Credentials {
	var credentials Credentials
	if accountType == "live" {
		credentials = Credentials{
			API_KEY: live.API_KEY,
			API_SECRET: live.API_SECRET,
			BASE_URL: live.BASE_URL,
		}
	} else if accountType == "paper" {
		credentials = Credentials{
			API_KEY: paper.API_KEY,
			API_SECRET: paper.API_SECRET,
			BASE_URL: paper.BASE_URL,
		}
	} 
	return credentials
}