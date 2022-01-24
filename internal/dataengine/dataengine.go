package dataengine

import (
	"log"

	"github.com/HaoxuanXu/TradingBot/configs"
	"github.com/alpacahq/alpaca-trade-api-go/v2/marketdata"
)

type MarketDataEngine struct {
	client marketdata.Client
}

func GetDataEngine(accountType string) *MarketDataEngine {
	engine := &MarketDataEngine{}
	engine.initialize(accountType)
	return engine
}

func (engine *MarketDataEngine) initialize(accountType string) {
	cred := configs.GetCredentials(accountType)
	engine.client = marketdata.NewClient(
		marketdata.ClientOpts{
			ApiKey:    cred.API_KEY,
			ApiSecret: cred.API_SECRET,
		})
}

func (engine *MarketDataEngine) GetMultiQuotes(symbols []string) map[string]marketdata.Quote {

	quotes, err := engine.client.GetLatestQuotes(symbols)
	if err != nil {
		log.Println(err)
	}
	return quotes
}
