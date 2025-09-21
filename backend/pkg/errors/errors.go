// Package errors provides custom error types and error handling utilities for the FinanceTracker application.
// It defines domain-specific errors with appropriate HTTP status codes and user-friendly messages.
package errors

import (
	"errors"
	"fmt"
	"net/http"
)

// AppError represents an application-specific error with HTTP status code and user-friendly message.
type AppError struct {
	Code       string `json:"code"`                 // Error code for client identification
	Message    string `json:"message"`              // User-friendly error message
	Detail     string `json:"detail,omitempty"`     // Detailed error information (optional)
	StatusCode int    `json:"-"`                    // HTTP status code
	Internal   error  `json:"-"`                    // Internal error for logging
}

// Error implements the error interface.
func (e *AppError) Error() string {
	if e.Internal != nil {
		return fmt.Sprintf("%s: %s (internal: %v)", e.Code, e.Message, e.Internal)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap returns the internal error.
func (e *AppError) Unwrap() error {
	return e.Internal
}

// Common error codes
const (
	// Client errors
	ErrCodeValidation     = "VALIDATION_ERROR"
	ErrCodeNotFound       = "NOT_FOUND"
	ErrCodeUnauthorized   = "UNAUTHORIZED"
	ErrCodeForbidden      = "FORBIDDEN"
	ErrCodeConflict       = "CONFLICT"
	ErrCodeBadRequest     = "BAD_REQUEST"
	ErrCodeTooManyRequest = "TOO_MANY_REQUESTS"

	// Server errors
	ErrCodeInternal       = "INTERNAL_ERROR"
	ErrCodeDatabase       = "DATABASE_ERROR"
	ErrCodeExternal       = "EXTERNAL_SERVICE_ERROR"
	ErrCodeUnavailable    = "SERVICE_UNAVAILABLE"
	ErrCodeNotImplemented = "NOT_IMPLEMENTED"
)

// Error constructors

// NewValidationError creates a validation error.
func NewValidationError(message string, detail ...string) *AppError {
	err := &AppError{
		Code:       ErrCodeValidation,
		Message:    message,
		StatusCode: http.StatusBadRequest,
	}
	if len(detail) > 0 {
		err.Detail = detail[0]
	}
	return err
}

// NewNotFoundError creates a not found error.
func NewNotFoundError(resource string) *AppError {
	return &AppError{
		Code:       ErrCodeNotFound,
		Message:    fmt.Sprintf("%s not found", resource),
		StatusCode: http.StatusNotFound,
	}
}

// NewUnauthorizedError creates an unauthorized error.
func NewUnauthorizedError(message string) *AppError {
	if message == "" {
		message = "Authentication required"
	}
	return &AppError{
		Code:       ErrCodeUnauthorized,
		Message:    message,
		StatusCode: http.StatusUnauthorized,
	}
}

// NewForbiddenError creates a forbidden error.
func NewForbiddenError(message string) *AppError {
	if message == "" {
		message = "Access denied"
	}
	return &AppError{
		Code:       ErrCodeForbidden,
		Message:    message,
		StatusCode: http.StatusForbidden,
	}
}

// NewConflictError creates a conflict error.
func NewConflictError(message string) *AppError {
	return &AppError{
		Code:       ErrCodeConflict,
		Message:    message,
		StatusCode: http.StatusConflict,
	}
}

// NewBadRequestError creates a bad request error.
func NewBadRequestError(message string) *AppError {
	return &AppError{
		Code:       ErrCodeBadRequest,
		Message:    message,
		StatusCode: http.StatusBadRequest,
	}
}

// NewTooManyRequestsError creates a rate limit error.
func NewTooManyRequestsError(message string) *AppError {
	if message == "" {
		message = "Too many requests, please try again later"
	}
	return &AppError{
		Code:       ErrCodeTooManyRequest,
		Message:    message,
		StatusCode: http.StatusTooManyRequests,
	}
}

// NewInternalError creates an internal server error.
func NewInternalError(message string, internal error) *AppError {
	if message == "" {
		message = "An internal error occurred"
	}
	return &AppError{
		Code:       ErrCodeInternal,
		Message:    message,
		StatusCode: http.StatusInternalServerError,
		Internal:   internal,
	}
}

// NewDatabaseError creates a database error.
func NewDatabaseError(message string, internal error) *AppError {
	if message == "" {
		message = "A database error occurred"
	}
	return &AppError{
		Code:       ErrCodeDatabase,
		Message:    message,
		StatusCode: http.StatusInternalServerError,
		Internal:   internal,
	}
}

// NewExternalServiceError creates an external service error.
func NewExternalServiceError(service string, internal error) *AppError {
	return &AppError{
		Code:       ErrCodeExternal,
		Message:    fmt.Sprintf("External service error: %s", service),
		StatusCode: http.StatusServiceUnavailable,
		Internal:   internal,
	}
}

// NewServiceUnavailableError creates a service unavailable error.
func NewServiceUnavailableError(message string) *AppError {
	if message == "" {
		message = "Service temporarily unavailable"
	}
	return &AppError{
		Code:       ErrCodeUnavailable,
		Message:    message,
		StatusCode: http.StatusServiceUnavailable,
	}
}

// NewNotImplementedError creates a not implemented error.
func NewNotImplementedError(feature string) *AppError {
	return &AppError{
		Code:       ErrCodeNotImplemented,
		Message:    fmt.Sprintf("Feature not implemented: %s", feature),
		StatusCode: http.StatusNotImplemented,
	}
}

// Helper functions

// IsAppError checks if an error is an AppError.
func IsAppError(err error) bool {
	var appErr *AppError
	return errors.As(err, &appErr)
}

// GetAppError extracts AppError from an error.
// If the error is not an AppError, it returns nil.
func GetAppError(err error) *AppError {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr
	}
	return nil
}

// GetStatusCode returns the HTTP status code from an error.
// If the error is not an AppError, it returns 500.
func GetStatusCode(err error) int {
	if appErr := GetAppError(err); appErr != nil {
		return appErr.StatusCode
	}
	return http.StatusInternalServerError
}

// Wrap wraps an error with additional context.
func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}
	
	if appErr := GetAppError(err); appErr != nil {
		// If it's already an AppError, preserve its properties
		wrapped := *appErr
		wrapped.Message = fmt.Sprintf("%s: %s", message, appErr.Message)
		return &wrapped
	}
	
	// Otherwise create a new internal error
	return NewInternalError(message, err)
}

// Domain-specific errors

// ErrUserNotFound is returned when a user is not found.
var ErrUserNotFound = NewNotFoundError("user")

// ErrAccountNotFound is returned when an account is not found.
var ErrAccountNotFound = NewNotFoundError("account")

// ErrTransactionNotFound is returned when a transaction is not found.
var ErrTransactionNotFound = NewNotFoundError("transaction")

// ErrCategoryNotFound is returned when a category is not found.
var ErrCategoryNotFound = NewNotFoundError("category")

// ErrBudgetNotFound is returned when a budget is not found.
var ErrBudgetNotFound = NewNotFoundError("budget")

// ErrInvalidCredentials is returned when login credentials are invalid.
var ErrInvalidCredentials = NewUnauthorizedError("Invalid email or password")

// ErrInsufficientBalance is returned when account balance is insufficient.
var ErrInsufficientBalance = NewBadRequestError("Insufficient account balance")

// ErrDuplicateEmail is returned when attempting to register with an existing email.
var ErrDuplicateEmail = NewConflictError("Email already registered")

// ErrInvalidToken is returned when a token is invalid or expired.
var ErrInvalidToken = NewUnauthorizedError("Invalid or expired token")