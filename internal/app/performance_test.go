package app

import (
	"testing"

	"kahn/internal/domain"
	"kahn/internal/services"
	"kahn/internal/ui/styles"

	"github.com/charmbracelet/bubbles/list"
)

// BenchmarkUpdateTaskLists measures the performance of task list updates
// This benchmark demonstrates the impact of our database query optimization (3 DB calls → 1 DB call)
func BenchmarkUpdateTaskLists(b *testing.B) {
	// Setup test data
	mockTaskRepo := services.NewMockTaskRepository()
	mockProjectRepo := services.NewMockProjectRepository()
	taskService := services.NewTaskService(mockTaskRepo, mockProjectRepo)

	// Create test project
	project := domain.NewProject("Test Project", "Test Description", "blue")
	mockProjectRepo.Create(project)

	// Create test tasks
	for i := 0; i < 50; i++ {
		task := domain.NewTask("Task ", "Description", project.ID)
		task.Status = domain.Status(i % 3)
		mockTaskRepo.Create(task)
	}

	// Setup navigation state
	taskLists := []list.Model{
		list.New(nil, list.NewDefaultDelegate(), 0, 0),
		list.New(nil, list.NewDefaultDelegate(), 0, 0),
		list.New(nil, list.NewDefaultDelegate(), 0, 0),
	}
	navState := NewNavigationState(taskLists)

	b.ResetTimer()
	b.ReportAllocs()

	// Benchmark the optimized UpdateTaskLists method
	for i := 0; i < b.N; i++ {
		navState.UpdateTaskLists(project, taskService)
	}
}

// BenchmarkUpdateDirtyLists measures the performance of incremental updates
// This benchmark demonstrates the impact of our dirty flag optimization (full rebuild → targeted update)
func BenchmarkUpdateDirtyLists(b *testing.B) {
	// Setup test data
	mockTaskRepo := services.NewMockTaskRepository()
	mockProjectRepo := services.NewMockProjectRepository()
	taskService := services.NewTaskService(mockTaskRepo, mockProjectRepo)

	// Create test project
	project := domain.NewProject("Test Project", "Test Description", "blue")
	mockProjectRepo.Create(project)

	// Create test tasks
	for i := 0; i < 50; i++ {
		task := domain.NewTask("Task ", "Description", project.ID)
		task.Status = domain.Status(i % 3)
		mockTaskRepo.Create(task)
	}

	// Setup navigation state
	taskLists := []list.Model{
		list.New(nil, list.NewDefaultDelegate(), 0, 0),
		list.New(nil, list.NewDefaultDelegate(), 0, 0),
		list.New(nil, list.NewDefaultDelegate(), 0, 0),
	}
	navState := NewNavigationState(taskLists)

	// Initial full update
	navState.UpdateTaskLists(project, taskService)

	b.ResetTimer()
	b.ReportAllocs()

	// Benchmark incremental updates (marking only one list dirty)
	for i := 0; i < b.N; i++ {
		navState.MarkListDirty(domain.NotStarted) // Only one list dirty
		navState.UpdateDirtyLists(project, taskService)
	}
}

// BenchmarkTaskTitleRendering measures the performance of task title rendering
// This benchmark demonstrates the impact of our style caching optimization (recreated styles → cached styles)
func BenchmarkTaskTitleRendering(b *testing.B) {
	// Create test task
	task := domain.NewTask("Test Task", "Test Description", "proj_123")
	task.Priority = domain.High

	// Create task wrapper
	taskWrapper := styles.NewTaskWithTitle(*task)

	b.ResetTimer()
	b.ReportAllocs()

	// Benchmark the Title() method (optimized with cached styles)
	for i := 0; i < b.N; i++ {
		_ = taskWrapper.Title()
	}
}
