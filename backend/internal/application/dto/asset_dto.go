package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	assetDomain "financetracker/internal/domain/asset/entity"
)

// AssetSnapshotResponse 資産スナップショットレスポンス
type AssetSnapshotResponse struct {
	ID           uuid.UUID                `json:"id"`
	UserID       uuid.UUID                `json:"user_id"`
	SnapshotDate time.Time                `json:"snapshot_date"`
	TotalAssets  decimal.Decimal          `json:"total_assets"`
	Accounts     []AccountBalanceResponse `json:"accounts"`
	CreatedAt    time.Time                `json:"created_at"`
}

// AccountBalanceResponse 口座残高レスポンス
type AccountBalanceResponse struct {
	AccountID   uuid.UUID       `json:"account_id"`
	AccountName string          `json:"account_name"`
	Balance     decimal.Decimal `json:"balance"`
	Currency    string          `json:"currency"`
}

// AssetSnapshotListResponse 資産スナップショット一覧レスポンス
type AssetSnapshotListResponse struct {
	Snapshots  []AssetSnapshotResponse `json:"snapshots"`
	TotalCount int                     `json:"total_count"`
}

// CreateAssetSnapshotRequest 資産スナップショット作成リクエスト
type CreateAssetSnapshotRequest struct {
	SnapshotDate time.Time               `json:"snapshot_date" binding:"required"`
	Accounts     []AccountBalanceRequest `json:"accounts" binding:"required"`
}

// AccountBalanceRequest 口座残高リクエスト
type AccountBalanceRequest struct {
	AccountID uuid.UUID       `json:"account_id" binding:"required"`
	Balance   decimal.Decimal `json:"balance" binding:"required"`
}

// AssetSnapshotFromDomain ドメインエンティティからDTOへの変換
func AssetSnapshotFromDomain(snapshot *assetDomain.AssetSnapshot) *AssetSnapshotResponse {
	if snapshot == nil {
		return nil
	}

	totalAssets := snapshot.TotalAssets()
	accountBreakdown := snapshot.AccountBreakdown()

	// 口座別内訳を変換
	accountBalances := accountBreakdown.Accounts()
	accounts := make([]AccountBalanceResponse, len(accountBalances))
	for i, ab := range accountBalances {
		accounts[i] = AccountBalanceResponse{
			AccountID:   ab.AccountID,
			AccountName: ab.AccountName,
			Balance:     decimal.NewFromInt(ab.Balance.Amount()),
			Currency:    ab.Balance.Currency(),
		}
	}

	return &AssetSnapshotResponse{
		ID:           snapshot.ID,
		UserID:       snapshot.UserID(),
		SnapshotDate: snapshot.SnapshotDate(),
		TotalAssets:  decimal.NewFromInt(totalAssets.Amount()),
		Accounts:     accounts,
		CreatedAt:    snapshot.CreatedAt,
	}
}

// AssetSnapshotsFromDomain ドメインエンティティのスライスからDTOへの変換
func AssetSnapshotsFromDomain(snapshots []assetDomain.AssetSnapshot) []AssetSnapshotResponse {
	result := make([]AssetSnapshotResponse, len(snapshots))
	for i, snapshot := range snapshots {
		if snapshotDTO := AssetSnapshotFromDomain(&snapshot); snapshotDTO != nil {
			result[i] = *snapshotDTO
		}
	}
	return result
}
