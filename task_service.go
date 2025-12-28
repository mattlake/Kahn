package main

import "strings"

type TaskService struct {
	taskRepo    TaskRepository
	projectRepo ProjectRepository
}

func NewTaskService(taskRepo TaskRepository, projectRepo ProjectRepository) *TaskService {
	return &TaskService{
		taskRepo:    taskRepo,
		projectRepo: projectRepo,
	}
}

func (ts *TaskService) CreateTask(name, description, projectID string) (*Task, error) {
	if strings.TrimSpace(name) == "" {
		return nil, &ValidationError{Field: "name", Message: "task name cannot be empty"}
	}

	if projectID == "" {
		return nil, &ValidationError{Field: "project_id", Message: "project ID cannot be empty"}
	}

	project, err := ts.projectRepo.GetByID(projectID)
	if err != nil {
		return nil, &RepositoryError{Operation: "get", Entity: "project", ID: projectID, Cause: err}
	}
	if project == nil {
		return nil, &ValidationError{Field: "project_id", Message: "project not found"}
	}

	task := NewTask(name, description, projectID)

	if err := ts.taskRepo.Create(task); err != nil {
		return nil, &RepositoryError{Operation: "create", Entity: "task", Cause: err}
	}

	return task, nil
}

func (ts *TaskService) UpdateTask(id, name, description string) (*Task, error) {
	if strings.TrimSpace(name) == "" {
		return nil, &ValidationError{Field: "name", Message: "task name cannot be empty"}
	}

	task, err := ts.taskRepo.GetByID(id)
	if err != nil {
		return nil, &RepositoryError{Operation: "get", Entity: "task", ID: id, Cause: err}
	}
	if task == nil {
		return nil, &ValidationError{Field: "id", Message: "task not found"}
	}

	task.Name = name
	task.Desc = description

	if err := ts.taskRepo.Update(task); err != nil {
		return nil, &RepositoryError{Operation: "update", Entity: "task", ID: id, Cause: err}
	}

	return task, nil
}

func (ts *TaskService) DeleteTask(id string) error {
	task, err := ts.taskRepo.GetByID(id)
	if err != nil {
		return &RepositoryError{Operation: "get", Entity: "task", ID: id, Cause: err}
	}
	if task == nil {
		return &ValidationError{Field: "id", Message: "task not found"}
	}

	if err := ts.taskRepo.Delete(id); err != nil {
		return &RepositoryError{Operation: "delete", Entity: "task", ID: id, Cause: err}
	}

	return nil
}

func (ts *TaskService) MoveTaskToNextStatus(id string) (*Task, error) {
	task, err := ts.taskRepo.GetByID(id)
	if err != nil {
		return nil, &RepositoryError{Operation: "get", Entity: "task", ID: id, Cause: err}
	}
	if task == nil {
		return nil, &ValidationError{Field: "id", Message: "task not found"}
	}

	var nextStatus Status
	switch task.Status {
	case NotStarted:
		nextStatus = InProgress
	case InProgress:
		nextStatus = Done
	case Done:
		nextStatus = NotStarted
	}

	if err := ts.taskRepo.UpdateStatus(id, nextStatus); err != nil {
		return nil, &RepositoryError{Operation: "update status", Entity: "task", ID: id, Cause: err}
	}

	task.Status = nextStatus
	return task, nil
}

func (ts *TaskService) MoveTaskToPreviousStatus(id string) (*Task, error) {
	task, err := ts.taskRepo.GetByID(id)
	if err != nil {
		return nil, &RepositoryError{Operation: "get", Entity: "task", ID: id, Cause: err}
	}
	if task == nil {
		return nil, &ValidationError{Field: "id", Message: "task not found"}
	}

	var prevStatus Status
	switch task.Status {
	case NotStarted:
		prevStatus = Done
	case InProgress:
		prevStatus = NotStarted
	case Done:
		prevStatus = InProgress
	}

	if err := ts.taskRepo.UpdateStatus(id, prevStatus); err != nil {
		return nil, &RepositoryError{Operation: "update status", Entity: "task", ID: id, Cause: err}
	}

	task.Status = prevStatus
	return task, nil
}

func (ts *TaskService) GetTask(id string) (*Task, error) {
	task, err := ts.taskRepo.GetByID(id)
	if err != nil {
		return nil, &RepositoryError{Operation: "get", Entity: "task", ID: id, Cause: err}
	}
	return task, nil
}

func (ts *TaskService) GetTasksByProject(projectID string) ([]Task, error) {
	if projectID == "" {
		return nil, &ValidationError{Field: "project_id", Message: "project ID cannot be empty"}
	}

	tasks, err := ts.taskRepo.GetByProjectID(projectID)
	if err != nil {
		return nil, &RepositoryError{Operation: "get by project", Entity: "tasks", ID: projectID, Cause: err}
	}

	return tasks, nil
}

func (ts *TaskService) GetTasksByStatus(projectID string, status Status) ([]Task, error) {
	if projectID == "" {
		return nil, &ValidationError{Field: "project_id", Message: "project ID cannot be empty"}
	}

	tasks, err := ts.taskRepo.GetByStatus(projectID, status)
	if err != nil {
		return nil, &RepositoryError{Operation: "get by status", Entity: "tasks", ID: projectID, Cause: err}
	}

	return tasks, nil
}
