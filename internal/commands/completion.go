package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/MayR-Labs/mayrlabs-go/internal/utils"
	"github.com/spf13/cobra"
)

// InstallCompletionCmd installs autocompletion for the specified shell
var InstallCompletionCmd = &cobra.Command{
	Use:   "install-completion [shell...]",
	Short: "Install autocompletion for your shell",
	Long:  "Install autocompletion for zsh, bash, fish, or powershell. You can select multiple shells.",
	RunE: func(cmd *cobra.Command, args []string) error {
		shells := args
		if len(shells) == 0 {
			selected, err := utils.PromptMultiSelect(
				"Select shells to install completion for:",
				[]string{"zsh", "bash", "fish", "powershell"},
			)
			if err != nil {
				return err
			}
			shells = selected
		}

		if len(shells) == 0 {
			return fmt.Errorf("no shells selected")
		}

		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get home directory: %w", err)
		}

		configDir := filepath.Join(homeDir, ".mayrlabs")
		if err := os.MkdirAll(configDir, 0755); err != nil {
			return fmt.Errorf("failed to create config directory: %w", err)
		}

		for _, shell := range shells {
			fmt.Printf("Installing completion for %s...\n", shell)

			var completionFile string
			var rcFile string
			var sourceCmd string

			switch shell {
			case "zsh":
				completionFile = filepath.Join(configDir, "completion.zsh")
				rcFile = filepath.Join(homeDir, ".zshrc")
				sourceCmd = fmt.Sprintf("source %s", completionFile)
				// Generate completion
				if err := cmd.Root().GenZshCompletionFile(completionFile); err != nil {
					fmt.Printf("❌ Failed to generate zsh completion: %v\n", err)
					continue
				}
			case "bash":
				completionFile = filepath.Join(configDir, "completion.bash")
				rcFile = filepath.Join(homeDir, ".bashrc")
				// On macOS, it might be .bash_profile
				if _, err := os.Stat(rcFile); os.IsNotExist(err) {
					profile := filepath.Join(homeDir, ".bash_profile")
					if _, err := os.Stat(profile); err == nil {
						rcFile = profile
					}
				}
				sourceCmd = fmt.Sprintf("source %s", completionFile)
				if err := cmd.Root().GenBashCompletionFile(completionFile); err != nil {
					fmt.Printf("❌ Failed to generate bash completion: %v\n", err)
					continue
				}
			case "fish":
				completionFile = filepath.Join(configDir, "completion.fish")
				configPath := filepath.Join(homeDir, ".config", "fish")
				rcFile = filepath.Join(configPath, "config.fish")
				sourceCmd = fmt.Sprintf("source %s", completionFile)

				// Ensure fish config dir exists
				_ = os.MkdirAll(configPath, 0755)

				if err := cmd.Root().GenFishCompletionFile(completionFile, true); err != nil {
					fmt.Printf("❌ Failed to generate fish completion: %v\n", err)
					continue
				}
			case "powershell":
				completionFile = filepath.Join(configDir, "completion.ps1")
				// PowerShell profile path is variable, but usually in Documents
				// This is a simplification
				documentsDir := filepath.Join(homeDir, "Documents")
				// Check for OneDrive
				oneDriveDocs := filepath.Join(homeDir, "OneDrive", "Documents")
				if _, err := os.Stat(oneDriveDocs); err == nil {
					documentsDir = oneDriveDocs
				}

				psDir := filepath.Join(documentsDir, "PowerShell")
				_ = os.MkdirAll(psDir, 0755)
				rcFile = filepath.Join(psDir, "Microsoft.PowerShell_profile.ps1")

				sourceCmd = fmt.Sprintf(". %s", completionFile)
				if err := cmd.Root().GenPowerShellCompletionFile(completionFile); err != nil {
					fmt.Printf("❌ Failed to generate powershell completion: %v\n", err)
					continue
				}
			default:
				fmt.Printf("⚠️  Unsupported shell: %s\n", shell)
				continue
			}

			// Add to RC file if not present
			if rcFile != "" {
				// Check if RC file exists, if not create it
				f, err := os.OpenFile(rcFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
				if err != nil {
					fmt.Printf("⚠️  Could not open RC file %s: %v\n", rcFile, err)
					fmt.Printf("   Please manually add this line to your config: %s\n", sourceCmd)
					continue
				}

				content, err := os.ReadFile(rcFile)
				if err == nil {
					if !strings.Contains(string(content), completionFile) {
						if _, err := fmt.Fprintf(f, "\n# mayrlabs completion\n%s\n", sourceCmd); err != nil {
							fmt.Printf("⚠️  Failed to write to RC file: %v\n", err)
						} else {
							fmt.Printf("✅ Added source command to %s\n", rcFile)
						}
					} else {
						fmt.Printf("ℹ️  Completion already configured in %s\n", rcFile)
					}
				}
				_ = f.Close()
			}

			fmt.Printf("✅ Completion installed for %s!\n", shell)
		}

		fmt.Println("\n🎉 Installation complete! Please restart your shell or source your config file to apply changes.")
		return nil
	},
}
