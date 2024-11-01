package cmd

import (
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/robzlabz/db-backup/internal/adapters/backupers"
	"github.com/robzlabz/db-backup/internal/adapters/repositories"
	"github.com/robzlabz/db-backup/internal/core/ports"
	"github.com/robzlabz/db-backup/internal/core/services"
	"github.com/spf13/cobra"
)

var scheduleCmd = &cobra.Command{
	Use:   "schedule",
	Short: "Menjalankan backup terjadwal",
	Long:  `Menjalankan backup database secara otomatis sesuai jadwal yang telah diatur`,
	Run: func(cmd *cobra.Command, args []string) {
		db, err := sqlx.Connect("sqlite3", "./backup.db")
		if err != nil {
			log.Fatalf("Gagal membuka database: %v", err)
		}
		defer db.Close()

		repo := repositories.NewSQLiteRepository(db)
		pgBackuper := backupers.NewPostgresBackuper()
		mysqlBackuper := backupers.NewMySQLBackuper()
		backupService := services.NewBackupService(repo, mysqlBackuper, pgBackuper)

		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()

		log.Println("Memulai backup scheduler...")
		log.Println("Memeriksa backup setiap 1 menit")

		// Jalankan pengecekan pertama kali
		checkAndExecuteBackups(backupService)

		// Kemudian jalankan setiap menit
		for range ticker.C {
			checkAndExecuteBackups(backupService)
		}
	},
}

func checkAndExecuteBackups(backupService ports.BackupService) {
	configs, err := backupService.GetAllConfigs()
	if err != nil {
		log.Printf("Gagal mengambil konfigurasi: %v", err)
		return
	}

	now := time.Now()
	for _, config := range configs {
		lastBackupTime := time.Unix(config.LastBackup, 0)
		nextBackupTime := lastBackupTime.Add(time.Duration(config.Interval) * time.Minute)

		if now.After(nextBackupTime) {
			log.Printf("Menjalankan backup untuk database %s (%s)", config.Database, config.Type)
			if err := backupService.ExecuteBackup(config); err != nil {
				log.Printf("Gagal backup database %s: %v", config.Database, err)
			} else {
				log.Printf("Berhasil backup database %s", config.Database)
				log.Printf("Backup berikutnya pada: %s",
					time.Now().Add(time.Duration(config.Interval)*time.Minute).Format("2006-01-02 15:04:05"))
			}
		}
	}
}

func init() {
	rootCmd.AddCommand(scheduleCmd)
}
