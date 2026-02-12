package database

import (
	"fmt"

	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/config"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/errors"
	"github.com/Daniel-Fonseca-da-Silva/dafon-cv-api/internal/models"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB holds the database connection
var DB *gorm.DB

// Connect establishes database connection
func Connect(cfg *config.Config, log *zap.Logger) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.Name,
	)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		return errors.WrapError(err, "failed to connect to database")
	}

	log.Info("Database connected successfully")

	return nil
}

// AutoMigrate runs database migrations
func AutoMigrate(log *zap.Logger) error {
	// Temporarily disable SQL logging during migrations to avoid cluttering logs
	originalLogger := DB.Config.Logger
	DB.Config.Logger = logger.Default.LogMode(logger.Silent)

	if err := DB.AutoMigrate(
		&models.User{},
		&models.Subscription{},
		&models.Curriculums{},
		&models.Work{},
		&models.Configuration{},
		&models.Session{},
		&models.Education{},
	); err != nil {
		// Restore original logger before returning error
		DB.Config.Logger = originalLogger
		return errors.WrapError(err, "failed to run migrations")
	}

	// Restore original logger after migrations
	DB.Config.Logger = originalLogger

	log.Info("Database migrations completed successfully")
	return nil
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return DB
}
