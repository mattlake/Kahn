package repository

import (
	"database/sql"
	"fmt"
	"kahn/internal/domain"
	"time"
)

type BaseRepository struct {
	db *sql.DB
}

func NewBaseRepository(db *sql.DB) *BaseRepository {
	return &BaseRepository{db: db}
}

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

func (b *BaseRepository) ScanTaskRows(rows *sql.Rows) ([]domain.Task, error) {
	var tasks []domain.Task
	for rows.Next() {
		var task domain.Task
		err := rows.Scan(
			&task.IntID, &task.ID, &task.ProjectID, &task.Name, &task.Desc,
			&task.Status, &task.Type, &task.Priority, &task.BlockedBy,
			&task.CreatedAt, &task.UpdatedAt,
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

func (b *BaseRepository) ScanSingleTask(row *sql.Row) (*domain.Task, error) {
	var task domain.Task
	err := row.Scan(
		&task.IntID, &task.ID, &task.ProjectID, &task.Name, &task.Desc,
		&task.Status, &task.Type, &task.Priority, &task.BlockedBy,
		&task.CreatedAt, &task.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, b.WrapDBError("get", "task", "", err)
	}
	return &task, nil
}

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

// Generic CRUD operations for extending repositories

// CreateGeneric executes a generic insert query with the provided parameters
func (b *BaseRepository) CreateGeneric(query string, values ...interface{}) error {
	_, err := b.db.Exec(query, values...)
	if err != nil {
		return b.WrapDBError("create", "entity", "", err)
	}
	return nil
}

// UpdateGeneric executes a generic update query with the provided parameters
func (b *BaseRepository) UpdateGeneric(query string, values ...interface{}) error {
	_, err := b.db.Exec(query, values...)
	if err != nil {
		return b.WrapDBError("update", "entity", "", err)
	}
	return nil
}

// UpdateTimestampedGeneric executes a generic update query that updates timestamps
func (b *BaseRepository) UpdateTimestampedGeneric(query string, values ...interface{}) error {
	// Add updated_at timestamp if not already included
	valuesWithTime := append(values, time.Now())
	_, err := b.db.Exec(query, valuesWithTime...)
	if err != nil {
		return b.WrapDBError("update", "entity", "", err)
	}
	return nil
}

// DeleteGeneric executes a generic delete query with the provided parameters
func (b *BaseRepository) DeleteGeneric(query string, values ...interface{}) error {
	result, err := b.db.Exec(query, values...)
	if err != nil {
		return b.WrapDBError("delete", "entity", "", err)
	}
	return b.HandleRowsAffected(result, "delete", "entity")
}

// QueryAndScanGeneric executes a generic query and scans results using the provided scan function
func (b *BaseRepository) QueryAndScanGeneric(query string, scanFunc func(*sql.Rows) (interface{}, error), values ...interface{}) (interface{}, error) {
	rows, err := b.db.Query(query, values...)
	if err != nil {
		return nil, b.WrapDBError("query", "entity", "", err)
	}
	defer rows.Close()

	return scanFunc(rows)
}

// QueryRowAndScanGeneric executes a generic single-row query and scans result
func (b *BaseRepository) QueryRowAndScanGeneric(query string, scanFunc func(*sql.Row) (interface{}, error), values ...interface{}) (interface{}, error) {
	row := b.db.QueryRow(query, values...)
	return scanFunc(row)
}
