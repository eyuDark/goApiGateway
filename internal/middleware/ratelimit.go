package middleware

import (
	"math"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type clientBucket struct {
	tokens   int
	lastSeen time.Time
}

var (
	mu      sync.RWMutex // NOT sync.Mutex!
	clients = make(map[string]*clientBucket)
)

var once sync.Once

func RateLimiter(rps int) gin.HandlerFunc {
	once.Do(func() {
		go func() {
			ticker := time.NewTicker(1 * time.Minute)
			for range ticker.C {
				cleanupInactiveClients(5 * time.Minute)
			}
		}()
	})
	return func(c *gin.Context) {
		ip := c.ClientIP()
		now := time.Now()
		mu.RLock()
		bucket, exists := clients[ip]
		mu.RUnlock()
		if !exists {
			mu.Lock()
			if bucket, exists = clients[ip]; exists {
				mu.Unlock()
			} else {
				bucket = &clientBucket{
					tokens:   rps,
					lastSeen: now,
				}
				clients[ip] = bucket
				mu.Unlock()
			}
		}
		mu.Lock()
		defer mu.Unlock()
		elapsed := now.Sub(bucket.lastSeen).Seconds()
		refill := int(elapsed * float64(rps))
		newTokens := min(bucket.tokens+refill, rps)
		bucket.lastSeen = now

		c.Header("X-RateLimit-Limit", strconv.Itoa(rps))
		c.Header("X-RateLimit-Remaining", strconv.Itoa(max(0, newTokens-1)))
		if newTokens > 0 {
			bucket.tokens = newTokens - 1
			c.Next()
			return
		}
		retryAfter := int(math.Ceil(1.0 / float64(rps)))


		c.Header("Retry-After", strconv.Itoa(retryAfter))
		c.Header("X-RateLimit-Remaining", "0")
		c.AbortWithStatus(http.StatusTooManyRequests)
		
	}
}
func cleanupInactiveClients(maxAge time.Duration) {
	cutoff := time.Now().Add(-maxAge)

	mu.Lock()
	defer mu.Unlock()

	for ip, bucket := range clients {
		if bucket.lastSeen.Before(cutoff) {
			delete(clients, ip)
		}
	}
}
