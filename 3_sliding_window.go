package ratelimiter

import (
	"sync"
	"time"
)

type SlidingWindowLimiter struct {
	limit    int32
	interval time.Duration
	mu       sync.Mutex

	currentTime   time.Time
	previousCount int32
	currentCount  int32
}

func NewSlidingWindowLimiter(limit int32, interval time.Duration) *SlidingWindowLimiter {
	limiter := &SlidingWindowLimiter{
		limit:       limit,
		interval:    interval,
		currentTime: time.Now(),
	}

	return limiter
}

func (l *SlidingWindowLimiter) Allow() bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	newPeriodTime := l.currentTime.Add(l.interval)
	if time.Now().After(newPeriodTime) {
		l.currentTime = time.Now()
		l.previousCount = l.currentCount
		l.currentCount = 0
	}

	interval := float64(l.interval)
	currentCount := float64(l.currentCount)
	previousCount := float64(l.previousCount)
	elapsed := time.Now().Sub(l.currentTime).Seconds()
	count := (previousCount * (interval - elapsed) / interval) + currentCount
	if int32(count) >= l.limit {
		return false
	}

	l.currentCount++
	return true
}
