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
				model.StrategyAssetType != "technology" {
				return true
			}
		} else {
			if model.LongExpensiveStockShortCheapStockPriceRatio < model.PriceRatioThreshold &&
				model.LongExpensiveStockShortCheapStockRepeatNumber >= model.LongExpensiveShortCheapRepeatNumThreshold {
				return true
			}
		}
	}
	return false
}

func GetExitSignal(model *model.PairTradingModel) bool {
	if model.IsShortExpensiveStockLongCheapStock &&
		model.LongExpensiveStockShortCheapStockRepeatNumber >= model.LongExpensiveShortCheapRepeatNumThreshold {
		model.ExitNetValue = model.CheapStockShortQuotePrice*model.CheapStockEntryVolume - model.ExpensiveStockLongQuotePrice*model.ExpensiveStockEntryVolume
		if model.ExitNetValue+model.EntryNetValue >= model.MinProfitThreshold.Applied {
			return true
		}
	} else if model.IsLongExpensiveStockShortCheapStock &&
		model.ShortExpensiveStockLongCheapStockRepeatNumber >= model.ShortExpensiveLongCheapRepeatNumThreshold {
		model.ExitNetValue = model.ExpensiveStockShortQuotePrice*model.ExpensiveStockEntryVolume - model.CheapStockLongQuotePrice*model.CheapStockEntryVolume
		if model.ExitNetValue+model.EntryNetValue >= model.MinProfitThreshold.Applied {
			return true
		}
	}
	return false
}

func GetTrimSignal(model *model.PairTradingModel) bool {
	if model.IsLongExpensiveStockShortCheapStock &&
		model.IsTrimmable {
		return true
	} else if model.IsShortExpensiveStockLongCheapStock &&
		model.IsTrimmable {
		return true
	}
	return false
}
