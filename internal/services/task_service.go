package services

import (
	"fmt"
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

func (ts *TaskService) CreateTask(name, description, projectID string, taskType domain.TaskType, priority domain.Priority, blockedByIntID *int) (*domain.Task, error) {

	_, err := ts.validator.ValidateProjectExists(ts.projectRepo, projectID)
	if err != nil {
		return nil, err
	}

	task := domain.NewTask(name, description, projectID)
	task.Type = taskType
	task.Priority = priority
	task.BlockedBy = blockedByIntID

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
	task, err := ts.validator.ValidateTaskExists(ts.taskRepo, id)
	if err != nil {
		return err
	}

	if err := ts.taskRepo.Delete(id); err != nil {
		return domain.NewRepositoryError("delete", "task", id, err)
	}

	// Unblock dependent tasks to trigger UI refresh.
	// Database constraint also handles this, but explicit call ensures UI state updates.
	if task.IntID != 0 {
		if err := ts.UnblockDependents(task.IntID); err != nil {
			// Ignore error; task was deleted successfully
			return nil
		}
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

	if nextStatus == domain.Done {
		if err := ts.UnblockDependents(task.IntID); err != nil {
			// Ignore error; status was updated successfully
			return task, nil
		}
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

	if prevStatus == domain.Done {
		if err := ts.UnblockDependents(task.IntID); err != nil {
			// Ignore error; status was updated successfully
			return task, nil
		}
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

	if status == domain.Done {
		if err := ts.UnblockDependents(task.IntID); err != nil {
			// Ignore error; status was updated successfully
			return task, nil
		}
	}

	task.Status = status
	return task, nil
}

// UnblockDependents clears the BlockedBy field for all tasks blocked by the given intID.
// Called when a task is moved to Done or deleted to ensure dependent tasks can proceed.
func (ts *TaskService) UnblockDependents(intID int) error {
	if intID == 0 {
		return nil
	}

	if err := ts.taskRepo.ClearBlockersForIntID(intID); err != nil {
		return domain.NewRepositoryError("clear blockers", "tasks", fmt.Sprintf("int_id=%d", intID), err)
	}

	return nil
}

// SetTaskBlockedBy sets or clears the BlockedBy field for a task
func (ts *TaskService) SetTaskBlockedBy(taskID string, blockedByIntID *int) (*domain.Task, error) {
	// Validate the task to be updated exists
	task, err := ts.validator.ValidateTaskExists(ts.taskRepo, taskID)
	if err != nil {
		return nil, err
	}

	// If setting a blocker, validate it exists and is in the same project
	if blockedByIntID != nil {
		// Get all tasks for the project to find the blocking task
		allTasks, err := ts.taskRepo.GetByProjectID(task.ProjectID)
		if err != nil {
			return nil, domain.NewRepositoryError("get tasks", "project", task.ProjectID, err)
		}

		// Find the blocking task by IntID
		var blockingTask *domain.Task
		for i := range allTasks {
			if allTasks[i].IntID == *blockedByIntID {
				blockingTask = &allTasks[i]
				break
			}
		}

		if blockingTask == nil {
			return nil, domain.NewValidationError("blocked_by", "blocking task not found")
		}

		// Validate that the blocking task is in the same project
		if blockingTask.ProjectID != task.ProjectID {
			return nil, domain.NewValidationError("blocked_by", "blocking task must be in the same project")
		}

		// Validate that task doesn't block itself
		if blockingTask.IntID == task.IntID {
			return nil, domain.NewValidationError("blocked_by", "task cannot block itself")
		}
	}

	// Update the BlockedBy field
	task.BlockedBy = blockedByIntID

	if err := task.Validate(); err != nil {
		return nil, err
	}

	if err := ts.taskRepo.Update(task); err != nil {
		return nil, domain.NewRepositoryError("update", "task", taskID, err)
	}

	return task, nil
}
