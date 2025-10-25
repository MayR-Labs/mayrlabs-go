package commands

import (
	"fmt"
	"os/exec"

	"github.com/MayR-Labs/mayrlabs-go/internal/utils"
	"github.com/spf13/cobra"
)

// JSCmd is the parent command for JavaScript operations
var JSCmd = &cobra.Command{
	Use:   "js",
	Short: "JavaScript-related commands",
	Long:  "Commands for managing JavaScript/Node.js projects",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// JSSetupPrettierCmd installs and configures Prettier
var JSSetupPrettierCmd = &cobra.Command{
	Use:   "setup-prettier",
	Short: "Install and configure Prettier with .prettierrc.yaml",
	Long:  "Set up Prettier code formatter with a default configuration file",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Installing Prettier...")

		// Install prettier
		installCmd := exec.Command("npm", "install", "--save-dev", "prettier")
		installCmd.Stdout = cmd.OutOrStdout()
		installCmd.Stderr = cmd.ErrOrStderr()

		if err := installCmd.Run(); err != nil {
			return fmt.Errorf("failed to install Prettier: %w", err)
		}

		// Create .prettierrc.yaml
		prettierConfig := `# Prettier configuration
printWidth: 100
tabWidth: 2
useTabs: false
semi: true
singleQuote: true
quoteProps: as-needed
jsxSingleQuote: false
trailingComma: es5
bracketSpacing: true
arrowParens: always
endOfLine: lf
`

		if err := utils.WriteFile(".prettierrc.yaml", prettierConfig); err != nil {
			return fmt.Errorf("failed to create .prettierrc.yaml: %w", err)
		}

		// Create .prettierignore
		prettierIgnore := `# Prettier ignore
node_modules/
dist/
build/
coverage/
*.min.js
*.min.css
package-lock.json
yarn.lock
`

		if err := utils.WriteFile(".prettierignore", prettierIgnore); err != nil {
			return fmt.Errorf("failed to create .prettierignore: %w", err)
		}

		fmt.Println("✅ Prettier configured successfully!")
		fmt.Println("   - .prettierrc.yaml created")
		fmt.Println("   - .prettierignore created")
		fmt.Println("\nRun 'mayrlabs js pretty' to format your code")
		return nil
	},
}

// JSPrettyCmd runs Prettier on the project
var JSPrettyCmd = &cobra.Command{
	Use:   "pretty",
	Short: "Run Prettier on the project",
	Long:  "Format JavaScript/TypeScript files using Prettier",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Running Prettier...")

		command := exec.Command("npx", "prettier", "--write", ".")
		command.Stdout = cmd.OutOrStdout()
		command.Stderr = cmd.ErrOrStderr()

		if err := command.Run(); err != nil {
			return fmt.Errorf("failed to run Prettier: %w", err)
		}

		fmt.Println("✅ Code formatted successfully!")
		return nil
	},
}

func init() {
	JSCmd.AddCommand(JSSetupPrettierCmd)
	JSCmd.AddCommand(JSPrettyCmd)
}
