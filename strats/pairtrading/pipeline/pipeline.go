package pipeline

import (
	"sync"

	"github.com/HaoxuanXu/TradingBot/db"
	"github.com/HaoxuanXu/TradingBot/internal/broker"
	"github.com/HaoxuanXu/TradingBot/strats/pairtrading/logging"
	"github.com/HaoxuanXu/TradingBot/strats/pairtrading/model"
	"github.com/HaoxuanXu/TradingBot/strats/pairtrading/transaction"
	"github.com/HaoxuanXu/TradingBot/tools/readwrite"
	"github.com/alpacahq/alpaca-trade-api-go/v2/alpaca"
)

func EntryShortExpensiveLongCheap(model *model.PairTradingModel, broker *broker.AlpacaBroker, wg *sync.WaitGroup) {
	expensiveStockOrderChannel := make(chan *alpaca.Order)
	cheapStockOrderChannel := make(chan *alpaca.Order)
	entryValue := broker.MaxPortfolioPercent * broker.PortfolioValue
	model.ExpensiveStockEntryVolume = float64(int((entryValue / 2.0) / model.ExpensiveStockShortQuotePrice))
	model.CheapStockEntryVolume = (model.ExpensiveStockEntryVolume * model.ExpensiveStockShortQuotePrice) / model.CheapStockLongQuotePrice
	wg.Add(2)
	go broker.SubmitOrderAsync(
		model.ExpensiveStockEntryVolume,
		model.ExpensiveStockSymbol,
		"sell",
		"market",
		"day",
		expensiveStockOrderChannel,
		wg,
	)
	go broker.SubmitOrderAsync(
		model.CheapStockEntryVolume,
		model.CheapStockSymbol,
		"buy",
		"market",
		"day",
		cheapStockOrderChannel,
		wg,
	)
	wg.Wait()

	model.IsShortExpensiveStockLongCheapStock = true
	model.IsLongExpensiveStockShortCheapStock = false 

	transaction.UpdateFieldsAfterTransaction(model, broker, <- cheapStockOrderChannel, <- expensiveStockOrderChannel)
	logging.LogTransaction(model, broker)
	transaction.VetPosition(broker, model)
	transaction.SlideRepeatAndPriceRatioArrays(model)

	close(expensiveStockOrderChannel)
	close(cheapStockOrderChannel)
}

func EntryLongExpensiveShortCheap(model *model.PairTradingModel, broker *broker.AlpacaBroker, wg *sync.WaitGroup) {
	expensiveStockOrderChannel := make(chan *alpaca.Order)
	cheapStockOrderChannel := make(chan *alpaca.Order)
	entryValue := broker.MaxPortfolioPercent * broker.PortfolioValue
	model.CheapStockEntryVolume = float64(int((entryValue / 2.0) / model.CheapStockShortQuotePrice))
	model.ExpensiveStockEntryVolume = (model.CheapStockEntryVolume * model.CheapStockShortQuotePrice) / model.ExpensiveStockLongQuotePrice
	wg.Add(2)
	go broker.SubmitOrderAsync(
		model.CheapStockEntryVolume,
		model.CheapStockSymbol,
		"sell",
		"market",
		"day",
		cheapStockOrderChannel,
		wg,
	)
	go broker.SubmitOrderAsync(
		model.ExpensiveStockEntryVolume,
		model.ExpensiveStockSymbol,
		"buy",
		"market",
		"day",
		expensiveStockOrderChannel,
		wg,
	)
	wg.Wait()

	model.IsLongExpensiveStockShortCheapStock = true
	model.IsShortExpensiveStockLongCheapStock = false 

	transaction.UpdateFieldsAfterTransaction(model, broker, <- cheapStockOrderChannel, <- expensiveStockOrderChannel)
	logging.LogTransaction(model, broker)
	transaction.VetPosition(broker, model)
	transaction.SlideRepeatAndPriceRatioArrays(model)

	close(cheapStockOrderChannel)
	close(expensiveStockOrderChannel)
}


func ExitShortExpensiveLongCheap(model *model.PairTradingModel, broker *broker.AlpacaBroker, wg *sync.WaitGroup) {
	expensiveStockOrderChannel := make(chan *alpaca.Order)
	cheapStockOrderChannel := make(chan *alpaca.Order)
	wg.Add(2)
	go broker.SubmitOrderAsync(
		model.CheapStockEntryVolume,
		model.CheapStockSymbol,
		"sell",
		"market",
		"day",
		cheapStockOrderChannel,
		wg,
	)
	go broker.SubmitOrderAsync(
		model.ExpensiveStockEntryVolume,
		model.ExpensiveStockSymbol,
		"buy",
		"market",
		"day",
		expensiveStockOrderChannel,
		wg,
	)
	wg.Wait()

	model.IsLongExpensiveStockShortCheapStock = false
	model.IsShortExpensiveStockLongCheapStock = false 

	transaction.UpdateFieldsAfterTransaction(model, broker, <- cheapStockOrderChannel, <- expensiveStockOrderChannel)
	logging.LogTransaction(model, broker)
	transaction.VetPosition(broker, model)
	transaction.SlideRepeatAndPriceRatioArrays(model)

	close(cheapStockOrderChannel)
	close(expensiveStockOrderChannel)
}


func ExitLongExpensiveShortCheap(model *model.PairTradingModel, broker *broker.AlpacaBroker, wg *sync.WaitGroup) {
	expensiveStockOrderChannel := make(chan *alpaca.Order)
	cheapStockOrderChannel := make(chan *alpaca.Order)
	wg.Add(2)
	go broker.SubmitOrderAsync(
		model.ExpensiveStockEntryVolume,
		model.ExpensiveStockSymbol,
		"sell",
		"market",
		"day",
		expensiveStockOrderChannel,
		wg,
	)
	go broker.SubmitOrderAsync(
		model.CheapStockEntryVolume,
		model.CheapStockSymbol,
		"buy",
		"market",
		"day",
		cheapStockOrderChannel,
		wg,
	)
	wg.Wait()

	model.IsLongExpensiveStockShortCheapStock = false
	model.IsShortExpensiveStockLongCheapStock = false 

	transaction.UpdateFieldsAfterTransaction(model, broker, <- cheapStockOrderChannel, <- expensiveStockOrderChannel)
	logging.LogTransaction(model, broker)
	transaction.VetPosition(broker, model)
	transaction.SlideRepeatAndPriceRatioArrays(model)

	close(cheapStockOrderChannel)
	close(expensiveStockOrderChannel)

}


func WriteRecord(model *model.PairTradingModel, strat string) {
	shortExpensiveLongCheapPath, longExpensiveShortCheapPath, repeatNumsPath := db.MapRecordPath(strat)
	readwrite.WriteIntSlice(model.RepeatArray, repeatNumsPath)
	readwrite.WriteFloatSlice(model.ShortExpensiveStockLongCheapStockPriceRatioRecord, shortExpensiveLongCheapPath)
	readwrite.WriteFloatSlice(model.LongExpensiveStockShortCheapStockPriceRatioRecord, longExpensiveShortCheapPath)
}