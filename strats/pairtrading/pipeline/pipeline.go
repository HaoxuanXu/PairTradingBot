package pipeline

import (
	"log"
	"time"

	"github.com/HaoxuanXu/TradingBot/db"
	"github.com/HaoxuanXu/TradingBot/internal/broker"
	"github.com/HaoxuanXu/TradingBot/strats/pairtrading/model"
	"github.com/HaoxuanXu/TradingBot/strats/pairtrading/transaction"
	"github.com/HaoxuanXu/TradingBot/tools/readwrite"
	"github.com/HaoxuanXu/TradingBot/tools/util"
	"github.com/alpacahq/alpaca-trade-api-go/v2/alpaca"
)

func EntryShortExpensiveLongCheap(model *model.PairTradingModel, broker *broker.AlpacaBroker, assetParams *db.AssetParamConfig) {
	entryValue := broker.MaxPortfolioPercent * broker.PortfolioValue
	// we try to use round lot here to imporve execution
	model.ExpensiveStockEntryVolume = float64(int((entryValue / 2.0) / (model.ExpensiveStockShortQuotePrice)))
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
	WriteRecord(model, assetParams)
}

func EntryLongExpensiveShortCheap(model *model.PairTradingModel, broker *broker.AlpacaBroker, assetParams *db.AssetParamConfig) {
	entryValue := broker.MaxPortfolioPercent * broker.PortfolioValue
	model.CheapStockEntryVolume = float64(int((entryValue / 2.0) / (model.CheapStockShortQuotePrice)))
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
	WriteRecord(model, assetParams)
}

func ExitShortExpensiveLongCheap(model *model.PairTradingModel, broker *broker.AlpacaBroker, assetParams *db.AssetParamConfig) {
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
	model.IsMinProfitAdjusted = false
	model.IsTrimmable = false
	model.TrimmedAmount = 0.0

	// Write the current data to disk
	WriteRecord(model, assetParams)
}

func ExitLongExpensiveShortCheap(model *model.PairTradingModel, broker *broker.AlpacaBroker, assetParams *db.AssetParamConfig) {
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
	model.IsMinProfitAdjusted = false
	model.IsTrimmable = false
	model.TrimmedAmount = 0.0

	// Write the current data to disk
	WriteRecord(model, assetParams)
}

func TrimPosition(model *model.PairTradingModel, broker *broker.AlpacaBroker, assetParams *db.AssetParamConfig) {
	var qty float64
	var profit float64
	var order *alpaca.Order
	if model.IsLongExpensiveStockShortCheapStock {
		qty = model.TrimmedAmount / model.ExpensiveStockFilledPrice
		order = broker.SubmitOrder(
			qty,
			model.ExpensiveStockSymbol,
			"sell",
			"market",
			"day",
		)
		profit = qty * (order.FilledAvgPrice.InexactFloat64() - model.ExpensiveStockFilledPrice)
	} else if model.IsShortExpensiveStockLongCheapStock {
		qty = model.TrimmedAmount / model.CheapStockFilledPrice
		order = broker.SubmitOrder(
			qty,
			model.CheapStockSymbol,
			"sell",
			"market",
			"day",
		)
		profit = qty * (order.FilledAvgPrice.InexactFloat64() - model.CheapStockFilledPrice)
	}

	time.Sleep(5 * time.Second) // make sure the positions are updated
	model.ExpensiveStockFilledQuantity = broker.GetPosition(model.ExpensiveStockSymbol).Qty.Abs().InexactFloat64()
	model.CheapStockFilledQuantity = broker.GetPosition(model.CheapStockSymbol).Qty.Abs().InexactFloat64()
	model.ExpensiveStockEntryVolume = model.ExpensiveStockFilledQuantity
	model.CheapStockEntryVolume = model.CheapStockFilledQuantity

	transaction.VetPosition(model)

	log.Printf("Position successfully trimmed. Trimming Profit: $%f", profit)

	WriteRecord(model, assetParams)

}

func UpdateSignalThresholds(model *model.PairTradingModel, broker *broker.AlpacaBroker, counter *util.Counter, wrappingUp bool, assetParams *db.AssetParamConfig) {
	if time.Since(counter.BaseTime) > time.Minute {
		transaction.SlideRepeatAndPriceRatioArrays(model)
		WriteRecord(model, assetParams)
		model.UpdateParameters()
		counter.BaseTime = time.Now()
		counter.Incrementer++
	}
	if counter.Incrementer == 1 {
		counter.RefreshIncrementer()
	}
	if model.IsMinProfitAdjusted {
		return
	}
	if wrappingUp {
		model.MinProfitThreshold.Applied = 0.0
	} else if time.Since(broker.LastTradeTime) > 10*time.Minute {
		model.MinProfitThreshold.Applied = model.MinProfitThreshold.Low
	} else if time.Since(broker.LastTradeTime) > 15*time.Minute {
		model.MinProfitThreshold.Applied = 0
	}
}

func WriteRecord(model *model.PairTradingModel, assetParams *db.AssetParamConfig) {
	readwrite.WriteIntSlice(&model.LongExpensiveShortCheapRepeatArray, assetParams.LongExpensiveShortCheapRepeatNumPath)
	readwrite.WriteIntSlice(&model.ShortExpensiveLongCheapRepeatArray, assetParams.ShortExpensiveLongCheapRepeatNumPath)
	readwrite.WriteFloatSlice(&model.ShortExpensiveStockLongCheapStockPriceRatioRecord, assetParams.ShortExensiveLongCheapPriceRatioPath)
	readwrite.WriteFloatSlice(&model.LongExpensiveStockShortCheapStockPriceRatioRecord, assetParams.LongExpensiveShortCheapPriceRatioPath)
}
