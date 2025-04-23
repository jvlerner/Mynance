package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	_ "github.com/lib/pq"
)

var (
	dbs  = make(map[string]*sql.DB)
	lock = sync.RWMutex{}
)

// InitDB initializes and stores a database connection under the given name
func InitDB(name, user, password, host, port, dbName string) {
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		user,
		password,
		host,
		port,
		dbName,
	)

	const maxRetries = 10
	const retryInterval = 3 * time.Second

	var db *sql.DB
	var err error

	for i := 1; i <= maxRetries || maxRetries == -1; i++ {
		db, err = sql.Open("postgres", connStr)
		if err != nil {
			log.Printf("[ERROR] [DB] Attempt %d: failed to open connection: %v", i, err)
			time.Sleep(retryInterval)
			continue
		}

		err = db.Ping()
		if err == nil {
			db.SetConnMaxLifetime(5 * time.Minute)
			db.SetConnMaxIdleTime(3 * time.Minute)
			db.SetMaxIdleConns(5)
			db.SetMaxOpenConns(10)

			lock.Lock()
			dbs[name] = db
			lock.Unlock()

			log.Printf("[INFO] [DB] Successfully connected to database '%s'.", name)
			return
		}

		log.Printf("[ERROR] [DB] Attempt %d: failed to ping database: %v", i, err)
		time.Sleep(retryInterval)
	}

	log.Fatalf("[ERROR] [DB] Failed to connect to database '%s' after multiple attempts.", name)
}

// GetDB retrieves a previously initialized DB by name
func GetDB(name string) *sql.DB {
	lock.RLock()
	defer lock.RUnlock()
	return dbs[name]
}

// CloseDB closes the DB connection and removes it from the map
func CloseDB(name string) {
	lock.Lock()
	defer lock.Unlock()

	if db, exists := dbs[name]; exists && db != nil {
		err := db.Close()
		if err != nil {
			log.Printf("[ERROR] [DB] Error closing DB '%s': %v", name, err)
		} else {
			log.Printf("[INFO] [DB] Closed DB '%s' successfully.", name)
		}
		delete(dbs, name)
	}
}

// CloseAll closes all registered database connections
func CloseAll() {
	lock.Lock()
	defer lock.Unlock()

	for name, db := range dbs {
		if db != nil {
			err := db.Close()
			if err != nil {
				log.Printf("[ERROR] [DB] Error closing DB '%s': %v", name, err)
			} else {
				log.Printf("[INFO] [DB] Closed DB '%s' successfully.", name)
			}
		}
	}
	dbs = make(map[string]*sql.DB)
}
