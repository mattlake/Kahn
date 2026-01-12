package app

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/require"

	"kahn/internal/config"
	"kahn/internal/database"
	"kahn/internal/domain"
)

// setupTestApp creates a fully initialized KahnModel for testing with an in-memory database
func setupTestApp(t *testing.T) (*KahnModel, func()) {
	t.Helper()

	// Create in-memory database configuration
	cfg := &config.Config{}
	cfg.Database.Path = ":memory:"
	cfg.Database.BusyTimeout = 5000
	cfg.Database.JournalMode = "WAL"
	cfg.Database.CacheSize = 10000
	cfg.Database.ForeignKeys = true

	// Initialize database
	db, err := database.NewDatabase(cfg)
	require.NoError(t, err, "Failed to create test database")

	// Create KahnModel
	model := NewKahnModel(db, "test-version")
	require.NotNil(t, model, "NewKahnModel should not return nil")

	// Cleanup function
	cleanup := func() {
		db.Close()
	}

	return model, cleanup
}

// setupTestAppWithData creates a KahnModel with predefined test data
func setupTestAppWithData(t *testing.T, projectCount, taskCount int) (*KahnModel, func()) {
	t.Helper()

	km, cleanup := setupTestApp(t)

	// Create additional projects if requested
	for i := 1; i < projectCount; i++ {
		err := km.projectManager.CreateProject(
			"Test Project "+string(rune('A'+i)),
			"Test description",
		)
		require.NoError(t, err, "Failed to create test project")
	}

	// Create tasks in the active project
	if taskCount > 0 {
		for i := 0; i < taskCount; i++ {
			err := km.CreateTask(
				"Test Task "+string(rune('A'+i)),
				"Test task description",
			)
			require.NoError(t, err, "Failed to create test task")
		}
	}

	return km, cleanup
}

// createTestTask is a helper to create a task in the app for testing
func createTestTask(t *testing.T, km *KahnModel, name, desc string) string {
	t.Helper()

	activeProj := km.projectManager.GetActiveProject()
	require.NotNil(t, activeProj, "Active project should not be nil")

	task, err := km.taskService.CreateTask(name, desc, activeProj.ID, domain.RegularTask, domain.Medium, nil)
	require.NoError(t, err, "Failed to create test task")

	// Refresh task lists
	km.RefreshTasksWithSearch()

	return task.ID
}

// createTestTaskWithPriority is a helper to create a task with specific priority
func createTestTaskWithPriority(t *testing.T, km *KahnModel, name, desc string, priority domain.Priority) string {
	t.Helper()

	activeProj := km.projectManager.GetActiveProject()
	require.NotNil(t, activeProj, "Active project should not be nil")

	task, err := km.taskService.CreateTask(name, desc, activeProj.ID, domain.RegularTask, priority, nil)
	require.NoError(t, err, "Failed to create test task with priority")

	// Refresh task lists
	km.RefreshTasksWithSearch()

	return task.ID
}

// createTestProject is a helper to create a project in the app for testing
func createTestProject(t *testing.T, km *KahnModel, name, desc string) string {
	t.Helper()

	project, err := km.projectService.CreateProject(name, desc)
	require.NoError(t, err, "Failed to create test project")

	return project.ID
}

// moveTaskToStatus is a helper to move a task to a specific status
func moveTaskToStatus(t *testing.T, km *KahnModel, taskID string, status domain.Status) {
	t.Helper()

	task, err := km.taskService.UpdateTaskStatus(taskID, status)
	require.NoError(t, err, "Failed to move task to %v", status)

	// Refresh the task lists
	km.RefreshTasksWithSearch()

	// Verify the update
	require.Equal(t, status, task.Status, "Task status should be updated")
}

// simulateKeyPress simulates a keyboard input and returns the updated model
func simulateKeyPress(km *KahnModel, keyString string) (tea.Model, tea.Cmd) {
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(keyString)}
	return km.Update(msg)
}

// simulateKeyType simulates a specific key type
func simulateKeyType(km *KahnModel, keyType tea.KeyType) (tea.Model, tea.Cmd) {
	msg := tea.KeyMsg{Type: keyType}
	return km.Update(msg)
}

// assertViewState is a helper to assert the current view state
func assertViewState(t *testing.T, km *KahnModel, expected ViewState) {
	t.Helper()

	actual := km.uiStateManager.GetCurrentViewState()
	require.Equal(t, expected, actual, "Expected view state %v, got %v", expected, actual)
}

// assertFormError is a helper to assert form error state
func assertFormError(t *testing.T, km *KahnModel, expectedMsg string) {
	t.Helper()

	formError, _ := km.uiStateManager.FormState().GetError()
	require.Contains(t, formError, expectedMsg, "Expected form error to contain '%s', got '%s'", expectedMsg, formError)
}

// assertNoFormError is a helper to assert no form error
func assertNoFormError(t *testing.T, km *KahnModel) {
	t.Helper()

	formError, _ := km.uiStateManager.FormState().GetError()
	require.Empty(t, formError, "Expected no form error, got '%s'", formError)
}

// assertActiveProject is a helper to assert the active project
func assertActiveProject(t *testing.T, km *KahnModel, expectedID string) {
	t.Helper()

	activeProj := km.projectManager.GetActiveProject()
	require.NotNil(t, activeProj, "Active project should not be nil")
	require.Equal(t, expectedID, activeProj.ID, "Expected active project %s, got %s", expectedID, activeProj.ID)
}

// assertTaskCount is a helper to assert the number of tasks in a status
func assertTaskCount(t *testing.T, km *KahnModel, status domain.Status, expectedCount int) {
	t.Helper()

	activeProj := km.projectManager.GetActiveProject()
	require.NotNil(t, activeProj, "Active project should not be nil")

	tasks := activeProj.GetTasksByStatus(status)
	require.Len(t, tasks, expectedCount, "Expected %d tasks in status %v, got %d", expectedCount, status, len(tasks))
}

// assertSearchActive is a helper to assert search state
func assertSearchActive(t *testing.T, km *KahnModel, expected bool) {
	t.Helper()

	actual := km.searchState.IsActive()
	require.Equal(t, expected, actual, "Expected search active=%v, got %v", expected, actual)
}

// assertSearchQuery is a helper to assert search query value
func assertSearchQuery(t *testing.T, km *KahnModel, expected string) {
	t.Helper()

	actual := km.searchState.GetQuery()
	require.Equal(t, expected, actual, "Expected search query '%s', got '%s'", expected, actual)
}

// assertProjectCount is a helper to assert the number of projects
func assertProjectCount(t *testing.T, km *KahnModel, expectedCount int) {
	t.Helper()

	actualCount := km.projectManager.GetProjectCount()
	require.Equal(t, expectedCount, actualCount, "Expected %d projects, got %d", expectedCount, actualCount)
}
