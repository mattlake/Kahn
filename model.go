package main

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"kahn/pkg/input"
)

type Model struct {
	Projects                 []Project
	ActiveProjectID          string
	Tasks                    []list.Model
	activeListIndex          Status
	showForm                 bool
	showProjectSwitch        bool
	showProjectForm          bool
	showProjectDeleteConfirm bool
	projectToDelete          string
	showTaskDeleteConfirm    bool
	taskToDelete             string
	showTaskEditForm         bool
	editingTaskID            string
	nameInput                textinput.Model
	descInput                textinput.Model
	projNameInput            textinput.Model
	projDescInput            textinput.Model
	focusedInput             int // 0 for name, 1 for desc
	focusedProjInput         int // 0 for name, 1 for desc
	width                    int
	height                   int
	database                 *Database
	projectDAO               *ProjectDAO
	taskDAO                  *TaskDAO
	inputHandler             *input.Handler
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m *Model) GetActiveProject() *Project {
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

	newTask := NewTask(name, description, m.ActiveProjectID)

	// Save to database
	if err := m.taskDAO.Create(newTask); err != nil {
		return err
	}

	// Add to active project in memory
	activeProj.AddTask(*newTask)
	m.updateTaskLists()

	return nil
}

func (m *Model) UpdateTask(id, name, description string) error {
	task, err := m.taskDAO.GetByID(id)
	if err != nil {
		return err
	}

	task.Name = name
	task.Desc = description

	if err := m.taskDAO.Update(task); err != nil {
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
	// Delete from database
	if err := m.taskDAO.Delete(id); err != nil {
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

	var nextStatus Status
	for _, task := range activeProj.Tasks {
		if task.ID == id {
			switch task.Status {
			case NotStarted:
				nextStatus = InProgress
			case InProgress:
				nextStatus = Done
			case Done:
				nextStatus = NotStarted
			}
			break
		}
	}

	if err := m.taskDAO.UpdateStatus(id, nextStatus); err != nil {
		return err
	}

	activeProj.UpdateTaskStatus(id, nextStatus)
	m.updateTaskLists()
	return nil
}

func (m *Model) MoveTaskToPreviousStatus(id string) error {
	activeProj := m.GetActiveProject()
	if activeProj == nil {
		return nil
	}

	var prevStatus Status
	for _, task := range activeProj.Tasks {
		if task.ID == id {
			switch task.Status {
			case NotStarted:
				prevStatus = Done
			case InProgress:
				prevStatus = NotStarted
			case Done:
				prevStatus = InProgress
			}
			break
		}
	}

	if err := m.taskDAO.UpdateStatus(id, prevStatus); err != nil {
		return err
	}

	activeProj.UpdateTaskStatus(id, prevStatus)
	m.updateTaskLists()
	return nil
}

func (m *Model) GetSelectedTask() (input.TaskInterface, bool) {
	selectedItem := m.Tasks[m.activeListIndex].SelectedItem()
	if selectedItem == nil {
		return nil, false
	}

	task, ok := selectedItem.(Task)
	if !ok {
		return nil, false
	}

	return &TaskWrapper{task: task}, true
}

func (m *Model) GetProjects() []input.ProjectInterface {
	var projects []input.ProjectInterface
	for i := range m.Projects {
		projects = append(projects, &ProjectWrapper{project: &m.Projects[i]})
	}
	return projects
}

func (m *Model) CreateProject(name, description string) error {
	newProject := NewProject(name, description, "#89b4fa") // Blue color

	if err := newProject.Validate(); err != nil {
		return err
	}

	if err := m.projectDAO.Create(newProject); err != nil {
		return err
	}

	m.Projects = append(m.Projects, *newProject)
	m.ActiveProjectID = newProject.ID
	m.updateTaskLists()

	return nil
}

func (m *Model) DeleteProject(id string) error {
	if err := m.projectDAO.Delete(id); err != nil {
		return err
	}

	if len(m.Projects) == 1 {
		m.Projects = []Project{}
		m.ActiveProjectID = ""
		m.Tasks[NotStarted].SetItems([]list.Item{})
		m.Tasks[InProgress].SetItems([]list.Item{})
		m.Tasks[Done].SetItems([]list.Item{})
	} else {
		var newProjects []Project
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

func (m *Model) ShowTaskForm() {
	m.showForm = true
	m.focusedInput = 0
	m.nameInput.Focus()
	m.descInput.Blur()
}

func (m *Model) ShowTaskEditForm(taskID string, name, description string) {
	m.showTaskEditForm = true
	m.editingTaskID = taskID
	m.nameInput.SetValue(name)
	m.descInput.SetValue(description)
	m.focusedInput = 0
	m.nameInput.Focus()
	m.descInput.Blur()
}

func (m *Model) ShowProjectForm() {
	m.showProjectForm = true
	m.focusedProjInput = 0
	m.projNameInput.Focus()
	m.projDescInput.Blur()
}

func (m *Model) ShowProjectSwitcher() {
	m.showProjectSwitch = true
}

func (m *Model) HideAllForms() {
	m.showForm = false
	m.showTaskEditForm = false
	m.showProjectForm = false
	m.showProjectSwitch = false
	m.showTaskDeleteConfirm = false
	m.showProjectDeleteConfirm = false
	m.nameInput.Reset()
	m.descInput.Reset()
	m.projNameInput.Reset()
	m.projDescInput.Reset()
	m.editingTaskID = ""
	m.taskToDelete = ""
	m.projectToDelete = ""
}

func (m *Model) NextList() {
	if m.activeListIndex == Done {
		m.activeListIndex = NotStarted
	} else {
		m.activeListIndex++
	}
}

func (m *Model) PrevList() {
	if m.activeListIndex == NotStarted {
		m.activeListIndex = Done
	} else {
		m.activeListIndex--
	}
}

func (m *Model) updateTaskLists() {
	activeProj := m.GetActiveProject()
	if activeProj == nil {
		return
	}

	m.Tasks[NotStarted].SetItems(convertTasksToListItems(activeProj.GetTasksByStatus(NotStarted)))
	m.Tasks[InProgress].SetItems(convertTasksToListItems(activeProj.GetTasksByStatus(InProgress)))
	m.Tasks[Done].SetItems(convertTasksToListItems(activeProj.GetTasksByStatus(Done)))
}

func (m *Model) executeTaskDeletion() *Model {
	if m.taskToDelete == "" {
		m.showTaskDeleteConfirm = false
		return m
	}

	// Delete task from database
	if err := m.taskDAO.Delete(m.taskToDelete); err != nil {
		// If database deletion fails, cancel the operation
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

	// Delete project from database
	if err := m.projectDAO.Delete(m.projectToDelete); err != nil {
		// If database deletion fails, cancel the operation
		m.showProjectDeleteConfirm = false
		m.projectToDelete = ""
		return m
	}

	// Handle edge case: deleting last project
	if len(m.Projects) == 1 {
		m.Projects = []Project{}
		m.ActiveProjectID = ""
		// Clear task lists
		m.Tasks[NotStarted].SetItems([]list.Item{})
		m.Tasks[InProgress].SetItems([]list.Item{})
		m.Tasks[Done].SetItems([]list.Item{})
	} else {
		// Find and remove the project from the slice
		var newProjects []Project
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
