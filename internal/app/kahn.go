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
	width           int
	height          int
	database        *database.Database
	taskService     *services.TaskService
	projectService  *services.ProjectService
	board           *components.Board
	projectSwitcher *components.ProjectSwitcher
	version         string

	// State managers
	uiStateManager *UIStateManager
	projectManager *ProjectManager
	navState       *NavigationState
	searchState    *SearchState
}

func (km KahnModel) Init() tea.Cmd {
	return nil
}

// renderForm renders the form view (task or project forms)
func (km *KahnModel) renderForm() string {
	formState := km.uiStateManager.FormState()
	comps := formState.GetActiveInputComponents()
	formError, formErrorField := formState.GetError()
	return comps.Render(formError, formErrorField, km.width, km.height)
}

// renderProjectSwitcher renders the project switcher with delete confirmation
func (km *KahnModel) renderProjectSwitcher() string {
	confirmState := km.uiStateManager.ConfirmationState()
	return km.projectSwitcher.RenderSwitcherWithError(
		km.projectManager.GetProjectsAsDomain(),
		km.projectManager.GetActiveProjectID(),
		confirmState.IsShowingProjectDeleteConfirm(),
		confirmState.GetProjectToDelete(),
		confirmState.GetProjectError(),
		km.width,
		km.height,
	)
}

// renderTaskDeleteConfirm renders the task deletion confirmation dialog
func (km *KahnModel) renderTaskDeleteConfirm() string {
	confirmState := km.uiStateManager.ConfirmationState()
	taskToDelete := km.findTaskForDeletion()
	return km.board.GetRenderer().RenderTaskDeleteConfirmWithError(taskToDelete, confirmState.GetTaskError(), km.width, km.height)
}

// renderProjectDeleteConfirm renders the project deletion confirmation dialog
func (km *KahnModel) renderProjectDeleteConfirm() string {
	confirmState := km.uiStateManager.ConfirmationState()
	return km.board.GetRenderer().RenderTaskDeleteConfirmWithError(nil, confirmState.GetProjectError(), km.width, km.height)
}

// renderNoProjects renders the no projects state
func (km *KahnModel) renderNoProjects() string {
	return km.board.GetRenderer().RenderNoProjectsBoard(km.width, km.height)
}

// renderBoard renders the main board view with task lists
func (km *KahnModel) renderBoard() string {
	if !km.projectManager.HasProjects() {
		return km.renderNoProjects()
	}

	activeProj := km.projectManager.GetActiveProject()
	if activeProj == nil {
		return ""
	}

	taskLists := km.getTaskListsForBoard()
	navState := km.navState
	return km.board.GetRenderer().RenderBoard(
		activeProj,
		taskLists,
		navState.GetActiveListIndex(),
		km.width,
		km.version,
		km.searchState.IsActive(),
		km.searchState.GetQuery(),
		km.searchState.GetMatchCount(),
	)
}

// View renders the appropriate view based on current UI state
func (km KahnModel) View() string {
	switch km.uiStateManager.GetCurrentViewState() {
	case FormView:
		return km.renderForm()
	case ProjectSwitchView:
		return km.renderProjectSwitcher()
	case TaskDeleteConfirmView:
		return km.renderTaskDeleteConfirm()
	case ProjectDeleteConfirmView:
		return km.renderProjectDeleteConfirm()
	case NoProjectsView:
		return km.renderNoProjects()
	default: // BoardView
		return km.renderBoard()
	}
}

func (km *KahnModel) GetActiveProject() *domain.Project {
	return km.projectManager.GetActiveProject()
}

// RefreshTasksWithSearch updates task lists applying the current search filter if active,
// or shows all tasks if search is inactive. Updates the match count when search is active.
func (km *KahnModel) RefreshTasksWithSearch() {
	activeProj := km.GetActiveProject()
	if activeProj == nil {
		return
	}

	if km.searchState.IsActive() {
		km.navState.UpdateTaskListsWithSearch(
			activeProj,
			km.taskService,
			km.searchState.GetQuery(),
		)

		// Update match count
		allTasks := activeProj.Tasks
		matchCount := domain.CountSearchMatches(allTasks, km.searchState.GetQuery())
		km.searchState.UpdateMatchCount(matchCount)
	} else {
		km.navState.UpdateTaskLists(activeProj, km.taskService)
	}
}

func (km *KahnModel) GetActiveProjectID() string {
	return km.projectManager.GetActiveProjectID()
}

func (km *KahnModel) CreateTask(name, description string) error {
	return km.CreateTaskWithPriority(name, description, domain.Low)
}

func (km *KahnModel) CreateTaskWithPriority(name, description string, priority domain.Priority) error {
	activeProj := km.GetActiveProject()
	if activeProj == nil {
		return nil
	}

	newTask, err := km.taskService.CreateTask(name, description, km.GetActiveProjectID(), domain.RegularTask, priority, nil)
	if err != nil {
		return err
	}

	activeProj.AddTask(*newTask)

	km.navState.MarkListDirty(domain.NotStarted)
	km.RefreshTasksWithSearch()

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
		km.RefreshTasksWithSearch()
	}

	return nil
}

func (km *KahnModel) DeleteTask(id string) error {
	if err := km.taskService.DeleteTask(id); err != nil {
		return err
	}

	activeProj := km.GetActiveProject()
	if activeProj != nil {
		activeProj.RemoveTask(id)

		// Refresh all columns to update visual indicators for unblocked tasks
		km.navState.MarkAllListsDirty()
		km.RefreshTasksWithSearch()
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
	var movingToComplete bool
	for _, t := range activeProj.Tasks {
		if t.ID == id {
			oldStatus = t.Status
			movingToComplete = (t.Status == domain.InProgress)
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

	// Refresh all columns to update visual indicators for unblocked tasks
	if movingToComplete {
		km.navState.MarkAllListsDirty()
	}

	km.RefreshTasksWithSearch()
	return nil
}

func (km *KahnModel) MoveTaskToPreviousStatus(id string) error {
	activeProj := km.GetActiveProject()
	if activeProj == nil {
		return nil
	}

	// Find current status before movement for dirty flags
	var oldStatus domain.Status
	var movingToComplete bool
	for _, t := range activeProj.Tasks {
		if t.ID == id {
			oldStatus = t.Status
			movingToComplete = (t.Status == domain.NotStarted)
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

	// Refresh all columns to update visual indicators for unblocked tasks
	if movingToComplete {
		km.navState.MarkAllListsDirty()
	}

	km.RefreshTasksWithSearch()
	return nil
}

// GetSelectedTask returns the currently selected task for internal use
func (km *KahnModel) getSelectedTask() (*styles.TaskWithTitle, bool) {
	selectedItem := km.navState.GetActiveList().SelectedItem()
	if selectedItem == nil {
		return nil, false
	}

	taskWrapper, ok := selectedItem.(styles.TaskWithTitle)
	if !ok {
		return nil, false
	}

	return &taskWrapper, true
}

func (km *KahnModel) CreateProject(name, description string) error {
	return km.projectManager.CreateProject(name, description)
}

func (km *KahnModel) DeleteProject(id string) error {
	return km.projectManager.DeleteProject(id)
}

func (km *KahnModel) SwitchToProject(id string) error {
	return km.projectManager.SwitchToProject(id)
}

func (km *KahnModel) GetFormError() string {
	errorMsg, _ := km.uiStateManager.FormState().GetError()
	return errorMsg
}

func (km *KahnModel) GetFormErrorField() string {
	_, errorField := km.uiStateManager.FormState().GetError()
	return errorField
}

func (km *KahnModel) ClearFormError() {
	km.uiStateManager.FormState().ClearError()
}

func (km *KahnModel) GetActiveInputComponents() *input.InputComponents {
	return km.uiStateManager.FormState().GetActiveInputComponents()
}

func (km *KahnModel) GetActiveFormType() input.FormType {
	return km.uiStateManager.FormState().GetActiveFormType()
}

func (km *KahnModel) SubmitCurrentForm() error {
	formState := km.uiStateManager.FormState()
	isValid, errorField, errorMsg := formState.ValidateForSubmit()
	if !isValid {
		formState.SetError(errorMsg, errorField)
		return fmt.Errorf("validation failed: %s", errorMsg)
	}

	formState.ClearError()
	name, desc, taskType, priority, blockedByIntID := formState.GetFormData()

	switch formState.GetActiveFormType() {
	case input.TaskCreateForm:
		newTask, err := km.taskService.CreateTask(name, desc, km.GetActiveProjectID(), taskType, priority, blockedByIntID)
		if err == nil {
			activeProj := km.GetActiveProject()
			if activeProj != nil {
				activeProj.Tasks = append(activeProj.Tasks, *newTask)
				km.RefreshTasksWithSearch()
			}
		}
		return err
	case input.TaskEditForm:
		taskID := formState.GetTaskID()
		// Update basic task fields
		err := km.UpdateTask(taskID, name, desc, priority, taskType)
		if err != nil {
			return err
		}
		// Update BlockedBy field separately
		_, err = km.taskService.SetTaskBlockedBy(taskID, blockedByIntID)
		if err != nil {
			return err
		}
		// Update the task in the active project and refresh display
		activeProj := km.GetActiveProject()
		if activeProj != nil {
			var taskStatus domain.Status
			for i, t := range activeProj.Tasks {
				if t.ID == taskID {
					taskStatus = t.Status // Save status for dirty flag
					activeProj.Tasks[i].BlockedBy = blockedByIntID
					break
				}
			}
			// Mark list dirty and refresh display to show updated blocked status
			km.navState.MarkListDirty(taskStatus)
			km.RefreshTasksWithSearch()
		}
		return nil
	case input.ProjectCreateForm:
		return km.projectManager.CreateProject(name, desc)
	}
	return nil
}

func (km *KahnModel) CancelCurrentForm() {
	km.uiStateManager.HideAllStates()
}

func (km *KahnModel) ShowTaskForm() {
	// Get available tasks for BlockedBy field
	availableTasks := km.getAvailableBlockerTasks("")
	km.uiStateManager.ShowTaskForm(availableTasks)
}

func (km *KahnModel) ShowTaskEditForm(taskID string, name, description string, priority domain.Priority, taskType domain.TaskType, blockedByIntID *int) {
	// Get available tasks for BlockedBy field (exclude current task)
	availableTasks := km.getAvailableBlockerTasks(taskID)
	km.uiStateManager.ShowTaskEditForm(taskID, name, description, priority, taskType, blockedByIntID, availableTasks)
}

func (km *KahnModel) ShowProjectForm() {
	km.uiStateManager.ShowProjectForm()
}

func (km *KahnModel) ShowProjectSwitcher() {
	km.uiStateManager.ShowProjectSwitcher()
}

func (km *KahnModel) ShowTaskDeleteConfirm(taskID string) {
	km.uiStateManager.ShowTaskDeleteConfirm(taskID)
}

func (km *KahnModel) HideAllForms() {
	km.uiStateManager.HideAllStates()
}

func (km *KahnModel) NextList() {
	km.navState.NextList()
}

func (km *KahnModel) PrevList() {
	km.navState.PrevList()
}

func (km *KahnModel) updateTaskLists() {
	km.navState.UpdateTaskListsConditional(km.GetActiveProject(), km.taskService)
}

// findTaskForDeletion finds the task to be deleted either from active project or selected item
func (km *KahnModel) findTaskForDeletion() *domain.Task {
	confirmState := km.uiStateManager.ConfirmationState()
	taskID := confirmState.GetTaskToDelete()

	// First try to find in active project
	activeProj := km.projectManager.GetActiveProject()
	if activeProj != nil {
		for _, task := range activeProj.Tasks {
			if task.ID == taskID {
				return &task
			}
		}
	}

	// Fallback to selected item
	if selectedItem := km.navState.GetActiveList().SelectedItem(); selectedItem != nil {
		if taskWrapper, ok := selectedItem.(styles.TaskWithTitle); ok {
			return &taskWrapper.Task
		}
	}

	return nil
}

// getTaskListsForBoard builds the task lists array needed for board rendering
func (km *KahnModel) getTaskListsForBoard() [3]list.Model {
	navState := km.navState
	return [3]list.Model{
		navState.Tasks[domain.NotStarted],
		navState.Tasks[domain.InProgress],
		navState.Tasks[domain.Done],
	}
}

// getAvailableBlockerTasks returns tasks that can block another task
// Filters to only NotStarted and InProgress tasks
// If excludeTaskID is provided, that task is excluded from the list
func (km *KahnModel) getAvailableBlockerTasks(excludeTaskID string) []domain.Task {
	activeProj := km.GetActiveProject()
	if activeProj == nil {
		return []domain.Task{}
	}

	var availableTasks []domain.Task
	for _, task := range activeProj.Tasks {
		// Skip the task being edited
		if task.ID == excludeTaskID {
			continue
		}
		// Only include NotStarted and InProgress tasks
		if task.Status == domain.NotStarted || task.Status == domain.InProgress {
			availableTasks = append(availableTasks, task)
		}
	}
	return availableTasks
}

func (km *KahnModel) executeTaskDeletion() tea.Model {
	confirmState := km.uiStateManager.ConfirmationState()
	taskToDelete := confirmState.GetTaskToDelete()
	if taskToDelete == "" {
		confirmState.ClearTaskDelete()
		return km
	}

	if err := km.taskService.DeleteTask(taskToDelete); err != nil {
		confirmState.SetTaskError("Failed to delete task: " + err.Error())
		return km
	}

	activeProj := km.GetActiveProject()
	if activeProj != nil {
		activeProj.RemoveTask(taskToDelete)
		km.RefreshTasksWithSearch()
	}

	confirmState.ClearTaskDelete()
	return km
}

func (km *KahnModel) executeProjectDeletion() tea.Model {
	confirmState := km.uiStateManager.ConfirmationState()
	projectToDelete := confirmState.GetProjectToDelete()
	if projectToDelete == "" {
		confirmState.ClearProjectDelete()
		return km
	}

	// The project manager handles the deletion logic
	err := km.projectManager.DeleteProject(projectToDelete)
	if err != nil {
		confirmState.SetProjectError("Failed to delete project: " + err.Error())
	}

	confirmState.ClearProjectDelete()
	return km
}

func (km *KahnModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Check if in search mode first
		if km.searchState.IsActive() {
			return km.handleSearchInput(msg)
		}
		if km.uiStateManager.FormState().IsShowingForm() {
			return km.handleFormInput(msg)
		}
		if km.navState.IsShowingProjectSwitch() || km.uiStateManager.ConfirmationState().IsShowingProjectDeleteConfirm() {
			return km.handleProjectSwitch(msg)
		}
		if km.uiStateManager.ConfirmationState().IsShowingTaskDeleteConfirm() {
			return km.handleTaskDeleteConfirm(msg)
		}
		return km.handleNormalMode(msg)
	case tea.WindowSizeMsg:
		return km.handleResize(msg)
	}
	return km, nil
}

func (km *KahnModel) GetTaskToDelete() string {
	return km.uiStateManager.ConfirmationState().GetTaskToDelete()
}

func (km *KahnModel) GetProjectToDelete() string {
	return km.uiStateManager.ConfirmationState().GetProjectToDelete()
}

func (km *KahnModel) IsShowingTaskDeleteConfirm() bool {
	return km.uiStateManager.ConfirmationState().IsShowingTaskDeleteConfirm()
}

func (km *KahnModel) IsShowingProjectDeleteConfirm() bool {
	return km.uiStateManager.ConfirmationState().IsShowingProjectDeleteConfirm()
}

func (km *KahnModel) GetTaskItems(status domain.Status) []list.Item {
	return km.navState.GetTaskItems(status)
}

func (km *KahnModel) GetActiveListIndex() domain.Status {
	return km.navState.GetActiveListIndex()
}

func (km *KahnModel) IsShowingForm() bool {
	return km.uiStateManager.FormState().IsShowingForm()
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

	taskComps := input.NewInputComponents()
	taskInputComponents := &taskComps
	projectComps := input.NewInputComponents()
	projectInputComponents := &projectComps

	// Create repositories
	taskRepo := repo.NewSQLiteTaskRepository(database.GetDB())
	projectRepo := repo.NewSQLiteProjectRepository(database.GetDB())

	// Create services
	taskService := services.NewTaskService(taskRepo, projectRepo)
	projectService := services.NewProjectService(projectRepo, taskRepo)

	// Create state management components
	formState := NewFormState(taskInputComponents, projectInputComponents)
	confirmState := NewConfirmationState()
	navState := NewNavigationState(taskLists)
	searchState := NewSearchState()

	// Create managers
	projectManager := NewProjectManager(projectService, taskService, navState)
	uiStateManager := NewUIStateManager(formState, confirmState, navState)

	// Initialize projects through project manager
	projectManager.InitializeProjects()

	// Apply list titles and styles
	taskLists[domain.NotStarted].Title = domain.NotStarted.ToString()
	taskLists[domain.InProgress].Title = domain.InProgress.ToString()
	taskLists[domain.Done].Title = domain.Done.ToString()

	// Update selection states after initialization
	// NotStarted is the active list by default, others are inactive
	taskLists[domain.NotStarted].SetItems(styles.UpdateTaskSelection(taskLists[domain.NotStarted].Items(), taskLists[domain.NotStarted].Index(), true))
	taskLists[domain.InProgress].SetItems(styles.UpdateTaskSelection(taskLists[domain.InProgress].Items(), taskLists[domain.InProgress].Index(), false))
	taskLists[domain.Done].SetItems(styles.UpdateTaskSelection(taskLists[domain.Done].Items(), taskLists[domain.Done].Index(), false))

	// Apply title styles to all lists based on active list index
	styles.ApplyFocusedTitleStyles(taskLists, domain.NotStarted) // Default to NotStarted as initial active list

	return &KahnModel{
		width:           80,
		height:          24,
		database:        database,
		taskService:     taskService,
		projectService:  projectService,
		board:           components.NewBoard(),
		projectSwitcher: components.NewProjectSwitcher(),
		version:         version,
		uiStateManager:  uiStateManager,
		projectManager:  projectManager,
		navState:        navState,
		searchState:     searchState,
	}
}
