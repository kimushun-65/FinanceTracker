// Package main provides the database migration command-line tool for the FinanceTracker application.
// It supports various migration operations like apply, check, diff, auto, and validate.
package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"financetracker/internal/infrastructure/gorm/model"
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

	// Get command
	if len(os.Args) < 2 {
		printUsage(cfg.AppEnv)
		return
	}

	command := os.Args[1]

	if err := runCommand(command, cfg, logger); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}

// printUsage prints the usage message based on the environment
func printUsage(appEnv string) {
	if appEnv == "development" {
		fmt.Println("Usage: migrate [apply|check|diff|validate|dev-setup]")
	} else {
		fmt.Println("Usage: migrate [apply|check|diff|validate]")
	}
	os.Exit(1)
}

// runCommand executes the specified command
func runCommand(command string, cfg *config.Config, logger *logger.Logger) error {
	switch command {
	case "apply":
		return runApplyCommand(logger)
	case "check":
		return runCheckCommand(logger)
	case "diff":
		return runDiffCommand(logger)
	case "validate":
		return runValidateCommand(logger)
	case "dev-setup":
		return runDevSetupCommand(cfg, logger)
	default:
		return fmt.Errorf("unknown command: %s", command)
	}
}

// runApplyCommand applies Atlas migrations
func runApplyCommand(logger *logger.Logger) error {
	logger.Info("Applying Atlas migrations...")

	// First, try to apply migrations normally
	err := runAtlasCommand("migrate", "apply",
		"--config", "file:///app/atlas.hcl",
		"--env", "dev",
		"--dir", "file:///app/cmd/migrate/migrations")

	if err != nil {
		// If it fails due to dirty database, try with baseline
		logger.Info("Initial migration failed, trying with baseline...")
		if err := runAtlasCommand("migrate", "apply",
			"--config", "file:///app/atlas.hcl",
			"--env", "dev",
			"--dir", "file:///app/cmd/migrate/migrations",
			"--baseline", "20240101000000"); err != nil {
			return fmt.Errorf("failed to apply migrations with baseline: %w", err)
		}
	}
	logger.Info("Migrations applied successfully")
	return nil
}

// runCheckCommand checks schema differences
func runCheckCommand(logger *logger.Logger) error {
	logger.Info("Checking schema differences...")
	if err := runAtlasCommand("schema", "diff", "--env", "dev", "--to", "file://schema.hcl"); err != nil {
		return fmt.Errorf("failed to check schema: %w", err)
	}
	return nil
}

// runDiffCommand generates migration diff
func runDiffCommand(logger *logger.Logger) error {
	logger.Info("Generating migration diff...")
	// Get migration name from args or use timestamp
	var migrationName string
	if len(os.Args) > 2 {
		migrationName = os.Args[2]
	} else {
		migrationName = fmt.Sprintf("migration_%s", time.Now().Format("20060102150405"))
	}

	if err := runAtlasCommand("migrate", "diff", migrationName, "--env", "dev"); err != nil {
		return fmt.Errorf("failed to generate diff: %w", err)
	}
	return nil
}

// runValidateCommand validates migrations
func runValidateCommand(logger *logger.Logger) error {
	logger.Info("Validating migrations...")
	if err := runAtlasCommand("migrate", "validate", "--env", "dev"); err != nil {
		return fmt.Errorf("failed to validate migrations: %w", err)
	}
	logger.Info("Migrations validated successfully")
	return nil
}

// runDevSetupCommand runs development setup with GORM AutoMigrate
func runDevSetupCommand(cfg *config.Config, logger *logger.Logger) error {
	if cfg.AppEnv != "development" {
		return fmt.Errorf("dev-setup command is only available in development environment")
	}
	logger.Info("Running development setup with GORM AutoMigrate...")
	if err := runGORMAutoMigration(cfg, logger); err != nil {
		return fmt.Errorf("failed to run development setup: %w", err)
	}
	logger.Info("Development setup completed successfully")
	return nil
}

// runAtlasCommand executes an Atlas CLI command
func runAtlasCommand(args ...string) error {
	// For migrate commands, we need to be in the directory where atlas.hcl is located
	cmd := exec.Command("atlas", args...)
	cmd.Dir = "/app"
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Set environment to ensure Atlas can find the config
	cmd.Env = append(os.Environ(), "ATLAS_CONFIG_PATH=/app/atlas.hcl")

	return cmd.Run()
}

// runGORMAutoMigration runs GORM's auto migration for initial development
func runGORMAutoMigration(cfg *config.Config, logger *logger.Logger) error {
	// Connect to database
	db, err := gorm.Open(postgres.Open(cfg.GetDSN()), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Run migrations
	logger.Info("Migrating database schema...")
	logger.Warn("Using GORM AutoMigrate for development only. Production should use Atlas migrations.")

	// Migrate in the correct order to handle foreign key dependencies
	if err := db.AutoMigrate(
		// Base tables first (no foreign keys)
		&model.User{},
		&model.CategoryMaster{},

		// Tables with foreign keys to base tables
		&model.Account{},
		&model.Category{},

		// Tables with foreign keys to above tables
		&model.Transaction{},
		&model.Budget{},
		&model.BudgetSuggestion{},
		&model.AssetSnapshot{},
		&model.AssetForecast{},
		&model.NotificationSetting{},
	); err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	// Note: Triggers are managed by Atlas migrations, not GORM AutoMigrate
	logger.Info("Note: Database triggers should be managed through Atlas migrations")

	// Add composite unique indexes
	if err := db.Exec(`
		CREATE UNIQUE INDEX IF NOT EXISTS idx_categories_user_master 
		ON categories(user_id, category_master_id)
	`).Error; err != nil {
		logger.Warn("Failed to create unique index on categories")
	}

	if err := db.Exec(`
		CREATE UNIQUE INDEX IF NOT EXISTS idx_asset_snapshots_user_date 
		ON asset_snapshots(user_id, date)
	`).Error; err != nil {
		logger.Warn("Failed to create unique index on asset_snapshots")
	}

	if err := db.Exec(`
		CREATE UNIQUE INDEX IF NOT EXISTS idx_notification_settings_user_type 
		ON notification_settings(user_id, notification_type, channel)
	`).Error; err != nil {
		logger.Warn("Failed to create unique index on notification_settings")
	}

	logger.Info("Database migration completed")
	return nil
}
