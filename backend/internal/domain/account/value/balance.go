package value

import (
	"financetracker/internal/domain/common"
	"financetracker/internal/domain/common/value"
)

// Balance 口座の残高情報を表現する値オブジェクト
type Balance struct {
	initialBalance value.Money // 初期残高
	currentBalance value.Money // 現在残高
}

// NewBalance 新しいBalanceインスタンスを作成
func NewBalance(initialBalance, currentBalance value.Money) (*Balance, error) {
	// 通貨の整合性チェック
	if initialBalance.Currency() != currentBalance.Currency() {
		return nil, common.NewValidationError("balance", "currency mismatch",
			"initial balance and current balance must have the same currency")
	}

	return &Balance{
		initialBalance: initialBalance,
		currentBalance: currentBalance,
	}, nil
}

// NewBalanceWithInitial 初期残高のみを指定してBalanceインスタンスを作成（現在残高は初期残高と同じ）
func NewBalanceWithInitial(initialBalance value.Money) *Balance {
	return &Balance{
		initialBalance: initialBalance,
		currentBalance: initialBalance,
	}
}

// NewZeroBalance ゼロ残高のBalanceインスタンスを作成
func NewZeroBalance() *Balance {
	zeroMoney := value.NewMoneyJPY(0)
	return &Balance{
		initialBalance: *zeroMoney,
		currentBalance: *zeroMoney,
	}
}

// InitialBalance 初期残高を取得
func (b Balance) InitialBalance() value.Money {
	return b.initialBalance
}

// CurrentBalance 現在残高を取得
func (b Balance) CurrentBalance() value.Money {
	return b.currentBalance
}

// Add 入金（現在残高に加算）
func (b Balance) Add(amount value.Money) (*Balance, error) {
	newBalance, err := b.currentBalance.Add(amount)
	if err != nil {
		return nil, err
	}

	return &Balance{
		initialBalance: b.initialBalance,
		currentBalance: *newBalance,
	}, nil
}

// Subtract 出金（現在残高から減算）
func (b Balance) Subtract(amount value.Money) (*Balance, error) {
	newBalance, err := b.currentBalance.Subtract(amount)
	if err != nil {
		return nil, err
	}

	return &Balance{
		initialBalance: b.initialBalance,
		currentBalance: *newBalance,
	}, nil
}

// SetCurrentBalance 現在残高を直接設定
func (b Balance) SetCurrentBalance(amount value.Money) (*Balance, error) {
	// 通貨の整合性チェック
	if b.initialBalance.Currency() != amount.Currency() {
		return nil, common.NewValidationError("balance", "currency mismatch",
			"new balance must have the same currency as initial balance")
	}

	return &Balance{
		initialBalance: b.initialBalance,
		currentBalance: amount,
	}, nil
}

// GetDifference 初期残高との差額を取得
func (b Balance) GetDifference() (*value.Money, error) {
	return b.currentBalance.Subtract(b.initialBalance)
}

// IsPositive プラス残高かどうかを判定
func (b Balance) IsPositive() bool {
	return b.currentBalance.IsPositive()
}

// IsNegative マイナス残高かどうかを判定
func (b Balance) IsNegative() bool {
	return b.currentBalance.IsNegative()
}

// IsZero ゼロ残高かどうかを判定
func (b Balance) IsZero() bool {
	return b.currentBalance.IsZero()
}

// CanWithdraw 指定金額の出金が可能かどうかを判定
func (b Balance) CanWithdraw(amount value.Money, allowNegative bool) (bool, error) {
	// 通貨の整合性チェック
	if b.currentBalance.Currency() != amount.Currency() {
		return false, common.NewValidationError("amount", amount.Currency(),
			"withdrawal amount must have the same currency as balance")
	}

	if allowNegative {
		return true, nil // マイナス残高を許可する場合は常に出金可能
	}

	// 出金後の残高を計算
	afterWithdrawal, err := b.currentBalance.Subtract(amount)
	if err != nil {
		return false, err
	}

	// 出金後の残高がゼロ以上かチェック
	return !afterWithdrawal.IsNegative(), nil
}

// HasGained 初期残高より増加しているかどうかを判定
func (b Balance) HasGained() bool {
	difference, err := b.GetDifference()
	if err != nil {
		return false
	}
	return difference.IsPositive()
}

// HasLost 初期残高より減少しているかどうかを判定
func (b Balance) HasLost() bool {
	difference, err := b.GetDifference()
	if err != nil {
		return false
	}
	return difference.IsNegative()
}

// GetGainLossPercentage 初期残高に対する増減率を取得（パーセント）
func (b Balance) GetGainLossPercentage() float64 {
	if b.initialBalance.IsZero() {
		return 0.0 // 初期残高がゼロの場合は計算不可
	}

	difference, err := b.GetDifference()
	if err != nil {
		return 0.0
	}

	return (float64(difference.Amount()) / float64(b.initialBalance.Amount())) * 100.0
}

// String 文字列表現
func (b Balance) String() string {
	return b.currentBalance.String()
}

// Equals 残高の同一性判定
func (b Balance) Equals(other Balance) bool {
	return b.initialBalance.Equals(other.initialBalance) &&
		b.currentBalance.Equals(other.currentBalance)
}
