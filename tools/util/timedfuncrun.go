package util

import (
	"fmt"
	"time"
)

func TimedFuncRun(runDuration time.Duration, runFunc func(), runInterval time.Duration) {
	beginTime := time.Now()
	fmt.Println("run for 1 minute")
	for time.Since(beginTime) <= runDuration {
		runFunc()
		time.Sleep(runInterval)
	}
	fmt.Println("finished")
}
