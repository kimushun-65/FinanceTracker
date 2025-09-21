// Package value カテゴリ関連の値オブジェクトを定義
package value

import (
	"fmt"

	"financetracker/internal/domain/common"
)

// CategoryType カテゴリタイプを表す値オブジェクト
type CategoryType struct {
	value string
}

// 定数定義
const (
	CategoryTypeIncome  = "income"  // 収入カテゴリ
	CategoryTypeExpense = "expense" // 支出カテゴリ
)

// NewCategoryType 新しいCategoryTypeを作成
func NewCategoryType(value string) (CategoryType, error) {
	if err := validateCategoryType(value); err != nil {
		return CategoryType{}, err
	}
	return CategoryType{value: value}, nil
}

// validateCategoryType カテゴリタイプの妥当性を検証
func validateCategoryType(value string) error {
	switch value {
	case CategoryTypeIncome, CategoryTypeExpense:
		return nil
	default:
		return common.NewDomainError(
			common.DomainErrorTypeInvalidValue,
			fmt.Sprintf("無効なカテゴリタイプです: %s", value),
		)
	}
}

// String 文字列表現を返す
func (t CategoryType) String() string {
	return t.value
}

// IsIncome 収入カテゴリかどうかを判定
func (t CategoryType) IsIncome() bool {
	return t.value == CategoryTypeIncome
}

// IsExpense 支出カテゴリかどうかを判定
func (t CategoryType) IsExpense() bool {
	return t.value == CategoryTypeExpense
}

// GetDisplayName 表示用名称を取得
func (t CategoryType) GetDisplayName() string {
	switch t.value {
	case CategoryTypeIncome:
		return "収入"
	case CategoryTypeExpense:
		return "支出"
	default:
		return ""
	}
}

// Equals 等価性をチェック
func (t CategoryType) Equals(other CategoryType) bool {
	return t.value == other.value
}
