package backupers

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/robzlabz/db-backup/internal/core/domain"
)

type mysqlBackuper struct{}

func NewMySQLBackuper() *mysqlBackuper {
	return &mysqlBackuper{}
}

func (b *mysqlBackuper) Backup(config domain.BackupConfig) error {
	host := config.Host
	if host == "localhost" {
		host = "127.0.0.1"
	}

	log.Printf("[Backuper][MySQLBackuper] Backing up database: %s@%s:%d/%s",
		config.User, host, config.Port, config.Database)

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

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("[Backuper][MySQLBackuper][Backup] Error: %v\nOutput: %s", err, string(output))
		return fmt.Errorf("backup error: %v", err)
	}

	log.Printf("[Backuper][MySQLBackuper][Backup] Backup completed: %s", config.OutputPath)

	now := time.Now().Format("20240101120000")
	filename := fmt.Sprintf("%s_%s.sql", config.Database, now)

	return os.WriteFile(filepath.Join(config.OutputPath, filename), output, 0644)
}
