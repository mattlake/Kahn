package repository

import (
	"database/sql"
	"testing"
	"time"

	"kahn/internal/database"
	"kahn/internal/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

// setupTestRepository creates an in-memory database with migrations and returns a repository
func setupTestRepository(t *testing.T) domain.TaskRepository {
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)

	// Run migrations
	dbWrapper := &database.Database{Db: db}
	err = dbWrapper.RunMigrations()
	require.NoError(t, err)

	return NewSQLiteTaskRepository(db)
}

// TaskTestData represents test data for task creation
type TaskTestData struct {
	Name      string
	Status    domain.Status
	TaskType  domain.TaskType
	Priority  domain.Priority
	CreatedAt time.Time
	UpdatedAt time.Time
}

// createTasksForStatus creates tasks with specific properties for testing status-based ordering
func createTasksForStatus(t *testing.T, repo domain.TaskRepository, projectID string, status domain.Status, tasks []TaskTestData) []string {
	var taskIDs []string

	for i, taskData := range tasks {
		task := &domain.Task{
			ID:        "task_" + string(rune(i)),
			ProjectID: projectID,
			Name:      taskData.Name,
			Status:    taskData.Status,
			Type:      taskData.TaskType,
			Priority:  taskData.Priority,
			CreatedAt: taskData.CreatedAt,
			UpdatedAt: taskData.UpdatedAt,
		}
		err := repo.Create(task)
		require.NoError(t, err)
		taskIDs = append(taskIDs, task.ID)
	}

	return taskIDs
}

// getNotStartedTestData returns test data for NotStarted status ordering tests
func getNotStartedTestData(baseTime time.Time) []TaskTestData {
	return []TaskTestData{
		{"Low Priority Old", domain.NotStarted, domain.RegularTask, domain.Low, baseTime.Add(-5 * time.Hour), baseTime.Add(-5 * time.Hour)},
		{"High Priority Old", domain.NotStarted, domain.Bug, domain.High, baseTime.Add(-4 * time.Hour), baseTime.Add(-4 * time.Hour)},
		{"Medium Priority Old", domain.NotStarted, domain.Feature, domain.Medium, baseTime.Add(-3 * time.Hour), baseTime.Add(-3 * time.Hour)},
		{"Low Priority New", domain.NotStarted, domain.RegularTask, domain.Low, baseTime.Add(-2 * time.Hour), baseTime.Add(-2 * time.Hour)},
		{"High Priority New", domain.NotStarted, domain.Bug, domain.High, baseTime.Add(-1 * time.Hour), baseTime.Add(-1 * time.Hour)},
	}
}

// getInProgressTestData returns test data for InProgress status ordering tests
func getInProgressTestData(baseTime time.Time) []TaskTestData {
	return []TaskTestData{
		{"In Progress Oldest", domain.InProgress, domain.Feature, domain.Low, baseTime.Add(-5 * time.Hour), baseTime.Add(-5 * time.Hour)},
		{"In Progress Middle", domain.InProgress, domain.RegularTask, domain.Medium, baseTime.Add(-4 * time.Hour), baseTime.Add(-3 * time.Hour)},
		{"In Progress Newest", domain.InProgress, domain.Bug, domain.High, baseTime.Add(-3 * time.Hour), baseTime.Add(-1 * time.Hour)},
	}
}

// getDoneTestData returns test data for Done status ordering tests
func getDoneTestData(baseTime time.Time) []TaskTestData {
	return []TaskTestData{
		{"Done Oldest", domain.Done, domain.Feature, domain.Low, baseTime.Add(-5 * time.Hour), baseTime.Add(-5 * time.Hour)},
		{"Done Middle", domain.Done, domain.RegularTask, domain.Medium, baseTime.Add(-4 * time.Hour), baseTime.Add(-3 * time.Hour)},
		{"Done Newest", domain.Done, domain.Bug, domain.High, baseTime.Add(-3 * time.Hour), baseTime.Add(-1 * time.Hour)},
	}
}

func TestTaskRepository_GetByStatus_NotStarted_Ordering(t *testing.T) {
	// Arrange
	repo := setupTestRepository(t)
	projectID := "test_project"

	now := time.Now()
	taskData := getNotStartedTestData(now)
	createTasksForStatus(t, repo, projectID, domain.NotStarted, taskData)

	// Act
	notStartedTasks, err := repo.GetByStatus(projectID, domain.NotStarted)

	// Assert
	require.NoError(t, err)
	assert.Len(t, notStartedTasks, 5, "Should have 5 NotStarted tasks")

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
}

func TestTaskRepository_GetByStatus_InProgress_Ordering(t *testing.T) {
	// Arrange
	repo := setupTestRepository(t)
	projectID := "test_project"

	now := time.Now()
	taskData := getInProgressTestData(now)
	createTasksForStatus(t, repo, projectID, domain.InProgress, taskData)

	// Act
	inProgressTasks, err := repo.GetByStatus(projectID, domain.InProgress)

	// Assert
	require.NoError(t, err)
	assert.Len(t, inProgressTasks, 3, "Should have 3 InProgress tasks")

	expectedOrder := []string{
		"In Progress Newest", // Updated 1 hour ago (newest)
		"In Progress Middle", // Updated 3 hours ago
		"In Progress Oldest", // Updated 5 hours ago (oldest)
	}

	for i, expectedName := range expectedOrder {
		assert.Equal(t, expectedName, inProgressTasks[i].Name,
			"InProgress task at position %d should be %s", i, expectedName)
	}
}

func TestTaskRepository_GetByStatus_Done_Ordering(t *testing.T) {
	// Arrange
	repo := setupTestRepository(t)
	projectID := "test_project"

	now := time.Now()
	taskData := getDoneTestData(now)
	createTasksForStatus(t, repo, projectID, domain.Done, taskData)

	// Act
	doneTasks, err := repo.GetByStatus(projectID, domain.Done)

	// Assert
	require.NoError(t, err)
	assert.Len(t, doneTasks, 3, "Should have 3 Done tasks")

	expectedOrder := []string{
		"Done Newest", // Updated 1 hour ago (newest)
		"Done Middle", // Updated 3 hours ago
		"Done Oldest", // Updated 5 hours ago (oldest)
	}

	for i, expectedName := range expectedOrder {
		assert.Equal(t, expectedName, doneTasks[i].Name,
			"Done task at position %d should be %s", i, expectedName)
	}
}

func TestTaskRepository_UpdateStatus_UpdatedAt(t *testing.T) {
	// Setup in-memory database with migrations
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer db.Close()

	// Run migrations
	dbWrapper := &database.Database{Db: db}
	err = dbWrapper.RunMigrations()
	require.NoError(t, err)

	// Create repository
	repo := NewSQLiteTaskRepository(db)
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
	err = repo.Create(task)
	require.NoError(t, err)

	// Verify initial UpdatedAt
	initialTask, err := repo.GetByID("test_task")
	require.NoError(t, err)
	assert.WithinDuration(t, initialTime, initialTask.UpdatedAt, time.Second, "Initial UpdatedAt should be preserved")

	// Wait a bit to ensure timestamp difference
	time.Sleep(10 * time.Millisecond)

	// Update status
	err = repo.UpdateStatus("test_task", domain.InProgress)
	require.NoError(t, err)

	// Verify UpdatedAt was updated
	updatedTask, err := repo.GetByID("test_task")
	require.NoError(t, err)
	assert.True(t, updatedTask.UpdatedAt.After(initialTask.UpdatedAt),
		"UpdatedAt should be updated after status change")
	assert.Equal(t, domain.InProgress, updatedTask.Status, "Status should be updated")
}

func TestTaskRepository_Update_UpdatedAt(t *testing.T) {
	// Setup in-memory database with migrations
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer db.Close()

	// Run migrations
	dbWrapper := &database.Database{Db: db}
	err = dbWrapper.RunMigrations()
	require.NoError(t, err)

	// Create repository
	repo := NewSQLiteTaskRepository(db)
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
	err = repo.Create(task)
	require.NoError(t, err)

	// Wait a bit to ensure timestamp difference
	time.Sleep(10 * time.Millisecond)

	// Update task
	task.Name = "Updated Task"
	task.Desc = "Updated Description"
	task.Priority = domain.High
	err = repo.Update(task)
	require.NoError(t, err)

	// Verify UpdatedAt was updated
	updatedTask, err := repo.GetByID("test_task")
	require.NoError(t, err)
	assert.True(t, updatedTask.UpdatedAt.After(initialTime),
		"UpdatedAt should be updated after task update")
	assert.Equal(t, "Updated Task", updatedTask.Name, "Name should be updated")
	assert.Equal(t, "Updated Description", updatedTask.Desc, "Description should be updated")
	assert.Equal(t, domain.High, updatedTask.Priority, "Priority should be updated")
}

func TestTaskRepository_CreateTasksWithDifferentTypes(t *testing.T) {
	// Arrange
	repo := setupTestRepository(t)
	projectID := "test_project"

	testTasks := []struct {
		name     string
		taskType domain.TaskType
	}{
		{"Regular Task", domain.RegularTask},
		{"Bug Task", domain.Bug},
		{"Feature Task", domain.Feature},
	}

	// Act & Assert - Create tasks with different types
	for i, taskData := range testTasks {
		task := &domain.Task{
			ID:        "task_" + string(rune(i)),
			ProjectID: projectID,
			Name:      taskData.name,
			Status:    domain.NotStarted,
			Type:      taskData.taskType,
			Priority:  domain.Low,
		}
		err := repo.Create(task)
		require.NoError(t, err, "Should be able to create task with type %v", taskData.taskType)
	}
}

func TestTaskRepository_RetrieveTaskByType(t *testing.T) {
	// Arrange
	repo := setupTestRepository(t)
	projectID := "test_project"

	// Create test tasks with different types
	testTasks := []struct {
		name     string
		taskType domain.TaskType
	}{
		{"Regular Task", domain.RegularTask},
		{"Bug Task", domain.Bug},
		{"Feature Task", domain.Feature},
	}

	taskIDs := []string{"task_regular", "task_bug", "task_feature"}
	for i, taskData := range testTasks {
		task := &domain.Task{
			ID:        taskIDs[i],
			ProjectID: projectID,
			Name:      taskData.name,
			Status:    domain.NotStarted,
			Type:      taskData.taskType,
			Priority:  domain.Low,
		}
		err := repo.Create(task)
		require.NoError(t, err)
	}

	// Act & Assert - Test GetByID preserves types
	for i, taskData := range testTasks {
		retrievedTask, err := repo.GetByID(taskIDs[i])
		require.NoError(t, err, "Should be able to retrieve task %d", i)
		require.NotNil(t, retrievedTask, "Retrieved task should not be nil")
		assert.Equal(t, taskData.name, retrievedTask.Name, "Task name should match")
		assert.Equal(t, taskData.taskType, retrievedTask.Type, "Task type should match")
	}
}

func TestTaskRepository_GetByProjectID_TypePreservation(t *testing.T) {
	// Arrange
	repo := setupTestRepository(t)
	projectID := "test_project"

	// Create test tasks with different types
	testTasks := []TaskTestData{
		{"Regular Task", domain.NotStarted, domain.RegularTask, domain.Low, time.Now(), time.Now()},
		{"Bug Task", domain.NotStarted, domain.Bug, domain.Medium, time.Now(), time.Now()},
		{"Feature Task", domain.NotStarted, domain.Feature, domain.High, time.Now(), time.Now()},
	}
	createTasksForStatus(t, repo, projectID, domain.NotStarted, testTasks)

	// Act
	allTasks, err := repo.GetByProjectID(projectID)

	// Assert
	require.NoError(t, err, "Should be able to get all tasks for project")
	assert.Len(t, allTasks, 3, "Should have 3 tasks")

	// Verify types are preserved
	taskTypes := make(map[domain.TaskType]int)
	for _, task := range allTasks {
		taskTypes[task.Type]++
	}
	assert.Equal(t, 1, taskTypes[domain.RegularTask], "Should have 1 RegularTask")
	assert.Equal(t, 1, taskTypes[domain.Bug], "Should have 1 Bug")
	assert.Equal(t, 1, taskTypes[domain.Feature], "Should have 1 Feature")
}

func TestTaskRepository_UpdateTaskType(t *testing.T) {
	// Arrange
	repo := setupTestRepository(t)
	projectID := "test_project"

	// Create a bug task
	bugTask := &domain.Task{
		ID:        "task_bug",
		ProjectID: projectID,
		Name:      "Bug Task",
		Status:    domain.NotStarted,
		Type:      domain.Bug,
		Priority:  domain.High,
	}
	err := repo.Create(bugTask)
	require.NoError(t, err)

	// Act - Change type from Bug to Feature
	bugTask.Type = domain.Feature
	err = repo.Update(bugTask)
	require.NoError(t, err, "Should be able to update task type")

	// Assert
	updatedTask, err := repo.GetByID("task_bug")
	require.NoError(t, err)
	assert.Equal(t, domain.Feature, updatedTask.Type, "Task type should be updated to Feature")
}

func TestTaskRepository_Migration_DefaultType(t *testing.T) {
	// Setup in-memory database with migrations
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	defer db.Close()

	// Run migrations
	dbWrapper := &database.Database{Db: db}
	err = dbWrapper.RunMigrations()
	require.NoError(t, err)

	// Create repository
	repo := NewSQLiteTaskRepository(db)
	projectID := "test_project"

	// Create a task without explicitly setting type (should default to RegularTask)
	task := &domain.Task{
		ID:        "task_default",
		ProjectID: projectID,
		Name:      "Default Type Task",
		Status:    domain.NotStarted,
		Priority:  domain.Low,
	}
	err = repo.Create(task)
	require.NoError(t, err, "Should be able to create task without explicit type")

	// Verify default type
	retrievedTask, err := repo.GetByID("task_default")
	require.NoError(t, err)
	assert.Equal(t, domain.RegularTask, retrievedTask.Type, "Default task type should be RegularTask")
}
