package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	categoryDomain "financetracker/internal/domain/category/entity"
	categoryRepo "financetracker/internal/domain/category/repository"
	categoryValue "financetracker/internal/domain/category/value"
	"financetracker/internal/domain/common"
	"financetracker/internal/infrastructure/gorm/model"
)

// CategoryRepository カテゴリリポジトリの実装
type CategoryRepository struct {
	db *gorm.DB
}

// NewCategoryRepository 新しいCategoryRepositoryを作成
func NewCategoryRepository(db *gorm.DB) categoryRepo.CategoryRepository {
	return &CategoryRepository{db: db}
}

// Save カテゴリを保存（新規作成または更新）
func (r *CategoryRepository) Save(ctx context.Context, category categoryDomain.Category) error {
	categoryModel := r.toModel(&category)

	// IDが存在する場合は更新、存在しない場合は作成
	result := r.db.WithContext(ctx).Save(categoryModel)
	if result.Error != nil {
		return fmt.Errorf("カテゴリの保存に失敗しました: %w", result.Error)
	}

	return nil
}

// FindByID IDでカテゴリを取得
func (r *CategoryRepository) FindByID(ctx context.Context, id uuid.UUID) (*categoryDomain.Category, error) {
	var categoryModel model.Category
	result := r.db.WithContext(ctx).First(&categoryModel, "id = ?", id)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("カテゴリの取得に失敗しました: %w", result.Error)
	}

	return r.toDomain(&categoryModel)
}

// FindByUserID ユーザーIDでカテゴリを取得
func (r *CategoryRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]categoryDomain.Category, error) {
	var categoryModels []model.Category
	result := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at ASC").Find(&categoryModels)

	if result.Error != nil {
		return nil, fmt.Errorf("カテゴリ一覧の取得に失敗しました: %w", result.Error)
	}

	// ドメインモデルに変換
	categories := make([]categoryDomain.Category, len(categoryModels))
	for i, categoryModel := range categoryModels {
		category, err := r.toDomain(&categoryModel)
		if err != nil {
			return nil, err
		}
		categories[i] = *category
	}

	return categories, nil
}

// FindByUserIDAndCategoryMasterID ユーザーIDとカテゴリマスターIDで取得
func (r *CategoryRepository) FindByUserIDAndCategoryMasterID(ctx context.Context, userID, categoryMasterID uuid.UUID) (*categoryDomain.Category, error) {
	var categoryModel model.Category
	result := r.db.WithContext(ctx).Where("user_id = ? AND category_master_id = ?", userID, categoryMasterID).First(&categoryModel)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("カテゴリの取得に失敗しました: %w", result.Error)
	}

	return r.toDomain(&categoryModel)
}

// FindActiveByUserID ユーザーのアクティブなカテゴリを取得
func (r *CategoryRepository) FindActiveByUserID(ctx context.Context, userID uuid.UUID) ([]categoryDomain.Category, error) {
	var categoryModels []model.Category
	result := r.db.WithContext(ctx).Where("user_id = ? AND is_active = ?", userID, true).Order("created_at ASC").Find(&categoryModels)

	if result.Error != nil {
		return nil, fmt.Errorf("アクティブカテゴリ一覧の取得に失敗しました: %w", result.Error)
	}

	// ドメインモデルに変換
	categories := make([]categoryDomain.Category, len(categoryModels))
	for i, categoryModel := range categoryModels {
		category, err := r.toDomain(&categoryModel)
		if err != nil {
			return nil, err
		}
		categories[i] = *category
	}

	return categories, nil
}

// Update カテゴリを更新
func (r *CategoryRepository) Update(ctx context.Context, category categoryDomain.Category) error {
	categoryModel := r.toModel(&category)

	result := r.db.WithContext(ctx).Model(&model.Category{}).
		Where("id = ?", category.ID).
		Updates(categoryModel)

	if result.Error != nil {
		return fmt.Errorf("カテゴリの更新に失敗しました: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("更新対象のカテゴリが見つかりません: %s", category.ID)
	}

	return nil
}

// ExistsByUserIDAndCategoryMasterID ユーザーとカテゴリマスターの組み合わせが存在するかチェック
func (r *CategoryRepository) ExistsByUserIDAndCategoryMasterID(ctx context.Context, userID, categoryMasterID uuid.UUID) (bool, error) {
	var count int64
	result := r.db.WithContext(ctx).Model(&model.Category{}).
		Where("user_id = ? AND category_master_id = ?", userID, categoryMasterID).
		Count(&count)

	if result.Error != nil {
		return false, fmt.Errorf("カテゴリの存在確認に失敗しました: %w", result.Error)
	}

	return count > 0, nil
}

// toModel ドメインモデルからGORMモデルへの変換
func (r *CategoryRepository) toModel(category *categoryDomain.Category) *model.Category {
	var customName *string
	if category.CustomName() != nil {
		name := category.CustomName().String()
		customName = &name
	}

	return &model.Category{
		Base: model.Base{
			ID:        category.ID,
			CreatedAt: category.CreatedAt,
			UpdatedAt: category.UpdatedAt,
		},
		UserID:           category.UserID(),
		CategoryMasterID: category.CategoryMasterID(),
		CustomName:       customName,
		IsActive:         category.IsActive(),
	}
}

// toDomain GORMモデルからドメインモデルへの変換
func (r *CategoryRepository) toDomain(categoryModel *model.Category) (*categoryDomain.Category, error) {
	// カスタム名称
	var customName *categoryValue.CategoryName
	if categoryModel.CustomName != nil && *categoryModel.CustomName != "" {
		name, err := categoryValue.NewCategoryName(*categoryModel.CustomName)
		if err != nil {
			return nil, fmt.Errorf("カスタム名称の作成に失敗しました: %w", err)
		}
		customName = &name
	}

	// ドメインエンティティを作成
	category, err := categoryDomain.NewCategory(
		categoryModel.UserID,
		categoryModel.CategoryMasterID,
		customName,
	)
	if err != nil {
		return nil, fmt.Errorf("カテゴリエンティティの作成に失敗しました: %w", err)
	}

	// BaseEntityを設定
	category.BaseEntity = common.BaseEntity{
		ID:        categoryModel.ID,
		CreatedAt: categoryModel.CreatedAt,
		UpdatedAt: categoryModel.UpdatedAt,
	}

	// アクティブ状態を設定
	if !categoryModel.IsActive {
		category.Deactivate()
	}

	return &category, nil
}
