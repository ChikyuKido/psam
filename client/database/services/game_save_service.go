package services

import (
	"psam_client/database"
	"psam_client/database/models"
)

func AddGame(gameName, dirPath string) error {
	game := models.GameSave{GameName: gameName, DirPath: dirPath}
	return database.DB.Create(&game).Error
}

func GetGame(gameName string) (*models.GameSave, error) {
	var game models.GameSave
	err := database.DB.Where("game_name = ?", gameName).First(&game).Error
	if err != nil {
		return nil, err
	}
	return &game, nil
}

func DeleteGame(gameName string) error {
	return database.DB.Where("game_name = ?", gameName).Delete(&models.GameSave{}).Error
}

func ListGames() ([]models.GameSave, error) {
	var games []models.GameSave
	err := database.DB.Find(&games).Error
	return games, err
}
