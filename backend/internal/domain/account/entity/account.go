package entity

import (
	accountValue "financetracker/internal/domain/account/value"
	"financetracker/internal/domain/common"
	"financetracker/internal/domain/common/value"
	userValue "financetracker/internal/domain/user/value"
)

// Account 口座エンティティ
type Account struct {
	common.BaseEntity
	userID  userValue.UserID
	name    accountValue.AccountName
	accType accountValue.AccountType
	balance accountValue.Balance
}

// NewAccount 新しいAccountエンティティを作成
func NewAccount(
	userID userValue.UserID,
	name accountValue.AccountName,
	accountType accountValue.AccountType,
	initialBalance value.Money,
) *Account {
	balance := accountValue.NewBalanceWithInitial(initialBalance)

	return &Account{
		BaseEntity: common.NewBaseEntity(),
		userID:     userID,
		name:       name,
		accType:    accountType,
		balance:    *balance,
	}
}

// ReconstructAccount 既存のデータからAccountエンティティを再構築（リポジトリから取得時に使用）
func ReconstructAccount(
	baseEntity common.BaseEntity,
	userID userValue.UserID,
	name accountValue.AccountName,
	accountType accountValue.AccountType,
	balance accountValue.Balance,
) *Account {
	return &Account{
		BaseEntity: baseEntity,
		userID:     userID,
		name:       name,
		accType:    accountType,
		balance:    balance,
	}
}

// UserID ユーザーIDを取得
func (a Account) UserID() userValue.UserID {
	return a.userID
}

// Name 口座名を取得
func (a Account) Name() accountValue.AccountName {
	return a.name
}

// Type 口座タイプを取得
func (a Account) Type() accountValue.AccountType {
	return a.accType
}

// Balance 残高情報を取得
func (a Account) Balance() accountValue.Balance {
	return a.balance
}

// CurrentBalance 現在残高を取得
func (a Account) CurrentBalance() value.Money {
	return a.balance.CurrentBalance()
}

// InitialBalance 初期残高を取得
func (a Account) InitialBalance() value.Money {
	return a.balance.InitialBalance()
}

// UpdateBalance 残高を更新（入金・出金）
func (a *Account) UpdateBalance(amount value.Money, isDeposit bool) error {
	var newBalance *accountValue.Balance
	var err error

	// クレジットカード口座の場合はマイナス残高を許可
	allowNegative := a.accType.IsCreditCard()

	if isDeposit {
		newBalance, err = a.balance.Add(amount)
	} else {
		// 出金の場合、残高不足をチェック（クレジットカード以外はマイナス不可）
		if !allowNegative {
			canWithdraw, checkErr := a.balance.CanWithdraw(amount, false)
			if checkErr != nil {
				return checkErr
			}

			if !canWithdraw {
				return common.NewBusinessRuleError("insufficient_balance",
					"insufficient balance for withdrawal")
			}
		}

		newBalance, err = a.balance.Subtract(amount)
	}

	if err != nil {
		return err
	}

	// クレジットカードの場合、結果の残高が正にならないよう確認
	if a.accType.IsCreditCard() && newBalance.CurrentBalance().IsPositive() {
		return common.NewBusinessRuleError("invalid_credit_card_balance",
			"credit card balance cannot be positive")
	}

	a.balance = *newBalance
	a.UpdateTimestamp()
	return nil
}

// SetBalance 残高を直接設定
func (a *Account) SetBalance(newBalance value.Money) error {
	// クレジットカードの場合は正の残高を許可しない
	if a.accType.IsCreditCard() && newBalance.IsPositive() {
		return common.NewBusinessRuleError("invalid_credit_card_balance",
			"credit card balance cannot be positive")
	}

	updatedBalance, err := a.balance.SetCurrentBalance(newBalance)
	if err != nil {
		return err
	}

	a.balance = *updatedBalance
	a.UpdateTimestamp()
	return nil
}

// UpdateName 口座名を更新
func (a *Account) UpdateName(newName accountValue.AccountName) {
	a.name = newName
	a.UpdateTimestamp()
}

// UpdateType 口座タイプを更新
func (a *Account) UpdateType(newType accountValue.AccountType) {
	a.accType = newType
	a.UpdateTimestamp()
}

// CanWithdraw 指定金額の出金が可能かどうかを判定
func (a Account) CanWithdraw(amount value.Money) (bool, error) {
	return a.balance.CanWithdraw(amount, false)
}

// GetDisplayInfo 表示用の口座情報を取得
func (a Account) GetDisplayInfo() AccountDisplayInfo {
	return AccountDisplayInfo{
		Name:            a.name.Value(),
		TypeDisplayName: a.accType.GetDisplayName(),
		Balance:         a.balance.CurrentBalance().Format(),
	}
}

// GetBalanceStatus 残高の状態を取得
func (a Account) GetBalanceStatus() BalanceStatus {
	status := BalanceStatusNormal

	if a.balance.IsNegative() {
		status = BalanceStatusNegative
	} else if a.balance.IsZero() {
		status = BalanceStatusZero
	}

	return BalanceStatus{
		Status:             status,
		HasGained:          a.balance.HasGained(),
		HasLost:            a.balance.HasLost(),
		GainLossPercentage: a.balance.GetGainLossPercentage(),
	}
}

// AccountDisplayInfo 表示用の口座情報
type AccountDisplayInfo struct {
	Name            string `json:"name"`
	TypeDisplayName string `json:"type_display_name"`
	Balance         string `json:"balance"`
}

// BalanceStatusType 残高状態のタイプ
type BalanceStatusType string

const (
	BalanceStatusNormal   BalanceStatusType = "normal"   // 通常（プラス残高）
	BalanceStatusZero     BalanceStatusType = "zero"     // ゼロ残高
	BalanceStatusNegative BalanceStatusType = "negative" // マイナス残高
)

// BalanceStatus 残高の状態情報
type BalanceStatus struct {
	Status             BalanceStatusType `json:"status"`
	HasGained          bool              `json:"has_gained"`
	HasLost            bool              `json:"has_lost"`
	GainLossPercentage float64           `json:"gain_loss_percentage"`
}
