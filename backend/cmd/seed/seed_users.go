package main

import (
	"log"

	"financetracker/internal/infrastructure/gorm/model"

	"gorm.io/gorm"
)

func seedTestUsers(db *gorm.DB) error {
	users := []model.User{
		{
			Auth0ID: "auth0|test_user_1",
			Email:   "test1@example.com",
			Name:    "テストユーザー1",
		},
		{
			Auth0ID: "auth0|test_user_2",
			Email:   "test2@example.com",
			Name:    "テストユーザー2",
		},
	}

	for _, user := range users {
		var existing model.User
		err := db.Where("auth0_id = ?", user.Auth0ID).First(&existing).Error

		if err == nil {
			log.Printf("User already exists: %s (Auth0ID: %s)\n", user.Name, user.Auth0ID)
			continue
		}

		if err != gorm.ErrRecordNotFound {
			return err
		}

		if err := db.Create(&user).Error; err != nil {
			return err
		}

		log.Printf("Created user: %s (Auth0ID: %s)\n", user.Name, user.Auth0ID)
	}

	return nil
}