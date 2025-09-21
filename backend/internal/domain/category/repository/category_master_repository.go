// Package repository カテゴリ関連のリポジトリインターフェースを定義
package repository

import (
	"context"

	"financetracker/internal/domain/category/entity"
	"financetracker/internal/domain/category/value"

	"github.com/google/uuid"
)

// CategoryMasterRepository カテゴリマスターリポジトリインターフェース
type CategoryMasterRepository interface {
	// Save カテゴリマスターを保存
	Save(ctx context.Context, categoryMaster entity.CategoryMaster) error

	// FindByID IDでカテゴリマスターを取得
	FindByID(ctx context.Context, id uuid.UUID) (*entity.CategoryMaster, error)

	// FindAll すべてのカテゴリマスターを取得
	FindAll(ctx context.Context) ([]entity.CategoryMaster, error)

	// FindByType タイプでカテゴリマスターを取得
	FindByType(ctx context.Context, categoryType value.CategoryType) ([]entity.CategoryMaster, error)

	// Update カテゴリマスターを更新
	Update(ctx context.Context, categoryMaster entity.CategoryMaster) error

	// Delete カテゴリマスターを削除
	Delete(ctx context.Context, id uuid.UUID) error
}
