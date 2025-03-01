package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"k8s.io/klog/v2"
)

// Logger is a middleware that logs HTTP requests
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Get status code
		statusCode := c.Writer.Status()
		method := c.Request.Method
		clientIP := c.ClientIP()

		if raw != "" {
			path = path + "?" + raw
		}

		// Log the request
		klog.Infof("[GIN] %s | %3d | %13v | %15s | %-7s %s",
			statusCode,
			latency,
			clientIP,
			method,
			path,
		)
	}
}

// Recovery is a middleware that recovers from panics
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Log the error
				klog.Errorf("Panic recovered: %v", err)

				// Return a 500 error
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": "Internal server error",
				})
			}
		}()

		c.Next()
	}
}

// RateLimit is a middleware that limits the number of requests per second
func RateLimit(rps int) gin.HandlerFunc {
	// Simple token bucket implementation
	// In a production environment, you would use a more sophisticated rate limiter
	// such as redis-based rate limiting or a dedicated rate limiting service
	type client struct {
		tokens         int
		lastRefillTime time.Time
	}

	clients := make(map[string]*client)

	return func(c *gin.Context) {
		// Get client IP
		clientIP := c.ClientIP()

		// Get or create client
		cl, exists := clients[clientIP]
		if !exists {
			cl = &client{
				tokens:         rps,
				lastRefillTime: time.Now(),
			}
			clients[clientIP] = cl
		}

		// Refill tokens if needed
		now := time.Now()
		elapsed := now.Sub(cl.lastRefillTime).Seconds()
		tokensToAdd := int(elapsed * float64(rps))
		if tokensToAdd > 0 {
			cl.tokens = min(cl.tokens+tokensToAdd, rps)
			cl.lastRefillTime = now
		}

		// Check if client has tokens
		if cl.tokens <= 0 {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": fmt.Sprintf("Rate limit exceeded. Maximum %d requests per second.", rps),
			})
			return
		}

		// Consume a token
		cl.tokens--

		c.Next()
	}
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// CORS is a middleware that adds CORS headers
func CORS(allowOrigins []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Check if the origin is allowed
		allowed := false
		for _, allowedOrigin := range allowOrigins {
			if allowedOrigin == "*" || allowedOrigin == origin {
				allowed = true
				break
			}
		}

		if allowed {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		// Handle preflight requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
