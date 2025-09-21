// Package model defines GORM models for the FinanceTracker application.
package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// AssetSnapshot represents a daily snapshot of user's assets.
type AssetSnapshot struct {
	Base
	UserID           uuid.UUID       `gorm:"type:uuid;not null;index"`
	Date             time.Time       `gorm:"type:date;not null;index"`
	TotalAssets      decimal.Decimal `gorm:"type:decimal(15,2);not null"`
	TotalLiabilities decimal.Decimal `gorm:"type:decimal(15,2);not null"`
	NetWorth         decimal.Decimal `gorm:"type:decimal(15,2);not null"`

	// Relations
	User User `gorm:"foreignKey:UserID"`
}

// TableName specifies the table name for AssetSnapshot.
func (AssetSnapshot) TableName() string {
	return "asset_snapshots"
}

// ForecastMethod represents the forecasting method used.
type ForecastMethod string

const (
	ForecastMethodLinear      ForecastMethod = "LINEAR"
	ForecastMethodExponential ForecastMethod = "EXPONENTIAL"
	ForecastMethodSeasonal    ForecastMethod = "SEASONAL"
)

// AssetForecast represents a future asset forecast.
type AssetForecast struct {
	Base
	UserID           uuid.UUID       `gorm:"type:uuid;not null;index"`
	ForecastDate     time.Time       `gorm:"type:date;not null;index"`
	ForecastedAmount decimal.Decimal `gorm:"type:decimal(15,2);not null"`
	ConfidenceLevel  decimal.Decimal `gorm:"type:decimal(3,2);not null"` // 0.00 to 1.00
	ForecastMethod   ForecastMethod  `gorm:"type:varchar(50);not null"`
	GeneratedAt      time.Time       `gorm:"not null;default:CURRENT_TIMESTAMP"`

	// Relations
	User User `gorm:"foreignKey:UserID"`
}

// TableName specifies the table name for AssetForecast.
func (AssetForecast) TableName() string {
	return "asset_forecasts"
}