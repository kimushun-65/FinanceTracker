// Package entity 資産スナップショット関連のエンティティを定義
package entity

import (
	"time"

	"financetracker/internal/domain/asset/value"
	"financetracker/internal/domain/common"
	commonValue "financetracker/internal/domain/common/value"

	"github.com/google/uuid"
)

// AssetSnapshot 資産スナップショットエンティティ
type AssetSnapshot struct {
	common.BaseEntity
	userID           uuid.UUID
	snapshotDate     time.Time
	totalAssets      commonValue.Money
	accountBreakdown value.AccountBreakdown
}

// NewAssetSnapshot 新しい資産スナップショットを作成
func NewAssetSnapshot(
	userID uuid.UUID,
	snapshotDate time.Time,
	totalAssets commonValue.Money,
	accountBreakdown value.AccountBreakdown,
) (AssetSnapshot, error) {
	if err := validateAssetSnapshotParams(userID, snapshotDate); err != nil {
		return AssetSnapshot{}, err
	}

	// 総資産額と内訳の合計が一致するかチェック
	if !totalAssets.Equals(accountBreakdown.TotalBalance()) {
		return AssetSnapshot{}, common.NewDomainError(
			common.DomainErrorTypeInvalidValue,
			"総資産額と口座別内訳の合計が一致しません",
		)
	}

	baseEntity := common.NewBaseEntity()

	return AssetSnapshot{
		BaseEntity:       baseEntity,
		userID:           userID,
		snapshotDate:     snapshotDate,
		totalAssets:      totalAssets,
		accountBreakdown: accountBreakdown,
	}, nil
}

// validateAssetSnapshotParams 資産スナップショットパラメータの妥当性を検証
func validateAssetSnapshotParams(userID uuid.UUID, snapshotDate time.Time) error {
	if userID == uuid.Nil {
		return common.NewDomainError(
			common.DomainErrorTypeInvalidValue,
			"ユーザーIDが必要です",
		)
	}

	if snapshotDate.IsZero() {
		return common.NewDomainError(
			common.DomainErrorTypeInvalidValue,
			"スナップショット日付が必要です",
		)
	}

	// 未来の日付は許可しない
	if snapshotDate.After(time.Now()) {
		return common.NewDomainError(
			common.DomainErrorTypeInvalidValue,
			"未来の日付でスナップショットを作成することはできません",
		)
	}

	return nil
}

// UserID ユーザーIDを取得
func (a AssetSnapshot) UserID() uuid.UUID {
	return a.userID
}

// SnapshotDate スナップショット日付を取得
func (a AssetSnapshot) SnapshotDate() time.Time {
	return a.snapshotDate
}

// TotalAssets 総資産額を取得
func (a AssetSnapshot) TotalAssets() commonValue.Money {
	return a.totalAssets
}

// AccountBreakdown 口座別内訳を取得
func (a AssetSnapshot) AccountBreakdown() value.AccountBreakdown {
	return a.accountBreakdown
}

// IsForDate 指定された日付のスナップショットかどうかを判定
func (a AssetSnapshot) IsForDate(date time.Time) bool {
	// 日付部分のみで比較（時刻は無視）
	y1, m1, d1 := a.snapshotDate.Date()
	y2, m2, d2 := date.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

// GetAccountBalance 特定の口座の残高を取得
func (a AssetSnapshot) GetAccountBalance(accountID uuid.UUID) (*commonValue.Money, bool) {
	balance, found := a.accountBreakdown.FindByAccountID(accountID)
	if !found {
		return nil, false
	}
	return &balance.Balance, true
}

// CalculateChangeFrom 別のスナップショットとの差額を計算
func (a AssetSnapshot) CalculateChangeFrom(previous AssetSnapshot) (*commonValue.Money, error) {
	// 同じユーザーのスナップショットかチェック
	if a.userID != previous.userID {
		return nil, common.NewDomainError(
			common.DomainErrorTypeInvalidValue,
			"異なるユーザーのスナップショットとは比較できません",
		)
	}

	// より新しいスナップショットかチェック
	if !a.snapshotDate.After(previous.snapshotDate) {
		return nil, common.NewDomainError(
			common.DomainErrorTypeInvalidValue,
			"比較対象は過去のスナップショットである必要があります",
		)
	}

	return a.totalAssets.Subtract(previous.totalAssets)
}

// ReconstructAssetSnapshot データベースから取得したデータからAssetSnapshotを再構築
func ReconstructAssetSnapshot(
	baseEntity common.BaseEntity,
	userID uuid.UUID,
	snapshotDate time.Time,
	totalAssets commonValue.Money,
	accountBreakdown value.AccountBreakdown,
) AssetSnapshot {
	return AssetSnapshot{
		BaseEntity:       baseEntity,
		userID:           userID,
		snapshotDate:     snapshotDate,
		totalAssets:      totalAssets,
		accountBreakdown: accountBreakdown,
	}
}
