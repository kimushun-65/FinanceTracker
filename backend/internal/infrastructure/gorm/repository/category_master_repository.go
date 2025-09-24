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
	commonValue "financetracker/internal/domain/common/value"
	"financetracker/internal/infrastructure/gorm/model"
)

// CategoryMasterRepository カテゴリマスターリポジトリの実装
type CategoryMasterRepository struct {
	db *gorm.DB
}

// NewCategoryMasterRepository 新しいCategoryMasterRepositoryを作成
func NewCategoryMasterRepository(db *gorm.DB) categoryRepo.CategoryMasterRepository {
	return &CategoryMasterRepository{db: db}
}

// Save カテゴリマスターを保存（新規作成または更新）
func (r *CategoryMasterRepository) Save(ctx context.Context, categoryMaster categoryDomain.CategoryMaster) error {
	categoryMasterModel := r.toModel(&categoryMaster)

	// IDが存在する場合は更新、存在しない場合は作成
	result := r.db.WithContext(ctx).Save(categoryMasterModel)
	if result.Error != nil {
		return fmt.Errorf("カテゴリマスターの保存に失敗しました: %w", result.Error)
	}

	return nil
}

// FindByID IDでカテゴリマスターを取得
func (r *CategoryMasterRepository) FindByID(ctx context.Context, id uuid.UUID) (*categoryDomain.CategoryMaster, error) {
	var categoryMasterModel model.CategoryMaster
	result := r.db.WithContext(ctx).First(&categoryMasterModel, "id = ?", id)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("カテゴリマスターの取得に失敗しました: %w", result.Error)
	}

	return r.toDomain(&categoryMasterModel)
}

// FindAll すべてのカテゴリマスターを取得
func (r *CategoryMasterRepository) FindAll(ctx context.Context) ([]categoryDomain.CategoryMaster, error) {
	var categoryMasterModels []model.CategoryMaster
	result := r.db.WithContext(ctx).Order("display_order ASC, created_at ASC").Find(&categoryMasterModels)

	if result.Error != nil {
		return nil, fmt.Errorf("カテゴリマスター一覧の取得に失敗しました: %w", result.Error)
	}

	// ドメインモデルに変換
	categoryMasters := make([]categoryDomain.CategoryMaster, len(categoryMasterModels))
	for i, categoryMasterModel := range categoryMasterModels {
		categoryMaster, err := r.toDomain(&categoryMasterModel)
		if err != nil {
			return nil, err
		}
		categoryMasters[i] = *categoryMaster
	}

	return categoryMasters, nil
}

// FindByType タイプでカテゴリマスターを取得
func (r *CategoryMasterRepository) FindByType(ctx context.Context, categoryType categoryValue.CategoryType) ([]categoryDomain.CategoryMaster, error) {
	var categoryMasterModels []model.CategoryMaster
	result := r.db.WithContext(ctx).
		Where("type = ?", categoryType.String()).
		Order("display_order ASC, created_at ASC").
		Find(&categoryMasterModels)

	if result.Error != nil {
		return nil, fmt.Errorf("カテゴリマスター一覧の取得に失敗しました: %w", result.Error)
	}

	// ドメインモデルに変換
	categoryMasters := make([]categoryDomain.CategoryMaster, len(categoryMasterModels))
	for i, categoryMasterModel := range categoryMasterModels {
		categoryMaster, err := r.toDomain(&categoryMasterModel)
		if err != nil {
			return nil, err
		}
		categoryMasters[i] = *categoryMaster
	}

	return categoryMasters, nil
}

// Update カテゴリマスターを更新
func (r *CategoryMasterRepository) Update(ctx context.Context, categoryMaster categoryDomain.CategoryMaster) error {
	categoryMasterModel := r.toModel(&categoryMaster)

	result := r.db.WithContext(ctx).Model(&model.CategoryMaster{}).
		Where("id = ?", categoryMaster.ID).
		Updates(categoryMasterModel)

	if result.Error != nil {
		return fmt.Errorf("カテゴリマスターの更新に失敗しました: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("更新対象のカテゴリマスターが見つかりません: %s", categoryMaster.ID)
	}

	return nil
}

// Delete カテゴリマスターを削除
func (r *CategoryMasterRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&model.CategoryMaster{}, "id = ?", id)
	if result.Error != nil {
		return fmt.Errorf("カテゴリマスターの削除に失敗しました: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("削除対象のカテゴリマスターが見つかりません: %s", id)
	}

	return nil
}

// toModel ドメインモデルからGORMモデルへの変換
func (r *CategoryMasterRepository) toModel(categoryMaster *categoryDomain.CategoryMaster) *model.CategoryMaster {
	var color *string
	if categoryMaster.Color() != nil {
		colorStr := categoryMaster.Color().String()
		color = &colorStr
	}

	icon := categoryMaster.Icon()
	var iconPtr *string
	if icon != "" {
		iconPtr = &icon
	}

	return &model.CategoryMaster{
		Base: model.Base{
			ID:        categoryMaster.ID,
			CreatedAt: categoryMaster.CreatedAt,
			UpdatedAt: categoryMaster.UpdatedAt,
		},
		Name:         categoryMaster.Name().String(),
		Type:         model.CategoryType(categoryMaster.Type().String()),
		Icon:         iconPtr,
		Color:        color,
		DisplayOrder: categoryMaster.DisplayOrder(),
	}
}

// toDomain GORMモデルからドメインモデルへの変換
func (r *CategoryMasterRepository) toDomain(categoryMasterModel *model.CategoryMaster) (*categoryDomain.CategoryMaster, error) {
	// カテゴリ名
	name, err := categoryValue.NewCategoryName(categoryMasterModel.Name)
	if err != nil {
		return nil, fmt.Errorf("カテゴリ名の作成に失敗しました: %w", err)
	}

	// カテゴリタイプ
	categoryType, err := categoryValue.NewCategoryType(string(categoryMasterModel.Type))
	if err != nil {
		return nil, fmt.Errorf("カテゴリタイプの作成に失敗しました: %w", err)
	}

	// カラー
	var color *commonValue.HexColor
	if categoryMasterModel.Color != nil && *categoryMasterModel.Color != "" {
		hexColor, err := commonValue.NewHexColor(*categoryMasterModel.Color)
		if err != nil {
			return nil, fmt.Errorf("カラーの作成に失敗しました: %w", err)
		}
		color = hexColor
	}

	// アイコン
	icon := ""
	if categoryMasterModel.Icon != nil {
		icon = *categoryMasterModel.Icon
	}

	// ドメインエンティティを作成
	categoryMaster, err := categoryDomain.NewCategoryMaster(
		name,
		categoryType,
		icon,
		color,
		categoryMasterModel.DisplayOrder,
	)
	if err != nil {
		return nil, fmt.Errorf("カテゴリマスターエンティティの作成に失敗しました: %w", err)
	}

	// BaseEntityを設定
	categoryMaster.BaseEntity = common.BaseEntity{
		ID:        categoryMasterModel.ID,
		CreatedAt: categoryMasterModel.CreatedAt,
		UpdatedAt: categoryMasterModel.UpdatedAt,
	}

	return &categoryMaster, nil
}
