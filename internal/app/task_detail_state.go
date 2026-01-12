package app

import (
	"kahn/internal/domain"
)

// TaskDetailState manages the state of the task detail modal view
type TaskDetailState struct {
	isShowing    bool
	task         *domain.Task
	blockingTask *domain.Task // The task that blocks the current task (if any)
}

// NewTaskDetailState creates a new task detail state manager
func NewTaskDetailState() *TaskDetailState {
	return &TaskDetailState{
		isShowing:    false,
		task:         nil,
		blockingTask: nil,
	}
}

// ShowTaskDetail shows the task detail modal with the given task
func (tds *TaskDetailState) ShowTaskDetail(task *domain.Task, blockingTask *domain.Task) {
	tds.isShowing = true
	tds.task = task
	tds.blockingTask = blockingTask
}

// HideTaskDetail hides the task detail modal
func (tds *TaskDetailState) HideTaskDetail() {
	tds.isShowing = false
	tds.task = nil
	tds.blockingTask = nil
}

// IsShowingTaskDetail returns true if the task detail modal is currently showing
func (tds *TaskDetailState) IsShowingTaskDetail() bool {
	return tds.isShowing
}

// GetTask returns the current task being displayed
func (tds *TaskDetailState) GetTask() *domain.Task {
	return tds.task
}

// GetBlockingTask returns the task that blocks the current task (if any)
func (tds *TaskDetailState) GetBlockingTask() *domain.Task {
	return tds.blockingTask
}
