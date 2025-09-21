// Package value 予算関連の値オブジェクトを定義
package value

import (
	"fmt"

	"financetracker/internal/domain/common"
)

// PeriodType 予算期間タイプを表す値オブジェクト
type PeriodType struct {
	value string
}

// 定数定義
const (
	PeriodTypeMonthly = "monthly" // 月次
	PeriodTypeYearly  = "yearly"  // 年次
)

// NewPeriodType 新しいPeriodTypeを作成
func NewPeriodType(value string) (PeriodType, error) {
	if err := validatePeriodType(value); err != nil {
		return PeriodType{}, err
	}
	return PeriodType{value: value}, nil
}

// validatePeriodType 期間タイプの妥当性を検証
func validatePeriodType(value string) error {
	switch value {
	case PeriodTypeMonthly, PeriodTypeYearly:
		return nil
	default:
		return common.NewDomainError(
			common.DomainErrorTypeInvalidValue,
			fmt.Sprintf("無効な期間タイプです: %s", value),
		)
	}
}

// String 文字列表現を返す
func (p PeriodType) String() string {
	return p.value
}

// IsMonthly 月次かどうかを判定
func (p PeriodType) IsMonthly() bool {
	return p.value == PeriodTypeMonthly
}

// IsYearly 年次かどうかを判定
func (p PeriodType) IsYearly() bool {
	return p.value == PeriodTypeYearly
}

// GetDisplayName 表示用名称を取得
func (p PeriodType) GetDisplayName() string {
	switch p.value {
	case PeriodTypeMonthly:
		return "月次"
	case PeriodTypeYearly:
		return "年次"
	default:
		return ""
	}
}

// Equals 等価性をチェック
func (p PeriodType) Equals(other PeriodType) bool {
	return p.value == other.value
}
