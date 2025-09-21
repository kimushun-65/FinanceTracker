package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"financetracker/internal/domain/account/entity"
	accountValue "financetracker/internal/domain/account/value"
	"financetracker/internal/domain/common/repository"
	"financetracker/internal/domain/common/value"
	userValue "financetracker/internal/domain/user/value"
)

// AccountRepository 口座リポジトリのインターフェース
type AccountRepository interface {
	repository.UserScopedRepository[*entity.Account]

	// FindByUserIDAndType ユーザーIDと口座タイプで検索
	FindByUserIDAndType(ctx context.Context, userID userValue.UserID, accountType accountValue.AccountType) ([]*entity.Account, error)

	// GetTotalBalanceByUserID ユーザーの総残高を取得（通貨別）
	GetTotalBalanceByUserID(ctx context.Context, userID userValue.UserID) (map[string]value.Money, error)

	// GetTotalAssetsByUserID ユーザーの総資産を取得
	GetTotalAssetsByUserID(ctx context.Context, userID userValue.UserID) (value.Money, error)
}

// AccountMovementRepository 口座移動履歴リポジトリのインターフェース
type AccountMovementRepository interface {
	repository.BaseRepository[*entity.AccountMovement]

	// FindByAccountID 口座IDで移動履歴を取得
	FindByAccountID(ctx context.Context, accountID uuid.UUID, pagination *repository.Pagination) (*repository.PagedResult[*entity.AccountMovement], error)

	// FindByUserID ユーザーIDで移動履歴を取得（全口座）
	FindByUserID(ctx context.Context, userID userValue.UserID, pagination *repository.Pagination) (*repository.PagedResult[*entity.AccountMovement], error)

	// FindByUserIDAndDateRange ユーザーIDと日付範囲で移動履歴を取得
	FindByUserIDAndDateRange(ctx context.Context, userID userValue.UserID, startDate, endDate time.Time, pagination *repository.Pagination) (*repository.PagedResult[*entity.AccountMovement], error)

	// FindByAccountIDAndDateRange 口座IDと日付範囲で移動履歴を取得
	FindByAccountIDAndDateRange(ctx context.Context, accountID uuid.UUID, startDate, endDate time.Time, pagination *repository.Pagination) (*repository.PagedResult[*entity.AccountMovement], error)

	// DeleteByAccountID 口座IDに関連する全ての移動履歴を削除
	DeleteByAccountID(ctx context.Context, accountID uuid.UUID) error

	// GetMovementSummary 移動履歴のサマリーを取得
	GetMovementSummary(ctx context.Context, userID userValue.UserID, startDate, endDate time.Time) (*MovementSummary, error)
}

// AccountFilter 口座検索フィルター
type AccountFilter struct {
	AccountType   *accountValue.AccountType `json:"account_type,omitempty"`   // 口座タイプでフィルタ
	NameQuery     *string                   `json:"name_query,omitempty"`     // 口座名の部分一致検索
	HasPositiveBalance *bool               `json:"has_positive_balance,omitempty"` // プラス残高のみ
	HasNegativeBalance *bool               `json:"has_negative_balance,omitempty"` // マイナス残高のみ
}

// ToMap フィルター条件をマップ形式で取得
func (f AccountFilter) ToMap() map[string]any {
	result := make(map[string]any)
	
	if f.AccountType != nil {
		result["account_type"] = f.AccountType.Value()
	}
	if f.NameQuery != nil {
		result["name_query"] = *f.NameQuery
	}
	if f.HasPositiveBalance != nil {
		result["has_positive_balance"] = *f.HasPositiveBalance
	}
	if f.HasNegativeBalance != nil {
		result["has_negative_balance"] = *f.HasNegativeBalance
	}
	
	return result
}

// IsEmpty フィルターが空かどうかを判定
func (f AccountFilter) IsEmpty() bool {
	return f.AccountType == nil && 
		   f.NameQuery == nil && 
		   f.HasPositiveBalance == nil && 
		   f.HasNegativeBalance == nil
}

// MovementFilter 移動履歴検索フィルター
type MovementFilter struct {
	repository.DateRangeFilter
	MovementType *entity.MovementType `json:"movement_type,omitempty"` // 移動タイプでフィルタ
	AmountMin    *int64               `json:"amount_min,omitempty"`    // 最小金額
	AmountMax    *int64               `json:"amount_max,omitempty"`    // 最大金額
	NoteQuery    *string              `json:"note_query,omitempty"`    // メモの部分一致検索
}

// ToMap フィルター条件をマップ形式で取得
func (f MovementFilter) ToMap() map[string]any {
	result := f.DateRangeFilter.ToMap()
	
	if f.MovementType != nil {
		result["movement_type"] = string(*f.MovementType)
	}
	if f.AmountMin != nil {
		result["amount_min"] = *f.AmountMin
	}
	if f.AmountMax != nil {
		result["amount_max"] = *f.AmountMax
	}
	if f.NoteQuery != nil {
		result["note_query"] = *f.NoteQuery
	}
	
	return result
}

// IsEmpty フィルターが空かどうかを判定
func (f MovementFilter) IsEmpty() bool {
	return f.DateRangeFilter.IsEmpty() && 
		   f.MovementType == nil && 
		   f.AmountMin == nil && 
		   f.AmountMax == nil && 
		   f.NoteQuery == nil
}

// MovementSummary 移動履歴のサマリー
type MovementSummary struct {
	TotalDeposits    value.Money `json:"total_deposits"`    // 総入金額
	TotalWithdrawals value.Money `json:"total_withdrawals"` // 総出金額
	NetMovement      value.Money `json:"net_movement"`      // 純移動額（入金 - 出金）
	MovementCount    int64       `json:"movement_count"`    // 移動回数
	DepositCount     int64       `json:"deposit_count"`     // 入金回数
	WithdrawalCount  int64       `json:"withdrawal_count"`  // 出金回数
}

// AccountSummary 口座サマリー
type AccountSummary struct {
	TotalAccounts int64       `json:"total_accounts"` // 総口座数
	TotalAssets   value.Money `json:"total_assets"`   // 総資産
}