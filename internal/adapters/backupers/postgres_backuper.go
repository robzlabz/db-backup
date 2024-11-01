package backupers

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"github.com/robzlabz/db-backup/internal/core/domain"
	"github.com/robzlabz/db-backup/pkg/logging"
	"github.com/robzlabz/db-backup/pkg/utils"
)

type PostgresBackuper struct{}

func NewPostgresBackuper() *PostgresBackuper {
	return &PostgresBackuper{}
}

func (b *PostgresBackuper) Backup(config domain.BackupConfig) error {
	logger := logging.Sugar()
	logger.Infow("Memulai backup PostgreSQL",
		"database", config.Database,
		"host", config.Host,
		"port", config.Port,
	)

	// Membuat nama file backup dengan timestamp
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	filename := fmt.Sprintf("%s_%s.sql", config.Database, timestamp)
	backupPath := filepath.Join(config.OutputPath, filename)
	zipPath := backupPath + ".zip"

	// Memastikan direktori backup ada
	if err := os.MkdirAll(config.OutputPath, 0755); err != nil {
		logger.Errorw("Gagal membuat direktori backup",
			"error", err,
			"path", config.OutputPath,
		)
		return fmt.Errorf("gagal membuat direktori backup: %v", err)
	}

	logger.Debugw("Menyiapkan backup",
		"output_file", backupPath,
	)

	// Menyiapkan environment variable untuk password
	env := os.Environ()
	env = append(env, fmt.Sprintf("PGPASSWORD=%s", config.Password))

	// Menyiapkan command pg_dump
	cmd := exec.Command("pg_dump",
		"-h", config.Host,
		"-p", strconv.Itoa(config.Port),
		"-U", config.User,
		"-d", config.Database,
		"-F", "p", // Format plain text SQL
		"-f", backupPath,
	)

	cmd.Env = env

	logger.Debug("Menjalankan pg_dump command")

	// Menjalankan backup
	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Errorw("Gagal menjalankan pg_dump",
			"error", err,
			"output", string(output),
		)
		return fmt.Errorf("gagal melakukan backup: %v, output: %s", err, string(output))
	}

	// Kompresi file backup
	logger.Debug("Memulai kompresi file backup")
	if err := utils.CompressFile(backupPath, zipPath); err != nil {
		logger.Errorw("Gagal mengkompresi file backup",
			"error", err,
			"source", backupPath,
			"destination", zipPath,
		)
		return fmt.Errorf("gagal mengkompresi file backup: %v", err)
	}

	// Get file info untuk ukuran file
	fileInfo, err := os.Stat(zipPath)
	if err != nil {
		logger.Warnw("Gagal mendapatkan informasi file backup",
			"error", err,
			"path", zipPath,
		)
	} else {
		logger.Infow("Backup selesai",
			"database", config.Database,
			"file", zipPath,
			"size_bytes", fileInfo.Size(),
			"duration", time.Since(time.Now()),
		)
	}

	return nil
}
