package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/ahmadammarm/amami/backend/pkg/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type IPRateLimiter struct {
	ips map[string]*rate.Limiter
	mu  *sync.RWMutex
	r   rate.Limit
	b   int
}

func NewIPRateLimiter(r rate.Limit, b int) *IPRateLimiter {
	i := &IPRateLimiter{
		ips: make(map[string]*rate.Limiter),
		mu:  &sync.RWMutex{},
		r:   r,
		b:   b,
	}

	return i
}

func (i *IPRateLimiter) AddIP(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	limiter := rate.NewLimiter(i.r, i.b)
	i.ips[ip] = limiter

	return limiter
}

func (i *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	i.mu.RLock()
	limiter, exists := i.ips[ip]
	i.mu.RUnlock()

	if !exists {
		return i.AddIP(ip)
	}

	return limiter
}

// RateLimiter returns a gin.HandlerFunc that applies IP-based rate limiting.
// reqs is the number of requests allowed per duration window.
func RateLimiter(reqs int, window time.Duration) gin.HandlerFunc {
	// Limit is requests per second
	r := rate.Every(window / time.Duration(reqs))
	
	// Create the rate limiter tracker
	limiter := NewIPRateLimiter(r, reqs)

	return func(c *gin.Context) {
		ip := c.ClientIP()
		l := limiter.GetLimiter(ip)
		
		if !l.Allow() {
			utils.ErrorResponse(c, http.StatusTooManyRequests, "Too many requests. Please try again later.")
			c.Abort()
			return
		}
		
		c.Next()
	}
}
