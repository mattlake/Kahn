package input

import (
	tea "github.com/charmbracelet/bubbletea"
)

// InputMode represents the current input mode
type InputMode int

const (
	NormalMode InputMode = iota
	TaskFormMode
	ProjectFormMode
	ProjectSwitchMode
	TaskEditFormMode
	TaskDeleteConfirmMode
	ProjectDeleteConfirmMode
)

// FocusType represents which input is focused
type FocusType int

const (
	NameFocus FocusType = iota
	DescriptionFocus
)

// ActionResult represents the result of an input action
type ActionResult struct {
	Mode         InputMode
	FocusType    FocusType
	Cmd          tea.Cmd
	Handled      bool
	ShouldUpdate bool
	ExitMode     bool
}

// ModelInterface defines the interface for the model that the input handler works with
type ModelInterface interface {
	// Task operations
	GetActiveProjectID() string
	CreateTask(name, description string) error
	UpdateTask(id, name, description string) error
	DeleteTask(id string) error
	MoveTaskToNextStatus(id string) error
	MoveTaskToPreviousStatus(id string) error
	GetSelectedTask() (TaskInterface, bool)

	// Project operations
	GetProjects() []ProjectInterface
	CreateProject(name, description string) error
	DeleteProject(id string) error
	SwitchToProject(id string) error
	GetSelectedProjectIndex() int

	// UI state operations
	ShowTaskForm()
	ShowTaskEditForm(taskID string, name, description string)
	ShowProjectForm()
	ShowProjectSwitcher()
	HideAllForms()

	// List operations
	NextList()
	PrevList()

	// Form operations (new)
	GetActiveInputComponents() *InputComponents
	GetActiveFormType() FormType
	SubmitCurrentForm() error
	CancelCurrentForm()
	GetFormError() string
	GetFormErrorField() string
	ClearFormError()
}

// TaskInterface defines the interface for a task
type TaskInterface interface {
	GetID() string
	GetName() string
	GetDescription() string
}

// ProjectInterface defines the interface for a project
type ProjectInterface interface {
	GetID() string
	GetName() string
}

// Handler manages input handling for the application
type Handler struct {
	mode      InputMode
	focusType FocusType
}

// NewHandler creates a new input handler
func NewHandler() *Handler {
	return &Handler{
		mode:      NormalMode,
		focusType: NameFocus,
	}
}

// HandleKeyMsg handles key messages and returns actions
func (h *Handler) HandleKeyMsg(msg tea.KeyMsg, model ModelInterface) ActionResult {
	// Route to mode-specific handlers
	switch h.mode {
	case TaskFormMode:
		return h.handleTaskFormKeys(msg, model)
	case TaskEditFormMode:
		return h.handleTaskEditFormKeys(msg, model)
	case ProjectFormMode:
		return h.handleProjectFormKeys(msg, model)
	case ProjectSwitchMode:
		return h.handleProjectSwitchKeys(msg, model)
	case TaskDeleteConfirmMode:
		return h.handleTaskDeleteConfirmKeys(msg, model)
	case ProjectDeleteConfirmMode:
		return h.handleProjectDeleteConfirmKeys(msg, model)
	default:
		return h.handleNormalModeKeys(msg, model)
	}
}

// handleNormalModeKeys handles keys when in normal mode
func (h *Handler) handleNormalModeKeys(msg tea.KeyMsg, model ModelInterface) ActionResult {
	switch msg.String() {
	case "l":
		model.NextList()
		return ActionResult{Handled: true}
	case "h":
		model.PrevList()
		return ActionResult{Handled: true}
	case "n":
		model.ShowTaskForm()
		h.mode = TaskFormMode
		h.focusType = NameFocus
		return ActionResult{Handled: true, Mode: TaskFormMode, FocusType: NameFocus}
	case "p":
		model.ShowProjectSwitcher()
		h.mode = ProjectSwitchMode
		return ActionResult{Handled: true, Mode: ProjectSwitchMode}
	case "e":
		if task, ok := model.GetSelectedTask(); ok {
			model.ShowTaskEditForm(task.GetID(), task.GetName(), task.GetDescription())
			h.mode = TaskEditFormMode
			h.focusType = NameFocus
			return ActionResult{Handled: true, Mode: TaskEditFormMode, FocusType: NameFocus}
		}
		return ActionResult{Handled: true}
	case "d":
		if _, ok := model.GetSelectedTask(); ok {
			h.mode = TaskDeleteConfirmMode
			return ActionResult{Handled: true, Mode: TaskDeleteConfirmMode}
		}
		return ActionResult{Handled: true}
	case "enter":
		// Let main Update function handle task progression
		return ActionResult{Handled: false}
	case "backspace":
		// Let main Update function handle task progression
		return ActionResult{Handled: false}
	}

	return ActionResult{Handled: false}
}

// handleTaskFormKeys handles keys when in task form mode
func (h *Handler) handleTaskFormKeys(msg tea.KeyMsg, model ModelInterface) ActionResult {
	switch msg.String() {
	case "esc":
		model.HideAllForms()
		h.mode = NormalMode
		return ActionResult{Handled: true, Mode: NormalMode, ExitMode: true}
	case "tab":
		return h.handleTabKey(model)
	case "enter":
		// Enter always means submit - no field advancement
		if err := model.SubmitCurrentForm(); err != nil {
			// Validation failed - stay in form mode, show inline error
			return ActionResult{Handled: true, ShouldUpdate: true, Mode: h.mode}
		}

		// Success - exit form mode
		model.CancelCurrentForm()
		h.mode = NormalMode
		return ActionResult{Handled: true, Mode: NormalMode, ExitMode: true}
	default:
		// Clear any previous errors when user types
		model.ClearFormError()
		return ActionResult{Handled: false} // Let textinput handle it
	}
}

// handleTaskEditFormKeys handles keys when in task edit form mode
func (h *Handler) handleTaskEditFormKeys(msg tea.KeyMsg, model ModelInterface) ActionResult {
	// Use the same logic as task form - consolidated handling
	return h.handleTaskFormKeys(msg, model)
}

func (h *Handler) handleTabKey(model ModelInterface) ActionResult {
	comps := model.GetActiveInputComponents()
	if comps.FocusedField == 0 {
		comps.FocusDesc()
		comps.BlurName()
		return ActionResult{Handled: true, Mode: h.mode, FocusType: DescriptionFocus}
	}
	comps.FocusName()
	comps.BlurDesc()
	return ActionResult{Handled: true, Mode: h.mode, FocusType: NameFocus}
}

// handleProjectFormKeys handles keys when in project form mode
func (h *Handler) handleProjectFormKeys(msg tea.KeyMsg, model ModelInterface) ActionResult {
	// Same logic as task form - consolidated handling
	return h.handleTaskFormKeys(msg, model)
}

// handleProjectSwitchKeys handles keys when in project switcher mode
func (h *Handler) handleProjectSwitchKeys(msg tea.KeyMsg, model ModelInterface) ActionResult {
	switch msg.String() {
	case "esc":
		model.HideAllForms()
		h.mode = NormalMode
		return ActionResult{Handled: true, Mode: NormalMode, ExitMode: true}
	case "n":
		model.HideAllForms()
		model.ShowProjectForm()
		h.mode = ProjectFormMode
		h.focusType = NameFocus
		return ActionResult{Handled: true, Mode: ProjectFormMode, FocusType: NameFocus}
	case "d":
		h.mode = ProjectDeleteConfirmMode
		return ActionResult{Handled: true, Mode: ProjectDeleteConfirmMode}
	case "j":
		projects := model.GetProjects()
		if len(projects) > 0 {
			currentIndex := model.GetSelectedProjectIndex()
			nextIndex := (currentIndex + 1) % len(projects)
			if err := model.SwitchToProject(projects[nextIndex].GetID()); err == nil {
				return ActionResult{Handled: true, ShouldUpdate: true}
			}
		}
		return ActionResult{Handled: true}
	case "k":
		projects := model.GetProjects()
		if len(projects) > 0 {
			currentIndex := model.GetSelectedProjectIndex()
			prevIndex := (currentIndex - 1 + len(projects)) % len(projects)
			if err := model.SwitchToProject(projects[prevIndex].GetID()); err == nil {
				return ActionResult{Handled: true, ShouldUpdate: true}
			}
		}
		return ActionResult{Handled: true}
	case "enter":
		model.HideAllForms()
		h.mode = NormalMode
		return ActionResult{Handled: true, Mode: NormalMode, ExitMode: true}
	default:
		// Handle number keys for project selection
		if len(msg.String()) == 1 && msg.String()[0] >= '1' && msg.String()[0] <= '9' {
			projects := model.GetProjects()
			index := int(msg.String()[0] - '1')
			if index < len(projects) {
				if err := model.SwitchToProject(projects[index].GetID()); err == nil {
					model.HideAllForms()
					h.mode = NormalMode
					return ActionResult{Handled: true, Mode: NormalMode, ExitMode: true, ShouldUpdate: true}
				}
			}
		}
		return ActionResult{Handled: true}
	}
}

// handleTaskDeleteConfirmKeys handles keys when in task delete confirmation mode
func (h *Handler) handleTaskDeleteConfirmKeys(msg tea.KeyMsg, model ModelInterface) ActionResult {
	switch msg.String() {
	case "y", "Y":
		if task, ok := model.GetSelectedTask(); ok {
			if err := model.DeleteTask(task.GetID()); err == nil {
				model.HideAllForms()
				h.mode = NormalMode
				return ActionResult{Handled: true, Mode: NormalMode, ExitMode: true, ShouldUpdate: true}
			}
		}
		return ActionResult{Handled: true}
	case "n", "N", "esc":
		model.HideAllForms()
		h.mode = NormalMode
		return ActionResult{Handled: true, Mode: NormalMode, ExitMode: true}
	}
	return ActionResult{Handled: false}
}

// handleProjectDeleteConfirmKeys handles keys when in project delete confirmation mode
func (h *Handler) handleProjectDeleteConfirmKeys(msg tea.KeyMsg, model ModelInterface) ActionResult {
	switch msg.String() {
	case "y", "Y":
		projects := model.GetProjects()
		if len(projects) > 0 {
			currentIndex := model.GetSelectedProjectIndex()
			if currentIndex < len(projects) {
				if err := model.DeleteProject(projects[currentIndex].GetID()); err == nil {
					model.HideAllForms()
					h.mode = NormalMode
					return ActionResult{Handled: true, Mode: NormalMode, ExitMode: true, ShouldUpdate: true}
				}
			}
		}
		return ActionResult{Handled: true}
	case "n", "N", "esc":
		model.HideAllForms()
		h.mode = ProjectSwitchMode
		return ActionResult{Handled: true, Mode: ProjectSwitchMode, ExitMode: true}
	}
	return ActionResult{Handled: false}
}

// SetMode sets the current input mode
func (h *Handler) SetMode(mode InputMode) {
	h.mode = mode
}

// GetMode returns the current input mode
func (h *Handler) GetMode() InputMode {
	return h.mode
}

// SetFocusType sets the current focus type
func (h *Handler) SetFocusType(focusType FocusType) {
	h.focusType = focusType
}

// GetFocusType returns the current focus type
func (h *Handler) GetFocusType() FocusType {
	return h.focusType
}
