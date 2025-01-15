package config

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq" 
)

func InitDB() *sql.DB {
  connStr := "postgresql://neondb_owner:6xB3UyVfjkON@ep-tight-union-a7q3zgwi.ap-southeast-2.aws.neon.tech/neondb?sslmode=require"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	db.SetMaxOpenConns(10)                 // Maximum number of open connections
	db.SetMaxIdleConns(5)                  // Maximum number of idle connections
	db.SetConnMaxLifetime(5 * time.Minute) // Maximum connection lifetime

	if err := db.Ping(); err != nil {
		log.Fatalf("Error pinging database: %v", err)
	}

	log.Println("Successfully connected to database")
	return db
}
