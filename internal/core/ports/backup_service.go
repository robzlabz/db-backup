package ports

import "github.com/robzlabz/db-backup/internal/core/domain"

type BackupService interface {
	AddConfig(config domain.BackupConfig) error
	GetAllConfigs() ([]domain.BackupConfig, error)
	ExecuteBackup(config domain.BackupConfig) error
}

type DatabaseBackuper interface {
	Backup(config domain.BackupConfig) error
}
