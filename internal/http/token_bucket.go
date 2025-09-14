package http

import (
	"sync"
	"time"
)

type (
	TokenBucket interface {
		Allow(userID string) bool
		LastRefill() time.Time
	}

	tokenBucket struct {
		maxTokens    int
		tokens       int
		lastRefilled time.Time
		refillRate   time.Duration
		mu           sync.Mutex
	}
)

func NewTokenBucket(maxTokens int, refillRate time.Duration) TokenBucket {
	return &tokenBucket{
		maxTokens:  maxTokens,
		tokens:     maxTokens,
		refillRate: refillRate,
	}
}

func (tb *tokenBucket) LastRefill() time.Time {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	return tb.lastRefilled
}

func (tb *tokenBucket) Allow(userID string) bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	if userID == "" {
		return false
	}

	now := time.Now()
	elapsed := now.Sub(tb.lastRefilled)

	// Refill tokens based on elapsed time
	tokensToAdd := int(elapsed / tb.refillRate)
	if tokensToAdd > 0 {
		tb.tokens += tokensToAdd
		if tb.tokens > tb.maxTokens {
			tb.tokens = tb.maxTokens
		}
		tb.lastRefilled = now
	}

	if tb.tokens > 0 {
		tb.tokens--
		return true
	}

	return false
}
