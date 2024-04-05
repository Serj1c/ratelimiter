package ratelimiter

import (
	"sync"
	"time"
)

type SlidingLogLimiter struct {
	limit    int32
	interval time.Duration
	logs     []time.Time
	mu       sync.Mutex
}

func NewSlidingLogLimiter(limit int32, interval time.Duration) *SlidingLogLimiter {
	limiter := &SlidingLogLimiter{
		limit:    limit,
		interval: interval,
		logs:     make([]time.Time, 0),
	}

	return limiter
}

func (l *SlidingLogLimiter) Allow() bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	lastPeriod := time.Now().Add(-l.interval)
	for len(l.logs) != 0 && l.logs[0].Before(lastPeriod) {
		l.logs = l.logs[1:]
	}

	l.logs = append(l.logs, time.Now())
	return len(l.logs) <= int(l.limit)
}
