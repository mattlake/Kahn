package database

import (
	"kahn/internal/config"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateDatabasePath(t *testing.T) {
	// Get current environment paths for dynamic test generation
	homeDir, _ := os.UserHomeDir()
	tempDir := os.TempDir()

	tests := []struct {
		name        string
		path        string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "Valid relative path",
			path:        "./test.db",
			expectError: false,
		},
		{
			name:        "Valid absolute path in home",
			path:        filepath.Join(homeDir, ".kahn", "test.db"),
			expectError: false,
		},
		{
			name:        "Valid temp directory path",
			path:        filepath.Join(tempDir, "kahn_test.db"),
			expectError: false,
		},
		{
			name:        "Valid absolute path in home",
			path:        filepath.Join(homeDir, ".kahn", "test.db"),
			expectError: false,
		},
		{
			name:        "Valid temp directory path",
			path:        filepath.Join(tempDir, "kahn_test.db"),
			expectError: false,
		},
		{
			name:        "Directory traversal with ..",
			path:        "../test.db",
			expectError: true,
			errorMsg:    "directory traversal attempts",
		},
		{
			name:        "Directory traversal with ../..",
			path:        "../../test.db",
			expectError: true,
			errorMsg:    "suspicious directory traversal pattern",
		},
		{
			name:        "Directory traversal with ../test/../",
			path:        "../test/../test.db",
			expectError: true,
			errorMsg:    "directory traversal attempts",
		},
		{
			name:        "Suspicious pattern",
			path:        "../../../etc/passwd",
			expectError: true,
			errorMsg:    "suspicious directory traversal pattern",
		},
		{
			name:        "Absolute path outside safe directories",
			path:        "/etc/passwd",
			expectError: true,
			errorMsg:    "must be within user home or temp directory",
		},
		{
			name:        "Another unsafe absolute path",
			path:        "/usr/local/bin/test.db",
			expectError: true,
			errorMsg:    "must be within user home or temp directory",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateDatabasePath(tt.path)

			if tt.expectError {
				assert.Error(t, err, "Should return error for path: %s", tt.path)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg, "Error message should contain expected text")
				}
			} else {
				assert.NoError(t, err, "Should not return error for valid path: %s", tt.path)
			}
		})
	}
}

func TestNewDatabase(t *testing.T) {
	config := createTestConfig()

	database, err := NewDatabase(config)
	require.NoError(t, err, "NewDatabase should not return error")
	require.NotNil(t, database, "Database should not be nil")

	// Test that database connection is working
	db := database.GetDB()
	require.NotNil(t, db, "Database connection should not be nil")

	// Test basic query
	err = db.Ping()
	assert.NoError(t, err, "Database should be pingable")

	// Clean up
	err = database.Close()
	assert.NoError(t, err, "Database should close without error")
}

func TestNewDatabase_SecurityValidation(t *testing.T) {
	// Create config with dangerous database path
	config := &config.Config{}
	config.Database.Path = "../../../etc/passwd" // Dangerous path
	config.Database.BusyTimeout = 5000
	config.Database.JournalMode = "WAL"
	config.Database.CacheSize = 10000
	config.Database.ForeignKeys = true

	database, err := NewDatabase(config)
	assert.Error(t, err, "NewDatabase should return error for dangerous path")
	assert.Nil(t, database, "Database should be nil for dangerous path")
	assert.Contains(t, err.Error(), "invalid database path", "Error should mention invalid path")
}

func TestDatabase_BeginTransaction(t *testing.T) {
	config := createTestConfig()
	database, err := NewDatabase(config)
	require.NoError(t, err, "NewDatabase should not return error")
	defer database.Close()

	tx, err := database.BeginTransaction()
	require.NoError(t, err, "BeginTransaction should not return error")
	require.NotNil(t, tx, "Transaction should not be nil")

	// Test transaction rollback
	err = tx.Rollback()
	assert.NoError(t, err, "Transaction rollback should not return error")
}

func TestDatabase_GetDB(t *testing.T) {
	config := createTestConfig()
	database, err := NewDatabase(config)
	require.NoError(t, err, "NewDatabase should not return error")
	defer database.Close()

	db := database.GetDB()
	assert.NotNil(t, db, "GetDB should return non-nil database connection")

	// Test that it's the same connection
	db2 := database.GetDB()
	assert.Equal(t, db, db2, "GetDB should return the same connection")
}

func TestDatabase_Close(t *testing.T) {
	config := createTestConfig()
	database, err := NewDatabase(config)
	require.NoError(t, err, "NewDatabase should not return error")

	// Test close
	err = database.Close()
	assert.NoError(t, err, "Close should not return error")

	// Test double close (should not panic)
	err = database.Close()
	assert.NoError(t, err, "Double close should not return error")
}

func TestDatabase_Close_NilDB(t *testing.T) {
	database := &Database{}

	err := database.Close()
	assert.NoError(t, err, "Close with nil DB should not return error")
}

func TestNewDatabase_InvalidPath(t *testing.T) {
	config := createTestConfig()
	config.Database.Path = "../../../etc/passwd" // Dangerous path

	database, err := NewDatabase(config)
	assert.Error(t, err, "NewDatabase should return error for invalid path")
	assert.Nil(t, database, "Database should be nil for invalid path")
	assert.Contains(t, err.Error(), "invalid database path", "Error should mention invalid path")
}

func createTestConfig() *config.Config {
	config := &config.Config{}
	config.Database.Path = ":memory:"
	config.Database.BusyTimeout = 5000
	config.Database.JournalMode = "WAL"
	config.Database.CacheSize = 10000
	config.Database.ForeignKeys = true
	return config
}
