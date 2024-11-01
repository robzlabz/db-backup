package utils

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"

	"github.com/robzlabz/db-backup/pkg/logging"
)

// CompressFile mengkompresi file ke format zip
func CompressFile(source, destination string) error {
	logger := logging.Sugar()
	logger.Debugw("Memulai kompresi file",
		"source", source,
		"destination", destination,
	)

	// Buat file zip
	zipfile, err := os.Create(destination)
	if err != nil {
		logger.Errorw("Gagal membuat file zip",
			"error", err,
			"destination", destination,
		)
		return err
	}
	defer zipfile.Close()

	// Buat writer zip
	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	// Buka file sumber
	file, err := os.Open(source)
	if err != nil {
		logger.Errorw("Gagal membuka file sumber",
			"error", err,
			"source", source,
		)
		return err
	}
	defer file.Close()

	// Buat file di dalam zip
	writer, err := archive.Create(filepath.Base(source))
	if err != nil {
		logger.Errorw("Gagal membuat entry dalam zip",
			"error", err,
			"filename", filepath.Base(source),
		)
		return err
	}

	// Salin isi file ke zip
	if _, err := io.Copy(writer, file); err != nil {
		logger.Errorw("Gagal menyalin file ke zip",
			"error", err,
			"source", source,
		)
		return err
	}

	logger.Infow("Berhasil mengkompresi file",
		"source", source,
		"destination", destination,
	)

	// Hapus file asli setelah kompresi
	if err := os.Remove(source); err != nil {
		logger.Warnw("Gagal menghapus file asli",
			"error", err,
			"source", source,
		)
	}

	return nil
}
