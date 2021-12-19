package main

import "github.com/HaoxuanXu/TradingBot/strats/pairtrading"

func JobRunner() {
	pairtrading.PairTradingJob("GLD", "IAU", "gold", "paper", 0.8)
}