package app

type ConfirmationState struct {
	showTaskDeleteConfirm    bool
	showProjectDeleteConfirm bool
	taskToDelete             string
	projectToDelete          string
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

func (cs *ConfirmationState) ClearTaskDelete() {
	cs.showTaskDeleteConfirm = false
	cs.taskToDelete = ""
}

func (cs *ConfirmationState) ClearProjectDelete() {
	cs.showProjectDeleteConfirm = false
	cs.projectToDelete = ""
}
