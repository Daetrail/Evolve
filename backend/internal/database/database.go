package database

import (
	"fmt"
	"database/sql"
	"path/filepath"
	"os"
	
	_ "modernc.org/sqlite"
)

func Open(name string) (*sql.DB, error) {
	dir, err := os.Getwd()
    if err != nil {
        return nil, fmt.Errorf("failed to get working directory: %w", err)
    }
	
    dataDir := filepath.Join(dir, "data")
    if err := os.MkdirAll(dataDir, 0755); err != nil {
        return nil, fmt.Errorf("failed to create data directory: %w", err)
    }
	
    dbPath := filepath.Join(dataDir, name)
    
    dsn := dbPath +
    	"?_pragma=journal_mode(WAL)" +
     	"&_pragma=foreign_keys(1)" +
      	"&_pragma=busy_timeout(5000)"
    
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	
	return db, nil
}

func SetupDB(db *sql.DB) error {
	err := setupUserDB(db)
	if err != nil {
		return err
	}
	err = setupExerciseDB(db) 
	if err != nil {
		return err
	}
	return nil
}