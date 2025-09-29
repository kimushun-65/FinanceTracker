package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"financetracker/internal/application/dto"
	categoryDomain "financetracker/internal/domain/category/entity"
	categoryRepo "financetracker/internal/domain/category/repository"
	categoryValue "financetracker/internal/domain/category/value"
	userValue "financetracker/internal/domain/user/value"
	"financetracker/pkg/errors"
	"financetracker/pkg/logger"
)

// CategoryService カテゴリー管理サービス
type CategoryService struct {
	categoryRepo       categoryRepo.CategoryRepository
	categoryMasterRepo categoryRepo.CategoryMasterRepository
	logger             *logger.Logger
}

// NewCategoryService 新しいCategoryServiceを作成
func NewCategoryService(
	categoryRepo categoryRepo.CategoryRepository,
	categoryMasterRepo categoryRepo.CategoryMasterRepository,
	logger *logger.Logger,
) *CategoryService {
	return &CategoryService{
		categoryRepo:       categoryRepo,
		categoryMasterRepo: categoryMasterRepo,
		logger:             logger,
	}
}

// CreateCategory ユーザーカテゴリーを作成
func (s *CategoryService) CreateCategory(ctx context.Context, userID uuid.UUID, req *dto.CreateCategoryRequest) (*dto.CategoryResponse, error) {
	// ユーザーIDを作成
	domainUserID := userValue.NewUserID(userID)

	// カテゴリーマスターの存在確認
	categoryMaster, err := s.categoryMasterRepo.FindByID(ctx, req.CategoryMasterID)
	if err != nil {
		s.logger.Error("カテゴリーマスター取得エラー",
			zap.Error(err),
			zap.String("categoryMasterID", req.CategoryMasterID.String()))
		return nil, errors.NewInternalError("カテゴリーマスター情報の取得に失敗しました", err)
	}

	if categoryMaster == nil {
		return nil, errors.NewNotFoundError(fmt.Sprintf("カテゴリーマスターが見つかりません: %s", req.CategoryMasterID))
	}

	// カスタム名の作成
	var customName *categoryValue.CategoryName
	if req.CustomName != nil {
		name, err := categoryValue.NewCategoryName(*req.CustomName)
		if err != nil {
			s.logger.Error("カスタム名作成エラー",
				zap.Error(err),
				zap.String("customName", *req.CustomName))
			return nil, errors.NewValidationError("無効なカスタム名です")
		}
		customName = &name
	}

	// 既存カテゴリーの重複チェック
	existingCategory, err := s.categoryRepo.FindByUserIDAndCategoryMasterID(ctx, domainUserID.Value(), req.CategoryMasterID)
	if err != nil {
		s.logger.Error("既存カテゴリー確認エラー",
			zap.Error(err),
			zap.String("userID", userID.String()),
			zap.String("categoryMasterID", req.CategoryMasterID.String()))
		return nil, errors.NewInternalError("カテゴリーの確認に失敗しました", err)
	}

	if existingCategory != nil {
		// 既存カテゴリーが存在する場合
		if existingCategory.IsActive() {
			return nil, errors.NewValidationError("このカテゴリーは既に作成済みです")
		}
		// 無効化されている場合は再有効化
		existingCategory.Activate()
		if customName != nil {
			existingCategory.UpdateCustomName(customName)
		}

		if err := s.categoryRepo.Save(ctx, *existingCategory); err != nil {
			s.logger.Error("カテゴリー更新エラー", zap.Error(err))
			return nil, errors.NewInternalError("カテゴリーの更新に失敗しました", err)
		}

		categoryDTO := dto.CategoryFromDomain(existingCategory)
		if categoryDTO != nil {
			categoryDTO.Name = categoryMaster.Name().String()
			categoryDTO.Icon = categoryMaster.Icon()
			categoryDTO.CategoryType = categoryMaster.Type().String()
			categoryDTO.DisplayOrder = categoryMaster.DisplayOrder()
			if categoryMaster.Color() != nil {
				color := categoryMaster.Color().String()
				categoryDTO.Color = &color
			}
		}
		return categoryDTO, nil
	}

	// 新規カテゴリーを作成
	category, err := categoryDomain.NewCategory(domainUserID.Value(), req.CategoryMasterID, customName)
	if err != nil {
		s.logger.Error("カテゴリー作成エラー", zap.Error(err))
		return nil, errors.NewValidationError("カテゴリーの作成に失敗しました")
	}

	// リポジトリに保存
	if err := s.categoryRepo.Save(ctx, category); err != nil {
		s.logger.Error("カテゴリー保存エラー", zap.Error(err))
		return nil, errors.NewInternalError("カテゴリーの保存に失敗しました", err)
	}

	// DTOに変換して返却
	categoryDTO := dto.CategoryFromDomain(&category)
	if categoryDTO != nil {
		categoryDTO.Name = categoryMaster.Name().String()
		categoryDTO.Icon = categoryMaster.Icon()
		categoryDTO.CategoryType = categoryMaster.Type().String()
		categoryDTO.DisplayOrder = categoryMaster.DisplayOrder()
		if categoryMaster.Color() != nil {
			color := categoryMaster.Color().String()
			categoryDTO.Color = &color
		}
	}
	return categoryDTO, nil
}

// GetCategory カテゴリー情報を取得
func (s *CategoryService) GetCategory(ctx context.Context, userID, categoryID uuid.UUID) (*dto.CategoryResponse, error) {
	// ユーザーIDを作成
	domainUserID := userValue.NewUserID(userID)

	// リポジトリからカテゴリーを取得
	category, err := s.categoryRepo.FindByID(ctx, categoryID)
	if err != nil {
		s.logger.Error("カテゴリー取得エラー",
			zap.Error(err),
			zap.String("categoryID", categoryID.String()))
		return nil, errors.NewInternalError("カテゴリー情報の取得に失敗しました", err)
	}

	if category == nil {
		return nil, errors.NewNotFoundError(fmt.Sprintf("カテゴリーが見つかりません: %s", categoryID))
	}

	// ユーザーのカテゴリーであることを確認
	if category.UserID() != domainUserID.Value() {
		return nil, errors.NewForbiddenError("このカテゴリーへのアクセス権限がありません")
	}

	// カテゴリマスターの情報を取得
	categoryMaster, err := s.categoryMasterRepo.FindByID(ctx, category.CategoryMasterID())
	if err != nil {
		s.logger.Error("カテゴリマスター取得エラー",
			zap.Error(err),
			zap.String("masterID", category.CategoryMasterID().String()))
		return nil, errors.NewInternalError("カテゴリマスター情報の取得に失敗しました", err)
	}

	// DTOに変換
	categoryDTO := dto.CategoryFromDomain(category)
	if categoryDTO != nil && categoryMaster != nil {
		categoryDTO.Name = categoryMaster.Name().String()
		categoryDTO.Icon = categoryMaster.Icon()
		categoryDTO.CategoryType = categoryMaster.Type().String()
		categoryDTO.DisplayOrder = categoryMaster.DisplayOrder()
		if categoryMaster.Color() != nil {
			color := categoryMaster.Color().String()
			categoryDTO.Color = &color
		}
	}
	return categoryDTO, nil
}

// GetCategoriesByUser ユーザーのカテゴリー一覧を取得
func (s *CategoryService) GetCategoriesByUser(ctx context.Context, userID uuid.UUID, params *dto.CategorySearchParams) (*dto.CategoryListResponse, error) {
	// ユーザーIDを作成
	domainUserID := userValue.NewUserID(userID)

	// リポジトリからカテゴリー一覧を取得
	categories, err := s.categoryRepo.FindByUserID(ctx, domainUserID.Value())
	if err != nil {
		s.logger.Error("カテゴリー一覧取得エラー",
			zap.Error(err),
			zap.String("userID", userID.String()))
		return nil, errors.NewInternalError("カテゴリー一覧の取得に失敗しました", err)
	}

	// フィルタリング
	var filteredCategories []*categoryDomain.Category
	for i := range categories {
		category := &categories[i]
		// アクティブフィルター
		if params.IsActive != nil && category.IsActive() != *params.IsActive {
			continue
		}
		// カテゴリーマスターIDフィルター
		if params.CategoryMasterID != nil && category.CategoryMasterID() != *params.CategoryMasterID {
			continue
		}
		filteredCategories = append(filteredCategories, category)
	}

	// カテゴリマスターを取得してマップを作成
	masterIDs := make([]uuid.UUID, 0, len(filteredCategories))
	masterIDMap := make(map[uuid.UUID]bool)
	for _, cat := range filteredCategories {
		if !masterIDMap[cat.CategoryMasterID()] {
			masterIDs = append(masterIDs, cat.CategoryMasterID())
			masterIDMap[cat.CategoryMasterID()] = true
		}
	}

	// カテゴリマスターを一括取得
	categoryMasters := make(map[uuid.UUID]*categoryDomain.CategoryMaster)
	for _, masterID := range masterIDs {
		master, err := s.categoryMasterRepo.FindByID(ctx, masterID)
		if err != nil {
			s.logger.Error("カテゴリマスター取得エラー",
				zap.Error(err),
				zap.String("masterID", masterID.String()))
			continue
		}
		if master != nil {
			categoryMasters[masterID] = master
		}
	}

	// DTOに変換（カテゴリマスター情報を含む）
	categoryResponses := make([]dto.CategoryResponse, len(filteredCategories))
	for i, category := range filteredCategories {
		categoryDTO := dto.CategoryFromDomain(category)
		if categoryDTO != nil {
			// カテゴリマスターの情報を設定
			if master, ok := categoryMasters[category.CategoryMasterID()]; ok {
				categoryDTO.Name = master.Name().String()
				categoryDTO.Icon = master.Icon()
				categoryDTO.CategoryType = master.Type().String()
				categoryDTO.DisplayOrder = master.DisplayOrder()
				if master.Color() != nil {
					color := master.Color().String()
					categoryDTO.Color = &color
				}
			}
			categoryResponses[i] = *categoryDTO
		}
	}

	return &dto.CategoryListResponse{
		Categories: categoryResponses,
		TotalCount: int64(len(categoryResponses)),
	}, nil
}

// UpdateCategory カテゴリー情報を更新
func (s *CategoryService) UpdateCategory(ctx context.Context, userID, categoryID uuid.UUID, req *dto.UpdateCategoryRequest) (*dto.CategoryResponse, error) {
	// ユーザーIDを作成
	domainUserID := userValue.NewUserID(userID)

	// リポジトリからカテゴリーを取得
	category, err := s.categoryRepo.FindByID(ctx, categoryID)
	if err != nil {
		s.logger.Error("カテゴリー取得エラー",
			zap.Error(err),
			zap.String("categoryID", categoryID.String()))
		return nil, errors.NewInternalError("カテゴリー情報の取得に失敗しました", err)
	}

	if category == nil {
		return nil, errors.NewNotFoundError(fmt.Sprintf("カテゴリーが見つかりません: %s", categoryID))
	}

	// ユーザーのカテゴリーであることを確認
	if category.UserID() != domainUserID.Value() {
		return nil, errors.NewForbiddenError("このカテゴリーへのアクセス権限がありません")
	}

	// 更新フィールドの適用
	if req.CustomName != nil {
		var customName *categoryValue.CategoryName
		if *req.CustomName != "" {
			name, err := categoryValue.NewCategoryName(*req.CustomName)
			if err != nil {
				s.logger.Error("カスタム名作成エラー",
					zap.Error(err),
					zap.String("customName", *req.CustomName))
				return nil, errors.NewValidationError("無効なカスタム名です")
			}
			customName = &name
		}
		category.UpdateCustomName(customName)
	}

	if req.IsActive != nil {
		if *req.IsActive {
			category.Activate()
		} else {
			category.Deactivate()
		}
	}

	// リポジトリで更新
	if err := s.categoryRepo.Save(ctx, *category); err != nil {
		s.logger.Error("カテゴリー更新エラー",
			zap.Error(err),
			zap.String("categoryID", categoryID.String()))
		return nil, errors.NewInternalError("カテゴリー情報の更新に失敗しました", err)
	}

	// カテゴリマスターの情報を取得
	categoryMaster, err := s.categoryMasterRepo.FindByID(ctx, category.CategoryMasterID())
	if err != nil {
		s.logger.Error("カテゴリマスター取得エラー",
			zap.Error(err),
			zap.String("masterID", category.CategoryMasterID().String()))
		return nil, errors.NewInternalError("カテゴリマスター情報の取得に失敗しました", err)
	}

	// DTOに変換して返却
	categoryDTO := dto.CategoryFromDomain(category)
	if categoryDTO != nil && categoryMaster != nil {
		categoryDTO.Name = categoryMaster.Name().String()
		categoryDTO.Icon = categoryMaster.Icon()
		categoryDTO.CategoryType = categoryMaster.Type().String()
		categoryDTO.DisplayOrder = categoryMaster.DisplayOrder()
		if categoryMaster.Color() != nil {
			color := categoryMaster.Color().String()
			categoryDTO.Color = &color
		}
	}
	return categoryDTO, nil
}

// DeleteCategory カテゴリーを削除（論理削除：無効化）
func (s *CategoryService) DeleteCategory(ctx context.Context, userID, categoryID uuid.UUID) error {
	// ユーザーIDを作成
	domainUserID := userValue.NewUserID(userID)

	// リポジトリからカテゴリーを取得
	category, err := s.categoryRepo.FindByID(ctx, categoryID)
	if err != nil {
		s.logger.Error("カテゴリー取得エラー",
			zap.Error(err),
			zap.String("categoryID", categoryID.String()))
		return errors.NewInternalError("カテゴリー情報の取得に失敗しました", err)
	}

	if category == nil {
		return errors.NewNotFoundError(fmt.Sprintf("カテゴリーが見つかりません: %s", categoryID))
	}

	// ユーザーのカテゴリーであることを確認
	if category.UserID() != domainUserID.Value() {
		return errors.NewForbiddenError("このカテゴリーへのアクセス権限がありません")
	}

	// カテゴリーを無効化
	category.Deactivate()

	// リポジトリで更新
	if err := s.categoryRepo.Save(ctx, *category); err != nil {
		s.logger.Error("カテゴリー削除エラー",
			zap.Error(err),
			zap.String("categoryID", categoryID.String()))
		return errors.NewInternalError("カテゴリーの削除に失敗しました", err)
	}

	return nil
}

// GetCategoryMasters カテゴリーマスター一覧を取得
func (s *CategoryService) GetCategoryMasters(ctx context.Context, params *dto.CategoryMasterSearchParams) (*dto.CategoryMasterListResponse, error) {
	// リポジトリからカテゴリーマスター一覧を取得
	categoryMasters, err := s.categoryMasterRepo.FindAll(ctx)
	if err != nil {
		s.logger.Error("カテゴリーマスター一覧取得エラー", zap.Error(err))
		return nil, errors.NewInternalError("カテゴリーマスター一覧の取得に失敗しました", err)
	}

	// フィルタリング
	var filteredCategoryMasters []*categoryDomain.CategoryMaster
	for i := range categoryMasters {
		categoryMaster := &categoryMasters[i]
		// カテゴリータイプフィルター
		if params.CategoryType != nil && categoryMaster.Type().String() != *params.CategoryType {
			continue
		}
		filteredCategoryMasters = append(filteredCategoryMasters, categoryMaster)
	}

	// DTOに変換
	return &dto.CategoryMasterListResponse{
		CategoryMasters: dto.CategoryMastersFromDomain(filteredCategoryMasters),
		TotalCount:      int64(len(filteredCategoryMasters)),
	}, nil
}

// GetCategoryMaster カテゴリーマスター情報を取得
func (s *CategoryService) GetCategoryMaster(ctx context.Context, categoryMasterID uuid.UUID) (*dto.CategoryMasterResponse, error) {
	// リポジトリからカテゴリーマスターを取得
	categoryMaster, err := s.categoryMasterRepo.FindByID(ctx, categoryMasterID)
	if err != nil {
		s.logger.Error("カテゴリーマスター取得エラー",
			zap.Error(err),
			zap.String("categoryMasterID", categoryMasterID.String()))
		return nil, errors.NewInternalError("カテゴリーマスター情報の取得に失敗しました", err)
	}

	if categoryMaster == nil {
		return nil, errors.NewNotFoundError(fmt.Sprintf("カテゴリーマスターが見つかりません: %s", categoryMasterID))
	}

	// DTOに変換
	return dto.CategoryMasterFromDomain(categoryMaster), nil
}
