package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"psam/database"
	"psam/database/models"
	"psam/server"

	"github.com/spf13/cobra"
)

func generateAPIKey() string {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal("Unable to generate API key:", err)
	}
	return hex.EncodeToString(b)
}

func main() {
	var rootCmd = &cobra.Command{Use: "app"}
	var serveCmd = &cobra.Command{
		Use:   "serve",
		Short: "Start the API server",
		Run: func(cmd *cobra.Command, args []string) {
			database.GetDB()
			server.StartServer()
		},
	}

	var createKeyCmd = &cobra.Command{
		Use:   "create-key",
		Short: "Create a new API key",
		Run: func(cmd *cobra.Command, args []string) {
			key := generateAPIKey()
			db := database.GetDB()
			apiKey := models.APIKey{Key: key}
			if err := db.Create(&apiKey).Error; err != nil {
				log.Fatal("Failed to create API key:", err)
			}
			fmt.Println("API Key created:", key)
		},
	}

	rootCmd.AddCommand(serveCmd, createKeyCmd)
	rootCmd.Execute()
}
