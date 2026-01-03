package repository

import (
	"database/sql"
	"fmt"
	"kahn/internal/domain"
)

// BaseRepository provides common database operations for repositories
type BaseRepository struct {
	db *sql.DB
}

// NewBaseRepository creates a new base repository with database connection
func NewBaseRepository(db *sql.DB) *BaseRepository {
	return &BaseRepository{db: db}
}

// WrapDBError standardizes error handling for database operations
func (b *BaseRepository) WrapDBError(operation, entity, id string, err error) error {
	if err == sql.ErrNoRows {
		return nil // Not found is not an error for Get operations
	}

	return &domain.RepositoryError{
		Operation: operation,
		Entity:    entity,
		ID:        id,
		Cause:     err,
	}
}

// HandleRowsAffected handles results for update/delete operations
func (b *BaseRepository) HandleRowsAffected(result sql.Result, operation, entity string) error {
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return b.WrapDBError("check rows affected", entity, "", err)
	}

	if rowsAffected == 0 {
		return &domain.RepositoryError{
			Operation: operation,
			Entity:    entity,
			ID:        "",
			Cause:     fmt.Errorf("no rows were affected"),
		}
	}

	return nil
}

// ScanTaskRows helper for scanning multiple task rows
func (b *BaseRepository) ScanTaskRows(rows *sql.Rows) ([]domain.Task, error) {
	var tasks []domain.Task
	for rows.Next() {
		var task domain.Task
		err := rows.Scan(
			&task.ID, &task.ProjectID, &task.Name, &task.Desc,
			&task.Status, &task.Type, &task.Priority, &task.CreatedAt, &task.UpdatedAt,
		)
		if err != nil {
			return nil, b.WrapDBError("scan", "task", "", err)
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, b.WrapDBError("iterate", "tasks", "", err)
	}

	return tasks, nil
}

// ScanProjectRows helper for scanning multiple project rows
func (b *BaseRepository) ScanProjectRows(rows *sql.Rows) ([]domain.Project, error) {
	var projects []domain.Project
	for rows.Next() {
		var project domain.Project
		err := rows.Scan(
			&project.ID, &project.Name, &project.Description, &project.Color,
			&project.CreatedAt, &project.UpdatedAt,
		)
		if err != nil {
			return nil, b.WrapDBError("scan", "project", "", err)
		}
		projects = append(projects, project)
	}

	if err := rows.Err(); err != nil {
		return nil, b.WrapDBError("iterate", "projects", "", err)
	}

	return projects, nil
}

// ScanSingleTask helper for scanning a single task row
func (b *BaseRepository) ScanSingleTask(row *sql.Row) (*domain.Task, error) {
	var task domain.Task
	err := row.Scan(
		&task.ID, &task.ProjectID, &task.Name, &task.Desc,
		&task.Status, &task.Type, &task.Priority, &task.CreatedAt, &task.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, b.WrapDBError("get", "task", "", err)
	}
	return &task, nil
}

// ScanSingleProject helper for scanning a single project row
func (b *BaseRepository) ScanSingleProject(row *sql.Row) (*domain.Project, error) {
	var project domain.Project
	err := row.Scan(
		&project.ID, &project.Name, &project.Description, &project.Color,
		&project.CreatedAt, &project.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, b.WrapDBError("get", "project", "", err)
	}
	return &project, nil
}
