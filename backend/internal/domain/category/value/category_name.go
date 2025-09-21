// Package value カテゴリ関連の値オブジェクトを定義
package value

import (
	"strings"

	"financetracker/internal/domain/common"
)

// CategoryName カテゴリ名を表す値オブジェクト
type CategoryName struct {
	value string
}

// 定数定義
const (
	minCategoryNameLength = 1   // カテゴリ名の最小文字数
	maxCategoryNameLength = 100 // カテゴリ名の最大文字数
)

// NewCategoryName 新しいCategoryNameを作成
func NewCategoryName(value string) (CategoryName, error) {
	if err := validateCategoryName(value); err != nil {
		return CategoryName{}, err
	}
	// 前後の空白を削除
	normalized := strings.TrimSpace(value)
	return CategoryName{value: normalized}, nil
}

// validateCategoryName カテゴリ名の妥当性を検証
func validateCategoryName(value string) error {
	normalized := strings.TrimSpace(value)

	if len(normalized) < minCategoryNameLength {
		return common.NewDomainError(
			common.DomainErrorTypeInvalidValue,
			"カテゴリ名は必須です",
		)
	}

	if len(normalized) > maxCategoryNameLength {
		return common.NewDomainError(
			common.DomainErrorTypeInvalidValue,
			"カテゴリ名が長すぎます（最大100文字）",
		)
	}

	// 空白のみは不可
	if strings.TrimSpace(value) == "" {
		return common.NewDomainError(
			common.DomainErrorTypeInvalidValue,
			"カテゴリ名は空白のみにすることはできません",
		)
	}

	return nil
}

// String 文字列表現を返す
func (n CategoryName) String() string {
	return n.value
}

// Equals 等価性をチェック
func (n CategoryName) Equals(other CategoryName) bool {
	return n.value == other.value
}
