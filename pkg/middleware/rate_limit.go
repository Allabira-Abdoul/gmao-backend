package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type rateLimiter struct {
	requests map[string][]time.Time
	mu       sync.Mutex
	limit    int
	window   time.Duration
}

func newRateLimiter(limit int, window time.Duration) *rateLimiter {
	rl := &rateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}

	// Simple cleanup goroutine to prevent memory leaks over time
	go func() {
		for {
			time.Sleep(window)
			rl.cleanup()
		}
	}()

	return rl
}

func (rl *rateLimiter) cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	for ip, times := range rl.requests {
		validTimes := make([]time.Time, 0)
		for _, t := range times {
			if now.Sub(t) < rl.window {
				validTimes = append(validTimes, t)
			}
		}

		if len(validTimes) == 0 {
			delete(rl.requests, ip)
		} else {
			rl.requests[ip] = validTimes
		}
	}
}

func (rl *rateLimiter) allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()

	// Initialize if not exists
	if _, exists := rl.requests[ip]; !exists {
		rl.requests[ip] = []time.Time{now}
		return true
	}

	// Filter out old requests
	validTimes := make([]time.Time, 0)
	for _, t := range rl.requests[ip] {
		if now.Sub(t) < rl.window {
			validTimes = append(validTimes, t)
		}
	}

	// Check if limit exceeded
	if len(validTimes) >= rl.limit {
		rl.requests[ip] = validTimes // keep valid times but deny
		return false
	}

	// Allow request and record time
	validTimes = append(validTimes, now)
	rl.requests[ip] = validTimes
	return true
}

// RateLimit returns a Gin middleware that limits the number of requests per IP address
// within a given time window.
func RateLimit(limit int, window time.Duration) gin.HandlerFunc {
	limiter := newRateLimiter(limit, window)

	return func(c *gin.Context) {
		ip := c.ClientIP()

		if !limiter.allow(ip) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "TOO_MANY_REQUESTS",
					"message": "Rate limit exceeded. Please try again later.",
				},
			})
			return
		}

		c.Next()
	}
}
