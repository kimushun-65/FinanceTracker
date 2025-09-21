// Package repository 予算関連のリポジトリインターフェースを定義
package repository

import (
	"context"
	"time"

	"financetracker/internal/domain/budget/entity"

	"github.com/google/uuid"
)

// BudgetRepository 予算リポジトリインターフェース
type BudgetRepository interface {
	// Save 予算を保存
	Save(ctx context.Context, budget entity.Budget) error

	// FindByID IDで予算を取得
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Budget, error)

	// FindByUserID ユーザーIDで予算を取得
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]entity.Budget, error)

	// FindByUserIDAndCategoryID ユーザーIDとカテゴリIDで予算を取得
	FindByUserIDAndCategoryID(ctx context.Context, userID, categoryID uuid.UUID) ([]entity.Budget, error)

	// FindActiveByUserID ユーザーのアクティブな予算を取得
	FindActiveByUserID(ctx context.Context, userID uuid.UUID) ([]entity.Budget, error)

	// FindByUserIDAndDate ユーザーIDと日付で有効な予算を取得
	FindByUserIDAndDate(ctx context.Context, userID uuid.UUID, date time.Time) ([]entity.Budget, error)

	// Update 予算を更新
	Update(ctx context.Context, budget entity.Budget) error

	// Delete 予算を削除
	Delete(ctx context.Context, id uuid.UUID) error
}
