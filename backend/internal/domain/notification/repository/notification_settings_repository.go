// Package repository 通知設定関連のリポジトリインターフェースを定義
package repository

import (
	"context"

	"financetracker/internal/domain/notification/entity"

	"github.com/google/uuid"
)

// NotificationSettingsRepository 通知設定リポジトリインターフェース
type NotificationSettingsRepository interface {
	// Save 通知設定を保存
	Save(ctx context.Context, settings entity.NotificationSettings) error

	// FindByUserID ユーザーIDで通知設定を取得
	FindByUserID(ctx context.Context, userID uuid.UUID) (*entity.NotificationSettings, error)

	// Update 通知設定を更新
	Update(ctx context.Context, settings entity.NotificationSettings) error

	// Delete 通知設定を削除
	Delete(ctx context.Context, userID uuid.UUID) error

	// ExistsByUserID ユーザーの通知設定が存在するかチェック
	ExistsByUserID(ctx context.Context, userID uuid.UUID) (bool, error)

	// FindAllEnabledMonthlyReports 有効な月次レポート設定をすべて取得
	FindAllEnabledMonthlyReports(ctx context.Context) ([]entity.NotificationSettings, error)

	// FindAllEnabledBudgetAlerts 有効な予算超過アラート設定をすべて取得
	FindAllEnabledBudgetAlerts(ctx context.Context) ([]entity.NotificationSettings, error)
}
