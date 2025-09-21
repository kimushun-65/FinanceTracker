// Package model defines GORM models for the FinanceTracker application.
package model

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// AccountType represents the type of account.
type AccountType string

const (
	AccountTypeCash       AccountType = "CASH"
	AccountTypeBank       AccountType = "BANK"
	AccountTypeCreditCard AccountType = "CREDIT_CARD"
	AccountTypeInvestment AccountType = "INVESTMENT"
	AccountTypeLoan       AccountType = "LOAN"
	AccountTypeOther      AccountType = "OTHER"
)

// Account represents a financial account.
type Account struct {
	Base
	UserID   uuid.UUID       `gorm:"type:uuid;not null;index"`
	Name     string          `gorm:"type:varchar(255);not null"`
	Type     AccountType     `gorm:"type:varchar(50);not null"`
	Balance  decimal.Decimal `gorm:"type:decimal(15,2);not null;default:0"`
	Currency string          `gorm:"type:varchar(3);not null;default:'JPY'"`
	IsActive bool            `gorm:"not null;default:true"`

	// Relations
	User         User          `gorm:"foreignKey:UserID"`
	Transactions []Transaction `gorm:"foreignKey:AccountID;constraint:OnDelete:RESTRICT"`
}

// TableName specifies the table name for Account.
func (Account) TableName() string {
	return "accounts"
}