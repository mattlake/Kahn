package services

import (
	"kahn/internal/domain"
	"sort"
	"time"
)

// MockTaskRepository implements domain.TaskRepository for testing
type MockTaskRepository struct {
	tasks     []domain.Task
	nextIntID int
}

func NewMockTaskRepository() *MockTaskRepository {
	return &MockTaskRepository{tasks: []domain.Task{}, nextIntID: 1}
}

func (r *MockTaskRepository) Create(task *domain.Task) error {
	// Simulate auto-increment for IntID
	task.IntID = r.nextIntID
	r.nextIntID++

	r.tasks = append(r.tasks, *task)
	return nil
}

func (r *MockTaskRepository) GetByID(id string) (*domain.Task, error) {
	for i, task := range r.tasks {
		if task.ID == id {
			// Return a copy to ensure we get current values
			taskCopy := r.tasks[i]
			return &taskCopy, nil
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

	// Apply same ordering logic as SQLite repository
	if status == domain.NotStarted {
		// Not Started: priority DESC, then created_at ASC (oldest highest priority first)
		sort.Slice(result, func(i, j int) bool {
			if result[i].Priority != result[j].Priority {
				return result[i].Priority > result[j].Priority // Higher priority first
			}
			return result[i].CreatedAt.Before(result[j].CreatedAt) // Older created first
		})
	} else {
		// In Progress and Done: updated_at DESC (newest changes first)
		sort.Slice(result, func(i, j int) bool {
			return result[i].UpdatedAt.After(result[j].UpdatedAt)
		})
	}

	return result, nil
}

func (r *MockTaskRepository) Update(task *domain.Task) error {
	for i, t := range r.tasks {
		if t.ID == task.ID {
			// Create a copy of provided task and set UpdatedAt
			updatedTask := *task
			updatedTask.UpdatedAt = time.Now()
			r.tasks[i] = updatedTask
			break
		}
	}
	return nil
}

func (r *MockTaskRepository) UpdateStatus(taskID string, status domain.Status) error {
	for i, task := range r.tasks {
		if task.ID == taskID {
			// Create a copy of task with updated values
			updatedTask := task
			updatedTask.Status = status
			updatedTask.UpdatedAt = time.Now()
			r.tasks[i] = updatedTask
			break
		}
	}
	return nil
}

func (r *MockTaskRepository) ClearBlockersForIntID(intID int) error {
	// Find and clear BlockedBy for all tasks that are blocked by this intID
	for i, task := range r.tasks {
		if task.BlockedBy != nil && *task.BlockedBy == intID {
			updatedTask := task
			updatedTask.BlockedBy = nil
			updatedTask.UpdatedAt = time.Now()
			r.tasks[i] = updatedTask
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
