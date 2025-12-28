package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"kahn/pkg/input"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle form modes first
		if m.showForm {
			// Ensure input handler mode is synchronized with UI state
			if m.inputHandler.GetMode() != input.TaskFormMode {
				m.inputHandler.SetMode(input.TaskFormMode)
			}

			// Task form mode key handling
			switch msg.String() {
			case "esc":
				m.showForm = false
				m.nameInput.Reset()
				m.descInput.Reset()
				return m, nil
			case "tab":
				if m.focusedInput == 0 {
					m.focusedInput = 1
					m.nameInput.Blur()
					m.descInput.Focus()
				} else {
					m.focusedInput = 0
					m.descInput.Blur()
					m.nameInput.Focus()
				}
				return m, nil
			case "enter":
				if m.nameInput.Value() != "" {
					// Use service layer for task creation
					newTask, err := m.taskService.CreateTask(m.nameInput.Value(), m.descInput.Value(), m.ActiveProjectID)
					if err != nil {
						// If service layer save fails, don't add to model
						return m, nil
					}

					// Add to active project in memory
					activeProj := m.GetActiveProject()
					if activeProj != nil {
						activeProj.AddTask(*newTask)
						m.updateTaskLists()
					}

					// Reset form
					m.showForm = false
					m.nameInput.Reset()
					m.descInput.Reset()
					m.focusedInput = 0
					m.nameInput.Focus()
				}
				return m, nil
			default:
				// Update the focused input
				if m.focusedInput == 0 {
					m.nameInput, _ = m.nameInput.Update(msg)
				} else {
					m.descInput, _ = m.descInput.Update(msg)
				}
				return m, nil
			}
		} else if m.showTaskEditForm {
			// Ensure input handler mode is synchronized with UI state
			if m.inputHandler.GetMode() != input.TaskEditFormMode {
				m.inputHandler.SetMode(input.TaskEditFormMode)
			}

			// Task edit form mode key handling
			switch msg.String() {
			case "esc":
				m.showTaskEditForm = false
				m.editingTaskID = ""
				m.nameInput.Reset()
				m.descInput.Reset()
				m.focusedInput = 0
				m.nameInput.Focus()
				return m, nil
			case "tab":
				if m.focusedInput == 0 {
					m.focusedInput = 1
					m.nameInput.Blur()
					m.descInput.Focus()
				} else {
					m.focusedInput = 0
					m.descInput.Blur()
					m.nameInput.Focus()
				}
				return m, nil
			case "enter":
				if m.nameInput.Value() != "" {
					// Update existing task
					if err := m.UpdateTask(m.editingTaskID, m.nameInput.Value(), m.descInput.Value()); err != nil {
						// If database update fails, don't update model
						return m, nil
					}

					// Reset form
					m.showTaskEditForm = false
					m.editingTaskID = ""
					m.nameInput.Reset()
					m.descInput.Reset()
					m.focusedInput = 0
					m.nameInput.Focus()
				}
				return m, nil
			default:
				// Update the focused input
				if m.focusedInput == 0 {
					m.nameInput, _ = m.nameInput.Update(msg)
				} else {
					m.descInput, _ = m.descInput.Update(msg)
				}
				return m, nil
			}
		} else if m.showProjectSwitch || m.showProjectDeleteConfirm {
			// Ensure input handler mode is synchronized with UI state
			if m.showProjectSwitch && m.inputHandler.GetMode() != input.ProjectSwitchMode {
				m.inputHandler.SetMode(input.ProjectSwitchMode)
			}
			if m.showProjectDeleteConfirm && m.inputHandler.GetMode() != input.ProjectDeleteConfirmMode {
				m.inputHandler.SetMode(input.ProjectDeleteConfirmMode)
			}

			// Handle confirmation dialog if active
			if m.showProjectDeleteConfirm {
				switch msg.String() {
				case "y", "Y":
					return m.executeProjectDeletion(), nil
				case "n", "N", "esc":
					m.showProjectDeleteConfirm = false
					m.projectToDelete = ""
					return m, nil
				}
				return m, nil
			}

			// Project switcher mode key handling
			switch msg.String() {
			case "esc":
				m.showProjectSwitch = false
				return m, nil
			case "d":
				// Trigger project deletion confirmation
				if len(m.Projects) > 0 {
					m.showProjectDeleteConfirm = true
					m.projectToDelete = m.ActiveProjectID
				}
				return m, nil
			case "n":
				m.showProjectSwitch = false
				m.showProjectForm = true
				m.focusedProjInput = 0
				m.projNameInput.Focus()
				m.projDescInput.Blur()
				return m, nil
			case "j":
				// Move down in project list
				for i, proj := range m.Projects {
					if proj.ID == m.ActiveProjectID {
						nextIndex := (i + 1) % len(m.Projects)
						m.ActiveProjectID = m.Projects[nextIndex].ID
						m.updateTaskLists()
						return m, nil
					}
				}
			case "k":
				// Move up in project list
				for i, proj := range m.Projects {
					if proj.ID == m.ActiveProjectID {
						prevIndex := (i - 1 + len(m.Projects)) % len(m.Projects)
						m.ActiveProjectID = m.Projects[prevIndex].ID
						m.updateTaskLists()
						return m, nil
					}
				}
			case "enter":
				// Select current project and close switcher
				m.showProjectSwitch = false
				return m, nil
			default:
				// Handle number keys for project selection
				if len(msg.String()) == 1 && msg.String()[0] >= '1' && msg.String()[0] <= '9' {
					index := int(msg.String()[0] - '1')
					if index < len(m.Projects) {
						m.ActiveProjectID = m.Projects[index].ID
						m.updateTaskLists()
						m.showProjectSwitch = false
					}
				}
				return m, nil
			}
		} else if m.showProjectForm {
			// Ensure input handler mode is synchronized with UI state
			if m.inputHandler.GetMode() != input.ProjectFormMode {
				m.inputHandler.SetMode(input.ProjectFormMode)
			}

			// Project form mode key handling
			switch msg.String() {
			case "esc":
				m.showProjectForm = false
				m.projNameInput.Reset()
				m.projDescInput.Reset()
				return m, nil
			case "tab":
				if m.focusedProjInput == 0 {
					m.focusedProjInput = 1
					m.projNameInput.Blur()
					m.projDescInput.Focus()
				} else {
					m.focusedProjInput = 0
					m.projDescInput.Blur()
					m.projNameInput.Focus()
				}
				return m, nil
			case "enter":
				if m.projNameInput.Value() != "" {
					// Use service layer for project creation
					newProject, err := m.projectService.CreateProject(m.projNameInput.Value(), m.projDescInput.Value())
					if err != nil {
						// If service layer save fails, don't add to model
						return m, nil
					}

					m.Projects = append(m.Projects, *newProject)
					m.ActiveProjectID = newProject.ID
					m.updateTaskLists()

					// Reset form
					m.showProjectForm = false
					m.projNameInput.Reset()
					m.projDescInput.Reset()
					m.focusedProjInput = 0
					m.projNameInput.Focus()
				}
				return m, nil
			default:
				// Update the focused input
				if m.focusedProjInput == 0 {
					m.projNameInput, _ = m.projNameInput.Update(msg)
				} else {
					m.projDescInput, _ = m.projDescInput.Update(msg)
				}
				return m, nil
			}
		} else if m.showTaskDeleteConfirm {
			// Ensure input handler mode is synchronized with UI state
			if m.inputHandler.GetMode() != input.TaskDeleteConfirmMode {
				m.inputHandler.SetMode(input.TaskDeleteConfirmMode)
			}

			// Handle task deletion confirmation dialog
			switch msg.String() {
			case "y", "Y":
				return m.executeTaskDeletion(), nil
			case "n", "N", "esc":
				m.showTaskDeleteConfirm = false
				m.taskToDelete = ""
				return m, nil
			}
			return m, nil
		} else {
			// Normal mode - handle navigation keys

			// Ensure input handler mode is synchronized with UI state
			if m.inputHandler.GetMode() != input.NormalMode {
				m.inputHandler.SetMode(input.NormalMode)
			}

			result := m.inputHandler.HandleKeyMsg(msg, &m)

			// Handle special keys through input handler
			if result.Handled && result.Cmd != nil {
				return m, result.Cmd
			}

			// Handle task progression directly
			switch msg.String() {
			case "q":
				return m, tea.Quit
			case "n":
				m.showForm = true
				m.focusedInput = 0
				m.nameInput.Focus()
				m.descInput.Blur()
				return m, nil
			case "p":
				m.showProjectSwitch = true
				return m, nil
			case "e":
				// Handle task editing - show edit form
				if selectedItem := m.Tasks[m.activeListIndex].SelectedItem(); selectedItem != nil {
					if task, ok := selectedItem.(Task); ok {
						m.ShowTaskEditForm(task.ID, task.Name, task.Desc)
						m.showTaskEditForm = true
						m.editingTaskID = task.ID
						m.nameInput.SetValue(task.Name)
						m.descInput.SetValue(task.Desc)
						m.focusedInput = 0
						m.nameInput.Focus()
						m.descInput.Blur()
					}
				}
				return m, nil
			case "d":
				// Handle task deletion - show confirmation dialog
				if selectedItem := m.Tasks[m.activeListIndex].SelectedItem(); selectedItem != nil {
					if task, ok := selectedItem.(Task); ok {
						m.showTaskDeleteConfirm = true
						m.taskToDelete = task.ID
					}
				}
				return m, nil
			case "enter":
				// Handle task selection - move to next status
				if selectedItem := m.Tasks[m.activeListIndex].SelectedItem(); selectedItem != nil {
					if task, ok := selectedItem.(Task); ok {
						m.MoveTaskToNextStatus(task.ID)
					}
				}
				return m, nil
			case "backspace":
				// Handle task selection - move to previous status
				if selectedItem := m.Tasks[m.activeListIndex].SelectedItem(); selectedItem != nil {
					if task, ok := selectedItem.(Task); ok {
						m.MoveTaskToPreviousStatus(task.ID)
					}
				}
				return m, nil
			}
		}
	case tea.WindowSizeMsg:
		// Store terminal dimensions
		m.width = msg.Width
		m.height = msg.Height

		h, v := defaultStyle.GetFrameSize()
		// Calculate equal column width (1/3 of terminal width)
		columnWidth := max(20, (msg.Width-(h*3))/3)
		m.Tasks[NotStarted].SetSize(columnWidth, msg.Height-v)
		m.Tasks[InProgress].SetSize(columnWidth, msg.Height-v)
		m.Tasks[Done].SetSize(columnWidth, msg.Height-v)
	}

	// Always update the active list if not in a form mode
	var cmd tea.Cmd
	if !m.showForm && !m.showProjectSwitch && !m.showProjectForm && !m.showTaskDeleteConfirm && !m.showTaskEditForm {
		m.Tasks[m.activeListIndex], cmd = m.Tasks[m.activeListIndex].Update(msg)
	}
	return m, cmd
}
