package main

import (
	"log"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Catppuccin Mocha palette
const (
	// Primary colors
	ColorMauve    = "#cba6f7"
	ColorBlue     = "#89b4fa"
	ColorLavender = "#b4befe"
	ColorSapphire = "#74c7ec"

	// Text colors
	ColorText     = "#cdd6f4"
	ColorSubtext1 = "#bac2de"
	ColorSubtext0 = "#a6adc8"

	// Surface colors
	ColorSurface0 = "#313244"
	ColorSurface1 = "#45475a"
	ColorSurface2 = "#585b70"
	ColorBase     = "#1e1e2e"

	// Border colors
	ColorOverlay2 = "#9399b2"
	ColorOverlay1 = "#7f849c"
	ColorOverlay0 = "#6c7086"

	// Status colors
	ColorGreen  = "#a6e3a1"
	ColorYellow = "#f9e2af"
	ColorRed    = "#f38ba8"
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
	showForm        bool
	nameInput       textinput.Model
	descInput       textinput.Model
	focusedInput    int // 0 for name, 1 for desc
}

func initializeInputs() (textinput.Model, textinput.Model) {
	name := textinput.New()
	name.Placeholder = "Task name"
	name.PlaceholderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(ColorSubtext0))
	name.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(ColorText))
	name.CursorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(ColorMauve))
	name.Focus()
	name.CharLimit = 50
	name.Width = 40

	desc := textinput.New()
	desc.Placeholder = "Task description"
	desc.PlaceholderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(ColorSubtext0))
	desc.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(ColorText))
	desc.CursorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(ColorMauve))
	desc.CharLimit = 100
	desc.Width = 40

	return name, desc
}

func NewModel() *Model {
	defaultList := list.New([]list.Item{}, list.NewDefaultDelegate(), 100, 0)
	defaultList.SetShowHelp(false)
	taskLists := []list.Model{defaultList, defaultList, defaultList}

	nameInput, descInput := initializeInputs()

	taskLists[NotStarted].Title = NotStarted.ToString()
	taskLists[NotStarted].SetItems(
		[]list.Item{
			Task{Name: "Task1", Desc: "Task1 Description", Status: NotStarted},
			Task{Name: "Task2", Desc: "Task2 Description", Status: NotStarted},
			Task{Name: "Task3", Desc: "Task3 Description", Status: NotStarted},
		})
	taskLists[NotStarted].Styles.Title = lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorBlue)).
		Bold(true).
		Align(lipgloss.Center)

	taskLists[InProgress].Title = InProgress.ToString()
	taskLists[InProgress].SetItems(
		[]list.Item{
			Task{Name: "Task1", Desc: "Task1 Description", Status: InProgress},
		})
	taskLists[InProgress].Styles.Title = lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorYellow)).
		Bold(true).
		Align(lipgloss.Center)

	taskLists[Done].Title = Done.ToString()
	taskLists[Done].SetItems(
		[]list.Item{
			Task{Name: "Task1", Desc: "Task1 Description", Status: Done},
			Task{Name: "Task2", Desc: "Task2 Description", Status: Done},
		})
	taskLists[Done].Styles.Title = lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorGreen)).
		Bold(true).
		Align(lipgloss.Center)

	return &Model{
		Tasks:     taskLists,
		nameInput: nameInput,
		descInput: descInput,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m Model) View() string {
	if m.showForm {
		return m.renderForm()
	}

	// Get the current width of the first list to use as column width
	columnWidth := m.Tasks[0].Width()

	// Apply width constraints to ensure equal columns
	notStartedView := defaultStyle.Width(columnWidth).Render(m.Tasks[NotStarted].View())
	inProgressView := defaultStyle.Width(columnWidth).Render(m.Tasks[InProgress].View())
	doneView := defaultStyle.Width(columnWidth).Render(m.Tasks[Done].View())

	focusedNotStartedView := focusedStyle.Width(columnWidth).Render(m.Tasks[NotStarted].View())
	focusedInProgressView := focusedStyle.Width(columnWidth).Render(m.Tasks[InProgress].View())
	focusedDoneView := focusedStyle.Width(columnWidth).Render(m.Tasks[Done].View())

	switch m.activeListIndex {
	case InProgress:
		return lipgloss.JoinHorizontal(
			lipgloss.Left,
			notStartedView,
			focusedInProgressView,
			doneView,
		)
	case Done:
		return lipgloss.JoinHorizontal(
			lipgloss.Left,
			notStartedView,
			inProgressView,
			focusedDoneView,
		)
	default:
		return lipgloss.JoinHorizontal(
			lipgloss.Left,
			focusedNotStartedView,
			inProgressView,
			doneView,
		)
	}
}

func (m Model) renderForm() string {
	formTitle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorMauve)).
		Bold(true).
		Align(lipgloss.Center).
		Width(50).
		Render("Add New Task")

	nameLabel := lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorText)).
		Bold(true).
		Render("Task Name:")

	descLabel := lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorText)).
		Bold(true).
		Render("Description:")

	instructions := lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorSubtext1)).
		Align(lipgloss.Center).
		Width(50).
		Render("Tab: Switch fields | Enter: Submit | Esc: Cancel")

	// Highlight focused input
	nameField := m.nameInput.View()
	descField := m.descInput.View()

	if m.focusedInput == 0 {
		nameField = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(ColorMauve)).
			Padding(0, 1).
			Render(nameField)
		descField = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(ColorOverlay1)).
			Padding(0, 1).
			Render(descField)
	} else {
		nameField = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(ColorOverlay1)).
			Padding(0, 1).
			Render(nameField)
		descField = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(ColorMauve)).
			Padding(0, 1).
			Render(descField)
	}

	formContent := lipgloss.JoinVertical(
		lipgloss.Left,
		"",
		formTitle,
		"",
		nameLabel,
		nameField,
		"",
		descLabel,
		descField,
		"",
		instructions,
	)

	form := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(ColorMauve)).
		Padding(2, 3).
		Width(60).
		Height(20).
		Render(formContent)

	return lipgloss.Place(
		50, 20,
		lipgloss.Center, lipgloss.Center,
		form,
	)
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

// Lipgloss styles with Catppuccin colors
var defaultStyle = lipgloss.NewStyle().
	Margin(1, 2).
	Border(lipgloss.HiddenBorder()).
	Padding(1, 2)

var focusedStyle = lipgloss.NewStyle().
	Margin(1, 2).
	Border(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color(ColorMauve)).
	Padding(1, 2)
