// Package entity カテゴリ関連のエンティティを定義
package entity

import (
	"financetracker/internal/domain/category/value"
	"financetracker/internal/domain/common"
	commonValue "financetracker/internal/domain/common/value"
)

// CategoryMaster カテゴリマスターエンティティ
type CategoryMaster struct {
	common.BaseEntity
	name         value.CategoryName
	categoryType value.CategoryType
	icon         string
	color        *commonValue.HexColor
	displayOrder int
}

// NewCategoryMaster 新しいカテゴリマスターを作成
func NewCategoryMaster(
	name value.CategoryName,
	categoryType value.CategoryType,
	icon string,
	color *commonValue.HexColor,
	displayOrder int,
) (CategoryMaster, error) {
	if err := validateCategoryMasterParams(displayOrder); err != nil {
		return CategoryMaster{}, err
	}

	baseEntity := common.NewBaseEntity()

	return CategoryMaster{
		BaseEntity:   baseEntity,
		name:         name,
		categoryType: categoryType,
		icon:         icon,
		color:        color,
		displayOrder: displayOrder,
	}, nil
}

// validateCategoryMasterParams カテゴリマスターパラメータの妥当性を検証
func validateCategoryMasterParams(displayOrder int) error {
	if displayOrder < 0 {
		return common.NewDomainError(
			common.DomainErrorTypeInvalidValue,
			"表示順序は0以上である必要があります",
		)
	}

	return nil
}

// Name カテゴリ名を取得
func (c CategoryMaster) Name() value.CategoryName {
	return c.name
}

// Type カテゴリタイプを取得
func (c CategoryMaster) Type() value.CategoryType {
	return c.categoryType
}

// Icon アイコンを取得
func (c CategoryMaster) Icon() string {
	return c.icon
}

// Color カラーを取得
func (c CategoryMaster) Color() *commonValue.HexColor {
	return c.color
}

// DisplayOrder 表示順序を取得
func (c CategoryMaster) DisplayOrder() int {
	return c.displayOrder
}

// UpdateName カテゴリ名を更新
func (c *CategoryMaster) UpdateName(name value.CategoryName) {
	c.name = name
	c.UpdateTimestamp()
}

// UpdateIcon アイコンを更新
func (c *CategoryMaster) UpdateIcon(icon string) {
	c.icon = icon
	c.UpdateTimestamp()
}

// UpdateColor カラーを更新
func (c *CategoryMaster) UpdateColor(color *commonValue.HexColor) {
	c.color = color
	c.UpdateTimestamp()
}

// UpdateDisplayOrder 表示順序を更新
func (c *CategoryMaster) UpdateDisplayOrder(displayOrder int) error {
	if err := validateCategoryMasterParams(displayOrder); err != nil {
		return err
	}

	c.displayOrder = displayOrder
	c.UpdateTimestamp()
	return nil
}
