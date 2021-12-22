package pairtrading

import (
	"log"
	"time"

	"github.com/HaoxuanXu/TradingBot/db"
	"github.com/HaoxuanXu/TradingBot/internal/broker"
	"github.com/HaoxuanXu/TradingBot/internal/dataengine"
	"github.com/HaoxuanXu/TradingBot/strats/pairtrading/logging"
	"github.com/HaoxuanXu/TradingBot/strats/pairtrading/model"
	"github.com/HaoxuanXu/TradingBot/strats/pairtrading/pipeline"
	"github.com/HaoxuanXu/TradingBot/strats/pairtrading/quotesprocessor"
	"github.com/HaoxuanXu/TradingBot/strats/pairtrading/signalcatcher"
	"github.com/HaoxuanXu/TradingBot/tools/util"
)

func PairTradingJob(assetType, accountType string, entryPercent float64) {
	// This job will not run if we are on weekends, so we will simply return if it is the weekends
	today := time.Now().Weekday().String()
	if today == "Saturday" || today == "Sunday" {
		log.Printf("Today is %s. We will not work today...\n", today)
		return
	}
	// initialize the data model struct and the broker struct
	tradingBroker := broker.GetBroker(accountType, entryPercent)
	dataEngine := dataengine.GetDataEngine(accountType)
	shortLongPath, longShortPath, repeatNumPath := db.MapRecordPath("gold")
	dataModel := model.GetModel(assetType, shortLongPath, longShortPath, repeatNumPath)

	// set up log file for today
	logging.SetLogging(assetType)

	// We will check if the market is open currently
	// If the market is not open, we will wait till it is open
	if tradingBroker.Clock.IsOpen {
		log.Println("Market is currently open")
	} else {
		timeToOpen := time.Until(tradingBroker.Clock.NextOpen)
		log.Printf("Waiting for %d hours until the market opens\n", int(timeToOpen.Hours()))
		time.Sleep(timeToOpen)
	}
	log.Println("Start Trading ...")
	// Start the main trading loop
	for time.Until(tradingBroker.Clock.NextClose) > 20*time.Minute {
		quotesprocessor.GetAndProcessPairQuotes(dataModel, dataEngine)
		if signalcatcher.GetEntrySignal(true, dataModel, tradingBroker) {
			pipeline.EntryShortExpensiveLongCheap(
				dataModel,
				tradingBroker,
			)
			// halt trading for a minute so the account is still treated as retail account
			util.TimedFuncRun(
				time.Minute,
				func() {
					quotesprocessor.GetAndProcessPairQuotes(dataModel, dataEngine)
				},
				10*time.Millisecond,
			)
		} else if signalcatcher.GetEntrySignal(false, dataModel, tradingBroker) {
			pipeline.EntryLongExpensiveShortCheap(
				dataModel,
				tradingBroker,
			)
			util.TimedFuncRun(
				time.Minute,
				func() {
					quotesprocessor.GetAndProcessPairQuotes(dataModel, dataEngine)
				},
				10*time.Millisecond,
			)
		} else if dataModel.IsShortExpensiveStockLongCheapStock && signalcatcher.GetExitSignal(dataModel) {
			pipeline.ExitShortExpensiveLongCheap(
				dataModel,
				tradingBroker,
			)
			util.TimedFuncRun(
				time.Minute,
				func() {
					quotesprocessor.GetAndProcessPairQuotes(dataModel, dataEngine)
				},
				10*time.Millisecond,
			)
		} else if dataModel.IsLongExpensiveStockShortCheapStock && signalcatcher.GetExitSignal(dataModel) {
			pipeline.ExitLongExpensiveShortCheap(
				dataModel,
				tradingBroker,
			)
			util.TimedFuncRun(
				time.Minute,
				func() {
					quotesprocessor.GetAndProcessPairQuotes(dataModel, dataEngine)
				},
				10*time.Millisecond,
			)
		} else {
			time.Sleep(10 * time.Millisecond)
		}
	}
	log.Println("Preparing to close the trading session ...")
	for time.Until(tradingBroker.Clock.NextClose) > time.Minute {
		quotesprocessor.GetAndProcessPairQuotes(dataModel, dataEngine)
		if !tradingBroker.HasPosition {
			break
		} else if dataModel.IsShortExpensiveStockLongCheapStock && signalcatcher.GetExitSignal(dataModel) {
			pipeline.ExitShortExpensiveLongCheap(
				dataModel,
				tradingBroker,
			)
			break
		} else if dataModel.IsLongExpensiveStockShortCheapStock && signalcatcher.GetExitSignal(dataModel) {
			pipeline.ExitLongExpensiveShortCheap(
				dataModel,
				tradingBroker,
			)
			break
		} else {
			time.Sleep(10 * time.Millisecond)
		}
	}

	// Close all positions and record data
	tradingBroker.CloseAllPositions()
	log.Printf("The amount you made today: $%.2f\n", tradingBroker.GetDailyProfit())
	log.Printf("The number of round trips you made today: %d\n", tradingBroker.TransactionNums)
	log.Printf("The number of losing trips you made todau: %d\n", dataModel.LoserNums)
	log.Println("Writing out record to json ...")
	pipeline.WriteRecord(dataModel, assetType)
	log.Println("Data successfully written to json!")
}
