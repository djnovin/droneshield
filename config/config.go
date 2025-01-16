package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

func InitDB() (*sql.DB, error) {
	user := getEnv("DB_USER", "neondb_owner")
	password := getEnv("DB_PASSWORD", "your_default_password")
	host := getEnv("DB_HOST", "ep-tight-union-a7q3zgwi.ap-southeast-2.aws.neon.tech")
	port := getEnv("DB_PORT", "5432")
	dbName := getEnv("DB_NAME", "neondb")
	sslMode := getEnv("DB_SSLMODE", "require")

	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s",
		user, password, host, port, dbName, sslMode)

  db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

  db.SetMaxOpenConns(10)
  db.SetMaxIdleConns(5)
  db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error pinging database: %w", err)
	}

	return db, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
