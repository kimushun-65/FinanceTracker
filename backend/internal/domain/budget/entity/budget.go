// Package entity 予算関連のエンティティを定義
package entity

import (
	"time"

	"financetracker/internal/domain/budget/value"
	"financetracker/internal/domain/common"
	commonValue "financetracker/internal/domain/common/value"

	"github.com/google/uuid"
)

// Budget 予算エンティティ
type Budget struct {
	common.BaseEntity
	userID     uuid.UUID
	categoryID uuid.UUID
	amount     commonValue.Money
	periodType value.PeriodType
	startDate  time.Time
	endDate    *time.Time
	isActive   bool
}

// NewBudget 新しい予算を作成
func NewBudget(
	userID uuid.UUID,
	categoryID uuid.UUID,
	amount commonValue.Money,
	periodType value.PeriodType,
	startDate time.Time,
	endDate *time.Time,
) (Budget, error) {
	if err := validateBudgetParams(userID, categoryID, amount, startDate, endDate); err != nil {
		return Budget{}, err
	}

	baseEntity := common.NewBaseEntity()

	return Budget{
		BaseEntity: baseEntity,
		userID:     userID,
		categoryID: categoryID,
		amount:     amount,
		periodType: periodType,
		startDate:  startDate,
		endDate:    endDate,
		isActive:   true, // 新規作成時はアクティブ
	}, nil
}

// validateBudgetParams 予算パラメータの妥当性を検証
func validateBudgetParams(
	userID uuid.UUID,
	categoryID uuid.UUID,
	amount commonValue.Money,
	startDate time.Time,
	endDate *time.Time,
) error {
	if userID == uuid.Nil {
		return common.NewDomainError(
			common.DomainErrorTypeInvalidValue,
			"ユーザーIDが必要です",
		)
	}

	if categoryID == uuid.Nil {
		return common.NewDomainError(
			common.DomainErrorTypeInvalidValue,
			"カテゴリIDが必要です",
		)
	}

	if amount.IsZero() || amount.IsNegative() {
		return common.NewDomainError(
			common.DomainErrorTypeInvalidValue,
			"予算額は0より大きい値である必要があります",
		)
	}

	if startDate.IsZero() {
		return common.NewDomainError(
			common.DomainErrorTypeInvalidValue,
			"開始日が必要です",
		)
	}

	if endDate != nil && endDate.Before(startDate) {
		return common.NewDomainError(
			common.DomainErrorTypeInvalidValue,
			"終了日は開始日より後である必要があります",
		)
	}

	return nil
}

// UserID ユーザーIDを取得
func (b Budget) UserID() uuid.UUID {
	return b.userID
}

// CategoryID カテゴリIDを取得
func (b Budget) CategoryID() uuid.UUID {
	return b.categoryID
}

// Amount 予算額を取得
func (b Budget) Amount() commonValue.Money {
	return b.amount
}

// PeriodType 期間タイプを取得
func (b Budget) PeriodType() value.PeriodType {
	return b.periodType
}

// StartDate 開始日を取得
func (b Budget) StartDate() time.Time {
	return b.startDate
}

// EndDate 終了日を取得
func (b Budget) EndDate() *time.Time {
	return b.endDate
}

// IsActive アクティブかどうかを判定
func (b Budget) IsActive() bool {
	return b.isActive
}

// UpdateAmount 予算額を更新
func (b *Budget) UpdateAmount(amount commonValue.Money) error {
	if amount.IsZero() || amount.IsNegative() {
		return common.NewDomainError(
			common.DomainErrorTypeInvalidValue,
			"予算額は0より大きい値である必要があります",
		)
	}

	b.amount = amount
	b.UpdateTimestamp()
	return nil
}

// UpdatePeriod 期間を更新
func (b *Budget) UpdatePeriod(startDate time.Time, endDate *time.Time) error {
	if err := validateBudgetParams(b.userID, b.categoryID, b.amount, startDate, endDate); err != nil {
		return err
	}

	b.startDate = startDate
	b.endDate = endDate
	b.UpdateTimestamp()
	return nil
}

// Activate 予算を有効化
func (b *Budget) Activate() {
	if !b.isActive {
		b.isActive = true
		b.UpdateTimestamp()
	}
}

// Deactivate 予算を無効化
func (b *Budget) Deactivate() {
	if b.isActive {
		b.isActive = false
		b.UpdateTimestamp()
	}
}

// IsValidForDate 指定日時が予算期間内かチェック
func (b Budget) IsValidForDate(date time.Time) bool {
	if date.Before(b.startDate) {
		return false
	}

	if b.endDate != nil && date.After(*b.endDate) {
		return false
	}

	return true
}
