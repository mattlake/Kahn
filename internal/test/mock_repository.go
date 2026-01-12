package test

import (
	"kahn/internal/domain"
	"sort"
	"time"
)

type GenericMockRepository[T any] struct {
	items       []T
	idExtractor func(T) string
}

func NewGenericMockRepository[T any](idExtractor func(T) string) *GenericMockRepository[T] {
	return &GenericMockRepository[T]{
		items:       []T{},
		idExtractor: idExtractor,
	}
}

func (r *GenericMockRepository[T]) Add(item T) error {
	r.items = append(r.items, item)
	return nil
}

func (r *GenericMockRepository[T]) GetByID(id string) (*T, error) {
	for i, item := range r.items {
		if r.idExtractor(item) == id {
			itemCopy := r.items[i]
			return &itemCopy, nil
		}
	}
	var zero T
	return &zero, domain.NewRepositoryError("get", "item", id, nil)
}

func (r *GenericMockRepository[T]) GetAll() []T {
	result := make([]T, len(r.items))
	copy(result, r.items)
	return result
}

func (r *GenericMockRepository[T]) Update(item T) error {
	itemID := r.idExtractor(item)
	for i, existing := range r.items {
		if r.idExtractor(existing) == itemID {
			r.items[i] = item
			return nil
		}
	}
	return domain.NewRepositoryError("update", "item", itemID, nil)
}

func (r *GenericMockRepository[T]) Delete(id string) error {
	for i, item := range r.items {
		if r.idExtractor(item) == id {
			r.items = append(r.items[:i], r.items[i+1:]...)
			return nil
		}
	}
	return domain.NewRepositoryError("delete", "item", id, nil)
}

func (r *GenericMockRepository[T]) Filter(predicate func(T) bool) []T {
	var result []T
	for _, item := range r.items {
		if predicate(item) {
			result = append(result, item)
		}
	}
	return result
}

func (r *GenericMockRepository[T]) Count() int {
	return len(r.items)
}

func (r *GenericMockRepository[T]) Clear() {
	r.items = []T{}
}

type MockTaskRepository struct {
	*GenericMockRepository[domain.Task]
	nextIntID int
}

func NewMockTaskRepository() *MockTaskRepository {
	repo := NewGenericMockRepository(func(task domain.Task) string {
		return task.ID
	})
	return &MockTaskRepository{
		GenericMockRepository: repo,
		nextIntID:             1,
	}
}

func (r *MockTaskRepository) Create(task *domain.Task) error {
	// Simulate auto-increment for IntID
	task.IntID = r.nextIntID
	r.nextIntID++

	taskCopy := *task
	return r.Add(taskCopy)
}

func (r *MockTaskRepository) GetByID(id string) (*domain.Task, error) {
	task, err := r.GenericMockRepository.GetByID(id)
	return task, err
}

func (r *MockTaskRepository) GetByProjectID(projectID string) ([]domain.Task, error) {
	tasks := r.Filter(func(task domain.Task) bool {
		return task.ProjectID == projectID
	})
	return tasks, nil
}

func (r *MockTaskRepository) GetByStatus(projectID string, status domain.Status) ([]domain.Task, error) {
	tasks := r.Filter(func(task domain.Task) bool {
		return task.ProjectID == projectID && task.Status == status
	})

	// Apply same ordering logic as SQLite repository
	if status == domain.NotStarted {
		// Not Started: priority DESC, then created_at ASC (oldest highest priority first)
		sort.Slice(tasks, func(i, j int) bool {
			if tasks[i].Priority != tasks[j].Priority {
				return tasks[i].Priority > tasks[j].Priority // Higher priority first
			}
			return tasks[i].CreatedAt.Before(tasks[j].CreatedAt) // Older created first
		})
	} else {
		// In Progress and Done: updated_at DESC (newest changes first)
		sort.Slice(tasks, func(i, j int) bool {
			return tasks[i].UpdatedAt.After(tasks[j].UpdatedAt)
		})
	}

	return tasks, nil
}

func (r *MockTaskRepository) Update(task *domain.Task) error {
	taskCopy := *task
	taskCopy.UpdatedAt = time.Now()
	return r.GenericMockRepository.Update(taskCopy)
}

func (r *MockTaskRepository) UpdateStatus(taskID string, status domain.Status) error {
	for i, task := range r.items {
		if task.ID == taskID {
			updatedTask := task
			updatedTask.Status = status
			updatedTask.UpdatedAt = time.Now()
			r.items[i] = updatedTask
			return nil
		}
	}
	return domain.NewRepositoryError("update", "task status", taskID, nil)
}

func (r *MockTaskRepository) ClearBlockersForIntID(intID int) error {
	for i, task := range r.items {
		if task.BlockedBy != nil && *task.BlockedBy == intID {
			updatedTask := task
			updatedTask.BlockedBy = nil
			updatedTask.UpdatedAt = time.Now()
			r.items[i] = updatedTask
		}
	}
	return nil
}

func (r *MockTaskRepository) Delete(id string) error {
	return r.GenericMockRepository.Delete(id)
}

type MockProjectRepository struct {
	*GenericMockRepository[domain.Project]
}

func NewMockProjectRepository() *MockProjectRepository {
	repo := NewGenericMockRepository(func(project domain.Project) string {
		return project.ID
	})
	return &MockProjectRepository{GenericMockRepository: repo}
}

func (r *MockProjectRepository) Create(project *domain.Project) error {
	projectCopy := *project
	return r.Add(projectCopy)
}

func (r *MockProjectRepository) GetByID(id string) (*domain.Project, error) {
	project, err := r.GenericMockRepository.GetByID(id)
	return project, err
}

func (r *MockProjectRepository) GetAll() ([]domain.Project, error) {
	return r.GenericMockRepository.GetAll(), nil
}

func (r *MockProjectRepository) Update(project *domain.Project) error {
	projectCopy := *project
	return r.GenericMockRepository.Update(projectCopy)
}

func (r *MockProjectRepository) Delete(id string) error {
	return r.GenericMockRepository.Delete(id)
}
