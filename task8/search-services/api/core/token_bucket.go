package core

import (
	"sync"
	"time"
)

// не проходит тесты
type TokenBucketLimiter struct {
	capacity      int
	tokens        int
	tokenInterval time.Duration
	lastRefill    time.Time
	mu            sync.Mutex
}

func NewTokenBucketLimiter(limit int, interval time.Duration) *TokenBucketLimiter {
	return &TokenBucketLimiter{
		capacity:      limit,
		tokens:        limit,
		tokenInterval: interval / time.Duration(limit),
		lastRefill:    time.Now(),
	}
}

func (tb *TokenBucketLimiter) Wait() {
	for {
		tb.mu.Lock()
		now := time.Now()

		elapsed := now.Sub(tb.lastRefill)
		newTokens := int(elapsed / tb.tokenInterval)
		if newTokens > 0 {
			tb.tokens += newTokens
			if tb.tokens > tb.capacity {
				tb.tokens = tb.capacity
			}

			tb.lastRefill = tb.lastRefill.Add(time.Duration(newTokens) * tb.tokenInterval)
		}

		if tb.tokens > 0 {
			tb.tokens--
			tb.mu.Unlock()
			return
		}

		remaining := tb.tokenInterval - elapsed%tb.tokenInterval
		tb.mu.Unlock()
		time.Sleep(remaining)
	}
}
