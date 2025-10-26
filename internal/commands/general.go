package commands

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/MayR-Labs/mayrlabs-go/internal/utils"
	"github.com/spf13/cobra"
)

// DNSClearCmd clears the DNS cache
var DNSClearCmd = &cobra.Command{
	Use:   "dns-clear",
	Short: "Clear the DNS cache (choose macOS, Linux, or Windows)",
	Long:  "Clear the DNS cache for your operating system to resolve DNS issues",
	RunE: func(cmd *cobra.Command, args []string) error {
		os := runtime.GOOS

		var cmdStr string
		var cmdArgs []string

		switch os {
		case "darwin": // macOS
			fmt.Println("Clearing DNS cache on macOS...")
			cmdStr = "sudo"
			cmdArgs = []string{"dscacheutil", "-flushcache"}
		case "linux":
			fmt.Println("Clearing DNS cache on Linux...")
			// Try systemd-resolved first
			cmdStr = "sudo"
			cmdArgs = []string{"systemd-resolve", "--flush-caches"}
		case "windows":
			fmt.Println("Clearing DNS cache on Windows...")
			cmdStr = "ipconfig"
			cmdArgs = []string{"/flushdns"}
		default:
			return fmt.Errorf("unsupported operating system: %s", os)
		}

		command := exec.Command(cmdStr, cmdArgs...)
		output, err := command.CombinedOutput()
		if err != nil {
			// On Linux, try alternative command
			if os == "linux" {
				cmdArgs = []string{"service", "nscd", "restart"}
				command = exec.Command("sudo", cmdArgs...)
				output, err = command.CombinedOutput()
				if err != nil {
					return fmt.Errorf(
						"failed to clear DNS cache: %w\nOutput: %s",
						err,
						string(output),
					)
				}
			} else {
				return fmt.Errorf("failed to clear DNS cache: %w\nOutput: %s", err, string(output))
			}
		}

		fmt.Println("✅ DNS cache cleared successfully!")
		if len(output) > 0 {
			fmt.Println(string(output))
		}
		return nil
	},
}

// CreateKeystoreCmd creates a new keystore interactively
var CreateKeystoreCmd = &cobra.Command{
	Use:   "create-keystore",
	Short: "Create a new keystore interactively",
	Long:  "Create a new Java keystore file for Android app signing",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("🔐 Creating Android Keystore...")

		alias, err := utils.PromptInput("Enter key alias: ")
		if err != nil {
			return err
		}

		keystoreName, err := utils.PromptInput(
			"Enter keystore filename (e.g., my-release-key.jks): ",
		)
		if err != nil {
			return err
		}

		validity, err := utils.PromptInput("Enter validity in days (optional, press Enter for default 10000): ")
		if err != nil {
			return err
		}
		if validity == "" {
			validity = "10000"
		}

		// Create directory for keystore if path contains directory
		if dir := strings.TrimSuffix(keystoreName, filepath.Base(keystoreName)); dir != "" {
			if err := os.MkdirAll(dir, 0o755); err != nil {
				return fmt.Errorf("failed to create directory: %w", err)
			}
		}

		fmt.Println("\nGenerating keystore...")

		// Build keytool command
		command := exec.Command("keytool",
			"-genkeypair",
			"-v",
			"-storetype", "JKS",
			"-keyalg", "RSA",
			"-keysize", "2048",
			"-validity", validity,
			"-alias", alias,
			"-keystore", keystoreName,
		)

		command.Stdout = cmd.OutOrStdout()
		command.Stderr = cmd.ErrOrStderr()
		command.Stdin = cmd.InOrStdin()

		if err := command.Run(); err != nil {
			return fmt.Errorf("failed to create keystore: %w", err)
		}

		fmt.Printf("\n✅ Keystore '%s' created successfully!\n", keystoreName)
		return nil
	},
}
