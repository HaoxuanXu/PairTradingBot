package util

import (
	"time"
)

func TimedFuncRun(runDuration time.Duration, runFunc func()) {
	beginTime := time.Now()
	for time.Since(beginTime) <= runDuration {
		runFunc()
	}
}
