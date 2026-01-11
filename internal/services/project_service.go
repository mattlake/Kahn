package services

import (
	"kahn/internal/domain"
)

type ProjectService struct {
	projectRepo domain.ProjectRepository
	taskRepo    domain.TaskRepository
	validator   *ServiceValidator
}

func NewProjectService(projectRepo domain.ProjectRepository, taskRepo domain.TaskRepository) *ProjectService {
	return &ProjectService{
		projectRepo: projectRepo,
		taskRepo:    taskRepo,
		validator:   NewServiceValidator(),
	}
}

func (ps *ProjectService) CreateProject(name, description string) (*domain.Project, error) {
	validator := domain.NewFieldValidator()
	if err := validator.ValidateNotEmpty("name", name, "project"); err != nil {
		return nil, err
	}

	project := domain.NewProject(name, description, "#89b4fa")

	if err := project.Validate(); err != nil {
		return nil, err
	}

	if err := ps.projectRepo.Create(project); err != nil {
		return nil, domain.NewRepositoryError("create", "project", project.ID, err)
	}

	return project, nil
}

func (ps *ProjectService) GetProject(id string) (*domain.Project, error) {
	if err := ps.validator.ValidateEntityID(id, "project"); err != nil {
		return nil, err
	}

	project, err := ps.projectRepo.GetByID(id)
	if err != nil {
		return nil, domain.NewRepositoryError("get", "project", id, err)
	}

	return project, nil
}

func (ps *ProjectService) GetAllProjects() ([]domain.Project, error) {
	projects, err := ps.projectRepo.GetAll()
	if err != nil {
		return nil, domain.NewRepositoryError("get all", "projects", "", err)
	}

	return projects, nil
}

func (ps *ProjectService) UpdateProject(id, name, description string) (*domain.Project, error) {
	validator := domain.NewFieldValidator()
	if err := validator.ValidateNotEmpty("name", name, "project"); err != nil {
		return nil, err
	}

	project, err := ps.validator.ValidateProjectExists(ps.projectRepo, id)
	if err != nil {
		return nil, err
	}

	project.Name = name
	project.Description = description

	if err := ps.projectRepo.Update(project); err != nil {
		return nil, domain.NewRepositoryError("update", "project", id, err)
	}

	return project, nil
}

func (ps *ProjectService) DeleteProject(id string) error {
	_, err := ps.validator.ValidateProjectExists(ps.projectRepo, id)
	if err != nil {
		return err
	}

	if err := ps.projectRepo.Delete(id); err != nil {
		return domain.NewRepositoryError("delete", "project", id, err)
	}

	return nil
}

func (ps *ProjectService) GetProjectWithTasks(id string) (*domain.Project, error) {
	if err := ps.validator.ValidateEntityID(id, "project"); err != nil {
		return nil, err
	}

	project, err := ps.projectRepo.GetByID(id)
	if err != nil {
		return nil, domain.NewRepositoryError("get", "project", id, err)
	}
	if project == nil {
		return nil, domain.NewValidationError("id", "project not found")
	}

	tasks, err := ps.taskRepo.GetByProjectID(id)
	if err != nil {
		return nil, domain.NewRepositoryError("get tasks for", "project", id, err)
	}

	project.Tasks = tasks
	return project, nil
}
