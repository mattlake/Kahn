package app

import (
	"kahn/internal/domain"
	"kahn/internal/ui/input"
)

type FormState struct {
	showForm          bool
	activeFormType    input.FormType
	taskComponents    *input.InputComponents
	projectComponents *input.InputComponents
	formError         string
	formErrorField    string
}

func NewFormState(taskComps, projectComps *input.InputComponents) *FormState {
	return &FormState{
		taskComponents:    taskComps,
		projectComponents: projectComps,
	}
}

func (fs *FormState) ShowTaskForm() {
	fs.taskComponents.SetupForTaskCreate()
	fs.activeFormType = input.TaskCreateForm
	fs.showForm = true
	fs.ClearError()
}

func (fs *FormState) ShowTaskEditForm(taskID string, name, description string, priority domain.Priority, taskType domain.TaskType) {
	fs.taskComponents.SetupForTaskEdit(taskID, name, description, priority, taskType)
	fs.activeFormType = input.TaskEditForm
	fs.showForm = true
	fs.ClearError()
}

func (fs *FormState) ShowProjectForm() {
	fs.projectComponents.SetupForProjectCreate()
	fs.activeFormType = input.ProjectCreateForm
	fs.showForm = true
	fs.ClearError()
}

func (fs *FormState) HideForm() {
	fs.showForm = false
	fs.ClearError()
	fs.taskComponents.Reset()
	fs.projectComponents.Reset()
}

func (fs *FormState) IsShowingForm() bool {
	return fs.showForm
}

func (fs *FormState) GetActiveFormType() input.FormType {
	return fs.activeFormType
}

// GetActiveInputComponents returns input components for active form
func (fs *FormState) GetActiveInputComponents() *input.InputComponents {
	if fs.activeFormType == input.TaskCreateForm || fs.activeFormType == input.TaskEditForm {
		return fs.taskComponents
	}
	return fs.projectComponents
}

func (fs *FormState) SetError(message, field string) {
	fs.formError = message
	fs.formErrorField = field
}

func (fs *FormState) ClearError() {
	fs.formError = ""
	fs.formErrorField = ""
}

func (fs *FormState) GetError() (string, string) {
	return fs.formError, fs.formErrorField
}

func (fs *FormState) ValidateForSubmit() (bool, string, string) {
	comps := fs.GetActiveInputComponents()
	return comps.ValidateForSubmit()
}

func (fs *FormState) GetFormData() (string, string, domain.TaskType, domain.Priority) {
	comps := fs.GetActiveInputComponents()
	name := comps.NameInput.Value()
	desc := comps.DescInput.Value()
	taskType := comps.TypeValue
	priority := comps.PriorityValue
	return name, desc, taskType, priority
}

func (fs *FormState) GetTaskID() string {
	if fs.activeFormType == input.TaskEditForm {
		return fs.taskComponents.GetTaskID()
	}
	return ""
}
