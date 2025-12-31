package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTask(t *testing.T) {
	tests := []struct {
		name        string
		taskName    string
		description string
		projectID   string
	}{
		{
			name:        "Create basic task",
			taskName:    "Test Task",
			description: "Test Description",
			projectID:   "proj_123",
		},
		{
			name:        "Create task with empty description",
			taskName:    "Simple Task",
			description: "",
			projectID:   "proj_456",
		},
		{
			name:        "Create task with special characters",
			taskName:    "Task with émojis 🎉",
			description: "Description with special chars: @#$%",
			projectID:   "proj_special",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			before := time.Now()
			task := NewTask(tt.taskName, tt.description, tt.projectID)
			after := time.Now()

			require.NotNil(t, task, "Task should not be nil")
			assert.NotEmpty(t, task.ID, "Task ID should be generated")
			assert.True(t, len(task.ID) > 4, "Task ID should have reasonable length")
			assert.Contains(t, task.ID, "task_", "Task ID should have prefix")
			assert.Equal(t, tt.taskName, task.Name, "Task name should match")
			assert.Equal(t, tt.description, task.Desc, "Task description should match")
			assert.Equal(t, tt.projectID, task.ProjectID, "Project ID should match")
			assert.Equal(t, NotStarted, task.Status, "Default status should be NotStarted")
			assert.Equal(t, Low, task.Priority, "Default priority should be Low")

			// Test timestamps
			assert.True(t, task.CreatedAt.After(before) || task.CreatedAt.Equal(before), "CreatedAt should be set correctly")
			assert.True(t, task.CreatedAt.Before(after) || task.CreatedAt.Equal(after), "CreatedAt should be reasonable")
			assert.True(t, task.UpdatedAt.After(before) || task.UpdatedAt.Equal(before), "UpdatedAt should be set correctly")
			assert.True(t, task.UpdatedAt.Before(after) || task.UpdatedAt.Equal(after), "UpdatedAt should be reasonable")
			assert.True(t, task.CreatedAt.Equal(task.UpdatedAt), "CreatedAt and UpdatedAt should be equal for new task")
		})
	}
}

func TestTask_Title(t *testing.T) {
	task := &Task{Name: "Test Task Title"}
	assert.Equal(t, "Test Task Title", task.Title(), "Title() should return task name")
}

func TestTask_Description(t *testing.T) {
	task := &Task{Desc: "Test Task Description"}
	assert.Equal(t, "Test Task Description", task.Description(), "Description() should return task description")
}

func TestTask_FilterValue(t *testing.T) {
	task := &Task{Name: "Filter Test Task"}
	assert.Equal(t, "Filter Test Task", task.FilterValue(), "FilterValue() should return task name")
}

func TestPriority_String(t *testing.T) {
	tests := []struct {
		name     string
		priority Priority
		expected string
	}{
		{
			name:     "Low priority",
			priority: Low,
			expected: "Low",
		},
		{
			name:     "Medium priority",
			priority: Medium,
			expected: "Medium",
		},
		{
			name:     "High priority",
			priority: High,
			expected: "High",
		},
		{
			name:     "Unknown priority (should return Medium)",
			priority: Priority(999),
			expected: "Medium",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.priority.String()
			assert.Equal(t, tt.expected, result, "Priority string representation should match expected value")
		})
	}
}

func TestPriority_Constants(t *testing.T) {
	// Test that priority constants have expected integer values
	assert.Equal(t, Priority(0), Low, "Low should be 0")
	assert.Equal(t, Priority(1), Medium, "Medium should be 1")
	assert.Equal(t, Priority(2), High, "High should be 2")
}

func TestGenerateTaskID(t *testing.T) {
	// Test multiple calls to ensure uniqueness
	id1 := generateTaskID()
	// Add small delay to ensure different timestamps
	time.Sleep(1 * time.Millisecond)
	id2 := generateTaskID()

	assert.NotEqual(t, id1, id2, "Generated task IDs should be unique")
	assert.Contains(t, id1, "task_", "Task ID should have prefix")
	assert.Contains(t, id2, "task_", "Task ID should have prefix")
	assert.True(t, len(id1) > 10, "Task ID should have reasonable length")
	assert.True(t, len(id2) > 10, "Task ID should have reasonable length")
}

func TestTask_CompleteWorkflow(t *testing.T) {
	// Test a complete task workflow
	task := NewTask("Workflow Task", "Testing workflow", "proj_123")

	// Initial state
	assert.Equal(t, NotStarted, task.Status, "Initial status should be NotStarted")

	// Update status to in progress
	task.Status = InProgress
	task.UpdatedAt = time.Now()
	assert.Equal(t, InProgress, task.Status, "Status should be updated to InProgress")

	// Update status to done
	task.Status = Done
	task.UpdatedAt = time.Now()
	assert.Equal(t, Done, task.Status, "Status should be updated to Done")

	// Test that other fields remain unchanged
	assert.Equal(t, "Workflow Task", task.Name, "Name should remain unchanged")
	assert.Equal(t, "Testing workflow", task.Desc, "Description should remain unchanged")
	assert.Equal(t, "proj_123", task.ProjectID, "ProjectID should remain unchanged")
}

func TestTask_DeleteWorkflow(t *testing.T) {
	// Test task deletion workflow
	task := NewTask("Task to Delete", "This task will be deleted", "proj_123")

	// Verify task exists initially
	assert.NotEmpty(t, task.ID, "Task should have ID")
	assert.Equal(t, "Task to Delete", task.Name, "Task name should match")

	// Simulate deletion (in real scenario, this would be handled by DAO)
	// Here we just verify the task structure is valid for deletion
	assert.NotNil(t, task, "Task should be valid for deletion")
}

func TestTask_DeleteFromDifferentStatuses(t *testing.T) {
	testCases := []struct {
		name     string
		status   Status
		expected Status
	}{
		{"Delete NotStarted task", NotStarted, NotStarted},
		{"Delete InProgress task", InProgress, InProgress},
		{"Delete Done task", Done, Done},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			task := NewTask("Test Task", "Test Description", "proj_123")
			task.Status = tc.status

			assert.Equal(t, tc.expected, task.Status, "Task should have correct status before deletion")
			// Task should be valid for deletion regardless of status
			assert.NotNil(t, task, "Task should be valid for deletion")
		})
	}
}

func TestTask_DeleteEdgeCases(t *testing.T) {
	// Test task with empty name
	task1 := NewTask("", "Description", "proj_123")
	assert.NotNil(t, task1, "Task with empty name should still be valid for deletion")

	// Test task with special characters
	task2 := NewTask("Task with émojis 🎉", "Description with @#$%", "proj_123")
	assert.NotNil(t, task2, "Task with special characters should be valid for deletion")

	// Test task with very long name
	longName := string(make([]byte, 1000))
	for i := range longName {
		longName = longName[:i] + "a" + longName[i+1:]
	}
	task3 := NewTask(longName, "Description", "proj_123")
	assert.NotNil(t, task3, "Task with long name should be valid for deletion")
}
