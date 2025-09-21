package main

import (
	"log"
	"time"

	"financetracker/internal/infrastructure/gorm/model"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

func seedTransactions(db *gorm.DB) error {
	var users []model.User
	if err := db.Find(&users).Error; err != nil {
		return err
	}

	// 最初のユーザーのみにサンプルトランザクションを作成
	if len(users) == 0 {
		return nil
	}

	user := users[0]

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

	// カテゴリをタイプ別にマッピング
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
		log.Println("No bank account found for seeding transactions")
		return nil
	}

	// 現在の月の初日
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	transactions := []model.Transaction{
		// 収入
		{
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
			UserID:      user.ID,
			AccountID:   bankAccount.ID,
			CategoryID:  categoryMap["住居費"].ID,
			Amount:      decimal.NewFromInt(85000),
			Type:        model.TransactionTypeExpense,
			Date:        startOfMonth.AddDate(0, 0, 26), // 27日
			Description: stringPtr("家賃"),
		},
		{
			UserID:      user.ID,
			AccountID:   bankAccount.ID,
			CategoryID:  categoryMap["光熱費"].ID,
			Amount:      decimal.NewFromInt(15000),
			Type:        model.TransactionTypeExpense,
			Date:        startOfMonth.AddDate(0, 0, 14), // 15日
			Description: stringPtr("電気・ガス・水道"),
		},
		{
			UserID:      user.ID,
			AccountID:   bankAccount.ID,
			CategoryID:  categoryMap["通信費"].ID,
			Amount:      decimal.NewFromInt(12000),
			Type:        model.TransactionTypeExpense,
			Date:        startOfMonth.AddDate(0, 0, 9), // 10日
			Description: stringPtr("携帯・インターネット"),
		},
		{
			UserID:      user.ID,
			AccountID:   bankAccount.ID,
			CategoryID:  categoryMap["食費"].ID,
			Amount:      decimal.NewFromInt(35000),
			Type:        model.TransactionTypeExpense,
			Date:        startOfMonth.AddDate(0, 0, 19), // 20日
			Description: stringPtr("スーパー・コンビニ"),
		},
		{
			UserID:      user.ID,
			AccountID:   bankAccount.ID,
			CategoryID:  categoryMap["外食費"].ID,
			Amount:      decimal.NewFromInt(8500),
			Type:        model.TransactionTypeExpense,
			Date:        startOfMonth.AddDate(0, 0, 5), // 6日
			Description: stringPtr("レストラン"),
		},
		{
			UserID:      user.ID,
			AccountID:   bankAccount.ID,
			CategoryID:  categoryMap["交通費"].ID,
			Amount:      decimal.NewFromInt(5200),
			Type:        model.TransactionTypeExpense,
			Date:        startOfMonth.AddDate(0, 0, 2), // 3日
			Description: stringPtr("電車・バス定期"),
		},
		{
			UserID:      user.ID,
			AccountID:   bankAccount.ID,
			CategoryID:  categoryMap["娯楽費"].ID,
			Amount:      decimal.NewFromInt(15000),
			Type:        model.TransactionTypeExpense,
			Date:        startOfMonth.AddDate(0, 0, 12), // 13日
			Description: stringPtr("映画・ゲーム"),
		},
		{
			UserID:      user.ID,
			AccountID:   bankAccount.ID,
			CategoryID:  categoryMap["日用品"].ID,
			Amount:      decimal.NewFromInt(7800),
			Type:        model.TransactionTypeExpense,
			Date:        startOfMonth.AddDate(0, 0, 7), // 8日
			Description: stringPtr("ドラッグストア"),
		},
		{
			UserID:      user.ID,
			AccountID:   bankAccount.ID,
			CategoryID:  categoryMap["衣服費"].ID,
			Amount:      decimal.NewFromInt(12000),
			Type:        model.TransactionTypeExpense,
			Date:        startOfMonth.AddDate(0, 0, 15), // 16日
			Description: stringPtr("春物衣類"),
		},
	}

	for _, transaction := range transactions {
		// 同じ日付・金額・カテゴリのトランザクションがあるか確認
		var existing model.Transaction
		err := db.Where("user_id = ? AND date = ? AND amount = ? AND category_id = ?",
			transaction.UserID, transaction.Date, transaction.Amount, transaction.CategoryID).First(&existing).Error

		if err == nil {
			log.Printf("Transaction already exists: %v on %s\n", transaction.Amount, transaction.Date.Format("2006-01-02"))
			continue
		}

		if err != gorm.ErrRecordNotFound {
			return err
		}

		if err := db.Create(&transaction).Error; err != nil {
			return err
		}

		log.Printf("Created transaction: %v on %s (%s)\n", transaction.Amount, transaction.Date.Format("2006-01-02"), *transaction.Description)
	}

	return nil
}