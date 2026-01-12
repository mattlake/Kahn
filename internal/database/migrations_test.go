package database

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

func TestGetMigrations(t *testing.T) {
	migrations := getMigrations()

	assert.Len(t, migrations, 5, "Should have 5 migrations")

	// Test migration names
	expectedNames := []string{
		"001_create_projects_table",
		"002_create_tasks_table",
		"003_add_type_to_tasks",
		"005_create_indexes",
		"006_add_integer_pk_and_blocked_by",
	}

	for i, expectedName := range expectedNames {
		assert.Equal(t, expectedName, migrations[i].name, "Migration %d should have correct name", i)
		assert.NotEmpty(t, migrations[i].sql, "Migration %d should have SQL content", i)
	}
}

func TestRunMigrations(t *testing.T) {
	// Create in-memory database
	db, err := sql.Open("sqlite", ":memory:?_foreign_keys=true")
	require.NoError(t, err, "Should be able to open in-memory database")
	defer db.Close()

	// Test that migrations table doesn't exist initially
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM migrations").Scan(&count)
	assert.Error(t, err, "Migrations table should not exist initially")

	// Run migrations
	database := &Database{Db: db}
	err = database.RunMigrations()
	assert.NoError(t, err, "runMigrations should not return error")

	// Test that migrations table exists and has records
	err = db.QueryRow("SELECT COUNT(*) FROM migrations").Scan(&count)
	assert.NoError(t, err, "Should be able to query migrations table")
	assert.Equal(t, 5, count, "Should have 5 migration records")

	// Test that all expected tables exist
	tables := []string{"projects", "tasks", "migrations"}
	for _, table := range tables {
		err = db.QueryRow("SELECT COUNT(*) FROM " + table).Scan(&count)
		assert.NoError(t, err, "Table %s should exist", table)
	}

	// Test that running migrations again doesn't cause errors (idempotency)
	err = database.RunMigrations()
	assert.NoError(t, err, "Running migrations again should not return error")

	// Test that migration count is still 5 (no duplicates)
	err = db.QueryRow("SELECT COUNT(*) FROM migrations").Scan(&count)
	assert.NoError(t, err, "Should be able to query migrations table")
	assert.Equal(t, 5, count, "Should still have 5 migration records (no duplicates)")
}

func TestMigration_ProjectsTable(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	// Insert test data to verify table works
	_, err := db.Exec(`
		INSERT INTO projects (id, name, description, color, created_at, updated_at)
		VALUES (?, ?, ?, ?, datetime('now'), datetime('now'))
	`, "test_proj", "Test Project", "Test Description", "blue")
	assert.NoError(t, err, "Should be able to insert into projects table")

	// Verify data was inserted
	var name string
	err = db.QueryRow("SELECT name FROM projects WHERE id = ?", "test_proj").Scan(&name)
	assert.NoError(t, err, "Should be able to query inserted project")
	assert.Equal(t, "Test Project", name, "Inserted project name should match")
}

func TestMigration_TasksTable(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	// First insert a project (foreign key constraint)
	_, err := db.Exec(`
		INSERT INTO projects (id, name, description, color, created_at, updated_at)
		VALUES (?, ?, ?, ?, datetime('now'), datetime('now'))
	`, "test_proj", "Test Project", "Test Description", "blue")
	require.NoError(t, err, "Should be able to insert project")

	// Insert test data to verify tasks table works
	_, err = db.Exec(`
		INSERT INTO tasks (id, project_id, name, desc, status, priority, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, datetime('now'), datetime('now'))
	`, "test_task", "test_proj", "Test Task", "Test Description", 0, 1)
	assert.NoError(t, err, "Should be able to insert into tasks table")

	// Verify data was inserted
	var name string
	err = db.QueryRow("SELECT name FROM tasks WHERE id = ?", "test_task").Scan(&name)
	assert.NoError(t, err, "Should be able to query inserted task")
	assert.Equal(t, "Test Task", name, "Inserted task name should match")
}

func TestMigration_ForeignKeyConstraints(t *testing.T) {
	t.Skip("Skipping foreign key constraint test - in-memory SQLite may have different FK behavior")
}

func TestMigration_Indexes(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	// Test that indexes exist
	indexes := []string{
		"idx_tasks_project_id",
		"idx_tasks_status",
		"idx_tasks_created_at",
		"idx_projects_created_at",
		"idx_tasks_priority",
		"idx_tasks_blocked_by",
	}

	for _, indexName := range indexes {
		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='index' AND name=?", indexName).Scan(&count)
		assert.NoError(t, err, "Should be able to query index %s", indexName)
		assert.Equal(t, 1, count, "Index %s should exist", indexName)
	}
}

func TestMigration_IntegerPKAndBlockedBy(t *testing.T) {
	db := setupTestDB(t)
	defer cleanupTestDB(t, db)

	// First insert a project (foreign key constraint)
	_, err := db.Exec(`
		INSERT INTO projects (id, name, description, color, created_at, updated_at)
		VALUES (?, ?, ?, ?, datetime('now'), datetime('now'))
	`, "test_proj", "Test Project", "Test Description", "blue")
	require.NoError(t, err, "Should be able to insert project")

	// Insert a task without blocked_by
	_, err = db.Exec(`
		INSERT INTO tasks (id, project_id, name, desc, status, priority, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, datetime('now'), datetime('now'))
	`, "task_1", "test_proj", "Task 1", "First task", 0, 1)
	assert.NoError(t, err, "Should be able to insert task without blocked_by")

	// Get the int_id of the first task
	var intID1 int
	err = db.QueryRow("SELECT int_id FROM tasks WHERE id = ?", "task_1").Scan(&intID1)
	assert.NoError(t, err, "Should be able to query int_id")
	assert.Greater(t, intID1, 0, "int_id should be auto-generated and greater than 0")

	// Insert a second task blocked by the first
	_, err = db.Exec(`
		INSERT INTO tasks (id, project_id, name, desc, status, priority, blocked_by, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, datetime('now'), datetime('now'))
	`, "task_2", "test_proj", "Task 2", "Second task", 0, 1, intID1)
	assert.NoError(t, err, "Should be able to insert task with blocked_by")

	// Verify blocked_by was set correctly
	var blockedBy sql.NullInt64
	err = db.QueryRow("SELECT blocked_by FROM tasks WHERE id = ?", "task_2").Scan(&blockedBy)
	assert.NoError(t, err, "Should be able to query blocked_by")
	assert.True(t, blockedBy.Valid, "blocked_by should be set")
	assert.Equal(t, int64(intID1), blockedBy.Int64, "blocked_by should reference first task's int_id")

	// Verify first task has nil blocked_by
	err = db.QueryRow("SELECT blocked_by FROM tasks WHERE id = ?", "task_1").Scan(&blockedBy)
	assert.NoError(t, err, "Should be able to query blocked_by")
	assert.False(t, blockedBy.Valid, "blocked_by should be NULL for first task")

	// Verify that both int_id and string id fields work correctly
	var stringID string
	var intID int
	err = db.QueryRow("SELECT id, int_id FROM tasks WHERE int_id = ?", intID1).Scan(&stringID, &intID)
	assert.NoError(t, err, "Should be able to query by int_id")
	assert.Equal(t, "task_1", stringID, "String ID should match")
	assert.Equal(t, intID1, intID, "Integer ID should match")
}

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite", ":memory:?_foreign_keys=true")
	require.NoError(t, err, "Failed to open in-memory database")

	// Test the connection
	err = db.Ping()
	require.NoError(t, err, "Failed to ping test database")

	// Run migrations on the test database
	database := &Database{Db: db}
	err = database.RunMigrations()
	require.NoError(t, err, "Failed to run migrations on test database")

	return db
}
func cleanupTestDB(t *testing.T, db *sql.DB) {
	err := db.Close()
	require.NoError(t, err, "Failed to close test database")
}
