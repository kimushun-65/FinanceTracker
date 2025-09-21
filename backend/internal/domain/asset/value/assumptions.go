// Package value 資産予測の前提条件を定義
package value

import (
	"encoding/json"

	"financetracker/internal/domain/common"
	commonValue "financetracker/internal/domain/common/value"
)

// MajorExpense 大型支出予定を表す
type MajorExpense struct {
	Description string
	Amount      commonValue.Money
	PlannedDate string // YYYY-MM形式
}

// Assumptions 予測の前提条件を表す値オブジェクト
type Assumptions struct {
	monthlyIncome    commonValue.Money
	monthlyExpense   commonValue.Money
	savingsRate      float64
	investmentReturn float64
	inflationRate    float64
	majorExpenses    []MajorExpense
}

// NewAssumptions 新しいAssumptionsを作成
func NewAssumptions(
	monthlyIncome commonValue.Money,
	monthlyExpense commonValue.Money,
	savingsRate float64,
	investmentReturn float64,
	inflationRate float64,
	majorExpenses []MajorExpense,
) (Assumptions, error) {
	if err := validateAssumptions(savingsRate, investmentReturn, inflationRate); err != nil {
		return Assumptions{}, err
	}

	// majorExpensesのディープコピー
	expensesCopy := make([]MajorExpense, len(majorExpenses))
	copy(expensesCopy, majorExpenses)

	return Assumptions{
		monthlyIncome:    monthlyIncome,
		monthlyExpense:   monthlyExpense,
		savingsRate:      savingsRate,
		investmentReturn: investmentReturn,
		inflationRate:    inflationRate,
		majorExpenses:    expensesCopy,
	}, nil
}

// validateAssumptions 前提条件の妥当性を検証
func validateAssumptions(savingsRate, investmentReturn, inflationRate float64) error {
	// 各率は-1から1の範囲
	if savingsRate < -1 || savingsRate > 1 {
		return common.NewDomainError(
			common.DomainErrorTypeInvalidValue,
			"貯蓄率は-1から1の範囲である必要があります",
		)
	}

	if investmentReturn < -1 || investmentReturn > 1 {
		return common.NewDomainError(
			common.DomainErrorTypeInvalidValue,
			"投資リターン率は-1から1の範囲である必要があります",
		)
	}

	if inflationRate < -1 || inflationRate > 1 {
		return common.NewDomainError(
			common.DomainErrorTypeInvalidValue,
			"インフレ率は-1から1の範囲である必要があります",
		)
	}

	return nil
}

// MonthlyIncome 月収予測を取得
func (a Assumptions) MonthlyIncome() commonValue.Money {
	return a.monthlyIncome
}

// MonthlyExpense 月支出予測を取得
func (a Assumptions) MonthlyExpense() commonValue.Money {
	return a.monthlyExpense
}

// SavingsRate 貯蓄率を取得
func (a Assumptions) SavingsRate() float64 {
	return a.savingsRate
}

// InvestmentReturn 投資リターン率を取得
func (a Assumptions) InvestmentReturn() float64 {
	return a.investmentReturn
}

// InflationRate インフレ率を取得
func (a Assumptions) InflationRate() float64 {
	return a.inflationRate
}

// MajorExpenses 大型支出予定を取得
func (a Assumptions) MajorExpenses() []MajorExpense {
	// ディープコピーして返す
	result := make([]MajorExpense, len(a.majorExpenses))
	copy(result, a.majorExpenses)
	return result
}

// CalculateMonthlySavings 月間貯蓄額を計算
func (a Assumptions) CalculateMonthlySavings() (*commonValue.Money, error) {
	// 収入 - 支出
	savings, err := a.monthlyIncome.Subtract(a.monthlyExpense)
	if err != nil {
		return nil, err
	}

	// 貯蓄率を適用
	savingsAmount := int64(float64(savings.Amount()) * a.savingsRate)

	result, err := commonValue.NewMoney(savingsAmount, savings.Currency())
	if err != nil {
		return nil, err
	}

	return result, nil
}

// ToJSON JSON形式で出力
func (a Assumptions) ToJSON() (string, error) {
	data := map[string]interface{}{
		"monthly_income": map[string]interface{}{
			"amount":   a.monthlyIncome.Amount(),
			"currency": a.monthlyIncome.Currency(),
		},
		"monthly_expense": map[string]interface{}{
			"amount":   a.monthlyExpense.Amount(),
			"currency": a.monthlyExpense.Currency(),
		},
		"savings_rate":      a.savingsRate,
		"investment_return": a.investmentReturn,
		"inflation_rate":    a.inflationRate,
		"major_expenses":    []map[string]interface{}{},
	}

	// 大型支出予定を追加
	for _, expense := range a.majorExpenses {
		data["major_expenses"] = append(
			data["major_expenses"].([]map[string]interface{}),
			map[string]interface{}{
				"description":  expense.Description,
				"amount":       expense.Amount.Amount(),
				"currency":     expense.Amount.Currency(),
				"planned_date": expense.PlannedDate,
			},
		)
	}

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}

// IsValid 前提条件の妥当性をチェック
func (a Assumptions) IsValid() bool {
	// 月収・支出は0以上
	if a.monthlyIncome.IsNegative() || a.monthlyExpense.IsNegative() {
		return false
	}

	// 各率のチェック
	if a.savingsRate < -1 || a.savingsRate > 1 ||
		a.investmentReturn < -1 || a.investmentReturn > 1 ||
		a.inflationRate < -1 || a.inflationRate > 1 {
		return false
	}

	return true
}
