package database

import (
	"database/sql"
	"fmt"
	"kahn/internal/config"
	"log"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	Db *sql.DB
}

// NewDatabase creates a new database connection with optimized settings
func NewDatabase(config *config.Config) (*Database, error) {
	// Ensure database directory exists
	dbDir := filepath.Dir(config.Database.Path)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return nil, fmt.Errorf("database directory cannot be created at %s\n"+
			"Suggestions:\n"+
			"  • Check parent directory exists\n"+
			"  • Verify write permissions\n"+
			"  • Try: mkdir -p %s", dbDir, dbDir)
	}

	// Build connection string with optimized settings
	dsn := buildConnectionString(config.Database)

	// Open database connection
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, handleDatabaseError(config.Database.Path, err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, handleDatabaseError(config.Database.Path, err)
	}

	// Set additional pragmas
	if err := setDatabasePragmas(db, config.Database); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to set database pragmas: %w", err)
	}

	database := &Database{Db: db}

	// Run migrations
	if err := database.RunMigrations(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return database, nil
}

// buildConnectionString creates an optimized SQLite connection string
func buildConnectionString(dbConfig struct {
	Path        string `mapstructure:"path"`
	BusyTimeout int    `mapstructure:"busy_timeout"`
	JournalMode string `mapstructure:"journal_mode"`
	CacheSize   int    `mapstructure:"cache_size"`
	ForeignKeys bool   `mapstructure:"foreign_keys"`
}) string {
	// Build DSN with performance optimizations
	dsn := fmt.Sprintf("%s?_journal_mode=%s&_busy_timeout=%d&_cache_size=%d&_foreign_keys=%t&_temp_store=memory",
		dbConfig.Path,
		dbConfig.JournalMode,
		dbConfig.BusyTimeout,
		dbConfig.CacheSize,
		dbConfig.ForeignKeys,
	)

	return dsn
}

// setDatabasePragmas sets additional database pragmas for performance
func setDatabasePragmas(db *sql.DB, _ struct {
	Path        string `mapstructure:"path"`
	BusyTimeout int    `mapstructure:"busy_timeout"`
	JournalMode string `mapstructure:"journal_mode"`
	CacheSize   int    `mapstructure:"cache_size"`
	ForeignKeys bool   `mapstructure:"foreign_keys"`
}) error {
	pragmas := map[string]any{
		"synchronous": "NORMAL",  // Better performance while maintaining safety
		"mmap_size":   268435456, // 256MB memory-mapped I/O
		"optimize":    "0x10002", // Enable optimizations
	}

	for pragma, value := range pragmas {
		_, err := db.Exec(fmt.Sprintf("PRAGMA %s = %v", pragma, value))
		if err != nil {
			log.Printf("Warning: failed to set pragma %s: %v", pragma, err)
		}
	}

	return nil
}

// handleDatabaseError provides detailed error messages with suggestions
func handleDatabaseError(dbPath string, err error) error {
	if os.IsPermission(err) {
		return fmt.Errorf("database permission denied at %s\n"+
			"Suggestions:\n"+
			"  • Check file permissions\n"+
			"  • Try a different location: kahn --db-path /tmp/kahn.db\n"+
			"  • Run with different user if needed", dbPath)
	}

	if os.IsNotExist(err) {
		return fmt.Errorf("database path does not exist at %s\n"+
			"Suggestions:\n"+
			"  • Check parent directory exists\n"+
			"  • Verify write permissions\n"+
			"  • Try: mkdir -p %s", dbPath, filepath.Dir(dbPath))
	}

	// Check for common SQLite issues
	errStr := err.Error()
	if strings.Contains(errStr, "unable to open database file") {
		return fmt.Errorf("unable to open database file at %s\n"+
			"Suggestions:\n"+
			"  • Verify the directory exists and is writable\n"+
			"  • Check if the database file is corrupted\n"+
			"  • Try removing the file and restarting: rm %s", dbPath, dbPath)
	}

	if strings.Contains(errStr, "database is locked") {
		return fmt.Errorf("database is locked at %s\n"+
			"Suggestions:\n"+
			"  • Another instance of Kahn may be running\n"+
			"  • Check for other processes using the database\n"+
			"  • Try restarting after a few seconds", dbPath)
	}

	return fmt.Errorf("database initialization failed at %s: %w\n"+
		"This is likely a configuration or system issue.\n"+
		"Try running with --help to see configuration options.", dbPath, err)
}

// Close closes the database connection
func (d *Database) Close() error {
	if d.Db != nil {
		return d.Db.Close()
	}
	return nil
}

// GetDB returns the underlying database connection
func (d *Database) GetDB() *sql.DB {
	return d.Db
}

// BeginTransaction starts a new database transaction
func (d *Database) BeginTransaction() (*sql.Tx, error) {
	return d.Db.Begin()
}

// runMigrations runs all pending database migrations
func (d *Database) RunMigrations() error {
	// Create migrations table if it doesn't exist
	_, err := d.Db.Exec(`
		CREATE TABLE IF NOT EXISTS migrations (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE,
			executed_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Run each migration
	migrations := getMigrations()
	for _, migration := range migrations {
		// Check if migration already ran
		var count int
		err := d.Db.QueryRow("SELECT COUNT(*) FROM migrations WHERE name = ?", migration.name).Scan(&count)
		if err != nil {
			return fmt.Errorf("failed to check migration %s: %w", migration.name, err)
		}

		if count > 0 {
			continue // Migration already executed
		}

		// Execute migration
		_, err = d.Db.Exec(migration.sql)
		if err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", migration.name, err)
		}

		// Record migration
		_, err = d.Db.Exec("INSERT INTO migrations (name) VALUES (?)", migration.name)
		if err != nil {
			return fmt.Errorf("failed to record migration %s: %w", migration.name, err)
		}

		log.Printf("Executed migration: %s", migration.name)
	}

	return nil
}
