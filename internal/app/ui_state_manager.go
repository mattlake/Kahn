package app

import (
	"kahn/internal/domain"
)

// ViewState represents the current UI view state
type ViewState int

const (
	BoardView ViewState = iota
	FormView
	ProjectSwitchView
	TaskDeleteConfirmView
	ProjectDeleteConfirmView
	NoProjectsView
)

// UIStateManager coordinates all UI states and provides a single source of truth
type UIStateManager struct {
	formState    *FormState
	confirmState *ConfirmationState
	navState     *NavigationState
}

// NewUIStateManager creates a new UI state manager
func NewUIStateManager(formState *FormState, confirmState *ConfirmationState, navState *NavigationState) *UIStateManager {
	return &UIStateManager{
		formState:    formState,
		confirmState: confirmState,
		navState:     navState,
	}
}

// GetCurrentViewState returns the current UI view state
func (usm *UIStateManager) GetCurrentViewState() ViewState {
	if usm.formState.IsShowingForm() {
		return FormView
	}
	if usm.navState.IsShowingProjectSwitch() {
		return ProjectSwitchView
	}
	if usm.confirmState.IsShowingTaskDeleteConfirm() {
		return TaskDeleteConfirmView
	}
	if usm.confirmState.IsShowingProjectDeleteConfirm() {
		return ProjectDeleteConfirmView
	}
	return BoardView
}

// IsShowingAnyForm returns true if any form or confirmation is active
func (usm *UIStateManager) IsShowingAnyForm() bool {
	return usm.formState.IsShowingForm() ||
		usm.navState.IsShowingProjectSwitch() ||
		usm.confirmState.IsShowingTaskDeleteConfirm() ||
		usm.confirmState.IsShowingProjectDeleteConfirm()
}

// HideAllStates hides all forms and confirmations
func (usm *UIStateManager) HideAllStates() {
	usm.formState.HideForm()
	usm.navState.HideProjectSwitch()
	usm.confirmState.HideAllConfirmations()
}

// ShowTaskForm shows the task creation form
func (usm *UIStateManager) ShowTaskForm(availableTasks []domain.Task) {
	usm.HideAllStates()
	usm.formState.ShowTaskForm(availableTasks)
}

// ShowTaskEditForm shows the task editing form
func (usm *UIStateManager) ShowTaskEditForm(taskID string, name, description string, priority domain.Priority, taskType domain.TaskType, blockedByIntID *int, availableTasks []domain.Task) {
	usm.HideAllStates()
	usm.formState.ShowTaskEditForm(taskID, name, description, priority, taskType, blockedByIntID, availableTasks)
}

// ShowProjectForm shows the project creation form
func (usm *UIStateManager) ShowProjectForm() {
	usm.HideAllStates()
	usm.formState.ShowProjectForm()
}

// ShowProjectSwitcher shows the project switcher
func (usm *UIStateManager) ShowProjectSwitcher() {
	usm.HideAllStates()
	usm.navState.ShowProjectSwitch()
}

// ShowTaskDeleteConfirm shows the task delete confirmation
func (usm *UIStateManager) ShowTaskDeleteConfirm(taskID string) {
	usm.HideAllStates()
	usm.confirmState.ShowTaskDeleteConfirm(taskID)
}

// ShowProjectDeleteConfirm shows the project delete confirmation
func (usm *UIStateManager) ShowProjectDeleteConfirm(projectID string) {
	usm.HideAllStates()
	usm.confirmState.ShowProjectDeleteConfirm(projectID)
}

// Getter methods for accessing specific state managers
func (usm *UIStateManager) FormState() *FormState {
	return usm.formState
}

func (usm *UIStateManager) ConfirmationState() *ConfirmationState {
	return usm.confirmState
}

func (usm *UIStateManager) NavigationState() *NavigationState {
	return usm.navState
}
