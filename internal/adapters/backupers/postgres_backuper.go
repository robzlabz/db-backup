package backupers

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"github.com/robzlabz/db-backup/internal/core/domain"
)

type PostgresBackuper struct{}

func NewPostgresBackuper() *PostgresBackuper {
	return &PostgresBackuper{}
}

func (b *PostgresBackuper) Backup(config domain.BackupConfig) error {
	// Mengambil konfigurasi yang diperlukan
	host := config.Host
	port := config.Port
	dbname := config.Database
	username := config.User
	password := config.Password
	backupDir := config.OutputPath

	// Membuat nama file backup dengan timestamp
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	filename := fmt.Sprintf("%s_%s.sql", dbname, timestamp)
	backupPath := filepath.Join(backupDir, filename)

	// Memastikan direktori backup ada
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return fmt.Errorf("gagal membuat direktori backup: %v", err)
	}

	// Menyiapkan environment variable untuk password
	env := os.Environ()
	env = append(env, fmt.Sprintf("PGPASSWORD=%s", password))

	// Menyiapkan command pg_dump
	cmd := exec.Command("pg_dump",
		"-h", host,
		"-p", strconv.Itoa(port),
		"-U", username,
		"-d", dbname,
		"-F", "p", // Format plain text SQL
		"-f", backupPath,
	)

	cmd.Env = env

	// Menjalankan backup
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("gagal melakukan backup: %v, output: %s", err, string(output))
	}

	return nil
}
