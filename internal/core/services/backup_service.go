package services

import (
	"fmt"
	"time"

	"github.com/robzlabz/db-backup/internal/core/domain"
	"github.com/robzlabz/db-backup/internal/core/ports"
)

type backupService struct {
	repository    ports.BackupRepository
	mysqlBackuper ports.DatabaseBackuper
	pgBackuper    ports.DatabaseBackuper
}

func (s *backupService) GetAllBackups() ([]domain.Backup, error) {
	return s.repository.GetAll()
}

func (s *backupService) GetBackup(id int64) (*domain.Backup, error) {
	return s.repository.GetByID(id)
}

func (s *backupService) CreateBackup(dbType string, config domain.BackupConfig) error {
	var backuper ports.DatabaseBackuper

	switch dbType {
	case "mysql":
		backuper = s.mysqlBackuper
	case "postgres":
		backuper = s.pgBackuper
	default:
		return fmt.Errorf("unsupported database type: %s", dbType)
	}

	if err := backuper.Backup(config); err != nil {
		return err
	}

	backup := &domain.Backup{
		DBType:    dbType,
		DBName:    config.Database,
		FilePath:  config.OutputPath,
		CreatedAt: time.Now(),
	}

	return s.repository.Save(backup)
}

func NewBackupService(
	repo ports.BackupRepository,
	mysqlBackuper ports.DatabaseBackuper,
	pgBackuper ports.DatabaseBackuper,
) ports.BackupService {
	return &backupService{
		repository:    repo,
		mysqlBackuper: mysqlBackuper,
		pgBackuper:    pgBackuper,
	}
}
