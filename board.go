package main

import (
	"fmt"
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
		Render("[p] Switch • [n] Add Task • [e] Edit Task • [d] Delete Task • [q] Quit")

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
	if m.showTaskDeleteConfirm {
		return m.renderTaskDeleteConfirm()
	}
	if m.showTaskEditForm {
		return m.renderTaskEditForm()
	}

	// Handle case when there are no projects
	if len(m.Projects) == 0 {
		return m.renderNoProjectsBoard()
	}

	// Render project header
	projectHeader := m.renderProjectHeader()

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

	return lipgloss.JoinVertical(
		lipgloss.Top,
		projectHeader,
		boardContent,
	)
}

func (m Model) renderNoProjectsBoard() string {
	title := lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorMauve)).
		Bold(true).
		Align(lipgloss.Center).
		Width(60).
		Render("No Projects")

	message := lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorText)).
		Align(lipgloss.Center).
		Width(60).
		Render("Create your first project to get started")

	instructions := lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorSubtext1)).
		Align(lipgloss.Center).
		Width(60).
		Render("[p] Create Project • [q] Quit")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		"",
		title,
		"",
		message,
		"",
		instructions,
	)

	form := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(ColorMauve)).
		Padding(2, 3).
		Width(70).
		Height(12).
		Render(content)

	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		form,
	)
}

func (m Model) renderTaskDeleteConfirm() string {
	// Find the task to delete
	var taskToDelete *Task
	activeProj := m.GetActiveProject()
	if activeProj != nil {
		for _, task := range activeProj.Tasks {
			if task.ID == m.taskToDelete {
				taskToDelete = &task
				break
			}
		}
	}

	if taskToDelete == nil {
		// Fallback to selected task if somehow taskToDelete is not set
		if selectedItem := m.Tasks[m.activeListIndex].SelectedItem(); selectedItem != nil {
			if task, ok := selectedItem.(Task); ok {
				taskToDelete = &task
			}
		}
		if taskToDelete == nil {
			return m.View() // Fallback to normal view
		}
	}

	title := lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorRed)).
		Bold(true).
		Align(lipgloss.Center).
		Width(60).
		Render("⚠️  Delete Task")

	taskName := lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorText)).
		Bold(true).
		Render(taskToDelete.Name)

	warningMessage := lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorText)).
		Align(lipgloss.Center).
		Width(60).
		Render(fmt.Sprintf("Delete task \"%s\"?", taskName))

	subWarning := lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorSubtext1)).
		Align(lipgloss.Center).
		Width(60).
		Render("This action cannot be undone.")

	instructions := lipgloss.NewStyle().
		Foreground(lipgloss.Color(ColorSubtext1)).
		Align(lipgloss.Center).
		Width(60).
		Render("[y] Yes, Delete • [n] No, Cancel")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		"",
		title,
		"",
		warningMessage,
		"",
		subWarning,
		"",
		instructions,
	)

	form := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(ColorRed)).
		Padding(2, 3).
		Width(70).
		Height(12).
		Render(content)

	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		form,
	)
}
