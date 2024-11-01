package main

import (
	"github.com/robzlabz/db-backup/cmd"
	"github.com/robzlabz/db-backup/pkg/logging"
)

func main() {
	// Initialize logger
	logging.InitLogger()
	defer logging.Logger.Sync()

	cmd.Execute()
}
