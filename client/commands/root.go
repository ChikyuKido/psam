package commands

import (
	"github.com/spf13/cobra"
	"psam_client/commands/game"
	"psam_client/commands/settings"
)

var rootCmd = &cobra.Command{
	Use:   "psam",
	Short: "Psam CLI - manage game saves",
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	game.Init()
	rootCmd.AddCommand(game.GameCmd)
	settings.Init()
	rootCmd.AddCommand(settings.SettingsCmd)
}
