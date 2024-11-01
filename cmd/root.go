package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "db-backup",
	Short: "Aplikasi untuk backup database",
	Long:  `Aplikasi CLI untuk melakukan backup database PostgreSQL dan MySQL secara otomatis`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
