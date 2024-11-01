package repositories

import (
	"database/sql"

	"github.com/robzlabz/db-backup/internal/core/domain"
)

type sqliteRepository struct {
	db *sql.DB
}

func NewSQLiteRepository(db *sql.DB) *sqliteRepository {
	return &sqliteRepository{db: db}
}

func (r *sqliteRepository) Save(backup *domain.Backup) error {
	query := `INSERT INTO backups (db_type, db_name, file_path, size, created_at)
			  VALUES (?, ?, ?, ?, ?)`

	result, err := r.db.Exec(query,
		backup.DBType,
		backup.DBName,
		backup.FilePath,
		backup.Size,
		backup.CreatedAt,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	backup.ID = id
	return nil
}

// Implementasi method lainnya...
