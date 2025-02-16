package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

// Use environment variables for configuration
var (
	DBUser     = getEnvOrDefault("DB_USER", "user")
	DBPassword = getEnvOrDefault("DB_PASSWORD", "password")
	DBName     = getEnvOrDefault("DB_NAME", "transformer")
	DBHost     = getEnvOrDefault("DB_HOST", "localhost")
	DBPort     = getEnvOrDefault("DB_PORT", "3306")
)

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// InitDB initializes the database connection
func InitDB() error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&collation=utf8mb4_unicode_ci",
		DBUser, DBPassword, DBHost, DBPort, DBName)

	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to MySQL: %v", err)
	}

	// Connection pool settings
	DB.SetMaxOpenConns(100)                 // Maximum number of open connections
	DB.SetMaxIdleConns(50)                  // Maximum number of idle connections
	DB.SetConnMaxLifetime(2 * time.Minute)  // Maximum connection lifetime
	DB.SetConnMaxIdleTime(10 * time.Minute) // Maximum idle time for connections

	// Verify connection
	if err := DB.Ping(); err != nil {
		return fmt.Errorf("MySQL connection failed: %v", err)
	}

	log.Println("Successfully connected to MySQL")
	return nil
}

// CloseDB closes the database connection
func CloseDB() {
	if DB != nil {
		if err := DB.Close(); err != nil {
			log.Printf("Error closing MySQL connection: %v", err)
			return
		}
		log.Println("MySQL connection closed")
	}
}
