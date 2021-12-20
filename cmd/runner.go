package main

import "github.com/HaoxuanXu/TradingBot/strats/pairtrading"

func JobRunner(assetType, accountType string) {
	pairtrading.PairTradingJob(assetType, accountType, 0.8)
}
