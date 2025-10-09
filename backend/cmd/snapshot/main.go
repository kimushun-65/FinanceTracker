// Package main is the entry point for the asset snapshot batch job.
// This command creates daily asset snapshots for all users and exits.
// It is designed to be run as an ECS Scheduled Task.
package main

import (
	"context"
	stdlog "log"
	"os"

	"financetracker/internal/di"
	"financetracker/pkg/logger"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		stdlog.Println("No .env file found")
	}

	// Initialize logger
	log := logger.NewLogger("snapshot-batch")
	defer func() {
		if err := log.Sync(); err != nil {
			stdlog.Printf("Failed to sync logger: %v", err)
		}
	}()

	log.Info("Starting asset snapshot batch job")

	// Initialize DI container
	container, err := di.NewContainer()
	if err != nil {
		log.Error("Failed to initialize DI container", zap.Error(err))
		os.Exit(1)
	}
	defer func() {
		if err := container.Close(); err != nil {
			log.Error("Failed to close DI container", zap.Error(err))
		}
	}()

	ctx := context.Background()

	// Get all users
	users, err := container.UserRepo.FindAll(ctx)
	if err != nil {
		log.Error("Failed to get users for snapshot creation", zap.Error(err))
		os.Exit(1)
	}

	log.Info("Found users for snapshot creation", zap.Int("user_count", len(users)))

	successCount := 0
	errorCount := 0

	// Create snapshot for each user
	for _, user := range users {
		err := container.AssetService.CreateDailySnapshot(ctx, user.GetID())
		if err != nil {
			log.Error("Failed to create snapshot for user",
				zap.String("user_id", user.GetID().String()),
				zap.Error(err))
			errorCount++
		} else {
			log.Info("Successfully created snapshot for user",
				zap.String("user_id", user.GetID().String()))
			successCount++
		}
	}

	log.Info("Asset snapshot batch job completed",
		zap.Int("success", successCount),
		zap.Int("errors", errorCount),
		zap.Int("total", len(users)))

	// Exit with error code if any snapshots failed
	if errorCount > 0 {
		os.Exit(1)
	}
}
