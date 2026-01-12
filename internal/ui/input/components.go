package input

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
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

type InputComponents struct {
	NameInput      textinput.Model
	DescInput      textarea.Model
	PriorityValue  domain.Priority // Track current priority value
	TypeValue      domain.TaskType // Track current task type value
	BlockedByValue *int            // Currently selected blocker (nil = None)
	availableTasks []domain.Task   // Tasks that can block this one
	blockedByIndex int             // Current index in availableTasks (-1 = None)
	formType       FormType
	taskID         string // for edit forms
	FocusedField   int    // 0=name, 1=desc, 2=priority, 3=type, 4=blockedBy (exported)
}

func NewInputComponents() InputComponents {
	ta := textarea.New()
	ta.SetWidth(40) // Must set width/height to initialize internal viewport
	ta.SetHeight(4)
	return InputComponents{
		formType:     TaskCreateForm,
		FocusedField: 0,
		TypeValue:    domain.RegularTask, // Default to RegularTask
		DescInput:    ta,                 // Use initialized textarea
	}
}

func (ic *InputComponents) SetupForTaskCreate() {
	ic.formType = TaskCreateForm
	ic.FocusedField = 0
	ic.taskID = ""
	ic.PriorityValue = domain.Low     // Default to Low priority
	ic.TypeValue = domain.RegularTask // Default to RegularTask type
	ic.BlockedByValue = nil           // Default to no blocker
	ic.blockedByIndex = -1            // -1 represents "None"
	ic.availableTasks = []domain.Task{}
	ic.NameInput = ic.createNameInput("Task name *")
	ic.DescInput = ic.createDescInput("Task description (optional)")
	ic.NameInput.Focus()
}

func (ic *InputComponents) SetupForTaskEdit(taskID, name, desc string, priority domain.Priority, taskType domain.TaskType, blockedByIntID *int) {
	ic.formType = TaskEditForm
	ic.FocusedField = 0
	ic.taskID = taskID
	ic.PriorityValue = priority
	ic.TypeValue = taskType
	ic.BlockedByValue = blockedByIntID
	ic.blockedByIndex = -1 // Will be set by SetAvailableTasks
	ic.availableTasks = []domain.Task{}
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

func (ic *InputComponents) createDescInput(placeholder string) textarea.Model {
	ta := textarea.New()
	ta.Placeholder = placeholder
	ta.CharLimit = 200
	ta.SetWidth(40)
	ta.SetHeight(4) // 4 lines as recommended
	ta.ShowLineNumbers = false

	// Style configuration - focused state
	ta.FocusedStyle.Placeholder = lipgloss.NewStyle().Foreground(lipgloss.Color(colors.Subtext0))
	ta.FocusedStyle.Text = lipgloss.NewStyle().Foreground(lipgloss.Color(colors.Text))
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle() // No cursor line highlight

	// Style configuration - blurred state
	ta.BlurredStyle.Placeholder = lipgloss.NewStyle().Foreground(lipgloss.Color(colors.Subtext0))
	ta.BlurredStyle.Text = lipgloss.NewStyle().Foreground(lipgloss.Color(colors.Text))

	// Cursor style
	ta.Cursor.Style = lipgloss.NewStyle().Foreground(lipgloss.Color(colors.Text))

	return ta
}

func (ic *InputComponents) Reset() {
	ic.NameInput.Reset()
	ic.DescInput.Reset()
	ic.PriorityValue = domain.Low
	ic.TypeValue = domain.RegularTask
	ic.BlockedByValue = nil
	ic.blockedByIndex = -1
	ic.availableTasks = []domain.Task{}
	ic.FocusedField = 0
	ic.taskID = ""
}

func (ic *InputComponents) FocusPriority() {
	ic.FocusedField = 2
}

func (ic *InputComponents) BlurPriority() {
	// No specific blur needed for priority field
}

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

func (ic *InputComponents) CycleTypeUp() {
	switch ic.TypeValue {
	case domain.RegularTask:
		ic.TypeValue = domain.Bug
	case domain.Bug:
		ic.TypeValue = domain.Feature
	case domain.Feature:
		ic.TypeValue = domain.RegularTask
	}
}

func (ic *InputComponents) CycleTypeDown() {
	switch ic.TypeValue {
	case domain.RegularTask:
		ic.TypeValue = domain.Feature
	case domain.Bug:
		ic.TypeValue = domain.RegularTask
	case domain.Feature:
		ic.TypeValue = domain.Bug
	}
}

// SetAvailableTasks sets the list of tasks that can block this task
func (ic *InputComponents) SetAvailableTasks(tasks []domain.Task) {
	ic.availableTasks = tasks

	// Find the index of the currently selected blocker
	if ic.BlockedByValue != nil {
		for i, task := range tasks {
			if task.IntID == *ic.BlockedByValue {
				ic.blockedByIndex = i
				return
			}
		}
		// If blocker not found in available tasks, reset to None
		ic.blockedByIndex = -1
		ic.BlockedByValue = nil
	} else {
		ic.blockedByIndex = -1
	}
}

// CycleBlockedByUp cycles to the next blocking task (up in list)
func (ic *InputComponents) CycleBlockedByUp() {
	if len(ic.availableTasks) == 0 {
		// No tasks available, stay at None
		ic.blockedByIndex = -1
		ic.BlockedByValue = nil
		return
	}

	ic.blockedByIndex++
	if ic.blockedByIndex >= len(ic.availableTasks) {
		// Wrap to None
		ic.blockedByIndex = -1
		ic.BlockedByValue = nil
	} else {
		// Set to the task at current index
		intID := ic.availableTasks[ic.blockedByIndex].IntID
		ic.BlockedByValue = &intID
	}
}

// CycleBlockedByDown cycles to the previous blocking task (down in list)
func (ic *InputComponents) CycleBlockedByDown() {
	if len(ic.availableTasks) == 0 {
		// No tasks available, stay at None
		ic.blockedByIndex = -1
		ic.BlockedByValue = nil
		return
	}

	ic.blockedByIndex--
	if ic.blockedByIndex < -1 {
		// Wrap to last task
		ic.blockedByIndex = len(ic.availableTasks) - 1
		intID := ic.availableTasks[ic.blockedByIndex].IntID
		ic.BlockedByValue = &intID
	} else if ic.blockedByIndex == -1 {
		// Set to None
		ic.BlockedByValue = nil
	} else {
		// Set to the task at current index
		intID := ic.availableTasks[ic.blockedByIndex].IntID
		ic.BlockedByValue = &intID
	}
}

// FocusBlockedBy focuses the blocked by field
func (ic *InputComponents) FocusBlockedBy() {
	ic.FocusedField = 4
}

// BlurBlockedBy blurs the blocked by field
func (ic *InputComponents) BlurBlockedBy() {
	// No specific blur needed for blocked by field
}

func (ic *InputComponents) IsTaskForm() bool {
	return ic.formType == TaskCreateForm || ic.formType == TaskEditForm
}

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

func (ic *InputComponents) Render(errorMsg string, errorField string, width, height int) string {
	title := ic.getFormTitle()
	nameLabel := ic.getNameLabel()
	descLabel := ic.getDescLabel()

	// Render fields with inline error highlighting
	nameField := ic.renderFieldWithError(0, errorMsg, errorField)
	descField := ic.renderFieldWithError(1, errorMsg, errorField)
	var priorityField string
	var typeField string
	var blockedByField string

	// Only show priority, type, and blocked by fields for task forms
	if ic.formType == TaskCreateForm || ic.formType == TaskEditForm {
		priorityField = ic.renderPriorityField(errorMsg, errorField)
		typeField = ic.renderTypeField(errorMsg, errorField)
		blockedByField = ic.renderBlockedByField(errorMsg, errorField)
	}

	instructions := ic.getInstructions()

	// Build form content
	var formContent string
	if ic.formType == TaskCreateForm || ic.formType == TaskEditForm {
		// Task forms include priority, type, and blocked by fields
		formContent = lipgloss.JoinVertical(
			lipgloss.Left,
			"", title, "",
			nameLabel, nameField, "",
			descLabel, descField, "",
			"Priority:", priorityField, "",
			"Type:", typeField, "",
			"Blocked By:", blockedByField, "",
			instructions,
		)
	} else {
		// Project forms (no priority, type, or blocked by fields)
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

// renderTypeField renders the task type field with visual indicators
func (ic *InputComponents) renderTypeField(errorMsg string, errorField string) string {
	typeOptions := []string{"Task", "Bug", "Feature"}
	currentIndex := int(ic.TypeValue)

	// Build display string like "Type: Bug »"
	display := fmt.Sprintf("Type: %s", typeOptions[currentIndex])

	if ic.FocusedField == 3 { // Type field focused
		display += " »" // Add indicator for focused field
	}

	// Determine field styling
	borderColor := colors.Text
	if errorMsg != "" && errorField == "type" {
		borderColor = colors.Red
	} else if ic.FocusedField == 3 {
		borderColor = colors.Green
	}

	// Get type color for text
	var textColor string
	switch ic.TypeValue {
	case domain.RegularTask:
		textColor = colors.Text // Default color
	case domain.Bug:
		textColor = colors.Text // Use same color as other fields
	case domain.Feature:
		textColor = colors.Text // Use same color as other fields
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

// renderBlockedByField renders the blocked by field with visual indicators
func (ic *InputComponents) renderBlockedByField(errorMsg string, errorField string) string {
	var display string

	if ic.BlockedByValue == nil {
		display = "Blocked By: (None)"
	} else {
		// Find the task name
		taskName := "(Unknown)"
		for _, task := range ic.availableTasks {
			if task.IntID == *ic.BlockedByValue {
				taskName = task.Name
				break
			}
		}
		display = fmt.Sprintf("Blocked By: %s", taskName)
	}

	if ic.FocusedField == 4 { // BlockedBy field focused
		display += " »" // Add indicator for focused field
	}

	// Determine field styling
	borderColor := colors.Text
	if errorMsg != "" && errorField == "blocked_by" {
		borderColor = colors.Red
	} else if ic.FocusedField == 4 {
		borderColor = colors.Green
	}

	// Render field with border
	fieldWithBorder := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(borderColor)).
		Padding(0, 1).
		Width(40).
		Render(lipgloss.NewStyle().Foreground(lipgloss.Color(colors.Text)).Render(display))

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
		return "Tab: Switch fields • ↑/↓: Change selection • Enter/Ctrl+Enter: Create Task • Esc: Cancel"
	case TaskEditForm:
		return "Tab: Switch fields • ↑/↓: Change selection • Enter/Ctrl+Enter: Save Changes • Esc: Cancel"
	case ProjectCreateForm:
		return "Tab: Switch fields • Enter/Ctrl+Enter: Create Project • Esc: Cancel"

	default:
		return "Tab: Switch fields • Enter/Ctrl+Enter: Submit • Esc: Cancel"
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

// FocusType focuses the type field
func (ic *InputComponents) FocusType() {
	ic.FocusedField = 3
}

// BlurType removes focus from the type field
func (ic *InputComponents) BlurType() {
	// No specific blur needed for type field
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
