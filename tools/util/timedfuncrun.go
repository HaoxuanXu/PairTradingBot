package util

import "time"

func TimedFuncRun(runDuration time.Duration, runFunc func(), runInterval time.Duration) {
	beginTime := time.Now()

	for time.Since(beginTime) <= runDuration {
		runFunc()
		time.Sleep(runInterval)
	}
}