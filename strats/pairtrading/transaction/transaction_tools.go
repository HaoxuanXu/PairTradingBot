package transaction

import (
	"log"
	"math"

	"github.com/HaoxuanXu/TradingBot/strats/pairtrading/model"
	"github.com/HaoxuanXu/TradingBot/strats/pairtrading/updater"
	"github.com/HaoxuanXu/TradingBot/tools/repeater"
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
		util.UpdateIntSlice(&m.RepeatArray, m.ShortExpensiveStockLongCheapStockRepeatNumber)
		util.UpdateFloatSlice(&m.ShortExpensiveStockLongCheapStockPriceRatioRecord, m.ShortExpensiveStockLongCheapStockPreviousRatio)
		m.ShortExpensiveStockLongCheapStockRepeatNumber = 1
		m.ShortExpensiveStockLongCheapStockPreviousRatio = m.ShortExpensiveStockLongCheapStockPriceRatio
	}

	if m.LongExpensiveStockShortCheapStockPriceRatio == m.LongExpensiveStockShortCheapStockPreviousRatio {
		m.LongExpensiveStockShortCheapStockRepeatNumber++
	} else {
		util.UpdateIntSlice(&m.RepeatArray, m.LongExpensiveStockShortCheapStockRepeatNumber)
		util.UpdateFloatSlice(&m.LongExpensiveStockShortCheapStockPriceRatioRecord, m.LongExpensiveStockShortCheapStockPreviousRatio)
		m.LongExpensiveStockShortCheapStockRepeatNumber = 1
		m.LongExpensiveStockShortCheapStockPreviousRatio = m.LongExpensiveStockShortCheapStockPriceRatio
	}
}

func UpdateFieldsAfterTransaction(m *model.PairTradingModel, cheapStockOrder, expensiveStockOrder *alpaca.Order) {
	m.CheapStockFilledPrice = math.Abs(cheapStockOrder.FilledAvgPrice.InexactFloat64())
	m.CheapStockFilledQuantity = math.Abs(cheapStockOrder.FilledQty.InexactFloat64())
	m.ExpensiveStockFilledPrice = math.Abs(expensiveStockOrder.FilledAvgPrice.InexactFloat64())
	m.ExpensiveStockFilledQuantity = math.Abs(expensiveStockOrder.FilledQty.InexactFloat64())
	m.MinProfitThreshold = m.CalculateMinProfitThreshold(1.0)
	m.ExpensiveStockEntryVolume = math.Abs(m.ExpensiveStockFilledQuantity)
	m.CheapStockEntryVolume = math.Abs(m.CheapStockFilledQuantity)

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
		log.Println("minimum profit adjusted")
	}
}

func SlideRepeatAndPriceRatioArrays(model *model.PairTradingModel) {
	model.RepeatArray = windowslider.SlideWindowInt(model.RepeatArray, model.DefaultRepeatArrayLength)
	model.RepeatNumThreshold = repeater.CalculateOptimalRepeatNum(model.RepeatArray)

	model.LongExpensiveStockShortCheapStockPriceRatioRecord = windowslider.SlideWindowFloat(
		model.LongExpensiveStockShortCheapStockPriceRatioRecord,
		model.DefaultPriceRatioArrayLength,
	)
	model.ShortExpensiveStockLongCheapStockPriceRatioRecord = windowslider.SlideWindowFloat(
		model.ShortExpensiveStockLongCheapStockPriceRatioRecord,
		model.DefaultPriceRatioArrayLength,
	)
	model.PriceRatioThreshold = updater.UpdatePriceRatioThreshold(
		model.LongExpensiveStockShortCheapStockPriceRatioRecord,
		model.ShortExpensiveStockLongCheapStockPriceRatioRecord,
	)

	// log.Printf("The price threshold now is %f; the repeat threshold now is %d; minprofitthreshold is %f\n",
	// 	model.PriceRatioThreshold, model.RepeatNumThreshold, model.MinProfitThreshold)
}
