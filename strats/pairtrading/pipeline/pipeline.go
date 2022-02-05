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
	model.ExpensiveStockEntryVolume = float64(int((entryValue / 2.0) / model.ExpensiveStockShortQuotePrice))
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

	transaction.UpdateFieldsAfterTransaction(model, CheapStockOrder, ExpensiveStockOrder)
	transaction.VetPosition(model)
	transaction.SlideRepeatAndPriceRatioArrays(model)
	transaction.RecordTransaction(model, broker)
}

func EntryLongExpensiveShortCheap(model *model.PairTradingModel, broker *broker.AlpacaBroker) {
	entryValue := broker.MaxPortfolioPercent * broker.PortfolioValue
	model.CheapStockEntryVolume = float64(int((entryValue / 2.0) / model.CheapStockShortQuotePrice))
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

	transaction.UpdateFieldsAfterTransaction(model, CheapStockOrder, ExpensiveStockOrder)
	transaction.VetPosition(model)
	transaction.SlideRepeatAndPriceRatioArrays(model)
	transaction.RecordTransaction(model, broker)
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

	transaction.UpdateFieldsAfterTransaction(model, CheapStockOrder, ExpensiveStockOrder)
	transaction.SlideRepeatAndPriceRatioArrays(model)
	transaction.RecordTransaction(model, broker)

	model.IsLongExpensiveStockShortCheapStock = false
	model.IsShortExpensiveStockLongCheapStock = false
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

	transaction.UpdateFieldsAfterTransaction(model, CheapStockOrder, ExpensiveStockOrder)
	transaction.SlideRepeatAndPriceRatioArrays(model)
	transaction.RecordTransaction(model, broker)

	model.IsLongExpensiveStockShortCheapStock = false
	model.IsShortExpensiveStockLongCheapStock = false
}

func UpdateSignalThresholds(model *model.PairTradingModel, baseTime *time.Time) {
	if time.Since(*baseTime) > time.Minute {
		transaction.SlideRepeatAndPriceRatioArrays(model)
		*baseTime = time.Now()
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
