// Package repository カテゴリ関連のリポジトリインターフェースを定義
package repository

import (
	"context"

	"financetracker/internal/domain/category/entity"

	"github.com/google/uuid"
)

// CategoryRepository ユーザーカテゴリリポジトリインターフェース
type CategoryRepository interface {
	// Save カテゴリを保存
	Save(ctx context.Context, category entity.Category) error

	// FindByID IDでカテゴリを取得
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Category, error)

	// FindByUserID ユーザーIDでカテゴリを取得
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]entity.Category, error)

	// FindByUserIDAndCategoryMasterID ユーザーIDとカテゴリマスターIDで取得
	FindByUserIDAndCategoryMasterID(ctx context.Context, userID, categoryMasterID uuid.UUID) (*entity.Category, error)

	// FindActiveByUserID ユーザーのアクティブなカテゴリを取得
	FindActiveByUserID(ctx context.Context, userID uuid.UUID) ([]entity.Category, error)

	// Update カテゴリを更新
	Update(ctx context.Context, category entity.Category) error

	// ExistsByUserIDAndCategoryMasterID ユーザーとカテゴリマスターの組み合わせが存在するかチェック
	ExistsByUserIDAndCategoryMasterID(ctx context.Context, userID, categoryMasterID uuid.UUID) (bool, error)
}
