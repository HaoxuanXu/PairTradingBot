package dataengine

import (
	"sync"

	"github.com/HaoxuanXu/TradingBot/configs"
	"github.com/alpacahq/alpaca-trade-api-go/v2/marketdata"
)

type MarketDataEngine struct {
	client   marketdata.Client
	channel1 chan marketdata.Quote
	channel2 chan marketdata.Quote
	WG       sync.WaitGroup
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
	engine.channel1 = make(chan marketdata.Quote)
	engine.channel2 = make(chan marketdata.Quote)
}

func (engine *MarketDataEngine) GetMultiQuotes(symbols []string) map[string]marketdata.Quote {

	quotes, _ := engine.client.GetLatestQuotes(symbols)
	return quotes
}
