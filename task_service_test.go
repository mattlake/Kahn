package main

import (
	"kahn/internal/domain"
	"testing"
)

func TestTaskService_CreateTask(t *testing.T) {
	// Setup
	taskRepo := &mockTaskRepository{}
	projectRepo := &mockProjectRepository{}

	// Create a test project
	testProject := domain.NewProject("Test Project", "Test Description", "#89b4fa")
	projectRepo.projects = []domain.Project{*testProject}

	service := NewTaskService(taskRepo, projectRepo)

	t.Run("successful task creation", func(t *testing.T) {
		// Act
		task, err := service.CreateTask("Test Task", "Test Description", testProject.ID)

		// Assert
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if task == nil {
			t.Error("Expected task to be created")
		}
		if task.Name != "Test Task" {
			t.Errorf("Expected task name 'Test Task', got '%s'", task.Name)
		}
		if task.Desc != "Test Description" {
			t.Errorf("Expected task description 'Test Description', got '%s'", task.Desc)
		}
		if task.ProjectID != testProject.ID {
			t.Errorf("Expected project ID '%s', got '%s'", testProject.ID, task.ProjectID)
		}
	})

	t.Run("empty name validation", func(t *testing.T) {
		// Act
		task, err := service.CreateTask("", "Test Description", "project-123")

		// Assert
		if err == nil {
			t.Error("Expected validation error for empty name")
		}
		if task != nil {
			t.Error("Expected no task to be created")
		}
	})

	t.Run("empty project ID validation", func(t *testing.T) {
		// Act
		task, err := service.CreateTask("Test Task", "Test Description", "")

		// Assert
		if err == nil {
			t.Error("Expected validation error for empty project ID")
		}
		if task != nil {
			t.Error("Expected no task to be created")
		}
	})
}

func TestTaskService_MoveTaskToNextStatus(t *testing.T) {
	// Setup
	taskRepo := &mockTaskRepository{}
	projectRepo := &mockProjectRepository{}
	service := NewTaskService(taskRepo, projectRepo)

	// Create test task
	task := domain.NewTask("Test Task", "Test Description", "project-123")
	task.Status = domain.NotStarted
	taskRepo.tasks = []domain.Task{*task}

	t.Run("move from NotStarted to InProgress", func(t *testing.T) {
		// Act
		updatedTask, err := service.MoveTaskToNextStatus(task.ID)

		// Assert
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if updatedTask.Status != domain.InProgress {
			t.Errorf("Expected status InProgress, got %v", updatedTask.Status)
		}
	})

	t.Run("move from InProgress to Done", func(t *testing.T) {
		// Setup
		task.Status = domain.InProgress

		// Act
		updatedTask, err := service.MoveTaskToNextStatus(task.ID)

		// Assert
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if updatedTask.Status != domain.Done {
			t.Errorf("Expected status Done, got %v", updatedTask.Status)
		}
	})

	t.Run("move from Done to NotStarted", func(t *testing.T) {
		// Setup
		task.Status = domain.Done

		// Act
		updatedTask, err := service.MoveTaskToNextStatus(task.ID)

		// Assert
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if updatedTask.Status != domain.NotStarted {
			t.Errorf("Expected status NotStarted, got %v", updatedTask.Status)
		}
	})
}

// Mock implementations for testing
type mockTaskRepository struct {
	tasks []domain.Task
}

func (r *mockTaskRepository) Create(task *domain.Task) error {
	r.tasks = append(r.tasks, *task)
	return nil
}

func (r *mockTaskRepository) GetByID(id string) (*domain.Task, error) {
	for i, task := range r.tasks {
		if task.ID == id {
			return &r.tasks[i], nil
		}
	}
	return &domain.Task{}, &domain.RepositoryError{Operation: "get", Entity: "task", ID: id}
}

func (r *mockTaskRepository) GetByProjectID(projectID string) ([]domain.Task, error) {
	var result []domain.Task
	for _, task := range r.tasks {
		if task.ProjectID == projectID {
			result = append(result, task)
		}
	}
	return result, nil
}

func (r *mockTaskRepository) GetByStatus(projectID string, status domain.Status) ([]domain.Task, error) {
	var result []domain.Task
	for _, task := range r.tasks {
		if task.ProjectID == projectID && task.Status == status {
			result = append(result, task)
		}
	}
	return result, nil
}

func (r *mockTaskRepository) Update(task *domain.Task) error {
	for i, t := range r.tasks {
		if t.ID == task.ID {
			r.tasks[i] = *task
			break
		}
	}
	return nil
}

func (r *mockTaskRepository) UpdateStatus(taskID string, status domain.Status) error {
	for i, task := range r.tasks {
		if task.ID == taskID {
			r.tasks[i].Status = status
			break
		}
	}
	return nil
}

func (r *mockTaskRepository) Delete(id string) error {
	for i, task := range r.tasks {
		if task.ID == id {
			r.tasks = append(r.tasks[:i], r.tasks[i+1:]...)
			return nil
		}
	}
	return &domain.RepositoryError{Operation: "delete", Entity: "task", ID: id}
}

type mockProjectRepository struct {
	projects []domain.Project
}

func (r *mockProjectRepository) Create(project *domain.Project) error {
	r.projects = append(r.projects, *project)
	return nil
}

func (r *mockProjectRepository) GetByID(id string) (*domain.Project, error) {
	for i, project := range r.projects {
		if project.ID == id {
			return &r.projects[i], nil
		}
	}
	return &domain.Project{}, &domain.RepositoryError{Operation: "get", Entity: "project", ID: id}
}

func (r *mockProjectRepository) GetAll() ([]domain.Project, error) {
	return r.projects, nil
}

func (r *mockProjectRepository) Update(project *domain.Project) error {
	for i, p := range r.projects {
		if p.ID == project.ID {
			r.projects[i] = *project
			break
		}
	}
	return nil
}

func (r *mockProjectRepository) Delete(id string) error {
	for i, project := range r.projects {
		if project.ID == id {
			r.projects = append(r.projects[:i], r.projects[i+1:]...)
			return nil
		}
	}
	return &domain.RepositoryError{Operation: "delete", Entity: "project", ID: id}
}
