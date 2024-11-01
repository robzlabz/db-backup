package backupers

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/robzlabz/db-backup/internal/core/domain"
	"github.com/robzlabz/db-backup/pkg/logging"
	"github.com/robzlabz/db-backup/pkg/utils"
)

type mysqlBackuper struct{}

func NewMySQLBackuper() *mysqlBackuper {
	return &mysqlBackuper{}
}

func (b *mysqlBackuper) Backup(config domain.BackupConfig) error {
	logger := logging.Sugar()
	host := config.Host
	if host == "localhost" {
		host = "127.0.0.1"
	}

	logger.Infow("[Backuper][MySQLBackuper] Memulai backup MySQL",
		"database", config.Database,
		"host", host,
		"port", config.Port,
	)

	// Membuat nama file backup dengan timestamp
	timestamp := time.Now().Format("20060102150405")
	filename := fmt.Sprintf("%s_%s.sql", config.Database, timestamp)
	backupPath := filepath.Join(config.OutputPath, filename)
	zipPath := backupPath + ".zip"

	// Memastikan direktori backup ada
	if err := os.MkdirAll(config.OutputPath, 0755); err != nil {
		logger.Errorw("[Backuper][MySQLBackuper] Gagal membuat direktori backup",
			"error", err,
			"path", config.OutputPath,
		)
		return fmt.Errorf("gagal membuat direktori backup: %v", err)
	}

	args := []string{
		"-h", host,
		"-P", fmt.Sprintf("%d", config.Port),
		"-u", config.User,
	}

	if config.Password != "" {
		args = append(args, fmt.Sprintf("-p%s", config.Password))
	}

	if config.Database != "" {
		args = append(args, config.Database)
	} else {
		return fmt.Errorf("database name is empty")
	}

	cmd := exec.Command("mysqldump", args...)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()
	cmd = exec.CommandContext(ctx, cmd.Path, cmd.Args[1:]...)

	logger.Debug("[Backuper][MySQLBackuper] Menjalankan mysqldump command")

	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Errorw("[Backuper][MySQLBackuper] Gagal menjalankan mysqldump",
			"error", err,
			"output", string(output),
		)
		return fmt.Errorf("backup error: %v", err)
	}

	// Tulis output ke file SQL
	if err := os.WriteFile(backupPath, output, 0644); err != nil {
		logger.Errorw("[Backuper][MySQLBackuper] Gagal menulis file backup",
			"error", err,
			"path", backupPath,
		)
		return err
	}

	// Kompresi file backup
	logger.Debug("[Backuper][MySQLBackuper] Memulai kompresi file backup")
	if err := utils.CompressFile(backupPath, zipPath); err != nil {
		logger.Errorw("[Backuper][MySQLBackuper] Gagal mengkompresi file backup",
			"error", err,
			"source", backupPath,
			"destination", zipPath,
		)
		return fmt.Errorf("gagal mengkompresi file backup: %v", err)
	}

	// Get file info untuk ukuran file
	fileInfo, err := os.Stat(zipPath)
	if err != nil {
		logger.Warnw("[Backuper][MySQLBackuper] Gagal mendapatkan informasi file backup",
			"error", err,
			"path", zipPath,
		)
	} else {
		logger.Infow("[Backuper][MySQLBackuper] Backup selesai",
			"database", config.Database,
			"file", zipPath,
			"size_bytes", fileInfo.Size(),
			"duration", time.Since(time.Now()),
		)
	}

	return nil
}
