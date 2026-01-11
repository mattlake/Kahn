package services

import (
	"kahn/internal/domain"
)

// ServiceValidator provides common validation utilities for services
type ServiceValidator struct{}

// NewServiceValidator creates a new ServiceValidator instance
func NewServiceValidator() *ServiceValidator {
	return &ServiceValidator{}
}

// ValidateEntityExists checks if an entity exists in the repository
func (v *ServiceValidator) ValidateEntityExists(repo interface{}, id, entityType string) (interface{}, error) {
	switch repo := repo.(type) {
	case domain.TaskRepository:
		task, err := repo.GetByID(id)
		if err != nil {
			return nil, domain.NewRepositoryError("get", entityType, id, err)
		}
		if task == nil {
			return nil, domain.NewValidationError("id", entityType+" not found")
		}
		return task, nil
	case domain.ProjectRepository:
		project, err := repo.GetByID(id)
		if err != nil {
			return nil, domain.NewRepositoryError("get", entityType, id, err)
		}
		if project == nil {
			return nil, domain.NewValidationError("id", entityType+" not found")
		}
		return project, nil
	default:
		return nil, domain.NewValidationError("repository", "unsupported repository type")
	}
}

// ValidateEntityID checks if an entity ID is not empty
func (v *ServiceValidator) ValidateEntityID(id, entityType string) error {
	if id == "" {
		return domain.NewEmptyValidationError("id", entityType)
	}
	return nil
}

// ValidateProjectExists checks if a project exists and returns it
func (v *ServiceValidator) ValidateProjectExists(projectRepo domain.ProjectRepository, projectID string) (*domain.Project, error) {
	if projectID == "" {
		return nil, domain.NewEmptyValidationError("project_id", "project")
	}

	project, err := projectRepo.GetByID(projectID)
	if err != nil {
		return nil, domain.NewRepositoryError("get", "project", projectID, err)
	}
	if project == nil {
		return nil, domain.NewValidationError("project_id", "project not found")
	}

	return project, nil
}

// ValidateTaskExists checks if a task exists and returns it
func (v *ServiceValidator) ValidateTaskExists(taskRepo domain.TaskRepository, taskID string) (*domain.Task, error) {
	if taskID == "" {
		return nil, domain.NewEmptyValidationError("id", "task")
	}

	task, err := taskRepo.GetByID(taskID)
	if err != nil {
		return nil, domain.NewRepositoryError("get", "task", taskID, err)
	}
	if task == nil {
		return nil, domain.NewValidationError("id", "task not found")
	}

	return task, nil
}
