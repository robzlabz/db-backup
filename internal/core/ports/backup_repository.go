package ports

import "github.com/robzlabz/db-backup/internal/core/domain"

type BackupRepository interface {
	SaveConfig(config domain.BackupConfig) error
	GetAllConfigs() ([]domain.BackupConfig, error)
	UpdateLastBackup(id int, timestamp int64) error
	Delete(id int) error
}
