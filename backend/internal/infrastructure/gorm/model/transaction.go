// Package model defines GORM models for the FinanceTracker application.
package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// TransactionType represents the type of transaction.
type TransactionType string

// TransactionType constants define the types of financial transactions.
const (
	TransactionTypeIncome  TransactionType = "INCOME"
	TransactionTypeExpense TransactionType = "EXPENSE"
)

// Transaction represents a financial transaction.
type Transaction struct {
	Base
	UserID      uuid.UUID       `gorm:"type:uuid;not null;index"`
	AccountID   uuid.UUID       `gorm:"type:uuid;not null;index"`
	CategoryID  uuid.UUID       `gorm:"type:uuid;not null;index"`
	Amount      decimal.Decimal `gorm:"type:decimal(15,2);not null"`
	Type        TransactionType `gorm:"type:varchar(20);not null"`
	Date        time.Time       `gorm:"type:date;not null;index"`
	Description *string         `gorm:"type:text"`

	// Relations
	User     User     `gorm:"foreignKey:UserID"`
	Account  Account  `gorm:"foreignKey:AccountID"`
	Category Category `gorm:"foreignKey:CategoryID"`
}

// TableName specifies the table name for Transaction.
func (Transaction) TableName() string {
	return "transactions"
}
