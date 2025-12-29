package input

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	"kahn/internal/ui/colors"
)

type FormType int

const (
	TaskCreateForm FormType = iota
	TaskEditForm
	ProjectCreateForm
)

// InputComponents holds the text input components for forms
type InputComponents struct {
	NameInput    textinput.Model
	DescInput    textinput.Model
	formType     FormType
	taskID       string // for edit forms
	FocusedField int    // 0=name, 1=desc (exported)
}

// NewInputComponents creates and initializes input components for forms
func NewInputComponents() InputComponents {
	return InputComponents{
		formType:     TaskCreateForm,
		FocusedField: 0,
	}
}

// SetupForTaskCreate configures components for task creation
func (ic *InputComponents) SetupForTaskCreate() {
	ic.formType = TaskCreateForm
	ic.FocusedField = 0
	ic.taskID = ""
	ic.NameInput = ic.createNameInput("Task name *")
	ic.DescInput = ic.createDescInput("Task description (optional)")
	ic.NameInput.Focus()
}

func (ic *InputComponents) SetupForTaskEdit(taskID, name, desc string) {
	ic.formType = TaskEditForm
	ic.FocusedField = 0
	ic.taskID = taskID
	ic.NameInput = ic.createNameInput("Task name *")
	ic.DescInput = ic.createDescInput("Task description (optional)")
	ic.NameInput.SetValue(name)
	ic.DescInput.SetValue(desc)
	ic.NameInput.Focus()
}

func (ic *InputComponents) SetupForProjectCreate() {
	ic.formType = ProjectCreateForm
	ic.FocusedField = 0
	ic.taskID = ""
	ic.NameInput = ic.createNameInput("Project name *")
	ic.DescInput = ic.createDescInput("Project description (optional)")
	ic.NameInput.Focus()
}

func (ic *InputComponents) createNameInput(placeholder string) textinput.Model {
	input := textinput.New()
	input.Placeholder = placeholder
	input.PlaceholderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(colors.Subtext0))
	input.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(colors.Text))
	input.Cursor.Style = lipgloss.NewStyle().Foreground(lipgloss.Color(colors.Mauve))
	input.CharLimit = 50
	input.Width = 40
	return input
}

func (ic *InputComponents) createDescInput(placeholder string) textinput.Model {
	input := textinput.New()
	input.Placeholder = placeholder
	input.PlaceholderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(colors.Subtext0))
	input.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(colors.Text))
	input.Cursor.Style = lipgloss.NewStyle().Foreground(lipgloss.Color(colors.Mauve))
	input.CharLimit = 200 // Larger for project descriptions
	input.Width = 40
	return input
}

// Reset resets all input components
func (ic *InputComponents) Reset() {
	ic.NameInput.Reset()
	ic.DescInput.Reset()
	ic.FocusedField = 0
	ic.taskID = ""
}

// Validate performs basic validation without field-specific errors
func (ic *InputComponents) Validate() error {
	_, _, errorMsg := ic.ValidateForSubmit()
	if errorMsg != "" {
		return fmt.Errorf("validation failed: %s", errorMsg)
	}
	return nil
}

// ValidateForSubmit performs validation when form is submitted
// Returns: (isValid, errorField, errorMessage)
func (ic *InputComponents) ValidateForSubmit() (bool, string, string) {
	name := strings.TrimSpace(ic.NameInput.Value())

	// Required field validation
	if name == "" {
		if ic.formType == TaskCreateForm || ic.formType == TaskEditForm {
			return false, "name", "Task name is required"
		} else {
			return false, "name", "Project name is required"
		}
	}

	// Length validation
	if len(name) > 50 {
		return false, "name", "Name too long (max 50 characters)"
	}

	// Project description validation
	if ic.formType == ProjectCreateForm {
		desc := ic.DescInput.Value()
		if len(desc) > 200 {
			return false, "description", "Project description too long (max 200 characters)"
		}
	}

	return true, "", ""
}

// Render renders the form with inline error display
func (ic *InputComponents) Render(errorMsg string, errorField string, width, height int) string {
	title := ic.getFormTitle()
	nameLabel := ic.getNameLabel()
	descLabel := ic.getDescLabel()

	// Render fields with inline error highlighting
	nameField := ic.renderFieldWithError(0, errorMsg, errorField)
	descField := ic.renderFieldWithError(1, errorMsg, errorField)

	instructions := ic.getInstructions()

	formContent := lipgloss.JoinVertical(
		lipgloss.Left,
		"", title, "",
		nameLabel, nameField, "",
		descLabel, descField, "",
		instructions,
	)

	return ic.wrapInBorder(formContent, width, height)
}

func (ic *InputComponents) renderFieldWithError(field int, errorMsg string, errorField string) string {
	var fieldView string
	var isFocused bool
	var fieldName string

	if field == 0 {
		fieldView = ic.NameInput.View()
		isFocused = ic.FocusedField == 0
		fieldName = "name"
	} else {
		fieldView = ic.DescInput.View()
		isFocused = ic.FocusedField == 1
		fieldName = "description"
	}

	// Determine field styling
	borderColor := colors.Overlay1
	if errorMsg != "" && errorField == fieldName {
		borderColor = colors.Red
	} else if isFocused {
		borderColor = colors.Mauve
	}

	// Render field with optional inline error
	fieldWithBorder := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(borderColor)).
		Padding(0, 1).
		Render(fieldView)

	// Add inline error message if this field has error
	if errorMsg != "" && errorField == fieldName {
		errorText := lipgloss.NewStyle().
			Foreground(lipgloss.Color(colors.Red)).
			Faint(true).
			Render("  " + errorMsg)

		return lipgloss.JoinVertical(
			lipgloss.Left,
			fieldWithBorder,
			errorText,
		)
	}

	return fieldWithBorder
}

func (ic *InputComponents) getFormTitle() string {
	switch ic.formType {
	case TaskCreateForm:
		return "Add New Task"
	case TaskEditForm:
		return "Edit Task"
	case ProjectCreateForm:
		return "New Project"

	default:
		return "Form"
	}
}

func (ic *InputComponents) getNameLabel() string {
	label := "Task Name"
	if ic.formType == ProjectCreateForm {
		label = "Project Name"
	}
	return label + ":"
}

func (ic *InputComponents) getDescLabel() string {
	if ic.formType == TaskCreateForm || ic.formType == TaskEditForm {
		return "Task Description:"
	}
	return "Project Description:"
}

func (ic *InputComponents) getInstructions() string {
	switch ic.formType {
	case TaskCreateForm:
		return "Tab: Switch fields • Enter: Create Task • Esc: Cancel"
	case TaskEditForm:
		return "Tab: Switch fields • Enter: Save Changes • Esc: Cancel"
	case ProjectCreateForm:
		return "Tab: Switch fields • Enter: Create Project • Esc: Cancel"

	default:
		return "Tab: Switch fields • Enter: Submit • Esc: Cancel"
	}
}

func (ic *InputComponents) wrapInBorder(formContent string, width, height int) string {
	form := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(colors.Mauve)).
		Padding(2, 3).
		Render(formContent)

	return lipgloss.Place(
		width, height,
		lipgloss.Center, lipgloss.Center,
		form,
	)
}

// FocusName focuses the name input
func (ic *InputComponents) FocusName() {
	ic.FocusedField = 0
	ic.NameInput.Focus()
	ic.DescInput.Blur()
}

// FocusDesc focuses the description input
func (ic *InputComponents) FocusDesc() {
	ic.FocusedField = 1
	ic.DescInput.Focus()
	ic.NameInput.Blur()
}

// BlurName blurs the name input
func (ic *InputComponents) BlurName() {
	ic.NameInput.Blur()
}

// BlurDesc blurs the description input
func (ic *InputComponents) BlurDesc() {
	ic.DescInput.Blur()
}

// Blur blurs all inputs
func (ic *InputComponents) Blur() {
	ic.NameInput.Blur()
	ic.DescInput.Blur()
}

// GetFormType returns the current form type
func (ic *InputComponents) GetFormType() FormType {
	return ic.formType
}

// GetTaskID returns the task ID for edit forms
func (ic *InputComponents) GetTaskID() string {
	return ic.taskID
}
