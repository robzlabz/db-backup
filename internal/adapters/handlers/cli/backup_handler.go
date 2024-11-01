package cli

import (
	"github.com/robzlabz/db-backup/internal/core/domain"
	"github.com/robzlabz/db-backup/internal/core/ports"
	"github.com/spf13/cobra"
)

type BackupHandler struct {
	backupService ports.BackupService
}

func NewBackupHandler(service ports.BackupService) *BackupHandler {
	return &BackupHandler{backupService: service}
}

func (h *BackupHandler) CreateBackupCommand() *cobra.Command {
	var config domain.BackupConfig

	cmd := &cobra.Command{
		Use:   "backup [mysql|postgres]",
		Short: "Backup database",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			dbType := args[0]
			return h.backupService.CreateBackup(dbType, config)
		},
	}

	cmd.Flags().StringVar(&config.Host, "host", "localhost", "Database host")
	cmd.Flags().IntVar(&config.Port, "port", 3306, "Database port")
	cmd.Flags().StringVar(&config.User, "user", "", "Database user")
	cmd.Flags().StringVar(&config.Password, "password", "", "Database password")
	cmd.Flags().StringVar(&config.Database, "database", "", "Database name")
	cmd.Flags().StringVar(&config.OutputPath, "output", "", "Output file path")

	return cmd
}
