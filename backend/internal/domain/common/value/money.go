package value

import (
	"fmt"
	"strconv"
	"strings"

	"financetracker/internal/domain/common"
)

// Money 金額を表現する値オブジェクト
// 整数で管理し、小数点以下は扱わない（円単位）
type Money struct {
	amount   int64
	currency string
}

// NewMoney 新しいMoneyインスタンスを作成
func NewMoney(amount int64, currency string) (*Money, error) {
	if currency == "" {
		currency = "JPY" // デフォルト通貨
	}

	if !isValidCurrency(currency) {
		return nil, common.NewValidationError("currency", currency, "invalid currency code")
	}

	return &Money{
		amount:   amount,
		currency: currency,
	}, nil
}

// NewMoneyJPY JPY通貨でMoneyインスタンスを作成
func NewMoneyJPY(amount int64) *Money {
	money, _ := NewMoney(amount, "JPY") // JPYは常に有効なので、エラーは発生しない
	return money
}

// Amount 金額を取得
func (m Money) Amount() int64 {
	return m.amount
}

// Currency 通貨を取得
func (m Money) Currency() string {
	return m.currency
}

// Add 加算 - 同じ通貨のみ計算可能
func (m Money) Add(other Money) (*Money, error) {
	if m.currency != other.currency {
		return nil, common.NewValidationError("currency", other.currency,
			fmt.Sprintf("cannot add different currencies: %s and %s", m.currency, other.currency))
	}

	return &Money{
		amount:   m.amount + other.amount,
		currency: m.currency,
	}, nil
}

// Subtract 減算 - 同じ通貨のみ計算可能
func (m Money) Subtract(other Money) (*Money, error) {
	if m.currency != other.currency {
		return nil, common.NewValidationError("currency", other.currency,
			fmt.Sprintf("cannot subtract different currencies: %s and %s", m.currency, other.currency))
	}

	return &Money{
		amount:   m.amount - other.amount,
		currency: m.currency,
	}, nil
}

// Multiply 乗算
func (m Money) Multiply(factor float64) *Money {
	return &Money{
		amount:   int64(float64(m.amount) * factor),
		currency: m.currency,
	}
}

// IsPositive 正の値かチェック
func (m Money) IsPositive() bool {
	return m.amount > 0
}

// IsZero ゼロかチェック
func (m Money) IsZero() bool {
	return m.amount == 0
}

// IsNegative 負の値かチェック
func (m Money) IsNegative() bool {
	return m.amount < 0
}

// Abs 絶対値を取得
func (m Money) Abs() *Money {
	amount := m.amount
	if amount < 0 {
		amount = -amount
	}
	return &Money{
		amount:   amount,
		currency: m.currency,
	}
}

// Format フォーマット済み文字列を取得（¥1,000形式）
func (m Money) Format() string {
	switch m.currency {
	case "JPY":
		return fmt.Sprintf("¥%s", formatNumber(m.amount))
	case "USD":
		return fmt.Sprintf("$%s", formatNumber(m.amount))
	case "EUR":
		return fmt.Sprintf("€%s", formatNumber(m.amount))
	default:
		return fmt.Sprintf("%s %s", m.currency, formatNumber(m.amount))
	}
}

// String 文字列表現
func (m Money) String() string {
	return m.Format()
}

// Equals 金額の同一性判定
func (m Money) Equals(other Money) bool {
	return m.amount == other.amount && m.currency == other.currency
}

// Compare 金額の比較 (-1: m < other, 0: m == other, 1: m > other)
func (m Money) Compare(other Money) (int, error) {
	if m.currency != other.currency {
		return 0, common.NewValidationError("currency", other.currency,
			fmt.Sprintf("cannot compare different currencies: %s and %s", m.currency, other.currency))
	}

	if m.amount < other.amount {
		return -1, nil
	} else if m.amount > other.amount {
		return 1, nil
	}
	return 0, nil
}

// formatNumber 数値を3桁区切りでフォーマット
func formatNumber(num int64) string {
	str := strconv.FormatInt(num, 10)
	if len(str) <= 3 {
		return str
	}

	var result strings.Builder
	for i, char := range str {
		if i > 0 && (len(str)-i)%3 == 0 {
			result.WriteString(",")
		}
		result.WriteRune(char)
	}
	return result.String()
}

// isValidCurrency 通貨コードの妥当性をチェック
func isValidCurrency(currency string) bool {
	validCurrencies := map[string]bool{
		"JPY": true,
		"USD": true,
		"EUR": true,
		"GBP": true,
		"CNY": true,
		"KRW": true,
	}
	return validCurrencies[currency]
}
