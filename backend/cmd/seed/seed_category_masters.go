package main

import (
	"financetracker/internal/infrastructure/gorm/model"

	"gorm.io/gorm"
)

func seedCategoryMasters(db *gorm.DB) error {
	categoryMasters := []model.CategoryMaster{
		// 収入カテゴリ
		{
			Name:         "給与",
			Type:         model.CategoryTypeIncome,
			Icon:         stringPtr("💰"),
			Color:        stringPtr("#4CAF50"),
			DisplayOrder: 1,
		},
		{
			Name:         "ボーナス",
			Type:         model.CategoryTypeIncome,
			Icon:         stringPtr("🎁"),
			Color:        stringPtr("#8BC34A"),
			DisplayOrder: 2,
		},
		{
			Name:         "副業",
			Type:         model.CategoryTypeIncome,
			Icon:         stringPtr("💼"),
			Color:        stringPtr("#66BB6A"),
			DisplayOrder: 3,
		},
		{
			Name:         "投資収入",
			Type:         model.CategoryTypeIncome,
			Icon:         stringPtr("📈"),
			Color:        stringPtr("#81C784"),
			DisplayOrder: 4,
		},
		{
			Name:         "その他収入",
			Type:         model.CategoryTypeIncome,
			Icon:         stringPtr("💵"),
			Color:        stringPtr("#A5D6A7"),
			DisplayOrder: 5,
		},
		// 支出カテゴリ
		{
			Name:         "食費",
			Type:         model.CategoryTypeExpense,
			Icon:         stringPtr("🍽️"),
			Color:        stringPtr("#F44336"),
			DisplayOrder: 6,
		},
		{
			Name:         "外食費",
			Type:         model.CategoryTypeExpense,
			Icon:         stringPtr("🍜"),
			Color:        stringPtr("#EF5350"),
			DisplayOrder: 7,
		},
		{
			Name:         "住居費",
			Type:         model.CategoryTypeExpense,
			Icon:         stringPtr("🏠"),
			Color:        stringPtr("#E91E63"),
			DisplayOrder: 8,
		},
		{
			Name:         "光熱費",
			Type:         model.CategoryTypeExpense,
			Icon:         stringPtr("💡"),
			Color:        stringPtr("#EC407A"),
			DisplayOrder: 9,
		},
		{
			Name:         "通信費",
			Type:         model.CategoryTypeExpense,
			Icon:         stringPtr("📱"),
			Color:        stringPtr("#AB47BC"),
			DisplayOrder: 10,
		},
		{
			Name:         "交通費",
			Type:         model.CategoryTypeExpense,
			Icon:         stringPtr("🚗"),
			Color:        stringPtr("#7E57C2"),
			DisplayOrder: 11,
		},
		{
			Name:         "医療費",
			Type:         model.CategoryTypeExpense,
			Icon:         stringPtr("🏥"),
			Color:        stringPtr("#5C6BC0"),
			DisplayOrder: 12,
		},
		{
			Name:         "教育費",
			Type:         model.CategoryTypeExpense,
			Icon:         stringPtr("📚"),
			Color:        stringPtr("#3F51B5"),
			DisplayOrder: 13,
		},
		{
			Name:         "娯楽費",
			Type:         model.CategoryTypeExpense,
			Icon:         stringPtr("🎮"),
			Color:        stringPtr("#2196F3"),
			DisplayOrder: 14,
		},
		{
			Name:         "衣服費",
			Type:         model.CategoryTypeExpense,
			Icon:         stringPtr("👔"),
			Color:        stringPtr("#03A9F4"),
			DisplayOrder: 15,
		},
		{
			Name:         "美容費",
			Type:         model.CategoryTypeExpense,
			Icon:         stringPtr("💄"),
			Color:        stringPtr("#00BCD4"),
			DisplayOrder: 16,
		},
		{
			Name:         "日用品",
			Type:         model.CategoryTypeExpense,
			Icon:         stringPtr("🧻"),
			Color:        stringPtr("#009688"),
			DisplayOrder: 17,
		},
		{
			Name:         "保険料",
			Type:         model.CategoryTypeExpense,
			Icon:         stringPtr("🛡️"),
			Color:        stringPtr("#795548"),
			DisplayOrder: 18,
		},
		{
			Name:         "税金",
			Type:         model.CategoryTypeExpense,
			Icon:         stringPtr("📋"),
			Color:        stringPtr("#607D8B"),
			DisplayOrder: 19,
		},
		{
			Name:         "貯金",
			Type:         model.CategoryTypeExpense,
			Icon:         stringPtr("🏦"),
			Color:        stringPtr("#FF9800"),
			DisplayOrder: 20,
		},
		{
			Name:         "投資",
			Type:         model.CategoryTypeExpense,
			Icon:         stringPtr("📊"),
			Color:        stringPtr("#FF5722"),
			DisplayOrder: 21,
		},
		{
			Name:         "その他支出",
			Type:         model.CategoryTypeExpense,
			Icon:         stringPtr("📦"),
			Color:        stringPtr("#9E9E9E"),
			DisplayOrder: 22,
		},
	}

	err := seedWithDuplicateCheck(
		db,
		categoryMasters,
		func(item model.CategoryMaster) *gorm.DB {
			return db.Where("name = ? AND type = ?", item.Name, item.Type)
		},
		func(item model.CategoryMaster) string {
			return "category master: " + item.Name + " (Type: " + string(item.Type) + ")"
		},
	)
	if err != nil {
		return err
	}

	return nil
}
