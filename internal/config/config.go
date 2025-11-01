package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds application configuration
type Config struct {
	Port       string
	Mode       string
	DB         DatabaseConfig
	Redis      RedisConfig
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

// RedisConfig holds Redis configuration
type RedisConfig struct {
	PublicURL         string
	Host              string
	Port              string
	Username          string
	Password          string
	DB                string
	MaxMemory         string
	MaxMemoryPolicy   string
	MemoryLimit       string
	MemoryReservation string
}

// WorkerPoolConfig holds worker pool configuration
type WorkerPoolConfig struct {
	NumWorkers int
	QueueSize  int
}

// EmailConfig holds email configuration
type EmailConfig struct {
	APIKey string
	From   string
}

// AppConfig holds application configuration
type AppConfig struct {
	URL         string
	StaticToken string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found, using system environment variables")
	}

	port := os.Getenv("PORT")

	mode := os.Getenv("GIN_MODE")

	// Email configuration
	emailAPIKey := os.Getenv("RESEND_API_KEY")
	emailFrom := os.Getenv("MAIL_FROM")

	// App configuration
	appURL := os.Getenv("APP_URL")
	staticToken := os.Getenv("BACKEND_APIKEY")

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
		Redis: RedisConfig{
			PublicURL:         os.Getenv("REDIS_PUBLIC_URL"),
			Host:              os.Getenv("REDIS_HOST"),
			Port:              os.Getenv("REDIS_PORT"),
			Username:          os.Getenv("REDISUSER"),
			Password:          os.Getenv("REDIS_PASSWORD"),
			DB:                os.Getenv("REDIS_DB"),
			MaxMemory:         os.Getenv("REDIS_MAX_MEMORY"),
			MaxMemoryPolicy:   os.Getenv("REDIS_MAX_MEMORY_POLICY"),
			MemoryLimit:       os.Getenv("REDIS_MEMORY_LIMIT"),
			MemoryReservation: os.Getenv("REDIS_MEMORY_RESERVATION"),
		},
		Email: EmailConfig{
			APIKey: emailAPIKey,
			From:   emailFrom,
		},
		App: AppConfig{
			URL:         appURL,
			StaticToken: staticToken,
		},
	}
}
