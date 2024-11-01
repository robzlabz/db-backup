package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"

	"github.com/labstack/echo/v4"
	"github.com/robzlabz/db-backup/internal/adapters/backupers"
	"github.com/robzlabz/db-backup/internal/adapters/handlers/cli"
	"github.com/robzlabz/db-backup/internal/adapters/handlers/web"
	"github.com/robzlabz/db-backup/internal/adapters/repositories"
	"github.com/robzlabz/db-backup/internal/core/services"
)

func initDB() *sql.DB {
	db, err := sql.Open("sqlite3", "./backup.db")
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func main() {
	// Initialize dependencies
	db := initDB()
	repo := repositories.NewSQLiteRepository(db)
	mysqlBackuper := backupers.NewMySQLBackuper()
	pgBackuper := backupers.NewPostgresBackuper()

	// Initialize service
	backupService := services.NewBackupService(repo, mysqlBackuper, pgBackuper)

	// Initialize handlers
	webHandler := web.NewBackupHandler(backupService)
	cliHandler := cli.NewBackupHandler(backupService)

	// Setup Echo
	e := echo.New()
	e.GET("/backups", webHandler.GetBackups)
	e.POST("/backups/:type", webHandler.CreateBackup)

	// Setup CLI
	rootCmd := cmd.NewRootCommand()
	rootCmd.AddCommand(cliHandler.CreateBackupCommand())

	// Start server
	go e.Start(":8080")

	// Execute CLI
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
