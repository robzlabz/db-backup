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
	ID         int    `json:"id"`
	Type       string `json:"type"`
	Host       string `json:"host"`
	Port       int    `json:"port"`
	Database   string `json:"database"`
	User       string `json:"user"`
	Password   string `json:"password"`
	Interval   int    `json:"interval"`
	OutputPath string `json:"output_path" db:"output_path"`
	LastBackup int64  `json:"last_backup" db:"last_backup"`
}
