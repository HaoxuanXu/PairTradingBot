package logging

import (
	"log"
	"math"

	"github.com/HaoxuanXu/TradingBot/internal/broker"
	"github.com/HaoxuanXu/TradingBot/strats/pairtrading/model"
)

func LogTransaction(model *model.PairTradingModel, broker *broker.AlpacaBroker) {
	if !broker.HasPosition {
		if model.IsLongExpensiveStockShortCheapStock {
			model.EntryNetValue = math.Abs(model.CheapStockFilledPrice*model.CheapStockFilledQuantity) - math.Abs(model.ExpensiveStockFilledPrice*model.ExpensiveStockFilledQuantity)
			log.Printf("long %s: %f shares; short %s: %f shares\n",
				model.ExpensiveStockSymbol,
				model.ExpensiveStockEntryVolume, model.CheapStockSymbol, model.CheapStockEntryVolume)
		} else {
			model.EntryNetValue = math.Abs(model.ExpensiveStockFilledPrice*model.ExpensiveStockFilledQuantity) - math.Abs(model.CheapStockFilledPrice*model.CheapStockFilledQuantity)
			log.Printf("short %s: %f shares; long %s: %f shares\n",
				model.ExpensiveStockSymbol, model.ExpensiveStockEntryVolume, model.CheapStockSymbol,
				model.CheapStockEntryVolume)
		}
		broker.HasPosition = true
	} else if broker.HasPosition {
		presumedProfit := model.ExitNetValue + model.EntryNetValue
		if model.IsShortExpensiveStockLongCheapStock {
			model.ExitNetValue = math.Abs(model.CheapStockFilledQuantity*model.CheapStockFilledPrice) - math.Abs(model.ExpensiveStockFilledQuantity*model.ExpensiveStockFilledPrice)
		} else {
			model.ExitNetValue = math.Abs(model.ExpensiveStockFilledPrice*model.ExpensiveStockFilledQuantity) - math.Abs(model.CheapStockFilledPrice*model.CheapStockFilledQuantity)
		}
		actualProfit := model.ExitNetValue + model.EntryNetValue

		if actualProfit < 0 {
			model.LoserNums++
		}
		log.Printf("position closed. Presumed Profit: $%f. Actual Profit: $%f\n",
			presumedProfit, actualProfit)
		broker.HasPosition = false
		broker.TransactionNums++
	}
}
