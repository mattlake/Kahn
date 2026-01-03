package app

import (
	"kahn/internal/domain"
	"kahn/internal/ui/input"
)

// FormState manages all form-related UI state
type FormState struct {
	showForm          bool
	activeFormType    input.FormType
	taskComponents    *input.InputComponents
	projectComponents *input.InputComponents
	formError         string
	formErrorField    string
}

// NewFormState creates a new form state
func NewFormState(taskComps, projectComps *input.InputComponents) *FormState {
	return &FormState{
		taskComponents:    taskComps,
		projectComponents: projectComps,
	}
}

// ShowTaskForm displays the task creation form
func (fs *FormState) ShowTaskForm() {
	fs.taskComponents.SetupForTaskCreate()
	fs.activeFormType = input.TaskCreateForm
	fs.showForm = true
	fs.ClearError()
}

// ShowTaskEditForm displays the task edit form
func (fs *FormState) ShowTaskEditForm(taskID string, name, description string, priority domain.Priority, taskType domain.TaskType) {
	fs.taskComponents.SetupForTaskEdit(taskID, name, description, priority, taskType)
	fs.activeFormType = input.TaskEditForm
	fs.showForm = true
	fs.ClearError()
}

// ShowProjectForm displays the project creation form
func (fs *FormState) ShowProjectForm() {
	fs.projectComponents.SetupForProjectCreate()
	fs.activeFormType = input.ProjectCreateForm
	fs.showForm = true
	fs.ClearError()
}

// HideForm hides all forms and resets state
func (fs *FormState) HideForm() {
	fs.showForm = false
	fs.ClearError()
	fs.taskComponents.Reset()
	fs.projectComponents.Reset()
}

// IsShowingForm returns whether any form is currently displayed
func (fs *FormState) IsShowingForm() bool {
	return fs.showForm
}

// GetActiveFormType returns the currently active form type
func (fs *FormState) GetActiveFormType() input.FormType {
	return fs.activeFormType
}

// GetActiveInputComponents returns the input components for the active form
func (fs *FormState) GetActiveInputComponents() *input.InputComponents {
	if fs.activeFormType == input.TaskCreateForm || fs.activeFormType == input.TaskEditForm {
		return fs.taskComponents
	}
	return fs.projectComponents
}

// SetError sets a form error with field information
func (fs *FormState) SetError(message, field string) {
	fs.formError = message
	fs.formErrorField = field
}

// ClearError clears any form errors
func (fs *FormState) ClearError() {
	fs.formError = ""
	fs.formErrorField = ""
}

// GetError returns the current form error
func (fs *FormState) GetError() (string, string) {
	return fs.formError, fs.formErrorField
}

// ValidateForSubmit validates the current form for submission
func (fs *FormState) ValidateForSubmit() (bool, string, string) {
	comps := fs.GetActiveInputComponents()
	return comps.ValidateForSubmit()
}

// GetFormData returns the current form data
func (fs *FormState) GetFormData() (string, string, domain.TaskType, domain.Priority) {
	comps := fs.GetActiveInputComponents()
	name := comps.NameInput.Value()
	desc := comps.DescInput.Value()
	taskType := comps.TypeValue
	priority := comps.PriorityValue
	return name, desc, taskType, priority
}

// GetTaskID returns the task ID for edit forms
func (fs *FormState) GetTaskID() string {
	if fs.activeFormType == input.TaskEditForm {
		return fs.taskComponents.GetTaskID()
	}
	return ""
}
