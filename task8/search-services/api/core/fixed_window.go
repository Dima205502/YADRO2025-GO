package core

import (
	"sync"
	"time"
)

// не проходит тесты
type FixedWindowLimiter struct {
	requests    int
	maxRequests int
	resetTime   time.Time
	interval    time.Duration
	lock        sync.Mutex
}

func NewFixedWindowLimiter(maxRequests int, interval time.Duration) *FixedWindowLimiter {
	return &FixedWindowLimiter{
		maxRequests: maxRequests,
		interval:    interval,
		resetTime:   time.Now().Add(interval),
	}
}

func (fw *FixedWindowLimiter) Wait() {
	for {
		fw.lock.Lock()
		now := time.Now()

		if now.After(fw.resetTime) {
			fw.requests = 0
			fw.resetTime = now.Add(fw.interval)
		}

		if fw.requests < fw.maxRequests {
			fw.requests++
			fw.lock.Unlock()
			return
		}

		waitTime := fw.resetTime.Sub(now)
		fw.lock.Unlock()
		time.Sleep(waitTime)
	}
}
