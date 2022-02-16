package pipeline

import (
	"time"

	"github.com/HaoxuanXu/TradingBot/db"
	"github.com/HaoxuanXu/TradingBot/internal/broker"
	"github.com/HaoxuanXu/TradingBot/strats/pairtrading/model"
	"github.com/HaoxuanXu/TradingBot/strats/pairtrading/transaction"
	"github.com/HaoxuanXu/TradingBot/tools/readwrite"
)

func EntryShortExpensiveLongCheap(model *model.PairTradingModel, broker *broker.AlpacaBroker) {
	entryValue := broker.MaxPortfolioPercent * broker.PortfolioValue
	// we try to use round lot here to imporve execution
	model.ExpensiveStockEntryVolume = float64(int((entryValue/2.0)/(model.ExpensiveStockShortQuotePrice*100)) * 100)
	if model.ExpensiveStockEntryVolume == 0 {
		model.ExpensiveStockEntryVolume = float64(int((entryValue / 2.0) / (model.ExpensiveStockShortQuotePrice)))
	}
	model.CheapStockEntryVolume = (model.ExpensiveStockEntryVolume * model.ExpensiveStockShortQuotePrice) / model.CheapStockLongQuotePrice
	go broker.SubmitOrderAsync(
		model.ExpensiveStockEntryVolume,
		model.ExpensiveStockSymbol,
		"sell",
		"market",
		"day",
		model.ExpensiveStockOrderChannel,
	)
	go broker.SubmitOrderAsync(
		model.CheapStockEntryVolume,
		model.CheapStockSymbol,
		"buy",
		"market",
		"day",
		model.CheapStockOrderChannel,
	)
	CheapStockOrder := <-model.CheapStockOrderChannel
	ExpensiveStockOrder := <-model.ExpensiveStockOrderChannel
	model.IsShortExpensiveStockLongCheapStock = true
	model.IsLongExpensiveStockShortCheapStock = false

	transaction.UpdateFieldsAfterTransaction(model, broker, CheapStockOrder, ExpensiveStockOrder)
	transaction.VetPosition(model)
	transaction.SlideRepeatAndPriceRatioArrays(model)
	transaction.RecordTransaction(model, broker)

	// Write the current data to disk
	WriteRecord(model)
}

func EntryLongExpensiveShortCheap(model *model.PairTradingModel, broker *broker.AlpacaBroker) {
	entryValue := broker.MaxPortfolioPercent * broker.PortfolioValue
	model.CheapStockEntryVolume = float64(int((entryValue/2.0)/(model.CheapStockShortQuotePrice*100)) * 100)
	if model.CheapStockEntryVolume == 0 {
		model.CheapStockEntryVolume = float64(int((entryValue / 2.0) / (model.CheapStockShortQuotePrice)))
	}
	model.ExpensiveStockEntryVolume = (model.CheapStockEntryVolume * model.CheapStockShortQuotePrice) / model.ExpensiveStockLongQuotePrice
	go broker.SubmitOrderAsync(
		model.CheapStockEntryVolume,
		model.CheapStockSymbol,
		"sell",
		"market",
		"day",
		model.CheapStockOrderChannel,
	)
	go broker.SubmitOrderAsync(
		model.ExpensiveStockEntryVolume,
		model.ExpensiveStockSymbol,
		"buy",
		"market",
		"day",
		model.ExpensiveStockOrderChannel,
	)
	CheapStockOrder := <-model.CheapStockOrderChannel
	ExpensiveStockOrder := <-model.ExpensiveStockOrderChannel
	model.IsLongExpensiveStockShortCheapStock = true
	model.IsShortExpensiveStockLongCheapStock = false

	transaction.UpdateFieldsAfterTransaction(model, broker, CheapStockOrder, ExpensiveStockOrder)
	transaction.VetPosition(model)
	transaction.SlideRepeatAndPriceRatioArrays(model)
	transaction.RecordTransaction(model, broker)

	// Write the current data to disk
	WriteRecord(model)
}

func ExitShortExpensiveLongCheap(model *model.PairTradingModel, broker *broker.AlpacaBroker) {
	go broker.SubmitOrderAsync(
		model.CheapStockEntryVolume,
		model.CheapStockSymbol,
		"sell",
		"market",
		"day",
		model.CheapStockOrderChannel,
	)
	go broker.SubmitOrderAsync(
		model.ExpensiveStockEntryVolume,
		model.ExpensiveStockSymbol,
		"buy",
		"market",
		"day",
		model.ExpensiveStockOrderChannel,
	)
	CheapStockOrder := <-model.CheapStockOrderChannel
	ExpensiveStockOrder := <-model.ExpensiveStockOrderChannel

	transaction.UpdateFieldsAfterTransaction(model, broker, CheapStockOrder, ExpensiveStockOrder)
	transaction.SlideRepeatAndPriceRatioArrays(model)
	transaction.RecordTransaction(model, broker)

	model.IsLongExpensiveStockShortCheapStock = false
	model.IsShortExpensiveStockLongCheapStock = false

	// Write the current data to disk
	WriteRecord(model)
}

func ExitLongExpensiveShortCheap(model *model.PairTradingModel, broker *broker.AlpacaBroker) {
	go broker.SubmitOrderAsync(
		model.ExpensiveStockEntryVolume,
		model.ExpensiveStockSymbol,
		"sell",
		"market",
		"day",
		model.ExpensiveStockOrderChannel,
	)
	go broker.SubmitOrderAsync(
		model.CheapStockEntryVolume,
		model.CheapStockSymbol,
		"buy",
		"market",
		"day",
		model.CheapStockOrderChannel,
	)
	CheapStockOrder := <-model.CheapStockOrderChannel
	ExpensiveStockOrder := <-model.ExpensiveStockOrderChannel

	transaction.UpdateFieldsAfterTransaction(model, broker, CheapStockOrder, ExpensiveStockOrder)
	transaction.SlideRepeatAndPriceRatioArrays(model)
	transaction.RecordTransaction(model, broker)

	model.IsLongExpensiveStockShortCheapStock = false
	model.IsShortExpensiveStockLongCheapStock = false

	// Write the current data to disk
	WriteRecord(model)
}

func UpdateSignalThresholds(model *model.PairTradingModel, broker *broker.AlpacaBroker, baseTime *time.Time, wrappingUp bool) {
	if time.Since(*baseTime) > time.Minute {
		transaction.SlideRepeatAndPriceRatioArrays(model)
		*baseTime = time.Now()
	}
	if wrappingUp {
		model.MinProfitThreshold = 0.0
	} else if time.Since(broker.LastTradeTime) > 15*time.Minute && time.Since(broker.LastTradeTime) < 16*time.Minute {
		model.MinProfitThreshold = model.CalculateMinProfitThreshold(1.0)
	} else if time.Since(broker.LastTradeTime) > 30*time.Minute {
		model.MinProfitThreshold = 0.0
	}
}

func WriteRecord(model *model.PairTradingModel) {
	shortExpensiveLongCheapPath, longExpensiveShortCheapPath,
		longExpensiveShortCheapRepeatNumsPath, shortExpensiveLongCheapRepeatNumsPath := db.MapRecordPath(model.StrategyAssetType)
	readwrite.WriteIntSlice(&model.LongExpensiveShortCheapRepeatArray, longExpensiveShortCheapRepeatNumsPath)
	readwrite.WriteIntSlice(&model.ShortExpensiveLongCheapRepeatArray, shortExpensiveLongCheapRepeatNumsPath)
	readwrite.WriteFloatSlice(&model.ShortExpensiveStockLongCheapStockPriceRatioRecord, shortExpensiveLongCheapPath)
	readwrite.WriteFloatSlice(&model.LongExpensiveStockShortCheapStockPriceRatioRecord, longExpensiveShortCheapPath)
}
