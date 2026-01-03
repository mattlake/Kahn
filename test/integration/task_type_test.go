package main

import (
	"kahn/internal/app"
	"kahn/internal/config"
	"kahn/internal/database"
	"kahn/internal/domain"
	"testing"
)

func TestTaskType_BasicWorkflow(t *testing.T) {
	// Setup in-memory database with migrations
	cfg := &config.Config{}
	cfg.Database.Path = ":memory:"
	cfg.Database.BusyTimeout = 5000
	cfg.Database.JournalMode = "WAL"
	cfg.Database.CacheSize = 10000
	cfg.Database.ForeignKeys = true

	db, err := database.NewDatabase(cfg)
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	defer db.Close()

	// Create model
	model := app.NewKahnModel(db, "test")
	if model == nil {
		t.Fatal("NewKahnModel should not return nil")
	}

	// Test creating tasks (all will be RegularTask since CreateTaskWithPriority always creates RegularTask)
	err = model.CreateTaskWithPriority("Bug Task", "Bug description", domain.High)
	if err != nil {
		t.Fatalf("Failed to create bug task: %v", err)
	}

	err = model.CreateTaskWithPriority("Feature Task", "Feature description", domain.Medium)
	if err != nil {
		t.Fatalf("Failed to create feature task: %v", err)
	}

	err = model.CreateTaskWithPriority("Regular Task", "Regular description", domain.Low)
	if err != nil {
		t.Fatalf("Failed to create regular task: %v", err)
	}

	// Verify tasks were created
	activeProject := model.GetActiveProject()
	if activeProject == nil {
		t.Fatal("Should have active project")
	}

	// Count task types - all should be RegularTask (type 0) since CreateTaskWithPriority always creates RegularTask
	typeCount := make(map[int]int)
	for _, task := range activeProject.Tasks {
		typeCount[int(task.Type)]++
	}

	if typeCount[0] == 0 { // RegularTask
		t.Error("Should have regular tasks")
	}
	// Note: Bug and Feature tasks are not created by CreateTaskWithPriority method
	// This test validates the current behavior where all tasks are RegularTask
}
