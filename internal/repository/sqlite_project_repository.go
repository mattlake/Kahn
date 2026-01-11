package repository

import (
	"database/sql"
	"kahn/internal/domain"
	"time"
)

type SQLiteProjectRepository struct {
	base *BaseRepository // Composition, not embedding
}

func NewSQLiteProjectRepository(db *sql.DB) *SQLiteProjectRepository {
	return &SQLiteProjectRepository{
		base: NewBaseRepository(db), // Composition
	}
}

func (r *SQLiteProjectRepository) Create(project *domain.Project) error {
	query := `
		INSERT INTO projects (id, name, description, color, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	return r.base.CreateGeneric(query, project.ID, project.Name, project.Description,
		project.Color, project.CreatedAt, project.UpdatedAt)
}

func (r *SQLiteProjectRepository) GetByID(id string) (*domain.Project, error) {
	query := `
		SELECT id, name, description, color, created_at, updated_at
		FROM projects WHERE id = ?
	`

	row := r.base.db.QueryRow(query, id)
	return r.base.ScanSingleProject(row)
}

func (r *SQLiteProjectRepository) GetAll() ([]domain.Project, error) {
	query := `
		SELECT id, name, description, color, created_at, updated_at
		FROM projects ORDER BY created_at DESC
	`

	rows, err := r.base.db.Query(query)
	if err != nil {
		return nil, r.base.WrapDBError("get", "projects", "", err)
	}
	defer rows.Close()

	return r.base.ScanProjectRows(rows)
}

func (r *SQLiteProjectRepository) Update(project *domain.Project) error {
	query := `
		UPDATE projects 
		SET name = ?, description = ?, color = ?, updated_at = ?
		WHERE id = ?
	`

	project.UpdatedAt = time.Now()
	_, err := r.base.db.Exec(query, project.Name, project.Description,
		project.Color, project.UpdatedAt, project.ID)
	if err != nil {
		return r.base.WrapDBError("update", "project", project.ID, err)
	}
	return nil
}

func (r *SQLiteProjectRepository) Delete(id string) error {
	query := `DELETE FROM projects WHERE id = ?`
	return r.base.DeleteGeneric(query, id)
}
