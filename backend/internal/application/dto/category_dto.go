package dto

import (
	"time"

	"github.com/google/uuid"

	categoryDomain "financetracker/internal/domain/category/entity"
)

// CategoryResponse カテゴリー情報レスポンス
type CategoryResponse struct {
	ID               uuid.UUID `json:"id"`
	UserID           uuid.UUID `json:"user_id"`
	CategoryMasterID uuid.UUID `json:"category_master_id"`
	CustomName       *string   `json:"custom_name,omitempty"`
	IsActive         bool      `json:"is_active"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// CreateCategoryRequest カテゴリー作成リクエスト
type CreateCategoryRequest struct {
	CategoryMasterID uuid.UUID `json:"category_master_id" binding:"required"`
	CustomName       *string   `json:"custom_name" binding:"omitempty,min=1,max=50"`
}

// UpdateCategoryRequest カテゴリー更新リクエスト
type UpdateCategoryRequest struct {
	CustomName *string `json:"custom_name" binding:"omitempty,min=1,max=50"`
	IsActive   *bool   `json:"is_active" binding:"omitempty"`
}

// CategoryMasterResponse カテゴリーマスター情報レスポンス
type CategoryMasterResponse struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	CategoryType string    `json:"category_type"`
	Icon         string    `json:"icon"`
	Color        *string   `json:"color,omitempty"`
	DisplayOrder int       `json:"display_order"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// CategoryListResponse カテゴリー一覧レスポンス
type CategoryListResponse struct {
	Categories []CategoryResponse `json:"categories"`
	TotalCount int64              `json:"total_count"`
}

// CategorySearchParams カテゴリー検索パラメータ
type CategorySearchParams struct {
	UserID           uuid.UUID  `form:"-"`
	CategoryMasterID *uuid.UUID `form:"category_master_id"`
	IsActive         *bool      `form:"is_active"`
	OrderBy          string     `form:"order_by,default=created_at desc"`
}

// CategoryMasterSearchParams カテゴリーマスター検索パラメータ
type CategoryMasterSearchParams struct {
	CategoryType *string `form:"category_type" binding:"omitempty,oneof=income expense"`
	OrderBy      string  `form:"order_by,default=display_order asc"`
}

// CategoryFromDomain ドメインエンティティからDTOへの変換
func CategoryFromDomain(category *categoryDomain.Category) *CategoryResponse {
	if category == nil {
		return nil
	}

	var customName *string
	if category.CustomName() != nil {
		name := category.CustomName().String()
		customName = &name
	}

	return &CategoryResponse{
		ID:               category.ID,
		UserID:           category.UserID(),
		CategoryMasterID: category.CategoryMasterID(),
		CustomName:       customName,
		IsActive:         category.IsActive(),
		CreatedAt:        category.CreatedAt,
		UpdatedAt:        category.UpdatedAt,
	}
}

// CategoryMasterFromDomain カテゴリーマスターエンティティからDTOへの変換
func CategoryMasterFromDomain(categoryMaster *categoryDomain.CategoryMaster) *CategoryMasterResponse {
	if categoryMaster == nil {
		return nil
	}

	var color *string
	if categoryMaster.Color() != nil {
		colorValue := categoryMaster.Color().String()
		color = &colorValue
	}

	return &CategoryMasterResponse{
		ID:           categoryMaster.ID,
		Name:         categoryMaster.Name().String(),
		CategoryType: categoryMaster.Type().String(),
		Icon:         categoryMaster.Icon(),
		Color:        color,
		DisplayOrder: categoryMaster.DisplayOrder(),
		CreatedAt:    categoryMaster.CreatedAt,
		UpdatedAt:    categoryMaster.UpdatedAt,
	}
}

// CategoriesFromDomain ドメインエンティティのスライスからDTOへの変換
func CategoriesFromDomain(categories []*categoryDomain.Category) []CategoryResponse {
	result := make([]CategoryResponse, len(categories))
	for i, category := range categories {
		if categoryDTO := CategoryFromDomain(category); categoryDTO != nil {
			result[i] = *categoryDTO
		}
	}
	return result
}

// CategoryMasterListResponse カテゴリーマスター一覧レスポンス
type CategoryMasterListResponse struct {
	CategoryMasters []CategoryMasterResponse `json:"category_masters"`
	TotalCount      int64                    `json:"total_count"`
}

// CategoryMastersFromDomain カテゴリーマスターエンティティのスライスからDTOへの変換
func CategoryMastersFromDomain(categoryMasters []*categoryDomain.CategoryMaster) []CategoryMasterResponse {
	result := make([]CategoryMasterResponse, len(categoryMasters))
	for i, categoryMaster := range categoryMasters {
		if categoryMasterDTO := CategoryMasterFromDomain(categoryMaster); categoryMasterDTO != nil {
			result[i] = *categoryMasterDTO
		}
	}
	return result
}
