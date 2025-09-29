package main

import (
	"log"
	"os"
	"time"

	"financetracker/internal/infrastructure/gorm/model"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

func seedBudgets(db *gorm.DB) error {
	var users []model.User
	if err := db.Find(&users).Error; err != nil {
		return err
	}

	if len(users) == 0 {
		return nil
	}

	// 対象ユーザー一覧を収集（追加で投入）
	auth0IDs := splitAndTrim(os.Getenv("SEED_TARGET_AUTH0_ID"))
	auth0IDs = append(auth0IDs,
		"google-oauth2|110905699660329788470",
		"auth0|dev-test-user-123",
	)

	// 実在するユーザーを抽出
	targets := []model.User{}
	seen := map[string]bool{}
	for _, id := range auth0IDs {
		for i := range users {
			if users[i].Auth0ID == id {
				if !seen[users[i].ID.String()] {
					targets = append(targets, users[i])
					seen[users[i].ID.String()] = true
					log.Printf("Will seed budgets for: %s (%s)", users[i].Name, users[i].Auth0ID)
				}
			}
		}
	}

	if len(targets) == 0 {
		targets = append(targets, users[0])
		log.Printf("No explicit targets found. Fallback to first user: %s", users[0].Name)
	}

	// 支出カテゴリのみの予算を作成
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	budgetData := map[string]int64{
		"食費":  40000,
		"外食費": 15000,
		"住居費": 85000,
		"光熱費": 20000,
		"通信費": 15000,
		"交通費": 10000,
		"娯楽費": 20000,
		"日用品": 10000,
		"衣服費": 15000,
		"美容費": 5000,
		"医療費": 10000,
		"教育費": 30000,
		"保険料": 25000,
		"貯金":  50000,
		"投資":  30000,
	}

	for _, target := range targets {
		// カテゴリを取得
		var categories []model.Category
		if err := db.Where("user_id = ?", target.ID).Preload("CategoryMaster").Find(&categories).Error; err != nil {
			return err
		}

		for i := range categories {
			if categories[i].CategoryMaster.Type != model.CategoryTypeExpense {
				continue
			}

			amount, hasBudget := budgetData[categories[i].CategoryMaster.Name]
			if !hasBudget {
				continue
			}

			budget := model.Budget{
				Base: model.Base{
					ID:        uuid.New(),
					CreatedAt: now,
					UpdatedAt: now,
				},
				UserID:     target.ID,
				CategoryID: categories[i].ID,
				Amount:     decimal.NewFromInt(amount),
				PeriodType: model.PeriodTypeMonthly,
				StartDate:  startOfMonth,
				IsActive:   true,
			}

			var existing model.Budget
			err := db.Where("user_id = ? AND category_id = ? AND is_active = ?",
				budget.UserID, budget.CategoryID, true).First(&existing).Error

			if err == nil {
				log.Printf("Active budget already exists for user %s, category: %s", target.Name, categories[i].CategoryMaster.Name)
				continue
			}

			if err != gorm.ErrRecordNotFound {
				return err
			}

			if err := db.Create(&budget).Error; err != nil {
				return err
			}

			log.Printf("Created budget for user %s, category: %s, amount: %d", target.Name, categories[i].CategoryMaster.Name, amount)
		}
	}

	return nil
}
