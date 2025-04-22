package main

import (
	"log"
	"psam_client/commands"
	"psam_client/database"
)

func main() {
	err := database.InitDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	commands.Execute()
}
