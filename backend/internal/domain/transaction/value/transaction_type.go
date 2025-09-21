// Package value トランザクション関連の値オブジェクトを定義
package value

import (
	"fmt"

	"financetracker/internal/domain/common"
)

// TransactionType トランザクションタイプを表す値オブジェクト
type TransactionType struct {
	value string
}

// 定数定義
const (
	TransactionTypeIncome  = "income"  // 収入
	TransactionTypeExpense = "expense" // 支出
)

// NewTransactionType 新しいTransactionTypeを作成
func NewTransactionType(value string) (TransactionType, error) {
	if err := validateTransactionType(value); err != nil {
		return TransactionType{}, err
	}
	return TransactionType{value: value}, nil
}

// validateTransactionType トランザクションタイプの妥当性を検証
func validateTransactionType(value string) error {
	switch value {
	case TransactionTypeIncome, TransactionTypeExpense:
		return nil
	default:
		return common.NewDomainError(
			common.DomainErrorTypeInvalidValue,
			fmt.Sprintf("無効なトランザクションタイプです: %s", value),
		)
	}
}

// String 文字列表現を返す
func (t TransactionType) String() string {
	return t.value
}

// IsIncome 収入かどうかを判定
func (t TransactionType) IsIncome() bool {
	return t.value == TransactionTypeIncome
}

// IsExpense 支出かどうかを判定
func (t TransactionType) IsExpense() bool {
	return t.value == TransactionTypeExpense
}

// Equals 等価性をチェック
func (t TransactionType) Equals(other TransactionType) bool {
	return t.value == other.value
}
