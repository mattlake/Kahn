package main

import (
	"fmt"
	"strings"

	"kahn/internal/domain"
	"kahn/pkg/input"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	Projects                 []domain.Project
	ActiveProjectID          string
	Tasks                    []list.Model
	activeListIndex          domain.Status
	showForm                 bool
	showProjectSwitch        bool
	showProjectDeleteConfirm bool
	projectToDelete          string
	showTaskDeleteConfirm    bool
	taskToDelete             string
	// New InputComponents system
	taskInputComponents    *input.InputComponents
	projectInputComponents *input.InputComponents
	activeFormType         input.FormType
	formError              string // validation error message
	formErrorField         string // which field has error
	width                  int
	height                 int
	database               *Database
	inputHandler           *input.Handler
	taskService            *TaskService
	projectService         *ProjectService
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m *Model) GetActiveProject() *domain.Project {
	for i, proj := range m.Projects {
		if proj.ID == m.ActiveProjectID {
			return &m.Projects[i]
		}
	}
	return nil
}

// ModelInterface implementation for input handler

func (m *Model) GetActiveProjectID() string {
	return m.ActiveProjectID
}

func (m *Model) CreateTask(name, description string) error {
	activeProj := m.GetActiveProject()
	if activeProj == nil {
		return nil // No active project
	}

	// Use service layer for business logic
	newTask, err := m.taskService.CreateTask(name, description, m.ActiveProjectID)
	if err != nil {
		return err
	}

	// Add to active project in memory
	activeProj.AddTask(*newTask)
	m.updateTaskLists()

	return nil
}

func (m *Model) UpdateTask(id, name, description string) error {
	// Use service layer for business logic
	task, err := m.taskService.UpdateTask(id, name, description)
	if err != nil {
		return err
	}

	// Update task in memory
	activeProj := m.GetActiveProject()
	if activeProj != nil {
		for i, t := range activeProj.Tasks {
			if t.ID == id {
				activeProj.Tasks[i].Name = task.Name
				activeProj.Tasks[i].Desc = task.Desc
				activeProj.Tasks[i].UpdatedAt = task.UpdatedAt
				break
			}
		}
		m.updateTaskLists()
	}

	return nil
}

func (m *Model) DeleteTask(id string) error {
	// Use service layer for business logic
	if err := m.taskService.DeleteTask(id); err != nil {
		return err
	}

	// Remove from active project in memory
	activeProj := m.GetActiveProject()
	if activeProj != nil {
		activeProj.RemoveTask(id)
		m.updateTaskLists()
	}

	return nil
}

func (m *Model) MoveTaskToNextStatus(id string) error {
	activeProj := m.GetActiveProject()
	if activeProj == nil {
		return nil
	}

	// Use service layer for business logic
	task, err := m.taskService.MoveTaskToNextStatus(id)
	if err != nil {
		return err
	}

	activeProj.UpdateTaskStatus(id, task.Status)
	m.updateTaskLists()
	return nil
}

func (m *Model) MoveTaskToPreviousStatus(id string) error {
	activeProj := m.GetActiveProject()
	if activeProj == nil {
		return nil
	}

	// Use service layer for business logic
	task, err := m.taskService.MoveTaskToPreviousStatus(id)
	if err != nil {
		return err
	}

	activeProj.UpdateTaskStatus(id, task.Status)
	m.updateTaskLists()
	return nil
}

func (m *Model) GetSelectedTask() (input.TaskInterface, bool) {
	selectedItem := m.Tasks[m.activeListIndex].SelectedItem()
	if selectedItem == nil {
		return nil, false
	}

	task, ok := selectedItem.(domain.Task)
	if !ok {
		return nil, false
	}

	return &domain.TaskWrapper{Task: task}, true
}

func (m *Model) GetProjects() []input.ProjectInterface {
	var projects []input.ProjectInterface
	for i := range m.Projects {
		projects = append(projects, &domain.ProjectWrapper{Project: &m.Projects[i]})
	}
	return projects
}

func (m *Model) CreateProject(name, description string) error {
	// Use service layer for business logic
	newProject, err := m.projectService.CreateProject(name, description)
	if err != nil {
		return err
	}

	m.Projects = append(m.Projects, *newProject)
	m.ActiveProjectID = newProject.ID
	m.updateTaskLists()

	return nil
}

func (m *Model) DeleteProject(id string) error {
	// Use service layer for business logic
	if err := m.projectService.DeleteProject(id); err != nil {
		return err
	}

	if len(m.Projects) == 1 {
		m.Projects = []domain.Project{}
		m.ActiveProjectID = ""
		m.Tasks[domain.NotStarted].SetItems([]list.Item{})
		m.Tasks[domain.InProgress].SetItems([]list.Item{})
		m.Tasks[domain.Done].SetItems([]list.Item{})
	} else {
		var newProjects []domain.Project
		var wasActiveProject bool
		for _, proj := range m.Projects {
			if proj.ID != id {
				newProjects = append(newProjects, proj)
			} else {
				wasActiveProject = (proj.ID == m.ActiveProjectID)
			}
		}
		m.Projects = newProjects

		if wasActiveProject && len(m.Projects) > 0 {
			m.ActiveProjectID = m.Projects[0].ID
			m.updateTaskLists()
		}
	}

	return nil
}

func (m *Model) SwitchToProject(id string) error {
	m.ActiveProjectID = id
	m.updateTaskLists()
	return nil
}

func (m *Model) GetSelectedProjectIndex() int {
	for i, proj := range m.Projects {
		if proj.ID == m.ActiveProjectID {
			return i
		}
	}
	return 0
}

func (m *Model) setFormError(message string, field string) {
	m.formError = message
	m.formErrorField = field
}

func (m *Model) ClearFormError() {
	m.formError = ""
	m.formErrorField = ""
}

func (m *Model) GetFormError() string {
	return m.formError
}

func (m *Model) GetFormErrorField() string {
	return m.formErrorField
}

func (m *Model) GetActiveInputComponents() *input.InputComponents {
	if m.activeFormType == input.TaskCreateForm || m.activeFormType == input.TaskEditForm {
		return m.taskInputComponents
	}
	return m.projectInputComponents
}

func (m *Model) GetActiveFormType() input.FormType {
	return m.activeFormType
}

func (m *Model) SubmitCurrentForm() error {
	comps := m.GetActiveInputComponents()

	// Validate at submit time
	isValid, errorField, errorMsg := comps.ValidateForSubmit()
	if !isValid {
		m.setFormError(errorMsg, errorField)
		return fmt.Errorf("validation failed: %s", errorMsg)
	}

	// Clear error and proceed with submission
	m.ClearFormError()
	name := strings.TrimSpace(comps.NameInput.Value())
	desc := comps.DescInput.Value()

	switch m.activeFormType {
	case input.TaskCreateForm:
		newTask, err := m.taskService.CreateTask(name, desc, m.ActiveProjectID)
		if err == nil {
			// Add to active project in memory
			activeProj := m.GetActiveProject()
			if activeProj != nil {
				activeProj.AddTask(*newTask)
				m.updateTaskLists()
			}
		}
		return err
	case input.TaskEditForm:
		err := m.UpdateTask(comps.GetTaskID(), name, desc)
		return err
	case input.ProjectCreateForm:
		newProject, err := m.projectService.CreateProject(name, desc)
		if err == nil {
			// Add to projects list and switch to it
			m.Projects = append(m.Projects, *newProject)
			m.ActiveProjectID = newProject.ID
			m.updateTaskLists()
		}
		return err
	}
	return nil
}

func (m *Model) CancelCurrentForm() {
	m.showForm = false
	m.ClearFormError()
	// Reset forms
	m.taskInputComponents.Reset()
	m.projectInputComponents.Reset()
}

func (m *Model) ShowTaskForm() {
	m.taskInputComponents.SetupForTaskCreate()
	m.activeFormType = input.TaskCreateForm
	m.showForm = true
	m.ClearFormError()
}

func (m *Model) ShowTaskEditForm(taskID string, name, description string) {
	m.taskInputComponents.SetupForTaskEdit(taskID, name, description)
	m.activeFormType = input.TaskEditForm
	m.showForm = true
	m.ClearFormError()
}

func (m *Model) ShowProjectForm() {
	m.projectInputComponents.SetupForProjectCreate()
	m.activeFormType = input.ProjectCreateForm
	m.showForm = true
	m.ClearFormError()
}

func (m *Model) ShowProjectSwitcher() {
	m.showProjectSwitch = true
}

func (m *Model) HideAllForms() {
	m.showForm = false
	m.showProjectSwitch = false
	m.showTaskDeleteConfirm = false
	m.showProjectDeleteConfirm = false
	m.ClearFormError()
	m.taskInputComponents.Reset()
	m.projectInputComponents.Reset()
	m.taskToDelete = ""
	m.projectToDelete = ""
}

func (m *Model) NextList() {
	if m.activeListIndex == domain.Done {
		m.activeListIndex = domain.NotStarted
	} else {
		m.activeListIndex++
	}
}

func (m *Model) PrevList() {
	if m.activeListIndex == domain.NotStarted {
		m.activeListIndex = domain.Done
	} else {
		m.activeListIndex--
	}
}

func (m *Model) updateTaskLists() {
	activeProj := m.GetActiveProject()
	if activeProj == nil {
		return
	}

	m.Tasks[domain.NotStarted].SetItems(convertTasksToListItems(activeProj.GetTasksByStatus(domain.NotStarted)))
	m.Tasks[domain.InProgress].SetItems(convertTasksToListItems(activeProj.GetTasksByStatus(domain.InProgress)))
	m.Tasks[domain.Done].SetItems(convertTasksToListItems(activeProj.GetTasksByStatus(domain.Done)))
}

func (m *Model) executeTaskDeletion() *Model {
	if m.taskToDelete == "" {
		m.showTaskDeleteConfirm = false
		return m
	}

	// Use service layer for business logic
	if err := m.taskService.DeleteTask(m.taskToDelete); err != nil {
		// If deletion fails, cancel the operation
		m.showTaskDeleteConfirm = false
		m.taskToDelete = ""
		return m
	}

	// Remove task from active project in memory
	activeProj := m.GetActiveProject()
	if activeProj != nil {
		activeProj.RemoveTask(m.taskToDelete)
		m.updateTaskLists()
	}

	// Reset confirmation state
	m.showTaskDeleteConfirm = false
	m.taskToDelete = ""

	return m
}

func (m *Model) executeProjectDeletion() *Model {
	if m.projectToDelete == "" {
		m.showProjectDeleteConfirm = false
		return m
	}

	// Use service layer for business logic
	if err := m.projectService.DeleteProject(m.projectToDelete); err != nil {
		// If deletion fails, cancel the operation
		m.showProjectDeleteConfirm = false
		m.projectToDelete = ""
		return m
	}

	// Handle edge case: deleting last project
	if len(m.Projects) == 1 {
		m.Projects = []domain.Project{}
		m.ActiveProjectID = ""
		// Clear task lists
		m.Tasks[domain.NotStarted].SetItems([]list.Item{})
		m.Tasks[domain.InProgress].SetItems([]list.Item{})
		m.Tasks[domain.Done].SetItems([]list.Item{})
	} else {
		// Find and remove the project from the slice
		var newProjects []domain.Project
		var wasActiveProject bool
		for _, proj := range m.Projects {
			if proj.ID != m.projectToDelete {
				newProjects = append(newProjects, proj)
			} else {
				wasActiveProject = (proj.ID == m.ActiveProjectID)
			}
		}
		m.Projects = newProjects

		// If we deleted the active project, switch to the next one
		if wasActiveProject && len(m.Projects) > 0 {
			// Switch to the first remaining project
			m.ActiveProjectID = m.Projects[0].ID
			m.updateTaskLists()
		}
	}

	// Reset confirmation state
	m.showProjectDeleteConfirm = false
	m.projectToDelete = ""

	return m
}
