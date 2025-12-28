package domain

import "fmt"

type TaskRepository interface {
	Create(task *Task) error
	GetByID(id string) (*Task, error)
	GetByProjectID(projectID string) ([]Task, error)
	GetByStatus(projectID string, status Status) ([]Task, error)
	Update(task *Task) error
	UpdateStatus(taskID string, status Status) error
	Delete(id string) error
}

type ProjectRepository interface {
	Create(project *Project) error
	GetByID(id string) (*Project, error)
	GetAll() ([]Project, error)
	Update(project *Project) error
	Delete(id string) error
}

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error on field '%s': %s", e.Field, e.Message)
}

type RepositoryError struct {
	Operation string
	Entity    string
	ID        string
	Cause     error
}

func (e *RepositoryError) Error() string {
	if e.ID != "" {
		return fmt.Sprintf("failed to %s %s with id '%s': %v", e.Operation, e.Entity, e.ID, e.Cause)
	}
	return fmt.Sprintf("failed to %s %s: %v", e.Operation, e.Entity, e.Cause)
}

func (e *RepositoryError) Unwrap() error {
	return e.Cause
}
