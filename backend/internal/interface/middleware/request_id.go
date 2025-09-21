// Package middleware provides HTTP middleware implementations for the FinanceTracker API.
package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RequestID returns a middleware that generates or extracts a unique request ID.
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if request ID exists in header
		requestID := c.GetHeader("X-Request-ID")

		// Generate new ID if not provided
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// Set request ID in context
		c.Set("RequestID", requestID)

		// Set request ID in response header
		c.Header("X-Request-ID", requestID)

		// Continue processing
		c.Next()
	}
}
