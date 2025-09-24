package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	budgetDomain "financetracker/internal/domain/budget/entity"
)

// BudgetResponse 予算情報レスポンス
type BudgetResponse struct {
	ID         uuid.UUID       `json:"id"`
	UserID     uuid.UUID       `json:"user_id"`
	CategoryID uuid.UUID       `json:"category_id"`
	Amount     decimal.Decimal `json:"amount"`
	Period     string          `json:"period"`
	StartDate  time.Time       `json:"start_date"`
	EndDate    *time.Time      `json:"end_date,omitempty"`
	IsActive   bool            `json:"is_active"`
	CreatedAt  time.Time       `json:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at"`
}

// CreateBudgetRequest 予算作成リクエスト
type CreateBudgetRequest struct {
	CategoryID uuid.UUID       `json:"category_id" binding:"required"`
	Amount     decimal.Decimal `json:"amount" binding:"required,gt=0"`
	Period     string          `json:"period" binding:"required,oneof=monthly yearly"`
	StartDate  time.Time       `json:"start_date" binding:"required"`
	EndDate    *time.Time      `json:"end_date" binding:"omitempty"`
}

// UpdateBudgetRequest 予算更新リクエスト
type UpdateBudgetRequest struct {
	Amount    *decimal.Decimal `json:"amount" binding:"omitempty,gt=0"`
	StartDate *time.Time       `json:"start_date" binding:"omitempty"`
	EndDate   *time.Time       `json:"end_date" binding:"omitempty"`
	IsActive  *bool            `json:"is_active" binding:"omitempty"`
}

// BudgetListResponse 予算一覧レスポンス
type BudgetListResponse struct {
	Budgets    []BudgetResponse `json:"budgets"`
	TotalCount int64            `json:"total_count"`
}

// BudgetSearchParams 予算検索パラメータ
type BudgetSearchParams struct {
	UserID     uuid.UUID  `form:"-"`
	CategoryID *uuid.UUID `form:"category_id"`
	Period     *string    `form:"period" binding:"omitempty,oneof=monthly yearly"`
	IsActive   *bool      `form:"is_active"`
	OrderBy    string     `form:"order_by,default=start_date desc"`
}

// BudgetFromDomain ドメインエンティティからDTOへの変換
func BudgetFromDomain(budget *budgetDomain.Budget) *BudgetResponse {
	if budget == nil {
		return nil
	}

	return &BudgetResponse{
		ID:         budget.ID,
		UserID:     budget.UserID(),
		CategoryID: budget.CategoryID(),
		Amount:     decimal.NewFromInt(budget.Amount().Amount()),
		Period:     budget.PeriodType().String(),
		StartDate:  budget.StartDate(),
		EndDate:    budget.EndDate(),
		IsActive:   budget.IsActive(),
		CreatedAt:  budget.CreatedAt,
		UpdatedAt:  budget.UpdatedAt,
	}
}

// BudgetsFromDomain ドメインエンティティのスライスからDTOへの変換
func BudgetsFromDomain(budgets []*budgetDomain.Budget) []BudgetResponse {
	result := make([]BudgetResponse, len(budgets))
	for i, budget := range budgets {
		if budgetDTO := BudgetFromDomain(budget); budgetDTO != nil {
			result[i] = *budgetDTO
		}
	}
	return result
}
