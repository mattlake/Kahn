package repository

import (
	"database/sql"
	"kahn/internal/domain"
	"time"
)

type SQLiteTaskRepository struct {
	base *BaseRepository // Composition, not embedding
}

func NewSQLiteTaskRepository(db *sql.DB) *SQLiteTaskRepository {
	return &SQLiteTaskRepository{
		base: NewBaseRepository(db), // Composition
	}
}

func (r *SQLiteTaskRepository) Create(task *domain.Task) error {
	query := `
		INSERT INTO tasks (id, project_id, name, desc, status, type, priority, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	return r.base.CreateGeneric(query, task.ID, task.ProjectID, task.Name, task.Desc,
		task.Status, task.Type, task.Priority, task.CreatedAt, task.UpdatedAt)
}

func (r *SQLiteTaskRepository) GetByID(id string) (*domain.Task, error) {
	query := `
		SELECT id, project_id, name, desc, status, type, priority, created_at, updated_at
		FROM tasks WHERE id = ?
	`

	row := r.base.db.QueryRow(query, id)
	return r.base.ScanSingleTask(row)
}

func (r *SQLiteTaskRepository) GetByProjectID(projectID string) ([]domain.Task, error) {
	query := `
		SELECT id, project_id, name, desc, status, type, priority, created_at, updated_at
		FROM tasks WHERE project_id = ? ORDER BY created_at DESC
	`

	rows, err := r.base.db.Query(query, projectID)
	if err != nil {
		return nil, r.base.WrapDBError("get", "tasks for project", projectID, err)
	}
	defer rows.Close()

	return r.base.ScanTaskRows(rows)
}

func (r *SQLiteTaskRepository) GetByStatus(projectID string, status domain.Status) ([]domain.Task, error) {
	var query string

	// Different ordering based on status
	if status == domain.NotStarted {
		// Not Started: priority DESC, then created_at ASC (oldest highest priority first)
		query = `
			SELECT id, project_id, name, desc, status, type, priority, created_at, updated_at
			FROM tasks WHERE project_id = ? AND status = ? 
			ORDER BY priority DESC, created_at ASC
		`
	} else {
		// In Progress and Done: updated_at DESC (newest changes first)
		query = `
			SELECT id, project_id, name, desc, status, type, priority, created_at, updated_at
			FROM tasks WHERE project_id = ? AND status = ? 
			ORDER BY updated_at DESC
		`
	}

	rows, err := r.base.db.Query(query, projectID, status)
	if err != nil {
		return nil, r.base.WrapDBError("get", "tasks by status", "", err)
	}
	defer rows.Close()

	return r.base.ScanTaskRows(rows)
}

func (r *SQLiteTaskRepository) Update(task *domain.Task) error {
	query := `
		UPDATE tasks 
		SET name = ?, desc = ?, status = ?, type = ?, priority = ?, updated_at = ?
		WHERE id = ?
	`

	task.UpdatedAt = time.Now()
	_, err := r.base.db.Exec(query, task.Name, task.Desc, task.Status,
		task.Type, task.Priority, task.UpdatedAt, task.ID)
	if err != nil {
		return r.base.WrapDBError("update", "task", task.ID, err)
	}
	return nil
}

func (r *SQLiteTaskRepository) UpdateStatus(taskID string, status domain.Status) error {
	query := `
		UPDATE tasks 
		SET status = ?, updated_at = ?
		WHERE id = ?
	`

	_, err := r.base.db.Exec(query, status, time.Now(), taskID)
	if err != nil {
		return r.base.WrapDBError("update", "task status", taskID, err)
	}
	return nil
}

func (r *SQLiteTaskRepository) Delete(id string) error {
	query := `DELETE FROM tasks WHERE id = ?`
	return r.base.DeleteGeneric(query, id)
}
