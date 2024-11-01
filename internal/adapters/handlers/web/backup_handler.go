package web

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/robzlabz/db-backup/internal/core/domain"
	"github.com/robzlabz/db-backup/internal/core/ports"
)

type BackupHandler struct {
	backupService ports.BackupService
}

func NewBackupHandler(service ports.BackupService) *BackupHandler {
	return &BackupHandler{backupService: service}
}

func (h *BackupHandler) GetBackups(c echo.Context) error {
	backups, err := h.backupService.GetAllBackups()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, backups)
}

func (h *BackupHandler) CreateBackup(c echo.Context) error {
	var config domain.BackupConfig
	if err := c.Bind(&config); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	dbType := c.Param("type")
	if err := h.backupService.CreateBackup(dbType, config); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Backup created successfully",
	})
}
