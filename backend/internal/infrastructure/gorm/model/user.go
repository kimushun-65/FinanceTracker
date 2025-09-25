// Package model defines GORM models for the FinanceTracker application.
package model

// User represents a user in the system.
type User struct {
	Base
	Auth0ID string `gorm:"column:auth0_id;type:varchar(255);not null;uniqueIndex"`
	Email   string `gorm:"type:varchar(255);not null;uniqueIndex"`
	Name    string `gorm:"type:varchar(255);not null"`

	// Relations
	Accounts             []Account             `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Categories           []Category            `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Transactions         []Transaction         `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Budgets              []Budget              `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	BudgetSuggestions    []BudgetSuggestion    `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	AssetSnapshots       []AssetSnapshot       `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	AssetForecasts       []AssetForecast       `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	NotificationSettings []NotificationSetting `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

// TableName specifies the table name for User.
func (User) TableName() string {
	return "users"
}
