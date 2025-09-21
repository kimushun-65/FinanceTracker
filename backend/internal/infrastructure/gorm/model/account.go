// Package model defines GORM models for the FinanceTracker application.
package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// AccountType represents the type of account.
type AccountType string

// AccountType constants define the types of accounts supported in the system.
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
	User         User              `gorm:"foreignKey:UserID"`
	Transactions []Transaction     `gorm:"foreignKey:AccountID;constraint:OnDelete:RESTRICT"`
	Movements    []AccountMovement `gorm:"foreignKey:AccountID;constraint:OnDelete:CASCADE"`
}

// TableName specifies the table name for Account.
func (Account) TableName() string {
	return "accounts"
}

// AccountMovement represents a movement in an account (deposit/withdrawal).
type AccountMovement struct {
	Base
	UserID     uuid.UUID       `gorm:"type:uuid;not null;index"`
	AccountID  uuid.UUID       `gorm:"type:uuid;not null;index"`
	Amount     decimal.Decimal `gorm:"type:decimal(15,2);not null"`
	OccurredAt time.Time       `gorm:"not null;index"`
	Note       string          `gorm:"type:varchar(255)"`

	// Relations
	User    User    `gorm:"foreignKey:UserID"`
	Account Account `gorm:"foreignKey:AccountID"`
}

// TableName specifies the table name for AccountMovement.
func (AccountMovement) TableName() string {
	return "account_movements"
}
