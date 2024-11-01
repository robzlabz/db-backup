package services

import (
	"time"

	"github.com/robzlabz/db-backup/internal/core/domain"
	"github.com/robzlabz/db-backup/internal/core/ports"
)

type backupService struct {
	repo          ports.BackupRepository
	mysqlBackuper ports.DatabaseBackuper
	pgBackuper    ports.DatabaseBackuper
}

func NewBackupService(
	repo ports.BackupRepository,
	mysqlBackuper ports.DatabaseBackuper,
	pgBackuper ports.DatabaseBackuper,
) ports.BackupService {
	return &backupService{
		repo:          repo,
		mysqlBackuper: mysqlBackuper,
		pgBackuper:    pgBackuper,
	}
}

func (s *backupService) AddConfig(config domain.BackupConfig) error {
	return s.repo.SaveConfig(config)
}

func (s *backupService) GetAllConfigs() ([]domain.BackupConfig, error) {
	return s.repo.GetAllConfigs()
}

func (s *backupService) ExecuteBackup(config domain.BackupConfig) error {
	// Cek apakah sudah waktunya backup
	now := time.Now().Unix()
	if now-config.LastBackup < int64(config.Interval*60) {
		return nil // Belum waktunya backup
	}

	var err error
	if config.Type == "MySQL" {
		err = s.mysqlBackuper.Backup(config)
	} else {
		err = s.pgBackuper.Backup(config)
	}

	if err != nil {
		return err
	}

	// Update waktu backup terakhir
	return s.repo.UpdateLastBackup(config.ID, now)
}
