package common

import (
	"fmt"

	"github.com/google/uuid"
)

// DomainError ドメイン固有のエラーを表現するインターフェース
type DomainError interface {
	error
	ErrorCode() string
	ErrorType() string
}

// NotFoundError リソースが見つからない場合のエラー
type NotFoundError struct {
	Resource string
	ID       uuid.UUID
	Message  string
}

func (e NotFoundError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return fmt.Sprintf("%s with ID %s not found", e.Resource, e.ID.String())
}

func (e NotFoundError) ErrorCode() string {
	return "NOT_FOUND"
}

func (e NotFoundError) ErrorType() string {
	return "NotFoundError"
}

// NewNotFoundError 新しいNotFoundErrorを作成
func NewNotFoundError(resource string, id uuid.UUID) *NotFoundError {
	return &NotFoundError{
		Resource: resource,
		ID:       id,
	}
}

// NewNotFoundErrorWithMessage メッセージ付きのNotFoundErrorを作成
func NewNotFoundErrorWithMessage(resource string, id uuid.UUID, message string) *NotFoundError {
	return &NotFoundError{
		Resource: resource,
		ID:       id,
		Message:  message,
	}
}

// ValidationError バリデーションエラー
type ValidationError struct {
	Field   string
	Value   interface{}
	Message string
}

func (e ValidationError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return fmt.Sprintf("validation failed for field '%s' with value '%v'", e.Field, e.Value)
}

func (e ValidationError) ErrorCode() string {
	return "VALIDATION_ERROR"
}

func (e ValidationError) ErrorType() string {
	return "ValidationError"
}

// NewValidationError 新しいValidationErrorを作成
func NewValidationError(field string, value interface{}, message string) *ValidationError {
	return &ValidationError{
		Field:   field,
		Value:   value,
		Message: message,
	}
}

// ConflictError リソースの競合エラー
type ConflictError struct {
	Resource string
	Field    string
	Value    interface{}
	Message  string
}

func (e ConflictError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return fmt.Sprintf("conflict in %s: %s '%v' already exists", e.Resource, e.Field, e.Value)
}

func (e ConflictError) ErrorCode() string {
	return "CONFLICT"
}

func (e ConflictError) ErrorType() string {
	return "ConflictError"
}

// NewConflictError 新しいConflictErrorを作成
func NewConflictError(resource, field string, value interface{}) *ConflictError {
	return &ConflictError{
		Resource: resource,
		Field:    field,
		Value:    value,
	}
}

// NewConflictErrorWithMessage メッセージ付きのConflictErrorを作成
func NewConflictErrorWithMessage(resource, field string, value interface{}, message string) *ConflictError {
	return &ConflictError{
		Resource: resource,
		Field:    field,
		Value:    value,
		Message:  message,
	}
}

// BusinessRuleError ビジネスルール違反エラー
type BusinessRuleError struct {
	Rule    string
	Message string
}

func (e BusinessRuleError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return fmt.Sprintf("business rule violation: %s", e.Rule)
}

func (e BusinessRuleError) ErrorCode() string {
	return "BUSINESS_RULE_VIOLATION"
}

func (e BusinessRuleError) ErrorType() string {
	return "BusinessRuleError"
}

// NewBusinessRuleError 新しいBusinessRuleErrorを作成
func NewBusinessRuleError(rule, message string) *BusinessRuleError {
	return &BusinessRuleError{
		Rule:    rule,
		Message: message,
	}
}
