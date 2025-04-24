package core

import (
	"context"
	"log/slog"
	"time"

	"golang.org/x/time/rate"
)

type UpdateStatus string
type LimiterAlgo string

const (
	StatusUpdateUnknown UpdateStatus = "unknown"
	StatusUpdateIdle    UpdateStatus = "idle"
	StatusUpdateRunning UpdateStatus = "running"
)

const (
	FixedWindowAlgo LimiterAlgo = "fixed window"
	SlidingLogAlgo  LimiterAlgo = "sliding log"
	TokenBucketAlgo LimiterAlgo = "token bucket"
	DefaultAlgo     LimiterAlgo = "default algo"
)

type UpdateStats struct {
	WordsTotal    int
	WordsUnique   int
	ComicsFetched int
	ComicsTotal   int
}

type Comics struct {
	ID    int
	URL   string
	Score int
}

// проходит тесты :)
type DefaultLimiter struct {
	limiter *rate.Limiter
}

func NewDefaultLimiter(limit int, interval time.Duration) *DefaultLimiter {
	return &DefaultLimiter{
		limiter: rate.NewLimiter(rate.Limit(limit), 1),
	}
}

func GetRateLimiter(algo LimiterAlgo, limit int, interval time.Duration) Waiter {
	switch algo {
	case FixedWindowAlgo:
		return NewFixedWindowLimiter(limit, interval)
	case SlidingLogAlgo:
		return NewSlidingLogLimiter(limit, interval)
	case TokenBucketAlgo:
		return NewTokenBucketLimiter(limit, interval)
	case DefaultAlgo:
		return NewDefaultLimiter(limit, interval)
	default:
		return NewDefaultLimiter(limit, interval)
	}
}

func (dl *DefaultLimiter) Wait() {
	if err := dl.limiter.Wait(context.Background()); err != nil {
		slog.Error("defaultLimiter", "Wait error", err)
	}
}
