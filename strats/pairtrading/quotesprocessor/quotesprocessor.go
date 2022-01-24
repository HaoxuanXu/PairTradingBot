package quotesprocessor

import (
	"log"
	"time"

	"github.com/HaoxuanXu/TradingBot/internal/broker"
	"github.com/HaoxuanXu/TradingBot/internal/dataengine"
	"github.com/HaoxuanXu/TradingBot/strats/pairtrading/model"
	"github.com/HaoxuanXu/TradingBot/strats/pairtrading/pipeline"
	"github.com/HaoxuanXu/TradingBot/strats/pairtrading/transaction"
)

func checkIfNotZeros(input []float64) bool {
	for _, val := range input {
		if val == 0.0 {
			return false
		}
	}
	return true
}

func GetAndProcessPairQuotes(model *model.PairTradingModel, dataEngine *dataengine.MarketDataEngine) {
	pairQuotes := dataEngine.GetMultiQuotes(
		[]string{model.ExpensiveStockSymbol, model.CheapStockSymbol},
	)
	model.QuotesConditions[model.CheapStockSymbol] = pairQuotes[model.CheapStockSymbol].Conditions
	model.QuotesConditions[model.ExpensiveStockSymbol] = pairQuotes[model.ExpensiveStockSymbol].Conditions

	model.CheapStockLongQuotePrice = pairQuotes[model.CheapStockSymbol].AskPrice
	model.CheapStockShortQuotePrice = pairQuotes[model.CheapStockSymbol].BidPrice
	model.ExpensiveStockLongQuotePrice = pairQuotes[model.ExpensiveStockSymbol].AskPrice
	model.ExpensiveStockShortQuotePrice = pairQuotes[model.ExpensiveStockSymbol].BidPrice

	if checkIfNotZeros([]float64{
		model.CheapStockLongQuotePrice,
		model.CheapStockShortQuotePrice,
		model.ExpensiveStockLongQuotePrice,
		model.ExpensiveStockShortQuotePrice,
	}) {
		model.LongExpensiveStockShortCheapStockPriceRatio = float64(model.ExpensiveStockLongQuotePrice / model.CheapStockShortQuotePrice)
		model.ShortExpensiveStockLongCheapStockPriceRatio = float64(model.ExpensiveStockShortQuotePrice / model.CheapStockLongQuotePrice)
		transaction.UpdateFieldsFromQuotes(model)
	}
}

func WarmUpData(timeUntil, assetType string, model *model.PairTradingModel, dataEngine *dataengine.MarketDataEngine, broker *broker.AlpacaBroker) {
	parsedTime, err := time.Parse("10:00:00", timeUntil)
	if err != nil {
		log.Println(err)
	}
	log.Printf("Start Warming up data until %s...\n", timeUntil)
	for time.Until(parsedTime) > 0 {
		GetAndProcessPairQuotes(model, dataEngine)
	}
	transaction.SlideRepeatAndPriceRatioArrays(model)
	model.UpdateParameters()
	pipeline.WriteRecord(model, assetType)
	log.Println("Data-warming complete!")
}
