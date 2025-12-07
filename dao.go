package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

type ProjectDAO struct {
	db *sql.DB
}

type TaskDAO struct {
	db *sql.DB
}

// NewProjectDAO creates a new ProjectDAO
func NewProjectDAO(db *sql.DB) *ProjectDAO {
	return &ProjectDAO{db: db}
}

// NewTaskDAO creates a new TaskDAO
func NewTaskDAO(db *sql.DB) *TaskDAO {
	return &TaskDAO{db: db}
}

// ProjectDAO methods

// Create creates a new project in the database
func (p *ProjectDAO) Create(project *Project) error {
	query := `
		INSERT INTO projects (id, name, description, color, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	_, err := p.db.Exec(query, project.ID, project.Name, project.Description,
		project.Color, project.CreatedAt, project.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create project: %w", err)
	}

	return nil
}

// GetByID retrieves a project by ID
func (p *ProjectDAO) GetByID(id string) (*Project, error) {
	query := `
		SELECT id, name, description, color, created_at, updated_at
		FROM projects WHERE id = ?
	`

	var project Project
	err := p.db.QueryRow(query, id).Scan(
		&project.ID, &project.Name, &project.Description, &project.Color,
		&project.CreatedAt, &project.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("project not found: %s", id)
		}
		return nil, fmt.Errorf("failed to get project: %w", err)
	}

	// Load tasks for this project
	taskDAO := NewTaskDAO(p.db)
	tasks, err := taskDAO.GetByProjectID(id)
	if err != nil {
		log.Printf("Warning: failed to load tasks for project %s: %v", id, err)
		tasks = []Task{}
	}

	project.Tasks = tasks
	return &project, nil
}

// GetAll retrieves all projects
func (p *ProjectDAO) GetAll() ([]Project, error) {
	query := `
		SELECT id, name, description, color, created_at, updated_at
		FROM projects ORDER BY created_at DESC
	`

	rows, err := p.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get projects: %w", err)
	}
	defer rows.Close()

	var projects []Project
	for rows.Next() {
		var project Project
		err := rows.Scan(
			&project.ID, &project.Name, &project.Description, &project.Color,
			&project.CreatedAt, &project.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan project: %w", err)
		}

		// Load tasks for this project
		taskDAO := NewTaskDAO(p.db)
		tasks, err := taskDAO.GetByProjectID(project.ID)
		if err != nil {
			log.Printf("Warning: failed to load tasks for project %s: %v", project.ID, err)
			tasks = []Task{}
		}

		project.Tasks = tasks
		projects = append(projects, project)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating projects: %w", err)
	}

	return projects, nil
}

// Update updates an existing project
func (p *ProjectDAO) Update(project *Project) error {
	query := `
		UPDATE projects 
		SET name = ?, description = ?, color = ?, updated_at = ?
		WHERE id = ?
	`

	project.UpdatedAt = time.Now()
	_, err := p.db.Exec(query, project.Name, project.Description,
		project.Color, project.UpdatedAt, project.ID)
	if err != nil {
		return fmt.Errorf("failed to update project: %w", err)
	}

	return nil
}

// Delete deletes a project by ID
func (p *ProjectDAO) Delete(id string) error {
	query := `DELETE FROM projects WHERE id = ?`

	result, err := p.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete project: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("project not found: %s", id)
	}

	return nil
}

// TaskDAO methods

// Create creates a new task in the database
func (t *TaskDAO) Create(task *Task) error {
	query := `
		INSERT INTO tasks (id, project_id, name, desc, status, priority, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := t.db.Exec(query, task.ID, task.ProjectID, task.Name, task.Desc,
		task.Status, task.Priority, task.CreatedAt, task.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create task: %w", err)
	}

	return nil
}

// GetByID retrieves a task by ID
func (t *TaskDAO) GetByID(id string) (*Task, error) {
	query := `
		SELECT id, project_id, name, desc, status, priority, created_at, updated_at
		FROM tasks WHERE id = ?
	`

	var task Task
	err := t.db.QueryRow(query, id).Scan(
		&task.ID, &task.ProjectID, &task.Name, &task.Desc,
		&task.Status, &task.Priority, &task.CreatedAt, &task.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("task not found: %s", id)
		}
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	return &task, nil
}

// GetByProjectID retrieves all tasks for a project
func (t *TaskDAO) GetByProjectID(projectID string) ([]Task, error) {
	query := `
		SELECT id, project_id, name, desc, status, priority, created_at, updated_at
		FROM tasks WHERE project_id = ? ORDER BY created_at DESC
	`

	rows, err := t.db.Query(query, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks for project: %w", err)
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		err := rows.Scan(
			&task.ID, &task.ProjectID, &task.Name, &task.Desc,
			&task.Status, &task.Priority, &task.CreatedAt, &task.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}
		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating tasks: %w", err)
	}

	return tasks, nil
}

// GetByStatus retrieves all tasks for a project with a specific status
func (t *TaskDAO) GetByStatus(projectID string, status Status) ([]Task, error) {
	query := `
		SELECT id, project_id, name, desc, status, priority, created_at, updated_at
		FROM tasks WHERE project_id = ? AND status = ? ORDER BY created_at DESC
	`

	rows, err := t.db.Query(query, projectID, status)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks by status: %w", err)
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		err := rows.Scan(
			&task.ID, &task.ProjectID, &task.Name, &task.Desc,
			&task.Status, &task.Priority, &task.CreatedAt, &task.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}
		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating tasks: %w", err)
	}

	return tasks, nil
}

// Update updates an existing task
func (t *TaskDAO) Update(task *Task) error {
	query := `
		UPDATE tasks 
		SET name = ?, desc = ?, status = ?, priority = ?, updated_at = ?
		WHERE id = ?
	`

	task.UpdatedAt = time.Now()
	_, err := t.db.Exec(query, task.Name, task.Desc, task.Status,
		task.Priority, task.UpdatedAt, task.ID)
	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	return nil
}

// UpdateStatus updates only the status of a task
func (t *TaskDAO) UpdateStatus(taskID string, status Status) error {
	query := `
		UPDATE tasks 
		SET status = ?, updated_at = ?
		WHERE id = ?
	`

	_, err := t.db.Exec(query, status, time.Now(), taskID)
	if err != nil {
		return fmt.Errorf("failed to update task status: %w", err)
	}

	return nil
}

// Delete deletes a task by ID
func (t *TaskDAO) Delete(id string) error {
	query := `DELETE FROM tasks WHERE id = ?`

	result, err := t.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("task not found: %s", id)
	}

	return nil
}
