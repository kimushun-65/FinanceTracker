// Package entity 予算提案関連のエンティティを定義
package entity

import (
	"financetracker/internal/domain/budget/value"
	"financetracker/internal/domain/common"
	commonValue "financetracker/internal/domain/common/value"

	"github.com/google/uuid"
)

// BudgetSuggestion 予算提案エンティティ
type BudgetSuggestion struct {
	common.BaseEntity
	userID            uuid.UUID
	month             string // YYYY-MM形式
	categoryID        uuid.UUID
	suggestedAmount   commonValue.Money
	currentBudget     commonValue.Money
	lastMonthActual   commonValue.Money
	threeMonthAverage commonValue.Money
	status            value.SuggestionStatus
	reason            string
	confidence        float64
}

// NewBudgetSuggestion 新しい予算提案を作成
func NewBudgetSuggestion(
	userID uuid.UUID,
	month string,
	categoryID uuid.UUID,
	suggestedAmount commonValue.Money,
	currentBudget commonValue.Money,
	lastMonthActual commonValue.Money,
	threeMonthAverage commonValue.Money,
	reason string,
	confidence float64,
) (BudgetSuggestion, error) {
	if err := validateBudgetSuggestionParams(
		userID, month, categoryID, suggestedAmount, confidence,
	); err != nil {
		return BudgetSuggestion{}, err
	}

	// 新規作成時のステータスは必ずpending
	status, _ := value.NewSuggestionStatus(value.SuggestionStatusPending)

	baseEntity := common.NewBaseEntity()

	return BudgetSuggestion{
		BaseEntity:        baseEntity,
		userID:            userID,
		month:             month,
		categoryID:        categoryID,
		suggestedAmount:   suggestedAmount,
		currentBudget:     currentBudget,
		lastMonthActual:   lastMonthActual,
		threeMonthAverage: threeMonthAverage,
		status:            status,
		reason:            reason,
		confidence:        confidence,
	}, nil
}

// validateBudgetSuggestionParams 予算提案パラメータの妥当性を検証
func validateBudgetSuggestionParams(
	userID uuid.UUID,
	month string,
	categoryID uuid.UUID,
	suggestedAmount commonValue.Money,
	confidence float64,
) error {
	if userID == uuid.Nil {
		return common.NewDomainError(
			common.DomainErrorTypeInvalidValue,
			"ユーザーIDが必要です",
		)
	}

	if month == "" {
		return common.NewDomainError(
			common.DomainErrorTypeInvalidValue,
			"対象年月が必要です",
		)
	}

	// YYYY-MM形式の簡易チェック
	if len(month) != 7 || month[4] != '-' {
		return common.NewDomainError(
			common.DomainErrorTypeInvalidValue,
			"対象年月はYYYY-MM形式である必要があります",
		)
	}

	if categoryID == uuid.Nil {
		return common.NewDomainError(
			common.DomainErrorTypeInvalidValue,
			"カテゴリIDが必要です",
		)
	}

	if suggestedAmount.IsNegative() {
		return common.NewDomainError(
			common.DomainErrorTypeInvalidValue,
			"提案金額は0以上である必要があります",
		)
	}

	if confidence < 0 || confidence > 1 {
		return common.NewDomainError(
			common.DomainErrorTypeInvalidValue,
			"信頼度は0以上1以下である必要があります",
		)
	}

	return nil
}

// UserID ユーザーIDを取得
func (b BudgetSuggestion) UserID() uuid.UUID {
	return b.userID
}

// Month 対象年月を取得
func (b BudgetSuggestion) Month() string {
	return b.month
}

// CategoryID カテゴリIDを取得
func (b BudgetSuggestion) CategoryID() uuid.UUID {
	return b.categoryID
}

// SuggestedAmount 提案金額を取得
func (b BudgetSuggestion) SuggestedAmount() commonValue.Money {
	return b.suggestedAmount
}

// CurrentBudget 現在の予算額を取得
func (b BudgetSuggestion) CurrentBudget() commonValue.Money {
	return b.currentBudget
}

// LastMonthActual 前月実績を取得
func (b BudgetSuggestion) LastMonthActual() commonValue.Money {
	return b.lastMonthActual
}

// ThreeMonthAverage 3ヶ月平均を取得
func (b BudgetSuggestion) ThreeMonthAverage() commonValue.Money {
	return b.threeMonthAverage
}

// Status ステータスを取得
func (b BudgetSuggestion) Status() value.SuggestionStatus {
	return b.status
}

// Reason 提案理由を取得
func (b BudgetSuggestion) Reason() string {
	return b.reason
}

// Confidence 信頼度を取得
func (b BudgetSuggestion) Confidence() float64 {
	return b.confidence
}

// Accept 提案を採用
func (b *BudgetSuggestion) Accept() error {
	newStatus, err := value.NewSuggestionStatus(value.SuggestionStatusAccepted)
	if err != nil {
		return err
	}

	if !b.status.CanTransitionTo(newStatus) {
		return common.NewDomainError(
			common.DomainErrorTypeBusinessRule,
			"この提案は既に処理済みです",
		)
	}

	b.status = newStatus
	b.UpdateTimestamp()
	return nil
}

// Reject 提案を却下
func (b *BudgetSuggestion) Reject() error {
	newStatus, err := value.NewSuggestionStatus(value.SuggestionStatusRejected)
	if err != nil {
		return err
	}

	if !b.status.CanTransitionTo(newStatus) {
		return common.NewDomainError(
			common.DomainErrorTypeBusinessRule,
			"この提案は既に処理済みです",
		)
	}

	b.status = newStatus
	b.UpdateTimestamp()
	return nil
}

// IsPending 検討中かどうかを判定
func (b BudgetSuggestion) IsPending() bool {
	return b.status.IsPending()
}

// IsFinal 最終状態かどうかを判定
func (b BudgetSuggestion) IsFinal() bool {
	return b.status.IsFinal()
}
