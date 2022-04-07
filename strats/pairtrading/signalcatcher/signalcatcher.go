package signalcatcher

import (
	"github.com/HaoxuanXu/TradingBot/internal/broker"
	"github.com/HaoxuanXu/TradingBot/strats/pairtrading/model"
)

func GetEntrySignal(shortExpensiveStock bool, model *model.PairTradingModel, broker *broker.AlpacaBroker) bool {
	if !broker.HasPosition {
		if shortExpensiveStock {
			if model.ShortExpensiveStockLongCheapStockPriceRatio > model.PriceRatioThreshold &&
				model.ShortExpensiveStockLongCheapStockRepeatNumber >= model.ShortExpensiveLongCheapRepeatNumThreshold &&
				model.FastShortExpensiveLongCheapRepeatNumThreshold == model.ShortExpensiveLongCheapRepeatNumThreshold {
				return true
			}
		} else {
			if model.LongExpensiveStockShortCheapStockPriceRatio < model.PriceRatioThreshold &&
				model.LongExpensiveStockShortCheapStockRepeatNumber >= model.LongExpensiveShortCheapRepeatNumThreshold &&
				model.FastLongExpensiveShortCheapRepeatNumThreshold == model.LongExpensiveShortCheapRepeatNumThreshold {
				return true
			}
		}
	}
	return false
}

func GetExitSignal(model *model.PairTradingModel) bool {
	if model.IsShortExpensiveStockLongCheapStock &&
		model.LongExpensiveStockShortCheapStockRepeatNumber >= model.LongExpensiveShortCheapRepeatNumThreshold &&
		model.FastLongExpensiveShortCheapRepeatNumThreshold == model.LongExpensiveShortCheapRepeatNumThreshold {
		model.ExitNetValue = model.CheapStockShortQuotePrice*model.CheapStockEntryVolume - model.ExpensiveStockLongQuotePrice*model.ExpensiveStockEntryVolume
		if model.ExitNetValue+model.EntryNetValue >= model.MinProfitThreshold {
			return true
		}
	} else if model.IsLongExpensiveStockShortCheapStock &&
		model.ShortExpensiveStockLongCheapStockRepeatNumber >= model.ShortExpensiveLongCheapRepeatNumThreshold &&
		model.FastShortExpensiveLongCheapRepeatNumThreshold == model.ShortExpensiveLongCheapRepeatNumThreshold {
		model.ExitNetValue = model.ExpensiveStockShortQuotePrice*model.ExpensiveStockEntryVolume - model.CheapStockLongQuotePrice*model.CheapStockEntryVolume
		if model.ExitNetValue+model.EntryNetValue >= model.MinProfitThreshold {
			return true
		}
	}
	return false
}
