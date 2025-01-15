package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

func InitDB() *sql.DB {
	user := getEnv("DB_USER", "neondb_owner")
	password := getEnv("DB_PASSWORD", "your_default_password")
	host := getEnv("DB_HOST", "ep-tight-union-a7q3zgwi.ap-southeast-2.aws.neon.tech")
	port := getEnv("DB_PORT", "5432")
	dbName := getEnv("DB_NAME", "neondb")
	sslMode := getEnv("DB_SSLMODE", "require")

	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s",
		user, password, host, port, dbName, sslMode)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	// Set connection pool parameters
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		log.Fatalf("Error pinging database: %v", err)
	}

	log.Println("Successfully connected to database")
	return db
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
