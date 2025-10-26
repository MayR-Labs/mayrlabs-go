package commands

import (
	"fmt"
	"os/exec"
	"runtime"

	"github.com/spf13/cobra"
)

const (
	osLinux   = "linux"
	osDarwin  = "darwin"
	osWindows = "windows"
)

// openURL opens a URL in the default browser
func openURL(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case osDarwin:
		cmd = "open"
		args = []string{url}
	case osLinux:
		cmd = "xdg-open"
		args = []string{url}
	case osWindows:
		cmd = "cmd"
		args = []string{"/c", "start", url}
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	command := exec.Command(cmd, args...)
	return command.Start()
}

// VisitCmd opens the MayR Labs website
var VisitCmd = &cobra.Command{
	Use:   "visit",
	Short: "Visit the MayR Labs website",
	Long:  "Open the MayR Labs website (https://mayrlabs.com) in your default browser",
	RunE: func(cmd *cobra.Command, args []string) error {
		url := "https://mayrlabs.com"
		fmt.Printf("Opening %s in your browser...\n", url)
		if err := openURL(url); err != nil {
			return fmt.Errorf("failed to open URL: %w", err)
		}
		fmt.Println("✅ Website opened successfully!")
		return nil
	},
}

// GithubCmd opens the GitHub repository
var GithubCmd = &cobra.Command{
	Use:   "github",
	Short: "Visit the mayrlabs-go GitHub repository",
	Long:  "Open the mayrlabs-go GitHub repository (https://github.com/MayR-Labs/mayrlabs-go) in your default browser",
	RunE: func(cmd *cobra.Command, args []string) error {
		url := "https://github.com/MayR-Labs/mayrlabs-go"
		fmt.Printf("Opening %s in your browser...\n", url)
		if err := openURL(url); err != nil {
			return fmt.Errorf("failed to open URL: %w", err)
		}
		fmt.Println("✅ Repository opened successfully!")
		return nil
	},
}
