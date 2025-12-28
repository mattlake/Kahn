package main

import (
	"testing"
)

func TestTaskService_CreateTask(t *testing.T) {
	// Setup
	taskRepo := &mockTaskRepository{}
	projectRepo := &mockProjectRepository{}

	// Create a test project
	testProject := NewProject("Test Project", "Test Description", "#89b4fa")
	projectRepo.projects = []Project{*testProject}

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
	task := NewTask("Test Task", "Test Description", "project-123")
	task.Status = NotStarted
	taskRepo.tasks = []Task{*task}

	t.Run("move from NotStarted to InProgress", func(t *testing.T) {
		// Act
		updatedTask, err := service.MoveTaskToNextStatus(task.ID)

		// Assert
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if updatedTask.Status != InProgress {
			t.Errorf("Expected status InProgress, got %v", updatedTask.Status)
		}
	})

	t.Run("move from InProgress to Done", func(t *testing.T) {
		// Setup
		task.Status = InProgress

		// Act
		updatedTask, err := service.MoveTaskToNextStatus(task.ID)

		// Assert
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if updatedTask.Status != Done {
			t.Errorf("Expected status Done, got %v", updatedTask.Status)
		}
	})

	t.Run("move from Done to NotStarted", func(t *testing.T) {
		// Setup
		task.Status = Done

		// Act
		updatedTask, err := service.MoveTaskToNextStatus(task.ID)

		// Assert
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if updatedTask.Status != NotStarted {
			t.Errorf("Expected status NotStarted, got %v", updatedTask.Status)
		}
	})
}

// Mock implementations for testing
type mockTaskRepository struct {
	tasks []Task
}

func (r *mockTaskRepository) Create(task *Task) error {
	r.tasks = append(r.tasks, *task)
	return nil
}

func (r *mockTaskRepository) GetByID(id string) (*Task, error) {
	for i, task := range r.tasks {
		if task.ID == id {
			return &r.tasks[i], nil
		}
	}
	return &Task{}, &RepositoryError{Operation: "get", Entity: "task", ID: id}
}

func (r *mockTaskRepository) GetByProjectID(projectID string) ([]Task, error) {
	var result []Task
	for _, task := range r.tasks {
		if task.ProjectID == projectID {
			result = append(result, task)
		}
	}
	return result, nil
}

func (r *mockTaskRepository) GetByStatus(projectID string, status Status) ([]Task, error) {
	var result []Task
	for _, task := range r.tasks {
		if task.ProjectID == projectID && task.Status == status {
			result = append(result, task)
		}
	}
	return result, nil
}

func (r *mockTaskRepository) Update(task *Task) error {
	for i, t := range r.tasks {
		if t.ID == task.ID {
			r.tasks[i] = *task
			break
		}
	}
	return nil
}

func (r *mockTaskRepository) UpdateStatus(taskID string, status Status) error {
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
	return &RepositoryError{Operation: "delete", Entity: "task", ID: id}
}

type mockProjectRepository struct {
	projects []Project
}

func (r *mockProjectRepository) Create(project *Project) error {
	r.projects = append(r.projects, *project)
	return nil
}

func (r *mockProjectRepository) GetByID(id string) (*Project, error) {
	for i, project := range r.projects {
		if project.ID == id {
			return &r.projects[i], nil
		}
	}
	return &Project{}, &RepositoryError{Operation: "get", Entity: "project", ID: id}
}

func (r *mockProjectRepository) GetAll() ([]Project, error) {
	return r.projects, nil
}

func (r *mockProjectRepository) Update(project *Project) error {
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
	return &RepositoryError{Operation: "delete", Entity: "project", ID: id}
}
