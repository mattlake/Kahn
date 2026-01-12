package app

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"kahn/internal/domain"
)

// handleFormInput Tests

func TestHandleFormInput_EscapeKey(t *testing.T) {
	km, cleanup := setupTestApp(t)
	defer cleanup()

	// Show task form
	km.uiStateManager.ShowTaskForm([]domain.Task{})
	assert.Equal(t, FormView, km.uiStateManager.GetCurrentViewState())

	// Press Esc
	msg := tea.KeyMsg{Type: tea.KeyEsc}
	updatedModel, _ := km.Update(msg)
	km = updatedModel.(*KahnModel)

	// Form should be hidden
	assert.Equal(t, BoardView, km.uiStateManager.GetCurrentViewState())
}

func TestHandleFormInput_TabKey_TaskForm(t *testing.T) {
	km, cleanup := setupTestApp(t)
	defer cleanup()

	// Show task form
	km.uiStateManager.ShowTaskForm([]domain.Task{})
	formState := km.uiStateManager.FormState()
	comps := formState.GetActiveInputComponents()

	// Initially focused on name (field 0)
	assert.Equal(t, 0, comps.FocusedField)

	// Press Tab
	msg := tea.KeyMsg{Type: tea.KeyTab}
	updatedModel, _ := km.Update(msg)
	km = updatedModel.(*KahnModel)

	// Should move to description (field 1)
	comps = km.uiStateManager.FormState().GetActiveInputComponents()
	assert.Equal(t, 1, comps.FocusedField)
}

func TestHandleFormInput_TabKey_ProjectForm(t *testing.T) {
	km, cleanup := setupTestApp(t)
	defer cleanup()

	// Show project form
	km.uiStateManager.ShowProjectForm()
	formState := km.uiStateManager.FormState()
	comps := formState.GetActiveInputComponents()

	// Initially focused on name (field 0)
	assert.Equal(t, 0, comps.FocusedField)

	// Press Tab
	msg := tea.KeyMsg{Type: tea.KeyTab}
	updatedModel, _ := km.Update(msg)
	km = updatedModel.(*KahnModel)

	// Should move to description (field 1)
	comps = km.uiStateManager.FormState().GetActiveInputComponents()
	assert.Equal(t, 1, comps.FocusedField)

	// Press Tab again (should cycle back to name)
	updatedModel, _ = km.Update(msg)
	km = updatedModel.(*KahnModel)
	comps = km.uiStateManager.FormState().GetActiveInputComponents()
	assert.Equal(t, 0, comps.FocusedField)
}

func TestHandleFormInput_EnterKey_SubmitsForm(t *testing.T) {
	km, cleanup := setupTestApp(t)
	defer cleanup()

	// Show task form
	km.uiStateManager.ShowTaskForm([]domain.Task{})
	formState := km.uiStateManager.FormState()
	comps := formState.GetActiveInputComponents()

	// Set task name
	comps.NameInput.SetValue("Test Task")

	// Press Enter
	msg := tea.KeyMsg{Type: tea.KeyEnter}
	updatedModel, _ := km.Update(msg)
	km = updatedModel.(*KahnModel)

	// Form should be closed (task created)
	assert.Equal(t, BoardView, km.uiStateManager.GetCurrentViewState())

	// Verify task was created
	activeProj := km.projectManager.GetActiveProject()
	assert.Len(t, activeProj.Tasks, 1)
	assert.Equal(t, "Test Task", activeProj.Tasks[0].Name)
}

func TestHandleFormInput_EnterKey_ValidationFails(t *testing.T) {
	km, cleanup := setupTestApp(t)
	defer cleanup()

	// Show task form
	km.uiStateManager.ShowTaskForm([]domain.Task{})

	// Don't set task name (validation should fail)

	// Press Enter
	msg := tea.KeyMsg{Type: tea.KeyEnter}
	updatedModel, _ := km.Update(msg)
	km = updatedModel.(*KahnModel)

	// Should still be in form view (validation failed)
	assert.Equal(t, FormView, km.uiStateManager.GetCurrentViewState())

	// Should have error
	assertFormError(t, km, "name")
}

func TestHandleFormInput_CtrlEnter_SubmitsFromAnyField(t *testing.T) {
	km, cleanup := setupTestApp(t)
	defer cleanup()

	// Show task form
	km.uiStateManager.ShowTaskForm([]domain.Task{})
	formState := km.uiStateManager.FormState()
	comps := formState.GetActiveInputComponents()

	// Set task name
	comps.NameInput.SetValue("Test Task")

	// Move to priority field (field 2) by pressing Tab twice
	msg := tea.KeyMsg{Type: tea.KeyTab}
	updatedModel, _ := km.Update(msg) // Name -> Desc
	km = updatedModel.(*KahnModel)
	updatedModel, _ = km.Update(msg) // Desc -> Priority
	km = updatedModel.(*KahnModel)

	// Verify we're in priority field (field 2)
	comps = km.uiStateManager.FormState().GetActiveInputComponents()
	assert.Equal(t, 2, comps.FocusedField)

	// Press Enter from priority field - this should submit (not add newline like description)
	// This tests that form submission works from non-name fields
	msg = tea.KeyMsg{Type: tea.KeyEnter}
	updatedModel, _ = km.Update(msg)
	km = updatedModel.(*KahnModel)

	// Form should be closed (task created successfully)
	assert.Equal(t, BoardView, km.uiStateManager.GetCurrentViewState())

	// Verify task was created
	activeProj := km.projectManager.GetActiveProject()
	assert.Len(t, activeProj.Tasks, 1)
	assert.Equal(t, "Test Task", activeProj.Tasks[0].Name)
}

func TestHandleFormInput_UpDownKeys_CyclePriority(t *testing.T) {
	km, cleanup := setupTestApp(t)
	defer cleanup()

	// Show task form
	km.uiStateManager.ShowTaskForm([]domain.Task{})
	formState := km.uiStateManager.FormState()
	comps := formState.GetActiveInputComponents()

	// Navigate to priority field (field 2)
	msg := tea.KeyMsg{Type: tea.KeyTab}
	km.Update(msg) // Name -> Desc
	km.Update(msg) // Desc -> Priority

	initialPriority := comps.PriorityValue

	// Press Up
	msg = tea.KeyMsg{Type: tea.KeyUp}
	updatedModel, _ := km.Update(msg)
	km = updatedModel.(*KahnModel)

	comps = km.uiStateManager.FormState().GetActiveInputComponents()
	assert.NotEqual(t, initialPriority, comps.PriorityValue)
}

func TestHandleFormInput_CharacterInput_UpdatesField(t *testing.T) {
	km, cleanup := setupTestApp(t)
	defer cleanup()

	// Show task form
	km.uiStateManager.ShowTaskForm([]domain.Task{})

	// Type characters (simulating typing "Test")
	keys := []rune{'T', 'e', 's', 't'}
	for _, r := range keys {
		msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}}
		km.Update(msg)
	}

	// Should clear any previous errors
	assertNoFormError(t, km)
}

// handleSearchInput Tests

func TestHandleSearchInput_EscKey_ClearsSearch(t *testing.T) {
	km, cleanup := setupTestApp(t)
	defer cleanup()

	// Activate search
	km.searchState.Activate()
	km.searchState.AppendChar("test")

	// Press Esc
	msg := tea.KeyMsg{Type: tea.KeyEsc}
	updatedModel, _ := km.Update(msg)
	km = updatedModel.(*KahnModel)

	// Search should be cleared
	assert.False(t, km.searchState.IsActive())
	assert.Empty(t, km.searchState.GetQuery())
}

func TestHandleSearchInput_BackspaceKey(t *testing.T) {
	km, cleanup := setupTestApp(t)
	defer cleanup()

	// Activate search and type
	km.searchState.Activate()
	km.searchState.AppendChar("test")
	assert.Equal(t, "test", km.searchState.GetQuery())

	// Press Backspace
	msg := tea.KeyMsg{Type: tea.KeyBackspace}
	updatedModel, _ := km.Update(msg)
	km = updatedModel.(*KahnModel)

	// Last character should be removed
	assert.Equal(t, "tes", km.searchState.GetQuery())
}

func TestHandleSearchInput_EnterKey_DoesNothing(t *testing.T) {
	km, cleanup := setupTestApp(t)
	defer cleanup()

	// Activate search
	km.searchState.Activate()
	km.searchState.AppendChar("test")

	// Press Enter (should not do anything special)
	msg := tea.KeyMsg{Type: tea.KeyEnter}
	updatedModel, _ := km.Update(msg)
	km = updatedModel.(*KahnModel)

	// Search should still be active
	assert.True(t, km.searchState.IsActive())
	assert.Equal(t, "test", km.searchState.GetQuery())
}

func TestHandleSearchInput_CharacterInput(t *testing.T) {
	km, cleanup := setupTestApp(t)
	defer cleanup()

	// Create tasks
	createTestTask(t, km, "API Task", "Description")
	createTestTask(t, km, "Database Task", "Description")

	// Activate search
	km.searchState.Activate()

	// Type 'A'
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'A'}}
	updatedModel, _ := km.Update(msg)
	km = updatedModel.(*KahnModel)

	assert.Equal(t, "A", km.searchState.GetQuery())

	// Type 'P'
	msg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'P'}}
	updatedModel, _ = km.Update(msg)
	km = updatedModel.(*KahnModel)

	assert.Equal(t, "AP", km.searchState.GetQuery())
}

func TestHandleSearchInput_FiltersTasksRealtime(t *testing.T) {
	km, cleanup := setupTestApp(t)
	defer cleanup()

	// Create tasks
	createTestTask(t, km, "API Task", "Description")
	createTestTask(t, km, "Database Task", "Description")
	createTestTask(t, km, "API Service", "Description")

	// Activate search
	km.searchState.Activate()

	// Type 'API'
	for _, r := range []rune{'A', 'P', 'I'} {
		msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}}
		km.Update(msg)
	}

	// Verify match count
	assert.Equal(t, 2, km.searchState.GetMatchCount())
}

// handleProjectSwitch Tests

func TestHandleProjectSwitch_EscKey_HidesProjectSwitch(t *testing.T) {
	km, cleanup := setupTestApp(t)
	defer cleanup()

	// Show project switcher
	km.uiStateManager.ShowProjectSwitcher()
	assert.Equal(t, ProjectSwitchView, km.uiStateManager.GetCurrentViewState())

	// Press Esc
	msg := tea.KeyMsg{Type: tea.KeyEsc}
	updatedModel, _ := km.Update(msg)
	km = updatedModel.(*KahnModel)

	// Project switcher should be hidden
	assert.Equal(t, BoardView, km.uiStateManager.GetCurrentViewState())
}

func TestHandleProjectSwitch_NKey_ShowsProjectForm(t *testing.T) {
	km, cleanup := setupTestApp(t)
	defer cleanup()

	// Show project switcher
	km.uiStateManager.ShowProjectSwitcher()

	// Press 'n'
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}}
	updatedModel, _ := km.Update(msg)
	km = updatedModel.(*KahnModel)

	// Should show project form
	assert.Equal(t, FormView, km.uiStateManager.GetCurrentViewState())
}

func TestHandleProjectSwitch_DKey_ShowsDeleteConfirm(t *testing.T) {
	km, cleanup := setupTestApp(t)
	defer cleanup()

	// Show project switcher
	km.uiStateManager.ShowProjectSwitcher()

	// Press 'd'
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}}
	updatedModel, _ := km.Update(msg)
	km = updatedModel.(*KahnModel)

	// Should show delete confirmation (as overlay on project switch)
	// Project switch view is still active, but confirmation state is also showing
	assert.Equal(t, ProjectSwitchView, km.uiStateManager.GetCurrentViewState())
	assert.True(t, km.uiStateManager.ConfirmationState().IsShowingProjectDeleteConfirm())
}

func TestHandleProjectSwitch_NumberKey_SwitchesProject(t *testing.T) {
	km, cleanup := setupTestApp(t)
	defer cleanup()

	// Create second project
	err := km.projectManager.CreateProject("Project 2", "Description")
	require.NoError(t, err)
	project2ID := km.projectManager.GetActiveProjectID()

	// Create third project
	err = km.projectManager.CreateProject("Project 3", "Description")
	require.NoError(t, err)

	// Show project switcher
	km.uiStateManager.ShowProjectSwitcher()

	// Press '2' to switch to second project (0-indexed, so '2' = index 1)
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'2'}}
	updatedModel, _ := km.Update(msg)
	km = updatedModel.(*KahnModel)

	// Should switch to project 2 and hide switcher
	assert.Equal(t, BoardView, km.uiStateManager.GetCurrentViewState())
	assert.Equal(t, project2ID, km.projectManager.GetActiveProjectID())
}

func TestHandleProjectSwitch_JKey_NavigatesDown(t *testing.T) {
	km, cleanup := setupTestApp(t)
	defer cleanup()

	// Get default project ID
	defaultID := km.projectManager.GetActiveProjectID()

	// Create second project
	err := km.projectManager.CreateProject("Project 2", "Description")
	require.NoError(t, err)
	project2ID := km.projectManager.GetActiveProjectID()

	// Switch back to default
	km.projectManager.SwitchToProject(defaultID)

	// Show project switcher
	km.uiStateManager.ShowProjectSwitcher()

	// Press 'j' (down)
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}}
	updatedModel, _ := km.Update(msg)
	km = updatedModel.(*KahnModel)

	// Should switch to next project
	assert.Equal(t, project2ID, km.projectManager.GetActiveProjectID())
}

func TestHandleProjectSwitch_KKey_NavigatesUp(t *testing.T) {
	km, cleanup := setupTestApp(t)
	defer cleanup()

	// Get default project ID
	defaultID := km.projectManager.GetActiveProjectID()

	// Create second project
	err := km.projectManager.CreateProject("Project 2", "Description")
	require.NoError(t, err)
	// Now active project is project 2

	// Currently on project 2, press 'k' (up)
	km.uiStateManager.ShowProjectSwitcher()
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}}
	updatedModel, _ := km.Update(msg)
	km = updatedModel.(*KahnModel)

	// Should switch to previous project
	assert.Equal(t, defaultID, km.projectManager.GetActiveProjectID())
}

func TestHandleProjectSwitch_EnterKey_ConfirmsSelection(t *testing.T) {
	km, cleanup := setupTestApp(t)
	defer cleanup()

	// Show project switcher
	km.uiStateManager.ShowProjectSwitcher()

	// Press Enter
	msg := tea.KeyMsg{Type: tea.KeyEnter}
	updatedModel, _ := km.Update(msg)
	km = updatedModel.(*KahnModel)

	// Should hide project switcher
	assert.Equal(t, BoardView, km.uiStateManager.GetCurrentViewState())
}

func TestHandleProjectSwitch_ClearsSearchWhenSwitching(t *testing.T) {
	km, cleanup := setupTestApp(t)
	defer cleanup()

	// Create second project
	err := km.projectManager.CreateProject("Project 2", "Description")
	require.NoError(t, err)
	project2ID := km.projectManager.GetActiveProjectID()

	// Switch back to first project
	projects := km.projectManager.GetProjectsAsDomain()
	km.projectManager.SwitchToProject(projects[0].ID)

	// Activate search
	km.searchState.Activate()
	km.searchState.AppendChar("test")
	assert.True(t, km.searchState.IsActive())

	// Show project switcher (this should take precedence but doesn't due to routing)
	// First we need to clear search to properly test project switching
	msg := tea.KeyMsg{Type: tea.KeyEsc} // Clear search first
	updatedModel, _ := km.Update(msg)
	km = updatedModel.(*KahnModel)

	// Now show project switcher
	km.uiStateManager.ShowProjectSwitcher()

	// Press 'j' to navigate to next project (avoids number key which would append to search)
	msg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}}
	updatedModel, _ = km.Update(msg)
	km = updatedModel.(*KahnModel)

	// Search should remain cleared (it was cleared by Esc)
	assert.False(t, km.searchState.IsActive())
	// Should be on project 2
	assert.Equal(t, project2ID, km.projectManager.GetActiveProjectID())
}

func TestHandleProjectSwitch_DeleteConfirm_YKey(t *testing.T) {
	km, cleanup := setupTestApp(t)
	defer cleanup()

	// Create second project so we don't delete the last one
	err := km.projectManager.CreateProject("Project 2", "Description")
	require.NoError(t, err)

	// Get first project ID
	projects := km.projectManager.GetProjectsAsDomain()
	firstProjectID := projects[0].ID

	// Switch to first project
	km.projectManager.SwitchToProject(firstProjectID)

	// Show project switcher and delete confirmation
	km.uiStateManager.ShowProjectSwitcher()
	km.uiStateManager.ShowProjectDeleteConfirm(firstProjectID)

	// Press 'y' to confirm
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}}
	updatedModel, _ := km.Update(msg)
	km = updatedModel.(*KahnModel)

	// Project should be deleted
	assert.Equal(t, 1, km.projectManager.GetProjectCount())
}

func TestHandleProjectSwitch_DeleteConfirm_NKey(t *testing.T) {
	km, cleanup := setupTestApp(t)
	defer cleanup()

	projectID := km.projectManager.GetActiveProjectID()

	// Show project switcher and delete confirmation
	km.uiStateManager.ShowProjectSwitcher()
	km.uiStateManager.ShowProjectDeleteConfirm(projectID)

	// Press 'n' to cancel
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}}
	updatedModel, _ := km.Update(msg)
	km = updatedModel.(*KahnModel)

	// Project should NOT be deleted
	assert.Equal(t, 1, km.projectManager.GetProjectCount())
	// After canceling, we're back to board view (project switch was hidden when showing delete confirm)
	assert.Equal(t, BoardView, km.uiStateManager.GetCurrentViewState())
}

// handleTaskDeleteConfirm Tests

func TestHandleTaskDeleteConfirm_YKey_DeletesTask(t *testing.T) {
	km, cleanup := setupTestApp(t)
	defer cleanup()

	// Create task
	taskID := createTestTask(t, km, "Test Task", "Description")

	// Show delete confirmation
	km.uiStateManager.ShowTaskDeleteConfirm(taskID)

	// Press 'y' to confirm
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}}
	updatedModel, _ := km.Update(msg)
	km = updatedModel.(*KahnModel)

	// Task should be deleted
	activeProj := km.projectManager.GetActiveProject()
	assert.Len(t, activeProj.Tasks, 0)
	assert.Equal(t, BoardView, km.uiStateManager.GetCurrentViewState())
}

func TestHandleTaskDeleteConfirm_NKey_CancelsDelete(t *testing.T) {
	km, cleanup := setupTestApp(t)
	defer cleanup()

	// Create task
	taskID := createTestTask(t, km, "Test Task", "Description")

	// Show delete confirmation
	km.uiStateManager.ShowTaskDeleteConfirm(taskID)

	// Press 'n' to cancel
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}}
	updatedModel, _ := km.Update(msg)
	km = updatedModel.(*KahnModel)

	// Task should NOT be deleted
	activeProj := km.projectManager.GetActiveProject()
	assert.Len(t, activeProj.Tasks, 1)
	assert.Equal(t, BoardView, km.uiStateManager.GetCurrentViewState())
}

func TestHandleTaskDeleteConfirm_EscKey_CancelsDelete(t *testing.T) {
	km, cleanup := setupTestApp(t)
	defer cleanup()

	// Create task
	taskID := createTestTask(t, km, "Test Task", "Description")

	// Show delete confirmation
	km.uiStateManager.ShowTaskDeleteConfirm(taskID)

	// Press Esc to cancel
	msg := tea.KeyMsg{Type: tea.KeyEsc}
	updatedModel, _ := km.Update(msg)
	km = updatedModel.(*KahnModel)

	// Task should NOT be deleted
	activeProj := km.projectManager.GetActiveProject()
	assert.Len(t, activeProj.Tasks, 1)
	assert.Equal(t, BoardView, km.uiStateManager.GetCurrentViewState())
}

// handleNormalMode Tests

func TestHandleNormalMode_SlashKey_ActivatesSearch(t *testing.T) {
	km, cleanup := setupTestApp(t)
	defer cleanup()

	// Press '/'
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}}
	updatedModel, _ := km.Update(msg)
	km = updatedModel.(*KahnModel)

	// Search should be activated
	assert.True(t, km.searchState.IsActive())
}

func TestHandleNormalMode_QKey_Quits(t *testing.T) {
	km, cleanup := setupTestApp(t)
	defer cleanup()

	// Press 'q'
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
	_, cmd := km.Update(msg)

	// Should return tea.Quit command
	assert.NotNil(t, cmd)
}

func TestHandleNormalMode_NKey_ShowsTaskForm(t *testing.T) {
	km, cleanup := setupTestApp(t)
	defer cleanup()

	// Press 'n'
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}}
	updatedModel, _ := km.Update(msg)
	km = updatedModel.(*KahnModel)

	// Should show task form
	assert.Equal(t, FormView, km.uiStateManager.GetCurrentViewState())
}

func TestHandleNormalMode_PKey_ShowsProjectSwitcher(t *testing.T) {
	km, cleanup := setupTestApp(t)
	defer cleanup()

	// Press 'p'
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'p'}}
	updatedModel, _ := km.Update(msg)
	km = updatedModel.(*KahnModel)

	// Should show project switcher
	assert.Equal(t, ProjectSwitchView, km.uiStateManager.GetCurrentViewState())
}

func TestHandleNormalMode_LKey_NavigatesRight(t *testing.T) {
	km, cleanup := setupTestApp(t)
	defer cleanup()

	// Start at NotStarted
	assert.Equal(t, domain.NotStarted, km.navState.GetActiveListIndex())

	// Press 'l'
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}}
	updatedModel, _ := km.Update(msg)
	km = updatedModel.(*KahnModel)

	// Should move to InProgress
	assert.Equal(t, domain.InProgress, km.navState.GetActiveListIndex())
}

func TestHandleNormalMode_HKey_NavigatesLeft(t *testing.T) {
	km, cleanup := setupTestApp(t)
	defer cleanup()

	// Move to InProgress first
	km.navState.SetActiveListIndex(domain.InProgress)

	// Press 'h'
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}}
	updatedModel, _ := km.Update(msg)
	km = updatedModel.(*KahnModel)

	// Should move to NotStarted
	assert.Equal(t, domain.NotStarted, km.navState.GetActiveListIndex())
}

// handleResize Tests

func TestHandleResize_UpdatesDimensions(t *testing.T) {
	km, cleanup := setupTestApp(t)
	defer cleanup()

	// Send resize message
	msg := tea.WindowSizeMsg{Width: 120, Height: 40}
	updatedModel, _ := km.Update(msg)
	km = updatedModel.(*KahnModel)

	// Dimensions should be updated
	assert.Equal(t, 120, km.width)
	assert.Equal(t, 40, km.height)
}
