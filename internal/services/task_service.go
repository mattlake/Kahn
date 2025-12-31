package services

import (
	"kahn/internal/domain"
	"strings"
)

type TaskService struct {
	taskRepo    domain.TaskRepository
	projectRepo domain.ProjectRepository
}

func NewTaskService(taskRepo domain.TaskRepository, projectRepo domain.ProjectRepository) *TaskService {
	return &TaskService{
		taskRepo:    taskRepo,
		projectRepo: projectRepo,
	}
}

func (ts *TaskService) CreateTask(name, description, projectID string, priority domain.Priority) (*domain.Task, error) {
	if strings.TrimSpace(name) == "" {
		return nil, &domain.ValidationError{Field: "name", Message: "task name cannot be empty"}
	}

	if projectID == "" {
		return nil, &domain.ValidationError{Field: "project_id", Message: "project ID cannot be empty"}
	}

	project, err := ts.projectRepo.GetByID(projectID)
	if err != nil {
		return nil, &domain.RepositoryError{Operation: "get", Entity: "project", ID: projectID, Cause: err}
	}
	if project == nil {
		return nil, &domain.ValidationError{Field: "project_id", Message: "project not found"}
	}

	task := domain.NewTask(name, description, projectID)
	task.Priority = priority // Set the specified priority

	if err := ts.taskRepo.Create(task); err != nil {
		return nil, &domain.RepositoryError{Operation: "create", Entity: "task", Cause: err}
	}

	return task, nil
}

func (ts *TaskService) UpdateTask(id, name, description string, priority domain.Priority) (*domain.Task, error) {
	if strings.TrimSpace(name) == "" {
		return nil, &domain.ValidationError{Field: "name", Message: "task name cannot be empty"}
	}

	task, err := ts.taskRepo.GetByID(id)
	if err != nil {
		return nil, &domain.RepositoryError{Operation: "get", Entity: "task", ID: id, Cause: err}
	}
	if task == nil {
		return nil, &domain.ValidationError{Field: "id", Message: "task not found"}
	}

	task.Name = name
	task.Desc = description
	task.Priority = priority // Update the priority

	if err := ts.taskRepo.Update(task); err != nil {
		return nil, &domain.RepositoryError{Operation: "update", Entity: "task", ID: id, Cause: err}
	}

	return task, nil
}

func (ts *TaskService) DeleteTask(id string) error {
	task, err := ts.taskRepo.GetByID(id)
	if err != nil {
		return &domain.RepositoryError{Operation: "get", Entity: "task", ID: id, Cause: err}
	}
	if task == nil {
		return &domain.ValidationError{Field: "id", Message: "task not found"}
	}

	if err := ts.taskRepo.Delete(id); err != nil {
		return &domain.RepositoryError{Operation: "delete", Entity: "task", ID: id, Cause: err}
	}

	return nil
}

func (ts *TaskService) MoveTaskToNextStatus(id string) (*domain.Task, error) {
	task, err := ts.taskRepo.GetByID(id)
	if err != nil {
		return nil, &domain.RepositoryError{Operation: "get", Entity: "task", ID: id, Cause: err}
	}
	if task == nil {
		return nil, &domain.ValidationError{Field: "id", Message: "task not found"}
	}

	var nextStatus domain.Status
	switch task.Status {
	case domain.NotStarted:
		nextStatus = domain.InProgress
	case domain.InProgress:
		nextStatus = domain.Done
	case domain.Done:
		nextStatus = domain.NotStarted
	}

	if err := ts.taskRepo.UpdateStatus(id, nextStatus); err != nil {
		return nil, &domain.RepositoryError{Operation: "update status", Entity: "task", ID: id, Cause: err}
	}

	task.Status = nextStatus
	return task, nil
}

func (ts *TaskService) MoveTaskToPreviousStatus(id string) (*domain.Task, error) {
	task, err := ts.taskRepo.GetByID(id)
	if err != nil {
		return nil, &domain.RepositoryError{Operation: "get", Entity: "task", ID: id, Cause: err}
	}
	if task == nil {
		return nil, &domain.ValidationError{Field: "id", Message: "task not found"}
	}

	var prevStatus domain.Status
	switch task.Status {
	case domain.NotStarted:
		prevStatus = domain.Done
	case domain.InProgress:
		prevStatus = domain.NotStarted
	case domain.Done:
		prevStatus = domain.InProgress
	}

	if err := ts.taskRepo.UpdateStatus(id, prevStatus); err != nil {
		return nil, &domain.RepositoryError{Operation: "update status", Entity: "task", ID: id, Cause: err}
	}

	task.Status = prevStatus
	return task, nil
}

func (ts *TaskService) GetTask(id string) (*domain.Task, error) {
	task, err := ts.taskRepo.GetByID(id)
	if err != nil {
		return nil, &domain.RepositoryError{Operation: "get", Entity: "task", ID: id, Cause: err}
	}
	return task, nil
}

func (ts *TaskService) GetTasksByProject(projectID string) ([]domain.Task, error) {
	if projectID == "" {
		return nil, &domain.ValidationError{Field: "project_id", Message: "project ID cannot be empty"}
	}

	tasks, err := ts.taskRepo.GetByProjectID(projectID)
	if err != nil {
		return nil, &domain.RepositoryError{Operation: "get by project", Entity: "tasks", ID: projectID, Cause: err}
	}

	return tasks, nil
}

func (ts *TaskService) GetTasksByStatus(projectID string, status domain.Status) ([]domain.Task, error) {
	if projectID == "" {
		return nil, &domain.ValidationError{Field: "project_id", Message: "project ID cannot be empty"}
	}

	tasks, err := ts.taskRepo.GetByStatus(projectID, status)
	if err != nil {
		return nil, &domain.RepositoryError{Operation: "get by status", Entity: "tasks", ID: projectID, Cause: err}
	}

	return tasks, nil
}

func (ts *TaskService) UpdateTaskStatus(id string, status domain.Status) (*domain.Task, error) {
	if id == "" {
		return nil, &domain.ValidationError{Field: "id", Message: "task ID cannot be empty"}
	}

	task, err := ts.taskRepo.GetByID(id)
	if err != nil {
		return nil, &domain.RepositoryError{Operation: "get", Entity: "task", ID: id, Cause: err}
	}
	if task == nil {
		return nil, &domain.ValidationError{Field: "id", Message: "task not found"}
	}

	if err := ts.taskRepo.UpdateStatus(id, status); err != nil {
		return nil, &domain.RepositoryError{Operation: "update status", Entity: "task", ID: id, Cause: err}
	}

	task.Status = status
	return task, nil
}
