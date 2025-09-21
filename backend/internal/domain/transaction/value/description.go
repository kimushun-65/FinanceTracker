// Package value トランザクション関連の値オブジェクトを定義
package value

import (
	"strings"

	"financetracker/internal/domain/common"
)

// Description トランザクションの説明を表す値オブジェクト
type Description struct {
	value string
}

// 定数定義
const (
	maxDescriptionLength = 500 // 説明の最大文字数
)

// NewDescription 新しいDescriptionを作成
func NewDescription(value string) (Description, error) {
	if err := validateDescription(value); err != nil {
		return Description{}, err
	}
	// 前後の空白を削除
	normalized := strings.TrimSpace(value)
	return Description{value: normalized}, nil
}

// validateDescription 説明の妥当性を検証
func validateDescription(value string) error {
	normalized := strings.TrimSpace(value)

	if len(normalized) > maxDescriptionLength {
		return common.NewDomainError(
			common.DomainErrorTypeInvalidValue,
			"説明が長すぎます（最大500文字）",
		)
	}

	return nil
}

// String 文字列表現を返す
func (d Description) String() string {
	return d.value
}

// IsEmpty 空文字かどうかを判定
func (d Description) IsEmpty() bool {
	return d.value == ""
}

// Equals 等価性をチェック
func (d Description) Equals(other Description) bool {
	return d.value == other.value
}
