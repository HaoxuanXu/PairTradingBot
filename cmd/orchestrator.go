package main

import (
	"fmt"
	"time"

	"github.com/go-co-op/gocron"
)

func main() {
	timeString := "9:40"
	loc, _ := time.LoadLocation("America/New_York")
	s := gocron.NewScheduler(loc)
	s.Every(1).Days().At(timeString).Do(JobRunner)
	fmt.Printf("Wait for the job to begin at %s AM EST\n", timeString)
	s.StartBlocking()
}