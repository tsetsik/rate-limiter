package http

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestLimiter_Allow(t *testing.T) {
	t.Parallel()

	t.Run("should return false for empty userID", func(t *testing.T) {
		t.Parallel()

		limiter := NewRateLimiter(5, time.Minute)
		got := limiter.Allow("")
		require.False(t, got)
	})

	t.Run("should return true for new userID", func(t *testing.T) {
		t.Parallel()

		limiter := NewRateLimiter(5, time.Minute)
		got := limiter.Allow("user1")

		require.True(t, got)
	})

	t.Run("should return true for max requests in a minute", func(t *testing.T) {
		t.Parallel()

		limiter := NewRateLimiter(5, time.Minute)
		for range 5 {
			got := limiter.Allow("user" + fmt.Sprint(time.Now().UnixNano()))
			require.True(t, got)
		}
	})

	t.Run("should return false for more requests than allowed", func(t *testing.T) {
		t.Parallel()

		limiter := NewRateLimiter(5, time.Minute)
		for i := 0; i <= 5; i++ {
			got := limiter.Allow("user1")
			if i < 5 {
				require.True(t, got)
			} else {
				require.False(t, got)
			}
		}
	})

	t.Run("there shouldn't be tokens after a minute", func(t *testing.T) {
		t.Parallel()

		limiter := NewRateLimiter(5, time.Second)

		got := limiter.Allow("user1")
		require.True(t, got)

		time.Sleep(1 * time.Second)

		got = limiter.Allow("user1")
		require.True(t, got)
	})
}
