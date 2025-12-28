package domain

// TaskWrapper implements input.TaskInterface for the main Task type
type TaskWrapper struct {
	Task Task
}

func (tw *TaskWrapper) GetID() string {
	return tw.Task.ID
}

func (tw *TaskWrapper) GetName() string {
	return tw.Task.Name
}

func (tw *TaskWrapper) GetDescription() string {
	return tw.Task.Desc
}

// ProjectWrapper implements input.ProjectInterface for the main Project type
type ProjectWrapper struct {
	Project *Project
}

func (pw *ProjectWrapper) GetID() string {
	return pw.Project.ID
}

func (pw *ProjectWrapper) GetName() string {
	return pw.Project.Name
}
