package main

import (
	"database/sql"
	"testing"
	"time"

	"github.com/charmbracelet/bubbletea"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
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

// createTestModelWithTasks creates a model with predefined tasks for testing
func createTestModelWithTasks(t *testing.T, taskNames []string, statuses []Status) *Model {
	config := createTestConfig()
	database, err := NewDatabase(config)
	require.NoError(t, err, "Failed to create test database")

	model := NewModel(database)
	require.NotNil(t, model, "NewModel should not return nil")

	// Create tasks for testing
	if len(taskNames) == len(statuses) && len(taskNames) > 0 {
		activeProj := model.GetActiveProject()
		require.NotNil(t, activeProj, "Should have active project")

		for i, taskName := range taskNames {
			task := NewTask(taskName, "Test description", activeProj.ID)
			task.Status = statuses[i]

			// Save to database
			err := model.taskDAO.Create(task)
			require.NoError(t, err, "Failed to create test task")

			// Add to project in memory
			activeProj.AddTask(*task)
		}

		model.updateTaskLists()
	}

	return model
}

// simulateKeyPress simulates a key press on the model
func simulateKeyPress(t *testing.T, model *Model, key string) *Model {
	var keyMsg tea.KeyMsg

	// Handle special keys
	switch key {
	case "tab":
		keyMsg = tea.KeyMsg{Type: tea.KeyTab}
	case "enter":
		keyMsg = tea.KeyMsg{Type: tea.KeyEnter}
	case "esc":
		keyMsg = tea.KeyMsg{Type: tea.KeyEsc}
	case "backspace":
		keyMsg = tea.KeyMsg{Type: tea.KeyBackspace}
	default:
		// Handle regular character keys
		if len(key) > 0 {
			keyMsg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{rune(key[0])}}
		}
	}

	newModel, cmd := model.Update(keyMsg)
	require.Nil(t, cmd, "Command should be nil for simple key press")

	// Type assertion to convert tea.Model back to Model (not *Model)
	resultModel, ok := newModel.(Model)
	require.True(t, ok, "Model should be of type Model")

	// Return pointer to the new model
	return &resultModel
}

// assertTaskDeletionState asserts the model is in correct deletion state
func assertTaskDeletionState(t *testing.T, model *Model, inConfirmation bool, taskID string) {
	assert.Equal(t, inConfirmation, model.showTaskDeleteConfirm, "Task deletion confirmation state should match")
	if inConfirmation {
		assert.Equal(t, taskID, model.taskToDelete, "Task to delete should match")
	} else {
		assert.Equal(t, "", model.taskToDelete, "Task to delete should be empty when not in confirmation")
	}
}

// assertTaskNotInLists asserts a task is not present in any task list
func assertTaskNotInLists(t *testing.T, model *Model, taskID string) {
	// Check all three status lists
	for status := NotStarted; status <= Done; status++ {
		items := model.Tasks[status].Items()
		for _, item := range items {
			if task, ok := item.(Task); ok {
				assert.NotEqual(t, taskID, task.ID, "Task should not be found in %s list", status.ToString())
			}
		}
	}
}

// assertTaskInList asserts a task is present in a specific status list
func assertTaskInList(t *testing.T, model *Model, taskID string, status Status) {
	items := model.Tasks[status].Items()
	found := false
	for _, item := range items {
		if task, ok := item.(Task); ok {
			if task.ID == taskID {
				found = true
				break
			}
		}
	}
	assert.True(t, found, "Task should be found in %s list", status.ToString())
}
