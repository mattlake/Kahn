package repository

import (
	"database/sql"
	"fmt"
	"kahn/internal/domain"
	"time"
)

type SQLiteTaskRepository struct {
	db *sql.DB
}

func NewSQLiteTaskRepository(db *sql.DB) *SQLiteTaskRepository {
	return &SQLiteTaskRepository{db: db}
}

func (r *SQLiteTaskRepository) Create(task *domain.Task) error {
	query := `
		INSERT INTO tasks (id, project_id, name, desc, status, priority, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.Exec(query, task.ID, task.ProjectID, task.Name, task.Desc,
		task.Status, task.Priority, task.CreatedAt, task.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create task: %w", err)
	}

	return nil
}

func (r *SQLiteTaskRepository) GetByID(id string) (*domain.Task, error) {
	query := `
		SELECT id, project_id, name, desc, status, priority, created_at, updated_at
		FROM tasks WHERE id = ?
	`

	var task domain.Task
	err := r.db.QueryRow(query, id).Scan(
		&task.ID, &task.ProjectID, &task.Name, &task.Desc,
		&task.Status, &task.Priority, &task.CreatedAt, &task.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	return &task, nil
}

func (r *SQLiteTaskRepository) GetByProjectID(projectID string) ([]domain.Task, error) {
	query := `
		SELECT id, project_id, name, desc, status, priority, created_at, updated_at
		FROM tasks WHERE project_id = ? ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks for project: %w", err)
	}
	defer rows.Close()

	var tasks []domain.Task
	for rows.Next() {
		var task domain.Task
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

func (r *SQLiteTaskRepository) GetByStatus(projectID string, status domain.Status) ([]domain.Task, error) {
	query := `
		SELECT id, project_id, name, desc, status, priority, created_at, updated_at
		FROM tasks WHERE project_id = ? AND status = ? ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, projectID, status)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks by status: %w", err)
	}
	defer rows.Close()

	var tasks []domain.Task
	for rows.Next() {
		var task domain.Task
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

func (r *SQLiteTaskRepository) Update(task *domain.Task) error {
	query := `
		UPDATE tasks 
		SET name = ?, desc = ?, status = ?, priority = ?, updated_at = ?
		WHERE id = ?
	`

	_, err := r.db.Exec(query, task.Name, task.Desc, task.Status,
		task.Priority, task.UpdatedAt, task.ID)
	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	return nil
}

func (r *SQLiteTaskRepository) UpdateStatus(taskID string, status domain.Status) error {
	query := `
		UPDATE tasks 
		SET status = ?, updated_at = ?
		WHERE id = ?
	`

	_, err := r.db.Exec(query, status, time.Now(), taskID)
	if err != nil {
		return fmt.Errorf("failed to update task status: %w", err)
	}

	return nil
}

func (r *SQLiteTaskRepository) Delete(id string) error {
	query := `DELETE FROM tasks WHERE id = ?`

	result, err := r.db.Exec(query, id)
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
