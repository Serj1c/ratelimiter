package ratelimiter

import (
	"context"
	"sync/atomic"
	"time"
)

type FixedWindowLimiter struct {
	current int32
	limit   int32
}

func NewFixedWindowLimiter(ctx context.Context, limit int32, interval time.Duration) *FixedWindowLimiter {
	limiter := &FixedWindowLimiter{
		current: 0,
		limit:   limit,
	}

	go limiter.startLimiterRefresh(ctx, interval)

	return limiter
}

func (l *FixedWindowLimiter) startLimiterRefresh(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	select {
	case <-ctx.Done():
		return
	case <-ticker.C:
		atomic.StoreInt32(&l.current, 0)
	}
}

func (l *FixedWindowLimiter) Allow() bool {
	curr := atomic.LoadInt32(&l.current)
	if curr > l.limit {
		return false
	}

	for !atomic.CompareAndSwapInt32(&l.current, curr, curr+1) {
		curr = atomic.LoadInt32(&l.current)
	}

	return curr < l.limit
}
