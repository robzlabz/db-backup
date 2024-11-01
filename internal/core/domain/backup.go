package domain

import "time"

type Backup struct {
	ID        int64
	DBType    string // mysql atau postgres
	DBName    string
	FilePath  string
	Size      int64
	CreatedAt time.Time
}

type BackupConfig struct {
	Host       string
	Port       int
	User       string
	Password   string
	Database   string
	OutputPath string
}
