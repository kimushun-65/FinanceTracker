package main

import (
	"log"
	"time"

	"financetracker/internal/infrastructure/gorm/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func seedTestUsers(db *gorm.DB) error {
	now := time.Now()
	users := []model.User{
		{
			Base: model.Base{
				ID:        uuid.New(),
				CreatedAt: now,
				UpdatedAt: now,
			},
			Auth0ID:       "auth0|test_user_1",
			Email:         "test1@example.com",
			Name:          "テストユーザー1",
			EmailVerified: true,
		},
		{
			Base: model.Base{
				ID:        uuid.New(),
				CreatedAt: now,
				UpdatedAt: now,
			},
			Auth0ID:       "auth0|test_user_2",
			Email:         "test2@example.com",
			Name:          "テストユーザー2",
			EmailVerified: true,
		},
		{
			Base: model.Base{
				ID:        uuid.New(),
				CreatedAt: now,
				UpdatedAt: now,
			},
			Auth0ID:       "auth0|dev-test-user-123",
			Email:         "dev-test@example.com",
			Name:          "開発テストユーザー",
			EmailVerified: false, // 開発用ユーザーは未検証として設定
		},
	}

	for i := range users {
		var existing model.User
		err := db.Where("auth0_id = ?", users[i].Auth0ID).First(&existing).Error

		if err == nil {
			log.Printf("User already exists: %s (Auth0ID: %s)\n", users[i].Name, users[i].Auth0ID)
			continue
		}

		if err != gorm.ErrRecordNotFound {
			return err
		}

		if err := db.Create(&users[i]).Error; err != nil {
			return err
		}

		log.Printf("Created user: %s (Auth0ID: %s)\n", users[i].Name, users[i].Auth0ID)
	}

	return nil
}
