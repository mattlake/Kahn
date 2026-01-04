package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewProject(t *testing.T) {
	tests := []struct {
		name        string
		projectName string
		description string
		color       string
	}{
		{
			name:        "Create basic project",
			projectName: "Test Project",
			description: "Test Description",
			color:       "blue",
		},
		{
			name:        "Create project with empty description",
			projectName: "Simple Project",
			description: "",
			color:       "red",
		},
		{
			name:        "Create project with special characters",
			projectName: "Project with émojis 🎉",
			description: "Description with special chars: @#$%",
			color:       "green",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			before := time.Now()
			project := NewProject(tt.projectName, tt.description, tt.color)
			after := time.Now()

			require.NotNil(t, project, "Project should not be nil")
			assert.NotEmpty(t, project.ID, "Project ID should be generated")
			assert.True(t, len(project.ID) > 4, "Project ID should have reasonable length")
			assert.Contains(t, project.ID, "proj_", "Project ID should have prefix")
			assert.Equal(t, tt.projectName, project.Name, "Project name should match")
			assert.Equal(t, tt.description, project.Description, "Project description should match")
			assert.Equal(t, tt.color, project.Color, "Project color should match")
			assert.Equal(t, []Task{}, project.Tasks, "Default tasks should be empty")

			// Test timestamps
			assert.True(t, project.CreatedAt.After(before) || project.CreatedAt.Equal(before), "CreatedAt should be set correctly")
			assert.True(t, project.CreatedAt.Before(after) || project.CreatedAt.Equal(after), "CreatedAt should be reasonable")
			assert.True(t, project.UpdatedAt.After(before) || project.UpdatedAt.Equal(before), "UpdatedAt should be set correctly")
			assert.True(t, project.UpdatedAt.Before(after) || project.UpdatedAt.Equal(after), "UpdatedAt should be reasonable")
			assert.True(t, project.CreatedAt.Equal(project.UpdatedAt), "CreatedAt and UpdatedAt should be equal for new project")
		})
	}
}

func TestGenerateProjectID(t *testing.T) {
	// Test multiple calls to ensure uniqueness
	id1 := generateProjectID()
	time.Sleep(1 * time.Millisecond) // Ensure different timestamps
	id2 := generateProjectID()

	assert.NotEqual(t, id1, id2, "Generated project IDs should be unique")
	assert.Contains(t, id1, "proj_", "Project ID should have prefix")
	assert.Contains(t, id2, "proj_", "Project ID should have prefix")
	assert.True(t, len(id1) > 10, "Project ID should have reasonable length")
	assert.True(t, len(id2) > 10, "Project ID should have reasonable length")
}

func TestProject_AddTask(t *testing.T) {
	project := createTestProject("Test Project", "Test Description", "blue")
	task := createTestTask("Test Task", "Test Description", "different_project", NotStarted)

	// Test adding task to project
	project.AddTask(*task)

	assert.Len(t, project.Tasks, 1, "Project should have 1 task")
	assert.Equal(t, task.ID, project.Tasks[0].ID, "Task ID should match")
	assert.Equal(t, task.Name, project.Tasks[0].Name, "Task name should match")
	assert.Equal(t, project.ID, project.Tasks[0].ProjectID, "Task ProjectID should be updated to project ID")
	assert.True(t, project.UpdatedAt.After(project.CreatedAt), "UpdatedAt should be updated after adding task")
}

func TestProject_AddMultipleTasks(t *testing.T) {
	project := createTestProject("Test Project", "Test Description", "blue")

	// Add multiple tasks
	task1 := createTestTask("Task 1", "Description 1", "proj_1", NotStarted)
	task2 := createTestTask("Task 2", "Description 2", "proj_2", InProgress)
	task3 := createTestTask("Task 3", "Description 3", "proj_3", Done)

	project.AddTask(*task1)
	project.AddTask(*task2)
	project.AddTask(*task3)

	assert.Len(t, project.Tasks, 3, "Project should have 3 tasks")
	assert.Equal(t, "Task 1", project.Tasks[0].Name, "First task name should match")
	assert.Equal(t, "Task 2", project.Tasks[1].Name, "Second task name should match")
	assert.Equal(t, "Task 3", project.Tasks[2].Name, "Third task name should match")

	// All tasks should have been same project ID
	for _, task := range project.Tasks {
		assert.Equal(t, project.ID, task.ProjectID, "All tasks should have project ID")
	}
}

func TestProject_RemoveTask(t *testing.T) {
	project := createTestProject("Test Project", "Test Description", "blue")

	// Add tasks
	task1 := createTestTask("Task 1", "Description 1", "proj_1", NotStarted)
	task2 := createTestTask("Task 2", "Description 2", "proj_2", InProgress)
	task3 := createTestTask("Task 3", "Description 3", "proj_3", Done)

	project.AddTask(*task1)
	project.AddTask(*task2)
	project.AddTask(*task3)

	// Test removing existing task
	removed := project.RemoveTask(task2.ID)
	assert.True(t, removed, "RemoveTask should return true for existing task")
	assert.Len(t, project.Tasks, 2, "Project should have 2 tasks after removal")
	assert.Equal(t, task1.ID, project.Tasks[0].ID, "First task should remain")
	assert.Equal(t, task3.ID, project.Tasks[1].ID, "Third task should remain")

	// Test removing non-existing task
	removed = project.RemoveTask("non_existing_task_id")
	assert.False(t, removed, "RemoveTask should return false for non-existing task")
	assert.Len(t, project.Tasks, 2, "Project should still have 2 tasks")
}

func TestProject_GetTasksByStatus(t *testing.T) {
	project := createTestProject("Test Project", "Test Description", "blue")

	// Add tasks with different statuses
	task1 := createTestTask("Task 1", "Description 1", "proj_1", NotStarted)
	task2 := createTestTask("Task 2", "Description 2", "proj_2", InProgress)
	task3 := createTestTask("Task 3", "Description 3", "proj_3", Done)
	task4 := createTestTask("Task 4", "Description 4", "proj_4", NotStarted)

	project.AddTask(*task1)
	project.AddTask(*task2)
	project.AddTask(*task3)
	project.AddTask(*task4)

	// Test filtering by status
	notStartedTasks := project.GetTasksByStatus(NotStarted)
	inProgressTasks := project.GetTasksByStatus(InProgress)
	doneTasks := project.GetTasksByStatus(Done)

	assert.Len(t, notStartedTasks, 2, "Should have 2 NotStarted tasks")
	assert.Len(t, inProgressTasks, 1, "Should have 1 InProgress task")
	assert.Len(t, doneTasks, 1, "Should have 1 Done task")

	// Verify task names
	assert.Equal(t, "Task 1", notStartedTasks[0].Name)
	assert.Equal(t, "Task 4", notStartedTasks[1].Name)
	assert.Equal(t, "Task 2", inProgressTasks[0].Name)
	assert.Equal(t, "Task 3", doneTasks[0].Name)
}

func TestProject_UpdateTaskStatus(t *testing.T) {
	project := createTestProject("Test Project", "Test Description", "blue")

	// Add a task
	task := createTestTask("Test Task", "Description", "proj_1", NotStarted)
	project.AddTask(*task)

	// Test updating status
	updated := project.UpdateTaskStatus(task.ID, InProgress)
	assert.True(t, updated, "UpdateTaskStatus should return true for existing task")
	assert.Equal(t, InProgress, project.Tasks[0].Status, "Task status should be updated")
	assert.True(t, project.Tasks[0].UpdatedAt.After(project.Tasks[0].CreatedAt), "Task UpdatedAt should be updated")

	// Test updating non-existing task
	updated = project.UpdateTaskStatus("non_existing_task_id", Done)
	assert.False(t, updated, "UpdateTaskStatus should return false for non-existing task")
}

func TestProject_Validate(t *testing.T) {
	tests := []struct {
		name        string
		project     *Project
		expectError bool
		errorMsg    string
	}{
		{
			name:        "Valid project",
			project:     createTestProject("Valid Project", "Valid Description", "blue"),
			expectError: false,
		},
		{
			name:        "Empty name",
			project:     createTestProject("", "Valid Description", "blue"),
			expectError: true,
			errorMsg:    "project name is required",
		},
		{
			name:        "Whitespace name",
			project:     createTestProject("   ", "Valid Description", "blue"),
			expectError: true,
			errorMsg:    "project name is required",
		},
		{
			name:        "Name too long",
			project:     createTestProject(string(make([]byte, 51)), "Valid Description", "blue"),
			expectError: true,
			errorMsg:    "project name too long",
		},
		{
			name:        "Description too long",
			project:     createTestProject("Valid Name", string(make([]byte, 201)), "blue"),
			expectError: true,
			errorMsg:    "project description too long",
		},
		{
			name:        "Maximum valid name length",
			project:     createTestProject(string(make([]byte, MaxProjectNameLength)), "Valid Description", "blue"),
			expectError: false,
		},
		{
			name:        "Maximum valid description length",
			project:     createTestProject("Valid Name", string(make([]byte, MaxProjectDescriptionLength)), "blue"),
			expectError: false,
		},
		{
			name:        "Name with spaces within length",
			project:     createTestProject("  My Project  ", "Valid Description", "blue"),
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.project.Validate()

			if tt.expectError {
				assert.Error(t, err, "Should return error")
				assert.Contains(t, err.Error(), tt.errorMsg, "Error message should match expected")
			} else {
				assert.NoError(t, err, "Should not return error")
			}
		})
	}
}

func TestProject_GetTasksByStatus_NotStarted_SortsByPriorityThenCreated(t *testing.T) {
	// Arrange
	project := createProjectWithTasksOfVaryingPriorities()

	// Act
	notStartedTasks := project.GetTasksByStatus(NotStarted)

	// Assert
	assert.Len(t, notStartedTasks, 5, "Should have 5 NotStarted tasks")

	expectedOrder := []string{
		"High Priority Old",   // High priority, oldest
		"High Priority New",   // High priority, newer
		"Medium Priority Old", // Medium priority, oldest
		"Low Priority Old",    // Low priority, oldest
		"Low Priority New",    // Low priority, newer
	}

	actualOrder := getTaskNamesFromTasks(notStartedTasks)
	assert.Equal(t, expectedOrder, actualOrder, "NotStarted tasks should be sorted by priority DESC, then created_at ASC")
}

func TestProject_GetTasksByStatus_InProgress_SortsByUpdatedDesc(t *testing.T) {
	// Arrange
	project := createProjectWithTasksOfVaryingPriorities()

	// Update task statuses with specific timing to test updated_at ordering
	// We need to find the actual task IDs after they're created
	var highPriorityOldID, lowPriorityNewID string
	for _, task := range project.Tasks {
		if task.Name == "High Priority Old" {
			highPriorityOldID = task.ID
		} else if task.Name == "Low Priority New" {
			lowPriorityNewID = task.ID
		}
	}

	updates := []TaskUpdate{
		{TaskID: highPriorityOldID, NewStatus: InProgress, Delay: 1 * time.Millisecond},
		{TaskID: lowPriorityNewID, NewStatus: InProgress, Delay: 1 * time.Millisecond},
	}
	updateTaskStatusesWithTimestamps(project, updates)

	// Act
	inProgressTasks := project.GetTasksByStatus(InProgress)

	// Assert
	assert.Len(t, inProgressTasks, 2, "Should have 2 InProgress tasks")

	// Most recently updated task should be first (lowPriorityNew was updated last)
	assert.Equal(t, "Low Priority New", inProgressTasks[0].Name, "Most recently updated should be first")
	assert.Equal(t, "High Priority Old", inProgressTasks[1].Name, "Less recently updated should be second")

	// Verify updated_at timestamps are in descending order
	assert.True(t, inProgressTasks[0].UpdatedAt.After(inProgressTasks[1].UpdatedAt),
		"InProgress tasks should be sorted by updated_at DESC")
}

func TestProject_GetTasksByStatus_Done_SortsByUpdatedDesc(t *testing.T) {
	// Arrange
	project := createProjectWithTasksOfVaryingPriorities()

	// Update one task to Done status
	var mediumPriorityOldID string
	for _, task := range project.Tasks {
		if task.Name == "Medium Priority Old" {
			mediumPriorityOldID = task.ID
			break
		}
	}
	project.UpdateTaskStatus(mediumPriorityOldID, Done)

	// Act
	doneTasks := project.GetTasksByStatus(Done)

	// Assert
	assert.Len(t, doneTasks, 1, "Should have 1 Done task")
	assert.Equal(t, "Medium Priority Old", doneTasks[0].Name)
}

func TestProject_GetTasksByStatus_HandlesMultiplePriorityLevels(t *testing.T) {
	// Arrange
	project := createProjectWithTasksOfVaryingPriorities()

	// Act
	notStartedTasks := project.GetTasksByStatus(NotStarted)

	// Assert
	// Verify all priority levels are present and correctly grouped
	priorities := make(map[Priority]int)
	for _, task := range notStartedTasks {
		priorities[task.Priority]++
	}

	assert.Equal(t, 2, priorities[High], "Should have 2 high priority tasks")
	assert.Equal(t, 1, priorities[Medium], "Should have 1 medium priority task")
	assert.Equal(t, 2, priorities[Low], "Should have 2 low priority tasks")
}

func TestSortTasks(t *testing.T) {
	// Test the SortTasks function directly
	baseTime := time.Now()

	tasks := []Task{
		{Name: "Low Priority", Priority: Low, CreatedAt: baseTime, UpdatedAt: baseTime},
		{Name: "High Priority", Priority: High, CreatedAt: baseTime.Add(1 * time.Hour), UpdatedAt: baseTime},
		{Name: "Medium Priority", Priority: Medium, CreatedAt: baseTime.Add(30 * time.Minute), UpdatedAt: baseTime.Add(2 * time.Hour)},
	}

	// Test NotStarted sorting
	sortedNotStarted := SortTasks(tasks, NotStarted)
	assert.Equal(t, "High Priority", sortedNotStarted[0].Name, "High priority should be first")
	assert.Equal(t, "Medium Priority", sortedNotStarted[1].Name, "Medium priority should be second")
	assert.Equal(t, "Low Priority", sortedNotStarted[2].Name, "Low priority should be third")

	// Test InProgress sorting
	sortedInProgress := SortTasks(tasks, InProgress)
	assert.Equal(t, "Medium Priority", sortedInProgress[0].Name, "Most recently updated should be first")
	assert.Equal(t, "Low Priority", sortedInProgress[1].Name, "Less recently updated should be second")
	assert.Equal(t, "High Priority", sortedInProgress[2].Name, "Least recently updated should be third")

	// Test Done sorting (same logic as InProgress)
	sortedDone := SortTasks(tasks, Done)
	assert.Equal(t, "Medium Priority", sortedDone[0].Name, "Most recently updated should be first")
}
