package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"kahn/internal/domain"
)

// CreateTask Tests

func TestCreateTask_Success(t *testing.T) {
	km, cleanup := setupTestApp(t)
	defer cleanup()

	// Create a task
	err := km.CreateTask("New Task", "Description")
	require.NoError(t, err)

	// Verify task was created
	activeProj := km.projectManager.GetActiveProject()
	assert.Len(t, activeProj.Tasks, 1)
	assert.Equal(t, "New Task", activeProj.Tasks[0].Name)
	assert.Equal(t, "Description", activeProj.Tasks[0].Desc)
	assert.Equal(t, domain.Low, activeProj.Tasks[0].Priority)
	assert.Equal(t, domain.NotStarted, activeProj.Tasks[0].Status)
}

func TestCreateTaskWithPriority_Success(t *testing.T) {
	km, cleanup := setupTestApp(t)
	defer cleanup()

	// Create a task with high priority
	err := km.CreateTaskWithPriority("High Priority Task", "Description", domain.High)
	require.NoError(t, err)

	// Verify task was created with correct priority
	activeProj := km.projectManager.GetActiveProject()
	assert.Len(t, activeProj.Tasks, 1)
	assert.Equal(t, "High Priority Task", activeProj.Tasks[0].Name)
	assert.Equal(t, domain.High, activeProj.Tasks[0].Priority)
}

func TestCreateTask_EmptyName_ReturnsError(t *testing.T) {
	km, cleanup := setupTestApp(t)
	defer cleanup()

	// Try to create task with empty name
	err := km.CreateTask("", "Description")
	assert.Error(t, err)

	// Verify task was not created
	activeProj := km.projectManager.GetActiveProject()
	assert.Len(t, activeProj.Tasks, 0)
}

func TestCreateTask_WhitespaceName_ReturnsError(t *testing.T) {
	km, cleanup := setupTestApp(t)
	defer cleanup()

	// Try to create task with whitespace-only name
	err := km.CreateTask("   ", "Description")
	assert.Error(t, err)

	// Verify task was not created
	activeProj := km.projectManager.GetActiveProject()
	assert.Len(t, activeProj.Tasks, 0)
}

func TestCreateTask_MultipleTasks(t *testing.T) {
	km, cleanup := setupTestApp(t)
	defer cleanup()

	// Create multiple tasks
	err := km.CreateTask("Task 1", "Description 1")
	require.NoError(t, err)

	err = km.CreateTaskWithPriority("Task 2", "Description 2", domain.High)
	require.NoError(t, err)

	err = km.CreateTaskWithPriority("Task 3", "Description 3", domain.Medium)
	require.NoError(t, err)

	// Verify all tasks were created
	activeProj := km.projectManager.GetActiveProject()
	assert.Len(t, activeProj.Tasks, 3)

	// Tasks are stored in reverse chronological order (newest first)
	assert.Equal(t, "Task 3", activeProj.Tasks[0].Name)
	assert.Equal(t, "Task 2", activeProj.Tasks[1].Name)
	assert.Equal(t, "Task 1", activeProj.Tasks[2].Name)
}

// UpdateTask Tests

func TestUpdateTask_Success(t *testing.T) {
	km, cleanup := setupTestApp(t)
	defer cleanup()

	// Create a task
	taskID := createTestTask(t, km, "Original Task", "Original Description")

	// Update the task
	err := km.UpdateTask(taskID, "Updated Task", "Updated Description", domain.High, domain.Bug)
	require.NoError(t, err)

	// Verify task was updated
	activeProj := km.projectManager.GetActiveProject()
	require.Len(t, activeProj.Tasks, 1)
	task := activeProj.Tasks[0]
	assert.Equal(t, "Updated Task", task.Name)
	assert.Equal(t, "Updated Description", task.Desc)
	assert.Equal(t, domain.High, task.Priority)
	assert.Equal(t, domain.Bug, task.Type)
}

func TestUpdateTask_NonExistentTask_ReturnsError(t *testing.T) {
	km, cleanup := setupTestApp(t)
	defer cleanup()

	// Try to update non-existent task
	err := km.UpdateTask("non-existent-id", "Updated Task", "Description", domain.High, domain.Bug)
	assert.Error(t, err)
}

func TestUpdateTask_EmptyName_ReturnsError(t *testing.T) {
	km, cleanup := setupTestApp(t)
	defer cleanup()

	// Create a task
	taskID := createTestTask(t, km, "Original Task", "Original Description")

	// Try to update with empty name
	err := km.UpdateTask(taskID, "", "Updated Description", domain.High, domain.Bug)
	assert.Error(t, err)

	// Verify task was not updated
	activeProj := km.projectManager.GetActiveProject()
	require.Len(t, activeProj.Tasks, 1)
	assert.Equal(t, "Original Task", activeProj.Tasks[0].Name)
}

func TestUpdateTask_PreservesStatus(t *testing.T) {
	km, cleanup := setupTestApp(t)
	defer cleanup()

	// Create a task and move it to InProgress
	taskID := createTestTask(t, km, "Test Task", "Description")
	err := km.MoveTaskToNextStatus(taskID)
	require.NoError(t, err)

	// Verify it's in InProgress
	activeProj := km.projectManager.GetActiveProject()
	require.Len(t, activeProj.Tasks, 1)
	assert.Equal(t, domain.InProgress, activeProj.Tasks[0].Status)

	// Update the task
	err = km.UpdateTask(taskID, "Updated Task", "Updated Description", domain.High, domain.Bug)
	require.NoError(t, err)

	// Verify status was preserved
	activeProj = km.projectManager.GetActiveProject()
	require.Len(t, activeProj.Tasks, 1)
	assert.Equal(t, domain.InProgress, activeProj.Tasks[0].Status)
}

// DeleteTask Tests

func TestDeleteTask_Success(t *testing.T) {
	km, cleanup := setupTestApp(t)
	defer cleanup()

	// Create a task
	taskID := createTestTask(t, km, "Task to Delete", "Description")

	// Delete the task
	err := km.DeleteTask(taskID)
	require.NoError(t, err)

	// Verify task was deleted
	activeProj := km.projectManager.GetActiveProject()
	assert.Len(t, activeProj.Tasks, 0)
}

func TestDeleteTask_NonExistentTask_ReturnsError(t *testing.T) {
	km, cleanup := setupTestApp(t)
	defer cleanup()

	// Try to delete non-existent task
	err := km.DeleteTask("non-existent-id")
	assert.Error(t, err)
}

func TestDeleteTask_MultipleTasksDeletesCorrectOne(t *testing.T) {
	km, cleanup := setupTestApp(t)
	defer cleanup()

	// Create multiple tasks
	taskID1 := createTestTask(t, km, "Task 1", "Description 1")
	taskID2 := createTestTask(t, km, "Task 2", "Description 2")
	createTestTask(t, km, "Task 3", "Description 3")

	// Delete the second task
	err := km.DeleteTask(taskID2)
	require.NoError(t, err)

	// Verify correct task was deleted
	activeProj := km.projectManager.GetActiveProject()
	assert.Len(t, activeProj.Tasks, 2)

	// Verify remaining tasks
	var remainingIDs []string
	for _, task := range activeProj.Tasks {
		remainingIDs = append(remainingIDs, task.ID)
	}
	assert.Contains(t, remainingIDs, taskID1)
	assert.NotContains(t, remainingIDs, taskID2)
}

// MoveTaskToNextStatus Tests

func TestMoveTaskToNextStatus_NotStartedToInProgress(t *testing.T) {
	km, cleanup := setupTestApp(t)
	defer cleanup()

	// Create a task (starts in NotStarted)
	taskID := createTestTask(t, km, "Test Task", "Description")

	// Move to next status
	err := km.MoveTaskToNextStatus(taskID)
	require.NoError(t, err)

	// Verify status changed to InProgress
	activeProj := km.projectManager.GetActiveProject()
	require.Len(t, activeProj.Tasks, 1)
	assert.Equal(t, domain.InProgress, activeProj.Tasks[0].Status)
}

func TestMoveTaskToNextStatus_InProgressToDone(t *testing.T) {
	km, cleanup := setupTestApp(t)
	defer cleanup()

	// Create a task and move to InProgress
	taskID := createTestTask(t, km, "Test Task", "Description")
	err := km.MoveTaskToNextStatus(taskID)
	require.NoError(t, err)

	// Move to next status (Done)
	err = km.MoveTaskToNextStatus(taskID)
	require.NoError(t, err)

	// Verify status changed to Done
	activeProj := km.projectManager.GetActiveProject()
	require.Len(t, activeProj.Tasks, 1)
	assert.Equal(t, domain.Done, activeProj.Tasks[0].Status)
}

func TestMoveTaskToNextStatus_DoneStaysDone(t *testing.T) {
	km, cleanup := setupTestApp(t)
	defer cleanup()

	// Create a task and move to Done
	taskID := createTestTask(t, km, "Test Task", "Description")
	km.MoveTaskToNextStatus(taskID) // NotStarted -> InProgress
	km.MoveTaskToNextStatus(taskID) // InProgress -> Done

	// Moving beyond Done cycles back to NotStarted
	err := km.MoveTaskToNextStatus(taskID)
	require.NoError(t, err)

	// Verify status cycled to NotStarted
	activeProj := km.projectManager.GetActiveProject()
	require.Len(t, activeProj.Tasks, 1)
	assert.Equal(t, domain.NotStarted, activeProj.Tasks[0].Status)
}

func TestMoveTaskToNextStatus_NonExistentTask_ReturnsError(t *testing.T) {
	km, cleanup := setupTestApp(t)
	defer cleanup()

	// Try to move non-existent task
	err := km.MoveTaskToNextStatus("non-existent-id")
	assert.Error(t, err)
}

// MoveTaskToPreviousStatus Tests

func TestMoveTaskToPreviousStatus_DoneToInProgress(t *testing.T) {
	km, cleanup := setupTestApp(t)
	defer cleanup()

	// Create a task and move to Done
	taskID := createTestTask(t, km, "Test Task", "Description")
	km.MoveTaskToNextStatus(taskID) // NotStarted -> InProgress
	km.MoveTaskToNextStatus(taskID) // InProgress -> Done

	// Move to previous status
	err := km.MoveTaskToPreviousStatus(taskID)
	require.NoError(t, err)

	// Verify status changed to InProgress
	activeProj := km.projectManager.GetActiveProject()
	require.Len(t, activeProj.Tasks, 1)
	assert.Equal(t, domain.InProgress, activeProj.Tasks[0].Status)
}

func TestMoveTaskToPreviousStatus_InProgressToNotStarted(t *testing.T) {
	km, cleanup := setupTestApp(t)
	defer cleanup()

	// Create a task and move to InProgress
	taskID := createTestTask(t, km, "Test Task", "Description")
	km.MoveTaskToNextStatus(taskID) // NotStarted -> InProgress

	// Move to previous status
	err := km.MoveTaskToPreviousStatus(taskID)
	require.NoError(t, err)

	// Verify status changed to NotStarted
	activeProj := km.projectManager.GetActiveProject()
	require.Len(t, activeProj.Tasks, 1)
	assert.Equal(t, domain.NotStarted, activeProj.Tasks[0].Status)
}

func TestMoveTaskToPreviousStatus_NotStartedStaysNotStarted(t *testing.T) {
	km, cleanup := setupTestApp(t)
	defer cleanup()

	// Create a task (starts in NotStarted)
	taskID := createTestTask(t, km, "Test Task", "Description")

	// Moving before NotStarted cycles back to Done
	err := km.MoveTaskToPreviousStatus(taskID)
	require.NoError(t, err)

	// Verify status cycled to Done
	activeProj := km.projectManager.GetActiveProject()
	require.Len(t, activeProj.Tasks, 1)
	assert.Equal(t, domain.Done, activeProj.Tasks[0].Status)
}

func TestMoveTaskToPreviousStatus_NonExistentTask_ReturnsError(t *testing.T) {
	km, cleanup := setupTestApp(t)
	defer cleanup()

	// Try to move non-existent task
	err := km.MoveTaskToPreviousStatus("non-existent-id")
	assert.Error(t, err)
}

// Search Persistence Tests

func TestCreateTask_PreservesSearch(t *testing.T) {
	km, cleanup := setupTestApp(t)
	defer cleanup()

	// Activate search
	km.searchState.Activate()
	km.searchState.AppendChar("test")
	assert.True(t, km.searchState.IsActive())
	assert.Equal(t, "test", km.searchState.GetQuery())

	// Create a task
	err := km.CreateTask("Test Task", "Description")
	require.NoError(t, err)

	// Verify search is still active
	assert.True(t, km.searchState.IsActive())
	assert.Equal(t, "test", km.searchState.GetQuery())
}

func TestDeleteTask_PreservesSearch(t *testing.T) {
	km, cleanup := setupTestApp(t)
	defer cleanup()

	// Create tasks
	taskID := createTestTask(t, km, "Task to Delete", "Description")
	createTestTask(t, km, "Another Task", "Description")

	// Activate search
	km.searchState.Activate()
	km.searchState.AppendChar("task")
	assert.True(t, km.searchState.IsActive())

	// Delete a task
	err := km.DeleteTask(taskID)
	require.NoError(t, err)

	// Verify search is still active
	assert.True(t, km.searchState.IsActive())
	assert.Equal(t, "task", km.searchState.GetQuery())
}

func TestMoveTaskStatus_PreservesSearch(t *testing.T) {
	km, cleanup := setupTestApp(t)
	defer cleanup()

	// Create a task
	taskID := createTestTask(t, km, "Test Task", "Description")

	// Activate search
	km.searchState.Activate()
	km.searchState.AppendChar("test")
	assert.True(t, km.searchState.IsActive())

	// Move task status
	err := km.MoveTaskToNextStatus(taskID)
	require.NoError(t, err)

	// Verify search is still active
	assert.True(t, km.searchState.IsActive())
	assert.Equal(t, "test", km.searchState.GetQuery())
}

// Edge Cases

func TestTaskOperations_EmptyProject(t *testing.T) {
	km, cleanup := setupTestApp(t)
	defer cleanup()

	// Verify empty project
	activeProj := km.projectManager.GetActiveProject()
	assert.Len(t, activeProj.Tasks, 0)

	// Operations on non-existent tasks should error
	err := km.UpdateTask("fake-id", "Name", "Desc", domain.Low, domain.RegularTask)
	assert.Error(t, err)

	err = km.DeleteTask("fake-id")
	assert.Error(t, err)

	err = km.MoveTaskToNextStatus("fake-id")
	assert.Error(t, err)

	err = km.MoveTaskToPreviousStatus("fake-id")
	assert.Error(t, err)
}

func TestTaskOperations_MultipleStatusTransitions(t *testing.T) {
	km, cleanup := setupTestApp(t)
	defer cleanup()

	// Create a task
	taskID := createTestTask(t, km, "Test Task", "Description")

	// Move forward through all statuses
	err := km.MoveTaskToNextStatus(taskID) // -> InProgress
	require.NoError(t, err)

	err = km.MoveTaskToNextStatus(taskID) // -> Done
	require.NoError(t, err)

	activeProj := km.projectManager.GetActiveProject()
	assert.Equal(t, domain.Done, activeProj.Tasks[0].Status)

	// Move backward through all statuses
	err = km.MoveTaskToPreviousStatus(taskID) // -> InProgress
	require.NoError(t, err)

	err = km.MoveTaskToPreviousStatus(taskID) // -> NotStarted
	require.NoError(t, err)

	activeProj = km.projectManager.GetActiveProject()
	assert.Equal(t, domain.NotStarted, activeProj.Tasks[0].Status)
}
