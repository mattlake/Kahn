package app

type ConfirmationState struct {
	showTaskDeleteConfirm    bool
	showProjectDeleteConfirm bool
	taskToDelete             string
	projectToDelete          string
	errorMessage             string
}

func NewConfirmationState() *ConfirmationState {
	return &ConfirmationState{}
}

func (cs *ConfirmationState) ShowTaskDeleteConfirm(taskID string) {
	cs.showTaskDeleteConfirm = true
	cs.taskToDelete = taskID
}

func (cs *ConfirmationState) ShowProjectDeleteConfirm(projectID string) {
	cs.showProjectDeleteConfirm = true
	cs.projectToDelete = projectID
}

func (cs *ConfirmationState) HideAllConfirmations() {
	cs.showTaskDeleteConfirm = false
	cs.showProjectDeleteConfirm = false
	cs.taskToDelete = ""
	cs.projectToDelete = ""
	cs.errorMessage = ""
}

func (cs *ConfirmationState) IsShowingTaskDeleteConfirm() bool {
	return cs.showTaskDeleteConfirm
}

func (cs *ConfirmationState) IsShowingProjectDeleteConfirm() bool {
	return cs.showProjectDeleteConfirm
}

func (cs *ConfirmationState) GetTaskToDelete() string {
	return cs.taskToDelete
}

func (cs *ConfirmationState) GetProjectToDelete() string {
	return cs.projectToDelete
}

func (cs *ConfirmationState) SetError(message string) {
	cs.errorMessage = message
}

func (cs *ConfirmationState) GetError() string {
	return cs.errorMessage
}

func (cs *ConfirmationState) HasError() bool {
	return cs.errorMessage != ""
}

func (cs *ConfirmationState) ClearTaskDelete() {
	cs.showTaskDeleteConfirm = false
	cs.taskToDelete = ""
	cs.errorMessage = ""
}

func (cs *ConfirmationState) ClearProjectDelete() {
	cs.showProjectDeleteConfirm = false
	cs.projectToDelete = ""
	cs.errorMessage = ""
}
