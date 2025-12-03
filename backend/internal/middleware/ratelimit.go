package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimiter implements token bucket rate limiting
type RateLimiter struct {
	mu           sync.RWMutex
	clients      map[string]*clientLimiter
	rate         int           // requests per window
	window       time.Duration // time window
	cleanupInterval time.Duration
}

type clientLimiter struct {
	tokens     int
	lastUpdate time.Time
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(rate int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		clients:         make(map[string]*clientLimiter),
		rate:            rate,
		window:          window,
		cleanupInterval: window * 2,
	}

	// Start cleanup goroutine
	go rl.cleanup()

	return rl
}

// Allow checks if a request is allowed
func (rl *RateLimiter) Allow(clientID string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	limiter, exists := rl.clients[clientID]

	if !exists {
		rl.clients[clientID] = &clientLimiter{
			tokens:     rl.rate - 1,
			lastUpdate: now,
		}
		return true
	}

	// Add tokens based on elapsed time
	elapsed := now.Sub(limiter.lastUpdate)
	tokensToAdd := int(elapsed / rl.window * time.Duration(rl.rate))
	
	if tokensToAdd > 0 {
		limiter.tokens = min(limiter.tokens+tokensToAdd, rl.rate)
		limiter.lastUpdate = now
	}

	if limiter.tokens > 0 {
		limiter.tokens--
		return true
	}

	return false
}

// cleanup removes old client limiters
func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(rl.cleanupInterval)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		for clientID, limiter := range rl.clients {
			if now.Sub(limiter.lastUpdate) > rl.cleanupInterval {
				delete(rl.clients, clientID)
			}
		}
		rl.mu.Unlock()
	}
}

// RateLimit middleware applies rate limiting
func RateLimit(rate int, window time.Duration) gin.HandlerFunc {
	limiter := NewRateLimiter(rate, window)

	return func(c *gin.Context) {
		clientID := c.ClientIP()
		if c.GetString("user_id") != "" {
			clientID = c.GetString("user_id")
		}

		if !limiter.Allow(clientID) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "rate limit exceeded",
				"retry_after": window.Seconds(),
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
