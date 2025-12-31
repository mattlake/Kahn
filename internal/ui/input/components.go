package input

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	"kahn/internal/domain"
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
	NameInput     textinput.Model
	DescInput     textinput.Model
	PriorityValue domain.Priority // Track current priority value
	formType      FormType
	taskID        string // for edit forms
	FocusedField  int    // 0=name, 1=desc, 2=priority (exported)
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
	ic.PriorityValue = domain.Low // Default to Low priority
	ic.NameInput = ic.createNameInput("Task name *")
	ic.DescInput = ic.createDescInput("Task description (optional)")
	ic.NameInput.Focus()
}

func (ic *InputComponents) SetupForTaskEdit(taskID, name, desc string, priority domain.Priority) {
	ic.formType = TaskEditForm
	ic.FocusedField = 0
	ic.taskID = taskID
	ic.PriorityValue = priority
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
	input.Cursor.Style = lipgloss.NewStyle().Foreground(lipgloss.Color(colors.Text))
	input.CharLimit = 50
	input.Width = 40
	return input
}

func (ic *InputComponents) createDescInput(placeholder string) textinput.Model {
	input := textinput.New()
	input.Placeholder = placeholder
	input.PlaceholderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(colors.Subtext0))
	input.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(colors.Text))
	input.Cursor.Style = lipgloss.NewStyle().Foreground(lipgloss.Color(colors.Text))
	input.CharLimit = 200 // Larger for project descriptions
	input.Width = 40
	return input
}

// Reset resets all input components
func (ic *InputComponents) Reset() {
	ic.NameInput.Reset()
	ic.DescInput.Reset()
	ic.PriorityValue = domain.Low
	ic.FocusedField = 0
	ic.taskID = ""
}

// FocusPriority sets focus to priority field
func (ic *InputComponents) FocusPriority() {
	ic.FocusedField = 2
}

// BlurPriority removes focus from priority field
func (ic *InputComponents) BlurPriority() {
	// No specific blur needed for priority field
}

// CyclePriorityUp cycles priority up (Low -> Medium -> High -> Low)
func (ic *InputComponents) CyclePriorityUp() {
	switch ic.PriorityValue {
	case domain.Low:
		ic.PriorityValue = domain.Medium
	case domain.Medium:
		ic.PriorityValue = domain.High
	case domain.High:
		ic.PriorityValue = domain.Low
	}
}

// CyclePriorityDown cycles priority down (High -> Medium -> Low -> High)
func (ic *InputComponents) CyclePriorityDown() {
	switch ic.PriorityValue {
	case domain.Low:
		ic.PriorityValue = domain.High
	case domain.Medium:
		ic.PriorityValue = domain.Low
	case domain.High:
		ic.PriorityValue = domain.Medium
	}
}

// IsTaskForm returns true if the current form is a task form
func (ic *InputComponents) IsTaskForm() bool {
	return ic.formType == TaskCreateForm || ic.formType == TaskEditForm
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
	var priorityField string

	// Only show priority field for task forms
	if ic.formType == TaskCreateForm || ic.formType == TaskEditForm {
		priorityField = ic.renderPriorityField(errorMsg, errorField)
	}

	instructions := ic.getInstructions()

	// Build form content
	var formContent string
	if ic.formType == TaskCreateForm || ic.formType == TaskEditForm {
		// Task forms include priority field
		formContent = lipgloss.JoinVertical(
			lipgloss.Left,
			"", title, "",
			nameLabel, nameField, "",
			descLabel, descField, "",
			"Priority:", priorityField, "",
			instructions,
		)
	} else {
		// Project forms (no priority field)
		formContent = lipgloss.JoinVertical(
			lipgloss.Left,
			"", title, "",
			nameLabel, nameField, "",
			descLabel, descField, "",
			instructions,
		)
	}

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
	borderColor := colors.Text
	if errorMsg != "" && errorField == fieldName {
		borderColor = colors.Red
	} else if isFocused {
		borderColor = colors.Green
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

// renderPriorityField renders the priority field with visual indicators
func (ic *InputComponents) renderPriorityField(errorMsg string, errorField string) string {
	priorityOptions := []string{"Low", "Medium", "High"}
	currentIndex := int(ic.PriorityValue)

	// Build display string like "Priority: Low »"
	display := fmt.Sprintf("Priority: %s", priorityOptions[currentIndex])

	if ic.FocusedField == 2 { // Priority field focused
		display += " »" // Add indicator for focused field
	}

	// Determine field styling
	borderColor := colors.Text
	if errorMsg != "" && errorField == "priority" {
		borderColor = colors.Red
	} else if ic.FocusedField == 2 {
		borderColor = colors.Green
	}

	// Get priority color for the text
	var textColor string
	switch ic.PriorityValue {
	case domain.Low:
		textColor = colors.Green
	case domain.Medium:
		textColor = colors.Peach
	case domain.High:
		textColor = colors.Red
	}

	// Render field with border and color
	fieldWithBorder := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(borderColor)).
		Padding(0, 1).
		Width(40).
		Render(lipgloss.NewStyle().Foreground(lipgloss.Color(textColor)).Render(display))

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
		BorderForeground(lipgloss.Color(colors.Green)).
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
