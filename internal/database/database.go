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

	log.Info("Database connected successfully",
		zap.String("host", cfg.DB.Host),
		zap.String("port", cfg.DB.Port),
		zap.String("database", cfg.DB.Name),
	)

	return nil
}

// AutoMigrate runs database migrations
func AutoMigrate(log *zap.Logger) error {
	if err := DB.AutoMigrate(&models.User{}, &models.Curriculums{}, &models.Work{}, &models.Configuration{}); err != nil {
		return errors.WrapError(err, "failed to run migrations")
	}

	log.Info("Database migrations completed successfully")
	return nil
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return DB
}
