package cmd

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/robzlabz/db-backup/internal/adapters/backupers"
	"github.com/robzlabz/db-backup/internal/adapters/repositories"
	"github.com/robzlabz/db-backup/internal/core/ports"
	"github.com/robzlabz/db-backup/internal/core/services"
	"github.com/robzlabz/db-backup/pkg/logging"
	"github.com/spf13/cobra"
)

var scheduleCmd = &cobra.Command{
	Use:   "schedule",
	Short: "Menjalankan backup terjadwal",
	Long:  `Menjalankan backup database secara otomatis sesuai jadwal yang telah diatur`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := logging.Sugar()

		db, err := sqlx.Connect("sqlite3", "./backup.db")
		if err != nil {
			logger.Errorw("[CMD][Schedule][Run] Gagal membuka database",
				"error", err,
			)
		}
		defer db.Close()

		repo := repositories.NewSQLiteRepository(db)
		pgBackuper := backupers.NewPostgresBackuper()
		mysqlBackuper := backupers.NewMySQLBackuper()
		backupService := services.NewBackupService(repo, mysqlBackuper, pgBackuper)

		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()

		logger.Infow("[CMD][Schedule][Run] Memulai backup scheduler...")
		logger.Infow("[CMD][Schedule][Run] Memeriksa backup setiap 1 menit")

		// Jalankan pengecekan pertama kali
		checkAndExecuteBackups(backupService)

		// Kemudian jalankan setiap menit
		for range ticker.C {
			checkAndExecuteBackups(backupService)
		}
	},
}

func checkAndExecuteBackups(backupService ports.BackupService) {
	logger := logging.Sugar()

	configs, err := backupService.GetAllConfigs()
	if err != nil {
		logger.Errorw("[CMD][Schedule][checkAndExecuteBackups] Gagal mengambil konfigurasi",
			"error", err,
		)
		return
	}

	now := time.Now()
	for _, config := range configs {
		lastBackupTime := time.Unix(config.LastBackup, 0)
		nextBackupTime := lastBackupTime.Add(time.Duration(config.Interval) * time.Minute)

		if now.After(nextBackupTime) {
			logger.Infow("[CMD][Schedule][checkAndExecuteBackups] Menjalankan backup",
				"database", config.Database,
				"type", config.Type,
			)

			if err := backupService.ExecuteBackup(config); err != nil {
				logger.Errorw("[CMD][Schedule][checkAndExecuteBackups] Gagal backup database",
					"database", config.Database,
					"error", err,
				)
			} else {
				logger.Infow("[CMD][Schedule][checkAndExecuteBackups] Berhasil backup database",
					"database", config.Database,
					"next_backup", time.Now().Add(time.Duration(config.Interval)*time.Minute).Format("2006-01-02 15:04:05"),
				)
			}
		}
	}
}

func init() {
	rootCmd.AddCommand(scheduleCmd)
}
