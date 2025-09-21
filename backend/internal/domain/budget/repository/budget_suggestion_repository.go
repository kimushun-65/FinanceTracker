// Package repository 予算提案関連のリポジトリインターフェースを定義
package repository

import (
	"context"

	"financetracker/internal/domain/budget/entity"

	"github.com/google/uuid"
)

// BudgetSuggestionRepository 予算提案リポジトリインターフェース
type BudgetSuggestionRepository interface {
	// Save 予算提案を保存
	Save(ctx context.Context, suggestion entity.BudgetSuggestion) error

	// FindByID IDで予算提案を取得
	FindByID(ctx context.Context, id uuid.UUID) (*entity.BudgetSuggestion, error)

	// FindByUserID ユーザーIDで予算提案を取得
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]entity.BudgetSuggestion, error)

	// FindByUserIDAndMonth ユーザーIDと年月で予算提案を取得
	FindByUserIDAndMonth(ctx context.Context, userID uuid.UUID, month string) ([]entity.BudgetSuggestion, error)

	// FindPendingByUserID ユーザーの検討中の提案を取得
	FindPendingByUserID(ctx context.Context, userID uuid.UUID) ([]entity.BudgetSuggestion, error)

	// FindLatestByUserIDAndCategoryID ユーザーIDとカテゴリIDで最新の提案を取得
	FindLatestByUserIDAndCategoryID(ctx context.Context, userID, categoryID uuid.UUID) (*entity.BudgetSuggestion, error)

	// Update 予算提案を更新
	Update(ctx context.Context, suggestion entity.BudgetSuggestion) error

	// Delete 予算提案を削除
	Delete(ctx context.Context, id uuid.UUID) error
}
