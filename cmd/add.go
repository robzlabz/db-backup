package cmd

import (
	"fmt"
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/manifoldco/promptui"
	"github.com/robzlabz/db-backup/internal/adapters/backupers"
	"github.com/robzlabz/db-backup/internal/adapters/repositories"
	"github.com/robzlabz/db-backup/internal/core/domain"
	"github.com/robzlabz/db-backup/internal/core/services"
	"github.com/robzlabz/db-backup/pkg/logging"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Menambahkan database baru untuk backup",
	Long:  `Menambahkan konfigurasi database baru yang akan di-backup secara otomatis`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := logging.Sugar()
		// Pilihan tipe database
		dbTypePrompt := promptui.Select{
			Label: "Pilih tipe database",
			Items: []string{"PostgreSQL", "MySQL"},
		}
		_, dbType, err := dbTypePrompt.Run()
		if err != nil {
			fmt.Printf("Gagal memilih tipe database: %v\n", err)
			return
		}

		// Input host
		hostPrompt := promptui.Prompt{
			Label:   "Host",
			Default: "localhost",
		}
		host, err := hostPrompt.Run()
		if err != nil {
			fmt.Printf("Gagal input host: %v\n", err)
			return
		}

		// Input port
		portPrompt := promptui.Prompt{
			Label: "Port",
			Default: func() string {
				if dbType == "PostgreSQL" {
					return "5432"
				}
				return "3306"
			}(),
			Validate: func(input string) error {
				_, err := strconv.Atoi(input)
				return err
			},
		}
		portStr, err := portPrompt.Run()
		if err != nil {
			fmt.Printf("Gagal input port: %v\n", err)
			return
		}
		port, _ := strconv.Atoi(portStr)

		// Input database name
		dbPrompt := promptui.Prompt{
			Label: "Nama Database",
		}
		dbname, err := dbPrompt.Run()
		if err != nil {
			logger.Errorf("Gagal input nama database: %v", err)
			return
		}

		// Input username
		userPrompt := promptui.Prompt{
			Label: "Username",
		}
		username, err := userPrompt.Run()
		if err != nil {
			logger.Errorf("Gagal input username: %v", err)
			return
		}

		// Input password
		passPrompt := promptui.Prompt{
			Label: "Password",
			Mask:  '*',
		}
		password, err := passPrompt.Run()
		if err != nil {
			logger.Errorf("Gagal input password: %v", err)
			return
		}

		// Input interval backup (dalam menit)
		intervalPrompt := promptui.Prompt{
			Label:   "Interval Backup (menit)",
			Default: "60",
			Validate: func(input string) error {
				_, err := strconv.Atoi(input)
				return err
			},
		}
		intervalStr, err := intervalPrompt.Run()
		if err != nil {
			logger.Errorf("Gagal input interval: %v", err)
			return
		}
		interval, _ := strconv.Atoi(intervalStr)

		config := domain.BackupConfig{
			Type:       dbType,
			Host:       host,
			Port:       port,
			Database:   dbname,
			User:       username,
			Password:   password,
			Interval:   interval,
			OutputPath: "./backup",
		}

		db, err := sqlx.Connect("sqlite3", "./backup.db")
		if err != nil {
			logger.Errorf("Gagal membuka database: %v", err)
			return
		}
		defer db.Close()

		repo := repositories.NewSQLiteRepository(db)
		pgBackuper := backupers.NewPostgresBackuper()
		mysqlBackuper := backupers.NewMySQLBackuper()
		backupService := services.NewBackupService(repo, mysqlBackuper, pgBackuper)
		if err := backupService.AddConfig(config); err != nil {
			logger.Errorf("Gagal menyimpan konfigurasi: %v", err)
			return
		}

		logger.Infof("Berhasil menambahkan konfigurasi database!")
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
