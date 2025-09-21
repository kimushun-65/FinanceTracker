package value

import (
	"strings"

	"financetracker/internal/domain/common"
)

// AccountName 口座名を表現する値オブジェクト
type AccountName struct {
	value string
}

// NewAccountName 新しいAccountNameインスタンスを作成
func NewAccountName(value string) (*AccountName, error) {
	name := strings.TrimSpace(value)
	
	if err := validateAccountName(name); err != nil {
		return nil, err
	}

	return &AccountName{
		value: name,
	}, nil
}

// Value 口座名の値を取得
func (a AccountName) Value() string {
	return a.value
}

// String 文字列表現
func (a AccountName) String() string {
	return a.value
}

// Equals 口座名の同一性判定
func (a AccountName) Equals(other AccountName) bool {
	return a.value == other.value
}

// IsEmpty 空の口座名かどうかを判定
func (a AccountName) IsEmpty() bool {
	return a.value == ""
}

// Length 口座名の文字数を取得
func (a AccountName) Length() int {
	return len([]rune(a.value)) // マルチバイト文字を考慮
}

// validateAccountName 口座名のバリデーション
func validateAccountName(name string) error {
	// 空文字チェック
	if name == "" {
		return common.NewValidationError("account_name", name, "account name is required")
	}

	// 長さチェック（日本語を考慮して50文字）
	if len([]rune(name)) > 50 {
		return common.NewValidationError("account_name", name, "account name must be 50 characters or less")
	}

	// 制御文字チェック
	for _, char := range name {
		if char < 32 || char == 127 { // ASCII制御文字
			return common.NewValidationError("account_name", name, "account name contains invalid characters")
		}
	}

	return nil
}