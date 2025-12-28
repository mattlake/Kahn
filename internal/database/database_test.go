package database

import (
	"kahn/internal/config"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

func TestNewDatabase_InvalidPath(t *testing.T) {
	// Create config with invalid database path
	config := &config.Config{}
	config.Database.Path = "/invalid/path/that/does/not/exist/kahn.db"
	config.Database.BusyTimeout = 5000
	config.Database.JournalMode = "WAL"
	config.Database.CacheSize = 10000
	config.Database.ForeignKeys = true

	database, err := NewDatabase(config)
	assert.Error(t, err, "NewDatabase should return error for invalid path")
	assert.Nil(t, database, "Database should be nil for invalid path")
	assert.Contains(t, err.Error(), "database directory cannot be created", "Error should mention directory creation")
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

func createTestConfig() *config.Config {
	config := &config.Config{}
	config.Database.Path = ":memory:"
	config.Database.BusyTimeout = 5000
	config.Database.JournalMode = "WAL"
	config.Database.CacheSize = 10000
	config.Database.ForeignKeys = true
	return config
}
