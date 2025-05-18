package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "loklm",
	Short: "Manage agent crafter environment",
	Long: `
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⠀⢀⣀⣀⣀⣀⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⢀⣤⣶⣿⣿⣿⣆⠘⠿⠟⢻⣿⣿⡇⢐⣷⣦⣄⡀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⢸⣿⣿⣿⣧⡄⠙⣿⣷⣶⣶⡿⠿⠿⢃⣼⡟⠻⣿⣿⣶⡄⠀⠀⠀⠀
⠀⠀⢰⣷⣌⠙⠉⣿⣿⡟⢀⣿⣿⡟⢁⣤⣤⣶⣾⣿⡇⠸⢿⣿⠿⢃⣴⡄⠀⠀
⠀⠀⢸⣿⣿⣿⣿⠿⠋⣠⣾⣿⣿⠀⣾⣿⣿⣛⠛⢿⣿⣶⣤⣤⣴⣿⣿⣿⡆⠀
⠀⣴⣤⣄⣀⣠⣤⣴⣾⣿⣿⣿⣿⣆⠘⠿⣿⣿⣷⡄⢹⣿⣿⠿⠟⢿⣿⣿⣿⠀
⠀⢸⣿⣿⡿⠛⠛⣻⣿⣿⣿⣿⣿⣿⣷⣦⣼⣿⣿⠃⣸⣿⠃⢰⣶⣾⣿⣿⡟⠀
⠀⠀⢿⡏⢠⣾⣿⣿⡿⠋⣠⣄⡉⢻⣿⣿⡿⠟⠁⠀⠛⠛⠀⠘⠿⠿⠿⠋⠀⠀
⠀⠀⠀⠁⠘⢿⣿⣿⣷⣤⣿⣿⠗⠀⣉⣥⣴⣶⡶⠋⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⣤⣀⡉⠛⠛⠋⣉⣠⣴⠿⢿⣿⠿⠋⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠈⠻⢿⣿⣿⣿⣿⡿⠋⣠⠟⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣤⡾⠋⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢰⡿⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀

		LokLM CLI is a developer-focused command-line interface designed to streamline the
		setup, execution, and management of a fully local AI development environment.`,
}

var configPath string

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&configPath, "config", "", "Path to the config file")
}
