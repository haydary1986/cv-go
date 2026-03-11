package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type visitor struct {
	count    int
	lastSeen time.Time
}

type RateLimiter struct {
	visitors map[string]*visitor
	mu       sync.RWMutex
	limit    int
	window   time.Duration
	done     chan struct{}
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*visitor),
		limit:    limit,
		window:   window,
		done:     make(chan struct{}),
	}
	go rl.cleanup()
	return rl
}

func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			rl.mu.Lock()
			for ip, v := range rl.visitors {
				if time.Since(v.lastSeen) > rl.window {
					delete(rl.visitors, ip)
				}
			}
			rl.mu.Unlock()
		case <-rl.done:
			return
		}
	}
}

func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		rl.mu.Lock()
		v, exists := rl.visitors[ip]
		if !exists {
			rl.visitors[ip] = &visitor{count: 1, lastSeen: time.Now()}
			rl.mu.Unlock()
			c.Next()
			return
		}

		if time.Since(v.lastSeen) > rl.window {
			v.count = 1
			v.lastSeen = time.Now()
			rl.mu.Unlock()
			c.Next()
			return
		}

		v.count++
		v.lastSeen = time.Now()

		if v.count > rl.limit {
			rl.mu.Unlock()
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded"})
			c.Abort()
			return
		}

		rl.mu.Unlock()
		c.Next()
	}
}

// LoginRateLimiter specifically for login attempts
type LoginRateLimiter struct {
	attempts map[string]*loginAttempt
	mu       sync.RWMutex
	done     chan struct{}
}

type loginAttempt struct {
	count    int
	lockedAt *time.Time
}

func NewLoginRateLimiter() *LoginRateLimiter {
	lrl := &LoginRateLimiter{
		attempts: make(map[string]*loginAttempt),
		done:     make(chan struct{}),
	}
	go lrl.cleanup()
	return lrl
}

func (lrl *LoginRateLimiter) cleanup() {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			lrl.mu.Lock()
			for email, a := range lrl.attempts {
				if a.lockedAt != nil && time.Since(*a.lockedAt) > 15*time.Minute {
					delete(lrl.attempts, email)
				} else if a.lockedAt == nil && a.count > 0 {
					// Clean up stale attempts older than 30 minutes
					delete(lrl.attempts, email)
				}
			}
			lrl.mu.Unlock()
		case <-lrl.done:
			return
		}
	}
}

func (lrl *LoginRateLimiter) Check(email string) bool {
	lrl.mu.RLock()
	defer lrl.mu.RUnlock()

	attempt, exists := lrl.attempts[email]
	if !exists {
		return true
	}

	if attempt.lockedAt != nil {
		if time.Since(*attempt.lockedAt) < 15*time.Minute {
			return false
		}
	}
	return true
}

func (lrl *LoginRateLimiter) RecordFailure(email string) {
	lrl.mu.Lock()
	defer lrl.mu.Unlock()

	attempt, exists := lrl.attempts[email]
	if !exists {
		lrl.attempts[email] = &loginAttempt{count: 1}
		return
	}

	attempt.count++
	if attempt.count >= 5 {
		now := time.Now()
		attempt.lockedAt = &now
	}
}

func (lrl *LoginRateLimiter) Reset(email string) {
	lrl.mu.Lock()
	defer lrl.mu.Unlock()
	delete(lrl.attempts, email)
}
