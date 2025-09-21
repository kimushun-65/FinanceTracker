package main

import (
	"log"

	"financetracker/internal/infrastructure/gorm/model"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

func seedAccounts(db *gorm.DB) error {
	var users []model.User
	if err := db.Find(&users).Error; err != nil {
		return err
	}

	for _, user := range users {
		accounts := []model.Account{
			{
				UserID:   user.ID,
				Name:     "現金",
				Type:     model.AccountTypeCash,
				Balance:  decimal.NewFromInt(50000),
				Currency: "JPY",
				IsActive: true,
			},
			{
				UserID:   user.ID,
				Name:     "みずほ銀行",
				Type:     model.AccountTypeBank,
				Balance:  decimal.NewFromInt(1500000),
				Currency: "JPY",
				IsActive: true,
			},
			{
				UserID:   user.ID,
				Name:     "楽天カード",
				Type:     model.AccountTypeCreditCard,
				Balance:  decimal.NewFromInt(-80000),
				Currency: "JPY",
				IsActive: true,
			},
			{
				UserID:   user.ID,
				Name:     "SBI証券",
				Type:     model.AccountTypeInvestment,
				Balance:  decimal.NewFromInt(3000000),
				Currency: "JPY",
				IsActive: true,
			},
		}

		for _, account := range accounts {
			var existing model.Account
			err := db.Where("user_id = ? AND name = ?", account.UserID, account.Name).First(&existing).Error

			if err == nil {
				log.Printf("Account already exists: %s (UserID: %s)\n", account.Name, account.UserID)
				continue
			}

			if err != gorm.ErrRecordNotFound {
				return err
			}

			if err := db.Create(&account).Error; err != nil {
				return err
			}

			log.Printf("Created account: %s (UserID: %s)\n", account.Name, account.UserID)
		}
	}

	return nil
}