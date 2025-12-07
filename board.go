package main

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

func NewModel(database *Database) *Model {
	defaultList := list.New([]list.Item{}, list.NewDefaultDelegate(), 100, 0)
	defaultList.SetShowHelp(false)
	taskLists := []list.Model{defaultList, defaultList, defaultList}

	nameInput, descInput := initializeInputs()
	projNameInput, projDescInput := initializeProjectInputs()

	projectDAO := NewProjectDAO(database.GetDB())
	projects, err := projectDAO.GetAll()
	if err != nil {
		projects = []Project{}
	}

	if len(projects) == 0 {
		newProject := NewProject("Default Project", "A default project for your tasks", ColorBlue)
		if err := projectDAO.Create(newProject); err != nil {
			projects = []Project{}
		} else {
			projects = []Project{*newProject}
		}
	}

	var activeProjectID string
	if len(projects) > 0 {
		activeProjectID = projects[0].ID
	}

	taskLists[NotStarted].Title = NotStarted.ToString()
	if len(projects) > 0 {
		taskLists[NotStarted].SetItems(convertTasksToListItems(projects[0].GetTasksByStatus(NotStarted)))
	}
	taskLists[NotStarted].Styles.Title = lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorBlue)).
		Bold(true).
		Align(lipgloss.Center)

	taskLists[InProgress].Title = InProgress.ToString()
	if len(projects) > 0 {
		taskLists[InProgress].SetItems(convertTasksToListItems(projects[0].GetTasksByStatus(InProgress)))
	}
	taskLists[InProgress].Styles.Title = lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorYellow)).
		Bold(true).
		Align(lipgloss.Center)

	taskLists[Done].Title = Done.ToString()
	if len(projects) > 0 {
		taskLists[Done].SetItems(convertTasksToListItems(projects[0].GetTasksByStatus(Done)))
	}
	taskLists[Done].Styles.Title = lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorGreen)).
		Bold(true).
		Align(lipgloss.Center)

	return &Model{
		Projects:        projects,
		ActiveProjectID: activeProjectID,
		Tasks:           taskLists,
		nameInput:       nameInput,
		descInput:       descInput,
		projNameInput:   projNameInput,
		projDescInput:   projDescInput,
		width:           80,
		height:          24,
		database:        database,
		projectDAO:      NewProjectDAO(database.GetDB()),
		taskDAO:         NewTaskDAO(database.GetDB()),
	}
}

func convertTasksToListItems(tasks []Task) []list.Item {
	items := make([]list.Item, len(tasks))
	for i, task := range tasks {
		items[i] = task
	}
	return items
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

	columnWidth := m.Tasks[0].Width()

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
	return boardContent
}
