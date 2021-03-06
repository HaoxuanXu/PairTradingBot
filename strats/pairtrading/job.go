package pairtrading

import (
	"log"
	"time"

	"github.com/HaoxuanXu/TradingBot/db"
	"github.com/HaoxuanXu/TradingBot/internal/broker"
	"github.com/HaoxuanXu/TradingBot/internal/dataengine"
	"github.com/HaoxuanXu/TradingBot/strats/pairtrading/model"
	"github.com/HaoxuanXu/TradingBot/strats/pairtrading/pipeline"
	"github.com/HaoxuanXu/TradingBot/strats/pairtrading/quotesprocessor"
	"github.com/HaoxuanXu/TradingBot/strats/pairtrading/signalcatcher"
	"github.com/HaoxuanXu/TradingBot/strats/pairtrading/transaction"
	"github.com/HaoxuanXu/TradingBot/tools/logging"
	"github.com/HaoxuanXu/TradingBot/tools/util"
)

func PairTradingJob(assetType, accountType, serverType string, entryPercent float64, startTime string) {
	// This job will not run if we are on weekends, so we will simply return if it is the weekends
	today := time.Now().Weekday().String()
	if today == "Saturday" || today == "Sunday" {
		log.Printf("Today is %s. We will not work today...\n", today)
		return
	}
	// initialize the data model struct and the broker struct
	tradingBroker := broker.GetBroker(accountType, serverType, entryPercent)
	dataEngine := dataengine.GetDataEngine(accountType, serverType)
	tradingAssetParamConfig := db.MapRecordPath(assetType)
	dataModel := model.GetModel(tradingAssetParamConfig)

	// set up log file for today
	logFile := logging.SetLogging(assetType)

	// We will check if the market is open currently
	// If the market is not open, we will wait till it is open
	if tradingBroker.Clock.IsOpen {
		log.Println("Market is currently open")
	} else {
		timeToOpen := time.Until(tradingBroker.Clock.NextOpen)
		log.Printf("Waiting for %d hours until the market opens\n", int(timeToOpen.Hours()))
		time.Sleep(timeToOpen)
	}
	// Warm up data for a specified period of time before trading
	quotesprocessor.WarmUpData(startTime, assetType, dataModel, tradingBroker, dataEngine, tradingAssetParamConfig)
	log.Printf("Start Trading   --  (longExpensiveShortCheapRepeatNum -> %d, shortExpensiveLongCheapRepeatNum -> %d, priceRatio -> %f)\n",
		dataModel.LongExpensiveShortCheapRepeatNumThreshold,
		dataModel.ShortExpensiveLongCheapRepeatNumThreshold,
		dataModel.PriceRatioThreshold,
	)
	tradingBroker.UpdateLastTradeTime()

	counter := util.GetCounter()
	// Check if we currently have trades pending
	transaction.CheckExistingPositions(dataModel, tradingBroker)
	// Start the main trading loop
	for time.Until(tradingBroker.Clock.NextClose) > 10*time.Minute {
		pipeline.UpdateSignalThresholds(dataModel, tradingBroker, counter, false, tradingAssetParamConfig)
		quotesprocessor.GetAndProcessPairQuotes(dataModel, dataEngine)
		if signalcatcher.GetEntrySignal(true, dataModel, tradingBroker) {
			pipeline.EntryShortExpensiveLongCheap(
				dataModel,
				tradingBroker,
				tradingAssetParamConfig,
			)
			// halt trading for a minute so the account is still treated as retail account
			util.TimedFuncRun(
				45*time.Second,
				func() {
					pipeline.UpdateSignalThresholds(dataModel, tradingBroker, counter, false, tradingAssetParamConfig)
					quotesprocessor.GetAndProcessPairQuotes(dataModel, dataEngine)
				},
				10,
			)
		} else if signalcatcher.GetEntrySignal(false, dataModel, tradingBroker) {
			pipeline.EntryLongExpensiveShortCheap(
				dataModel,
				tradingBroker,
				tradingAssetParamConfig,
			)
			util.TimedFuncRun(
				45*time.Second,
				func() {
					pipeline.UpdateSignalThresholds(dataModel, tradingBroker, counter, false, tradingAssetParamConfig)
					quotesprocessor.GetAndProcessPairQuotes(dataModel, dataEngine)
				},
				10,
			)
		} else if dataModel.IsShortExpensiveStockLongCheapStock && signalcatcher.GetExitSignal(dataModel) {
			pipeline.ExitShortExpensiveLongCheap(
				dataModel,
				tradingBroker,
				tradingAssetParamConfig,
			)
			util.TimedFuncRun(
				45*time.Second,
				func() {
					pipeline.UpdateSignalThresholds(dataModel, tradingBroker, counter, false, tradingAssetParamConfig)
					quotesprocessor.GetAndProcessPairQuotes(dataModel, dataEngine)
				},
				10,
			)
		} else if dataModel.IsLongExpensiveStockShortCheapStock && signalcatcher.GetExitSignal(dataModel) {
			pipeline.ExitLongExpensiveShortCheap(
				dataModel,
				tradingBroker,
				tradingAssetParamConfig,
			)
			util.TimedFuncRun(
				45*time.Second,
				func() {
					pipeline.UpdateSignalThresholds(dataModel, tradingBroker, counter, false, tradingAssetParamConfig)
					quotesprocessor.GetAndProcessPairQuotes(dataModel, dataEngine)
				},
				10,
			)
		} else if dataModel.IsTrimmable {
			pipeline.TrimPosition(
				dataModel,
				tradingBroker,
				tradingAssetParamConfig,
			)
			util.TimedFuncRun(
				45*time.Second,
				func() {
					pipeline.UpdateSignalThresholds(dataModel, tradingBroker, counter, false, tradingAssetParamConfig)
					quotesprocessor.GetAndProcessPairQuotes(dataModel, dataEngine)
				},
				10,
			)
		} else {
			time.Sleep(10 * time.Millisecond)
		}
	}
	log.Println("Preparing to close the trading session ...")
	for time.Until(tradingBroker.Clock.NextClose) > time.Minute {
		pipeline.UpdateSignalThresholds(dataModel, tradingBroker, counter, true, tradingAssetParamConfig)
		quotesprocessor.GetAndProcessPairQuotes(dataModel, dataEngine)
		if !tradingBroker.HasPosition {
			break
		} else if dataModel.IsShortExpensiveStockLongCheapStock && signalcatcher.GetExitSignal(dataModel) {
			pipeline.ExitShortExpensiveLongCheap(
				dataModel,
				tradingBroker,
				tradingAssetParamConfig,
			)
			break
		} else if dataModel.IsLongExpensiveStockShortCheapStock && signalcatcher.GetExitSignal(dataModel) {
			pipeline.ExitLongExpensiveShortCheap(
				dataModel,
				tradingBroker,
				tradingAssetParamConfig,
			)
			break
		} else if dataModel.IsTrimmable {
			pipeline.TrimPosition(
				dataModel,
				tradingBroker,
				tradingAssetParamConfig,
			)
			util.TimedFuncRun(
				45*time.Second,
				func() {
					pipeline.UpdateSignalThresholds(dataModel, tradingBroker, counter, false, tradingAssetParamConfig)
					quotesprocessor.GetAndProcessPairQuotes(dataModel, dataEngine)
				},
				10,
			)
		} else {
			time.Sleep(10 * time.Millisecond)
		}
	}

	// Close all positions and record data
	tradingBroker.CloseAllPositions()
	log.Printf("The amount you made today: $%.2f\n", tradingBroker.GetDailyProfit())
	log.Printf("The number of round trips you made today: %d\n", tradingBroker.TransactionNums)
	log.Printf("The number of losing trips you made today: %d\n", dataModel.LoserNums)
	log.Println("Writing out record to json ...")
	pipeline.WriteRecord(dataModel, tradingAssetParamConfig)
	log.Println("Data successfully written to json!")
	logFile.Close()
}
