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
		task, err := service.CreateTask("Test Task", "Test Description", testProject.ID, domain.RegularTask, domain.Low, nil)

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
		task, err := service.CreateTask("", "Test Description", "project-123", domain.RegularTask, domain.Low, nil)

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
		task, err := service.CreateTask("Test Task", "Test Description", "", domain.RegularTask, domain.Low, nil)

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

func TestTaskService_CreateTask_WithDifferentTypes(t *testing.T) {
	// Setup
	taskRepo := NewMockTaskRepository()
	projectRepo := NewMockProjectRepository()
	testProject := domain.NewProject("Test Project", "Test Description", "blue")
	projectRepo.Create(testProject)
	service := NewTaskService(taskRepo, projectRepo)

	tests := []struct {
		name        string
		taskName    string
		description string
		taskType    domain.TaskType
		priority    domain.Priority
	}{
		{"create task with RegularTask type", "Regular Task", "Description", domain.RegularTask, domain.Low},
		{"create task with Bug type", "Bug Task", "Bug Description", domain.Bug, domain.High},
		{"create task with Feature type", "Feature Task", "Feature Description", domain.Feature, domain.Medium},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			task, err := service.CreateTask(tt.taskName, tt.description, testProject.ID, tt.taskType, tt.priority, nil)

			// Assert
			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
			if task == nil {
				t.Error("Expected task to be created")
			}
			if task.Type != tt.taskType {
				t.Errorf("Expected task type %v, got %v", tt.taskType, task.Type)
			}
		})
	}
}

func TestTaskService_UpdateTask_TypeChange(t *testing.T) {
	// Setup
	taskRepo := NewMockTaskRepository()
	projectRepo := NewMockProjectRepository()
	testProject := domain.NewProject("Test Project", "Test Description", "blue")
	projectRepo.Create(testProject)
	service := NewTaskService(taskRepo, projectRepo)

	// Create original task
	originalTask, err := service.CreateTask("Original Task", "Description", testProject.ID, domain.RegularTask, domain.Low, nil)
	if err != nil {
		t.Fatalf("Failed to create original task: %v", err)
	}

	// Act
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
}

func TestTaskService_GetTasksByProject_TypePreservation(t *testing.T) {
	// Setup
	taskRepo := NewMockTaskRepository()
	projectRepo := NewMockProjectRepository()
	testProject := domain.NewProject("Test Project", "Test Description", "blue")
	projectRepo.Create(testProject)
	service := NewTaskService(taskRepo, projectRepo)

	// Create tasks with different types
	task1, _ := service.CreateTask("Task 1", "Regular", testProject.ID, domain.RegularTask, domain.Low, nil)
	task2, _ := service.CreateTask("Task 2", "Bug", testProject.ID, domain.Bug, domain.Medium, nil)
	task3, _ := service.CreateTask("Task 3", "Feature", testProject.ID, domain.Feature, domain.High, nil)

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
}

func TestTaskService_GetTasksByStatus_TypePreservation(t *testing.T) {
	// Setup
	taskRepo := NewMockTaskRepository()
	projectRepo := NewMockProjectRepository()
	testProject := domain.NewProject("Test Project", "Test Description", "blue")
	projectRepo.Create(testProject)
	service := NewTaskService(taskRepo, projectRepo)

	// Create tasks with different types and same status
	task1, _ := service.CreateTask("Not Started Regular", "Regular", testProject.ID, domain.RegularTask, domain.Low, nil)
	task2, _ := service.CreateTask("Not Started Bug", "Bug", testProject.ID, domain.Bug, domain.Medium, nil)
	task3, _ := service.CreateTask("Not Started Feature", "Feature", testProject.ID, domain.Feature, domain.High, nil)

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
}

func TestTaskService_AutoUnblock(t *testing.T) {
	// Setup
	taskRepo := NewMockTaskRepository()
	projectRepo := NewMockProjectRepository()
	testProject := domain.NewProject("Test Project", "Test Description", "blue")
	projectRepo.Create(testProject)
	service := NewTaskService(taskRepo, projectRepo)

	t.Run("moving blocker to Done unblocks dependent task", func(t *testing.T) {
		// Create Task A (blocker) and Task B (blocked by A)
		taskA, _ := service.CreateTask("Task A", "Blocker task", testProject.ID, domain.RegularTask, domain.Medium, nil)
		taskB, _ := service.CreateTask("Task B", "Blocked task", testProject.ID, domain.RegularTask, domain.Low, &taskA.IntID)

		// Verify Task B is blocked
		if taskB.BlockedBy == nil || *taskB.BlockedBy != taskA.IntID {
			t.Error("Task B should be blocked by Task A")
		}

		// Move Task A to InProgress
		service.MoveTaskToNextStatus(taskA.ID)

		// Verify Task B is still blocked (blocker not Done yet)
		taskB, _ = service.GetTask(taskB.ID)
		if taskB.BlockedBy == nil {
			t.Error("Task B should still be blocked (Task A is InProgress)")
		}

		// Move Task A to Done
		service.MoveTaskToNextStatus(taskA.ID)

		// Verify Task B is now unblocked
		taskB, _ = service.GetTask(taskB.ID)
		if taskB.BlockedBy != nil {
			t.Error("Task B should be unblocked after Task A moved to Done")
		}
	})

	t.Run("multiple blocked tasks are unblocked", func(t *testing.T) {
		// Create Task X (blocker) and Tasks Y, Z (both blocked by X)
		taskX, _ := service.CreateTask("Task X", "Blocker", testProject.ID, domain.RegularTask, domain.High, nil)
		taskY, _ := service.CreateTask("Task Y", "Blocked 1", testProject.ID, domain.RegularTask, domain.Low, &taskX.IntID)
		taskZ, _ := service.CreateTask("Task Z", "Blocked 2", testProject.ID, domain.RegularTask, domain.Low, &taskX.IntID)

		// Move Task X to Done
		service.MoveTaskToNextStatus(taskX.ID)
		service.MoveTaskToNextStatus(taskX.ID)

		// Verify both Y and Z are unblocked
		taskY, _ = service.GetTask(taskY.ID)
		taskZ, _ = service.GetTask(taskZ.ID)

		if taskY.BlockedBy != nil {
			t.Error("Task Y should be unblocked")
		}
		if taskZ.BlockedBy != nil {
			t.Error("Task Z should be unblocked")
		}
	})

	t.Run("moving blocker back from Done does not re-block", func(t *testing.T) {
		// Create Task C (blocker) and Task D (blocked by C)
		taskC, _ := service.CreateTask("Task C", "Blocker", testProject.ID, domain.RegularTask, domain.Medium, nil)
		taskD, _ := service.CreateTask("Task D", "Blocked", testProject.ID, domain.RegularTask, domain.Low, &taskC.IntID)

		// Move Task C to Done (unblocks Task D)
		service.MoveTaskToNextStatus(taskC.ID)
		service.MoveTaskToNextStatus(taskC.ID)

		// Verify Task D is unblocked
		taskD, _ = service.GetTask(taskD.ID)
		if taskD.BlockedBy != nil {
			t.Error("Task D should be unblocked")
		}

		// Move Task C back to InProgress
		service.MoveTaskToPreviousStatus(taskC.ID)

		// Verify Task D remains unblocked (no re-blocking)
		taskD, _ = service.GetTask(taskD.ID)
		if taskD.BlockedBy != nil {
			t.Error("Task D should remain unblocked (no re-blocking)")
		}
	})

	t.Run("cascade scenario - only direct dependents unblocked", func(t *testing.T) {
		// Create Task E, F, G where E blocks F, F blocks G
		taskE, _ := service.CreateTask("Task E", "Root blocker", testProject.ID, domain.RegularTask, domain.High, nil)
		taskF, _ := service.CreateTask("Task F", "Intermediate", testProject.ID, domain.RegularTask, domain.Medium, &taskE.IntID)
		taskG, _ := service.CreateTask("Task G", "Final", testProject.ID, domain.RegularTask, domain.Low, &taskF.IntID)

		// Move Task E to Done
		service.MoveTaskToNextStatus(taskE.ID)
		service.MoveTaskToNextStatus(taskE.ID)

		// Verify Task F is unblocked (direct dependent)
		taskF, _ = service.GetTask(taskF.ID)
		if taskF.BlockedBy != nil {
			t.Error("Task F should be unblocked (direct dependent of E)")
		}

		// Verify Task G is still blocked by F (no cascade)
		taskG, _ = service.GetTask(taskG.ID)
		if taskG.BlockedBy == nil || *taskG.BlockedBy != taskF.IntID {
			t.Error("Task G should still be blocked by F (no cascade unblocking)")
		}
	})

	t.Run("using backwards movement to Done also unblocks", func(t *testing.T) {
		// Create Task H (blocker) and Task I (blocked by H)
		taskH, _ := service.CreateTask("Task H", "Blocker", testProject.ID, domain.RegularTask, domain.Medium, nil)
		taskI, _ := service.CreateTask("Task I", "Blocked", testProject.ID, domain.RegularTask, domain.Low, &taskH.IntID)

		// Verify Task I is blocked
		if taskI.BlockedBy == nil || *taskI.BlockedBy != taskH.IntID {
			t.Error("Task I should be blocked by Task H")
		}

		// Move Task H backwards from NotStarted to Done (cycle)
		service.MoveTaskToPreviousStatus(taskH.ID)

		// Verify Task I is now unblocked
		taskI, _ = service.GetTask(taskI.ID)
		if taskI.BlockedBy != nil {
			t.Error("Task I should be unblocked after Task H moved to Done via backwards movement")
		}
	})

	t.Run("deleting blocker unblocks dependent task", func(t *testing.T) {
		// Create Task J (blocker) and Task K (blocked by J)
		taskJ, _ := service.CreateTask("Task J", "Blocker to be deleted", testProject.ID, domain.RegularTask, domain.High, nil)
		taskK, _ := service.CreateTask("Task K", "Blocked by J", testProject.ID, domain.RegularTask, domain.Low, &taskJ.IntID)

		// Verify Task K is blocked
		if taskK.BlockedBy == nil || *taskK.BlockedBy != taskJ.IntID {
			t.Error("Task K should be blocked by Task J")
		}

		// Delete Task J (the blocker)
		service.DeleteTask(taskJ.ID)

		// Verify Task K is now unblocked
		taskK, _ = service.GetTask(taskK.ID)
		if taskK.BlockedBy != nil {
			t.Error("Task K should be unblocked after Task J was deleted")
		}
	})
}
