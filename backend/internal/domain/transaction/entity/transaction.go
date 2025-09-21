// Package entity トランザクション関連のエンティティを定義
package entity

import (
	"time"

	"financetracker/internal/domain/common"
	"financetracker/internal/domain/common/value"
	transactionValue "financetracker/internal/domain/transaction/value"

	"github.com/google/uuid"
)

// Transaction トランザクションエンティティ
type Transaction struct {
	common.BaseEntity
	userID      uuid.UUID
	accountID   uuid.UUID
	categoryID  uuid.UUID
	amount      value.Money
	txType      transactionValue.TransactionType
	date        time.Time
	description transactionValue.Description
}

// NewTransaction 新しいトランザクションを作成
func NewTransaction(
	userID uuid.UUID,
	accountID uuid.UUID,
	categoryID uuid.UUID,
	amount value.Money,
	txType transactionValue.TransactionType,
	date time.Time,
	description transactionValue.Description,
) (Transaction, error) {
	if err := validateTransactionParams(userID, accountID, categoryID, amount, date); err != nil {
		return Transaction{}, err
	}

	baseEntity := common.NewBaseEntity()

	return Transaction{
		BaseEntity:  baseEntity,
		userID:      userID,
		accountID:   accountID,
		categoryID:  categoryID,
		amount:      amount,
		txType:      txType,
		date:        date,
		description: description,
	}, nil
}

// validateTransactionParams トランザクションパラメータの妥当性を検証
func validateTransactionParams(
	userID uuid.UUID,
	accountID uuid.UUID,
	categoryID uuid.UUID,
	amount value.Money,
	date time.Time,
) error {
	if userID == uuid.Nil {
		return common.NewDomainError(
			common.DomainErrorTypeInvalidValue,
			"ユーザーIDが必要です",
		)
	}

	if accountID == uuid.Nil {
		return common.NewDomainError(
			common.DomainErrorTypeInvalidValue,
			"アカウントIDが必要です",
		)
	}

	if categoryID == uuid.Nil {
		return common.NewDomainError(
			common.DomainErrorTypeInvalidValue,
			"カテゴリIDが必要です",
		)
	}

	if amount.IsZero() {
		return common.NewDomainError(
			common.DomainErrorTypeInvalidValue,
			"金額は0より大きい必要があります",
		)
	}

	if date.IsZero() {
		return common.NewDomainError(
			common.DomainErrorTypeInvalidValue,
			"取引日が必要です",
		)
	}

	return nil
}

// UserID ユーザーIDを取得
func (t Transaction) UserID() uuid.UUID {
	return t.userID
}

// AccountID アカウントIDを取得
func (t Transaction) AccountID() uuid.UUID {
	return t.accountID
}

// CategoryID カテゴリIDを取得
func (t Transaction) CategoryID() uuid.UUID {
	return t.categoryID
}

// Amount 金額を取得
func (t Transaction) Amount() value.Money {
	return t.amount
}

// Type トランザクションタイプを取得
func (t Transaction) Type() transactionValue.TransactionType {
	return t.txType
}

// Date 取引日を取得
func (t Transaction) Date() time.Time {
	return t.date
}

// Description 説明を取得
func (t Transaction) Description() transactionValue.Description {
	return t.description
}

// UpdateAmount 金額を更新
func (t *Transaction) UpdateAmount(amount value.Money) error {
	if amount.IsZero() {
		return common.NewDomainError(
			common.DomainErrorTypeInvalidValue,
			"金額は0より大きい必要があります",
		)
	}

	t.amount = amount
	t.UpdateTimestamp()
	return nil
}

// UpdateDescription 説明を更新
func (t *Transaction) UpdateDescription(description transactionValue.Description) {
	t.description = description
	t.UpdateTimestamp()
}

// UpdateDate 取引日を更新
func (t *Transaction) UpdateDate(date time.Time) error {
	if date.IsZero() {
		return common.NewDomainError(
			common.DomainErrorTypeInvalidValue,
			"取引日が必要です",
		)
	}

	t.date = date
	t.UpdateTimestamp()
	return nil
}

// IsIncome 収入かどうかを判定
func (t Transaction) IsIncome() bool {
	return t.txType.IsIncome()
}

// IsExpense 支出かどうかを判定
func (t Transaction) IsExpense() bool {
	return t.txType.IsExpense()
}
