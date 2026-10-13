package database

import "database/sql"

func setupExerciseDB(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS exercises (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			muscles_worked TEXT NOT NULL
		);
	`)
	
	return err
}