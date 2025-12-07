package main

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	Tasks           []list.Model
	activeListIndex Status
	showForm        bool
	nameInput       textinput.Model
	descInput       textinput.Model
	focusedInput    int // 0 for name, 1 for desc
	width           int
	height          int
}

func (m Model) Init() tea.Cmd {
	return nil
}
