package storage

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	"k8s.io/klog/v2"
)

// SQLiteDB represents a SQLite database connection
type SQLiteDB struct {
	db *sql.DB
}

// NewSQLiteDB creates a new SQLite database connection
func NewSQLiteDB(dbPath string) (*SQLiteDB, error) {
	// Ensure directory exists
	dir := filepath.Dir(dbPath)
	if dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create database directory: %w", err)
		}
	}

	// Open database connection
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Initialize database
	if err := initializeDatabase(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	klog.InfoS("Connected to SQLite database", "path", dbPath)
	return &SQLiteDB{db: db}, nil
}

// Close closes the database connection
func (s *SQLiteDB) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

// initializeDatabase creates the necessary tables if they don't exist
func initializeDatabase(db *sql.DB) error {
	// Create RSS sources table
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS rss_sources (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			url TEXT NOT NULL UNIQUE,
			description TEXT,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			last_fetched_at TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create rss_sources table: %w", err)
	}

	// Create RSS contents table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS rss_contents (
			id TEXT PRIMARY KEY,
			source_id TEXT NOT NULL,
			title TEXT NOT NULL,
			link TEXT NOT NULL,
			description TEXT,
			content TEXT,
			published_at TIMESTAMP NOT NULL,
			fetched_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP,
			author TEXT,
			FOREIGN KEY (source_id) REFERENCES rss_sources(id) ON DELETE CASCADE
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create rss_contents table: %w", err)
	}

	// Create content categories table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS content_categories (
			content_id TEXT NOT NULL,
			category TEXT NOT NULL,
			PRIMARY KEY (content_id, category),
			FOREIGN KEY (content_id) REFERENCES rss_contents(id) ON DELETE CASCADE
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create content_categories table: %w", err)
	}

	// Create fetch jobs table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS fetch_jobs (
			id TEXT PRIMARY KEY,
			status TEXT NOT NULL,
			started_at TIMESTAMP NOT NULL,
			completed_at TIMESTAMP,
			items_processed INTEGER DEFAULT 0,
			source_id TEXT,
			days INTEGER DEFAULT 1,
			FOREIGN KEY (source_id) REFERENCES rss_sources(id) ON DELETE SET NULL
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create fetch_jobs table: %w", err)
	}

	// Create job errors table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS job_errors (
			job_id TEXT NOT NULL,
			error_message TEXT NOT NULL,
			timestamp TIMESTAMP NOT NULL,
			FOREIGN KEY (job_id) REFERENCES fetch_jobs(id) ON DELETE CASCADE
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create job_errors table: %w", err)
	}

	// Create recommendation feedback table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS recommendation_feedback (
			id TEXT PRIMARY KEY,
			content_id TEXT NOT NULL,
			feedback_type TEXT NOT NULL,
			comment TEXT,
			created_at TIMESTAMP NOT NULL,
			FOREIGN KEY (content_id) REFERENCES rss_contents(id) ON DELETE CASCADE
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create recommendation_feedback table: %w", err)
	}

	return nil
}
