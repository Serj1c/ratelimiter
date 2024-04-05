package ratelimiter

import (
	"context"
	"time"
)

type LeakyBucketLimiter struct {
	tokenChan chan struct{}
}

func NewLeakyBucketLimiter(ctx context.Context, limit int64, interval time.Duration) *LeakyBucketLimiter {
	limiter := &LeakyBucketLimiter{
		tokenChan: make(chan struct{}, limit),
	}

	leakingInterval := interval.Nanoseconds() / limit
	go limiter.startBucketLeaking(ctx, time.Duration(leakingInterval))

	return limiter
}

func (l *LeakyBucketLimiter) startBucketLeaking(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			select {
			case <-l.tokenChan:
			default:
			}
		}
	}
}

func (l *LeakyBucketLimiter) Allow() bool {
	select {
	case l.tokenChan <- struct{}{}:
		return true
	default:
		return false
	}
}
