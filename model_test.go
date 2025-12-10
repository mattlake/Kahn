package main

import (
	"testing"

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
	assert.False(t, model.showProjectForm, "showProjectForm should be false initially")
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
