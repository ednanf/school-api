package ratelimit

import (
	"testing"

	"golang.org/x/time/rate"
)

// TestGetLimiterReturnsSameLimiterForSameIP tests if the implementation creates a new limiter for every request,
// allowing clients to bypass the rate limit simply by making more requests.
func TestGetLimiterReturnsSameLimiterForSameIP(t *testing.T) {
	// Creates a limiter that replenishes one token per second, with a maximum burst of two tokens
	limiter := NewIPRateLimiter(rate.Limit(1), 2)

	// Call twice with the same IP
	first := limiter.GetLimiter("192.168.1.10")
	second := limiter.GetLimiter("192.168.1.10")

	// Check if both calls return the exact same limiter instance
	if first != second {
		t.Fatal("expected the same IP to reuse its rate limiter")
	}
}
