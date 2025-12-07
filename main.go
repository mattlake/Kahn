package main

import (
	"log"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Status int

const (
	NotStarted Status = iota
	InProgress
	Done
)

func main() {
	m := NewModel()
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

type Model struct {
	Tasks           []list.Model
	activeListIndex Status
}

func NewModel() *Model {
	defaultList := list.New([]list.Item{}, list.NewDefaultDelegate(), 100, 0)
	defaultList.SetShowHelp(false)
	taskLists := []list.Model{defaultList, defaultList, defaultList}

	taskLists[NotStarted].Title = NotStarted.ToString()
	taskLists[NotStarted].SetItems(
		[]list.Item{
			Task{Name: "Task1", Desc: "Task1 Description", Status: NotStarted},
			Task{Name: "Task2", Desc: "Task2 Description", Status: NotStarted},
			Task{Name: "Task3", Desc: "Task3 Description", Status: NotStarted},
		})

	taskLists[InProgress].Title = InProgress.ToString()
	taskLists[InProgress].SetItems(
		[]list.Item{
			Task{Name: "Task1", Desc: "Task1 Description", Status: InProgress},
		})

	taskLists[Done].Title = Done.ToString()
	taskLists[Done].SetItems(
		[]list.Item{
			Task{Name: "Task1", Desc: "Task1 Description", Status: Done},
			Task{Name: "Task2", Desc: "Task2 Description", Status: Done},
		})

	return &Model{
		Tasks: taskLists,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "l":
			m.NextList()
		case "h":
			m.Prevlist()
		}
	case tea.WindowSizeMsg:
		h, v := defaultStyle.GetFrameSize()
		m.Tasks[NotStarted].SetSize(msg.Width-h, msg.Height-v)
		m.Tasks[InProgress].SetSize(msg.Width-h, msg.Height-v)
		m.Tasks[Done].SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.Tasks[m.activeListIndex], cmd = m.Tasks[m.activeListIndex].Update(msg)
	return m, cmd
}

func (m Model) View() string {
	switch m.activeListIndex {
	case InProgress:
		return lipgloss.JoinHorizontal(
			lipgloss.Left,
			defaultStyle.Render(m.Tasks[NotStarted].View()),
			focusedStyle.Render(m.Tasks[InProgress].View()),
			defaultStyle.Render(m.Tasks[Done].View()),
		)
	case Done:
		return lipgloss.JoinHorizontal(
			lipgloss.Left,
			defaultStyle.Render(m.Tasks[NotStarted].View()),
			defaultStyle.Render(m.Tasks[InProgress].View()),
			focusedStyle.Render(m.Tasks[Done].View()),
		)
	default:
		return lipgloss.JoinHorizontal(
			lipgloss.Left,
			focusedStyle.Render(m.Tasks[NotStarted].View()),
			defaultStyle.Render(m.Tasks[InProgress].View()),
			defaultStyle.Render(m.Tasks[Done].View()),
		)
	}
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

func (s Status) ToString() string {
	switch s {
	case NotStarted:
		return "Not Started"
	case InProgress:
		return "In Progress"
	case Done:
		return "Done"
	default:
		return "Placeholder"
	}
}

// Lipgloss
var defaultStyle = lipgloss.NewStyle().
	Margin(1, 2).
	Border(lipgloss.HiddenBorder()).
	Padding(1, 2)

var focusedStyle = lipgloss.NewStyle().
	Margin(1, 2).
	Border(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("21")).
	Padding(1, 2)
