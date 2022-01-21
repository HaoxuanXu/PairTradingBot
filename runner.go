package main

import (
	"github.com/HaoxuanXu/TradingBot/strats/pairtrading"
)

func JobRunner(assetType, accountType string, entryPercent float64, startTime string) {
	pairtrading.PairTradingJob(assetType, accountType, entryPercent, startTime)
}
