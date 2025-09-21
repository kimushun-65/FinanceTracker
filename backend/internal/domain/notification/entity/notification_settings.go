// Package entity 通知設定関連のエンティティを定義
package entity

import (
	"time"

	"financetracker/internal/domain/common"
	commonValue "financetracker/internal/domain/common/value"

	"github.com/google/uuid"
)

// NotificationSettings 通知設定エンティティ
type NotificationSettings struct {
	common.BaseEntity
	userID                     uuid.UUID
	monthlyReportEnabled       bool
	monthlyReportSendDay       int
	monthlyReportSendTime      time.Time
	monthlyReportEmail         commonValue.Email
	budgetExceededAlertEnabled bool
	budgetExceededAlertEmail   commonValue.Email
}

// NewNotificationSettings 新しい通知設定を作成
func NewNotificationSettings(
	userID uuid.UUID,
	monthlyReportEnabled bool,
	monthlyReportSendDay int,
	monthlyReportSendTime time.Time,
	monthlyReportEmail commonValue.Email,
	budgetExceededAlertEnabled bool,
	budgetExceededAlertEmail commonValue.Email,
) (NotificationSettings, error) {
	if err := validateNotificationSettingsParams(
		userID, monthlyReportSendDay,
	); err != nil {
		return NotificationSettings{}, err
	}

	baseEntity := common.NewBaseEntity()

	return NotificationSettings{
		BaseEntity:                 baseEntity,
		userID:                     userID,
		monthlyReportEnabled:       monthlyReportEnabled,
		monthlyReportSendDay:       monthlyReportSendDay,
		monthlyReportSendTime:      monthlyReportSendTime,
		monthlyReportEmail:         monthlyReportEmail,
		budgetExceededAlertEnabled: budgetExceededAlertEnabled,
		budgetExceededAlertEmail:   budgetExceededAlertEmail,
	}, nil
}

// validateNotificationSettingsParams 通知設定パラメータの妥当性を検証
func validateNotificationSettingsParams(
	userID uuid.UUID,
	monthlyReportSendDay int,
) error {
	if userID == uuid.Nil {
		return common.NewDomainError(
			common.DomainErrorTypeInvalidValue,
			"ユーザーIDが必要です",
		)
	}

	if monthlyReportSendDay < 1 || monthlyReportSendDay > 31 {
		return common.NewDomainError(
			common.DomainErrorTypeInvalidValue,
			"送信日は1〜31の範囲である必要があります",
		)
	}

	return nil
}

// UserID ユーザーIDを取得
func (n NotificationSettings) UserID() uuid.UUID {
	return n.userID
}

// IsMonthlyReportEnabled 月次レポートが有効かどうかを取得
func (n NotificationSettings) IsMonthlyReportEnabled() bool {
	return n.monthlyReportEnabled
}

// MonthlyReportSendDay 月次レポート送信日を取得
func (n NotificationSettings) MonthlyReportSendDay() int {
	return n.monthlyReportSendDay
}

// MonthlyReportSendTime 月次レポート送信時刻を取得
func (n NotificationSettings) MonthlyReportSendTime() time.Time {
	return n.monthlyReportSendTime
}

// MonthlyReportEmail 月次レポート送信先メールアドレスを取得
func (n NotificationSettings) MonthlyReportEmail() commonValue.Email {
	return n.monthlyReportEmail
}

// IsBudgetExceededAlertEnabled 予算超過アラートが有効かどうかを取得
func (n NotificationSettings) IsBudgetExceededAlertEnabled() bool {
	return n.budgetExceededAlertEnabled
}

// BudgetExceededAlertEmail 予算超過アラート送信先メールアドレスを取得
func (n NotificationSettings) BudgetExceededAlertEmail() commonValue.Email {
	return n.budgetExceededAlertEmail
}

// UpdateMonthlyReportSettings 月次レポート設定を更新
func (n *NotificationSettings) UpdateMonthlyReportSettings(
	enabled bool,
	sendDay int,
	sendTime time.Time,
	email commonValue.Email,
) error {
	if sendDay < 1 || sendDay > 31 {
		return common.NewDomainError(
			common.DomainErrorTypeInvalidValue,
			"送信日は1〜31の範囲である必要があります",
		)
	}

	n.monthlyReportEnabled = enabled
	n.monthlyReportSendDay = sendDay
	n.monthlyReportSendTime = sendTime
	n.monthlyReportEmail = email
	n.UpdateTimestamp()

	return nil
}

// UpdateBudgetExceededAlertSettings 予算超過アラート設定を更新
func (n *NotificationSettings) UpdateBudgetExceededAlertSettings(
	enabled bool,
	email commonValue.Email,
) {
	n.budgetExceededAlertEnabled = enabled
	n.budgetExceededAlertEmail = email
	n.UpdateTimestamp()
}

// EnableMonthlyReport 月次レポートを有効化
func (n *NotificationSettings) EnableMonthlyReport() {
	n.monthlyReportEnabled = true
	n.UpdateTimestamp()
}

// DisableMonthlyReport 月次レポートを無効化
func (n *NotificationSettings) DisableMonthlyReport() {
	n.monthlyReportEnabled = false
	n.UpdateTimestamp()
}

// EnableBudgetExceededAlert 予算超過アラートを有効化
func (n *NotificationSettings) EnableBudgetExceededAlert() {
	n.budgetExceededAlertEnabled = true
	n.UpdateTimestamp()
}

// DisableBudgetExceededAlert 予算超過アラートを無効化
func (n *NotificationSettings) DisableBudgetExceededAlert() {
	n.budgetExceededAlertEnabled = false
	n.UpdateTimestamp()
}

// ShouldSendMonthlyReportToday 今日が月次レポート送信日かチェック
func (n NotificationSettings) ShouldSendMonthlyReportToday() bool {
	if !n.monthlyReportEnabled {
		return false
	}

	today := time.Now().Day()

	// 月末日を超える設定の場合、月末日に送信
	lastDayOfMonth := time.Now().AddDate(0, 1, -time.Now().Day()).Day()
	if n.monthlyReportSendDay > lastDayOfMonth {
		return today == lastDayOfMonth
	}

	return today == n.monthlyReportSendDay
}

// GetFormattedSendTime 送信時刻をHH:mm形式で取得
func (n NotificationSettings) GetFormattedSendTime() string {
	return n.monthlyReportSendTime.Format("15:04")
}
