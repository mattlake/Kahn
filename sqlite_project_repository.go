package main

import (
	"database/sql"
	"fmt"
	"log"
)

type SQLiteProjectRepository struct {
	db *sql.DB
}

func NewSQLiteProjectRepository(db *sql.DB) *SQLiteProjectRepository {
	return &SQLiteProjectRepository{db: db}
}

func (r *SQLiteProjectRepository) Create(project *Project) error {
	query := `
		INSERT INTO projects (id, name, description, color, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.Exec(query, project.ID, project.Name, project.Description,
		project.Color, project.CreatedAt, project.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create project: %w", err)
	}

	return nil
}

func (r *SQLiteProjectRepository) GetByID(id string) (*Project, error) {
	query := `
		SELECT id, name, description, color, created_at, updated_at
		FROM projects WHERE id = ?
	`

	var project Project
	err := r.db.QueryRow(query, id).Scan(
		&project.ID, &project.Name, &project.Description, &project.Color,
		&project.CreatedAt, &project.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get project: %w", err)
	}

	return &project, nil
}

func (r *SQLiteProjectRepository) GetAll() ([]Project, error) {
	query := `
		SELECT id, name, description, color, created_at, updated_at
		FROM projects ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query)
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

		projects = append(projects, project)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating projects: %w", err)
	}

	return projects, nil
}

func (r *SQLiteProjectRepository) Update(project *Project) error {
	query := `
		UPDATE projects 
		SET name = ?, description = ?, color = ?, updated_at = ?
		WHERE id = ?
	`

	_, err := r.db.Exec(query, project.Name, project.Description,
		project.Color, project.UpdatedAt, project.ID)
	if err != nil {
		return fmt.Errorf("failed to update project: %w", err)
	}

	return nil
}

func (r *SQLiteProjectRepository) Delete(id string) error {
	query := `DELETE FROM projects WHERE id = ?`

	result, err := r.db.Exec(query, id)
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

func (r *SQLiteProjectRepository) GetByIDWithTasks(id string) (*Project, error) {
	query := `
		SELECT id, name, description, color, created_at, updated_at
		FROM projects WHERE id = ?
	`

	var project Project
	err := r.db.QueryRow(query, id).Scan(
		&project.ID, &project.Name, &project.Description, &project.Color,
		&project.CreatedAt, &project.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get project: %w", err)
	}

	// Load tasks for this project
	taskRepo := NewSQLiteTaskRepository(r.db)
	tasks, err := taskRepo.GetByProjectID(id)
	if err != nil {
		log.Printf("Warning: failed to load tasks for project %s: %v", id, err)
		tasks = []Task{}
	}

	project.Tasks = tasks
	return &project, nil
}
