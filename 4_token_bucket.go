package ratelimiter

import (
	"context"
	"time"
)

type TokenBucketLimiter struct {
	tokenChan chan struct{}
}

func NewTokenBucketLimiter(ctx context.Context, limit int64, interval time.Duration) *TokenBucketLimiter {
	limiter := &TokenBucketLimiter{
		tokenChan: make(chan struct{}, limit),
	}

	for i := 0; i < int(limit); i++ {
		limiter.tokenChan <- struct{}{}
	}

	replenishmentInterval := interval.Nanoseconds() / limit
	go limiter.startBucketReplenishment(ctx, time.Duration(replenishmentInterval))

	return limiter
}

func (l *TokenBucketLimiter) startBucketReplenishment(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			select {
			case l.tokenChan <- struct{}{}:
			default:
			}
		}
	}
}

func (l *TokenBucketLimiter) Allow() bool {
	select {
	case <-l.tokenChan:
		return true
	default:
		return false
	}
}
