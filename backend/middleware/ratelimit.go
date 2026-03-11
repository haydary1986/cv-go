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
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*visitor),
		limit:    limit,
		window:   window,
	}
	go rl.cleanup()
	return rl
}

func (rl *RateLimiter) cleanup() {
	for {
		time.Sleep(time.Minute)
		rl.mu.Lock()
		for ip, v := range rl.visitors {
			if time.Since(v.lastSeen) > rl.window {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
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
}

type loginAttempt struct {
	count    int
	lockedAt *time.Time
}

func NewLoginRateLimiter() *LoginRateLimiter {
	return &LoginRateLimiter{
		attempts: make(map[string]*loginAttempt),
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
