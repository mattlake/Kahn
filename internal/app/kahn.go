package app

import (
	"fmt"
	"strings"

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
	Projects                 []domain.Project
	ActiveProjectID          string
	Tasks                    []list.Model
	activeListIndex          domain.Status
	showForm                 bool
	showProjectSwitch        bool
	showProjectDeleteConfirm bool
	projectToDelete          string
	showTaskDeleteConfirm    bool
	taskToDelete             string

	TaskService *services.TaskService

	taskInputComponents    *input.InputComponents
	projectInputComponents *input.InputComponents
	activeFormType         input.FormType
	formError              string // validation error message
	formErrorField         string // which field has error
	width                  int
	height                 int
	database               *database.Database
	inputHandler           *input.Handler
	taskService            *services.TaskService
	projectService         *services.ProjectService
	board                  *components.Board
	projectSwitcher        *components.ProjectSwitcher
}

func (km KahnModel) Init() tea.Cmd {
	return nil
}

func (km KahnModel) View() string {
	if km.showForm {
		comps := km.GetActiveInputComponents()
		return comps.Render(km.formError, km.formErrorField, km.width, km.height)
	}
	if km.showProjectSwitch {
		return km.projectSwitcher.RenderSwitcher(km.Projects, km.ActiveProjectID, km.showProjectDeleteConfirm, km.projectToDelete, km.width, km.height)
	}
	if km.showTaskDeleteConfirm {
		var taskToDelete *domain.Task
		activeProj := km.GetActiveProject()
		if activeProj != nil {
			for _, task := range activeProj.Tasks {
				if task.ID == km.taskToDelete {
					taskToDelete = &task
					break
				}
			}
		}

		if taskToDelete == nil {
			if selectedItem := km.Tasks[km.activeListIndex].SelectedItem(); selectedItem != nil {
				if task, ok := selectedItem.(domain.Task); ok {
					taskToDelete = &task
				}
			}
		}

		return km.board.GetRenderer().RenderTaskDeleteConfirm(taskToDelete, km.width, km.height)
	}

	if len(km.Projects) == 0 {
		return km.board.GetRenderer().RenderNoProjectsBoard(km.width, km.height)
	}

	activeProj := km.GetActiveProject()
	if activeProj != nil {
		var taskLists [3]list.Model
		taskLists[domain.NotStarted] = km.Tasks[domain.NotStarted]
		taskLists[domain.InProgress] = km.Tasks[domain.InProgress]
		taskLists[domain.Done] = km.Tasks[domain.Done]

		return km.board.GetRenderer().RenderBoard(activeProj, taskLists, km.activeListIndex, km.width)
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
	activeProj := km.GetActiveProject()
	if activeProj == nil {
		return nil
	}

	newTask, err := km.taskService.CreateTask(name, description, km.ActiveProjectID)
	if err != nil {
		return err
	}

	activeProj.AddTask(*newTask)
	km.updateTaskLists()

	return nil
}

func (km *KahnModel) UpdateTask(id, name, description string) error {
	task, err := km.taskService.UpdateTask(id, name, description)
	if err != nil {
		return err
	}

	activeProj := km.GetActiveProject()
	if activeProj != nil {
		for i, t := range activeProj.Tasks {
			if t.ID == id {
				activeProj.Tasks[i].Name = task.Name
				activeProj.Tasks[i].Desc = task.Desc
				activeProj.Tasks[i].UpdatedAt = task.UpdatedAt
				break
			}
		}
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
		activeProj.RemoveTask(id)
		km.updateTaskLists()
	}

	return nil
}

func (km *KahnModel) MoveTaskToNextStatus(id string) error {
	activeProj := km.GetActiveProject()
	if activeProj == nil {
		return nil
	}

	task, err := km.taskService.MoveTaskToNextStatus(id)
	if err != nil {
		return err
	}

	activeProj.UpdateTaskStatus(id, task.Status)
	km.updateTaskLists()
	return nil
}

func (km *KahnModel) MoveTaskToPreviousStatus(id string) error {
	activeProj := km.GetActiveProject()
	if activeProj == nil {
		return nil
	}

	task, err := km.taskService.MoveTaskToPreviousStatus(id)
	if err != nil {
		return err
	}

	activeProj.UpdateTaskStatus(id, task.Status)
	km.updateTaskLists()
	return nil
}

func (km *KahnModel) GetSelectedTask() (input.TaskInterface, bool) {
	selectedItem := km.Tasks[km.activeListIndex].SelectedItem()
	if selectedItem == nil {
		return nil, false
	}

	task, ok := selectedItem.(domain.Task)
	if !ok {
		return nil, false
	}

	return &domain.TaskWrapper{Task: task}, true
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
		km.Tasks[domain.NotStarted].SetItems([]list.Item{})
		km.Tasks[domain.InProgress].SetItems([]list.Item{})
		km.Tasks[domain.Done].SetItems([]list.Item{})
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

func (km *KahnModel) setFormError(message string, field string) {
	km.formError = message
	km.formErrorField = field
}

func (km *KahnModel) ClearFormError() {
	km.formError = ""
	km.formErrorField = ""
}

func (km *KahnModel) GetFormError() string {
	return km.formError
}

func (km *KahnModel) GetFormErrorField() string {
	return km.formErrorField
}

func (km *KahnModel) GetActiveInputComponents() *input.InputComponents {
	if km.activeFormType == input.TaskCreateForm || km.activeFormType == input.TaskEditForm {
		return km.taskInputComponents
	}
	return km.projectInputComponents
}

func (km *KahnModel) GetActiveFormType() input.FormType {
	return km.activeFormType
}

func (km *KahnModel) SubmitCurrentForm() error {
	comps := km.GetActiveInputComponents()

	isValid, errorField, errorMsg := comps.ValidateForSubmit()
	if !isValid {
		km.setFormError(errorMsg, errorField)
		return fmt.Errorf("validation failed: %s", errorMsg)
	}

	km.ClearFormError()
	name := strings.TrimSpace(comps.NameInput.Value())
	desc := comps.DescInput.Value()

	switch km.activeFormType {
	case input.TaskCreateForm:
		newTask, err := km.taskService.CreateTask(name, desc, km.ActiveProjectID)
		if err == nil {
			activeProj := km.GetActiveProject()
			if activeProj != nil {
				activeProj.AddTask(*newTask)
				km.updateTaskLists()
			}
		}
		return err
	case input.TaskEditForm:
		err := km.UpdateTask(comps.GetTaskID(), name, desc)
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
	km.showForm = false
	km.ClearFormError()

	km.taskInputComponents.Reset()
	km.projectInputComponents.Reset()
}

func (km *KahnModel) ShowTaskForm() {
	km.taskInputComponents.SetupForTaskCreate()
	km.activeFormType = input.TaskCreateForm
	km.showForm = true
	km.ClearFormError()
}

func (km *KahnModel) ShowTaskEditForm(taskID string, name, description string) {
	km.taskInputComponents.SetupForTaskEdit(taskID, name, description)
	km.activeFormType = input.TaskEditForm
	km.showForm = true
	km.ClearFormError()
}

func (km *KahnModel) ShowProjectForm() {
	km.projectInputComponents.SetupForProjectCreate()
	km.activeFormType = input.ProjectCreateForm
	km.showForm = true
	km.ClearFormError()
}

func (km *KahnModel) ShowProjectSwitcher() {
	km.showProjectSwitch = true
}

func (km *KahnModel) HideAllForms() {
	km.showForm = false
	km.showProjectSwitch = false
	km.showTaskDeleteConfirm = false
	km.showProjectDeleteConfirm = false
	km.ClearFormError()
	km.taskInputComponents.Reset()
	km.projectInputComponents.Reset()
	km.taskToDelete = ""
	km.projectToDelete = ""
}

func (km *KahnModel) NextList() {
	if km.activeListIndex == domain.Done {
		km.activeListIndex = domain.NotStarted
	} else {
		km.activeListIndex++
	}
}

func (km *KahnModel) PrevList() {
	if km.activeListIndex == domain.NotStarted {
		km.activeListIndex = domain.Done
	} else {
		km.activeListIndex--
	}
}

func (km *KahnModel) updateTaskLists() {
	activeProj := km.GetActiveProject()
	if activeProj == nil {
		return
	}

	km.Tasks[domain.NotStarted].SetItems(convertTasksToListItems(activeProj.GetTasksByStatus(domain.NotStarted)))
	km.Tasks[domain.InProgress].SetItems(convertTasksToListItems(activeProj.GetTasksByStatus(domain.InProgress)))
	km.Tasks[domain.Done].SetItems(convertTasksToListItems(activeProj.GetTasksByStatus(domain.Done)))
}

func (km *KahnModel) executeTaskDeletion() tea.Model {
	if km.taskToDelete == "" {
		km.showTaskDeleteConfirm = false
		return km
	}

	if err := km.taskService.DeleteTask(km.taskToDelete); err != nil {
		km.showTaskDeleteConfirm = false
		km.taskToDelete = ""
		return km
	}

	activeProj := km.GetActiveProject()
	if activeProj != nil {
		activeProj.RemoveTask(km.taskToDelete)
		km.updateTaskLists()
	}

	km.showTaskDeleteConfirm = false
	km.taskToDelete = ""

	return km
}

func (km *KahnModel) executeProjectDeletion() tea.Model {
	if km.projectToDelete == "" {
		km.showProjectDeleteConfirm = false
		return km
	}

	if err := km.projectService.DeleteProject(km.projectToDelete); err != nil {
		km.showProjectDeleteConfirm = false
		km.projectToDelete = ""
		return km
	}

	if len(km.Projects) == 1 {
		km.Projects = []domain.Project{}
		km.ActiveProjectID = ""

		km.Tasks[domain.NotStarted].SetItems([]list.Item{})
		km.Tasks[domain.InProgress].SetItems([]list.Item{})
		km.Tasks[domain.Done].SetItems([]list.Item{})
	} else {
		var newProjects []domain.Project
		var wasActiveProject bool
		for _, proj := range km.Projects {
			if proj.ID != km.projectToDelete {
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

	km.showProjectDeleteConfirm = false
	km.projectToDelete = ""

	return km
}

func (km *KahnModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if km.showForm {
			action := km.inputHandler.HandleKeyMsg(msg, km)
			if action.Handled {
				return km, action.Cmd
			}

			comps := km.GetActiveInputComponents()
			if comps.FocusedField == 0 {
				updatedName, cmd := comps.NameInput.Update(msg)
				comps.NameInput = updatedName
				return km, cmd
			} else {
				updatedDesc, cmd := comps.DescInput.Update(msg)
				comps.DescInput = updatedDesc
				return km, cmd
			}
		}
		if km.showProjectSwitch || km.showProjectDeleteConfirm {
			if km.showProjectSwitch && km.inputHandler.GetMode() != input.ProjectSwitchMode {
				km.inputHandler.SetMode(input.ProjectSwitchMode)
			}
			if km.showProjectDeleteConfirm && km.inputHandler.GetMode() != input.ProjectDeleteConfirmMode {
				km.inputHandler.SetMode(input.ProjectDeleteConfirmMode)
			}

			if km.showProjectDeleteConfirm {
				switch msg.String() {
				case "y", "Y":
					return km.executeProjectDeletion(), nil
				case "n", "N", "esc":
					km.showProjectDeleteConfirm = false
					km.projectToDelete = ""
					return km, nil
				}
				return km, nil
			}

			switch msg.String() {
			case "esc":
				km.showProjectSwitch = false
				return km, nil
			case "d":
				if len(km.Projects) > 0 {
					km.showProjectDeleteConfirm = true
					km.projectToDelete = km.ActiveProjectID
				}
				return km, nil
			case "n":
				km.showProjectSwitch = false
				km.ShowProjectForm()
				return km, nil
			case "j":
				for i, proj := range km.Projects {
					if proj.ID == km.ActiveProjectID {
						nextIndex := (i + 1) % len(km.Projects)
						km.ActiveProjectID = km.Projects[nextIndex].ID
						km.updateTaskLists()
						return km, nil
					}
				}
			case "k":
				for i, proj := range km.Projects {
					if proj.ID == km.ActiveProjectID {
						prevIndex := (i - 1 + len(km.Projects)) % len(km.Projects)
						km.ActiveProjectID = km.Projects[prevIndex].ID
						km.updateTaskLists()
						return km, nil
					}
				}
			case "enter":
				km.showProjectSwitch = false
				return km, nil
			default:
				if len(msg.String()) == 1 && msg.String()[0] >= '1' && msg.String()[0] <= '9' {
					index := int(msg.String()[0] - '1')
					if index < len(km.Projects) {
						km.ActiveProjectID = km.Projects[index].ID
						km.updateTaskLists()
						km.showProjectSwitch = false
					}
				}
				return km, nil
			}

		} else if km.showTaskDeleteConfirm {
			if km.inputHandler.GetMode() != input.TaskDeleteConfirmMode {
				km.inputHandler.SetMode(input.TaskDeleteConfirmMode)
			}

			switch msg.String() {
			case "y", "Y":
				return km.executeTaskDeletion(), nil
			case "n", "N", "esc":
				km.showTaskDeleteConfirm = false
				km.taskToDelete = ""
				return km, nil
			}
			return km, nil
		} else {
			if km.inputHandler.GetMode() != input.NormalMode {
				km.inputHandler.SetMode(input.NormalMode)
			}

			result := km.inputHandler.HandleKeyMsg(msg, km)

			if result.Handled && result.Cmd != nil {
				return km, result.Cmd
			}

			switch msg.String() {
			case "q":
				return km, tea.Quit
			case "n":
				km.ShowTaskForm()
				return km, nil
			case "p":
				km.showProjectSwitch = true
				return km, nil
			case "e":
				if selectedItem := km.Tasks[km.activeListIndex].SelectedItem(); selectedItem != nil {
					if task, ok := selectedItem.(domain.Task); ok {
						km.ShowTaskEditForm(task.ID, task.Name, task.Desc)
						km.showForm = true
					}
				}
				return km, nil
			case "d":
				if selectedItem := km.Tasks[km.activeListIndex].SelectedItem(); selectedItem != nil {
					if task, ok := selectedItem.(domain.Task); ok {
						km.showTaskDeleteConfirm = true
						km.taskToDelete = task.ID
					}
				}
				return km, nil
			case "enter":
				if selectedItem := km.Tasks[km.activeListIndex].SelectedItem(); selectedItem != nil {
					if task, ok := selectedItem.(domain.Task); ok {
						km.MoveTaskToNextStatus(task.ID)
					}
				}
				return km, nil
			case "backspace":
				if selectedItem := km.Tasks[km.activeListIndex].SelectedItem(); selectedItem != nil {
					if task, ok := selectedItem.(domain.Task); ok {
						km.MoveTaskToPreviousStatus(task.ID)
					}
				}
				return km, nil
			}
		}
	case tea.WindowSizeMsg:
		km.width = msg.Width
		km.height = msg.Height

		h, v := styles.DefaultStyle.GetFrameSize()
		columnWidth := max(20, (msg.Width-(h*3))/3)
		km.Tasks[domain.NotStarted].SetSize(columnWidth, msg.Height-v)
		km.Tasks[domain.InProgress].SetSize(columnWidth, msg.Height-v)
		km.Tasks[domain.Done].SetSize(columnWidth, msg.Height-v)
	}

	var cmd tea.Cmd
	if !km.showForm && !km.showProjectSwitch && !km.showTaskDeleteConfirm {
		km.Tasks[km.activeListIndex], cmd = km.Tasks[km.activeListIndex].Update(msg)
	}
	return km, cmd
}

func (km *KahnModel) GetTaskToDelete() string {
	return km.taskToDelete
}

func (km *KahnModel) GetProjectToDelete() string {
	return km.projectToDelete
}

func (km *KahnModel) IsShowingTaskDeleteConfirm() bool {
	return km.showTaskDeleteConfirm
}

func (km *KahnModel) IsShowingProjectDeleteConfirm() bool {
	return km.showProjectDeleteConfirm
}

func (km *KahnModel) GetTaskItems(status domain.Status) []list.Item {
	return km.Tasks[status].Items()
}

func (km *KahnModel) GetActiveListIndex() domain.Status {
	return km.activeListIndex
}

func (km *KahnModel) IsShowingForm() bool {
	return km.showForm
}

func (km *KahnModel) IsShowingProjectSwitch() bool {
	return km.showProjectSwitch
}

func (km *KahnModel) ExecuteTaskDeletion() *KahnModel {
	result := km.executeTaskDeletion()
	model, _ := result.(*KahnModel)
	return model
}

func convertTasksToListItems(tasks []domain.Task) []list.Item {
	items := make([]list.Item, len(tasks))
	for i, task := range tasks {
		items[i] = task
	}
	return items
}

func NewKahnModel(database *database.Database) *KahnModel {
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

	return &KahnModel{
		Projects:               projects,
		ActiveProjectID:        activeProjectID,
		Tasks:                  taskLists,
		taskInputComponents:    taskInputComponents,
		projectInputComponents: projectInputComponents,
		board:                  components.NewBoard(),
		projectSwitcher:        components.NewProjectSwitcher(),
		width:                  80,
		height:                 24,
		database:               database,
		taskService:            taskService,
		projectService:         projectService,
		inputHandler:           input.NewHandler(),
	}
}
