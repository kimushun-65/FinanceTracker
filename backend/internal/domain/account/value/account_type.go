package value

import (
	"strings"

	"financetracker/internal/domain/common"
)

// AccountType 口座種別を表現する値オブジェクト
type AccountType struct {
	value string
}

// 定義済みの口座タイプ
const (
	AccountTypeChecking   = "checking"    // 普通預金
	AccountTypeInvestment = "investment"  // 投資
	AccountTypeCash       = "cash"        // 現金
)

// NewAccountType 新しいAccountTypeインスタンスを作成
func NewAccountType(value string) (*AccountType, error) {
	accountType := strings.TrimSpace(strings.ToLower(value))
	
	if err := validateAccountType(accountType); err != nil {
		return nil, err
	}

	return &AccountType{
		value: accountType,
	}, nil
}

// Value 口座タイプの値を取得
func (a AccountType) Value() string {
	return a.value
}

// GetDisplayName 表示用名称を取得
func (a AccountType) GetDisplayName() string {
	switch a.value {
	case AccountTypeChecking:
		return "普通預金"
	case AccountTypeCash:
		return "現金"
	case AccountTypeInvestment:
		return "投資"
	default:
		return a.value
	}
}

// String 文字列表現
func (a AccountType) String() string {
	return a.value
}

// Equals 口座タイプの同一性判定
func (a AccountType) Equals(other AccountType) bool {
	return a.value == other.value
}

// validateAccountType 口座タイプのバリデーション
func validateAccountType(accountType string) error {
	// 空文字チェック
	if accountType == "" {
		return common.NewValidationError("account_type", accountType, "account type is required")
	}

	// 有効な口座タイプかチェック
	validTypes := map[string]bool{
		AccountTypeChecking:   true,
		AccountTypeInvestment: true,
		AccountTypeCash:       true,
	}

	if !validTypes[accountType] {
		return common.NewValidationError("account_type", accountType, 
			"invalid account type. Must be one of: checking, cash, investment")
	}

	return nil
}

// 事前定義された口座タイプのインスタンス
var (
	Checking   = &AccountType{value: AccountTypeChecking}
	Investment = &AccountType{value: AccountTypeInvestment}
	Cash       = &AccountType{value: AccountTypeCash}
)

// GetAllAccountTypes 利用可能な全ての口座タイプを取得
func GetAllAccountTypes() []*AccountType {
	return []*AccountType{
		Checking,
		Cash,
		Investment,
	}
}