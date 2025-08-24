package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type RateLimiter struct {
	clients map[string]*clientInfo
	mu      sync.RWMutex
	cleanup *time.Ticker
}

type clientInfo struct {
	requests    int
	lastRequest time.Time
}

func NewRateLimiter() *RateLimiter {
	rl := &RateLimiter{
		clients: make(map[string]*clientInfo),
		cleanup: time.NewTicker(time.Hour),
	}
	
	// Cleanup old entries
	go func() {
		for range rl.cleanup.C {
			rl.mu.Lock()
			cutoff := time.Now().Add(-time.Hour)
			for ip, info := range rl.clients {
				if info.lastRequest.Before(cutoff) {
					delete(rl.clients, ip)
				}
			}
			rl.mu.Unlock()
		}
	}()
	
	return rl
}

func (rl *RateLimiter) FileUploadMiddleware(maxRequests int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		
		rl.mu.Lock()
		defer rl.mu.Unlock()
		
		now := time.Now()
		info, exists := rl.clients[clientIP]
		
		if !exists {
			rl.clients[clientIP] = &clientInfo{
				requests:    1,
				lastRequest: now,
			}
			c.Next()
			return
		}
		
		// Reset counter if window has passed
		if now.Sub(info.lastRequest) > window {
			info.requests = 1
			info.lastRequest = now
			c.Next()
			return
		}
		
		// Check if limit exceeded
		if info.requests >= maxRequests {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many upload requests. Please try again later.",
			})
			c.Abort()
			return
		}
		
		info.requests++
		info.lastRequest = now
		c.Next()
	}
}

func (rl *RateLimiter) Stop() {
	rl.cleanup.Stop()
}