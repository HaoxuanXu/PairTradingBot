package pipeline

import (
	"github.com/HaoxuanXu/TradingBot/db"
	"github.com/HaoxuanXu/TradingBot/internal/broker"
	"github.com/HaoxuanXu/TradingBot/strats/pairtrading/logging"
	"github.com/HaoxuanXu/TradingBot/strats/pairtrading/model"
	"github.com/HaoxuanXu/TradingBot/strats/pairtrading/transaction"
	"github.com/HaoxuanXu/TradingBot/tools/readwrite"
	"github.com/alpacahq/alpaca-trade-api-go/v2/alpaca"
)

func EntryShortExpensiveLongCheap(model *model.PairTradingModel, broker *broker.AlpacaBroker) {
	ExpensiveStockOrderChannel := make(chan *alpaca.Order)
	CheapStockOrderChannel := make(chan *alpaca.Order)
	entryValue := broker.MaxPortfolioPercent * broker.PortfolioValue
	model.ExpensiveStockEntryVolume = float64(int((entryValue / 2.0) / model.ExpensiveStockShortQuotePrice))
	model.CheapStockEntryVolume = (model.ExpensiveStockEntryVolume * model.ExpensiveStockShortQuotePrice) / model.CheapStockLongQuotePrice
	go broker.SubmitOrderAsync(
		model.ExpensiveStockEntryVolume,
		model.ExpensiveStockSymbol,
		"sell",
		"market",
		"day",
		ExpensiveStockOrderChannel,
	)
	go broker.SubmitOrderAsync(
		model.CheapStockEntryVolume,
		model.CheapStockSymbol,
		"buy",
		"market",
		"day",
		CheapStockOrderChannel,
	)
	CheapStockOrder := <-CheapStockOrderChannel
	ExpensiveStockOrder := <-ExpensiveStockOrderChannel
	model.IsShortExpensiveStockLongCheapStock = true
	model.IsLongExpensiveStockShortCheapStock = false

	transaction.UpdateFieldsAfterTransaction(model, broker, CheapStockOrder, ExpensiveStockOrder)
	logging.LogTransaction(model, broker)
	transaction.VetPosition(broker, model)
	transaction.SlideRepeatAndPriceRatioArrays(model)

	close(ExpensiveStockOrderChannel)
	close(CheapStockOrderChannel)
}

func EntryLongExpensiveShortCheap(model *model.PairTradingModel, broker *broker.AlpacaBroker) {
	ExpensiveStockOrderChannel := make(chan *alpaca.Order)
	CheapStockOrderChannel := make(chan *alpaca.Order)
	entryValue := broker.MaxPortfolioPercent * broker.PortfolioValue
	model.CheapStockEntryVolume = float64(int((entryValue / 2.0) / model.CheapStockShortQuotePrice))
	model.ExpensiveStockEntryVolume = (model.CheapStockEntryVolume * model.CheapStockShortQuotePrice) / model.ExpensiveStockLongQuotePrice
	go broker.SubmitOrderAsync(
		model.CheapStockEntryVolume,
		model.CheapStockSymbol,
		"sell",
		"market",
		"day",
		CheapStockOrderChannel,
	)
	go broker.SubmitOrderAsync(
		model.ExpensiveStockEntryVolume,
		model.ExpensiveStockSymbol,
		"buy",
		"market",
		"day",
		ExpensiveStockOrderChannel,
	)
	CheapStockOrder := <-CheapStockOrderChannel
	ExpensiveStockOrder := <-ExpensiveStockOrderChannel
	model.IsLongExpensiveStockShortCheapStock = true
	model.IsShortExpensiveStockLongCheapStock = false

	transaction.UpdateFieldsAfterTransaction(model, broker, CheapStockOrder, ExpensiveStockOrder)
	logging.LogTransaction(model, broker)
	transaction.VetPosition(broker, model)
	transaction.SlideRepeatAndPriceRatioArrays(model)

	close(ExpensiveStockOrderChannel)
	close(CheapStockOrderChannel)
}

func ExitShortExpensiveLongCheap(model *model.PairTradingModel, broker *broker.AlpacaBroker) {
	ExpensiveStockOrderChannel := make(chan *alpaca.Order)
	CheapStockOrderChannel := make(chan *alpaca.Order)
	go broker.SubmitOrderAsync(
		model.CheapStockEntryVolume,
		model.CheapStockSymbol,
		"sell",
		"market",
		"day",
		CheapStockOrderChannel,
	)
	go broker.SubmitOrderAsync(
		model.ExpensiveStockEntryVolume,
		model.ExpensiveStockSymbol,
		"buy",
		"market",
		"day",
		ExpensiveStockOrderChannel,
	)
	CheapStockOrder := <-CheapStockOrderChannel
	ExpensiveStockOrder := <-ExpensiveStockOrderChannel
	model.IsLongExpensiveStockShortCheapStock = false
	model.IsShortExpensiveStockLongCheapStock = false

	transaction.UpdateFieldsAfterTransaction(model, broker, CheapStockOrder, ExpensiveStockOrder)
	logging.LogTransaction(model, broker)
	transaction.VetPosition(broker, model)
	transaction.SlideRepeatAndPriceRatioArrays(model)

	close(ExpensiveStockOrderChannel)
	close(CheapStockOrderChannel)
}

func ExitLongExpensiveShortCheap(model *model.PairTradingModel, broker *broker.AlpacaBroker) {
	ExpensiveStockOrderChannel := make(chan *alpaca.Order)
	CheapStockOrderChannel := make(chan *alpaca.Order)
	go broker.SubmitOrderAsync(
		model.ExpensiveStockEntryVolume,
		model.ExpensiveStockSymbol,
		"sell",
		"market",
		"day",
		ExpensiveStockOrderChannel,
	)
	go broker.SubmitOrderAsync(
		model.CheapStockEntryVolume,
		model.CheapStockSymbol,
		"buy",
		"market",
		"day",
		CheapStockOrderChannel,
	)
	CheapStockOrder := <-CheapStockOrderChannel
	ExpensiveStockOrder := <-ExpensiveStockOrderChannel
	model.IsLongExpensiveStockShortCheapStock = false
	model.IsShortExpensiveStockLongCheapStock = false

	transaction.UpdateFieldsAfterTransaction(model, broker, CheapStockOrder, ExpensiveStockOrder)
	logging.LogTransaction(model, broker)
	transaction.VetPosition(broker, model)
	transaction.SlideRepeatAndPriceRatioArrays(model)

	close(ExpensiveStockOrderChannel)
	close(CheapStockOrderChannel)
}

func WriteRecord(model *model.PairTradingModel, strat string) {
	shortExpensiveLongCheapPath, longExpensiveShortCheapPath, repeatNumsPath := db.MapRecordPath(strat)
	readwrite.WriteIntSlice(model.RepeatArray, repeatNumsPath)
	readwrite.WriteFloatSlice(model.ShortExpensiveStockLongCheapStockPriceRatioRecord, shortExpensiveLongCheapPath)
	readwrite.WriteFloatSlice(model.LongExpensiveStockShortCheapStockPriceRatioRecord, longExpensiveShortCheapPath)
}
