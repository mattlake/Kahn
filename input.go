package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.showForm {
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
					// Create new task and add to active project
					newTask := NewTask(m.nameInput.Value(), m.descInput.Value(), m.ActiveProjectID)

					// Add to active project
					activeProj := m.GetActiveProject()
					if activeProj != nil {
						activeProj.AddTask(*newTask)
						m.updateTaskLists()
					} else {
						// Debug: No active project found
						return m, nil
					}

					// Reset form
					m.showForm = false
					m.nameInput.Reset()
					m.descInput.Reset()
					m.focusedInput = 0
					m.nameInput.Focus()
				}
				return m, nil
			}

			// Update the focused input
			if m.focusedInput == 0 {
				m.nameInput, _ = m.nameInput.Update(msg)
			} else {
				m.descInput, _ = m.descInput.Update(msg)
			}
		} else if m.showProjectSwitch {
			// Project switcher mode key handling
			switch msg.String() {
			case "esc":
				m.showProjectSwitch = false
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
					// Create new project
					newProject := NewProject(m.projNameInput.Value(), m.projDescInput.Value(), ColorBlue)

					if err := newProject.Validate(); err == nil {
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
				}
				return m, nil
			}

			// Update the focused input
			if m.focusedProjInput == 0 {
				m.projNameInput, _ = m.projNameInput.Update(msg)
			} else {
				m.projDescInput, _ = m.projDescInput.Update(msg)
			}
		} else {
			// Normal mode key handling
			switch msg.String() {
			case "q":
				return m, tea.Quit
			case "l":
				m.NextList()
			case "h":
				m.Prevlist()
			case "a":
				m.showForm = true
				m.focusedInput = 0
				m.nameInput.Focus()
				m.descInput.Blur()
				return m, nil
			case "p":
				m.showProjectSwitch = true
				return m, nil
			case "enter":
				// Handle task selection - move to next status
				if selectedItem := m.Tasks[m.activeListIndex].SelectedItem(); selectedItem != nil {
					if task, ok := selectedItem.(Task); ok {
						m.moveTaskToNextStatus(task)
					}
				}
				return m, nil
			case "backspace":
				// Handle task selection - move to previous status
				if selectedItem := m.Tasks[m.activeListIndex].SelectedItem(); selectedItem != nil {
					if task, ok := selectedItem.(Task); ok {
						m.moveTaskToPreviousStatus(task)
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

	var cmd tea.Cmd
	if !m.showForm && !m.showProjectSwitch && !m.showProjectForm {
		m.Tasks[m.activeListIndex], cmd = m.Tasks[m.activeListIndex].Update(msg)
	}
	return m, cmd
}

func (m *Model) NextList() {
	if m.activeListIndex == Done {
		m.activeListIndex = NotStarted
	} else {
		m.activeListIndex++
	}
}

func (m *Model) Prevlist() {
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

	// Update task lists with active project tasks
	m.Tasks[NotStarted].SetItems(convertTasksToListItems(activeProj.GetTasksByStatus(NotStarted)))
	m.Tasks[InProgress].SetItems(convertTasksToListItems(activeProj.GetTasksByStatus(InProgress)))
	m.Tasks[Done].SetItems(convertTasksToListItems(activeProj.GetTasksByStatus(Done)))
}

func (m *Model) moveTaskToNextStatus(task Task) {
	activeProj := m.GetActiveProject()
	if activeProj == nil {
		return
	}

	// Determine next status
	var nextStatus Status
	switch task.Status {
	case NotStarted:
		nextStatus = InProgress
	case InProgress:
		nextStatus = Done
	case Done:
		// Cycle back to NotStarted
		nextStatus = NotStarted
	}

	// Update task status
	activeProj.UpdateTaskStatus(task.ID, nextStatus)
	m.updateTaskLists()
}

func (m *Model) moveTaskToPreviousStatus(task Task) {
	activeProj := m.GetActiveProject()
	if activeProj == nil {
		return
	}

	// Determine previous status
	var prevStatus Status
	switch task.Status {
	case NotStarted:
		// Cycle back to Done
		prevStatus = Done
	case InProgress:
		prevStatus = NotStarted
	case Done:
		prevStatus = InProgress
	}

	// Update task status
	activeProj.UpdateTaskStatus(task.ID, prevStatus)
	m.updateTaskLists()
}
