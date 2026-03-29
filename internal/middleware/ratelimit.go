package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	apierrors "github.com/nielwyn/inventory-system/pkg/errors"
	"github.com/nielwyn/inventory-system/pkg/response"
	"golang.org/x/time/rate"
)

type ipLimiter struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type RateLimiter struct {
	mu       sync.Mutex
	limiters map[string]*ipLimiter
	r        rate.Limit
	b        int
}

func NewRateLimiter(r rate.Limit, b int) *RateLimiter {
	rl := &RateLimiter{
		limiters: make(map[string]*ipLimiter),
		r:        r,
		b:        b,
	}
	go rl.cleanup()
	return rl
}

func (rl *RateLimiter) cleanup() {
	for {
		time.Sleep(time.Minute)
		rl.mu.Lock()
		for ip, l := range rl.limiters {
			if time.Since(l.lastSeen) > 3*time.Minute {
				delete(rl.limiters, ip)
			}
		}
		rl.mu.Unlock()
	}
}

func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		rl.mu.Lock()
		if _, exists := rl.limiters[ip]; !exists {
			rl.limiters[ip] = &ipLimiter{limiter: rate.NewLimiter(rl.r, rl.b)}
		}
		rl.limiters[ip].lastSeen = time.Now()
		l := rl.limiters[ip].limiter
		rl.mu.Unlock()

		if !l.Allow() {
			response.Error(c, http.StatusTooManyRequests, apierrors.CodeRateLimited, "Rate limit exceeded")
			c.Abort()
			return
		}

		c.Next()
	}
}
