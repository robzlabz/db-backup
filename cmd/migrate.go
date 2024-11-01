package cmd

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migrations",
	RunE: func(cmd *cobra.Command, args []string) error {
		db, err := sqlx.Connect("sqlite3", "./backup.db")
		if err != nil {
			return err
		}
		defer db.Close()

		return initDB(db)
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
}

func initDB(db *sqlx.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS backup_configs (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		type TEXT NOT NULL,
		host TEXT NOT NULL,
		port INTEGER NOT NULL,
		database TEXT NOT NULL,
		username TEXT NOT NULL,
		password TEXT NOT NULL,
		interval INTEGER NOT NULL,
		output_path TEXT NOT NULL,
		last_backup INTEGER DEFAULT 0
	);`

	_, err := db.Exec(schema)
	return err
}
