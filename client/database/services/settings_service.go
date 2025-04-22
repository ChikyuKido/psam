package services

import (
	"errors"
	"gorm.io/gorm"
	"psam_client/database"
	"psam_client/database/models"
)

func SetAPIKey(key string) error {
	var settings models.Settings
	err := database.DB.First(&settings).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		settings.APIKey = key
		return database.DB.Create(&settings).Error
	}
	settings.APIKey = key
	return database.DB.Save(&settings).Error
}

func GetAPIKey() (string, error) {
	var settings models.Settings
	err := database.DB.First(&settings).Error
	if err != nil {
		return "", err
	}
	return settings.APIKey, nil
}
