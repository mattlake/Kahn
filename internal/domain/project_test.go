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
			project:     createTestProject(string(make([]byte, 50)), "Valid Description", "blue"),
			expectError: false,
		},
		{
			name:        "Maximum valid description length",
			project:     createTestProject("Valid Name", string(make([]byte, 200)), "blue"),
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

func TestProject_CompleteWorkflow(t *testing.T) {
	// Test a complete project workflow
	project := NewProject("Workflow Project", "Testing complete workflow", "green")

	// Initial state
	assert.Empty(t, project.Tasks, "Initial project should have no tasks")

	// Add tasks
	task1 := NewTask("First Task", "First task description", project.ID)
	time.Sleep(1 * time.Millisecond) // Ensure different timestamps
	task2 := NewTask("Second Task", "Second task description", project.ID)

	project.AddTask(*task1)
	project.AddTask(*task2)

	assert.Len(t, project.Tasks, 2, "Project should have 2 tasks")

	// Update task statuses
	project.UpdateTaskStatus(task1.ID, InProgress)
	project.UpdateTaskStatus(task2.ID, Done)

	notStartedTasks := project.GetTasksByStatus(NotStarted)
	inProgressTasks := project.GetTasksByStatus(InProgress)
	doneTasks := project.GetTasksByStatus(Done)

	assert.Len(t, notStartedTasks, 0, "Should have no NotStarted tasks")
	assert.Len(t, inProgressTasks, 1, "Should have 1 InProgress task")
	assert.Len(t, doneTasks, 1, "Should have 1 Done task")

	// Remove a task
	removed := project.RemoveTask(task1.ID)
	assert.True(t, removed, "Task should be removed")
	assert.Len(t, project.Tasks, 1, "Project should have 1 task remaining")

	remainingTasks := project.GetTasksByStatus(Done)
	assert.Len(t, remainingTasks, 1, "Should have 1 remaining Done task")
}
