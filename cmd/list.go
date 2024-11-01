package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/jmoiron/sqlx"
	"github.com/robzlabz/db-backup/internal/adapters/backupers"
	"github.com/robzlabz/db-backup/internal/adapters/repositories"
	"github.com/robzlabz/db-backup/internal/core/services"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Menampilkan daftar database yang terjadwal",
	Long:  `Menampilkan daftar database yang akan di-backup secara otomatis beserta intervalnya`,
	Run: func(cmd *cobra.Command, args []string) {
		db, err := sqlx.Connect("sqlite3", "./backup.db")
		if err != nil {
			fmt.Printf("Gagal membuka database: %v\n", err)
			return
		}
		defer db.Close()

		repo := repositories.NewSQLiteRepository(db)
		pgBackuper := backupers.NewPostgresBackuper()
		mysqlBackuper := backupers.NewMySQLBackuper()
		backupService := services.NewBackupService(repo, mysqlBackuper, pgBackuper)
		configs, err := backupService.GetAllConfigs()
		if err != nil {
			fmt.Printf("Gagal mengambil daftar database: %v\n", err)
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.TabIndent)
		fmt.Fprintln(w, "TIPE\tDATABASE\tHOST\tPORT\tUSERNAME\tINTERVAL (MENIT)")
		fmt.Fprintln(w, "----\t--------\t----\t----\t--------\t---------------")

		for _, config := range configs {
			fmt.Fprintf(w, "%s\t%s\t%s\t%d\t%s\t%d\n",
				config.Type,
				config.Database,
				config.Host,
				config.Port,
				config.User,
				config.Interval,
			)
		}
		w.Flush()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
