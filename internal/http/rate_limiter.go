package http

import (
	"fmt"
	"sync"
	"time"
)

type (
	RateLimiter interface {
		Allow(ip string) bool
		Cleanup()
	}
)

type rateLimiter struct {
	maxTokens    int
	maxDuration  time.Duration
	tokenBuckets map[string]TokenBucket
}

func NewRateLimiter(maxRequest int, maxDuration time.Duration) RateLimiter {
	rateLimiter := &rateLimiter{
		maxTokens:    maxRequest,
		maxDuration:  maxDuration,
		tokenBuckets: make(map[string]TokenBucket),
	}

	return rateLimiter
}

func (r *rateLimiter) Allow(userID string) bool {
	var mutex sync.Mutex
	mutex.Lock()
	defer mutex.Unlock()

	if userID == "" {
		return false
	}

	tokenBucket, bucketExists := r.tokenBuckets[userID]
	if !bucketExists {
		tokenBucket = NewTokenBucket(r.maxTokens, r.maxDuration)
		r.tokenBuckets[userID] = tokenBucket
	}

	return tokenBucket.Allow(userID)
}

func (r *rateLimiter) TokenBuckets() map[string]TokenBucket {
	return r.tokenBuckets
}

func (r *rateLimiter) Cleanup() {
	ticker := time.NewTicker(r.maxDuration)

	for range ticker.C {
		var mutex sync.Mutex
		mutex.Lock()

		timeToChek := time.Now().Add(-r.maxDuration)
		for userID, bucket := range r.tokenBuckets {
			if bucket.LastRefill().Before(timeToChek) {
				fmt.Println("Cleaning up token bucket for user:", userID)
				delete(r.tokenBuckets, userID)
			}
		}

		mutex.Unlock()
	}
}
