package main

import "github.com/HaoxuanXu/TradingBot/strats/pairtrading"

func JobRunner(assetType string) {
	pairtrading.PairTradingJob(assetType, "paper", 0.8)
}
