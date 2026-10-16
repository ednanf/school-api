package ratelimit

import (
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// Stores a rate limiter and the time the IP was last active.
type clientLimiter struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type IPRateLimiter struct {
	ips map[string]*clientLimiter // map that keeps track of a separate token bucket (a rate.Limiter) for each individual IP address
	mu  sync.Mutex                // mutex ensures that only one goroutine can access or modify the ips map at a time
	r   rate.Limit                // limit - steady-state rate at which tokens are added to the bucket
	b   int                       // burst - maximum capacity of the bucket
}

func NewIPRateLimiter(r rate.Limit, b int) *IPRateLimiter {
	limiter := &IPRateLimiter{
		ips: make(map[string]*clientLimiter),
		r:   r,
		b:   b,
	}

	// Periodically remove IPs that have been inactive for 30 minutes.
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()

		for range ticker.C {
			limiter.mu.Lock()

			cutoff := time.Now().Add(-30 * time.Minute)

			for ip, client := range limiter.ips {
				if client.lastSeen.Before(cutoff) {
					delete(limiter.ips, ip)
				}
			}

			limiter.mu.Unlock()
		}
	}()

	return limiter
}

func (i *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	client, exists := i.ips[ip]

	if !exists {
		client = &clientLimiter{
			limiter: rate.NewLimiter(i.r, i.b),
		}
		i.ips[ip] = client
	}

	// Each request refreshes the IP's activity timestamp
	client.lastSeen = time.Now()

	return client.limiter
}
