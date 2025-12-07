package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m KahnModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.showForm {
			// Form mode key handling
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
					// Create new task and add to Not Started list
					newTask := Task{
						Name:   m.nameInput.Value(),
						Desc:   m.descInput.Value(),
						Status: NotStarted,
					}

					currentItems := m.Tasks[NotStarted].Items()
					m.Tasks[NotStarted].SetItems(append(currentItems, newTask))

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
			}
		}
	case tea.WindowSizeMsg:
		// Store terminal dimensions
		m.width = msg.Width
		m.height = msg.Height

		h, v := defaultStyle.GetFrameSize()
		// Calculate equal column width (1/3 of terminal width)
		columnWidth := (msg.Width - (h * 3)) / 3
		if columnWidth < 20 {
			columnWidth = 20 // Minimum width
		}
		m.Tasks[NotStarted].SetSize(columnWidth, msg.Height-v)
		m.Tasks[InProgress].SetSize(columnWidth, msg.Height-v)
		m.Tasks[Done].SetSize(columnWidth, msg.Height-v)
	}

	var cmd tea.Cmd
	if !m.showForm {
		m.Tasks[m.activeListIndex], cmd = m.Tasks[m.activeListIndex].Update(msg)
	}
	return m, cmd
}

func (m *KahnModel) NextList() {
	if m.activeListIndex == Done {
		m.activeListIndex = NotStarted
	} else {
		m.activeListIndex++
	}
}

func (m *KahnModel) Prevlist() {
	if m.activeListIndex == NotStarted {
		m.activeListIndex = Done
	} else {
		m.activeListIndex--
	}
}
