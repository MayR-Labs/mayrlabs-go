package cmd

import (
	"github.com/spf13/cobra"

	"github.com/MayR-Labs/mayrlabs-go/internal/commands"
)

var rootCmd = &cobra.Command{
	Use:   "mayrlabs",
	Short: "🧰 MayR Labs CLI - Streamline your development workflow",
	Long: `MayR Labs CLI is a lightweight, cross-platform command-line tool 
built with Go to streamline common development, configuration, and 
automation tasks across projects.

It provides developers with a unified interface for generating configs, 
formatting code, managing environments, handling CI/CD, and keeping 
project structure consistent — all in seconds.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			_ = cmd.Help()
		}
	},
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// General commands
	rootCmd.AddCommand(commands.CreateKeystoreCmd)
	rootCmd.AddCommand(commands.DNSClearCmd)
	rootCmd.AddCommand(commands.CICmd)
	rootCmd.AddCommand(commands.FormatCmd)
	rootCmd.AddCommand(commands.AddLicenseCmd)
	rootCmd.AddCommand(commands.EditorConfigCmd)
	rootCmd.AddCommand(commands.HashCmd)
	rootCmd.AddCommand(commands.UUIDCmd)
	rootCmd.AddCommand(commands.PasswordCmd)
	rootCmd.AddCommand(commands.QuoteCmd)

	// Group commands
	rootCmd.AddCommand(commands.GitCmd)
	rootCmd.AddCommand(commands.EnvCmd)
	rootCmd.AddCommand(commands.ChangelogCmd)
	rootCmd.AddCommand(commands.FlutterCmd)
	rootCmd.AddCommand(commands.PHPCmd)
	rootCmd.AddCommand(commands.JSCmd)
}
