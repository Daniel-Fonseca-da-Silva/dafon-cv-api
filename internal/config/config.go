package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds application configuration
type Config struct {
	Port       string
	Mode       string
	DB         DatabaseConfig
	WorkerPool WorkerPoolConfig
	Email      EmailConfig
	App        AppConfig
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

// WorkerPoolConfig holds worker pool configuration
type WorkerPoolConfig struct {
	NumWorkers int
	QueueSize  int
}

// EmailConfig holds email configuration
type EmailConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string
}

// AppConfig holds application configuration
type AppConfig struct {
	URL string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found, using system environment variables")
	}

	port := os.Getenv("PORT")

	mode := os.Getenv("GIN_MODE")

	// Worker pool configuration
	numWorkers := getEnvAsInt("WORKER_POOL_SIZE", "5")
	queueSize := getEnvAsInt("WORKER_QUEUE_SIZE", "100")

	// Email configuration
	emailHost := os.Getenv("MAIL_HOST")

	emailPort := os.Getenv("MAIL_PORT")

	emailUsername := os.Getenv("MAIL_USERNAME")

	emailPassword := os.Getenv("MAIL_PASSWORD")

	emailFrom := os.Getenv("MAIL_FROM")

	// App configuration
	appURL := os.Getenv("APP_URL")

	return &Config{
		Port: port,
		Mode: mode,
		DB: DatabaseConfig{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Name:     os.Getenv("DB_NAME"),
			SSLMode:  os.Getenv("DB_SSL_MODE"),
		},
		WorkerPool: WorkerPoolConfig{
			NumWorkers: numWorkers,
			QueueSize:  queueSize,
		},
		Email: EmailConfig{
			Host:     emailHost,
			Port:     emailPort,
			Username: emailUsername,
			Password: emailPassword,
			From:     emailFrom,
		},
		App: AppConfig{
			URL: appURL,
		},
	}
}

// getEnvAsInt gets environment variable as integer with fallback
func getEnvAsInt(key, fallback string) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	if intValue, err := strconv.Atoi(fallback); err == nil {
		return intValue
	}
	return 5 // Default fallback
}
