package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/robzlabz/db-backup/internal/adapters/backupers"
	"github.com/robzlabz/db-backup/internal/adapters/repositories"
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

		fmt.Println("Memulai backup terjadwal...")
		for {
			configs, err := backupService.GetAllConfigs()
			if err != nil {
				log.Printf("Gagal mengambil konfigurasi: %v", err)
				continue
			}

			for _, config := range configs {
				if err := backupService.ExecuteBackup(config); err != nil {
					log.Printf("Gagal backup database %s: %v", config.Database, err)
				} else {
					log.Printf("Berhasil backup database %s", config.Database)
				}
			}

			time.Sleep(time.Minute * 5) // Cek setiap 5 menit
		}
	},
}

func init() {
	rootCmd.AddCommand(scheduleCmd)
}
