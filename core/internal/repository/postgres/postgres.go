package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

type Postgres struct {
	Db *sql.DB
}

func New(dsn string) (*Postgres, error) {
	var db *sql.DB
	var err error

	maxRetries := 10
	for i := 0; i < maxRetries; i++ {
		db, err = sql.Open("postgres", dsn)
		if err != nil {
			log.Printf("Failed to open postgres connection (attempt %d/%d): %v", i+1, maxRetries, err)
			if i < maxRetries-1 {
				time.Sleep(time.Duration(i+1) * time.Second)
				continue
			}
			return nil, fmt.Errorf("failed to open postgres after %d attempts: %w", maxRetries, err)
		}

		db.SetMaxOpenConns(25)
		db.SetMaxIdleConns(25)
		db.SetConnMaxLifetime(5 * time.Minute)

		if err := db.Ping(); err != nil {
			log.Printf("Failed to ping postgres (attempt %d/%d): %v", i+1, maxRetries, err)
			if i < maxRetries-1 {
				time.Sleep(time.Duration(i+1) * time.Second)
				continue
			}
			return nil, fmt.Errorf("failed to ping postgres after %d attempts: %w", maxRetries, err)
		}

		log.Printf("Successfully connected to postgres on attempt %d", i+1)
		break
	}

	return &Postgres{Db: db}, nil
}
