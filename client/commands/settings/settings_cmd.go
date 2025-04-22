package settings

import (
	"fmt"
	"github.com/spf13/cobra"
	"psam_client/database/services"
)

var SettingsCmd = &cobra.Command{
	Use:   "settings",
	Short: "Manage API settings",
}

var setApiKeyCmd = &cobra.Command{
	Use:   "key [api_key]",
	Short: "Set API key",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := services.SetAPIKey(args[0])
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println("API key set successfully!")
		}
	},
}
var setURL = &cobra.Command{
	Use:   "url [url]",
	Short: "Set server url",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		err := services.SetURL(args[0])
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println("URL set successfully!")
		}
	},
}

func Init() {
	SettingsCmd.AddCommand(setApiKeyCmd)
	SettingsCmd.AddCommand(setURL)
}
