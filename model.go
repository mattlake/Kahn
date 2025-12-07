package main

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	Projects          []Project
	ActiveProjectID   string
	Tasks             []list.Model
	activeListIndex   Status
	showForm          bool
	showProjectSwitch bool
	showProjectForm   bool
	nameInput         textinput.Model
	descInput         textinput.Model
	projNameInput     textinput.Model
	projDescInput     textinput.Model
	focusedInput      int // 0 for name, 1 for desc
	focusedProjInput  int // 0 for name, 1 for desc
	width             int
	height            int
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
