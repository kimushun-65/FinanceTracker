// Package model defines GORM models for the FinanceTracker application.
package model

import "github.com/google/uuid"

// CategoryType represents the type of category.
type CategoryType string

const (
	CategoryTypeIncome  CategoryType = "INCOME"
	CategoryTypeExpense CategoryType = "EXPENSE"
)

// CategoryMaster represents a master category template.
type CategoryMaster struct {
	Base
	Name         string       `gorm:"type:varchar(100);not null"`
	Type         CategoryType `gorm:"type:varchar(20);not null;index"`
	Icon         *string      `gorm:"type:varchar(50)"`
	Color        *string      `gorm:"type:varchar(7)"` // Hex color code
	DisplayOrder int          `gorm:"not null;default:0"`

	// Relations
	Categories []Category `gorm:"foreignKey:CategoryMasterID;constraint:OnDelete:RESTRICT"`
}

// TableName specifies the table name for CategoryMaster.
func (CategoryMaster) TableName() string {
	return "category_masters"
}

// Category represents a user-specific category.
type Category struct {
	Base
	UserID           uuid.UUID  `gorm:"type:uuid;not null;index"`
	CategoryMasterID uuid.UUID  `gorm:"type:uuid;not null;index"`
	CustomName       *string    `gorm:"type:varchar(100)"`
	IsActive         bool       `gorm:"not null;default:true"`

	// Relations
	User           User           `gorm:"foreignKey:UserID"`
	CategoryMaster CategoryMaster `gorm:"foreignKey:CategoryMasterID"`
	Transactions   []Transaction  `gorm:"foreignKey:CategoryID;constraint:OnDelete:RESTRICT"`
	Budgets        []Budget       `gorm:"foreignKey:CategoryID;constraint:OnDelete:RESTRICT"`
}

// TableName specifies the table name for Category.
func (Category) TableName() string {
	return "categories"
}