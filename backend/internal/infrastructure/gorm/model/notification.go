// Package model defines GORM models for the FinanceTracker application.
package model

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// NotificationType represents the type of notification.
type NotificationType string

const (
	NotificationTypeBudgetExceeded  NotificationType = "BUDGET_EXCEEDED"
	NotificationTypePaymentDue      NotificationType = "PAYMENT_DUE"
	NotificationTypeLowBalance      NotificationType = "LOW_BALANCE"
	NotificationTypeMonthlySummary  NotificationType = "MONTHLY_SUMMARY"
)

// NotificationChannel represents the notification delivery channel.
type NotificationChannel string

const (
	NotificationChannelEmail NotificationChannel = "EMAIL"
	NotificationChannelPush  NotificationChannel = "PUSH"
	NotificationChannelSMS   NotificationChannel = "SMS"
)

// NotificationSetting represents user's notification preferences.
type NotificationSetting struct {
	Base
	UserID           uuid.UUID           `gorm:"type:uuid;not null;index"`
	NotificationType NotificationType    `gorm:"type:varchar(50);not null"`
	IsEnabled        bool                `gorm:"not null;default:true"`
	Channel          NotificationChannel `gorm:"type:varchar(20);not null"`
	ThresholdValue   *decimal.Decimal    `gorm:"type:decimal(15,2)"`

	// Relations
	User User `gorm:"foreignKey:UserID"`
}

// TableName specifies the table name for NotificationSetting.
func (NotificationSetting) TableName() string {
	return "notification_settings"
}