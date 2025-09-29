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

func seedTransactions(db *gorm.DB) error {
	var users []model.User
	if err := db.Find(&users).Error; err != nil {
		return err
	}

	// 対象ユーザー一覧を収集（追加で投入）
	// 優先: 環境変数(カンマ区切り可) -> Googleアカウント -> dev-test-user
	auth0IDs := []string{}
    if v := os.Getenv("SEED_TARGET_AUTH0_ID"); v != "" {
        // そのまま追加（カンマ区切りにも対応）
        auth0IDs = append(auth0IDs, splitAndTrim(v)...)
    }
	auth0IDs = append(auth0IDs,
		"google-oauth2|110905699660329788470",
		"auth0|dev-test-user-123",
	)

	// DBに存在する対象ユーザーを集める
	targets := []model.User{}
	seen := map[string]bool{}
	for _, id := range auth0IDs {
		var u model.User
		if err := db.Where("auth0_id = ?", id).First(&u).Error; err == nil {
			if !seen[u.ID.String()] {
				targets = append(targets, u)
				seen[u.ID.String()] = true
				log.Printf("Will seed transactions for: %s (%s)", u.Name, u.Auth0ID)
			}
		}
	}

	if len(targets) == 0 {
		if len(users) == 0 {
			return nil
		}
		// フォールバック: 最初のユーザーのみ
		targets = append(targets, users[0])
		log.Printf("No explicit targets found. Fallback to first user: %s", users[0].Name)
	}

	for _, user := range targets {
		// 口座を取得
		var accounts []model.Account
		if err := db.Where("user_id = ?", user.ID).Find(&accounts).Error; err != nil {
			return err
		}

		// カテゴリを取得
		var categories []model.Category
		if err := db.Where("user_id = ?", user.ID).Preload("CategoryMaster").Find(&categories).Error; err != nil {
			return err
		}

		// カテゴリを名前でマッピング
		categoryMap := make(map[string]*model.Category)
		for i := range categories {
			categoryMap[categories[i].CategoryMaster.Name] = &categories[i]
		}

		// 銀行口座を取得
		var bankAccount *model.Account
		for i := range accounts {
			if accounts[i].Type == model.AccountTypeBank {
				bankAccount = &accounts[i]
				break
			}
		}

		if bankAccount == nil {
			log.Printf("No bank account found for seeding transactions (user: %s)", user.Name)
			continue
		}

		// 現在の月の初日
		now := time.Now()
		startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

		transactions := []model.Transaction{
			// 収入
			{
				Base: model.Base{
					ID:        uuid.New(),
					CreatedAt: now,
					UpdatedAt: now,
				},
				UserID:      user.ID,
				AccountID:   bankAccount.ID,
				CategoryID:  categoryMap["給与"].ID,
				Amount:      decimal.NewFromInt(280000),
				Type:        model.TransactionTypeIncome,
				Date:        startOfMonth.AddDate(0, 0, 24), // 25日
				Description: stringPtr("4月分給与"),
			},
			// 支出
			{
				Base: model.Base{
					ID:        uuid.New(),
					CreatedAt: now,
					UpdatedAt: now,
				},
				UserID:      user.ID,
				AccountID:   bankAccount.ID,
				CategoryID:  categoryMap["住居費"].ID,
				Amount:      decimal.NewFromInt(85000),
				Type:        model.TransactionTypeExpense,
				Date:        startOfMonth.AddDate(0, 0, 26), // 27日
				Description: stringPtr("家賃"),
			},
			{
				Base: model.Base{
					ID:        uuid.New(),
					CreatedAt: now,
					UpdatedAt: now,
				},
				UserID:      user.ID,
				AccountID:   bankAccount.ID,
				CategoryID:  categoryMap["光熱費"].ID,
				Amount:      decimal.NewFromInt(15000),
				Type:        model.TransactionTypeExpense,
				Date:        startOfMonth.AddDate(0, 0, 14), // 15日
				Description: stringPtr("電気・ガス・水道"),
			},
			{
				Base: model.Base{
					ID:        uuid.New(),
					CreatedAt: now,
					UpdatedAt: now,
				},
				UserID:      user.ID,
				AccountID:   bankAccount.ID,
				CategoryID:  categoryMap["通信費"].ID,
				Amount:      decimal.NewFromInt(12000),
				Type:        model.TransactionTypeExpense,
				Date:        startOfMonth.AddDate(0, 0, 9), // 10日
				Description: stringPtr("携帯・インターネット"),
			},
			{
				Base: model.Base{
					ID:        uuid.New(),
					CreatedAt: now,
					UpdatedAt: now,
				},
				UserID:      user.ID,
				AccountID:   bankAccount.ID,
				CategoryID:  categoryMap["食費"].ID,
				Amount:      decimal.NewFromInt(35000),
				Type:        model.TransactionTypeExpense,
				Date:        startOfMonth.AddDate(0, 0, 19), // 20日
				Description: stringPtr("スーパー・コンビニ"),
			},
			{
				Base: model.Base{
					ID:        uuid.New(),
					CreatedAt: now,
					UpdatedAt: now,
				},
				UserID:      user.ID,
				AccountID:   bankAccount.ID,
				CategoryID:  categoryMap["外食費"].ID,
				Amount:      decimal.NewFromInt(8500),
				Type:        model.TransactionTypeExpense,
				Date:        startOfMonth.AddDate(0, 0, 5), // 6日
				Description: stringPtr("レストラン"),
			},
			{
				Base: model.Base{
					ID:        uuid.New(),
					CreatedAt: now,
					UpdatedAt: now,
				},
				UserID:      user.ID,
				AccountID:   bankAccount.ID,
				CategoryID:  categoryMap["交通費"].ID,
				Amount:      decimal.NewFromInt(5200),
				Type:        model.TransactionTypeExpense,
				Date:        startOfMonth.AddDate(0, 0, 2), // 3日
				Description: stringPtr("電車・バス定期"),
			},
			{
				Base: model.Base{
					ID:        uuid.New(),
					CreatedAt: now,
					UpdatedAt: now,
				},
				UserID:      user.ID,
				AccountID:   bankAccount.ID,
				CategoryID:  categoryMap["娯楽費"].ID,
				Amount:      decimal.NewFromInt(15000),
				Type:        model.TransactionTypeExpense,
				Date:        startOfMonth.AddDate(0, 0, 12), // 13日
				Description: stringPtr("映画・ゲーム"),
			},
			{
				Base: model.Base{
					ID:        uuid.New(),
					CreatedAt: now,
					UpdatedAt: now,
				},
				UserID:      user.ID,
				AccountID:   bankAccount.ID,
				CategoryID:  categoryMap["日用品"].ID,
				Amount:      decimal.NewFromInt(7800),
				Type:        model.TransactionTypeExpense,
				Date:        startOfMonth.AddDate(0, 0, 7), // 8日
				Description: stringPtr("ドラッグストア"),
			},
			{
				Base: model.Base{
					ID:        uuid.New(),
					CreatedAt: now,
					UpdatedAt: now,
				},
				UserID:      user.ID,
				AccountID:   bankAccount.ID,
				CategoryID:  categoryMap["衣服費"].ID,
				Amount:      decimal.NewFromInt(12000),
				Type:        model.TransactionTypeExpense,
				Date:        startOfMonth.AddDate(0, 0, 15), // 16日
				Description: stringPtr("春物衣類"),
			},
		}

		for i := range transactions {
			// 同じ日付・金額・カテゴリのトランザクションがあるか確認
			var existing model.Transaction
			err := db.Where("user_id = ? AND date = ? AND amount = ? AND category_id = ?",
				transactions[i].UserID, transactions[i].Date, transactions[i].Amount, transactions[i].CategoryID).First(&existing).Error

			if err == nil {
				log.Printf("Transaction already exists for %s: %v on %s", user.Name, transactions[i].Amount, transactions[i].Date.Format("2006-01-02"))
				continue
			}

			if err != gorm.ErrRecordNotFound {
				return err
			}

			if err := db.Create(&transactions[i]).Error; err != nil {
				return err
			}

			log.Printf("Created transaction for %s: %v on %s (%s)", user.Name, transactions[i].Amount, transactions[i].Date.Format("2006-01-02"), *transactions[i].Description)
		}
	}

	return nil
}
