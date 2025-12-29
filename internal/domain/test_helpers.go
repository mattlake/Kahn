package domain

import "time"

// createTestTask creates a test task with the given parameters
func createTestTask(name, description, projectID string, status Status) *Task {
	task := &Task{
		ID:        "test_task_" + time.Now().Format("20060102150405.000000000"),
		Name:      name,
		Desc:      description,
		ProjectID: projectID,
		Status:    status,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Priority:  Medium,
	}
	return task
}

// createTestProject creates a test project with the given parameters
func createTestProject(name, description, color string) *Project {
	project := &Project{
		ID:          "test_proj_" + time.Now().Format("20060102150405.000000000"),
		Name:        name,
		Description: description,
		Color:       color,
		Tasks:       []Task{},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	return project
}
