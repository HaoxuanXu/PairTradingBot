package main

import (
	"github.com/HaoxuanXu/TradingBot/strats/pairtrading"
)

func JobRunner(assetType, accountType string, entryPercent float64) {
	pairtrading.PairTradingJob(assetType, accountType, entryPercent)
}
