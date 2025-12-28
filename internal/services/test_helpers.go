package services

import (
	"kahn/internal/domain"
)

// MockTaskRepository implements domain.TaskRepository for testing
type MockTaskRepository struct {
	tasks []domain.Task
}

func NewMockTaskRepository() *MockTaskRepository {
	return &MockTaskRepository{tasks: []domain.Task{}}
}

func (r *MockTaskRepository) Create(task *domain.Task) error {
	r.tasks = append(r.tasks, *task)
	return nil
}

func (r *MockTaskRepository) GetByID(id string) (*domain.Task, error) {
	for i, task := range r.tasks {
		if task.ID == id {
			return &r.tasks[i], nil
		}
	}
	return &domain.Task{}, &domain.RepositoryError{Operation: "get", Entity: "task", ID: id}
}

func (r *MockTaskRepository) GetByProjectID(projectID string) ([]domain.Task, error) {
	var result []domain.Task
	for _, task := range r.tasks {
		if task.ProjectID == projectID {
			result = append(result, task)
		}
	}
	return result, nil
}

func (r *MockTaskRepository) GetByStatus(projectID string, status domain.Status) ([]domain.Task, error) {
	var result []domain.Task
	for _, task := range r.tasks {
		if task.ProjectID == projectID && task.Status == status {
			result = append(result, task)
		}
	}
	return result, nil
}

func (r *MockTaskRepository) Update(task *domain.Task) error {
	for i, t := range r.tasks {
		if t.ID == task.ID {
			r.tasks[i] = *task
			break
		}
	}
	return nil
}

func (r *MockTaskRepository) UpdateStatus(taskID string, status domain.Status) error {
	for i, task := range r.tasks {
		if task.ID == taskID {
			r.tasks[i].Status = status
			break
		}
	}
	return nil
}

func (r *MockTaskRepository) Delete(id string) error {
	for i, task := range r.tasks {
		if task.ID == id {
			r.tasks = append(r.tasks[:i], r.tasks[i+1:]...)
			return nil
		}
	}
	return &domain.RepositoryError{Operation: "delete", Entity: "task", ID: id}
}

// MockProjectRepository implements domain.ProjectRepository for testing
type MockProjectRepository struct {
	projects []domain.Project
}

func NewMockProjectRepository() *MockProjectRepository {
	return &MockProjectRepository{projects: []domain.Project{}}
}

func (r *MockProjectRepository) Create(project *domain.Project) error {
	r.projects = append(r.projects, *project)
	return nil
}

func (r *MockProjectRepository) GetByID(id string) (*domain.Project, error) {
	for i, project := range r.projects {
		if project.ID == id {
			return &r.projects[i], nil
		}
	}
	return &domain.Project{}, &domain.RepositoryError{Operation: "get", Entity: "project", ID: id}
}

func (r *MockProjectRepository) GetAll() ([]domain.Project, error) {
	return r.projects, nil
}

func (r *MockProjectRepository) Update(project *domain.Project) error {
	for i, p := range r.projects {
		if p.ID == project.ID {
			r.projects[i] = *project
			break
		}
	}
	return nil
}

func (r *MockProjectRepository) Delete(id string) error {
	for i, project := range r.projects {
		if project.ID == id {
			r.projects = append(r.projects[:i], r.projects[i+1:]...)
			return nil
		}
	}
	return &domain.RepositoryError{Operation: "delete", Entity: "project", ID: id}
}
