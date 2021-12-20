package quotesprocessor

import (
	"github.com/HaoxuanXu/TradingBot/internal/dataengine"
	"github.com/HaoxuanXu/TradingBot/strats/pairtrading/model"
	"github.com/HaoxuanXu/TradingBot/strats/pairtrading/transaction"
)

func GetAndProcessPairQuotes(model *model.PairTradingModel, dataEngine *dataengine.MarketDataEngine) {
	pairQuotes := dataEngine.GetMultiQuotes(
		[]string{model.ExpensiveStockSymbol, model.CheapStockSymbol},
	)

	model.CheapStockLongQuotePrice = pairQuotes[model.CheapStockSymbol].AskPrice
	model.CheapStockShortQuotePrice = pairQuotes[model.CheapStockSymbol].BidPrice
	model.ExpensiveStockLongQuotePrice = pairQuotes[model.ExpensiveStockSymbol].AskPrice
	model.ExpensiveStockShortQuotePrice = pairQuotes[model.ExpensiveStockSymbol].BidPrice

	model.LongExpensiveStockShortCheapStockPriceRatio = float64(model.ExpensiveStockLongQuotePrice / model.CheapStockShortQuotePrice)
	model.ShortExpensiveStockLongCheapStockPriceRatio = float64(model.ExpensiveStockShortQuotePrice / model.CheapStockLongQuotePrice)

	transaction.UpdateFieldsFromQuotes(model)

}
