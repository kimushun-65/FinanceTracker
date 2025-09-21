// Package repository 資産スナップショット関連のリポジトリインターフェースを定義
package repository

import (
	"context"
	"time"

	"financetracker/internal/domain/asset/entity"

	"github.com/google/uuid"
)

// AssetSnapshotRepository 資産スナップショットリポジトリインターフェース
type AssetSnapshotRepository interface {
	// Save 資産スナップショットを保存
	Save(ctx context.Context, snapshot entity.AssetSnapshot) error

	// FindByID IDで資産スナップショットを取得
	FindByID(ctx context.Context, id uuid.UUID) (*entity.AssetSnapshot, error)

	// FindByUserIDAndDate ユーザーIDと日付で資産スナップショットを取得
	FindByUserIDAndDate(ctx context.Context, userID uuid.UUID, date time.Time) (*entity.AssetSnapshot, error)

	// FindByUserID ユーザーIDで資産スナップショットを取得（日付降順）
	FindByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]entity.AssetSnapshot, error)

	// FindByUserIDAndDateRange ユーザーIDと日付範囲で資産スナップショットを取得
	FindByUserIDAndDateRange(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) ([]entity.AssetSnapshot, error)

	// FindLatestByUserID ユーザーの最新のスナップショットを取得
	FindLatestByUserID(ctx context.Context, userID uuid.UUID) (*entity.AssetSnapshot, error)

	// ExistsByUserIDAndDate ユーザーIDと日付でスナップショットが存在するかチェック
	ExistsByUserIDAndDate(ctx context.Context, userID uuid.UUID, date time.Time) (bool, error)

	// Delete 資産スナップショットを削除
	Delete(ctx context.Context, id uuid.UUID) error

	// CountByUserID ユーザーのスナップショット数を取得
	CountByUserID(ctx context.Context, userID uuid.UUID) (int64, error)
}
