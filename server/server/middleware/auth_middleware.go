package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"psam/database"
	"psam/database/models"
	"strings"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid authorization header"})
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		var apiKey models.APIKey
		db := database.GetDB()
		if err := db.Where("key = ?", token).First(&apiKey).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid API Key"})
			return
		}

		c.Next()
	}
}
