package middleware

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type clientLimiter struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var (
	clients   = make(map[string]*clientLimiter)
	mu        sync.Mutex
	rateLimit = rate.Every(time.Minute / 200) // 200 req/min
	burst     = 100
)

// StartRateLimiter starts a cron that cleans old clients every 5 minutes.
func StartRateLimiter() {
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		for range ticker.C {
			mu.Lock()
			for ip, cl := range clients {
				if time.Since(cl.lastSeen) > 10*time.Minute {
					delete(clients, ip)
				}
			}
			mu.Unlock()
		}
	}()
}

func getLimiter(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	if client, exists := clients[ip]; exists {
		client.lastSeen = time.Now()
		return client.limiter
	}

	limiter := rate.NewLimiter(rateLimit, burst)
	clients[ip] = &clientLimiter{
		limiter:  limiter,
		lastSeen: time.Now(),
	}
	return limiter
}

// RateLimit returns a Gin middleware that limits requests per IP using token bucket
func RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter := getLimiter(ip)

		remaining := limiter.Burst() - int(limiter.Tokens())

		// Headers para controle
		c.Header("X-RateLimit-Limit", strconv.Itoa(burst))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(remaining))

		if !limiter.Allow() {
			c.AbortWithStatus(http.StatusTooManyRequests)
			return
		}

		c.Next()
	}
}
