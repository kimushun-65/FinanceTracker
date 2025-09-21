package main

import (
	"financetracker/internal/infrastructure/gorm/model"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

func seedAccounts(db *gorm.DB) error {
	var users []model.User
	if err := db.Find(&users).Error; err != nil {
		return err
	}

	for i := range users {
		accounts := []model.Account{
			{
				UserID:   users[i].ID,
				Name:     "現金",
				Type:     model.AccountTypeCash,
				Balance:  decimal.NewFromInt(50000),
				Currency: "JPY",
				IsActive: true,
			},
			{
				UserID:   users[i].ID,
				Name:     "みずほ銀行",
				Type:     model.AccountTypeBank,
				Balance:  decimal.NewFromInt(1500000),
				Currency: "JPY",
				IsActive: true,
			},
			{
				UserID:   users[i].ID,
				Name:     "楽天カード",
				Type:     model.AccountTypeCreditCard,
				Balance:  decimal.NewFromInt(-80000),
				Currency: "JPY",
				IsActive: true,
			},
			{
				UserID:   users[i].ID,
				Name:     "SBI証券",
				Type:     model.AccountTypeInvestment,
				Balance:  decimal.NewFromInt(3000000),
				Currency: "JPY",
				IsActive: true,
			},
		}

		err := seedWithDuplicateCheck(
			db,
			accounts,
			func(item model.Account) *gorm.DB {
				return db.Where("user_id = ? AND name = ?", item.UserID, item.Name)
			},
			func(item model.Account) string {
				return "account: " + item.Name + " (UserID: " + item.UserID.String() + ")"
			},
		)
		if err != nil {
			return err
		}
	}

	return nil
}
