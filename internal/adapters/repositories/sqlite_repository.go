package repositories

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/robzlabz/db-backup/internal/core/domain"
	"github.com/robzlabz/db-backup/internal/core/ports"
	"github.com/robzlabz/db-backup/pkg/logging"
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
		logging.Errorf("[Repository][SQLiteRepository][SaveConfig] Error: %v", err)
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
		logging.Errorf("[Repository][SQLiteRepository][GetAllConfigs] Error: %v", err)
		return nil, err
	}

	return configs, nil
}

func (r *SQLiteRepository) UpdateLastBackup(id int, timestamp int64) error {
	query := `UPDATE backup_configs SET last_backup = ? WHERE id = ?`
	_, err := r.db.Exec(query, timestamp, id)
	if err != nil {
		logging.Errorf("[Repository][SQLiteRepository][UpdateLastBackup] Error: %v", err)
		return err
	}

	return nil
}

func (r *SQLiteRepository) Delete(id int) error {
	query := `DELETE FROM backup_configs WHERE id = ?`
	result, err := r.db.Exec(query, id)
	if err != nil {
		logging.Errorf("[Repository][SQLiteRepository][Delete] Error calling query: %v", err)
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		logging.Errorf("[Repository][SQLiteRepository][Delete] Error getting rows affected: %v", err)
		return err
	}

	if affected == 0 {
		logging.Errorf("[Repository][SQLiteRepository][Delete] No configuration found with ID %d", id)
		return fmt.Errorf("no configuration found with ID %d", id)
	}

	return nil
}
