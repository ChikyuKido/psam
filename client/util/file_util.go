package util

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func DoesFileExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

func ZipFile(dir string) (string, error) {
	tmpFile, err := os.CreateTemp(os.TempDir(), "archive-*.zip")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	defer tmpFile.Close()

	zipWriter := zip.NewWriter(tmpFile)
	defer zipWriter.Close()

	err = filepath.Walk(dir, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			relPath, err := filepath.Rel(dir, filePath)
			if err != nil {
				return err
			}

			zipFile, err := zipWriter.Create(relPath)
			if err != nil {
				return err
			}

			file, err := os.Open(filePath)
			if err != nil {
				return err
			}
			defer file.Close()

			_, err = io.Copy(zipFile, file)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return "", fmt.Errorf("failed to walk through directory: %w", err)
	}

	return tmpFile.Name(), nil
}
func ClearDir(dir string) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("failed to read directory %s: %w", dir, err)
	}

	for _, file := range files {
		filePath := filepath.Join(dir, file.Name())
		if file.IsDir() {
			if err := ClearDir(filePath); err != nil {
				return err
			}
			if err := os.Remove(filePath); err != nil {
				return fmt.Errorf("failed to remove directory %s: %w", filePath, err)
			}
		} else {
			if err := os.Remove(filePath); err != nil {
				return fmt.Errorf("failed to remove file %s: %w", filePath, err)
			}
		}
	}
	return nil
}
func UnzipFile(zipBytes []byte, destDir string) error {
	if err := ClearDir(destDir); err != nil {
		return fmt.Errorf("failed to clear destination directory: %w", err)
	}
	zipReader, err := zip.NewReader(bytes.NewReader(zipBytes), int64(len(zipBytes)))
	if err != nil {
		return fmt.Errorf("failed to read zip file: %w", err)
	}

	for _, file := range zipReader.File {
		destPath := filepath.Join(destDir, file.Name)

		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(destPath, os.ModePerm); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", destPath, err)
			}
			continue
		}

		dir := filepath.Dir(destPath)
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return fmt.Errorf("failed to create directory structure for file %s: %w", destPath, err)
		}
		destFile, err := os.Create(destPath)
		if err != nil {
			return fmt.Errorf("failed to create file %s: %w", destPath, err)
		}
		defer destFile.Close()

		zipFile, err := file.Open()
		if err != nil {
			return fmt.Errorf("failed to open file %s in zip archive: %w", file.Name, err)
		}
		defer zipFile.Close()

		_, err = io.Copy(destFile, zipFile)
		if err != nil {
			return fmt.Errorf("failed to copy file contents for %s: %w", file.Name, err)
		}
	}

	return nil
}
