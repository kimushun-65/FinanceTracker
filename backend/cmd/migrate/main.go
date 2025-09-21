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
		log.Fatal("Usage: migrate [apply|check|diff|auto|validate|gorm]")
	}

	command := os.Args[1]

	switch command {
	case "apply":
		logger.Info("Applying Atlas migrations...")
		// Change working directory to where atlas.hcl is located
		if err := os.Chdir("/app"); err != nil {
			logger.Error("Failed to change directory")
			os.Exit(1)
		}
		
		// Debug: Check current directory and files
		cwd, _ := os.Getwd()
		logger.Info("Current working directory: " + cwd)
		
		// Check if atlas.hcl exists
		if _, err := os.Stat("atlas.hcl"); err != nil {
			logger.Error("atlas.hcl not found in current directory")
		} else {
			logger.Info("atlas.hcl found")
		}
		
		// Check if migrations directory exists
		if _, err := os.Stat("cmd/migrate/migrations"); err != nil {
			logger.Error("migrations directory not found at cmd/migrate/migrations")
		} else {
			logger.Info("migrations directory found")
		}
		
		// Run with verbose output for debugging
		if err := runAtlasCommand("migrate", "apply", "--env", "dev", "-v"); err != nil {
			logger.Error("Failed to apply migrations")
			// Also try listing the migrations to see what Atlas can find
			logger.Info("Attempting to list migrations...")
			_ = runAtlasCommand("migrate", "status", "--env", "dev")
			os.Exit(1)
		}
		logger.Info("Migrations applied successfully")

	case "check":
		logger.Info("Checking schema differences...")
		if err := runAtlasCommand("schema", "diff", "--env", "dev", "--to", "file://schema.hcl"); err != nil {
			logger.Error("Failed to check schema")
			os.Exit(1)
		}

	case "diff":
		logger.Info("Generating migration diff...")
		// Get migration name from args or use timestamp
		var migrationName string
		if len(os.Args) > 2 {
			migrationName = os.Args[2]
		} else {
			migrationName = fmt.Sprintf("migration_%s", time.Now().Format("20060102150405"))
		}
		
		if err := runAtlasCommand("migrate", "diff", migrationName, "--env", "dev"); err != nil {
			logger.Error("Failed to generate diff")
			os.Exit(1)
		}

	case "auto":
		logger.Info("Running auto migration with GORM...")
		if err := runGORMAutoMigration(cfg, logger); err != nil {
			logger.Error("Failed to run auto migration")
			os.Exit(1)
		}
		logger.Info("Auto migration completed successfully")

	case "validate":
		logger.Info("Validating migrations...")
		if err := runAtlasCommand("migrate", "validate", "--env", "dev"); err != nil {
			logger.Error("Failed to validate migrations")
			os.Exit(1)
		}
		logger.Info("Migrations validated successfully")

	case "gorm":
		logger.Info("Running GORM auto migration...")
		if err := runGORMAutoMigration(cfg, logger); err != nil {
			logger.Error("Failed to run GORM auto migration")
			os.Exit(1)
		}
		logger.Info("GORM auto migration completed successfully")

	default:
		log.Fatalf("Unknown command: %s", command)
	}
}

// runAtlasCommand executes an Atlas CLI command
func runAtlasCommand(args ...string) error {
	cmd := exec.Command("atlas", args...)
	// Working directory should be where atlas.hcl is located
	// since atlas.hcl uses relative paths
	cmd.Dir = "/app" 
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
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
	
	// First create the trigger function
	logger.Info("Creating trigger function...")
	triggerFunctionSQL := `
		CREATE OR REPLACE FUNCTION trigger_set_timestamp()
		RETURNS TRIGGER AS $$
		BEGIN
			NEW.updated_at = CURRENT_TIMESTAMP;
			RETURN NEW;
		END;
		$$ LANGUAGE plpgsql;
	`
	if err := db.Exec(triggerFunctionSQL).Error; err != nil {
		return fmt.Errorf("failed to create trigger function: %w", err)
	}
	
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

	// Create updated_at triggers for all tables
	tables := []string{
		"users", "accounts", "category_masters", "categories",
		"transactions", "budgets", "budget_suggestions",
		"asset_snapshots", "asset_forecasts", "notification_settings",
	}

	for _, table := range tables {
		triggerSQL := fmt.Sprintf(`
			CREATE TRIGGER update_%s_updated_at 
			BEFORE UPDATE ON %s 
			FOR EACH ROW 
			EXECUTE FUNCTION trigger_set_timestamp();
		`, table, table)
		
		if err := db.Exec(triggerSQL).Error; err != nil {
			// Ignore error if trigger already exists
			logger.Info(fmt.Sprintf("Trigger for %s might already exist", table))
		}
	}

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
