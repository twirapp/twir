package middlewares

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/twirapp/twir/apps/api-gql/internal/server/rate_limiter"
)

func (m *Middlewares) RateLimit(key string, capacity int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		resp, err := m.rateLimiter.Use(
			c.Request.Context(),
			&rate_limiter.LeakyOptions{
				KeyPrefix:       fmt.Sprintf("api-gql:ratelimiter:%s:%s", c.ClientIP(), key),
				MaximumCapacity: capacity,
				WindowSeconds:   int(window.Seconds()),
			},
			1,
		)
		if err != nil {
			c.AbortWithStatus(500)
			return
		}

		c.Header("X-Rate-Limit-Bucket", key)
		c.Header("X-Rate-Limit-Limit", fmt.Sprint(capacity))
		c.Header("X-Rate-Limit-Remaining", fmt.Sprint(resp.RemainingTokens))
		c.Header("X-Rate-Limit-Reset", fmt.Sprint(resp.ResetAt.Unix()))
		if !resp.Success {
			c.AbortWithStatus(429)
			return
		}

		c.Next()
	}
}
