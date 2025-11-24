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

// AliasSelfCmd creates a permanent alias for mayrlabs command
var AliasSelfCmd = &cobra.Command{
	Use:   "alias-self [alias-name]",
	Short: "Create a permanent alias for mayrlabs command",
	Long:  "Create a shell alias for mayrlabs command in your shell configuration file",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var aliasName string
		var err error

		// Interactive mode
		if len(args) == 0 {
			aliasName, err = utils.SurveyInput("Enter alias name (e.g., ml, mr):", "")
			if err != nil {
				return err
			}
			if aliasName == "" {
				return fmt.Errorf("alias name cannot be empty")
			}
		} else {
			aliasName = args[0]
		}

		// Validate alias name
		if strings.Contains(aliasName, " ") || strings.Contains(aliasName, "/") {
			return fmt.Errorf("invalid alias name: %s", aliasName)
		}

		// Detect shell configuration file
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get home directory: %w", err)
		}

		// Determine which shell config file to use
		shell := os.Getenv("SHELL")
		var configFile string

		switch {
		case strings.Contains(shell, "zsh"):
			configFile = filepath.Join(homeDir, ".zshrc")
		case strings.Contains(shell, "bash"):
			if runtime.GOOS == "darwin" {
				// macOS uses .bash_profile
				configFile = filepath.Join(homeDir, ".bash_profile")
			} else {
				configFile = filepath.Join(homeDir, ".bashrc")
			}
		case strings.Contains(shell, "fish"):
			configFile = filepath.Join(homeDir, ".config", "fish", "config.fish")
		default:
			// Default to .bashrc
			configFile = filepath.Join(homeDir, ".bashrc")
		}

		// Ask for confirmation
		confirm, err := utils.SurveyConfirm(fmt.Sprintf("Add alias '%s' to %s?", aliasName, configFile), true)
		if err != nil {
			return err
		}
		if !confirm {
			return fmt.Errorf("operation cancelled")
		}

		// Get the current executable path
		execPath, err := os.Executable()
		if err != nil {
			return fmt.Errorf("failed to get executable path: %w", err)
		}

		// Prepare alias line based on shell
		var aliasLine string
		if strings.Contains(configFile, "config.fish") {
			aliasLine = fmt.Sprintf("\nalias %s='%s'\n", aliasName, execPath)
		} else {
			aliasLine = fmt.Sprintf("\nalias %s='%s'\n", aliasName, execPath)
		}

		// Append to config file
		file, err := os.OpenFile(configFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
		if err != nil {
			return fmt.Errorf("failed to open config file: %w", err)
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				return
			}
		}(file)

		if _, err := file.WriteString(aliasLine); err != nil {
			return fmt.Errorf("failed to write alias: %w", err)
		}

		fmt.Printf("✅ Alias '%s' added to %s\n", aliasName, configFile)
		fmt.Printf("Run 'source %s' or restart your terminal to use the alias\n", configFile)
		return nil
	},
}

// UpgradeCmd upgrades mayrlabs to the latest version
var UpgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade mayrlabs to the latest version",
	Long:  "Download and install the latest version of mayrlabs CLI",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("🚀 Checking for updates...")

		// Check if installed via go install
		execPath, err := os.Executable()
		if err != nil {
			return fmt.Errorf("failed to get executable path: %w", err)
		}

		// Check if in GOPATH/bin
		gopath := os.Getenv("GOPATH")
		if gopath == "" {
			gopath = filepath.Join(os.Getenv("HOME"), "go")
		}
		gopathBin := filepath.Join(gopath, "bin")

		var upgradeMethod string

		if strings.Contains(execPath, gopathBin) {
			upgradeMethod = "go install"
		} else {
			upgradeMethod = "curl script"
		}

		fmt.Printf("\nDetected installation method: %s\n", upgradeMethod)
		fmt.Println("\nTo upgrade mayrlabs, run:")

		if upgradeMethod == "go install" {
			fmt.Println("\n  go install github.com/MayR-Labs/mayrlabs-go@latest")
		} else {
			fmt.Println("\n  curl -sSL https://raw.githubusercontent.com/MayR-Labs/mayrlabs-go/main/install.sh | bash")
			fmt.Println("\n  Or download from: https://github.com/MayR-Labs/mayrlabs-go/releases/latest")
		}

		// Ask if user wants to proceed with automatic upgrade
		if upgradeMethod == "go install" {
			proceed, err := utils.SurveyConfirm("Do you want to upgrade now using 'go install'?", true)
			if err != nil {
				return err
			}

			if proceed {
				fmt.Println("\n📦 Upgrading mayrlabs...")
				// Note: We can't really upgrade ourselves while running, so we'll just show the command
				fmt.Println("\nPlease run this command in your terminal:")
				fmt.Println("  go install github.com/MayR-Labs/mayrlabs-go@latest")
			}
		}

		return nil
	},
}
