package pipeline

import (
	"log"
	"sync"

	"github.com/HaoxuanXu/TradingBot/db"
	"github.com/HaoxuanXu/TradingBot/internal/broker"
	"github.com/HaoxuanXu/TradingBot/strats/pairtrading/logging"
	"github.com/HaoxuanXu/TradingBot/strats/pairtrading/model"
	"github.com/HaoxuanXu/TradingBot/strats/pairtrading/transaction"
	"github.com/HaoxuanXu/TradingBot/tools/readwrite"
)

func EntryShortExpensiveLongCheap(model *model.PairTradingModel, broker *broker.AlpacaBroker, wg sync.WaitGroup) {
	entryValue := broker.MaxPortfolioPercent * broker.PortfolioValue
	model.ExpensiveStockEntryVolume = float64(int((entryValue / 2.0) / model.ExpensiveStockShortQuotePrice))
	model.CheapStockEntryVolume = (model.ExpensiveStockEntryVolume * model.ExpensiveStockShortQuotePrice) / model.CheapStockLongQuotePrice
	log.Println("start goroutine")
	wg.Add(2)
	go broker.SubmitOrderAsync(
		model.ExpensiveStockEntryVolume,
		model.ExpensiveStockSymbol,
		"sell",
		"market",
		"day",
		&wg,
	)
	go broker.SubmitOrderAsync(
		model.CheapStockEntryVolume,
		model.CheapStockSymbol,
		"buy",
		"market",
		"day",
		&wg,
	)
	wg.Wait()
	log.Println("goroutine complete")
	model.IsShortExpensiveStockLongCheapStock = true
	model.IsLongExpensiveStockShortCheapStock = false

	transaction.UpdateFieldsAfterTransaction(model, broker, broker.OrderMap[model.CheapStockSymbol], broker.OrderMap[model.ExpensiveStockSymbol])
	log.Println("begin logging")
	logging.LogTransaction(model, broker)
	transaction.VetPosition(broker, model)
	transaction.SlideRepeatAndPriceRatioArrays(model)
}

func EntryLongExpensiveShortCheap(model *model.PairTradingModel, broker *broker.AlpacaBroker, wg sync.WaitGroup) {
	entryValue := broker.MaxPortfolioPercent * broker.PortfolioValue
	model.CheapStockEntryVolume = float64(int((entryValue / 2.0) / model.CheapStockShortQuotePrice))
	model.ExpensiveStockEntryVolume = (model.CheapStockEntryVolume * model.CheapStockShortQuotePrice) / model.ExpensiveStockLongQuotePrice
	log.Println("start goroutine")
	wg.Add(2)
	go broker.SubmitOrderAsync(
		model.CheapStockEntryVolume,
		model.CheapStockSymbol,
		"sell",
		"market",
		"day",
		&wg,
	)
	go broker.SubmitOrderAsync(
		model.ExpensiveStockEntryVolume,
		model.ExpensiveStockSymbol,
		"buy",
		"market",
		"day",
		&wg,
	)
	wg.Wait()
	log.Println("goroutine complete")
	model.IsLongExpensiveStockShortCheapStock = true
	model.IsShortExpensiveStockLongCheapStock = false

	transaction.UpdateFieldsAfterTransaction(model, broker, broker.OrderMap[model.CheapStockSymbol], broker.OrderMap[model.ExpensiveStockSymbol])
	logging.LogTransaction(model, broker)
	transaction.VetPosition(broker, model)
	transaction.SlideRepeatAndPriceRatioArrays(model)
}

func ExitShortExpensiveLongCheap(model *model.PairTradingModel, broker *broker.AlpacaBroker, wg sync.WaitGroup) {
	log.Println("start goroutine")
	wg.Add(2)
	go broker.SubmitOrderAsync(
		model.CheapStockEntryVolume,
		model.CheapStockSymbol,
		"sell",
		"market",
		"day",
		&wg,
	)
	go broker.SubmitOrderAsync(
		model.ExpensiveStockEntryVolume,
		model.ExpensiveStockSymbol,
		"buy",
		"market",
		"day",
		&wg,
	)
	wg.Wait()
	log.Println("goroutine complete")
	model.IsLongExpensiveStockShortCheapStock = false
	model.IsShortExpensiveStockLongCheapStock = false

	transaction.UpdateFieldsAfterTransaction(model, broker, broker.OrderMap[model.CheapStockSymbol], broker.OrderMap[model.ExpensiveStockSymbol])
	logging.LogTransaction(model, broker)
	transaction.VetPosition(broker, model)
	transaction.SlideRepeatAndPriceRatioArrays(model)
}

func ExitLongExpensiveShortCheap(model *model.PairTradingModel, broker *broker.AlpacaBroker, wg sync.WaitGroup) {
	log.Println("start goroutine")
	wg.Add(2)
	go broker.SubmitOrderAsync(
		model.ExpensiveStockEntryVolume,
		model.ExpensiveStockSymbol,
		"sell",
		"market",
		"day",
		&wg,
	)
	go broker.SubmitOrderAsync(
		model.CheapStockEntryVolume,
		model.CheapStockSymbol,
		"buy",
		"market",
		"day",
		&wg,
	)
	wg.Wait()
	log.Println("goroutine complete")
	model.IsLongExpensiveStockShortCheapStock = false
	model.IsShortExpensiveStockLongCheapStock = false

	transaction.UpdateFieldsAfterTransaction(model, broker, broker.OrderMap[model.CheapStockSymbol], broker.OrderMap[model.ExpensiveStockSymbol])
	logging.LogTransaction(model, broker)
	transaction.VetPosition(broker, model)
	transaction.SlideRepeatAndPriceRatioArrays(model)

}

func WriteRecord(model *model.PairTradingModel, strat string) {
	shortExpensiveLongCheapPath, longExpensiveShortCheapPath, repeatNumsPath := db.MapRecordPath(strat)
	readwrite.WriteIntSlice(model.RepeatArray, repeatNumsPath)
	readwrite.WriteFloatSlice(model.ShortExpensiveStockLongCheapStockPriceRatioRecord, shortExpensiveLongCheapPath)
	readwrite.WriteFloatSlice(model.LongExpensiveStockShortCheapStockPriceRatioRecord, longExpensiveShortCheapPath)
}
