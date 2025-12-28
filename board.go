package main

import (
	"github.com/charmbracelet/bubbles/list"
	"kahn/internal/database"
	"kahn/internal/domain"
	repo "kahn/internal/repository"
	"kahn/internal/services"
	"kahn/internal/ui/components"
	"kahn/internal/ui/styles"
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
	taskService := services.NewTaskService(taskRepo, projectRepo)
	projectService := services.NewProjectService(projectRepo, taskRepo)

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

	// Apply list titles and styles
	taskLists[domain.NotStarted].Title = domain.NotStarted.ToString()
	if len(projects) > 0 {
		taskLists[domain.NotStarted].SetItems(convertTasksToListItems(projects[0].GetTasksByStatus(domain.NotStarted)))
	}

	taskLists[domain.InProgress].Title = domain.InProgress.ToString()
	if len(projects) > 0 {
		taskLists[domain.InProgress].SetItems(convertTasksToListItems(projects[0].GetTasksByStatus(domain.InProgress)))
	}

	taskLists[domain.Done].Title = domain.Done.ToString()
	if len(projects) > 0 {
		taskLists[domain.Done].SetItems(convertTasksToListItems(projects[0].GetTasksByStatus(domain.Done)))
	}

	// Apply proper title styles to all lists
	styles.ApplyListTitleStyles(taskLists)

	return &Model{
		Projects:               projects,
		ActiveProjectID:        activeProjectID,
		Tasks:                  taskLists,
		taskInputComponents:    taskInputComponents,
		projectInputComponents: projectInputComponents,
		board:                  components.NewBoard(),
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

func (m Model) View() string {
	if m.showForm {
		comps := m.GetActiveInputComponents()
		return comps.Render(m.formError, m.formErrorField, m.width, m.height)
	}
	if m.showProjectSwitch {
		return m.renderProjectSwitcher()
	}
	if m.showTaskDeleteConfirm {
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
		}

		return m.board.GetRenderer().RenderTaskDeleteConfirm(taskToDelete, m.width, m.height)
	}

	// Handle case when there are no projects
	if len(m.Projects) == 0 {
		return m.board.GetRenderer().RenderNoProjectsBoard(m.width, m.height)
	}

	// Render main board
	activeProj := m.GetActiveProject()
	if activeProj != nil {
		var taskLists [3]list.Model
		taskLists[domain.NotStarted] = m.Tasks[domain.NotStarted]
		taskLists[domain.InProgress] = m.Tasks[domain.InProgress]
		taskLists[domain.Done] = m.Tasks[domain.Done]

		return m.board.GetRenderer().RenderBoard(activeProj, taskLists, m.activeListIndex, m.width)
	}

	return ""
}
