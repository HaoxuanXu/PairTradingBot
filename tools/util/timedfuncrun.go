package util

import (
	"time"
)

func TimedFuncRun(runDuration time.Duration, runFunc func(), interval time.Duration) {
	beginTime := time.Now()
	for time.Since(beginTime) <= runDuration {
		runFunc()
		time.Sleep(interval)
	}
}
