// Package middleware — rate limiter for sensitive endpoints (e.g. /auth/*).
//
// Implementation: token bucket per client IP, in-memory, sweepable. Suitable
// for single-instance deployments; multi-instance setups should swap to a
// Redis-backed counter (out of scope for v1.0.0).
package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type bucket struct {
	tokens   int
	lastSeen time.Time
}

// RateLimiter returns a middleware allowing at most `rate` requests per
// `window` from the same client IP. Excess requests get HTTP 429.
func RateLimiter(rate int, window time.Duration) gin.HandlerFunc {
	if rate <= 0 {
		rate = 60
	}
	if window <= 0 {
		window = time.Minute
	}

	var mu sync.Mutex
	buckets := make(map[string]*bucket)

	// Periodic sweep of stale buckets to bound memory.
	go func() {
		ticker := time.NewTicker(window * 5)
		defer ticker.Stop()
		for range ticker.C {
			now := time.Now()
			mu.Lock()
			for k, b := range buckets {
				if now.Sub(b.lastSeen) > window*5 {
					delete(buckets, k)
				}
			}
			mu.Unlock()
		}
	}()

	return func(c *gin.Context) {
		ip := clientIP(c)
		now := time.Now()

		mu.Lock()
		b, ok := buckets[ip]
		if !ok {
			b = &bucket{tokens: rate, lastSeen: now}
			buckets[ip] = b
		} else {
			elapsed := now.Sub(b.lastSeen)
			if elapsed >= window {
				b.tokens = rate
			} else {
				refill := int(elapsed.Nanoseconds() * int64(rate) / window.Nanoseconds())
				if refill > 0 {
					b.tokens += refill
					if b.tokens > rate {
						b.tokens = rate
					}
				}
			}
			b.lastSeen = now
		}

		if b.tokens <= 0 {
			mu.Unlock()
			c.Header("Retry-After", "60")
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
			return
		}
		b.tokens--
		mu.Unlock()

		c.Next()
	}
}

func clientIP(c *gin.Context) string {
	// Honour X-Forwarded-For when behind a trusted proxy. For local/dev we
	// fall back to RemoteAddr.
	if xff := c.GetHeader("X-Forwarded-For"); xff != "" {
		// keep the left-most non-empty value
		for i := 0; i < len(xff); i++ {
			if xff[i] == ',' {
				return xff[:i]
			}
		}
		return xff
	}
	host, _, err := net.SplitHostPort(c.Request.RemoteAddr)
	if err != nil {
		return c.Request.RemoteAddr
	}
	return host
}
