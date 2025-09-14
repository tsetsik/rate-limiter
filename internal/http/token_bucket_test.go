package http

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestToken_Allow(t *testing.T) {
	t.Parallel()

	t.Run("should return false for empty userID", func(t *testing.T) {
		t.Parallel()

		limiter := NewTokenBucket(5, time.Minute)
		got := limiter.Allow("")

		require.False(t, got)
	})

	t.Run("should return true for new userID", func(t *testing.T) {
		t.Parallel()

		limiter := NewTokenBucket(5, time.Minute)
		got := limiter.Allow("user1")

		require.True(t, got)
	})

	t.Run("should return true for 5 parallel requests with rate for one minute", func(t *testing.T) {
		t.Parallel()

		max := 5
		limiter := NewTokenBucket(5, time.Minute)

		for i := 0; i <= max; i++ {
			got := limiter.Allow("user1")
			if i < max {
				require.True(t, got)
			} else {
				require.False(t, got)
			}
		}
	})
}
