package main

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

func NewModel() *Model {
	defaultList := list.New([]list.Item{}, list.NewDefaultDelegate(), 100, 0)
	defaultList.SetShowHelp(false)
	taskLists := []list.Model{defaultList, defaultList, defaultList}

	nameInput, descInput := initializeInputs()
	projNameInput, projDescInput := initializeProjectInputs()

	// Create default project with sample tasks
	defaultProject := NewProject("Default Project", "A default project for your tasks", ColorBlue)
	defaultProject.Tasks = []Task{
		*NewTask("Task1", "Task1 Description", defaultProject.ID),
		*NewTask("Task2", "Task2 Description", defaultProject.ID),
		*NewTask("Task3", "Task3 Description", defaultProject.ID),
	}

	// Create another sample project
	workProject := NewProject("Work Project", "Work-related tasks", ColorGreen)
	workProject.Tasks = []Task{
		*NewTask("Task1", "Task1 Description", workProject.ID),
	}

	projects := []Project{*defaultProject, *workProject}

	// Initialize task lists with default project tasks
	taskLists[NotStarted].Title = NotStarted.ToString()
	taskLists[NotStarted].SetItems(convertTasksToListItems(defaultProject.GetTasksByStatus(NotStarted)))
	taskLists[NotStarted].Styles.Title = lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorBlue)).
		Bold(true).
		Align(lipgloss.Center)

	taskLists[InProgress].Title = InProgress.ToString()
	taskLists[InProgress].SetItems(convertTasksToListItems(defaultProject.GetTasksByStatus(InProgress)))
	taskLists[InProgress].Styles.Title = lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorYellow)).
		Bold(true).
		Align(lipgloss.Center)

	taskLists[Done].Title = Done.ToString()
	taskLists[Done].SetItems(convertTasksToListItems(defaultProject.GetTasksByStatus(Done)))
	taskLists[Done].Styles.Title = lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorGreen)).
		Bold(true).
		Align(lipgloss.Center)

	return &Model{
		Projects:        projects,
		ActiveProjectID: defaultProject.ID,
		Tasks:           taskLists,
		nameInput:       nameInput,
		descInput:       descInput,
		projNameInput:   projNameInput,
		projDescInput:   projDescInput,
		width:           80, // default
		height:          24, // default
	}
}

func convertTasksToListItems(tasks []Task) []list.Item {
	items := make([]list.Item, len(tasks))
	for i, task := range tasks {
		items[i] = task
	}
	return items
}

func (m Model) renderProjectHeader() string {
	activeProj := m.GetActiveProject()
	if activeProj == nil {
		return ""
	}

	// Create a more prominent project indicator
	projectLabel := lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorSubtext1)).
		Render("Project:")

	projectNameText := lipgloss.NewStyle().
		Foreground(lipgloss.Color(activeProj.Color)).
		Bold(true).
		Render(activeProj.Name)

	helpText := lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorSubtext1)).
		Render("[p] Switch • [a] Add Task • [q] Quit")

	// Create a more prominent header with better visual hierarchy
	headerContent := lipgloss.JoinHorizontal(
		lipgloss.Left,
		projectLabel,
		lipgloss.NewStyle().Render(" "),
		projectNameText,
		lipgloss.NewStyle().Width(m.width-len(activeProj.Name)-len("Project: ")-25).Render(""),
		helpText,
	)

	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(activeProj.Color)).
		Padding(0, 1).
		Background(lipgloss.Color(ColorSurface0)).
		Width(m.width).
		Render(headerContent)
}

func (m Model) View() string {
	if m.showForm {
		return m.renderForm()
	}
	if m.showProjectSwitch {
		return m.renderProjectSwitcher()
	}
	if m.showProjectForm {
		return m.renderProjectForm()
	}

	// Render project header
	projectHeader := m.renderProjectHeader()

	// Get current width of first list to use as column width
	columnWidth := m.Tasks[0].Width()

	// Apply width constraints to ensure equal columns
	notStartedView := defaultStyle.Width(columnWidth).Render(m.Tasks[NotStarted].View())
	inProgressView := defaultStyle.Width(columnWidth).Render(m.Tasks[InProgress].View())
	doneView := defaultStyle.Width(columnWidth).Render(m.Tasks[Done].View())

	focusedNotStartedView := focusedStyle.Width(columnWidth).Render(m.Tasks[NotStarted].View())
	focusedInProgressView := focusedStyle.Width(columnWidth).Render(m.Tasks[InProgress].View())
	focusedDoneView := focusedStyle.Width(columnWidth).Render(m.Tasks[Done].View())

	boardContent := ""
	switch m.activeListIndex {
	case InProgress:
		boardContent = lipgloss.JoinHorizontal(
			lipgloss.Left,
			notStartedView,
			focusedInProgressView,
			doneView,
		)
	case Done:
		boardContent = lipgloss.JoinHorizontal(
			lipgloss.Left,
			notStartedView,
			inProgressView,
			focusedDoneView,
		)
	default:
		boardContent = lipgloss.JoinHorizontal(
			lipgloss.Left,
			focusedNotStartedView,
			inProgressView,
			doneView,
		)
	}

	return lipgloss.JoinVertical(
		lipgloss.Top,
		projectHeader,
		boardContent,
	)
}
