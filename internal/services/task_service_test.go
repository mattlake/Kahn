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
		task, err := service.CreateTask("Test Task", "Test Description", testProject.ID, domain.RegularTask, domain.Low)

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
		task, err := service.CreateTask("", "Test Description", "project-123", domain.RegularTask, domain.Low)

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
		task, err := service.CreateTask("Test Task", "Test Description", "", domain.RegularTask, domain.Low)

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

func TestTaskService_TaskType_Functionality(t *testing.T) {
	// Setup mocks
	taskRepo := NewMockTaskRepository()
	projectRepo := NewMockProjectRepository()

	// Create test project
	testProject := domain.NewProject("Test Project", "Test Description", "blue")
	projectRepo.Create(testProject)

	service := NewTaskService(taskRepo, projectRepo)

	t.Run("create task with RegularTask type", func(t *testing.T) {
		// Act
		task, err := service.CreateTask("Regular Task", "Description", testProject.ID, domain.RegularTask, domain.Low)

		// Assert
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if task == nil {
			t.Error("Expected task to be created")
		}
		if task.Type != domain.RegularTask {
			t.Errorf("Expected task type RegularTask, got %v", task.Type)
		}
	})

	t.Run("create task with Bug type", func(t *testing.T) {
		// Act
		task, err := service.CreateTask("Bug Task", "Bug Description", testProject.ID, domain.Bug, domain.High)

		// Assert
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if task == nil {
			t.Error("Expected task to be created")
		}
		if task.Type != domain.Bug {
			t.Errorf("Expected task type Bug, got %v", task.Type)
		}
	})

	t.Run("create task with Feature type", func(t *testing.T) {
		// Act
		task, err := service.CreateTask("Feature Task", "Feature Description", testProject.ID, domain.Feature, domain.Medium)

		// Assert
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if task == nil {
			t.Error("Expected task to be created")
		}
		if task.Type != domain.Feature {
			t.Errorf("Expected task type Feature, got %v", task.Type)
		}
	})

	t.Run("update task type", func(t *testing.T) {
		// Setup - create a task with RegularTask type
		originalTask, err := service.CreateTask("Original Task", "Description", testProject.ID, domain.RegularTask, domain.Low)
		if err != nil {
			t.Fatalf("Failed to create original task: %v", err)
		}

		// Act - update to Bug type
		updatedTask, err := service.UpdateTask(originalTask.ID, "Updated Task", "Updated Description", domain.Bug, domain.High)

		// Assert
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if updatedTask == nil {
			t.Error("Expected task to be updated")
		}
		if updatedTask.Type != domain.Bug {
			t.Errorf("Expected updated task type Bug, got %v", updatedTask.Type)
		}
		if updatedTask.Name != "Updated Task" {
			t.Errorf("Expected updated name 'Updated Task', got '%s'", updatedTask.Name)
		}
		if updatedTask.Desc != "Updated Description" {
			t.Errorf("Expected updated description 'Updated Description', got '%s'", updatedTask.Desc)
		}
		if updatedTask.Priority != domain.High {
			t.Errorf("Expected updated priority High, got %v", updatedTask.Priority)
		}
	})

	t.Run("get tasks by project with different types", func(t *testing.T) {
		// Setup - create tasks with different types
		task1, _ := service.CreateTask("Task 1", "Regular", testProject.ID, domain.RegularTask, domain.Low)
		task2, _ := service.CreateTask("Task 2", "Bug", testProject.ID, domain.Bug, domain.Medium)
		task3, _ := service.CreateTask("Task 3", "Feature", testProject.ID, domain.Feature, domain.High)

		// Act
		tasks, err := service.GetTasksByProject(testProject.ID)

		// Assert
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if len(tasks) < 3 {
			t.Errorf("Expected at least 3 tasks, got %d", len(tasks))
		}

		// Verify we have all task types
		typeCount := make(map[domain.TaskType]int)
		for _, task := range tasks {
			if task.ID == task1.ID || task.ID == task2.ID || task.ID == task3.ID {
				typeCount[task.Type]++
			}
		}

		if typeCount[domain.RegularTask] == 0 {
			t.Error("Expected at least one RegularTask")
		}
		if typeCount[domain.Bug] == 0 {
			t.Error("Expected at least one Bug task")
		}
		if typeCount[domain.Feature] == 0 {
			t.Error("Expected at least one Feature task")
		}
	})

	t.Run("get tasks by status with different types", func(t *testing.T) {
		// Setup - create tasks with different types and same status
		task1, _ := service.CreateTask("Not Started Regular", "Regular", testProject.ID, domain.RegularTask, domain.Low)
		task2, _ := service.CreateTask("Not Started Bug", "Bug", testProject.ID, domain.Bug, domain.Medium)
		task3, _ := service.CreateTask("Not Started Feature", "Feature", testProject.ID, domain.Feature, domain.High)

		// Act
		tasks, err := service.GetTasksByStatus(testProject.ID, domain.NotStarted)

		// Assert
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if len(tasks) < 3 {
			t.Errorf("Expected at least 3 NotStarted tasks, got %d", len(tasks))
		}

		// Verify we have all task types in NotStarted status
		typeCount := make(map[domain.TaskType]int)
		for _, task := range tasks {
			if task.ID == task1.ID || task.ID == task2.ID || task.ID == task3.ID {
				typeCount[task.Type]++
			}
		}

		if typeCount[domain.RegularTask] == 0 {
			t.Error("Expected at least one RegularTask in NotStarted")
		}
		if typeCount[domain.Bug] == 0 {
			t.Error("Expected at least one Bug task in NotStarted")
		}
		if typeCount[domain.Feature] == 0 {
			t.Error("Expected at least one Feature task in NotStarted")
		}
	})
}
