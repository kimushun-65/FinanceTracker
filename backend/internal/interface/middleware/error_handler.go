// Package middleware provides HTTP middleware implementations for the FinanceTracker API.
package middleware

import (
	"net/http"

	"financetracker/pkg/errors"
	"financetracker/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ErrorHandler returns a middleware that handles errors in a consistent format.
func ErrorHandler(log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Process request
		c.Next()

		// Check if there are any errors
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			requestID := c.GetString("RequestID")

			// Check if it's an AppError
			if appErr := errors.GetAppError(err); appErr != nil {
				// Log the error with appropriate level
				logFields := []zap.Field{
					zap.String("request_id", requestID),
					zap.String("error_code", appErr.Code),
					zap.Error(appErr.Internal),
				}

				if appErr.StatusCode >= 500 {
					log.Error("Application error", logFields...)
				} else {
					log.Warn("Client error", logFields...)
				}

				// Build error response
				response := gin.H{
					"code":       appErr.Code,
					"message":    appErr.Message,
					"request_id": requestID,
				}

				// Add detail if present
				if appErr.Detail != "" {
					response["detail"] = appErr.Detail
				}

				// Send response
				c.JSON(appErr.StatusCode, response)
				return
			}

			// Handle validation errors
			if c.Writer.Status() == http.StatusBadRequest {
				log.Warn("Validation error",
					zap.String("request_id", requestID),
					zap.Error(err),
				)

				c.JSON(http.StatusBadRequest, gin.H{
					"code":       "VALIDATION_ERROR",
					"message":    "Invalid request data",
					"detail":     err.Error(),
					"request_id": requestID,
				})
				return
			}

			// Handle unexpected errors
			log.Error("Unexpected error",
				zap.String("request_id", requestID),
				zap.Error(err),
				zap.String("path", c.Request.URL.Path),
				zap.String("method", c.Request.Method),
			)

			// Don't expose internal error details in production
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":       "INTERNAL_ERROR",
				"message":    "An unexpected error occurred",
				"request_id": requestID,
			})
		}
	}
}

// Recovery returns a middleware that recovers from panics and returns a 500 error.
func Recovery(log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				requestID := c.GetString("RequestID")

				// Log the panic
				log.Error("Panic recovered",
					zap.String("request_id", requestID),
					zap.Any("error", err),
					zap.String("path", c.Request.URL.Path),
					zap.String("method", c.Request.Method),
					zap.Stack("stack"),
				)

				// Return 500 error
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"code":       "INTERNAL_ERROR",
					"message":    "An unexpected error occurred",
					"request_id": requestID,
				})
			}
		}()

		c.Next()
	}
}
