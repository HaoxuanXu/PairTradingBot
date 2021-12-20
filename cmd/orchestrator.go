package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/go-co-op/gocron"
)

func main() {
	assetType := flag.String("assettype", "gold", "this is the asset type we will run the job with")
	accountType := flag.String("accounttype", "paper", "This determines if we will run the job on paper or live accounts")
	startTime := flag.String("starttime", "9:40", "this is the time we will start trading each day")
	flag.Parse()

	loc, _ := time.LoadLocation("America/New_York")
	s := gocron.NewScheduler(loc)
	s.Every(1).Days().At(*startTime).Do(func() {
		JobRunner(*assetType, *accountType)
	})
	fmt.Printf("Planning to run the %s job on the %s account...\n", *assetType, *accountType)
	fmt.Printf("Wait for the job to begin at %s EST\n", *startTime)
	s.StartBlocking()
}
