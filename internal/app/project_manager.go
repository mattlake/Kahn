package app

import (
	"kahn/internal/domain"
	"kahn/internal/services"
	"kahn/internal/ui/input"
)

// ProjectManager handles all project operations and state management
type ProjectManager struct {
	projects        []domain.Project
	activeProjectID string
	projectService  *services.ProjectService
	taskService     *services.TaskService
	taskListManager *TaskListManager
}

// NewProjectManager creates a new project manager
func NewProjectManager(projectService *services.ProjectService, taskService *services.TaskService, taskListManager *TaskListManager) *ProjectManager {
	return &ProjectManager{
		projectService:  projectService,
		taskService:     taskService,
		taskListManager: taskListManager,
	}
}

// InitializeProjects loads projects from the service and sets up the active project
func (pm *ProjectManager) InitializeProjects() error {
	projects, err := pm.projectService.GetAllProjects()
	if err != nil {
		projects = []domain.Project{}
	}

	for i := range projects {
		tasks, err := pm.taskService.GetTasksByProject(projects[i].ID)
		if err != nil {
			projects[i].Tasks = []domain.Task{}
		} else {
			projects[i].Tasks = tasks
		}
	}

	if len(projects) == 0 {
		newProject, err := pm.projectService.CreateProject("Default Project", "A default project for your tasks")
		if err != nil {
			projects = []domain.Project{}
		} else {
			projects = []domain.Project{*newProject}
		}
	}

	pm.projects = projects
	if len(projects) > 0 {
		pm.activeProjectID = projects[0].ID
	}

	// Initialize task lists for the active project
	activeProj := pm.GetActiveProject()
	if activeProj != nil {
		pm.taskListManager.UpdateTaskLists(activeProj)
	}

	return nil
}

// GetProjects returns all projects as project interfaces
func (pm *ProjectManager) GetProjects() []input.ProjectInterface {
	var projects []input.ProjectInterface
	for i := range pm.projects {
		projects = append(projects, &domain.ProjectWrapper{Project: &pm.projects[i]})
	}
	return projects
}

// GetActiveProject returns the currently active project
func (pm *ProjectManager) GetActiveProject() *domain.Project {
	for i, proj := range pm.projects {
		if proj.ID == pm.activeProjectID {
			return &pm.projects[i]
		}
	}
	return nil
}

// GetActiveProjectID returns the ID of the currently active project
func (pm *ProjectManager) GetActiveProjectID() string {
	return pm.activeProjectID
}

// GetSelectedProjectIndex returns the index of the active project in the projects slice
func (pm *ProjectManager) GetSelectedProjectIndex() int {
	for i, proj := range pm.projects {
		if proj.ID == pm.activeProjectID {
			return i
		}
	}
	return 0
}

// SwitchToProject switches to the specified project ID
func (pm *ProjectManager) SwitchToProject(id string) error {
	pm.activeProjectID = id
	pm.taskListManager.UpdateTaskLists(pm.GetActiveProject())
	return nil
}

// CreateProject creates a new project and makes it active
func (pm *ProjectManager) CreateProject(name, description string) error {
	newProject, err := pm.projectService.CreateProject(name, description)
	if err != nil {
		return err
	}

	pm.projects = append(pm.projects, *newProject)
	pm.activeProjectID = newProject.ID
	pm.taskListManager.UpdateTaskLists(pm.GetActiveProject())

	return nil
}

// DeleteProject deletes a project and handles the active project logic
func (pm *ProjectManager) DeleteProject(id string) error {
	if err := pm.projectService.DeleteProject(id); err != nil {
		return err
	}

	if len(pm.projects) == 1 {
		pm.projects = []domain.Project{}
		pm.activeProjectID = ""
		// Clear all task lists
		pm.taskListManager.MarkAllListsDirty()
		activeProj := pm.GetActiveProject()
		pm.taskListManager.UpdateTaskListsConditional(activeProj)
	} else {
		var newProjects []domain.Project
		var wasActiveProject bool
		for _, proj := range pm.projects {
			if proj.ID != id {
				newProjects = append(newProjects, proj)
			} else {
				wasActiveProject = (proj.ID == pm.activeProjectID)
			}
		}
		pm.projects = newProjects

		if wasActiveProject && len(pm.projects) > 0 {
			pm.activeProjectID = pm.projects[0].ID
			pm.taskListManager.UpdateTaskLists(pm.GetActiveProject())
		}
	}

	return nil
}

// HasProjects returns true if there are any projects
func (pm *ProjectManager) HasProjects() bool {
	return len(pm.projects) > 0
}

// GetProjectCount returns the number of projects
func (pm *ProjectManager) GetProjectCount() int {
	return len(pm.projects)
}

// GetProjectsAsDomain returns projects as domain.Project slice
func (pm *ProjectManager) GetProjectsAsDomain() []domain.Project {
	return pm.projects
}
