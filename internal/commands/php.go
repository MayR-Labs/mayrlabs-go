package commands

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

// PHPCmd is the parent command for PHP operations
var PHPCmd = &cobra.Command{
	Use:   "php",
	Short: "PHP-related commands",
	Long:  "Commands for managing PHP projects (linting, formatting, etc.)",
	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.Help()
		if err != nil {
			return
		}
	},
}

// PHPCSFixCmd runs PHP CodeSniffer/CS-Fixer
var PHPCSFixCmd = &cobra.Command{
	Use:   "cs-fix",
	Short: "Run PHP CodeSniffer/CS-Fixer",
	Long:  "Fix PHP code style issues using PHP-CS-Fixer",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Running PHP-CS-Fixer...")

		command := exec.Command("php-cs-fixer", "fix", ".")
		command.Stdout = cmd.OutOrStdout()
		command.Stderr = cmd.ErrOrStderr()

		if err := command.Run(); err != nil {
			return fmt.Errorf("failed to run PHP-CS-Fixer: %w", err)
		}

		fmt.Println("✅ PHP code style fixed!")
		return nil
	},
}

// PHPLintCmd lints PHP files
var PHPLintCmd = &cobra.Command{
	Use:   "lint",
	Short: "Lint PHP files",
	Long:  "Check PHP files for syntax errors",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Linting PHP files...")

		// Use find and php -l to check syntax
		command := exec.Command("sh", "-c", "find . -name '*.php' -exec php -l {} \\;")
		output, err := command.CombinedOutput()

		fmt.Println(string(output))

		if err != nil {
			return fmt.Errorf("PHP linting failed: %w", err)
		}

		fmt.Println("✅ All PHP files are valid!")
		return nil
	},
}

func init() {
	PHPCmd.AddCommand(PHPCSFixCmd)
	PHPCmd.AddCommand(PHPLintCmd)
}
