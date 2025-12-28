package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
	"kahn/internal/database"
	"kahn/internal/domain"
	repo "kahn/internal/repository"
	"kahn/pkg/colors"
	"kahn/pkg/input"
)

func NewModel(database *database.Database) *Model {
	defaultList := list.New([]list.Item{}, list.NewDefaultDelegate(), 100, 0)
	defaultList.SetShowHelp(false)
	taskLists := []list.Model{defaultList, defaultList, defaultList}

	taskInputComponents := &input.InputComponents{}
	projectInputComponents := &input.InputComponents{}

	// Create repositories
	taskRepo := repo.NewSQLiteTaskRepository(database.GetDB())
	projectRepo := repo.NewSQLiteProjectRepository(database.GetDB())

	// Create services
	taskService := NewTaskService(taskRepo, projectRepo)
	projectService := NewProjectService(projectRepo, taskRepo)

	projects, err := projectService.GetAllProjects()
	if err != nil {
		projects = []domain.Project{}
	}

	for i := range projects {
		tasks, err := taskService.GetTasksByProject(projects[i].ID)
		if err != nil {
			projects[i].Tasks = []domain.Task{}
		} else {
			projects[i].Tasks = tasks
		}
	}

	if len(projects) == 0 {
		newProject, err := projectService.CreateProject("Default Project", "A default project for your tasks")
		if err != nil {
			projects = []domain.Project{}
		} else {
			projects = []domain.Project{*newProject}
		}
	}

	var activeProjectID string
	if len(projects) > 0 {
		activeProjectID = projects[0].ID
	}

	taskLists[domain.NotStarted].Title = domain.NotStarted.ToString()
	if len(projects) > 0 {
		taskLists[domain.NotStarted].SetItems(convertTasksToListItems(projects[0].GetTasksByStatus(domain.NotStarted)))
	}
	taskLists[domain.NotStarted].Styles.Title = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colors.Blue)).
		Bold(true).
		Align(lipgloss.Center)

	taskLists[domain.InProgress].Title = domain.InProgress.ToString()
	if len(projects) > 0 {
		taskLists[domain.InProgress].SetItems(convertTasksToListItems(projects[0].GetTasksByStatus(domain.InProgress)))
	}
	taskLists[domain.InProgress].Styles.Title = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colors.Yellow)).
		Bold(true).
		Align(lipgloss.Center)

	taskLists[domain.Done].Title = domain.Done.ToString()
	if len(projects) > 0 {
		taskLists[domain.Done].SetItems(convertTasksToListItems(projects[0].GetTasksByStatus(domain.Done)))
	}
	taskLists[domain.Done].Styles.Title = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colors.Green)).
		Bold(true).
		Align(lipgloss.Center)

	return &Model{
		Projects:               projects,
		ActiveProjectID:        activeProjectID,
		Tasks:                  taskLists,
		taskInputComponents:    taskInputComponents,
		projectInputComponents: projectInputComponents,
		width:                  80,
		height:                 24,
		database:               database,
		taskService:            taskService,
		projectService:         projectService,
		inputHandler:           input.NewHandler(),
	}
}

func convertTasksToListItems(tasks []domain.Task) []list.Item {
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
		Foreground(lipgloss.Color(colors.Subtext1)).
		Render("Project:")

	projectNameText := lipgloss.NewStyle().
		Foreground(lipgloss.Color(activeProj.Color)).
		Bold(true).
		Render(activeProj.Name)

	helpText := lipgloss.NewStyle().
		Foreground(lipgloss.Color(colors.Subtext1)).
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
		Background(lipgloss.Color(colors.Surface0)).
		Width(m.width).
		Render(headerContent)
}

func (m Model) View() string {
	if m.showForm {
		comps := m.GetActiveInputComponents()
		return comps.Render(m.formError, m.formErrorField, m.width, m.height)
	}
	if m.showProjectSwitch {
		return m.renderProjectSwitcher()
	}
	if m.showTaskDeleteConfirm {
		return m.renderTaskDeleteConfirm()
	}

	// Handle case when there are no projects
	if len(m.Projects) == 0 {
		return m.renderNoProjectsBoard()
	}

	// Render project header
	projectHeader := m.renderProjectHeader()

	columnWidth := m.Tasks[0].Width()

	notStartedView := defaultStyle.Width(columnWidth).Render(m.Tasks[domain.NotStarted].View())
	inProgressView := defaultStyle.Width(columnWidth).Render(m.Tasks[domain.InProgress].View())
	doneView := defaultStyle.Width(columnWidth).Render(m.Tasks[domain.Done].View())

	focusedNotStartedView := focusedStyle.Width(columnWidth).Render(m.Tasks[domain.NotStarted].View())
	focusedInProgressView := focusedStyle.Width(columnWidth).Render(m.Tasks[domain.InProgress].View())
	focusedDoneView := focusedStyle.Width(columnWidth).Render(m.Tasks[domain.Done].View())

	boardContent := ""
	switch m.activeListIndex {
	case domain.InProgress:
		boardContent = lipgloss.JoinHorizontal(
			lipgloss.Left,
			notStartedView,
			focusedInProgressView,
			doneView,
		)
	case domain.Done:
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
		Foreground(lipgloss.Color(colors.Mauve)).
		Bold(true).
		Align(lipgloss.Center).
		Width(60).
		Render("No Projects")

	message := lipgloss.NewStyle().
		Foreground(lipgloss.Color(colors.Text)).
		Align(lipgloss.Center).
		Width(60).
		Render("Create your first project to get started")

	instructions := lipgloss.NewStyle().
		Foreground(lipgloss.Color(colors.Subtext1)).
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
		BorderForeground(lipgloss.Color(colors.Mauve)).
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
	var taskToDelete *domain.Task
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
			if task, ok := selectedItem.(domain.Task); ok {
				taskToDelete = &task
			}
		}
		if taskToDelete == nil {
			return m.View() // Fallback to normal view
		}
	}

	title := lipgloss.NewStyle().
		Foreground(lipgloss.Color(colors.Red)).
		Bold(true).
		Align(lipgloss.Center).
		Width(60).
		Render("⚠️  Delete Task")

	taskName := lipgloss.NewStyle().
		Foreground(lipgloss.Color(colors.Text)).
		Bold(true).
		Render(taskToDelete.Name)

	warningMessage := lipgloss.NewStyle().
		Foreground(lipgloss.Color(colors.Text)).
		Align(lipgloss.Center).
		Width(60).
		Render(fmt.Sprintf("Delete task \"%s\"?", taskName))

	subWarning := lipgloss.NewStyle().
		Foreground(lipgloss.Color(colors.Subtext1)).
		Align(lipgloss.Center).
		Width(60).
		Render("This action cannot be undone.")

	instructions := lipgloss.NewStyle().
		Foreground(lipgloss.Color(colors.Subtext1)).
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
		BorderForeground(lipgloss.Color(colors.Red)).
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
