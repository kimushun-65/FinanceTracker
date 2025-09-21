// Package middleware provides HTTP middleware implementations for the FinanceTracker API.
package middleware

import (
	"net/http"
	"strings"

	"financetracker/pkg/config"

	"github.com/gin-gonic/gin"
)

// CORS returns a middleware that handles Cross-Origin Resource Sharing.
func CORS(cfg *config.Config) gin.HandlerFunc {
	// Parse allowed origins from config
	allowedOrigins := make(map[string]bool)
	if cfg.AllowedOrigins != "" {
		for _, origin := range strings.Split(cfg.AllowedOrigins, ",") {
			allowedOrigins[strings.TrimSpace(origin)] = true
		}
	}

	// Default to allowing localhost in development
	if cfg.Environment == "development" && len(allowedOrigins) == 0 {
		allowedOrigins["http://localhost:3000"] = true
		allowedOrigins["http://localhost:3001"] = true
	}

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Check if origin is allowed
		if origin != "" {
			// Check exact match
			if allowedOrigins[origin] {
				c.Header("Access-Control-Allow-Origin", origin)
			} else if allowedOrigins["*"] {
				// Allow all origins if * is specified
				c.Header("Access-Control-Allow-Origin", "*")
			} else {
				// Check for wildcard subdomain matching (e.g., *.example.com)
				for allowedOrigin := range allowedOrigins {
					if strings.HasPrefix(allowedOrigin, "*.") {
						domain := strings.TrimPrefix(allowedOrigin, "*")
						if strings.HasSuffix(origin, domain) {
							c.Header("Access-Control-Allow-Origin", origin)
							break
						}
					}
				}
			}
		}

		// Set other CORS headers
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Request-ID, X-Auth-Token")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		c.Header("Access-Control-Max-Age", "86400") // 24 hours

		// Handle preflight requests
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
