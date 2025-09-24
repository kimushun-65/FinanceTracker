package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	accountDomain "financetracker/internal/domain/account/entity"
)

// AccountResponse 口座情報レスポンス
type AccountResponse struct {
	ID          uuid.UUID       `json:"id"`
	UserID      uuid.UUID       `json:"user_id"`
	Name        string          `json:"name"`
	AccountType string          `json:"account_type"`
	Balance     decimal.Decimal `json:"balance"`
	Currency    string          `json:"currency"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

// CreateAccountRequest 口座作成リクエスト
type CreateAccountRequest struct {
	Name           string           `json:"name" binding:"required,min=1,max=100"`
	AccountType    string           `json:"account_type" binding:"required,oneof=checking investment cash"`
	InitialBalance *decimal.Decimal `json:"initial_balance" binding:"omitempty,min=0"`
	Currency       string           `json:"currency" binding:"required,len=3"`
}

// UpdateAccountRequest 口座更新リクエスト
type UpdateAccountRequest struct {
	Name        *string `json:"name" binding:"omitempty,min=1,max=100"`
	AccountType *string `json:"account_type" binding:"omitempty,oneof=checking investment cash"`
}

// AccountListResponse 口座一覧レスポンス
type AccountListResponse struct {
	Accounts     []AccountResponse `json:"accounts"`
	TotalCount   int64             `json:"total_count"`
	TotalBalance decimal.Decimal   `json:"total_balance"`
}

// AccountMovementRequest 口座残高変動リクエスト
type AccountMovementRequest struct {
	Amount     decimal.Decimal `json:"amount" binding:"required,nonzero"`
	OccurredAt *time.Time      `json:"occurred_at" binding:"omitempty"`
	Note       string          `json:"note" binding:"required,min=1,max=200"`
}

// AccountMovementResponse 口座残高変動レスポンス
type AccountMovementResponse struct {
	ID         uuid.UUID       `json:"id"`
	AccountID  uuid.UUID       `json:"account_id"`
	Amount     decimal.Decimal `json:"amount"`
	OccurredAt time.Time       `json:"occurred_at"`
	Note       string          `json:"note"`
	CreatedAt  time.Time       `json:"created_at"`
}

// AccountFromDomain ドメインエンティティからDTOへの変換
func AccountFromDomain(account *accountDomain.Account) *AccountResponse {
	if account == nil {
		return nil
	}

	balance := account.CurrentBalance()
	return &AccountResponse{
		ID:          account.ID,
		UserID:      account.UserID().Value(),
		Name:        account.Name().String(),
		AccountType: account.Type().Value(),
		Balance:     decimal.NewFromInt(balance.Amount()),
		Currency:    balance.Currency(),
		CreatedAt:   account.CreatedAt,
		UpdatedAt:   account.UpdatedAt,
	}
}

// AccountsFromDomain ドメインエンティティのスライスからDTOへの変換
func AccountsFromDomain(accounts []*accountDomain.Account) []AccountResponse {
	result := make([]AccountResponse, len(accounts))
	for i, account := range accounts {
		if accountDTO := AccountFromDomain(account); accountDTO != nil {
			result[i] = *accountDTO
		}
	}
	return result
}

// AccountMovementFromDomain 口座残高変動ドメインエンティティからDTOへの変換
func AccountMovementFromDomain(movement *accountDomain.AccountMovement) *AccountMovementResponse {
	if movement == nil {
		return nil
	}

	amount := movement.Amount()
	return &AccountMovementResponse{
		ID:         movement.ID,
		AccountID:  movement.AccountID().ID,
		Amount:     decimal.NewFromInt(amount.Amount()),
		OccurredAt: movement.OccurredAt(),
		Note:       movement.Note(),
		CreatedAt:  movement.CreatedAt,
	}
}

// AccountSearchParams 口座検索パラメータ
type AccountSearchParams struct {
	UserID      uuid.UUID `form:"-"`
	AccountType *string   `form:"account_type" binding:"omitempty,oneof=checking investment cash"`
	IsActive    *bool     `form:"is_active"`
	OrderBy     string    `form:"order_by,default=display_order asc"`
}

// AccountBalanceSummary 口座残高サマリー
type AccountBalanceSummary struct {
	TotalBalance   decimal.Decimal            `json:"total_balance"`
	ByAccountType  map[string]decimal.Decimal `json:"by_account_type"`
	ByCurrency     map[string]decimal.Decimal `json:"by_currency"`
	ActiveAccounts int                        `json:"active_accounts"`
	TotalAccounts  int                        `json:"total_accounts"`
}
