package main

import (
	"kahn/internal/app"
	"kahn/internal/config"
	"kahn/internal/database"
	"testing"
)

func TestKahnModelCreation(t *testing.T) {
	// This is a simple test to verify the refactored structure works
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

	// Test that we can create a KahnModel
	model := app.NewKahnModel(db)
	if model == nil {
		t.Fatal("NewKahnModel should not return nil")
	}

	// Test basic functionality
	activeProject := model.GetActiveProject()
	if activeProject == nil {
		t.Fatal("Should have at least one project (default)")
	}

	// Test task creation
	err = model.CreateTask("Test Task", "Test Description")
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}
}
