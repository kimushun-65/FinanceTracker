package database

import (
	"fmt"

	"financetracker/pkg/config"
	"financetracker/pkg/logger"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Database represents the database connection
type Database struct {
	DB     *gorm.DB
	logger *logger.Logger
}

// New creates a new database connection
func New(cfg *config.Config, logger *logger.Logger) (*Database, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort, cfg.DBSSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return &Database{
		DB:     db,
		logger: logger,
	}, nil
}

// Close closes the database connection
func (d *Database) Close() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// Migrate runs database migrations
func (d *Database) Migrate(models ...interface{}) error {
	return d.DB.AutoMigrate(models...)
}

// Health checks the database health
func (d *Database) Health() error {
	sqlDB, err := d.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}
