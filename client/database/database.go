package database

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
	"psam_client/database/models"
)

var DB *gorm.DB

func InitDB() error {
	var err error
	dir := "~/.local/psam"
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory %w", err)
	}
	DB, err = gorm.Open(sqlite.Open("~/.local/psam/games.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return err
	}
	return DB.AutoMigrate(&models.GameSave{}, &models.Settings{})
}
