package backupers

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/robzlabz/db-backup/internal/core/domain"
)

type mysqlBackuper struct{}

func NewMySQLBackuper() *mysqlBackuper {
	return &mysqlBackuper{}
}

func (b *mysqlBackuper) Backup(config domain.BackupConfig) error {
	cmd := exec.Command("mysqldump",
		"-h", config.Host,
		"-P", fmt.Sprintf("%d", config.Port),
		"-u", config.User,
		fmt.Sprintf("-p%s", config.Password),
		config.Database,
	)

	output, err := cmd.Output()
	if err != nil {
		return err
	}

	return os.WriteFile(config.OutputPath, output, 0644)
}
