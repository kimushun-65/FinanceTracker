package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	transactionDomain "financetracker/internal/domain/transaction/entity"
)

// TransactionResponse 取引情報レスポンス
type TransactionResponse struct {
	ID              uuid.UUID       `json:"id"`
	UserID          uuid.UUID       `json:"user_id"`
	AccountID       uuid.UUID       `json:"account_id"`
	CategoryID      uuid.UUID       `json:"category_id"`
	TransactionType string          `json:"transaction_type"`
	Amount          decimal.Decimal `json:"amount"`
	Description     string          `json:"description"`
	TransactionDate time.Time       `json:"transaction_date"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}

// CreateTransactionRequest 取引作成リクエスト
type CreateTransactionRequest struct {
	AccountID       uuid.UUID       `json:"account_id" binding:"required"`
	CategoryID      uuid.UUID       `json:"category_id" binding:"required"`
	TransactionType string          `json:"transaction_type" binding:"required,oneof=income expense"`
	Amount          decimal.Decimal `json:"amount" binding:"required"`
	Description     string          `json:"description" binding:"required,min=1,max=500"`
	TransactionDate *time.Time      `json:"transaction_date" binding:"omitempty"`
}

// UpdateTransactionRequest 取引更新リクエスト
type UpdateTransactionRequest struct {
	CategoryID      *uuid.UUID       `json:"category_id" binding:"omitempty"`
	Amount          *decimal.Decimal `json:"amount" binding:"omitempty"`
	Description     *string          `json:"description" binding:"omitempty,min=1,max=500"`
	TransactionDate *time.Time       `json:"transaction_date" binding:"omitempty"`
}

// TransactionListResponse 取引一覧レスポンス
type TransactionListResponse struct {
	Transactions []TransactionResponse `json:"transactions"`
	TotalCount   int64                 `json:"total_count"`
	TotalIncome  decimal.Decimal       `json:"total_income"`
	TotalExpense decimal.Decimal       `json:"total_expense"`
	Page         int                   `json:"page"`
	PerPage      int                   `json:"per_page"`
}

// TransactionSearchParams 取引検索パラメータ
type TransactionSearchParams struct {
	UserID          uuid.UUID        `form:"-"`
	AccountID       *uuid.UUID       `form:"account_id"`
	CategoryID      *uuid.UUID       `form:"category_id"`
	TransactionType *string          `form:"transaction_type" binding:"omitempty,oneof=income expense"`
	DateFrom        *time.Time       `form:"date_from"`
	DateTo          *time.Time       `form:"date_to"`
	AmountMin       *decimal.Decimal `form:"amount_min" binding:"omitempty"`
	AmountMax       *decimal.Decimal `form:"amount_max" binding:"omitempty"`
	Description     *string          `form:"description"`
	Page            int              `form:"page,default=1" binding:"min=1"`
	PerPage         int              `form:"per_page,default=20" binding:"min=1,max=100"`
	OrderBy         string           `form:"order_by,default=transaction_date desc"`
}

// TransactionSummary 取引サマリー
type TransactionSummary struct {
	Period       string                       `json:"period"`
	TotalIncome  decimal.Decimal              `json:"total_income"`
	TotalExpense decimal.Decimal              `json:"total_expense"`
	NetAmount    decimal.Decimal              `json:"net_amount"`
	ByCategory   []CategoryTransactionSummary `json:"by_category"`
}

// CategoryTransactionSummary カテゴリー別取引サマリー
type CategoryTransactionSummary struct {
	CategoryID   uuid.UUID       `json:"category_id"`
	CategoryName string          `json:"category_name"`
	TotalAmount  decimal.Decimal `json:"total_amount"`
	Count        int             `json:"count"`
	Percentage   float64         `json:"percentage"`
}

// TransactionFromDomain ドメインエンティティからDTOへの変換
func TransactionFromDomain(transaction *transactionDomain.Transaction) *TransactionResponse {
	if transaction == nil {
		return nil
	}

	amount := transaction.Amount()
	return &TransactionResponse{
		ID:              transaction.ID,
		UserID:          transaction.UserID(),
		AccountID:       transaction.AccountID(),
		CategoryID:      transaction.CategoryID(),
		TransactionType: transaction.Type().String(),
		Amount:          decimal.NewFromInt(amount.Amount()),
		Description:     transaction.Description().String(),
		TransactionDate: transaction.Date(),
		CreatedAt:       transaction.CreatedAt,
		UpdatedAt:       transaction.UpdatedAt,
	}
}

// TransactionsFromDomain ドメインエンティティのスライスからDTOへの変換
func TransactionsFromDomain(transactions []*transactionDomain.Transaction) []TransactionResponse {
	result := make([]TransactionResponse, len(transactions))
	for i, transaction := range transactions {
		if transactionDTO := TransactionFromDomain(transaction); transactionDTO != nil {
			result[i] = *transactionDTO
		}
	}
	return result
}

// MonthlyTransactionSummary 月次取引サマリー
type MonthlyTransactionSummary struct {
	Year         int                    `json:"year"`
	Month        int                    `json:"month"`
	TotalIncome  decimal.Decimal        `json:"total_income"`
	TotalExpense decimal.Decimal        `json:"total_expense"`
	NetAmount    decimal.Decimal        `json:"net_amount"`
	DailyData    []DailyTransactionData `json:"daily_data"`
}

// DailyTransactionData 日別取引データ
type DailyTransactionData struct {
	Date         time.Time       `json:"date"`
	TotalIncome  decimal.Decimal `json:"total_income"`
	TotalExpense decimal.Decimal `json:"total_expense"`
	NetAmount    decimal.Decimal `json:"net_amount"`
	Count        int             `json:"count"`
}

// CategorySummaryResponse カテゴリー別サマリーレスポンス
type CategorySummaryResponse struct {
	Period      PeriodInfo              `json:"period"`
	TotalAmount decimal.Decimal         `json:"total_amount"`
	ByCategory  []CategorySummaryDetail `json:"by_category"`
}

// PeriodInfo 期間情報
type PeriodInfo struct {
	From time.Time `json:"from"`
	To   time.Time `json:"to"`
}

// CategorySummaryDetail カテゴリー別詳細
type CategorySummaryDetail struct {
	CategoryID       uuid.UUID       `json:"category_id"`
	CategoryName     string          `json:"category_name"`
	CategoryIcon     string          `json:"category_icon"`
	TotalAmount      decimal.Decimal `json:"total_amount"`
	TransactionCount int             `json:"transaction_count"`
	Percentage       decimal.Decimal `json:"percentage"`
}
