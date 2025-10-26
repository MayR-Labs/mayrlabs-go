package commands

import (
	"bufio"
	"fmt"
	"sort"
	"strings"

	"github.com/MayR-Labs/mayrlabs-go/internal/utils"
	"github.com/spf13/cobra"
)

// EnvCmd is the parent command for environment operations
var EnvCmd = &cobra.Command{
	Use:   "env",
	Short: "Environment file management commands",
	Long:  "Commands for managing .env files, validating variables, and syncing examples",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Interactive mode - prompt user to choose subcommand
		choice, err := utils.PromptSelect(
			"What would you like to do?",
			[]string{"update-example", "validate", "arrange"},
		)
		if err != nil {
			return err
		}

		switch choice {
		case "update-example":
			// Prompt for source file
			source, err := utils.PromptInput("Enter source file (default: .env): ")
			if err != nil {
				return err
			}
			if source == "" {
				source = ".env"
			}
			return EnvUpdateExampleCmd.RunE(cmd, []string{source})
		case "validate":
			return EnvValidateCmd.RunE(cmd, []string{})
		case "arrange":
			// Prompt for file
			file, err := utils.PromptInput("Enter file to arrange (default: .env): ")
			if err != nil {
				return err
			}
			if file == "" {
				file = ".env"
			}
			return EnvArrangeCmd.RunE(cmd, []string{file})
		default:
			return fmt.Errorf("invalid choice")
		}
	},
}

// EnvUpdateExampleCmd syncs .env.example with .env or .env.staging
var EnvUpdateExampleCmd = &cobra.Command{
	Use:   "update-example [source]",
	Short: "Sync .env.example with .env or .env.staging",
	Long:  "Update .env.example file with keys from .env or .env.staging, removing sensitive values",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		source := ".env"
		if len(args) > 0 {
			source = args[0]
		}

		if !utils.FileExists(source) {
			return fmt.Errorf("source file %s does not exist", source)
		}

		// Read source file
		content, err := utils.ReadFile(source)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", source, err)
		}

		// Parse and sanitize
		lines := strings.Split(content, "\n")
		var exampleLines []string

		for _, line := range lines {
			line = strings.TrimSpace(line)

			// Keep comments and empty lines
			if line == "" || strings.HasPrefix(line, "#") {
				exampleLines = append(exampleLines, line)
				continue
			}

			// Parse key=value
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				key := parts[0]
				// Remove the value, keep only the key
				exampleLines = append(exampleLines, fmt.Sprintf("%s=", key))
			}
		}

		// Write to .env.example
		exampleContent := strings.Join(exampleLines, "\n")
		if err := utils.WriteFile(".env.example", exampleContent); err != nil {
			return fmt.Errorf("failed to write .env.example: %w", err)
		}

		fmt.Println("✅ .env.example updated successfully!")
		return nil
	},
}

// EnvValidateCmd checks for missing keys, invalid values, or duplicates
var EnvValidateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Check for missing keys, invalid values, or duplicated variables",
	Long:  "Validate .env file for common issues like duplicates, empty values, and syntax errors",
	RunE: func(cmd *cobra.Command, args []string) error {
		if !utils.FileExists(".env") {
			return fmt.Errorf(".env file does not exist")
		}

		content, err := utils.ReadFile(".env")
		if err != nil {
			return fmt.Errorf("failed to read .env: %w", err)
		}

		lines := strings.Split(content, "\n")
		seen := make(map[string]int)
		issues := 0

		for i, line := range lines {
			lineNum := i + 1
			line = strings.TrimSpace(line)

			// Skip empty lines and comments
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}

			// Check for valid key=value format
			if !strings.Contains(line, "=") {
				fmt.Printf("❌ Line %d: Invalid format (missing '='): %s\n", lineNum, line)
				issues++
				continue
			}

			parts := strings.SplitN(line, "=", 2)
			key := parts[0]
			value := ""
			if len(parts) > 1 {
				value = parts[1]
			}

			// Check for duplicate keys
			if prevLine, exists := seen[key]; exists {
				fmt.Printf(
					"❌ Line %d: Duplicate key '%s' (first seen on line %d)\n",
					lineNum,
					key,
					prevLine,
				)
				issues++
			}
			seen[key] = lineNum

			// Check for empty values (warning only)
			if value == "" {
				fmt.Printf("⚠️  Line %d: Empty value for key '%s'\n", lineNum, key)
			}
		}

		if issues == 0 {
			fmt.Println("✅ .env file is valid!")
		} else {
			fmt.Printf("\n❌ Found %d issue(s) in .env file\n", issues)
		}

		return nil
	},
}

// EnvArrangeCmd sorts and groups environment keys
var EnvArrangeCmd = &cobra.Command{
	Use:   "arrange [file]",
	Short: "Sort and group environment keys by prefix",
	Long:  "Organize .env file by grouping keys with common prefixes (e.g., APP_*, DB_*)",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		file := ".env"
		if len(args) > 0 {
			file = args[0]
		}

		if !utils.FileExists(file) {
			return fmt.Errorf("file %s does not exist", file)
		}

		content, err := utils.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", file, err)
		}

		// Parse file
		type EnvLine struct {
			Original string
			Key      string
			Value    string
			Prefix   string
		}

		var entries []EnvLine
		var comments []string

		scanner := bufio.NewScanner(strings.NewReader(content))
		for scanner.Scan() {
			line := scanner.Text()
			trimmed := strings.TrimSpace(line)

			if trimmed == "" || strings.HasPrefix(trimmed, "#") {
				comments = append(comments, line)
				continue
			}

			parts := strings.SplitN(trimmed, "=", 2)
			if len(parts) == 2 {
				key := parts[0]
				value := parts[1]
				prefix := "OTHER"

				// Extract prefix (e.g., APP_ from APP_NAME)
				if idx := strings.Index(key, "_"); idx > 0 {
					prefix = key[:idx]
				}

				entries = append(entries, EnvLine{
					Original: line,
					Key:      key,
					Value:    value,
					Prefix:   prefix,
				})
			}
		}

		// Sort by prefix, then by key
		sort.Slice(entries, func(i, j int) bool {
			if entries[i].Prefix != entries[j].Prefix {
				return entries[i].Prefix < entries[j].Prefix
			}
			return entries[i].Key < entries[j].Key
		})

		// Build output
		var output []string

		// Add comments at the top
		if len(comments) > 0 {
			output = append(output, comments...)
			output = append(output, "")
		}

		// Group by prefix
		currentPrefix := ""
		for _, entry := range entries {
			if entry.Prefix != currentPrefix {
				if currentPrefix != "" {
					output = append(output, "")
				}
				output = append(output, fmt.Sprintf("# %s Configuration", entry.Prefix))
				currentPrefix = entry.Prefix
			}
			output = append(output, fmt.Sprintf("%s=%s", entry.Key, entry.Value))
		}

		// Write back
		arranged := strings.Join(output, "\n")
		if err := utils.WriteFile(file, arranged); err != nil {
			return fmt.Errorf("failed to write %s: %w", file, err)
		}

		fmt.Printf("✅ %s arranged successfully!\n", file)
		return nil
	},
}

func init() {
	EnvCmd.AddCommand(EnvUpdateExampleCmd)
	EnvCmd.AddCommand(EnvValidateCmd)
	EnvCmd.AddCommand(EnvArrangeCmd)
}
