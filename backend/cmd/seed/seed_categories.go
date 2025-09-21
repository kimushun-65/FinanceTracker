package main

import (
	"log"

	"financetracker/internal/infrastructure/gorm/model"

	"gorm.io/gorm"
)

func seedCategories(db *gorm.DB) error {
	var users []model.User
	if err := db.Find(&users).Error; err != nil {
		return err
	}

	var categoryMasters []model.CategoryMaster
	if err := db.Find(&categoryMasters).Error; err != nil {
		return err
	}

	for _, user := range users {
		for _, cm := range categoryMasters {
			category := model.Category{
				UserID:           user.ID,
				CategoryMasterID: cm.ID,
				IsActive:         true,
			}

			var existing model.Category
			err := db.Where("user_id = ? AND category_master_id = ?", category.UserID, category.CategoryMasterID).First(&existing).Error

			if err == nil {
				log.Printf("Category already exists for user: %s, master: %s\n", user.ID, cm.Name)
				continue
			}

			if err != gorm.ErrRecordNotFound {
				return err
			}

			if err := db.Create(&category).Error; err != nil {
				return err
			}

			log.Printf("Created category for user: %s, master: %s\n", user.ID, cm.Name)
		}
	}

	return nil
}