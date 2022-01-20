package quotesprocessor

import (
	"time"

	"github.com/HaoxuanXu/TradingBot/internal/dataengine"
	"github.com/HaoxuanXu/TradingBot/strats/pairtrading/model"
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

	model.QuotesTimestamps[model.CheapStockSymbol] = time.Since(pairQuotes[model.CheapStockSymbol].Timestamp).Milliseconds()
	model.QuotesTimestamps[model.ExpensiveStockSymbol] = time.Since(pairQuotes[model.ExpensiveStockSymbol].Timestamp).Milliseconds()

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
