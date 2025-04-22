package server

import (
	"github.com/gin-gonic/gin"
	"psam/server/middleware"
)

func StartServer() {
	r := gin.Default()

	r.Use(middleware.AuthMiddleware())

	v1 := r.Group("/api/v1/game")
	{
		v1.GET("/listGames", listGamesHandler)
		v1.GET("/getGameDetails/:gameName", getGameDetailsHandler)
		v1.POST("/uploadSave/:gameName/:version", uploadSaveHandler)
		v1.GET("/getSave/:gameName/:version", getSaveHandler)
		v1.DELETE("/deleteSave", deleteSaveHandler)
	}

	r.Run(":8080")
}
