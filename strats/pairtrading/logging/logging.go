package logging

import (
	"log"
	"time"

	"github.com/HaoxuanXu/TradingBot/internal/broker"
	"github.com/HaoxuanXu/TradingBot/strats/pairtrading/model"
)

func LogTransaction(model *model.PairTradingModel, broker *broker.AlpacaBroker) {
	dt := time.Now()
	if !broker.HasPosition {
		if model.IsLongExpensiveStockShortCheapStock {
			model.EntryNetValue = model.CheapStockFilledPrice*model.CheapStockFilledQuantity - model.ExpensiveStockFilledPrice*model.ExpensiveStockFilledQuantity
			log.Printf("[%d:%d:%d] -- long %s: %f shares; short %s: %f shares\n",
				dt.Hour(), dt.Minute(), dt.Second(), model.ExpensiveStockSymbol,
				model.ExpensiveStockEntryVolume, model.CheapStockSymbol, model.CheapStockEntryVolume)
		} else {
			model.EntryNetValue = model.ExpensiveStockFilledPrice*model.ExpensiveStockFilledQuantity - model.CheapStockFilledPrice*model.CheapStockFilledQuantity
			log.Printf("[%d:%d:%d] -- short %s: %f shares; long %s: %f shares\n", dt.Hour(),
				dt.Minute(), dt.Second(), model.ExpensiveStockSymbol, model.ExpensiveStockEntryVolume, model.CheapStockSymbol,
				model.CheapStockEntryVolume)
		}
		broker.HasPosition = true
	} else if broker.HasPosition {
		presumedProfit := model.ExitNetValue + model.EntryNetValue
		if model.IsLongExpensiveStockShortCheapStock {
			model.ExitNetValue = model.CheapStockFilledQuantity*model.CheapStockFilledPrice - model.ExpensiveStockFilledQuantity*model.ExpensiveStockFilledPrice
		} else {
			model.ExitNetValue = model.ExpensiveStockFilledPrice*model.ExpensiveStockFilledQuantity - model.CheapStockFilledPrice*model.CheapStockFilledQuantity
		}
		actualProfit := model.ExitNetValue + model.EntryNetValue

		if actualProfit < 0 {
			model.LoserNums++
		}
		log.Printf("[%d:%d:%d] -- position closed. Presumed Profit: $%f. Actual Profit: $%f\n", dt.Hour(),
			dt.Minute(), dt.Second(), presumedProfit, actualProfit)
		broker.HasPosition = false
		broker.TransactionNums++
	}
}