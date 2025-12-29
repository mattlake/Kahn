package services

import (
	"kahn/internal/domain"
	"strings"
)

type ProjectService struct {
	projectRepo domain.ProjectRepository
	taskRepo    domain.TaskRepository
}

func NewProjectService(projectRepo domain.ProjectRepository, taskRepo domain.TaskRepository) *ProjectService {
	return &ProjectService{
		projectRepo: projectRepo,
		taskRepo:    taskRepo,
	}
}

func (ps *ProjectService) CreateProject(name, description string) (*domain.Project, error) {
	if strings.TrimSpace(name) == "" {
		return nil, &domain.ValidationError{Field: "name", Message: "project name cannot be empty"}
	}

	project := domain.NewProject(name, description, "#89b4fa")

	if err := project.Validate(); err != nil {
		return nil, err
	}

	if err := ps.projectRepo.Create(project); err != nil {
		return nil, &domain.RepositoryError{Operation: "create", Entity: "project", Cause: err}
	}

	return project, nil
}

func (ps *ProjectService) GetProject(id string) (*domain.Project, error) {
	if id == "" {
		return nil, &domain.ValidationError{Field: "id", Message: "project ID cannot be empty"}
	}

	project, err := ps.projectRepo.GetByID(id)
	if err != nil {
		return nil, &domain.RepositoryError{Operation: "get", Entity: "project", ID: id, Cause: err}
	}

	return project, nil
}

func (ps *ProjectService) GetAllProjects() ([]domain.Project, error) {
	projects, err := ps.projectRepo.GetAll()
	if err != nil {
		return nil, &domain.RepositoryError{Operation: "get all", Entity: "projects", Cause: err}
	}

	return projects, nil
}

func (ps *ProjectService) UpdateProject(id, name, description string) (*domain.Project, error) {
	if strings.TrimSpace(name) == "" {
		return nil, &domain.ValidationError{Field: "name", Message: "project name cannot be empty"}
	}

	project, err := ps.projectRepo.GetByID(id)
	if err != nil {
		return nil, &domain.RepositoryError{Operation: "get", Entity: "project", ID: id, Cause: err}
	}
	if project == nil {
		return nil, &domain.ValidationError{Field: "id", Message: "project not found"}
	}

	project.Name = name
	project.Description = description

	if err := ps.projectRepo.Update(project); err != nil {
		return nil, &domain.RepositoryError{Operation: "update", Entity: "project", ID: id, Cause: err}
	}

	return project, nil
}

func (ps *ProjectService) DeleteProject(id string) error {
	if id == "" {
		return &domain.ValidationError{Field: "id", Message: "project ID cannot be empty"}
	}

	project, err := ps.projectRepo.GetByID(id)
	if err != nil {
		return &domain.RepositoryError{Operation: "get", Entity: "project", ID: id, Cause: err}
	}
	if project == nil {
		return &domain.ValidationError{Field: "id", Message: "project not found"}
	}

	if err := ps.projectRepo.Delete(id); err != nil {
		return &domain.RepositoryError{Operation: "delete", Entity: "project", ID: id, Cause: err}
	}

	return nil
}

func (ps *ProjectService) GetProjectWithTasks(id string) (*domain.Project, error) {
	if id == "" {
		return nil, &domain.ValidationError{Field: "id", Message: "project ID cannot be empty"}
	}

	project, err := ps.projectRepo.GetByID(id)
	if err != nil {
		return nil, &domain.RepositoryError{Operation: "get", Entity: "project", ID: id, Cause: err}
	}
	if project == nil {
		return nil, &domain.ValidationError{Field: "id", Message: "project not found"}
	}

	tasks, err := ps.taskRepo.GetByProjectID(id)
	if err != nil {
		return nil, &domain.RepositoryError{Operation: "get tasks for", Entity: "project", ID: id, Cause: err}
	}

	project.Tasks = tasks
	return project, nil
}
