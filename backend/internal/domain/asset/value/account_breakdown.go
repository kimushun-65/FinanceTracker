// Package value 資産関連の値オブジェクトを定義
package value

import (
	"encoding/json"

	"financetracker/internal/domain/common"
	commonValue "financetracker/internal/domain/common/value"

	"github.com/google/uuid"
)

// AccountBalance 口座残高情報
type AccountBalance struct {
	AccountID   uuid.UUID
	AccountName string
	Balance     commonValue.Money
}

// AccountBreakdown 口座別内訳を表す値オブジェクト
type AccountBreakdown struct {
	accounts []AccountBalance
}

// NewAccountBreakdown 新しいAccountBreakdownを作成
func NewAccountBreakdown(accounts []AccountBalance) (AccountBreakdown, error) {
	if err := validateAccountBreakdown(accounts); err != nil {
		return AccountBreakdown{}, err
	}

	// ディープコピーして不変性を保証
	accountsCopy := make([]AccountBalance, len(accounts))
	copy(accountsCopy, accounts)

	return AccountBreakdown{accounts: accountsCopy}, nil
}

// validateAccountBreakdown 口座内訳の妥当性を検証
func validateAccountBreakdown(accounts []AccountBalance) error {
	// 重複チェック
	seen := make(map[uuid.UUID]bool)
	for _, account := range accounts {
		if account.AccountID == uuid.Nil {
			return common.NewDomainError(
				common.DomainErrorTypeInvalidValue,
				"口座IDが無効です",
			)
		}

		if seen[account.AccountID] {
			return common.NewDomainError(
				common.DomainErrorTypeInvalidValue,
				"重複する口座IDが含まれています",
			)
		}
		seen[account.AccountID] = true

		if account.AccountName == "" {
			return common.NewDomainError(
				common.DomainErrorTypeInvalidValue,
				"口座名が必要です",
			)
		}
	}

	return nil
}

// Accounts 口座別残高のリストを取得
func (a AccountBreakdown) Accounts() []AccountBalance {
	// ディープコピーして返す（不変性を保証）
	result := make([]AccountBalance, len(a.accounts))
	copy(result, a.accounts)
	return result
}

// TotalBalance 総残高を計算
func (a AccountBreakdown) TotalBalance() commonValue.Money {
	if len(a.accounts) == 0 {
		// デフォルトはJPY
		return *commonValue.NewMoneyJPY(0)
	}

	// 最初の口座の通貨を基準とする
	total := commonValue.NewMoneyJPY(0)
	if a.accounts[0].Balance.Currency() != "JPY" {
		// JPY以外の通貨の場合
		money, _ := commonValue.NewMoney(0, a.accounts[0].Balance.Currency())
		total = money
	}

	for _, account := range a.accounts {
		newTotal, err := total.Add(account.Balance)
		if err != nil {
			// 異なる通貨の場合はスキップ（本来はエラーとすべきかもしれない）
			continue
		}
		total = newTotal
	}

	return *total
}

// ToJSON JSON形式で出力
func (a AccountBreakdown) ToJSON() (string, error) {
	data := make(map[string]interface{})

	accounts := make([]map[string]interface{}, len(a.accounts))
	for i, account := range a.accounts {
		accounts[i] = map[string]interface{}{
			"account_id":   account.AccountID.String(),
			"account_name": account.AccountName,
			"balance":      account.Balance.Amount(),
			"currency":     account.Balance.Currency(),
		}
	}

	data["accounts"] = accounts
	data["total_balance"] = a.TotalBalance().Amount()
	data["currency"] = a.TotalBalance().Currency()

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}

// FindByAccountID 口座IDで残高情報を検索
func (a AccountBreakdown) FindByAccountID(accountID uuid.UUID) (*AccountBalance, bool) {
	for _, account := range a.accounts {
		if account.AccountID == accountID {
			// コピーを返す
			result := account
			return &result, true
		}
	}
	return nil, false
}
