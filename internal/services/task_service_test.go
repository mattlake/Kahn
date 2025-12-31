package services

import (
	"kahn/internal/domain"
	"testing"
)

func TestTaskService_CreateTask(t *testing.T) {
	// Setup
	taskRepo := NewMockTaskRepository()
	projectRepo := NewMockProjectRepository()

	// Create a test project
	testProject := domain.NewProject("Test Project", "Test Description", "#89b4fa")
	projectRepo.Create(testProject)

	service := NewTaskService(taskRepo, projectRepo)

	t.Run("successful task creation", func(t *testing.T) {
		// Act
		task, err := service.CreateTask("Test Task", "Test Description", testProject.ID, domain.Low)

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
		task, err := service.CreateTask("", "Test Description", "project-123", domain.Low)

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
		task, err := service.CreateTask("Test Task", "Test Description", "", domain.Low)

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
	taskRepo := NewMockTaskRepository()
	projectRepo := NewMockProjectRepository()
	service := NewTaskService(taskRepo, projectRepo)

	// Create test task
	task := domain.NewTask("Test Task", "Test Description", "project-123")
	task.Status = domain.NotStarted
	taskRepo.Create(task)

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
