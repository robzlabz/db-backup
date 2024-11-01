package repositories

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/robzlabz/db-backup/internal/core/domain"
	"github.com/robzlabz/db-backup/internal/core/ports"
)

type SQLiteRepository struct {
	db *sqlx.DB
}

func NewSQLiteRepository(db *sqlx.DB) ports.BackupRepository {
	return &SQLiteRepository{db: db}
}

func (r *SQLiteRepository) SaveConfig(config domain.BackupConfig) error {
	query := `
		INSERT INTO backup_configs (
			type, host, port, database, username, password,
			interval, output_path, last_backup
		) VALUES (
			:type, :host, :port, :database, :user, :password,
			:interval, :output_path, :last_backup
		)
	`
	_, err := r.db.NamedExec(query, config)
	if err != nil {
		log.Printf("[Repository][SQLiteRepository][SaveConfig] Error: %v", err)
		return err
	}

	return nil
}

func (r *SQLiteRepository) GetAllConfigs() ([]domain.BackupConfig, error) {
	var configs []domain.BackupConfig
	query := `
		SELECT id, type, host, port, database, username as user,
		       password, interval, output_path, last_backup
		FROM backup_configs
	`
	err := r.db.Select(&configs, query)
	if err != nil {
		log.Printf("[Repository][SQLiteRepository][GetAllConfigs] Error: %v", err)
		return nil, err
	}

	return configs, nil
}

func (r *SQLiteRepository) UpdateLastBackup(id int, timestamp int64) error {
	query := `UPDATE backup_configs SET last_backup = ? WHERE id = ?`
	_, err := r.db.Exec(query, timestamp, id)
	if err != nil {
		log.Printf("[Repository][SQLiteRepository][UpdateLastBackup] Error: %v", err)
		return err
	}

	return nil
}

// InitDB membuat tabel jika belum ada
func (r *SQLiteRepository) InitDB() error {
	schema := `
	CREATE TABLE IF NOT EXISTS backup_configs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		type TEXT NOT NULL,
		host TEXT NOT NULL,
		port INTEGER NOT NULL,
		database TEXT NOT NULL,
		username TEXT NOT NULL,
		password TEXT NOT NULL,
		interval INTEGER NOT NULL,
		output_path TEXT NOT NULL,
		last_backup INTEGER DEFAULT 0
	);`

	_, err := r.db.Exec(schema)
	if err != nil {
		log.Printf("[Repository][SQLiteRepository][InitDB] Error: %v", err)
		return err
	}

	return nil
}

func (r *SQLiteRepository) Delete(id int) error {
	query := `DELETE FROM backup_configs WHERE id = ?`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return fmt.Errorf("no configuration found with ID %s", id)
	}

	return nil
}
