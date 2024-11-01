package ports

import "github.com/robzlabz/db-backup/internal/core/domain"

type BackupService interface {
	CreateBackup(dbType string, config domain.BackupConfig) error
	GetAllBackups() ([]domain.Backup, error)
	GetBackup(id int64) (*domain.Backup, error)
}
