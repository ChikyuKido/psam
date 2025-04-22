package game

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const saveRoot = "./saves"

type GameDetails struct {
	GameName string
	Versions map[string][]string
}

func ensureDir(path string) error {
	return os.MkdirAll(path, os.ModePerm)
}

func ListGames() ([]string, error) {
	var games = make([]string, 0)

	entries, err := os.ReadDir(saveRoot)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			games = append(games, entry.Name())
		}
	}
	return games, nil
}

func GetGameDetails(gameName string) (GameDetails, error) {
	gamePath := filepath.Join(saveRoot, gameName)
	details := GameDetails{
		GameName: gameName,
		Versions: make(map[string][]string),
	}

	versions, err := os.ReadDir(gamePath)
	if err != nil {
		return details, err
	}

	for _, version := range versions {
		if version.IsDir() {
			versionPath := filepath.Join(gamePath, version.Name())
			files, _ := os.ReadDir(versionPath)

			for _, f := range files {
				if !f.IsDir() && strings.HasSuffix(f.Name(), ".zip") {
					details.Versions[version.Name()] = append(details.Versions[version.Name()], strings.TrimSuffix(f.Name(), ".zip"))
				}
			}
		}
	}

	return details, nil
}

func AddSave(gameName, version, zipPath string) error {
	timestamp := time.Now().Format("2006-01-02T15-04-05")
	destDir := filepath.Join(saveRoot, gameName, version)

	if err := ensureDir(destDir); err != nil {
		return err
	}

	destPath := filepath.Join(destDir, fmt.Sprintf("%s.zip", timestamp))
	srcFile, err := os.Open(zipPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	return err
}

func GetSave(gameName, version, timestamp string) (string, error) {
	dirPath := filepath.Join(saveRoot, gameName, version)

	if timestamp == "0" {
		// Get latest
		files, err := os.ReadDir(dirPath)
		if err != nil {
			return "", err
		}

		var latest string
		for _, f := range files {
			if !f.IsDir() && strings.HasSuffix(f.Name(), ".zip") {
				latest = f.Name()
			}
		}

		if latest == "" {
			return "", errors.New("no saves found")
		}
		return filepath.Join(dirPath, latest), nil
	}

	filePath := filepath.Join(dirPath, fmt.Sprintf("%s.zip", timestamp))
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return "", errors.New("save file not found")
	}

	return filePath, nil
}

func DeleteSave(gameName, version, timestamp string) error {
	if gameName == "" {
		return errors.New("game name is required")
	}

	// Delete the entire game
	if version == "" {
		return os.RemoveAll(filepath.Join(saveRoot, gameName))
	}

	// Delete the entire version
	if timestamp == "" {
		return os.RemoveAll(filepath.Join(saveRoot, gameName, version))
	}

	// Delete specific save file
	filePath := filepath.Join(saveRoot, gameName, version, fmt.Sprintf("%s.zip", timestamp))
	return os.Remove(filePath)
}
