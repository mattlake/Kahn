package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"kahn/internal/domain"
	"kahn/internal/ui/input"
	"kahn/internal/ui/styles"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Handle form modes first
		if m.showForm {
			// Use InputHandler for all form key handling
			action := m.inputHandler.HandleKeyMsg(msg, &m)
			if action.Handled {
				return m, action.Cmd
			}

			// Handle textinput updates for unhandled keys
			comps := m.GetActiveInputComponents()
			if comps.FocusedField == 0 {
				updatedName, cmd := comps.NameInput.Update(msg)
				comps.NameInput = updatedName
				return m, cmd
			} else {
				updatedDesc, cmd := comps.DescInput.Update(msg)
				comps.DescInput = updatedDesc
				return m, cmd
			}
		}
		if m.showProjectSwitch || m.showProjectDeleteConfirm {
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
				m.ShowProjectForm()
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
				m.ShowTaskForm()
				return m, nil
			case "p":
				m.showProjectSwitch = true
				return m, nil
			case "e":
				// Handle task editing - show edit form
				if selectedItem := m.Tasks[m.activeListIndex].SelectedItem(); selectedItem != nil {
					if task, ok := selectedItem.(domain.Task); ok {
						m.ShowTaskEditForm(task.ID, task.Name, task.Desc)
						m.showForm = true // Show the unified form
					}
				}
				return m, nil
			case "d":
				// Handle task deletion - show confirmation dialog
				if selectedItem := m.Tasks[m.activeListIndex].SelectedItem(); selectedItem != nil {
					if task, ok := selectedItem.(domain.Task); ok {
						m.showTaskDeleteConfirm = true
						m.taskToDelete = task.ID
					}
				}
				return m, nil
			case "enter":
				// Handle task selection - move to next status
				if selectedItem := m.Tasks[m.activeListIndex].SelectedItem(); selectedItem != nil {
					if task, ok := selectedItem.(domain.Task); ok {
						m.MoveTaskToNextStatus(task.ID)
					}
				}
				return m, nil
			case "backspace":
				// Handle task selection - move to previous status
				if selectedItem := m.Tasks[m.activeListIndex].SelectedItem(); selectedItem != nil {
					if task, ok := selectedItem.(domain.Task); ok {
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

		h, v := styles.DefaultStyle.GetFrameSize()
		// Calculate equal column width (1/3 of terminal width)
		columnWidth := max(20, (msg.Width-(h*3))/3)
		m.Tasks[domain.NotStarted].SetSize(columnWidth, msg.Height-v)
		m.Tasks[domain.InProgress].SetSize(columnWidth, msg.Height-v)
		m.Tasks[domain.Done].SetSize(columnWidth, msg.Height-v)
	}

	// Always update the active list if not in a form mode
	var cmd tea.Cmd
	if !m.showForm && !m.showProjectSwitch && !m.showTaskDeleteConfirm {
		m.Tasks[m.activeListIndex], cmd = m.Tasks[m.activeListIndex].Update(msg)
	}
	return m, cmd
}
