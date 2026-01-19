package gincontext

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
)

const GinContextKey = "GqlGinContext"

func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), GinContextKey, c)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func GetGinContext(ctx context.Context) (*gin.Context, error) {
	ginContext := ctx.Value(GinContextKey)
	if ginContext == nil {
		err := fmt.Errorf("could not retrieve gin.Context")
		return nil, err
	}

	gc, ok := ginContext.(*gin.Context)
	if !ok {
		err := fmt.Errorf("gin.Context has wrong type")
		return nil, err
	}
	return gc, nil
}

// GetBaseUrl returns the base URL for the application based on request headers or fallback config
// It checks X-Forwarded-Host header first, then falls back to c.Request.Host, and finally to config
func GetBaseUrl(c *gin.Context, baseUrl string) string {
	var host string
	var scheme string

	// Check X-Forwarded-Host header first
	if forwardedHost := c.GetHeader("X-Forwarded-Host"); forwardedHost != "" {
		host = forwardedHost
	} else if c.Request.Host != "" {
		// Fall back to Request.Host
		host = c.Request.Host
	}

	// Determine scheme
	if forwardedProto := c.GetHeader("X-Forwarded-Proto"); forwardedProto != "" {
		scheme = forwardedProto
	} else if c.Request.TLS != nil {
		scheme = "https"
	} else {
		scheme = "http"
	}

	// If we have a host from headers, construct URL
	if host != "" {
		return fmt.Sprintf("%s://%s", scheme, host)
	}

	// Fall back to config
	return baseUrl
}

// GetBaseUrlFromContext is a helper that gets base URL from context by extracting gin.Context
func GetBaseUrlFromContext(ctx context.Context, baseUrl string) (string, error) {
	gc, err := GetGinContext(ctx)
	if err != nil {
		// If we can't get gin context, fall back to config
		return baseUrl, nil
	}

	return GetBaseUrl(gc, baseUrl), nil
}
