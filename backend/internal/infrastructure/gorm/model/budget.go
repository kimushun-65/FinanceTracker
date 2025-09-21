// Package model defines GORM models for the FinanceTracker application.
package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// PeriodType represents the budget period type.
type PeriodType string

// PeriodType constants define the budget period types for budget planning.
const (
	PeriodTypeMonthly PeriodType = "MONTHLY"
	PeriodTypeYearly  PeriodType = "YEARLY"
)

// Budget represents a budget for a category.
type Budget struct {
	Base
	UserID     uuid.UUID       `gorm:"type:uuid;not null;index"`
	CategoryID uuid.UUID       `gorm:"type:uuid;not null;index"`
	Amount     decimal.Decimal `gorm:"type:decimal(15,2);not null"`
	PeriodType PeriodType      `gorm:"type:varchar(20);not null"`
	StartDate  time.Time       `gorm:"type:date;not null"`
	EndDate    *time.Time      `gorm:"type:date"`
	IsActive   bool            `gorm:"not null;default:true"`

	// Relations
	User     User     `gorm:"foreignKey:UserID"`
	Category Category `gorm:"foreignKey:CategoryID"`
}

// TableName specifies the table name for Budget.
func (Budget) TableName() string {
	return "budgets"
}

// BudgetSuggestion represents an AI-generated budget suggestion.
type BudgetSuggestion struct {
	Base
	UserID          uuid.UUID       `gorm:"type:uuid;not null;index"`
	CategoryID      uuid.UUID       `gorm:"type:uuid;not null;index"`
	SuggestedAmount decimal.Decimal `gorm:"type:decimal(15,2);not null"`
	CurrentAverage  decimal.Decimal `gorm:"type:decimal(15,2);not null"`
	ConfidenceScore decimal.Decimal `gorm:"type:decimal(3,2);not null"` // 0.00 to 1.00
	Reason          *string         `gorm:"type:text"`
	PeriodType      PeriodType      `gorm:"type:varchar(20);not null;default:'MONTHLY'"`
	GeneratedAt     time.Time       `gorm:"not null;default:CURRENT_TIMESTAMP"`
	AppliedAt       *time.Time      `gorm:"type:timestamptz"`

	// Relations
	User     User     `gorm:"foreignKey:UserID"`
	Category Category `gorm:"foreignKey:CategoryID"`
}

// TableName specifies the table name for BudgetSuggestion.
func (BudgetSuggestion) TableName() string {
	return "budget_suggestions"
}
