package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"

	assetDomain "financetracker/internal/domain/asset/entity"
	assetRepo "financetracker/internal/domain/asset/repository"
	assetValue "financetracker/internal/domain/asset/value"
	"financetracker/internal/domain/common"
	commonValue "financetracker/internal/domain/common/value"
	"financetracker/internal/infrastructure/gorm/model"
)

// AssetSnapshotRepository 資産スナップショットリポジトリの実装
type AssetSnapshotRepository struct {
	db *gorm.DB
}

// NewAssetSnapshotRepository 新しいAssetSnapshotRepositoryを作成
func NewAssetSnapshotRepository(db *gorm.DB) assetRepo.AssetSnapshotRepository {
	return &AssetSnapshotRepository{db: db}
}

// Save 資産スナップショットを保存
func (r *AssetSnapshotRepository) Save(ctx context.Context, snapshot assetDomain.AssetSnapshot) error {
	snapshotModel, err := r.toModel(&snapshot)
	if err != nil {
		return fmt.Errorf("モデルへの変換に失敗しました: %w", err)
	}

	// IDが存在する場合は更新、存在しない場合は作成
	result := r.db.WithContext(ctx).Save(snapshotModel)
	if result.Error != nil {
		return fmt.Errorf("資産スナップショットの保存に失敗しました: %w", result.Error)
	}

	return nil
}

// FindByID IDで資産スナップショットを取得
func (r *AssetSnapshotRepository) FindByID(ctx context.Context, id uuid.UUID) (*assetDomain.AssetSnapshot, error) {
	var snapshotModel model.AssetSnapshot
	result := r.db.WithContext(ctx).First(&snapshotModel, "id = ?", id)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("資産スナップショットの取得に失敗しました: %w", result.Error)
	}

	return r.toDomain(&snapshotModel)
}

// FindByUserIDAndDate ユーザーIDと日付で資産スナップショットを取得
func (r *AssetSnapshotRepository) FindByUserIDAndDate(ctx context.Context, userID uuid.UUID, date time.Time) (*assetDomain.AssetSnapshot, error) {
	var snapshotModel model.AssetSnapshot
	// 日付部分のみで比較
	dateOnly := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)

	result := r.db.WithContext(ctx).Where("user_id = ? AND date = ?", userID, dateOnly).First(&snapshotModel)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("資産スナップショットの取得に失敗しました: %w", result.Error)
	}

	return r.toDomain(&snapshotModel)
}

// FindByUserID ユーザーIDで資産スナップショットを取得（日付降順）
func (r *AssetSnapshotRepository) FindByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]assetDomain.AssetSnapshot, error) {
	var snapshotModels []model.AssetSnapshot
	query := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("date DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Find(&snapshotModels).Error; err != nil {
		return nil, fmt.Errorf("資産スナップショット一覧の取得に失敗しました: %w", err)
	}

	// ドメインモデルに変換
	snapshots := make([]assetDomain.AssetSnapshot, 0, len(snapshotModels))
	for _, snapshotModel := range snapshotModels {
		snapshot, err := r.toDomain(&snapshotModel)
		if err != nil {
			return nil, err
		}
		snapshots = append(snapshots, *snapshot)
	}

	return snapshots, nil
}

// FindByUserIDAndDateRange ユーザーIDと日付範囲で資産スナップショットを取得
func (r *AssetSnapshotRepository) FindByUserIDAndDateRange(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) ([]assetDomain.AssetSnapshot, error) {
	var snapshotModels []model.AssetSnapshot

	// 日付部分のみで比較
	startDateOnly := time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, time.UTC)
	endDateOnly := time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 0, time.UTC)

	result := r.db.WithContext(ctx).
		Where("user_id = ? AND date >= ? AND date <= ?", userID, startDateOnly, endDateOnly).
		Order("date ASC").
		Find(&snapshotModels)

	if result.Error != nil {
		return nil, fmt.Errorf("資産スナップショット一覧の取得に失敗しました: %w", result.Error)
	}

	// ドメインモデルに変換
	snapshots := make([]assetDomain.AssetSnapshot, 0, len(snapshotModels))
	for _, snapshotModel := range snapshotModels {
		snapshot, err := r.toDomain(&snapshotModel)
		if err != nil {
			return nil, err
		}
		snapshots = append(snapshots, *snapshot)
	}

	return snapshots, nil
}

// FindLatestByUserID ユーザーの最新のスナップショットを取得
func (r *AssetSnapshotRepository) FindLatestByUserID(ctx context.Context, userID uuid.UUID) (*assetDomain.AssetSnapshot, error) {
	var snapshotModel model.AssetSnapshot
	result := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("date DESC").
		First(&snapshotModel)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("最新の資産スナップショットの取得に失敗しました: %w", result.Error)
	}

	return r.toDomain(&snapshotModel)
}

// ExistsByUserIDAndDate ユーザーIDと日付でスナップショットが存在するかチェック
func (r *AssetSnapshotRepository) ExistsByUserIDAndDate(ctx context.Context, userID uuid.UUID, date time.Time) (bool, error) {
	var count int64
	dateOnly := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)

	result := r.db.WithContext(ctx).
		Model(&model.AssetSnapshot{}).
		Where("user_id = ? AND date = ?", userID, dateOnly).
		Count(&count)

	if result.Error != nil {
		return false, fmt.Errorf("資産スナップショットの存在確認に失敗しました: %w", result.Error)
	}

	return count > 0, nil
}

// Delete 資産スナップショットを削除
func (r *AssetSnapshotRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&model.AssetSnapshot{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("資産スナップショットの削除に失敗しました: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("削除対象の資産スナップショットが見つかりません: %s", id)
	}

	return nil
}

// CountByUserID ユーザーのスナップショット数を取得
func (r *AssetSnapshotRepository) CountByUserID(ctx context.Context, userID uuid.UUID) (int64, error) {
	var count int64
	result := r.db.WithContext(ctx).
		Model(&model.AssetSnapshot{}).
		Where("user_id = ?", userID).
		Count(&count)

	if result.Error != nil {
		return 0, fmt.Errorf("資産スナップショット数の取得に失敗しました: %w", result.Error)
	}

	return count, nil
}

// toModel ドメインモデルからGORMモデルへの変換
func (r *AssetSnapshotRepository) toModel(snapshot *assetDomain.AssetSnapshot) (*model.AssetSnapshot, error) {
	totalAssets := snapshot.TotalAssets()

	// AccountBreakdownをJSON文字列に変換
	accountBreakdownJSON, err := snapshot.AccountBreakdown().ToJSON()
	if err != nil {
		return nil, fmt.Errorf("口座別内訳のJSON変換に失敗しました: %w", err)
	}

	return &model.AssetSnapshot{
		Base: model.Base{
			ID:        snapshot.ID,
			CreatedAt: snapshot.CreatedAt,
			UpdatedAt: snapshot.UpdatedAt,
		},
		UserID:           snapshot.UserID(),
		Date:             snapshot.SnapshotDate(),
		TotalAssets:      decimal.NewFromInt(totalAssets.Amount()),
		TotalLiabilities: decimal.Zero, // 現時点では負債は扱わない
		NetWorth:         decimal.NewFromInt(totalAssets.Amount()),
		AccountBreakdown: accountBreakdownJSON,
	}, nil
}

// toDomain GORMモデルからドメインモデルへの変換
func (r *AssetSnapshotRepository) toDomain(snapshotModel *model.AssetSnapshot) (*assetDomain.AssetSnapshot, error) {
	// 総資産を作成
	totalAssets, err := commonValue.NewMoney(snapshotModel.TotalAssets.IntPart(), "JPY")
	if err != nil {
		return nil, fmt.Errorf("総資産の作成に失敗しました: %w", err)
	}

	// AccountBreakdownをJSONからパース
	var accountBreakdownData map[string]interface{}
	if snapshotModel.AccountBreakdown != "" {
		if err := json.Unmarshal([]byte(snapshotModel.AccountBreakdown), &accountBreakdownData); err != nil {
			return nil, fmt.Errorf("口座別内訳のJSONパースに失敗しました: %w", err)
		}
	}

	// AccountBalanceのスライスを作成
	var accountBalances []assetValue.AccountBalance
	if accounts, ok := accountBreakdownData["accounts"].([]interface{}); ok {
		for _, acc := range accounts {
			if accountMap, ok := acc.(map[string]interface{}); ok {
				accountIDStr, _ := accountMap["account_id"].(string)
				accountName, _ := accountMap["account_name"].(string)
				balanceFloat, _ := accountMap["balance"].(float64)
				currency, _ := accountMap["currency"].(string)

				accountID, err := uuid.Parse(accountIDStr)
				if err != nil {
					continue
				}

				balance, err := commonValue.NewMoney(int64(balanceFloat), currency)
				if err != nil {
					continue
				}

				accountBalances = append(accountBalances, assetValue.AccountBalance{
					AccountID:   accountID,
					AccountName: accountName,
					Balance:     *balance,
				})
			}
		}
	}

	// AccountBreakdownを作成
	accountBreakdown, err := assetValue.NewAccountBreakdown(accountBalances)
	if err != nil {
		return nil, fmt.Errorf("口座別内訳の作成に失敗しました: %w", err)
	}

	// BaseEntity
	baseEntity := common.BaseEntity{
		ID:        snapshotModel.ID,
		CreatedAt: snapshotModel.CreatedAt,
		UpdatedAt: snapshotModel.UpdatedAt,
	}

	// ドメインエンティティを再構築
	snapshot := assetDomain.ReconstructAssetSnapshot(
		baseEntity,
		snapshotModel.UserID,
		snapshotModel.Date,
		*totalAssets,
		accountBreakdown,
	)

	return &snapshot, nil
}
