package ports

import "github.com/robzlabz/db-backup/internal/core/domain"

type BackupRepository interface {
	Save(backup *domain.Backup) error
	GetAll() ([]domain.Backup, error)
	GetByID(id int64) (*domain.Backup, error)
}

type DatabaseBackuper interface {
	Backup(config domain.BackupConfig) error
}
