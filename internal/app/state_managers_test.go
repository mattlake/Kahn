package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"kahn/internal/domain"
)

// UIStateManager Tests

func TestUIStateManager_GetCurrentViewState_BoardView(t *testing.T) {
	km, cleanup := setupTestApp(t)
	defer cleanup()

	viewState := km.uiStateManager.GetCurrentViewState()
	assert.Equal(t, BoardView, viewState)
}

func TestUIStateManager_GetCurrentViewState_FormView(t *testing.T) {
	km, cleanup := setupTestApp(t)
	defer cleanup()

	// Show task form
	km.uiStateManager.ShowTaskForm([]domain.Task{})

	viewState := km.uiStateManager.GetCurrentViewState()
	assert.Equal(t, FormView, viewState)
}

func TestUIStateManager_GetCurrentViewState_ProjectSwitchView(t *testing.T) {
	km, cleanup := setupTestApp(t)
	defer cleanup()

	// Show project switcher
	km.uiStateManager.ShowProjectSwitcher()

	viewState := km.uiStateManager.GetCurrentViewState()
	assert.Equal(t, ProjectSwitchView, viewState)
}

func TestUIStateManager_GetCurrentViewState_TaskDeleteConfirmView(t *testing.T) {
	km, cleanup := setupTestApp(t)
	defer cleanup()

	// Show task delete confirmation
	km.uiStateManager.ShowTaskDeleteConfirm("test-task-id")

	viewState := km.uiStateManager.GetCurrentViewState()
	assert.Equal(t, TaskDeleteConfirmView, viewState)
}

func TestUIStateManager_GetCurrentViewState_ProjectDeleteConfirmView(t *testing.T) {
	km, cleanup := setupTestApp(t)
	defer cleanup()

	// Show project delete confirmation
	km.uiStateManager.ShowProjectDeleteConfirm("test-project-id")

	viewState := km.uiStateManager.GetCurrentViewState()
	assert.Equal(t, ProjectDeleteConfirmView, viewState)
}

func TestUIStateManager_IsShowingAnyForm_False(t *testing.T) {
	km, cleanup := setupTestApp(t)
	defer cleanup()

	isShowing := km.uiStateManager.IsShowingAnyForm()
	assert.False(t, isShowing)
}

func TestUIStateManager_IsShowingAnyForm_True_TaskForm(t *testing.T) {
	km, cleanup := setupTestApp(t)
	defer cleanup()

	km.uiStateManager.ShowTaskForm([]domain.Task{})

	isShowing := km.uiStateManager.IsShowingAnyForm()
	assert.True(t, isShowing)
}

func TestUIStateManager_IsShowingAnyForm_True_ProjectSwitch(t *testing.T) {
	km, cleanup := setupTestApp(t)
	defer cleanup()

	km.uiStateManager.ShowProjectSwitcher()

	isShowing := km.uiStateManager.IsShowingAnyForm()
	assert.True(t, isShowing)
}

func TestUIStateManager_IsShowingAnyForm_True_TaskDeleteConfirm(t *testing.T) {
	km, cleanup := setupTestApp(t)
	defer cleanup()

	km.uiStateManager.ShowTaskDeleteConfirm("test-task-id")

	isShowing := km.uiStateManager.IsShowingAnyForm()
	assert.True(t, isShowing)
}

func TestUIStateManager_HideAllStates(t *testing.T) {
	km, cleanup := setupTestApp(t)
	defer cleanup()

	// Show multiple states
	km.uiStateManager.ShowTaskForm([]domain.Task{})
	km.uiStateManager.ShowProjectSwitcher()
	km.uiStateManager.ShowTaskDeleteConfirm("test-task-id")

	// Hide all
	km.uiStateManager.HideAllStates()

	// Verify all hidden
	assert.False(t, km.uiStateManager.IsShowingAnyForm())
	assert.Equal(t, BoardView, km.uiStateManager.GetCurrentViewState())
}

func TestUIStateManager_ShowTaskForm_HidesOtherStates(t *testing.T) {
	km, cleanup := setupTestApp(t)
	defer cleanup()

	// Show project switcher first
	km.uiStateManager.ShowProjectSwitcher()
	assert.Equal(t, ProjectSwitchView, km.uiStateManager.GetCurrentViewState())

	// Show task form
	km.uiStateManager.ShowTaskForm([]domain.Task{})

	// Verify only task form is showing
	assert.Equal(t, FormView, km.uiStateManager.GetCurrentViewState())
	assert.False(t, km.uiStateManager.NavigationState().IsShowingProjectSwitch())
}

func TestUIStateManager_ShowTaskEditForm(t *testing.T) {
	km, cleanup := setupTestApp(t)
	defer cleanup()

	// Create a task first
	taskID := createTestTask(t, km, "Test Task", "Description")

	// Show edit form
	task, err := km.taskService.GetTask(taskID)
	require.NoError(t, err)

	km.uiStateManager.ShowTaskEditForm(
		task.ID,
		task.Name,
		task.Desc,
		task.Priority,
		task.Type,
		task.BlockedBy,
		[]domain.Task{},
	)

	// Verify form is showing
	assert.Equal(t, FormView, km.uiStateManager.GetCurrentViewState())
	assert.Equal(t, task.ID, km.uiStateManager.FormState().GetTaskID())
}

func TestUIStateManager_ShowProjectForm(t *testing.T) {
	km, cleanup := setupTestApp(t)
	defer cleanup()

	// Show project form
	km.uiStateManager.ShowProjectForm()

	// Verify form is showing
	assert.Equal(t, FormView, km.uiStateManager.GetCurrentViewState())
	assert.True(t, km.uiStateManager.FormState().IsShowingForm())
}

func TestUIStateManager_ShowProjectSwitcher_HidesOtherStates(t *testing.T) {
	km, cleanup := setupTestApp(t)
	defer cleanup()

	// Show task form first
	km.uiStateManager.ShowTaskForm([]domain.Task{})
	assert.Equal(t, FormView, km.uiStateManager.GetCurrentViewState())

	// Show project switcher
	km.uiStateManager.ShowProjectSwitcher()

	// Verify only project switcher is showing
	assert.Equal(t, ProjectSwitchView, km.uiStateManager.GetCurrentViewState())
	assert.False(t, km.uiStateManager.FormState().IsShowingForm())
}

func TestUIStateManager_ShowTaskDeleteConfirm_HidesOtherStates(t *testing.T) {
	km, cleanup := setupTestApp(t)
	defer cleanup()

	// Show form first
	km.uiStateManager.ShowTaskForm([]domain.Task{})

	// Show task delete confirm
	taskID := "test-task-id"
	km.uiStateManager.ShowTaskDeleteConfirm(taskID)

	// Verify only delete confirm is showing
	assert.Equal(t, TaskDeleteConfirmView, km.uiStateManager.GetCurrentViewState())
	assert.False(t, km.uiStateManager.FormState().IsShowingForm())
	assert.Equal(t, taskID, km.uiStateManager.ConfirmationState().GetTaskToDelete())
}

func TestUIStateManager_ShowProjectDeleteConfirm_HidesOtherStates(t *testing.T) {
	km, cleanup := setupTestApp(t)
	defer cleanup()

	// Show form first
	km.uiStateManager.ShowTaskForm([]domain.Task{})

	// Show project delete confirm
	projectID := "test-project-id"
	km.uiStateManager.ShowProjectDeleteConfirm(projectID)

	// Verify only delete confirm is showing
	assert.Equal(t, ProjectDeleteConfirmView, km.uiStateManager.GetCurrentViewState())
	assert.False(t, km.uiStateManager.FormState().IsShowingForm())
	assert.Equal(t, projectID, km.uiStateManager.ConfirmationState().GetProjectToDelete())
}

// ProjectManager Tests

func TestProjectManager_InitializeProjects_CreatesDefault(t *testing.T) {
	km, cleanup := setupTestApp(t)
	defer cleanup()

	// Verify default project was created
	assert.True(t, km.projectManager.HasProjects())
	assert.Equal(t, 1, km.projectManager.GetProjectCount())

	activeProj := km.projectManager.GetActiveProject()
	require.NotNil(t, activeProj)
	assert.Equal(t, "Default Project", activeProj.Name)
}

func TestProjectManager_GetActiveProject(t *testing.T) {
	km, cleanup := setupTestApp(t)
	defer cleanup()

	activeProj := km.projectManager.GetActiveProject()
	require.NotNil(t, activeProj)
	assert.NotEmpty(t, activeProj.ID)
	assert.NotEmpty(t, activeProj.Name)
}

func TestProjectManager_GetActiveProjectID(t *testing.T) {
	km, cleanup := setupTestApp(t)
	defer cleanup()

	projectID := km.projectManager.GetActiveProjectID()
	assert.NotEmpty(t, projectID)

	activeProj := km.projectManager.GetActiveProject()
	assert.Equal(t, activeProj.ID, projectID)
}

func TestProjectManager_CreateProject(t *testing.T) {
	km, cleanup := setupTestApp(t)
	defer cleanup()

	initialCount := km.projectManager.GetProjectCount()

	// Create new project
	err := km.projectManager.CreateProject("New Project", "New description")
	require.NoError(t, err)

	// Verify project count increased
	assert.Equal(t, initialCount+1, km.projectManager.GetProjectCount())

	// Verify new project is active
	activeProj := km.projectManager.GetActiveProject()
	assert.Equal(t, "New Project", activeProj.Name)
	assert.Equal(t, "New description", activeProj.Description)
}

func TestProjectManager_SwitchToProject(t *testing.T) {
	km, cleanup := setupTestApp(t)
	defer cleanup()

	// Create second project
	err := km.projectManager.CreateProject("Project 2", "Description")
	require.NoError(t, err)
	project2ID := km.projectManager.GetActiveProjectID()

	// Create third project (becomes active)
	err = km.projectManager.CreateProject("Project 3", "Description")
	require.NoError(t, err)

	// Switch back to project 2
	err = km.projectManager.SwitchToProject(project2ID)
	require.NoError(t, err)

	// Verify active project
	assert.Equal(t, project2ID, km.projectManager.GetActiveProjectID())
	activeProj := km.projectManager.GetActiveProject()
	assert.Equal(t, "Project 2", activeProj.Name)
}

func TestProjectManager_DeleteProject_LastProject(t *testing.T) {
	km, cleanup := setupTestApp(t)
	defer cleanup()

	// Get the default project ID
	projectID := km.projectManager.GetActiveProjectID()

	// Delete the only project
	err := km.projectManager.DeleteProject(projectID)
	require.NoError(t, err)

	// Verify no projects remain
	assert.False(t, km.projectManager.HasProjects())
	assert.Equal(t, 0, km.projectManager.GetProjectCount())
	assert.Nil(t, km.projectManager.GetActiveProject())
	assert.Empty(t, km.projectManager.GetActiveProjectID())
}

func TestProjectManager_DeleteProject_NonActiveProject(t *testing.T) {
	km, cleanup := setupTestApp(t)
	defer cleanup()

	// Create second project
	err := km.projectManager.CreateProject("Project 2", "Description")
	require.NoError(t, err)
	project2ID := km.projectManager.GetActiveProjectID()

	// Create third project (becomes active)
	err = km.projectManager.CreateProject("Project 3", "Description")
	require.NoError(t, err)
	project3ID := km.projectManager.GetActiveProjectID()

	// Delete project 2 (non-active)
	err = km.projectManager.DeleteProject(project2ID)
	require.NoError(t, err)

	// Verify project 3 is still active
	assert.Equal(t, 2, km.projectManager.GetProjectCount())
	assert.Equal(t, project3ID, km.projectManager.GetActiveProjectID())
}

func TestProjectManager_DeleteProject_ActiveProject(t *testing.T) {
	km, cleanup := setupTestApp(t)
	defer cleanup()

	// Get default project
	firstProjectID := km.projectManager.GetActiveProjectID()

	// Create second project (becomes active)
	err := km.projectManager.CreateProject("Project 2", "Description")
	require.NoError(t, err)
	secondProjectID := km.projectManager.GetActiveProjectID()

	// Delete active project (second)
	err = km.projectManager.DeleteProject(secondProjectID)
	require.NoError(t, err)

	// Verify first project became active again
	assert.Equal(t, 1, km.projectManager.GetProjectCount())
	assert.Equal(t, firstProjectID, km.projectManager.GetActiveProjectID())
}

func TestProjectManager_HasProjects_True(t *testing.T) {
	km, cleanup := setupTestApp(t)
	defer cleanup()

	assert.True(t, km.projectManager.HasProjects())
}

func TestProjectManager_HasProjects_False(t *testing.T) {
	km, cleanup := setupTestApp(t)
	defer cleanup()

	// Delete the default project
	projectID := km.projectManager.GetActiveProjectID()
	err := km.projectManager.DeleteProject(projectID)
	require.NoError(t, err)

	assert.False(t, km.projectManager.HasProjects())
}

func TestProjectManager_GetProjectCount(t *testing.T) {
	tests := []struct {
		name            string
		additionalProjs int
		expectedCount   int
	}{
		{"default only", 0, 1},
		{"two projects", 1, 2},
		{"five projects", 4, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			km, cleanup := setupTestApp(t)
			defer cleanup()

			// Create additional projects
			for i := 0; i < tt.additionalProjs; i++ {
				err := km.projectManager.CreateProject("Project "+string(rune('A'+i)), "Description")
				require.NoError(t, err)
			}

			assert.Equal(t, tt.expectedCount, km.projectManager.GetProjectCount())
		})
	}
}

func TestProjectManager_GetProjectsAsDomain(t *testing.T) {
	km, cleanup := setupTestApp(t)
	defer cleanup()

	// Create additional projects
	err := km.projectManager.CreateProject("Project A", "Description A")
	require.NoError(t, err)
	err = km.projectManager.CreateProject("Project B", "Description B")
	require.NoError(t, err)

	projects := km.projectManager.GetProjectsAsDomain()
	assert.Len(t, projects, 3) // Default + 2 new

	// Verify all projects have IDs and names
	for _, proj := range projects {
		assert.NotEmpty(t, proj.ID)
		assert.NotEmpty(t, proj.Name)
	}
}

// NavigationState Tests

func TestNavigationState_GetActiveListIndex_DefaultIsNotStarted(t *testing.T) {
	km, cleanup := setupTestApp(t)
	defer cleanup()

	index := km.navState.GetActiveListIndex()
	assert.Equal(t, domain.NotStarted, index)
}

func TestNavigationState_SetActiveListIndex(t *testing.T) {
	km, cleanup := setupTestApp(t)
	defer cleanup()

	// Set to InProgress
	km.navState.SetActiveListIndex(domain.InProgress)
	assert.Equal(t, domain.InProgress, km.navState.GetActiveListIndex())

	// Set to Done
	km.navState.SetActiveListIndex(domain.Done)
	assert.Equal(t, domain.Done, km.navState.GetActiveListIndex())
}

func TestNavigationState_NextList(t *testing.T) {
	tests := []struct {
		name     string
		current  domain.Status
		expected domain.Status
	}{
		{"NotStarted to InProgress", domain.NotStarted, domain.InProgress},
		{"InProgress to Done", domain.InProgress, domain.Done},
		{"Done wraps to NotStarted", domain.Done, domain.NotStarted},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			km, cleanup := setupTestApp(t)
			defer cleanup()

			// Set starting position
			km.navState.SetActiveListIndex(tt.current)

			// Move to next
			km.navState.NextList()

			// Verify new position
			assert.Equal(t, tt.expected, km.navState.GetActiveListIndex())
		})
	}
}

func TestNavigationState_PrevList(t *testing.T) {
	tests := []struct {
		name     string
		current  domain.Status
		expected domain.Status
	}{
		{"InProgress to NotStarted", domain.InProgress, domain.NotStarted},
		{"Done to InProgress", domain.Done, domain.InProgress},
		{"NotStarted wraps to Done", domain.NotStarted, domain.Done},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			km, cleanup := setupTestApp(t)
			defer cleanup()

			// Set starting position
			km.navState.SetActiveListIndex(tt.current)

			// Move to previous
			km.navState.PrevList()

			// Verify new position
			assert.Equal(t, tt.expected, km.navState.GetActiveListIndex())
		})
	}
}

func TestNavigationState_ShowProjectSwitch(t *testing.T) {
	km, cleanup := setupTestApp(t)
	defer cleanup()

	assert.False(t, km.navState.IsShowingProjectSwitch())

	km.navState.ShowProjectSwitch()

	assert.True(t, km.navState.IsShowingProjectSwitch())
}

func TestNavigationState_HideProjectSwitch(t *testing.T) {
	km, cleanup := setupTestApp(t)
	defer cleanup()

	km.navState.ShowProjectSwitch()
	assert.True(t, km.navState.IsShowingProjectSwitch())

	km.navState.HideProjectSwitch()

	assert.False(t, km.navState.IsShowingProjectSwitch())
}

func TestNavigationState_UpdateTaskLists(t *testing.T) {
	km, cleanup := setupTestApp(t)
	defer cleanup()

	// Create some tasks
	createTestTask(t, km, "Task 1", "Description 1")
	createTestTask(t, km, "Task 2", "Description 2")

	// Manually move one task to InProgress
	activeProj := km.projectManager.GetActiveProject()
	task2ID := activeProj.Tasks[1].ID
	moveTaskToStatus(t, km, task2ID, domain.InProgress)

	// Verify task lists updated
	notStartedItems := km.navState.GetTaskItems(domain.NotStarted)
	inProgressItems := km.navState.GetTaskItems(domain.InProgress)

	assert.Len(t, notStartedItems, 1)
	assert.Len(t, inProgressItems, 1)
}

func TestNavigationState_UpdateTaskLists_NilProject(t *testing.T) {
	km, cleanup := setupTestApp(t)
	defer cleanup()

	// Should handle nil project gracefully
	km.navState.UpdateTaskLists(nil, km.taskService)

	// Should not panic
}

func TestNavigationState_UpdateTaskListsWithSearch(t *testing.T) {
	km, cleanup := setupTestApp(t)
	defer cleanup()

	// Create tasks with different names
	createTestTask(t, km, "API Task", "Description")
	createTestTask(t, km, "Database Task", "Description")
	createTestTask(t, km, "API Service", "Description")

	activeProj := km.projectManager.GetActiveProject()

	// Apply search filter for "API"
	km.navState.UpdateTaskListsWithSearch(activeProj, km.taskService, "API")

	// Verify only API tasks are shown
	notStartedItems := km.navState.GetTaskItems(domain.NotStarted)
	assert.Len(t, notStartedItems, 2) // Should have 2 tasks with "API"
}

func TestNavigationState_UpdateTaskListsWithSearch_EmptyQuery(t *testing.T) {
	km, cleanup := setupTestApp(t)
	defer cleanup()

	// Create tasks
	createTestTask(t, km, "Task 1", "Description")
	createTestTask(t, km, "Task 2", "Description")

	activeProj := km.projectManager.GetActiveProject()

	// Apply empty search (should show all)
	km.navState.UpdateTaskListsWithSearch(activeProj, km.taskService, "")

	// Verify all tasks shown
	notStartedItems := km.navState.GetTaskItems(domain.NotStarted)
	assert.Len(t, notStartedItems, 2)
}

func TestNavigationState_MarkListDirty(t *testing.T) {
	km, cleanup := setupTestApp(t)
	defer cleanup()

	// Mark NotStarted as dirty
	km.navState.MarkListDirty(domain.NotStarted)

	assert.True(t, km.navState.IsListDirty(domain.NotStarted))
	assert.False(t, km.navState.IsListDirty(domain.InProgress))
	assert.False(t, km.navState.IsListDirty(domain.Done))
}

func TestNavigationState_MarkAllListsDirty(t *testing.T) {
	km, cleanup := setupTestApp(t)
	defer cleanup()

	km.navState.MarkAllListsDirty()

	assert.True(t, km.navState.IsListDirty(domain.NotStarted))
	assert.True(t, km.navState.IsListDirty(domain.InProgress))
	assert.True(t, km.navState.IsListDirty(domain.Done))
}

func TestNavigationState_IsListDirty_NoDirtyFlags(t *testing.T) {
	km, cleanup := setupTestApp(t)
	defer cleanup()

	// No dirty flags set
	assert.False(t, km.navState.IsListDirty(domain.NotStarted))
}

func TestNavigationState_UpdateDirtyLists(t *testing.T) {
	km, cleanup := setupTestApp(t)
	defer cleanup()

	// Create tasks
	createTestTask(t, km, "Task 1", "Description")
	createTestTask(t, km, "Task 2", "Description")

	activeProj := km.projectManager.GetActiveProject()

	// Mark only NotStarted as dirty
	km.navState.MarkListDirty(domain.NotStarted)

	// Update dirty lists
	km.navState.UpdateDirtyLists(activeProj, km.taskService)

	// Verify dirty flags cleared
	assert.False(t, km.navState.IsListDirty(domain.NotStarted))
}

func TestNavigationState_GetActiveList(t *testing.T) {
	km, cleanup := setupTestApp(t)
	defer cleanup()

	// Default should be NotStarted
	activeList := km.navState.GetActiveList()
	assert.NotNil(t, activeList)

	// Move to InProgress
	km.navState.SetActiveListIndex(domain.InProgress)
	activeList = km.navState.GetActiveList()
	assert.NotNil(t, activeList)
}

func TestNavigationState_GetTaskItems(t *testing.T) {
	km, cleanup := setupTestApp(t)
	defer cleanup()

	// Create task
	createTestTask(t, km, "Test Task", "Description")

	// Get items for NotStarted
	items := km.navState.GetTaskItems(domain.NotStarted)
	assert.Len(t, items, 1)

	// Get items for InProgress (empty)
	items = km.navState.GetTaskItems(domain.InProgress)
	assert.Len(t, items, 0)
}

// FormState Tests

func TestFormState_IsShowingForm_InitiallyFalse(t *testing.T) {
	km, cleanup := setupTestApp(t)
	defer cleanup()

	assert.False(t, km.uiStateManager.FormState().IsShowingForm())
}

func TestFormState_ShowTaskForm(t *testing.T) {
	km, cleanup := setupTestApp(t)
	defer cleanup()

	formState := km.uiStateManager.FormState()
	formState.ShowTaskForm([]domain.Task{})

	assert.True(t, formState.IsShowingForm())
}

func TestFormState_ShowProjectForm(t *testing.T) {
	km, cleanup := setupTestApp(t)
	defer cleanup()

	formState := km.uiStateManager.FormState()
	formState.ShowProjectForm()

	assert.True(t, formState.IsShowingForm())
}

func TestFormState_HideForm(t *testing.T) {
	km, cleanup := setupTestApp(t)
	defer cleanup()

	formState := km.uiStateManager.FormState()

	// Show form first
	formState.ShowTaskForm([]domain.Task{})
	assert.True(t, formState.IsShowingForm())

	// Hide form
	formState.HideForm()
	assert.False(t, formState.IsShowingForm())
}

func TestFormState_ShowTaskEditForm(t *testing.T) {
	km, cleanup := setupTestApp(t)
	defer cleanup()

	// Create a task
	taskID := createTestTask(t, km, "Original Task", "Original Description")
	task, err := km.taskService.GetTask(taskID)
	require.NoError(t, err)

	formState := km.uiStateManager.FormState()

	// Show edit form
	formState.ShowTaskEditForm(
		task.ID,
		task.Name,
		task.Desc,
		task.Priority,
		task.Type,
		task.BlockedBy,
		[]domain.Task{},
	)

	assert.True(t, formState.IsShowingForm())
	assert.Equal(t, task.ID, formState.GetTaskID())
}

func TestFormState_GetTaskID_EmptyForCreateForm(t *testing.T) {
	km, cleanup := setupTestApp(t)
	defer cleanup()

	formState := km.uiStateManager.FormState()
	formState.ShowTaskForm([]domain.Task{})

	taskID := formState.GetTaskID()
	assert.Empty(t, taskID)
}

func TestFormState_SetError(t *testing.T) {
	km, cleanup := setupTestApp(t)
	defer cleanup()

	formState := km.uiStateManager.FormState()
	formState.ShowTaskForm([]domain.Task{})

	// Set error
	formState.SetError("Task name is required", "name")

	errorMsg, errorField := formState.GetError()
	assert.Equal(t, "Task name is required", errorMsg)
	assert.Equal(t, "name", errorField)
}

func TestFormState_ClearError(t *testing.T) {
	km, cleanup := setupTestApp(t)
	defer cleanup()

	formState := km.uiStateManager.FormState()
	formState.ShowTaskForm([]domain.Task{})

	// Set error
	formState.SetError("Task name is required", "name")

	// Clear error
	formState.ClearError()

	errorMsg, errorField := formState.GetError()
	assert.Empty(t, errorMsg)
	assert.Empty(t, errorField)
}

func TestFormState_GetError_InitiallyEmpty(t *testing.T) {
	km, cleanup := setupTestApp(t)
	defer cleanup()

	formState := km.uiStateManager.FormState()

	errorMsg, errorField := formState.GetError()
	assert.Empty(t, errorMsg)
	assert.Empty(t, errorField)
}

func TestFormState_ShowTaskForm_ClearsError(t *testing.T) {
	km, cleanup := setupTestApp(t)
	defer cleanup()

	formState := km.uiStateManager.FormState()

	// Show form with error
	formState.ShowTaskForm([]domain.Task{})
	formState.SetError("Previous error", "field")

	// Show form again
	formState.ShowTaskForm([]domain.Task{})

	// Error should be cleared
	errorMsg, _ := formState.GetError()
	assert.Empty(t, errorMsg)
}

func TestFormState_HideForm_ClearsError(t *testing.T) {
	km, cleanup := setupTestApp(t)
	defer cleanup()

	formState := km.uiStateManager.FormState()

	// Show form with error
	formState.ShowTaskForm([]domain.Task{})
	formState.SetError("Error message", "field")

	// Hide form
	formState.HideForm()

	// Error should be cleared
	errorMsg, _ := formState.GetError()
	assert.Empty(t, errorMsg)
}
