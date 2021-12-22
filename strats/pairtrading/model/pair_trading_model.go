package model

import (
	"github.com/HaoxuanXu/TradingBot/strats/pairtrading/updater"
	"github.com/HaoxuanXu/TradingBot/tools/readwrite"
	"github.com/HaoxuanXu/TradingBot/tools/repeater"
	"github.com/alpacahq/alpaca-trade-api-go/v2/alpaca"
)

type PairTradingModel struct {
	ExpensiveStockSymbol                              string
	CheapStockSymbol                                  string
	EntryNetValue                                     float64
	ExitNetValue                                      float64
	LoserNums                                         int
	MinProfitThreshold                                float64
	PriceRatioThreshold                               float64
	CheapStockEntryVolume                             float64
	ExpensiveStockEntryVolume                         float64
	ExpensiveStockFilledQuantity                      float64
	CheapStockFilledQuantity                          float64
	ExpensiveStockFilledPrice                         float64
	CheapStockFilledPrice                             float64
	ExpensiveStockShortQuotePrice                     float64
	ExpensiveStockLongQuotePrice                      float64
	CheapStockShortQuotePrice                         float64
	CheapStockLongQuotePrice                          float64
	IsShortExpensiveStockLongCheapStock               bool
	IsLongExpensiveStockShortCheapStock               bool
	ShortExpensiveStockLongCheapStockPriceRatio       float64
	LongExpensiveStockShortCheapStockPriceRatio       float64
	ShortExpensiveStockLongCheapStockPreviousRatio    float64
	LongExpensiveStockShortCheapStockPreviousRatio    float64
	ShortExpensiveStockLongCheapStockRepeatNumber     int
	LongExpensiveStockShortCheapStockRepeatNumber     int
	ShortExpensiveStockLongCheapStockPriceRatioRecord []float64
	LongExpensiveStockShortCheapStockPriceRatioRecord []float64
	RepeatArray                                       []int
	RepeatNumThreshold                                int
	DefaultRepeatArrayLength                          int
	DefaultPriceRatioArrayLength                      int
	ExpensiveStockOrderChannel                        chan *alpaca.Order
	CheapStockOrderChannel                            chan *alpaca.Order
}

func (model *PairTradingModel) getStockSymbols(assetType string) (string, string) {
	if assetType == "gold" {
		return "GLD", "IAU"
	}
	return "", ""
}

func GetModel(assetType, shortLongPath, longShortPath, repeatNumPath string) *PairTradingModel {
	dataModel := &PairTradingModel{}
	dataModel.initialize(assetType, shortLongPath, longShortPath, repeatNumPath)
	return dataModel
}

func (model *PairTradingModel) CalculateMinProfitThreshold(baseNum float64) float64 {
	return baseNum * (model.ExpensiveStockFilledPrice*model.ExpensiveStockFilledQuantity + model.CheapStockFilledPrice*model.CheapStockFilledQuantity) / 120000
}

func (model *PairTradingModel) initialize(assetType, shortLongPath, longShortPath, repeatNumPath string) {
	model.ExpensiveStockSymbol, model.CheapStockSymbol = model.getStockSymbols(assetType)
	model.ShortExpensiveStockLongCheapStockPriceRatioRecord = readwrite.ReadRecordFloat(shortLongPath)
	model.LongExpensiveStockShortCheapStockPriceRatioRecord = readwrite.ReadRecordFloat(longShortPath)
	model.RepeatArray = readwrite.ReadRecordInt(repeatNumPath)
	model.ShortExpensiveStockLongCheapStockRepeatNumber = 0
	model.LongExpensiveStockShortCheapStockRepeatNumber = 0
	model.LongExpensiveStockShortCheapStockPriceRatio = 0.0
	model.ShortExpensiveStockLongCheapStockPriceRatio = 0.0
	model.LongExpensiveStockShortCheapStockPreviousRatio = 0.0
	model.ShortExpensiveStockLongCheapStockPreviousRatio = 0.0
	model.IsLongExpensiveStockShortCheapStock = false
	model.IsShortExpensiveStockLongCheapStock = false
	model.CheapStockLongQuotePrice = 0.0
	model.CheapStockShortQuotePrice = 0.0
	model.ExpensiveStockLongQuotePrice = 0.0
	model.ExpensiveStockShortQuotePrice = 0.0
	model.CheapStockFilledPrice = 0.0
	model.ExpensiveStockFilledPrice = 0.0
	model.CheapStockFilledQuantity = 0.0
	model.ExpensiveStockFilledQuantity = 0.0
	model.ExpensiveStockEntryVolume = 0.0
	model.CheapStockEntryVolume = 0.0
	model.PriceRatioThreshold = updater.UpdatePriceRatioThreshold(
		model.LongExpensiveStockShortCheapStockPriceRatioRecord,
		model.ShortExpensiveStockLongCheapStockPriceRatioRecord,
	)
	model.RepeatNumThreshold = repeater.CalculateOptimalRepeatNum(model.RepeatArray)
	model.DefaultRepeatArrayLength = 5000
	model.DefaultPriceRatioArrayLength = 10000
	model.EntryNetValue = 0.0
	model.ExitNetValue = 0.0
	model.LoserNums = 0
	model.ExpensiveStockOrderChannel = make(chan *alpaca.Order)
	model.CheapStockOrderChannel = make(chan *alpaca.Order)
	model.MinProfitThreshold = 0.0
}
