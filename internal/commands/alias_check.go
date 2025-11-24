package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/MayR-Labs/mayrlabs-go/internal/utils"
	"github.com/spf13/cobra"
)

// AliasListCmd lists all aliases
var AliasListCmd = &cobra.Command{
	Use:   "alias-list",
	Short: "List all mayrlabs aliases",
	Long:  "Display all mayrlabs aliases from your shell configuration file",
	RunE: func(cmd *cobra.Command, args []string) error {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get home directory: %w", err)
		}

		// Get executable path
		execPath, err := os.Executable()
		if err != nil {
			return fmt.Errorf("failed to get executable path: %w", err)
		}

		// Determine shell config files to check
		shell := os.Getenv("SHELL")
		var configFiles []string

		switch {
		case strings.Contains(shell, "zsh"):
			configFiles = []string{filepath.Join(homeDir, ".zshrc")}
		case strings.Contains(shell, "bash"):
			if runtime.GOOS == "darwin" {
				configFiles = []string{
					filepath.Join(homeDir, ".bash_profile"),
					filepath.Join(homeDir, ".bashrc"),
				}
			} else {
				configFiles = []string{filepath.Join(homeDir, ".bashrc")}
			}
		case strings.Contains(shell, "fish"):
			configFiles = []string{filepath.Join(homeDir, ".config", "fish", "config.fish")}
		default:
			configFiles = []string{filepath.Join(homeDir, ".bashrc")}
		}

		fmt.Println("🔍 Searching for mayrlabs aliases...")
		fmt.Println()

		found := false
		for _, configFile := range configFiles {
			if !utils.FileExists(configFile) {
				continue
			}

			content, err := utils.ReadFile(configFile)
			if err != nil {
				continue
			}

			lines := strings.Split(content, "\n")
			var aliases []string

			for _, line := range lines {
				line = strings.TrimSpace(line)
				// Check if line contains alias and mayrlabs
				if strings.HasPrefix(line, "alias ") && (strings.Contains(line, execPath) || strings.Contains(line, "mayrlabs")) {
					aliases = append(aliases, line)
				}
			}

			if len(aliases) > 0 {
				found = true
				fmt.Printf("📄 %s:\n", configFile)
				for _, alias := range aliases {
					fmt.Printf("  %s\n", alias)
				}
				fmt.Println()
			}
		}

		if !found {
			fmt.Println("No mayrlabs aliases found")
			fmt.Println("Use 'mayrlabs alias' or 'mayrlabs ai-alias' to create one")
		}

		return nil
	},
}

// AliasClearCmd clears aliases with confirmation
var AliasClearCmd = &cobra.Command{
	Use:   "alias-clear",
	Short: "Clear all mayrlabs aliases with confirmation",
	Long:  "Remove all mayrlabs aliases from your shell configuration file after PIN confirmation",
	RunE: func(cmd *cobra.Command, args []string) error {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get home directory: %w", err)
		}

		// Get executable path
		execPath, err := os.Executable()
		if err != nil {
			return fmt.Errorf("failed to get executable path: %w", err)
		}

		// Determine shell config files to check
		shell := os.Getenv("SHELL")
		var configFiles []string

		switch {
		case strings.Contains(shell, "zsh"):
			configFiles = []string{filepath.Join(homeDir, ".zshrc")}
		case strings.Contains(shell, "bash"):
			if runtime.GOOS == "darwin" {
				configFiles = []string{
					filepath.Join(homeDir, ".bash_profile"),
					filepath.Join(homeDir, ".bashrc"),
				}
			} else {
				configFiles = []string{filepath.Join(homeDir, ".bashrc")}
			}
		case strings.Contains(shell, "fish"):
			configFiles = []string{filepath.Join(homeDir, ".config", "fish", "config.fish")}
		default:
			configFiles = []string{filepath.Join(homeDir, ".bashrc")}
		}

		// Confirm with PIN
		if err := utils.ConfirmWithPIN("⚠️  This will remove all mayrlabs aliases from your shell configuration"); err != nil {
			return fmt.Errorf("incorrect PIN. Operation cancelled")
		}

		// Remove aliases from config files
		removed := false
		for _, configFile := range configFiles {
			if !utils.FileExists(configFile) {
				continue
			}

			content, err := utils.ReadFile(configFile)
			if err != nil {
				continue
			}

			lines := strings.Split(content, "\n")
			var newLines []string
			var removedFromFile int

			for _, line := range lines {
				trimmed := strings.TrimSpace(line)
				// Skip lines that contain alias and mayrlabs
				if strings.HasPrefix(trimmed, "alias ") && (strings.Contains(trimmed, execPath) || strings.Contains(trimmed, "mayrlabs")) {
					removedFromFile++
					continue
				}
				newLines = append(newLines, line)
			}

			if removedFromFile > 0 {
				// Write back the file
				newContent := strings.Join(newLines, "\n")
				if err := os.WriteFile(configFile, []byte(newContent), 0o644); err != nil {
					return fmt.Errorf("failed to update %s: %w", configFile, err)
				}
				fmt.Printf("✅ Removed %d alias(es) from %s\n", removedFromFile, configFile)
				removed = true
			}
		}

		if !removed {
			fmt.Println("No mayrlabs aliases found to remove")
		} else {
			fmt.Println("\n⚠️  Please restart your terminal or run 'source <config-file>' to apply changes")
		}

		return nil
	},
}
