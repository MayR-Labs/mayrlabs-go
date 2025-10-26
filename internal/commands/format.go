package commands

import (
	"fmt"
	"os/exec"

	"github.com/MayR-Labs/mayrlabs-go/internal/utils"
	"github.com/spf13/cobra"
)

// EditorConfigCmd generates .editorconfig file
var EditorConfigCmd = &cobra.Command{
	Use:   "editor-config [language]",
	Short: "Generate .editorconfig for a specific language",
	Long:  "Create an .editorconfig file with best practices for the specified language",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		force, err := cmd.Flags().GetBool("force")
		if err != nil {
			return err
		}

		if utils.FileExists(".editorconfig") && !force {
			return fmt.Errorf(".editorconfig already exists. Use --force to overwrite")
		}

		language := "general"
		if len(args) > 0 {
			language = args[0]
		}

		content := generateEditorConfig(language)
		if err := utils.WriteFile(".editorconfig", content); err != nil {
			return fmt.Errorf("failed to write .editorconfig: %w", err)
		}

		fmt.Printf("✅ .editorconfig created for %s!\n", language)
		return nil
	},
}

func init() {
	EditorConfigCmd.Flags().BoolP("force", "f", false, "Overwrite existing .editorconfig")
}

func generateEditorConfig(language string) string {
	base := `# EditorConfig is awesome: https://EditorConfig.org

# top-most EditorConfig file
root = true

# Unix-style newlines with a newline ending every file
[*]
end_of_line = lf
insert_final_newline = true
charset = utf-8
trim_trailing_whitespace = true
`

	switch language {
	case "go", "golang":
		return base + `
[*.go]
indent_style = tab
indent_size = 4

[go.mod]
indent_style = tab
indent_size = 4
`
	case "javascript", "js", "typescript", "ts":
		return base + `
[*.{js,jsx,ts,tsx}]
indent_style = space
indent_size = 2

[*.json]
indent_style = space
indent_size = 2
`
	case "python", "py":
		return base + `
[*.py]
indent_style = space
indent_size = 4
max_line_length = 120
`
	case "php":
		return base + `
[*.php]
indent_style = space
indent_size = 4
`
	case "dart", "flutter":
		return base + `
[*.dart]
indent_style = space
indent_size = 2
`
	default:
		return base + `
[*]
indent_style = space
indent_size = 2
`
	}
}

// FormatCmd formats project files
var FormatCmd = &cobra.Command{
	Use:   "format [language]",
	Short: "Format project files for a given language",
	Long:  "Auto-format your project files using the appropriate formatter for the language",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		language := ""
		if len(args) > 0 {
			language = args[0]
		} else {
			var err error
			fmt.Println("Available languages: go, javascript, python, php, dart")
			language, err = utils.PromptInput("Select language: ")
			if err != nil {
				return err
			}
		}

		return formatCode(language)
	},
}

func formatCode(language string) error {
	var cmd *exec.Cmd

	switch language {
	case "go", "golang":
		fmt.Println("Formatting Go code with gofmt...")
		cmd = exec.Command("gofmt", "-w", ".")

	case "javascript", "js":
		fmt.Println("Formatting JavaScript with prettier...")
		cmd = exec.Command("npx", "prettier", "--write", "**/*.js")

	case "typescript", "ts":
		fmt.Println("Formatting TypeScript with prettier...")
		cmd = exec.Command("npx", "prettier", "--write", "**/*.ts")

	case "python", "py":
		fmt.Println("Formatting Python with black...")
		cmd = exec.Command("black", ".")

	case "php":
		fmt.Println("Formatting PHP with php-cs-fixer...")
		cmd = exec.Command("php-cs-fixer", "fix", ".")

	case "dart", "flutter":
		fmt.Println("Formatting Dart with dart format...")
		cmd = exec.Command("dart", "format", ".")

	default:
		return fmt.Errorf("unsupported language: %s", language)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to format code: %w\nOutput: %s", err, string(output))
	}

	fmt.Println("✅ Code formatted successfully!")
	if len(output) > 0 {
		fmt.Println(string(output))
	}
	return nil
}
