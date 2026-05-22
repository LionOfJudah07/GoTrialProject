package services

import (
	"io"
	"os"
	"path/filepath"
	"time"
)

type BackupService struct {
	dataDir string
}

func NewBackupService(dataDir string) *BackupService {
	return &BackupService{dataDir: dataDir}
}

func (s *BackupService) Export() (string, error) {
	src := filepath.Join(s.dataDir, "data.db")
	dst := src + ".backup." + time.Now().Format("20060102-150405")
	if err := copyFile(src, dst); err != nil {
		return "", err
	}
	return dst, nil
}

func (s *BackupService) Restore(backupPath string) error {
	dst := filepath.Join(s.dataDir, "data.db")
	return copyFile(backupPath, dst)
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}