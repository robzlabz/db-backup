package cmd

import (
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"

	"github.com/jmoiron/sqlx"
	"github.com/manifoldco/promptui"
	"github.com/robzlabz/db-backup/internal/adapters/backupers"
	"github.com/robzlabz/db-backup/internal/adapters/repositories"
	"github.com/robzlabz/db-backup/internal/core/services"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Menghapus konfigurasi backup database",
	Long:  `Menghapus konfigurasi backup database dari daftar backup terjadwal`,
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

		// Mengambil dan menampilkan daftar konfigurasi
		configs, err := backupService.GetAllConfigs()
		if err != nil {
			fmt.Printf("Gagal mengambil daftar database: %v\n", err)
			return
		}

		if len(configs) == 0 {
			fmt.Println("Tidak ada konfigurasi backup yang tersedia")
			return
		}

		// Menampilkan tabel konfigurasi
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', tabwriter.TabIndent)
		fmt.Fprintln(w, "ID\tTIPE\tDATABASE\tHOST\tPORT\tUSERNAME\tINTERVAL (MENIT)")
		fmt.Fprintln(w, "--\t----\t--------\t----\t----\t--------\t---------------")

		for _, config := range configs {
			fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%d\t%s\t%d\n",
				config.ID,
				config.Type,
				config.Database,
				config.Host,
				config.Port,
				config.User,
				config.Interval,
			)
		}
		w.Flush()

		// Prompt untuk memilih ID yang akan dihapus
		idPrompt := promptui.Prompt{
			Label: "Masukkan ID konfigurasi yang akan dihapus",
			Validate: func(input string) error {
				_, err := strconv.Atoi(input)
				if err != nil {
					return fmt.Errorf("ID harus berupa angka")
				}
				return nil
			},
		}

		idStr, err := idPrompt.Run()
		if err != nil {
			fmt.Printf("Dibatalkan: %v\n", err)
			return
		}

		id, _ := strconv.Atoi(idStr)

		// Konfirmasi penghapusan
		confirmPrompt := promptui.Prompt{
			Label:     "Apakah Anda yakin ingin menghapus konfigurasi ini",
			IsConfirm: true,
		}

		result, err := confirmPrompt.Run()
		if err != nil || result != "y" {
			fmt.Println("Penghapusan dibatalkan")
			return
		}

		// Melakukan penghapusan
		if err := repo.Delete(id); err != nil {
			fmt.Printf("Gagal menghapus konfigurasi: %v\n", err)
			return
		}

		fmt.Printf("Berhasil menghapus konfigurasi dengan ID %d\n", id)
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
