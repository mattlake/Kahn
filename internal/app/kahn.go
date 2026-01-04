package app

import (
	"fmt"

	"kahn/internal/database"
	"kahn/internal/domain"
	repo "kahn/internal/repository"
	"kahn/internal/services"
	"kahn/internal/ui/components"
	"kahn/internal/ui/input"
	"kahn/internal/ui/styles"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type KahnModel struct {
	Projects        []domain.Project
	ActiveProjectID string
	width           int
	height          int
	database        *database.Database
	inputHandler    *input.Handler
	taskService     *services.TaskService
	projectService  *services.ProjectService
	board           *components.Board
	projectSwitcher *components.ProjectSwitcher
	version         string

	formState    *FormState
	confirmState *ConfirmationState
	navState     *NavigationState
}

func (km KahnModel) Init() tea.Cmd {
	return nil
}

func (km KahnModel) View() string {
	if km.formState.IsShowingForm() {
		comps := km.formState.GetActiveInputComponents()
		formError, formErrorField := km.formState.GetError()
		return comps.Render(formError, formErrorField, km.width, km.height)
	}
	if km.navState.IsShowingProjectSwitch() {
		return km.projectSwitcher.RenderSwitcherWithError(km.Projects, km.ActiveProjectID, km.confirmState.IsShowingProjectDeleteConfirm(), km.confirmState.GetProjectToDelete(), km.confirmState.GetError(), km.width, km.height)
	}
	if km.confirmState.IsShowingTaskDeleteConfirm() {
		var taskToDelete *domain.Task
		activeProj := km.GetActiveProject()
		if activeProj != nil {
			for _, task := range activeProj.Tasks {
				if task.ID == km.confirmState.GetTaskToDelete() {
					taskToDelete = &task
					break
				}
			}
		}

		if taskToDelete == nil {
			if selectedItem := km.navState.GetActiveList().SelectedItem(); selectedItem != nil {
				if taskWrapper, ok := selectedItem.(styles.TaskWithTitle); ok {
					taskToDelete = &taskWrapper.Task
				}
			}
		}

		return km.board.GetRenderer().RenderTaskDeleteConfirmWithError(taskToDelete, km.confirmState.GetError(), km.width, km.height)
	}

	if len(km.Projects) == 0 {
		return km.board.GetRenderer().RenderNoProjectsBoard(km.width, km.height)
	}

	activeProj := km.GetActiveProject()
	if activeProj != nil {
		var taskLists [3]list.Model
		taskLists[domain.NotStarted] = km.navState.Tasks[domain.NotStarted]
		taskLists[domain.InProgress] = km.navState.Tasks[domain.InProgress]
		taskLists[domain.Done] = km.navState.Tasks[domain.Done]

		return km.board.GetRenderer().RenderBoard(activeProj, taskLists, km.navState.GetActiveListIndex(), km.width, km.version)
	}

	return ""
}

func (km *KahnModel) GetActiveProject() *domain.Project {
	for i, proj := range km.Projects {
		if proj.ID == km.ActiveProjectID {
			return &km.Projects[i]
		}
	}
	return nil
}

func (km *KahnModel) GetActiveProjectID() string {
	return km.ActiveProjectID
}

func (km *KahnModel) CreateTask(name, description string) error {
	return km.CreateTaskWithPriority(name, description, domain.Low)
}

func (km *KahnModel) CreateTaskWithPriority(name, description string, priority domain.Priority) error {
	activeProj := km.GetActiveProject()
	if activeProj == nil {
		return nil
	}

	newTask, err := km.taskService.CreateTask(name, description, km.ActiveProjectID, domain.RegularTask, priority)
	if err != nil {
		return err
	}

	activeProj.AddTask(*newTask)

	km.navState.MarkListDirty(domain.NotStarted)
	km.updateTaskLists()

	return nil
}

func (km *KahnModel) UpdateTask(id, name, description string, priority domain.Priority, taskType domain.TaskType) error {
	task, err := km.taskService.UpdateTask(id, name, description, taskType, priority)
	if err != nil {
		return err
	}

	activeProj := km.GetActiveProject()
	if activeProj != nil {
		var taskStatus domain.Status
		for i, t := range activeProj.Tasks {
			if t.ID == id {
				taskStatus = t.Status // Save status before update
				activeProj.Tasks[i].Name = task.Name
				activeProj.Tasks[i].Desc = task.Desc
				activeProj.Tasks[i].Priority = task.Priority
				activeProj.Tasks[i].UpdatedAt = task.UpdatedAt
				break
			}
		}

		km.navState.MarkListDirty(taskStatus)
		km.updateTaskLists()
	}

	return nil
}

func (km *KahnModel) DeleteTask(id string) error {
	if err := km.taskService.DeleteTask(id); err != nil {
		return err
	}

	activeProj := km.GetActiveProject()
	if activeProj != nil {
		// Find task status before deletion for dirty flag
		var taskStatus domain.Status
		for _, t := range activeProj.Tasks {
			if t.ID == id {
				taskStatus = t.Status
				break
			}
		}
		activeProj.RemoveTask(id)

		km.navState.MarkListDirty(taskStatus)
		km.updateTaskLists()
	}

	return nil
}

func (km *KahnModel) MoveTaskToNextStatus(id string) error {
	activeProj := km.GetActiveProject()
	if activeProj == nil {
		return nil
	}

	// Find current status before movement for dirty flags
	var oldStatus domain.Status
	for _, t := range activeProj.Tasks {
		if t.ID == id {
			oldStatus = t.Status
			break
		}
	}

	task, err := km.taskService.MoveTaskToNextStatus(id)
	if err != nil {
		return err
	}

	activeProj.UpdateTaskStatus(id, task.Status)
	km.navState.MarkListDirty(oldStatus)
	km.navState.MarkListDirty(task.Status)
	km.updateTaskLists()
	return nil
}

func (km *KahnModel) MoveTaskToPreviousStatus(id string) error {
	activeProj := km.GetActiveProject()
	if activeProj == nil {
		return nil
	}

	// Find current status before movement for dirty flags
	var oldStatus domain.Status
	for _, t := range activeProj.Tasks {
		if t.ID == id {
			oldStatus = t.Status
			break
		}
	}

	task, err := km.taskService.MoveTaskToPreviousStatus(id)
	if err != nil {
		return err
	}

	activeProj.UpdateTaskStatus(id, task.Status)
	km.navState.MarkListDirty(oldStatus)
	km.navState.MarkListDirty(task.Status)
	km.updateTaskLists()
	return nil
}

func (km *KahnModel) GetSelectedTask() (input.TaskInterface, bool) {
	selectedItem := km.navState.GetActiveList().SelectedItem()
	if selectedItem == nil {
		return nil, false
	}

	taskWrapper, ok := selectedItem.(styles.TaskWithTitle)
	if !ok {
		return nil, false
	}

	return &domain.TaskWrapper{Task: taskWrapper.Task}, true
}

func (km *KahnModel) GetProjects() []input.ProjectInterface {
	var projects []input.ProjectInterface
	for i := range km.Projects {
		projects = append(projects, &domain.ProjectWrapper{Project: &km.Projects[i]})
	}
	return projects
}

func (km *KahnModel) CreateProject(name, description string) error {
	newProject, err := km.projectService.CreateProject(name, description)
	if err != nil {
		return err
	}

	km.Projects = append(km.Projects, *newProject)
	km.ActiveProjectID = newProject.ID
	km.updateTaskLists()

	return nil
}

func (km *KahnModel) DeleteProject(id string) error {
	if err := km.projectService.DeleteProject(id); err != nil {
		return err
	}

	if len(km.Projects) == 1 {
		km.Projects = []domain.Project{}
		km.ActiveProjectID = ""
		km.navState.Tasks[domain.NotStarted].SetItems([]list.Item{})
		km.navState.Tasks[domain.InProgress].SetItems([]list.Item{})
		km.navState.Tasks[domain.Done].SetItems([]list.Item{})
	} else {
		var newProjects []domain.Project
		var wasActiveProject bool
		for _, proj := range km.Projects {
			if proj.ID != id {
				newProjects = append(newProjects, proj)
			} else {
				wasActiveProject = (proj.ID == km.ActiveProjectID)
			}
		}
		km.Projects = newProjects

		if wasActiveProject && len(km.Projects) > 0 {
			km.ActiveProjectID = km.Projects[0].ID
			km.updateTaskLists()
		}
	}

	return nil
}

func (km *KahnModel) SwitchToProject(id string) error {
	km.ActiveProjectID = id
	km.updateTaskLists()
	return nil
}

func (km *KahnModel) GetSelectedProjectIndex() int {
	for i, proj := range km.Projects {
		if proj.ID == km.ActiveProjectID {
			return i
		}
	}
	return 0
}

func (km *KahnModel) GetFormError() string {
	errorMsg, _ := km.formState.GetError()
	return errorMsg
}

func (km *KahnModel) GetFormErrorField() string {
	_, errorField := km.formState.GetError()
	return errorField
}

func (km *KahnModel) ClearFormError() {
	km.formState.ClearError()
}

func (km *KahnModel) GetActiveInputComponents() *input.InputComponents {
	return km.formState.GetActiveInputComponents()
}

func (km *KahnModel) GetActiveFormType() input.FormType {
	return km.formState.GetActiveFormType()
}

func (km *KahnModel) SubmitCurrentForm() error {
	isValid, errorField, errorMsg := km.formState.ValidateForSubmit()
	if !isValid {
		km.formState.SetError(errorMsg, errorField)
		return fmt.Errorf("validation failed: %s", errorMsg)
	}

	km.formState.ClearError()
	name, desc, taskType, priority := km.formState.GetFormData()

	switch km.formState.GetActiveFormType() {
	case input.TaskCreateForm:
		newTask, err := km.taskService.CreateTask(name, desc, km.ActiveProjectID, taskType, priority)
		if err == nil {
			activeProj := km.GetActiveProject()
			if activeProj != nil {
				activeProj.Tasks = append(activeProj.Tasks, *newTask)
				km.updateTaskLists()
			}
		}
	case input.TaskEditForm:
		taskID := km.formState.GetTaskID()
		err := km.UpdateTask(taskID, name, desc, priority, taskType)
		return err
	case input.ProjectCreateForm:
		newProject, err := km.projectService.CreateProject(name, desc)
		if err == nil {
			km.Projects = append(km.Projects, *newProject)
			km.ActiveProjectID = newProject.ID
			km.updateTaskLists()
		}
		return err
	}
	return nil
}

func (km *KahnModel) CancelCurrentForm() {
	km.formState.HideForm()
}

func (km *KahnModel) ShowTaskForm() {
	km.formState.ShowTaskForm()
}

func (km *KahnModel) ShowTaskEditForm(taskID string, name, description string, priority domain.Priority, taskType domain.TaskType) {
	km.formState.ShowTaskEditForm(taskID, name, description, priority, taskType)
}

func (km *KahnModel) ShowProjectForm() {
	km.formState.ShowProjectForm()
}

func (km *KahnModel) ShowProjectSwitcher() {
	km.navState.ShowProjectSwitch()
}

func (km *KahnModel) ShowTaskDeleteConfirm(taskID string) {
	km.confirmState.ShowTaskDeleteConfirm(taskID)
}

func (km *KahnModel) HideAllForms() {
	km.formState.HideForm()
	km.navState.HideProjectSwitch()
	km.confirmState.HideAllConfirmations()
}

func (km *KahnModel) NextList() {
	km.navState.NextList()
}

func (km *KahnModel) PrevList() {
	km.navState.PrevList()
}

func (km *KahnModel) updateTaskLists() {

	if km.navState.dirtyFlags != nil && len(km.navState.dirtyFlags) > 0 {
		km.navState.UpdateDirtyLists(km.GetActiveProject(), km.taskService)
	} else {
		km.navState.UpdateTaskLists(km.GetActiveProject(), km.taskService)
	}
}

func (km *KahnModel) executeTaskDeletion() tea.Model {
	taskToDelete := km.confirmState.GetTaskToDelete()
	if taskToDelete == "" {
		km.confirmState.ClearTaskDelete()
		return km
	}

	if err := km.taskService.DeleteTask(taskToDelete); err != nil {
		km.confirmState.SetError("Failed to delete task: " + err.Error())
		return km
	}

	activeProj := km.GetActiveProject()
	if activeProj != nil {
		activeProj.RemoveTask(taskToDelete)
		km.updateTaskLists()
	}

	km.confirmState.ClearTaskDelete()
	return km
}

func (km *KahnModel) executeProjectDeletion() tea.Model {
	projectToDelete := km.confirmState.GetProjectToDelete()
	if projectToDelete == "" {
		km.confirmState.ClearProjectDelete()
		return km
	}

	if err := km.projectService.DeleteProject(projectToDelete); err != nil {
		km.confirmState.SetError("Failed to delete project: " + err.Error())
		return km
	}

	if len(km.Projects) == 1 {
		km.Projects = []domain.Project{}
		km.ActiveProjectID = ""

		km.navState.Tasks[domain.NotStarted].SetItems([]list.Item{})
		km.navState.Tasks[domain.InProgress].SetItems([]list.Item{})
		km.navState.Tasks[domain.Done].SetItems([]list.Item{})
	} else {
		var newProjects []domain.Project
		var wasActiveProject bool
		for _, proj := range km.Projects {
			if proj.ID != projectToDelete {
				newProjects = append(newProjects, proj)
			} else {
				wasActiveProject = (proj.ID == km.ActiveProjectID)
			}
		}
		km.Projects = newProjects

		if wasActiveProject && len(km.Projects) > 0 {
			km.ActiveProjectID = km.Projects[0].ID
			km.updateTaskLists()
		}
	}

	km.confirmState.ClearProjectDelete()
	return km
}

func (km *KahnModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if km.formState.IsShowingForm() {
			return km.handleFormInput(msg)
		}
		if km.navState.IsShowingProjectSwitch() || km.confirmState.IsShowingProjectDeleteConfirm() {
			return km.handleProjectSwitch(msg)
		}
		if km.confirmState.IsShowingTaskDeleteConfirm() {
			return km.handleTaskDeleteConfirm(msg)
		}
		return km.handleNormalMode(msg)
	case tea.WindowSizeMsg:
		return km.handleResize(msg)
	}
	return km, nil
}

func (km *KahnModel) GetTaskToDelete() string {
	return km.confirmState.GetTaskToDelete()
}

func (km *KahnModel) GetProjectToDelete() string {
	return km.confirmState.GetProjectToDelete()
}

func (km *KahnModel) IsShowingTaskDeleteConfirm() bool {
	return km.confirmState.IsShowingTaskDeleteConfirm()
}

func (km *KahnModel) IsShowingProjectDeleteConfirm() bool {
	return km.confirmState.IsShowingProjectDeleteConfirm()
}

func (km *KahnModel) GetTaskItems(status domain.Status) []list.Item {
	return km.navState.Tasks[status].Items()
}

func (km *KahnModel) GetActiveListIndex() domain.Status {
	return km.navState.GetActiveListIndex()
}

func (km *KahnModel) IsShowingForm() bool {
	return km.formState.IsShowingForm()
}

func (km *KahnModel) IsShowingProjectSwitch() bool {
	return km.navState.IsShowingProjectSwitch()
}

func NewKahnModel(database *database.Database, version string) *KahnModel {
	// Create delegates for different list states
	activeDelegate := styles.NewActiveListDelegate()
	inactiveDelegate := styles.NewInactiveListDelegate()

	// Create lists with appropriate delegates - first list is active by default
	activeList := list.New([]list.Item{}, activeDelegate, 100, 0)
	activeList.SetShowHelp(false)

	inactiveList := list.New([]list.Item{}, inactiveDelegate, 100, 0)
	inactiveList.SetShowHelp(false)

	taskLists := []list.Model{activeList, inactiveList, inactiveList}

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
		notStartedTasks, err := taskService.GetTasksByStatus(projects[0].ID, domain.NotStarted)
		if err != nil {
			projects[0].Tasks = []domain.Task{}
		} else {
			projects[0].Tasks = notStartedTasks
		}
		taskLists[domain.NotStarted].SetItems(convertTasksToListItems(projects[0].GetTasksByStatus(domain.NotStarted)))
	}

	taskLists[domain.InProgress].Title = domain.InProgress.ToString()
	if len(projects) > 0 {
		inProgressTasks, err := taskService.GetTasksByStatus(projects[0].ID, domain.InProgress)
		if err != nil {
			projects[0].Tasks = []domain.Task{}
		} else {
			projects[0].Tasks = inProgressTasks
		}
		taskLists[domain.InProgress].SetItems(convertTasksToListItems(projects[0].GetTasksByStatus(domain.InProgress)))
	}

	taskLists[domain.Done].Title = domain.Done.ToString()
	if len(projects) > 0 {
		doneTasks, err := taskService.GetTasksByStatus(projects[0].ID, domain.Done)
		if err != nil {
			projects[0].Tasks = []domain.Task{}
		} else {
			projects[0].Tasks = doneTasks
		}
		taskLists[domain.Done].SetItems(convertTasksToListItems(projects[0].GetTasksByStatus(domain.Done)))
	}

	// Update selection states after initialization
	// NotStarted is the active list by default, others are inactive
	taskLists[domain.NotStarted].SetItems(styles.UpdateTaskSelection(taskLists[domain.NotStarted].Items(), taskLists[domain.NotStarted].Index(), true))
	taskLists[domain.InProgress].SetItems(styles.UpdateTaskSelection(taskLists[domain.InProgress].Items(), taskLists[domain.InProgress].Index(), false))
	taskLists[domain.Done].SetItems(styles.UpdateTaskSelection(taskLists[domain.Done].Items(), taskLists[domain.Done].Index(), false))

	// Apply title styles to all lists based on active list index
	styles.ApplyFocusedTitleStyles(taskLists, domain.NotStarted) // Default to NotStarted as initial active list

	// Initialize state management
	formState := NewFormState(taskInputComponents, projectInputComponents)
	confirmState := NewConfirmationState()
	navState := NewNavigationState(taskLists)

	return &KahnModel{
		Projects:        projects,
		ActiveProjectID: activeProjectID,
		width:           80,
		height:          24,
		database:        database,
		inputHandler:    input.NewHandler(),
		taskService:     taskService,
		projectService:  projectService,
		board:           components.NewBoard(),
		projectSwitcher: components.NewProjectSwitcher(),
		version:         version,
		formState:       formState,
		confirmState:    confirmState,
		navState:        navState,
	}
}
