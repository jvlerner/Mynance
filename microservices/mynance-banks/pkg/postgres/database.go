package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	const maxRetries = 10
	const retryInterval = 3 * time.Second

	var err error
	for i := 1; i <= maxRetries || maxRetries == -1; i++ {
		DB, err = sql.Open("postgres", connStr)
		if err != nil {
			log.Printf("[ERROR] [DB] Attempt %d: failed to open connection: %v", i, err)
			time.Sleep(retryInterval)
			continue
		}

		err = DB.Ping()
		if err == nil {
			DB.SetConnMaxLifetime(5 * time.Minute)
			DB.SetConnMaxIdleTime(3 * time.Minute)
			DB.SetMaxIdleConns(5)
			DB.SetMaxOpenConns(10)
			log.Println("[INFO] [DB] Successfully connected to the database.")
			return
		}

		log.Printf("[ERROR] [DB] Attempt %d: failed to ping database: %v", i, err)
		time.Sleep(retryInterval)
	}

	log.Fatal("[ERROR] [DB] Failed to connect to the database after multiple attempts.")
}

// CloseDB closes the database connection
func CloseDB() {
	if DB != nil {
		err := DB.Close()
		if err != nil {
			log.Printf("[ERROR] [DB] Error closing database connection: %v", err)
		} else {
			log.Println("[INFO] [DB] Database connection closed successfully.")
		}
	}
}
