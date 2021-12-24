package logging

import (
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"time"

	"github.com/HaoxuanXu/TradingBot/db"
	"github.com/HaoxuanXu/TradingBot/internal/broker"
	"github.com/HaoxuanXu/TradingBot/strats/pairtrading/model"
)

func SetLogging(assetType string) *os.File {
	dt := time.Now()
	logName := fmt.Sprintf("%d-%d-%d", dt.Year(), dt.Month(), dt.Day()) + "_" + "TradingLog.log"
	fullLogPath := db.MapLogPath(assetType) + logName
	logFile, err := os.OpenFile(fullLogPath, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	multiWrite := io.MultiWriter(os.Stdout, logFile)

	log.SetOutput(multiWrite)
	log.Printf("logging the trading record to %s\n", fullLogPath)
	return logFile
}

func LogTransaction(model *model.PairTradingModel, broker *broker.AlpacaBroker) {
	if !broker.HasPosition {
		if model.IsLongExpensiveStockShortCheapStock {
			model.EntryNetValue = math.Abs(model.CheapStockFilledPrice*model.CheapStockFilledQuantity) - math.Abs(model.ExpensiveStockFilledPrice*model.ExpensiveStockFilledQuantity)
			log.Printf("long %s: %f shares; short %s: %f shares   --  (repeatNum -> %d, priceRatio -> %f)\n",
				model.ExpensiveStockSymbol,
				model.ExpensiveStockEntryVolume,
				model.CheapStockSymbol,
				model.CheapStockEntryVolume,
				model.RepeatNumThreshold,
				model.PriceRatioThreshold,
			)
		} else {
			model.EntryNetValue = math.Abs(model.ExpensiveStockFilledPrice*model.ExpensiveStockFilledQuantity) - math.Abs(model.CheapStockFilledPrice*model.CheapStockFilledQuantity)
			log.Printf("short %s: %f shares; long %s: %f shares   --   (repeatNum -> %d, priceRatio -> %f)\n",
				model.ExpensiveStockSymbol,
				model.ExpensiveStockEntryVolume,
				model.CheapStockSymbol,
				model.CheapStockEntryVolume,
				model.RepeatNumThreshold,
				model.PriceRatioThreshold,
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
		log.Printf("position closed. Presumed Profit: $%f. Actual Profit: $%f   --   (repeatNum -> %d, priceRatio -> %f)\n",
			presumedProfit,
			actualProfit,
			model.RepeatNumThreshold,
			model.PriceRatioThreshold,
		)
		broker.HasPosition = false
		broker.TransactionNums++
	}
}
