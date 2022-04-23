package transaction

import (
	"log"
	"math"
	"time"

	"github.com/HaoxuanXu/TradingBot/internal/broker"
	"github.com/HaoxuanXu/TradingBot/strats/pairtrading/model"
	"github.com/HaoxuanXu/TradingBot/tools/util"
	"github.com/HaoxuanXu/TradingBot/tools/windowslider"
	"github.com/alpacahq/alpaca-trade-api-go/v2/alpaca"
)

func UpdateFieldsFromQuotes(m *model.PairTradingModel) {
	m.ShortExpensiveStockLongCheapStockPriceRatio = m.ExpensiveStockShortQuotePrice / m.CheapStockLongQuotePrice
	m.LongExpensiveStockShortCheapStockPriceRatio = m.ExpensiveStockLongQuotePrice / m.CheapStockShortQuotePrice

	if m.ShortExpensiveStockLongCheapStockPriceRatio == m.ShortExpensiveStockLongCheapStockPreviousRatio {
		m.ShortExpensiveStockLongCheapStockRepeatNumber++
	} else {
		util.UpdateIntSlice(&m.ShortExpensiveLongCheapRepeatArray, m.ShortExpensiveStockLongCheapStockRepeatNumber)
		m.ShortExpensiveStockLongCheapStockPreviousRepeatNumber = m.ShortExpensiveStockLongCheapStockRepeatNumber
		m.ShortExpensiveStockLongCheapStockRepeatNumber = 1
		m.ShortExpensiveStockLongCheapStockPreviousRatio = m.ShortExpensiveStockLongCheapStockPriceRatio
		util.UpdateFloatSlice(&m.ShortExpensiveStockLongCheapStockPriceRatioRecord, m.ShortExpensiveStockLongCheapStockPreviousRatio)
	}

	if m.LongExpensiveStockShortCheapStockPriceRatio == m.LongExpensiveStockShortCheapStockPreviousRatio {
		m.LongExpensiveStockShortCheapStockRepeatNumber++
	} else {
		util.UpdateIntSlice(&m.LongExpensiveShortCheapRepeatArray, m.LongExpensiveStockShortCheapStockRepeatNumber)
		m.LongExpensiveStockShortCheapStockPreviousRepeatNumber = m.LongExpensiveStockShortCheapStockRepeatNumber
		m.LongExpensiveStockShortCheapStockRepeatNumber = 1
		m.LongExpensiveStockShortCheapStockPreviousRatio = m.LongExpensiveStockShortCheapStockPriceRatio
		util.UpdateFloatSlice(&m.LongExpensiveStockShortCheapStockPriceRatioRecord, m.LongExpensiveStockShortCheapStockPreviousRatio)
	}
}

func UpdateFieldsAfterTransaction(m *model.PairTradingModel, broker *broker.AlpacaBroker, cheapStockOrder, expensiveStockOrder *alpaca.Order) {
	m.CheapStockFilledPrice = math.Abs(cheapStockOrder.FilledAvgPrice.InexactFloat64())
	m.CheapStockFilledQuantity = math.Abs(cheapStockOrder.FilledQty.InexactFloat64())
	m.ExpensiveStockFilledPrice = math.Abs(expensiveStockOrder.FilledAvgPrice.InexactFloat64())
	m.ExpensiveStockFilledQuantity = math.Abs(expensiveStockOrder.FilledQty.InexactFloat64())
	m.MinProfitThreshold = m.CalculateMinProfitThreshold(1.0)
	m.ExpensiveStockEntryVolume = math.Abs(m.ExpensiveStockFilledQuantity)
	m.CheapStockEntryVolume = math.Abs(m.CheapStockFilledQuantity)
	broker.LastTradeTime = time.Now()
}

func VetPosition(model *model.PairTradingModel) {
	var longPosition float64
	var shortPosition float64

	if model.IsLongExpensiveStockShortCheapStock {
		longPosition = model.ExpensiveStockFilledPrice * model.ExpensiveStockFilledQuantity
		shortPosition = model.CheapStockFilledPrice * model.CheapStockFilledQuantity
	} else {
		longPosition = model.CheapStockFilledPrice * model.CheapStockFilledQuantity
		shortPosition = model.ExpensiveStockFilledPrice * model.ExpensiveStockFilledQuantity
	}

	overboughtPercent := (longPosition - shortPosition) / longPosition

	if overboughtPercent > 0.00003 {
		model.MinProfitThreshold = model.MinProfitThreshold - (longPosition - shortPosition)
		model.IsMinProfitAdjusted = true
		log.Println("minimum profit adjusted")
	}
}

func SlideRepeatAndPriceRatioArrays(model *model.PairTradingModel) {
	model.LongExpensiveShortCheapRepeatArray = windowslider.SlideWindowInt(
		model.LongExpensiveShortCheapRepeatArray,
		model.DefaultRepeatArrayLength,
	)
	model.ShortExpensiveLongCheapRepeatArray = windowslider.SlideWindowInt(
		model.ShortExpensiveLongCheapRepeatArray,
		model.DefaultRepeatArrayLength,
	)
	model.LongExpensiveStockShortCheapStockPriceRatioRecord = windowslider.SlideWindowFloat(
		model.LongExpensiveStockShortCheapStockPriceRatioRecord,
		model.DefaultPriceRatioArrayLength,
	)
	model.ShortExpensiveStockLongCheapStockPriceRatioRecord = windowslider.SlideWindowFloat(
		model.ShortExpensiveStockLongCheapStockPriceRatioRecord,
		model.DefaultPriceRatioArrayLength,
	)
	model.UpdateParameters()
}

func RecordTransaction(model *model.PairTradingModel, broker *broker.AlpacaBroker) {
	if !broker.HasPosition {
		if model.IsLongExpensiveStockShortCheapStock {
			model.EntryNetValue = math.Abs(model.CheapStockFilledPrice*model.CheapStockFilledQuantity) - math.Abs(model.ExpensiveStockFilledPrice*model.ExpensiveStockFilledQuantity)
			log.Printf("long %s: %f shares; short %s: %f shares -- (current long repeat num: %d, previous long repeat num: %d, timestamp diffs: %f)\n",
				model.ExpensiveStockSymbol,
				model.ExpensiveStockEntryVolume,
				model.CheapStockSymbol,
				model.CheapStockEntryVolume,
				model.LongExpensiveStockShortCheapStockRepeatNumber,
				model.LongExpensiveStockShortCheapStockPreviousRepeatNumber,
				model.QuoteTimestampDifferenceMilliseconds,
			)
		} else {
			model.EntryNetValue = math.Abs(model.ExpensiveStockFilledPrice*model.ExpensiveStockFilledQuantity) - math.Abs(model.CheapStockFilledPrice*model.CheapStockFilledQuantity)
			log.Printf("short %s: %f shares; long %s: %f shares -- (current short repeat num: %d, previous short repeat num: %d, timestamp diffs: %f)\n",
				model.ExpensiveStockSymbol,
				model.ExpensiveStockEntryVolume,
				model.CheapStockSymbol,
				model.CheapStockEntryVolume,
				model.ShortExpensiveStockLongCheapStockRepeatNumber,
				model.ShortExpensiveStockLongCheapStockPreviousRepeatNumber,
				model.QuoteTimestampDifferenceMilliseconds,
			)
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
		log.Printf("position closed. Presumed Profit: $%f. Actual Profit: $%f -- (cur long repeat: %d, cur short repeat: %d) timestamp diffs: %f \n",
			presumedProfit,
			actualProfit,
			model.LongExpensiveStockShortCheapStockRepeatNumber,
			model.ShortExpensiveStockLongCheapStockRepeatNumber,
			model.QuoteTimestampDifferenceMilliseconds,
		)
		broker.HasPosition = false
		broker.TransactionNums++
	}
}

func CheckExistingPositions(model *model.PairTradingModel, broker *broker.AlpacaBroker) {
	expensiveStockPosition := broker.GetPosition(model.ExpensiveStockSymbol)
	cheapStockPosition := broker.GetPosition(model.CheapStockSymbol)

	if expensiveStockPosition != nil && cheapStockPosition != nil {
		broker.HasPosition = true
		if expensiveStockPosition.Side == "long" {
			model.IsLongExpensiveStockShortCheapStock = true
			model.IsShortExpensiveStockLongCheapStock = false
		} else {
			model.IsShortExpensiveStockLongCheapStock = true
			model.IsLongExpensiveStockShortCheapStock = false
		}

		model.CheapStockFilledPrice = cheapStockPosition.EntryPrice.Abs().InexactFloat64()
		model.CheapStockFilledQuantity = cheapStockPosition.Qty.Abs().InexactFloat64()
		model.ExpensiveStockFilledPrice = expensiveStockPosition.EntryPrice.Abs().InexactFloat64()
		model.ExpensiveStockFilledQuantity = expensiveStockPosition.Qty.Abs().InexactFloat64()
		model.MinProfitThreshold = model.CalculateMinProfitThreshold(1.0)
		model.ExpensiveStockEntryVolume = math.Abs(model.ExpensiveStockFilledQuantity)
		model.CheapStockEntryVolume = math.Abs(model.CheapStockFilledQuantity)

		VetPosition(model)
	}
}
