package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/HaoxuanXu/TradingBot/internal/broker"
	"github.com/go-co-op/gocron"
)

func main() {
	assetType := flag.String("assettype", "gold", "this is the asset type we will run the job with")
	accountType := flag.String("accounttype", "paper", "This determines if we will run the job on paper or live accounts")
	startTime := flag.String("starttime", "9:45", "this is the time we will start trading each day")
	entryPercent := flag.Float64("entrypercent", 0.12, "this is the percent of portfolio value we will commit")
	flag.Parse()

	loc, _ := time.LoadLocation("America/New_York")
	s := gocron.NewScheduler(loc)
	s.Every(1).Days().At(*startTime).Do(func() {
		JobRunner(*assetType, *accountType, *entryPercent)
	})
	fmt.Printf("Planning to run the %s job for %.1f%% of the portfolio on the %s account...\n", *assetType, *entryPercent*100, *accountType)
	fmt.Printf("Wait for the job to begin at %s EST\n", *startTime)
	if broker.GetBroker(*accountType, *entryPercent).Clock.IsOpen {
		JobRunner(*assetType, *accountType, *entryPercent)
	}
	s.StartBlocking()

}