package main

import (
	"testing"
)

func TestProjectService_CreateProject(t *testing.T) {
	// Setup
	taskRepo := &mockTaskRepository{}
	projectRepo := &mockProjectRepository{}
	service := NewProjectService(projectRepo, taskRepo)

	t.Run("successful project creation", func(t *testing.T) {
		// Act
		project, err := service.CreateProject("Test Project", "Test Description")

		// Assert
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if project == nil {
			t.Error("Expected project to be created")
		}
		if project.Name != "Test Project" {
			t.Errorf("Expected project name 'Test Project', got '%s'", project.Name)
		}
		if project.Description != "Test Description" {
			t.Errorf("Expected project description 'Test Description', got '%s'", project.Description)
		}
		if project.Color != "#89b4fa" {
			t.Errorf("Expected project color '#89b4fa', got '%s'", project.Color)
		}
	})

	t.Run("empty name validation", func(t *testing.T) {
		// Act
		project, err := service.CreateProject("", "Test Description")

		// Assert
		if err == nil {
			t.Error("Expected validation error for empty name")
		}
		if project != nil {
			t.Error("Expected no project to be created")
		}
	})
}

func TestProjectService_GetAllProjects(t *testing.T) {
	// Setup
	taskRepo := &mockTaskRepository{}
	projectRepo := &mockProjectRepository{}
	service := NewProjectService(projectRepo, taskRepo)

	// Setup test data
	testProjects := []Project{
		*NewProject("Project 1", "Description 1", "#89b4fa"),
		*NewProject("Project 2", "Description 2", "#89b4fa"),
	}
	projectRepo.projects = testProjects

	t.Run("get all projects", func(t *testing.T) {
		// Act
		projects, err := service.GetAllProjects()

		// Assert
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if len(projects) != 2 {
			t.Errorf("Expected 2 projects, got %d", len(projects))
		}
		if projects[0].Name != "Project 1" {
			t.Errorf("Expected first project name 'Project 1', got '%s'", projects[0].Name)
		}
	})
}

func TestProjectService_DeleteProject(t *testing.T) {
	// Setup
	taskRepo := &mockTaskRepository{}
	projectRepo := &mockProjectRepository{}
	service := NewProjectService(projectRepo, taskRepo)

	// Setup test data
	testProject := NewProject("Test Project", "Test Description", "#89b4fa")
	projectRepo.projects = []Project{*testProject}

	t.Run("successful project deletion", func(t *testing.T) {
		// Act
		err := service.DeleteProject(testProject.ID)

		// Assert
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("delete non-existent project", func(t *testing.T) {
		// Act
		err := service.DeleteProject("non-existent-id")

		// Assert
		if err == nil {
			t.Error("Expected error for non-existent project")
		}
	})
}
