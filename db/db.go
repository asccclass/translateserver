package db

import (
	"database/sql"
	// "log" // removed unused logging if needed, but keeping for DB.Exec error

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func Init() error {
	var err error
	DB, err = sql.Open("sqlite", "./jobs.db")
	if err != nil {
		return err
	}

	// Create table
	query := `
	CREATE TABLE IF NOT EXISTS jobs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		filename TEXT,
		filepath TEXT,
		email TEXT,
		status TEXT DEFAULT 'pending', -- pending, processing, completed, failed
		output_text TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err = DB.Exec(query)
	return err
}

func Close() {
	if DB != nil {
		DB.Close()
	}
}
