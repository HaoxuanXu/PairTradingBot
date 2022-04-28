package util

import (
	"time"
)

func TimedFuncRun(runDuration time.Duration, runFunc func(), interval int) {
	beginTime := time.Now()
	for time.Since(beginTime) <= runDuration {
		runFunc()
		if interval > 0 {
			time.Sleep(time.Duration(interval) * time.Millisecond)
		}

	}
}
