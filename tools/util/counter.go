package util

import (
	"sync"
	"time"
)

type Counter struct {
	BaseTime    time.Time
	Incrementer int
}

var lock = &sync.Mutex{}
var counter *Counter

func GetCounter() *Counter {
	lock.Lock()
	defer lock.Unlock()

	if counter == nil {
		counter = &Counter{
			BaseTime:    time.Now(),
			Incrementer: 0,
		}
	}
	return counter
}
