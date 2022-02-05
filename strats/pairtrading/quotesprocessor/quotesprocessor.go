package quotesprocessor

import (
	"log"
	"strconv"
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

func WarmUpData(timeDuration, assetType string, model *model.PairTradingModel, dataEngine *dataengine.MarketDataEngine, broker *broker.AlpacaBroker) {
	now := time.Now()
	timeDurationInt, _ := strconv.Atoi(timeDuration)
	loc, _ := time.LoadLocation("America/New_York")
	marketOpen := time.Date(now.Year(), now.Month(), now.Day(), 9, 30, 0, 0, loc)
	log.Printf("Start warming data until %s minutes after the market opens...", timeDuration)
	for time.Since(marketOpen) < time.Duration(timeDurationInt)*time.Minute {
		GetAndProcessPairQuotes(model, dataEngine)
	}
	transaction.SlideRepeatAndPriceRatioArrays(model)
	model.UpdateParameters()
	pipeline.WriteRecord(model)
	log.Println("Data-warming complete!")
}
