package app

// ConfirmationState manages confirmation dialogs for deletions
type ConfirmationState struct {
	showTaskDeleteConfirm    bool
	showProjectDeleteConfirm bool
	taskToDelete             string
	projectToDelete          string
}

// NewConfirmationState creates a new confirmation state
func NewConfirmationState() *ConfirmationState {
	return &ConfirmationState{}
}

// ShowTaskDeleteConfirm shows the task deletion confirmation
func (cs *ConfirmationState) ShowTaskDeleteConfirm(taskID string) {
	cs.showTaskDeleteConfirm = true
	cs.taskToDelete = taskID
}

// ShowProjectDeleteConfirm shows the project deletion confirmation
func (cs *ConfirmationState) ShowProjectDeleteConfirm(projectID string) {
	cs.showProjectDeleteConfirm = true
	cs.projectToDelete = projectID
}

// HideAllConfirmations hides all confirmation dialogs
func (cs *ConfirmationState) HideAllConfirmations() {
	cs.showTaskDeleteConfirm = false
	cs.showProjectDeleteConfirm = false
	cs.taskToDelete = ""
	cs.projectToDelete = ""
}

// IsShowingTaskDeleteConfirm returns whether task delete confirmation is shown
func (cs *ConfirmationState) IsShowingTaskDeleteConfirm() bool {
	return cs.showTaskDeleteConfirm
}

// IsShowingProjectDeleteConfirm returns whether project delete confirmation is shown
func (cs *ConfirmationState) IsShowingProjectDeleteConfirm() bool {
	return cs.showProjectDeleteConfirm
}

// GetTaskToDelete returns the ID of the task to delete
func (cs *ConfirmationState) GetTaskToDelete() string {
	return cs.taskToDelete
}

// GetProjectToDelete returns the ID of the project to delete
func (cs *ConfirmationState) GetProjectToDelete() string {
	return cs.projectToDelete
}

// ClearTaskDelete clears the task deletion state
func (cs *ConfirmationState) ClearTaskDelete() {
	cs.showTaskDeleteConfirm = false
	cs.taskToDelete = ""
}

// ClearProjectDelete clears the project deletion state
func (cs *ConfirmationState) ClearProjectDelete() {
	cs.showProjectDeleteConfirm = false
	cs.projectToDelete = ""
}
