// Package main provides the database seeding functionality for the FinanceTracker application.
// It initializes sample data for development and testing purposes.
package main

import (
	"log"

	"financetracker/pkg/config"
	"financetracker/pkg/logger"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Initialize logger
	logger := logger.New()
	defer func() {
		if err := logger.Sync(); err != nil {
			log.Printf("Failed to sync logger: %v", err)
		}
	}()

	// Load configuration
	cfg := config.Load()

	// Connect to database
	db, err := gorm.Open(postgres.Open(cfg.GetDSN()), &gorm.Config{})
	if err != nil {
		logger.Error("Failed to connect to database")
		return
	}

	logger.Info("Connected to database")

	// Seed data in order
	logger.Info("Seeding database...")

	// 1. Seed category masters
	if err := seedCategoryMasters(db); err != nil {
		logger.Error("Failed to seed category masters")
		return
	}

	// 2. Seed test users
	if err := seedTestUsers(db); err != nil {
		logger.Error("Failed to seed test users")
		return
	}

	// 3. Seed accounts
	if err := seedAccounts(db); err != nil {
		logger.Error("Failed to seed accounts")
		return
	}

	// 4. Seed categories (user-specific)
	if err := seedCategories(db); err != nil {
		logger.Error("Failed to seed categories")
		return
	}

	// 5. Seed transactions
	if err := seedTransactions(db); err != nil {
		logger.Error("Failed to seed transactions")
		return
	}

	// 6. Seed budgets
	if err := seedBudgets(db); err != nil {
		logger.Error("Failed to seed budgets")
		return
	}

	logger.Info("Database seeding completed successfully!")
}
