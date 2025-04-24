package core

import (
	"sync"
	"time"
)

// не проходит тесты
type SlidingLogLimiter struct {
	timestamps  []time.Time
	maxRequests int
	interval    time.Duration
	lock        sync.Mutex
}

func NewSlidingLogLimiter(maxRequests int, interval time.Duration) *SlidingLogLimiter {
	return &SlidingLogLimiter{
		maxRequests: maxRequests,
		interval:    interval,
		timestamps:  make([]time.Time, 0),
	}
}

func (swl *SlidingLogLimiter) Wait() {
	for {
		swl.lock.Lock()
		now := time.Now()

		validTimestamps := make([]time.Time, 0, len(swl.timestamps))
		for _, ts := range swl.timestamps {
			if now.Sub(ts) <= swl.interval {
				validTimestamps = append(validTimestamps, ts)
			}
		}
		swl.timestamps = validTimestamps

		if len(swl.timestamps) < swl.maxRequests {
			swl.timestamps = append(swl.timestamps, now)
			swl.lock.Unlock()
			return
		}

		sleepDuration := swl.interval - now.Sub(swl.timestamps[0])
		swl.lock.Unlock()

		if sleepDuration < 0 {
			sleepDuration = 0
		}
		time.Sleep(sleepDuration)
	}
}
