package app

type ConfirmationState struct {
	taskDeleteConfirm    *GenericConfirmationState[string]
	projectDeleteConfirm *GenericConfirmationState[string]
}

func NewConfirmationState() *ConfirmationState {
	return &ConfirmationState{
		taskDeleteConfirm:    NewGenericConfirmationState[string](),
		projectDeleteConfirm: NewGenericConfirmationState[string](),
	}
}

func (cs *ConfirmationState) ShowTaskDeleteConfirm(taskID string) {
	cs.taskDeleteConfirm.ShowConfirm(taskID)
	cs.projectDeleteConfirm.HideConfirm()
}

func (cs *ConfirmationState) ShowProjectDeleteConfirm(projectID string) {
	cs.projectDeleteConfirm.ShowConfirm(projectID)
	cs.taskDeleteConfirm.HideConfirm()
}

func (cs *ConfirmationState) HideAllConfirmations() {
	cs.taskDeleteConfirm.HideConfirm()
	cs.projectDeleteConfirm.HideConfirm()
}

func (cs *ConfirmationState) IsShowingTaskDeleteConfirm() bool {
	return cs.taskDeleteConfirm.IsShowingConfirm()
}

func (cs *ConfirmationState) IsShowingProjectDeleteConfirm() bool {
	return cs.projectDeleteConfirm.IsShowingConfirm()
}

func (cs *ConfirmationState) GetTaskToDelete() string {
	return cs.taskDeleteConfirm.GetItemToDelete()
}

func (cs *ConfirmationState) GetProjectToDelete() string {
	return cs.projectDeleteConfirm.GetItemToDelete()
}

func (cs *ConfirmationState) SetTaskError(message string) {
	cs.taskDeleteConfirm.SetError(message)
}

func (cs *ConfirmationState) SetProjectError(message string) {
	cs.projectDeleteConfirm.SetError(message)
}

func (cs *ConfirmationState) GetTaskError() string {
	return cs.taskDeleteConfirm.GetError()
}

func (cs *ConfirmationState) GetProjectError() string {
	return cs.projectDeleteConfirm.GetError()
}

func (cs *ConfirmationState) HasTaskError() bool {
	return cs.taskDeleteConfirm.HasError()
}

func (cs *ConfirmationState) HasProjectError() bool {
	return cs.projectDeleteConfirm.HasError()
}

func (cs *ConfirmationState) ClearTaskDelete() {
	cs.taskDeleteConfirm.HideConfirm()
}

func (cs *ConfirmationState) ClearProjectDelete() {
	cs.projectDeleteConfirm.HideConfirm()
}

// Legacy compatibility methods
func (cs *ConfirmationState) GetError() string {
	if cs.taskDeleteConfirm.IsShowingConfirm() {
		return cs.taskDeleteConfirm.GetError()
	}
	return cs.projectDeleteConfirm.GetError()
}

func (cs *ConfirmationState) SetError(message string) {
	if cs.taskDeleteConfirm.IsShowingConfirm() {
		cs.taskDeleteConfirm.SetError(message)
	} else {
		cs.projectDeleteConfirm.SetError(message)
	}
}

func (cs *ConfirmationState) HasError() bool {
	return cs.taskDeleteConfirm.HasError() || cs.projectDeleteConfirm.HasError()
}

func (cs *ConfirmationState) ClearError() {
	cs.taskDeleteConfirm.ClearError()
	cs.projectDeleteConfirm.ClearError()
}
