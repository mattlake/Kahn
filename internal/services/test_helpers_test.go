package services

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"kahn/internal/domain"
)

func TestMockTaskRepository_GetByStatus_Ordering(t *testing.T) {
	repo := NewMockTaskRepository()
	projectID := "test_project"

	// Create test tasks with different priorities and creation times
	now := time.Now()
	tasks := []struct {
		name      string
		status    domain.Status
		priority  domain.Priority
		createdAt time.Time
		updatedAt time.Time
	}{
		// Not Started tasks with different priorities and creation times
		{"Low Priority Old", domain.NotStarted, domain.Low, now.Add(-5 * time.Hour), now.Add(-5 * time.Hour)},
		{"High Priority Old", domain.NotStarted, domain.High, now.Add(-4 * time.Hour), now.Add(-4 * time.Hour)},
		{"Medium Priority Old", domain.NotStarted, domain.Medium, now.Add(-3 * time.Hour), now.Add(-3 * time.Hour)},
		{"Low Priority New", domain.NotStarted, domain.Low, now.Add(-2 * time.Hour), now.Add(-2 * time.Hour)},
		{"High Priority New", domain.NotStarted, domain.High, now.Add(-1 * time.Hour), now.Add(-1 * time.Hour)},

		// In Progress tasks with different update times
		{"In Progress Oldest", domain.InProgress, domain.Low, now.Add(-5 * time.Hour), now.Add(-5 * time.Hour)},
		{"In Progress Middle", domain.InProgress, domain.Medium, now.Add(-4 * time.Hour), now.Add(-3 * time.Hour)},
		{"In Progress Newest", domain.InProgress, domain.High, now.Add(-3 * time.Hour), now.Add(-1 * time.Hour)},

		// Done tasks with different update times
		{"Done Oldest", domain.Done, domain.Low, now.Add(-5 * time.Hour), now.Add(-5 * time.Hour)},
		{"Done Middle", domain.Done, domain.Medium, now.Add(-4 * time.Hour), now.Add(-3 * time.Hour)},
		{"Done Newest", domain.Done, domain.High, now.Add(-3 * time.Hour), now.Add(-1 * time.Hour)},
	}

	// Insert test tasks
	for i, taskData := range tasks {
		task := &domain.Task{
			ID:        "task_" + string(rune(i)),
			ProjectID: projectID,
			Name:      taskData.name,
			Status:    taskData.status,
			Priority:  taskData.priority,
			CreatedAt: taskData.createdAt,
			UpdatedAt: taskData.updatedAt,
		}
		err := repo.Create(task)
		require.NoError(t, err)
	}

	// Test Not Started ordering: priority DESC, then created_at ASC
	notStartedTasks, err := repo.GetByStatus(projectID, domain.NotStarted)
	require.NoError(t, err)
	assert.Len(t, notStartedTasks, 5, "Should have 5 NotStarted tasks")

	// Verify order: High Priority tasks first (oldest creation time first), then Medium, then Low
	expectedOrder := []string{
		"High Priority Old",   // High priority, oldest
		"High Priority New",   // High priority, newer
		"Medium Priority Old", // Medium priority, oldest
		"Low Priority Old",    // Low priority, oldest
		"Low Priority New",    // Low priority, newer
	}

	for i, expectedName := range expectedOrder {
		assert.Equal(t, expectedName, notStartedTasks[i].Name,
			"Task at position %d should be %s", i, expectedName)
	}

	// Test In Progress ordering: updated_at DESC
	inProgressTasks, err := repo.GetByStatus(projectID, domain.InProgress)
	require.NoError(t, err)
	assert.Len(t, inProgressTasks, 3, "Should have 3 InProgress tasks")

	expectedInProgressOrder := []string{
		"In Progress Newest", // Updated 1 hour ago (newest)
		"In Progress Middle", // Updated 3 hours ago
		"In Progress Oldest", // Updated 5 hours ago (oldest)
	}

	for i, expectedName := range expectedInProgressOrder {
		assert.Equal(t, expectedName, inProgressTasks[i].Name,
			"InProgress task at position %d should be %s", i, expectedName)
	}

	// Test Done ordering: updated_at DESC
	doneTasks, err := repo.GetByStatus(projectID, domain.Done)
	require.NoError(t, err)
	assert.Len(t, doneTasks, 3, "Should have 3 Done tasks")

	expectedDoneOrder := []string{
		"Done Newest", // Updated 1 hour ago (newest)
		"Done Middle", // Updated 3 hours ago
		"Done Oldest", // Updated 5 hours ago (oldest)
	}

	for i, expectedName := range expectedDoneOrder {
		assert.Equal(t, expectedName, doneTasks[i].Name,
			"Done task at position %d should be %s", i, expectedName)
	}
}

func TestMockTaskRepository_UpdateStatus_UpdatedAt(t *testing.T) {
	repo := NewMockTaskRepository()
	projectID := "test_project"

	// Create initial task
	initialTime := time.Now().Add(-1 * time.Hour).UTC()
	task := &domain.Task{
		ID:        "test_task",
		ProjectID: projectID,
		Name:      "Test Task",
		Status:    domain.NotStarted,
		Priority:  domain.Low,
		CreatedAt: initialTime,
		UpdatedAt: initialTime,
	}
	err := repo.Create(task)
	require.NoError(t, err)

	// Get task before update to compare
	taskBeforeUpdate, err := repo.GetByID("test_task")
	require.NoError(t, err)
	t.Logf("Initial UpdatedAt: %v (diff from now: %v)", taskBeforeUpdate.UpdatedAt, time.Since(taskBeforeUpdate.UpdatedAt))

	// Wait a bit to ensure timestamp difference
	time.Sleep(100 * time.Millisecond)

	// Update status
	err = repo.UpdateStatus("test_task", domain.InProgress)
	require.NoError(t, err)

	// Verify UpdatedAt was updated
	updatedTask, err := repo.GetByID("test_task")
	require.NoError(t, err)

	// Check if time difference is sufficient (allow for some timing precision issues)
	timeDiff := updatedTask.UpdatedAt.Sub(taskBeforeUpdate.UpdatedAt)
	assert.True(t, timeDiff > 0,
		"UpdatedAt should be updated after status change. Diff: %v", timeDiff)
	assert.Equal(t, domain.InProgress, updatedTask.Status, "Status should be updated")

	// Wait a bit more
	time.Sleep(100 * time.Millisecond)

	// Update status again
	err = repo.UpdateStatus("test_task", domain.Done)
	require.NoError(t, err)

	// Verify UpdatedAt was updated again
	doneTask, err := repo.GetByID("test_task")
	require.NoError(t, err)

	timeDiff2 := doneTask.UpdatedAt.Sub(updatedTask.UpdatedAt)
	assert.True(t, timeDiff2 > 0,
		"UpdatedAt should be updated again after second status change. Diff: %v", timeDiff2)
	assert.Equal(t, domain.Done, doneTask.Status, "Status should be updated to Done")
}

func TestMockTaskRepository_Update_UpdatedAt(t *testing.T) {
	repo := NewMockTaskRepository()
	projectID := "test_project"

	// Create initial task
	initialTime := time.Now().Add(-1 * time.Hour).UTC()
	task := &domain.Task{
		ID:        "test_task",
		ProjectID: projectID,
		Name:      "Test Task",
		Desc:      "Original Description",
		Status:    domain.NotStarted,
		Priority:  domain.Low,
		CreatedAt: initialTime,
		UpdatedAt: initialTime,
	}
	err := repo.Create(task)
	require.NoError(t, err)

	// Get task before update to compare
	taskBeforeUpdate, err := repo.GetByID("test_task")
	require.NoError(t, err)

	// Wait a bit to ensure timestamp difference
	time.Sleep(100 * time.Millisecond)

	// Update task
	task.Name = "Updated Task"
	task.Desc = "Updated Description"
	task.Priority = domain.High
	err = repo.Update(task)
	require.NoError(t, err)

	// Verify UpdatedAt was updated
	updatedTask, err := repo.GetByID("test_task")
	require.NoError(t, err)

	timeDiff := updatedTask.UpdatedAt.Sub(taskBeforeUpdate.UpdatedAt)
	assert.True(t, timeDiff > 0,
		"UpdatedAt should be updated after task update. Diff: %v", timeDiff)
	assert.Equal(t, "Updated Task", updatedTask.Name, "Name should be updated")
	assert.Equal(t, "Updated Description", updatedTask.Desc, "Description should be updated")
	assert.Equal(t, domain.High, updatedTask.Priority, "Priority should be updated")
}
