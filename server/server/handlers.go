package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path/filepath"
	"psam/game"
)

func listGamesHandler(c *gin.Context) {
	games, err := game.ListGames()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"games": games})
}

func getGameDetailsHandler(c *gin.Context) {
	gameName := c.Param("gameName")

	details, err := game.GetGameDetails(gameName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, details)
}

func uploadSaveHandler(c *gin.Context) {
	gameName := c.Param("gameName")
	version := c.Param("version")
	file, err := c.FormFile("save")
	if err != nil || gameName == "" || version == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing file, gameName or version"})
		return
	}

	tempPath := filepath.Join(os.TempDir(), file.Filename)
	if err := c.SaveUploadedFile(file, tempPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save uploaded file"})
		return
	}

	if err := game.AddSave(gameName, version, tempPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	_ = os.Remove(tempPath)

	c.JSON(http.StatusOK, gin.H{"status": "save uploaded successfully"})
}

func getSaveHandler(c *gin.Context) {
	gameName := c.Param("gameName")
	version := c.Param("version")
	timestamp := c.DefaultQuery("timestamp", "0")

	if gameName == "" || version == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "gameName and version are required"})
		return
	}

	path, err := game.GetSave(gameName, version, timestamp)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	fileBytes, err := os.ReadFile(path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read file"})
		return
	}

	c.Data(http.StatusOK, "application/octet-stream", fileBytes)
}

func deleteSaveHandler(c *gin.Context) {
	var body struct {
		GameName  string `json:"gameName"`
		Version   string `json:"version"`
		Timestamp string `json:"timestamp"`
	}

	if err := c.BindJSON(&body); err != nil || body.GameName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "gameName required"})
		return
	}

	if err := game.DeleteSave(body.GameName, body.Version, body.Timestamp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "delete successful"})
}
