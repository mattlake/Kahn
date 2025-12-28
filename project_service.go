package main

import (
	"strings"
)

type ProjectService struct {
	projectRepo ProjectRepository
	taskRepo    TaskRepository
}

func NewProjectService(projectRepo ProjectRepository, taskRepo TaskRepository) *ProjectService {
	return &ProjectService{
		projectRepo: projectRepo,
		taskRepo:    taskRepo,
	}
}

func (ps *ProjectService) CreateProject(name, description string) (*Project, error) {
	if strings.TrimSpace(name) == "" {
		return nil, &ValidationError{Field: "name", Message: "project name cannot be empty"}
	}

	project := NewProject(name, description, "#89b4fa")

	if err := project.Validate(); err != nil {
		return nil, err
	}

	if err := ps.projectRepo.Create(project); err != nil {
		return nil, &RepositoryError{Operation: "create", Entity: "project", Cause: err}
	}

	return project, nil
}

func (ps *ProjectService) GetProject(id string) (*Project, error) {
	if id == "" {
		return nil, &ValidationError{Field: "id", Message: "project ID cannot be empty"}
	}

	project, err := ps.projectRepo.GetByID(id)
	if err != nil {
		return nil, &RepositoryError{Operation: "get", Entity: "project", ID: id, Cause: err}
	}

	return project, nil
}

func (ps *ProjectService) GetAllProjects() ([]Project, error) {
	projects, err := ps.projectRepo.GetAll()
	if err != nil {
		return nil, &RepositoryError{Operation: "get all", Entity: "projects", Cause: err}
	}

	return projects, nil
}

func (ps *ProjectService) UpdateProject(id, name, description string) (*Project, error) {
	if strings.TrimSpace(name) == "" {
		return nil, &ValidationError{Field: "name", Message: "project name cannot be empty"}
	}

	project, err := ps.projectRepo.GetByID(id)
	if err != nil {
		return nil, &RepositoryError{Operation: "get", Entity: "project", ID: id, Cause: err}
	}
	if project == nil {
		return nil, &ValidationError{Field: "id", Message: "project not found"}
	}

	project.Name = name
	project.Description = description

	if err := ps.projectRepo.Update(project); err != nil {
		return nil, &RepositoryError{Operation: "update", Entity: "project", ID: id, Cause: err}
	}

	return project, nil
}

func (ps *ProjectService) DeleteProject(id string) error {
	if id == "" {
		return &ValidationError{Field: "id", Message: "project ID cannot be empty"}
	}

	project, err := ps.projectRepo.GetByID(id)
	if err != nil {
		return &RepositoryError{Operation: "get", Entity: "project", ID: id, Cause: err}
	}
	if project == nil {
		return &ValidationError{Field: "id", Message: "project not found"}
	}

	if err := ps.projectRepo.Delete(id); err != nil {
		return &RepositoryError{Operation: "delete", Entity: "project", ID: id, Cause: err}
	}

	return nil
}

func (ps *ProjectService) GetProjectWithTasks(id string) (*Project, error) {
	if id == "" {
		return nil, &ValidationError{Field: "id", Message: "project ID cannot be empty"}
	}

	project, err := ps.projectRepo.GetByID(id)
	if err != nil {
		return nil, &RepositoryError{Operation: "get", Entity: "project", ID: id, Cause: err}
	}
	if project == nil {
		return nil, &ValidationError{Field: "id", Message: "project not found"}
	}

	tasks, err := ps.taskRepo.GetByProjectID(id)
	if err != nil {
		return nil, &RepositoryError{Operation: "get tasks for", Entity: "project", ID: id, Cause: err}
	}

	project.Tasks = tasks
	return project, nil
}

func (ps *ProjectService) GetProjectsCount() (int, error) {
	projects, err := ps.projectRepo.GetAll()
	if err != nil {
		return 0, &RepositoryError{Operation: "get all", Entity: "projects", Cause: err}
	}

	return len(projects), nil
}

func (ps *ProjectService) ValidateProjectExists(id string) error {
	if id == "" {
		return &ValidationError{Field: "id", Message: "project ID cannot be empty"}
	}

	project, err := ps.projectRepo.GetByID(id)
	if err != nil {
		return &RepositoryError{Operation: "get", Entity: "project", ID: id, Cause: err}
	}
	if project == nil {
		return &ValidationError{Field: "id", Message: "project not found"}
	}

	return nil
}
