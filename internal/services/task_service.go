package services

import (
	"kahn/internal/domain"
)

type TaskService struct {
	taskRepo    domain.TaskRepository
	projectRepo domain.ProjectRepository
	validator   *ServiceValidator
}

func NewTaskService(taskRepo domain.TaskRepository, projectRepo domain.ProjectRepository) *TaskService {
	return &TaskService{
		taskRepo:    taskRepo,
		projectRepo: projectRepo,
		validator:   NewServiceValidator(),
	}
}

func (ts *TaskService) CreateTask(name, description, projectID string, taskType domain.TaskType, priority domain.Priority) (*domain.Task, error) {

	_, err := ts.validator.ValidateProjectExists(ts.projectRepo, projectID)
	if err != nil {
		return nil, err
	}

	task := domain.NewTask(name, description, projectID)
	task.Type = taskType
	task.Priority = priority

	if err := task.Validate(); err != nil {
		return nil, err
	}

	if err := ts.taskRepo.Create(task); err != nil {
		return nil, domain.NewRepositoryError("create", "task", task.ID, err)
	}

	return task, nil
}

func (ts *TaskService) UpdateTask(id, name, description string, taskType domain.TaskType, priority domain.Priority) (*domain.Task, error) {
	task, err := ts.validator.ValidateTaskExists(ts.taskRepo, id)
	if err != nil {
		return nil, err
	}

	// Update task fields
	task.Name = name
	task.Desc = description
	task.Type = taskType
	task.Priority = priority

	if err := task.Validate(); err != nil {
		return nil, err
	}

	if err := ts.taskRepo.Update(task); err != nil {
		return nil, domain.NewRepositoryError("update", "task", id, err)
	}

	return task, nil
}

func (ts *TaskService) DeleteTask(id string) error {
	_, err := ts.validator.ValidateTaskExists(ts.taskRepo, id)
	if err != nil {
		return err
	}

	if err := ts.taskRepo.Delete(id); err != nil {
		return domain.NewRepositoryError("delete", "task", id, err)
	}

	return nil
}

func (ts *TaskService) MoveTaskToNextStatus(id string) (*domain.Task, error) {
	task, err := ts.validator.ValidateTaskExists(ts.taskRepo, id)
	if err != nil {
		return nil, err
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
		return nil, domain.NewRepositoryError("update status", "task", id, err)
	}

	task.Status = nextStatus
	return task, nil
}

func (ts *TaskService) MoveTaskToPreviousStatus(id string) (*domain.Task, error) {
	task, err := ts.validator.ValidateTaskExists(ts.taskRepo, id)
	if err != nil {
		return nil, err
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
		return nil, domain.NewRepositoryError("update status", "task", id, err)
	}

	task.Status = prevStatus
	return task, nil
}

func (ts *TaskService) GetTask(id string) (*domain.Task, error) {
	task, err := ts.taskRepo.GetByID(id)
	if err != nil {
		return nil, domain.NewRepositoryError("get", "task", id, err)
	}
	return task, nil
}

func (ts *TaskService) GetTasksByProject(projectID string) ([]domain.Task, error) {
	if err := ts.validator.ValidateEntityID(projectID, "project"); err != nil {
		return nil, err
	}

	tasks, err := ts.taskRepo.GetByProjectID(projectID)
	if err != nil {
		return nil, domain.NewRepositoryError("get by project", "tasks", projectID, err)
	}

	return tasks, nil
}

func (ts *TaskService) GetTasksByStatus(projectID string, status domain.Status) ([]domain.Task, error) {
	if err := ts.validator.ValidateEntityID(projectID, "project"); err != nil {
		return nil, err
	}

	tasks, err := ts.taskRepo.GetByStatus(projectID, status)
	if err != nil {
		return nil, domain.NewRepositoryError("get by status", "tasks", projectID, err)
	}

	return tasks, nil
}

func (ts *TaskService) UpdateTaskStatus(id string, status domain.Status) (*domain.Task, error) {
	if err := ts.validator.ValidateEntityID(id, "task"); err != nil {
		return nil, err
	}

	task, err := ts.validator.ValidateTaskExists(ts.taskRepo, id)
	if err != nil {
		return nil, err
	}

	if err := ts.taskRepo.UpdateStatus(id, status); err != nil {
		return nil, domain.NewRepositoryError("update status", "task", id, err)
	}

	task.Status = status
	return task, nil
}
