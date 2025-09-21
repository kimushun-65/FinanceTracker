// Package middleware provides HTTP middleware implementations for the FinanceTracker API.
package middleware

import (
	"bytes"
	"io"
	"time"

	"financetracker/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// responseBodyWriter wraps gin.ResponseWriter to capture response body for logging.
type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// Write captures the response body.
func (w responseBodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// Logger returns a middleware that logs HTTP requests and responses.
func Logger(log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Read request body if needed (be careful with large bodies)
		var requestBody []byte
		if c.Request.Body != nil && c.Request.ContentLength > 0 && c.Request.ContentLength < 1024*1024 { // Only log bodies < 1MB
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// Create response body writer
		blw := &responseBodyWriter{
			ResponseWriter: c.Writer,
			body:           bytes.NewBuffer(nil),
		}
		c.Writer = blw

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Get request ID
		requestID := c.GetString("RequestID")
		if requestID == "" {
			requestID = c.GetHeader("X-Request-ID")
		}

		// Get user ID if authenticated
		userID := c.GetString("UserID")

		// Build log fields
		fields := []zap.Field{
			zap.String("request_id", requestID),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", raw),
			zap.Int("status", c.Writer.Status()),
			zap.Duration("latency", latency),
			zap.String("client_ip", c.ClientIP()),
			zap.String("user_agent", c.Request.UserAgent()),
			zap.Int("body_size", c.Writer.Size()),
		}

		// Add user ID if available
		if userID != "" {
			fields = append(fields, zap.String("user_id", userID))
		}

		// Add error if exists
		if len(c.Errors) > 0 {
			fields = append(fields, zap.String("error", c.Errors.String()))
		}

		// Log based on status code
		switch {
		case c.Writer.Status() >= 500:
			log.Error("Server error", fields...)
			// Include request/response bodies for server errors
			if len(requestBody) > 0 {
				log.Debug("Request body", zap.ByteString("request_body", requestBody))
			}
			if blw.body.Len() > 0 && blw.body.Len() < 1024*1024 { // Only log response bodies < 1MB
				log.Debug("Response body", zap.String("response_body", blw.body.String()))
			}
		case c.Writer.Status() >= 400:
			log.Warn("Client error", fields...)
		case c.Writer.Status() >= 300:
			log.Info("Redirection", fields...)
		default:
			log.Info("Request completed", fields...)
		}

		// Log slow requests
		if latency > 5*time.Second {
			log.Warn("Slow request detected",
				zap.String("request_id", requestID),
				zap.String("method", c.Request.Method),
				zap.String("path", path),
				zap.Duration("latency", latency),
			)
		}
	}
}