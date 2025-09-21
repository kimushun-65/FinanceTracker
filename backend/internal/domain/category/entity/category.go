// Package entity カテゴリ関連のエンティティを定義
package entity

import (
	"financetracker/internal/domain/category/value"
	"financetracker/internal/domain/common"

	"github.com/google/uuid"
)

// Category ユーザーカテゴリエンティティ
type Category struct {
	common.BaseEntity
	userID           uuid.UUID
	categoryMasterID uuid.UUID
	customName       *value.CategoryName // カスタム名称（任意）
	isActive         bool
}

// NewCategory 新しいカテゴリを作成
func NewCategory(
	userID uuid.UUID,
	categoryMasterID uuid.UUID,
	customName *value.CategoryName,
) (Category, error) {
	if err := validateCategoryParams(userID, categoryMasterID); err != nil {
		return Category{}, err
	}

	baseEntity := common.NewBaseEntity()

	return Category{
		BaseEntity:       baseEntity,
		userID:           userID,
		categoryMasterID: categoryMasterID,
		customName:       customName,
		isActive:         true, // 新規作成時はアクティブ
	}, nil
}

// validateCategoryParams カテゴリパラメータの妥当性を検証
func validateCategoryParams(userID, categoryMasterID uuid.UUID) error {
	if userID == uuid.Nil {
		return common.NewDomainError(
			common.DomainErrorTypeInvalidValue,
			"ユーザーIDが必要です",
		)
	}

	if categoryMasterID == uuid.Nil {
		return common.NewDomainError(
			common.DomainErrorTypeInvalidValue,
			"カテゴリマスターIDが必要です",
		)
	}

	return nil
}

// UserID ユーザーIDを取得
func (c Category) UserID() uuid.UUID {
	return c.userID
}

// CategoryMasterID カテゴリマスターIDを取得
func (c Category) CategoryMasterID() uuid.UUID {
	return c.categoryMasterID
}

// CustomName カスタム名称を取得
func (c Category) CustomName() *value.CategoryName {
	return c.customName
}

// IsActive アクティブかどうかを判定
func (c Category) IsActive() bool {
	return c.isActive
}

// UpdateCustomName カスタム名称を更新
func (c *Category) UpdateCustomName(customName *value.CategoryName) {
	c.customName = customName
	c.UpdateTimestamp()
}

// Activate カテゴリを有効化
func (c *Category) Activate() {
	if !c.isActive {
		c.isActive = true
		c.UpdateTimestamp()
	}
}

// Deactivate カテゴリを無効化（使用中のカテゴリは削除不可のため無効化のみ）
func (c *Category) Deactivate() {
	if c.isActive {
		c.isActive = false
		c.UpdateTimestamp()
	}
}

// HasCustomName カスタム名称を持っているかチェック
func (c Category) HasCustomName() bool {
	return c.customName != nil
}
