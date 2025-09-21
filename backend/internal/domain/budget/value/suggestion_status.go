// Package value 予算提案関連の値オブジェクトを定義
package value

import (
	"fmt"

	"financetracker/internal/domain/common"
)

// SuggestionStatus 予算提案ステータスを表す値オブジェクト
type SuggestionStatus struct {
	value string
}

// 定数定義
const (
	SuggestionStatusPending  = "pending"  // 検討中
	SuggestionStatusAccepted = "accepted" // 採用
	SuggestionStatusRejected = "rejected" // 却下
)

// NewSuggestionStatus 新しいSuggestionStatusを作成
func NewSuggestionStatus(value string) (SuggestionStatus, error) {
	if err := validateSuggestionStatus(value); err != nil {
		return SuggestionStatus{}, err
	}
	return SuggestionStatus{value: value}, nil
}

// validateSuggestionStatus ステータスの妥当性を検証
func validateSuggestionStatus(value string) error {
	switch value {
	case SuggestionStatusPending, SuggestionStatusAccepted, SuggestionStatusRejected:
		return nil
	default:
		return common.NewDomainError(
			common.DomainErrorTypeInvalidValue,
			fmt.Sprintf("無効な提案ステータスです: %s", value),
		)
	}
}

// String 文字列表現を返す
func (s SuggestionStatus) String() string {
	return s.value
}

// IsPending 検討中かどうかを判定
func (s SuggestionStatus) IsPending() bool {
	return s.value == SuggestionStatusPending
}

// IsAccepted 採用されたかどうかを判定
func (s SuggestionStatus) IsAccepted() bool {
	return s.value == SuggestionStatusAccepted
}

// IsRejected 却下されたかどうかを判定
func (s SuggestionStatus) IsRejected() bool {
	return s.value == SuggestionStatusRejected
}

// IsFinal 最終状態（採用または却下）かどうかを判定
func (s SuggestionStatus) IsFinal() bool {
	return s.value == SuggestionStatusAccepted || s.value == SuggestionStatusRejected
}

// CanTransitionTo 指定されたステータスへの遷移が可能かチェック
func (s SuggestionStatus) CanTransitionTo(newStatus SuggestionStatus) bool {
	// 最終状態からの遷移は不可
	if s.IsFinal() {
		return false
	}

	// pendingからaccepted/rejectedへの遷移のみ許可
	if s.IsPending() && (newStatus.IsAccepted() || newStatus.IsRejected()) {
		return true
	}

	return false
}

// GetDisplayName 表示用名称を取得
func (s SuggestionStatus) GetDisplayName() string {
	switch s.value {
	case SuggestionStatusPending:
		return "検討中"
	case SuggestionStatusAccepted:
		return "採用"
	case SuggestionStatusRejected:
		return "却下"
	default:
		return ""
	}
}

// Equals 等価性をチェック
func (s SuggestionStatus) Equals(other SuggestionStatus) bool {
	return s.value == other.value
}
