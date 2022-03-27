package main

import (
	"github.com/HaoxuanXu/TradingBot/strats/pairtrading"
)

func JobRunner(assetType, accountType, serverType string, entryPercent float64, startTime string) {
	pairtrading.PairTradingJob(assetType, accountType, serverType, entryPercent, startTime)
}
