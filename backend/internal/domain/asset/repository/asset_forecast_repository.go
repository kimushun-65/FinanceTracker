// Package repository 資産予測関連のリポジトリインターフェースを定義
package repository

import (
	"context"
	"time"

	"financetracker/internal/domain/asset/entity"

	"github.com/google/uuid"
)

// AssetForecastRepository 資産予測リポジトリインターフェース
type AssetForecastRepository interface {
	// Save 資産予測を保存
	Save(ctx context.Context, forecast entity.AssetForecast) error

	// FindByID IDで資産予測を取得
	FindByID(ctx context.Context, id uuid.UUID) (*entity.AssetForecast, error)

	// FindByUserID ユーザーIDで資産予測を取得（作成日降順）
	FindByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]entity.AssetForecast, error)

	// FindLatestByUserID ユーザーの最新の予測を取得
	FindLatestByUserID(ctx context.Context, userID uuid.UUID) (*entity.AssetForecast, error)

	// FindByUserIDAndHorizon ユーザーIDと予測期間で取得
	FindByUserIDAndHorizon(ctx context.Context, userID uuid.UUID, horizonMonths int) ([]entity.AssetForecast, error)

	// FindValidByUserID 有効な予測を取得（1ヶ月以内に作成されたもの）
	FindValidByUserID(ctx context.Context, userID uuid.UUID) ([]entity.AssetForecast, error)

	// FindByUserIDAndDateRange 生成日の範囲で予測を取得
	FindByUserIDAndDateRange(ctx context.Context, userID uuid.UUID, startDate, endDate time.Time) ([]entity.AssetForecast, error)

	// Delete 資産予測を削除
	Delete(ctx context.Context, id uuid.UUID) error

	// DeleteOldForecasts 古い予測を削除
	DeleteOldForecasts(ctx context.Context, userID uuid.UUID, keepDays int) error
}
