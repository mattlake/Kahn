package main

import (
	"testing"

	"github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewModel(t *testing.T) {
	// Create test database
	config := createTestConfig()
	database, err := NewDatabase(config)
	require.NoError(t, err, "NewDatabase should not return error")
	defer database.Close()

	// Create model
	model := NewModel(database)
	require.NotNil(t, model, "NewModel should not return nil")

	// Test initial state
	assert.NotEmpty(t, model.Projects, "Should have at least one project (default)")
	assert.NotEmpty(t, model.ActiveProjectID, "Should have active project ID")
	assert.Equal(t, model.Projects[0].ID, model.ActiveProjectID, "Active project should be first project")
	assert.Len(t, model.Tasks, 3, "Should have 3 task lists (NotStarted, InProgress, Done)")
	assert.Equal(t, NotStarted, model.activeListIndex, "Active list should be NotStarted initially")
	assert.False(t, model.showForm, "showForm should be false initially")
	assert.False(t, model.showProjectSwitch, "showProjectSwitch should be false initially")
	// showProjectForm field removed, so this test is no longer relevant
}

func TestModel_GetActiveProject(t *testing.T) {
	// Create test database
	config := createTestConfig()
	database, err := NewDatabase(config)
	require.NoError(t, err, "NewDatabase should not return error")
	defer database.Close()

	model := NewModel(database)
	require.NotNil(t, model, "NewModel should not return nil")

	// Test getting active project
	activeProject := model.GetActiveProject()
	assert.NotNil(t, activeProject, "GetActiveProject should not return nil")
	assert.Equal(t, model.ActiveProjectID, activeProject.ID, "Active project ID should match")
	assert.Equal(t, model.Projects[0].Name, activeProject.Name, "Active project name should match")
}

func TestModel_GetActiveProject_NoProjects(t *testing.T) {
	model := &Model{
		Projects:        []Project{},
		ActiveProjectID: "",
	}

	activeProject := model.GetActiveProject()
	assert.Nil(t, activeProject, "GetActiveProject should return nil when no projects")
}

func TestModel_GetActiveProject_NonExistentID(t *testing.T) {
	project1 := createTestProject("Project 1", "Description 1", "blue")
	project2 := createTestProject("Project 2", "Description 2", "red")

	model := &Model{
		Projects:        []Project{*project1, *project2},
		ActiveProjectID: "non_existent_id",
	}

	activeProject := model.GetActiveProject()
	assert.Nil(t, activeProject, "GetActiveProject should return nil when active ID not found")
}

func TestConvertTasksToListItems(t *testing.T) {
	tasks := []Task{
		*createTestTask("Task 1", "Description 1", "proj_1", NotStarted),
		*createTestTask("Task 2", "Description 2", "proj_2", InProgress),
		*createTestTask("Task 3", "Description 3", "proj_3", Done),
	}

	items := convertTasksToListItems(tasks)

	assert.Len(t, items, 3, "Should convert all tasks to list items")

	for i, task := range tasks {
		assert.Equal(t, task.FilterValue(), items[i].FilterValue(), "Task filter value should match item filter value")
	}
}

func TestConvertTasksToListItems_Empty(t *testing.T) {
	tasks := []Task{}

	items := convertTasksToListItems(tasks)

	assert.Len(t, items, 0, "Should return empty list for empty tasks")
}

func TestConvertTasksToListItems_Nil(t *testing.T) {
	items := convertTasksToListItems(nil)

	assert.Len(t, items, 0, "Should return empty list for nil tasks")
}

func TestModel_TaskDeletionKeyHandling(t *testing.T) {
	model := createTestModelWithTasks(t, []string{"Task 1", "Task 2"}, []Status{NotStarted, InProgress})

	// Test d key when task is selected
	model.Tasks[NotStarted].Select(0) // Select first task
	model = simulateKeyPress(t, model, "d")

	// Get the actual task ID that was set for deletion
	activeProj := model.GetActiveProject()
	expectedTaskID := activeProj.Tasks[0].ID
	assertTaskDeletionState(t, model, true, expectedTaskID)
}

func TestModel_TaskDeletionKeyHandling_NoTaskSelected(t *testing.T) {
	model := createTestModelWithTasks(t, []string{}, []Status{}) // Create model with no tasks

	// Test d key when no tasks exist
	model = simulateKeyPress(t, model, "d")

	// Should not be in deletion state since no tasks exist
	assertTaskDeletionState(t, model, false, "")
}

func TestModel_TaskDeletionConfirmation_Yes(t *testing.T) {
	model := createTestModelWithTasks(t, []string{"Task 1"}, []Status{NotStarted})

	// Setup deletion state
	taskID := model.GetActiveProject().Tasks[0].ID
	model.showTaskDeleteConfirm = true
	model.taskToDelete = taskID

	// Confirm deletion by calling executeTaskDeletion directly
	model = model.executeTaskDeletion()

	// Verify task is deleted
	assertTaskDeletionState(t, model, false, "")
	assertTaskNotInLists(t, model, taskID)
}

func TestModel_TaskDeletionConfirmation_No(t *testing.T) {
	model := createTestModelWithTasks(t, []string{"Task 1"}, []Status{NotStarted})

	// Setup deletion state
	taskID := model.GetActiveProject().Tasks[0].ID
	model.showTaskDeleteConfirm = true
	model.taskToDelete = taskID

	// Cancel deletion by simulating 'n' key
	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}}
	newModel, _ := model.Update(keyMsg)
	resultModel := newModel.(Model)
	model = &resultModel

	// Verify task is not deleted
	assertTaskDeletionState(t, model, false, "")
	assertTaskInList(t, model, taskID, NotStarted)
}

func TestModel_TaskDeletionConfirmation_Escape(t *testing.T) {
	model := createTestModelWithTasks(t, []string{"Task 1"}, []Status{NotStarted})

	// Setup deletion state
	taskID := model.GetActiveProject().Tasks[0].ID
	model.showTaskDeleteConfirm = true
	model.taskToDelete = taskID

	// Cancel deletion with escape
	escapeKey := tea.KeyMsg{Type: tea.KeyEsc}
	newModel, _ := model.Update(escapeKey)
	resultModel := newModel.(Model)
	model = &resultModel

	// Verify task is not deleted
	assertTaskDeletionState(t, model, false, "")
	assertTaskInList(t, model, taskID, NotStarted)
}

func TestModel_TaskDeletionExecution(t *testing.T) {
	model := createTestModelWithTasks(t, []string{"Task 1", "Task 2"}, []Status{NotStarted, InProgress})

	// Get task to delete
	taskToDelete := model.GetActiveProject().Tasks[0]
	taskID := taskToDelete.ID

	// Execute deletion
	model.showTaskDeleteConfirm = true
	model.taskToDelete = taskID
	model = model.executeTaskDeletion()

	// Verify task is deleted from database and memory
	assertTaskDeletionState(t, model, false, "")
	assertTaskNotInLists(t, model, taskID)

	// Verify other task is still there
	otherTask := model.GetActiveProject().Tasks[0]
	assert.NotEqual(t, taskID, otherTask.ID, "Other task should still exist")
}

func TestModel_TaskDeletionExecution_DatabaseError(t *testing.T) {
	model := createTestModelWithTasks(t, []string{"Task 1"}, []Status{NotStarted})

	// Get task to delete
	taskID := model.GetActiveProject().Tasks[0].ID

	// Setup deletion state with invalid task ID (will cause database error)
	model.showTaskDeleteConfirm = true
	model.taskToDelete = "invalid_task_id"
	model = model.executeTaskDeletion()

	// Verify state is reset on error
	assertTaskDeletionState(t, model, false, "")

	// Verify original task is still there
	assertTaskInList(t, model, taskID, NotStarted)
}

func TestModel_TaskDeletionEdgeCases(t *testing.T) {
	// Test deletion when taskToDelete is empty
	model := createTestModelWithTasks(t, []string{"Task 1"}, []Status{NotStarted})
	model.showTaskDeleteConfirm = true
	model.taskToDelete = ""
	model = model.executeTaskDeletion()

	assertTaskDeletionState(t, model, false, "")
}

func TestModel_TaskDeletionFromDifferentStatuses(t *testing.T) {
	testCases := []struct {
		name     string
		status   Status
		taskName string
	}{
		{"Delete from NotStarted", NotStarted, "NotStarted Task"},
		{"Delete from InProgress", InProgress, "InProgress Task"},
		{"Delete from Done", Done, "Done Task"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			model := createTestModelWithTasks(t, []string{tc.taskName}, []Status{tc.status})

			// Get task to delete
			taskID := model.GetActiveProject().Tasks[0].ID

			// Execute deletion
			model.showTaskDeleteConfirm = true
			model.taskToDelete = taskID
			model = model.executeTaskDeletion()

			// Verify task is deleted
			assertTaskNotInLists(t, model, taskID)
		})
	}
}
