package updater

import "github.com/alpacahq/alpaca-trade-api-go/v2/marketdata"

func UpdatePriceRatio(expensiveQuote, cheapQuote *marketdata.Quote) (float64, float64) {
	expensiveLong := expensiveQuote.BidPrice
	expensiveShort := expensiveQuote.AskPrice

	cheapLong := cheapQuote.BidPrice
	cheapShort := cheapQuote.AskPrice

	expensiveLongCheapShort := expensiveLong / cheapShort
	expensiveShortCheapLong := expensiveShort / cheapLong

	return expensiveShortCheapLong, expensiveLongCheapShort
}