package main

// TaskWrapper implements input.TaskInterface for the main Task type
type TaskWrapper struct {
	task Task
}

func (tw *TaskWrapper) GetID() string {
	return tw.task.ID
}

func (tw *TaskWrapper) GetName() string {
	return tw.task.Name
}

func (tw *TaskWrapper) GetDescription() string {
	return tw.task.Desc
}

// ProjectWrapper implements input.ProjectInterface for the main Project type
type ProjectWrapper struct {
	project *Project
}

func (pw *ProjectWrapper) GetID() string {
	return pw.project.ID
}

func (pw *ProjectWrapper) GetName() string {
	return pw.project.Name
}
