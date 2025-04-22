package game

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"psam_client/database/services"
	"psam_client/util"
)

var GameCmd = &cobra.Command{
	Use:   "game",
	Short: "Manage game saves",
}
var addGame = &cobra.Command{
	Use:   "add [game] [path]",
	Short: "Adds a game to the local db",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		game := args[0]
		path := args[1]
		dbGame, _ := services.GetGame(game)
		if dbGame != nil {
			fmt.Println("Game already exists. Overwrite the old dir path")
			err := services.DeleteGame(game)
			if err != nil {
				fmt.Println("Failed to delete old game data")
				return
			}
		}
		err := services.AddGame(game, path)
		if err != nil {
			fmt.Println("Error adding game:", err)
			return
		}
		fmt.Println("Successfully added game: " + game)
	},
}
var listGame = &cobra.Command{
	Use:   "list [location]",
	Short: "List the games on the client or server",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		location := args[0]
		if location == "client" {
			games, err := services.ListGames()
			if err != nil {
				fmt.Println("Error listing games:", err)
				return
			}
			fmt.Println("Games on the client:")
			for _, game := range games {
				fmt.Printf("> %s: %s\n", game.GameName, game.DirPath)
			}
		} else if location == "server" {
			request, err := util.MyClient.DoRequest("GET", "http://localhost:8080/api/v1/game/listGames", nil, "")
			if err != nil {
				fmt.Println("Error getting game list:", err)
				return
			}
			var data map[string]interface{}
			err = json.Unmarshal(request, &data)
			if err != nil {
				fmt.Println("Error unmarshalling game list:", err)
				return
			}
			fmt.Println("Games on the server:")
			for _, game := range data["games"].([]interface{}) {
				fmt.Println("> " + game.(string))
			}
		} else {
			fmt.Println("Location not supported. Use server or client")
		}
	},
}
var getGameDetails = &cobra.Command{
	Use:   "details [game]",
	Short: "Get game details for a game on the server",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		game := args[0]
		request, err := util.MyClient.DoRequest("GET", "http://localhost:8080/api/v1/game/getGameDetails/"+game, nil, "")
		if err != nil {
			fmt.Println("Error getting game list:", err)
			return
		}
		fmt.Println(string(request))
	},
}
var putGame = &cobra.Command{
	Use:   "put [game] [version]",
	Short: "Put a game to the server",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		gameName := args[0]
		version := args[1]
		game, err := services.GetGame(gameName)
		if err != nil {
			fmt.Println("Error getting game data:", err)
			return
		}
		if !util.DoesFileExist(game.DirPath) {
			fmt.Println("Game Path does not exist. Please update it")
			return
		}
		file, err := util.ZipFile(game.DirPath)
		if err != nil {
			fmt.Println("Error zipping file:", err)
			return
		}
		_, err = util.MyClient.UploadFile(fmt.Sprintf("http://localhost:8080/api/v1/game/uploadSave/%s/%s", gameName, version), "save", file, nil)
		if err != nil {
			fmt.Println("Error uploading file:", err)
			return
		}
		fmt.Println("Successfully uploaded file")
	},
}
var getGame = &cobra.Command{
	Use:   "get [game] [version]",
	Short: "Gets a game from the server",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		gameName := args[0]
		version := args[1]
		game, err := services.GetGame(gameName)
		if err != nil {
			fmt.Println("Error getting game data:", err)
			return
		}
		if !util.DoesFileExist(game.DirPath) {
			fmt.Println("Game Path does not exist. Please update it")
			return
		}
		file, err := util.MyClient.DoRequest("GET", fmt.Sprintf("http://localhost:8080/api/v1/game/getSave/%s/%s", gameName, version), nil, "")

		err = util.UnzipFile(file, game.DirPath)
		if err != nil {
			fmt.Println("Error unzipping file:", err)
			return
		}
		fmt.Println("Successfully downloaded file")
	},
}

func Init() {
	GameCmd.AddCommand(addGame)
	GameCmd.AddCommand(getGame)
	GameCmd.AddCommand(putGame)
	GameCmd.AddCommand(listGame)
	GameCmd.AddCommand(getGameDetails)
}
