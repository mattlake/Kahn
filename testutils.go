package main

import (
	"database/sql"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
)

// setupTestDB creates an in-memory SQLite database for testing
func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err, "Failed to open in-memory database")

	// Test the connection
	err = db.Ping()
	require.NoError(t, err, "Failed to ping test database")

	// Run migrations on the test database
	database := &Database{db: db}
	err = database.runMigrations()
	require.NoError(t, err, "Failed to run migrations on test database")

	return db
}

// createTestProject creates a test project with the given parameters
func createTestProject(name, description, color string) *Project {
	now := time.Now()
	return &Project{
		ID:          "test_proj_" + name,
		Name:        name,
		Description: description,
		Color:       color,
		CreatedAt:   now,
		UpdatedAt:   now,
		Tasks:       []Task{},
	}
}

// createTestTask creates a test task with the given parameters
func createTestTask(name, description, projectID string, status Status) *Task {
	now := time.Now()
	return &Task{
		ID:        "test_task_" + name,
		ProjectID: projectID,
		Name:      name,
		Desc:      description,
		Status:    status,
		CreatedAt: now,
		UpdatedAt: now,
		Priority:  Medium,
		Tags:      []string{},
	}
}

// createTestConfig creates a test configuration with in-memory database
func createTestConfig() *Config {
	config := &Config{}
	config.Database.Path = ":memory:"
	config.Database.BusyTimeout = 5000
	config.Database.JournalMode = "WAL"
	config.Database.CacheSize = 10000
	config.Database.ForeignKeys = true
	return config
}

// cleanupTestDB closes the database connection and handles cleanup
func cleanupTestDB(t *testing.T, db *sql.DB) {
	err := db.Close()
	require.NoError(t, err, "Failed to close test database")
}

// assertProjectEqual compares two projects for equality, ignoring time differences within a threshold
func assertProjectEqual(t *testing.T, expected, actual *Project) {
	require.Equal(t, expected.ID, actual.ID, "Project ID should match")
	require.Equal(t, expected.Name, actual.Name, "Project name should match")
	require.Equal(t, expected.Description, actual.Description, "Project description should match")
	require.Equal(t, expected.Color, actual.Color, "Project color should match")
	require.Equal(t, len(expected.Tasks), len(actual.Tasks), "Number of tasks should match")
}

// assertTaskEqual compares two tasks for equality, ignoring time differences within a threshold
func assertTaskEqual(t *testing.T, expected, actual *Task) {
	require.Equal(t, expected.ID, actual.ID, "Task ID should match")
	require.Equal(t, expected.ProjectID, actual.ProjectID, "Task ProjectID should match")
	require.Equal(t, expected.Name, actual.Name, "Task name should match")
	require.Equal(t, expected.Desc, actual.Desc, "Task description should match")
	require.Equal(t, expected.Status, actual.Status, "Task status should match")
	require.Equal(t, expected.Priority, actual.Priority, "Task priority should match")

	// Handle nil tags comparison
	if expected.Tags == nil && actual.Tags == nil {
		// Both nil, that's fine
	} else if expected.Tags == nil {
		require.Equal(t, []string{}, actual.Tags, "Actual tags should be empty when expected is nil")
	} else if actual.Tags == nil {
		require.Equal(t, expected.Tags, []string{}, "Expected tags should be empty when actual is nil")
	} else {
		require.Equal(t, expected.Tags, actual.Tags, "Task tags should match")
	}
}

// insertTestProject inserts a test project into the database and returns it
func insertTestProject(t *testing.T, db *sql.DB, project *Project) *Project {
	dao := NewProjectDAO(db)
	err := dao.Create(project)
	require.NoError(t, err, "Failed to insert test project")
	return project
}

// insertTestTask inserts a test task into the database and returns it
func insertTestTask(t *testing.T, db *sql.DB, task *Task) *Task {
	dao := NewTaskDAO(db)
	err := dao.Create(task)
	require.NoError(t, err, "Failed to insert test task")
	return task
}

// countTableRows returns the number of rows in a given table
func countTableRows(t *testing.T, db *sql.DB, tableName string) int {
	var count int
	query := "SELECT COUNT(*) FROM " + tableName
	err := db.QueryRow(query).Scan(&count)
	require.NoError(t, err, "Failed to count rows in table %s", tableName)
	return count
}
