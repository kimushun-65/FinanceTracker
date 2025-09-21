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
	_, err := gorm.Open(postgres.Open(cfg.GetDSN()), &gorm.Config{})
	if err != nil {
		logger.Fatal("Failed to connect to database", err)
	}

	logger.Info("Connected to database")

	// TODO: Add seed data here
	// For now, just a placeholder
	logger.Info("Seeding database...")

	// Example seed data structure (to be implemented)
	// 1. Create default categories
	// 2. Create demo accounts
	// 3. Create sample transactions
	// 4. Create sample budgets

	logger.Info("Database seeding completed successfully!")
}
