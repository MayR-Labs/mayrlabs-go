package commands

import (
	"fmt"
	"time"

	"github.com/MayR-Labs/mayrlabs-go/internal/utils"
	"github.com/spf13/cobra"
)

// ChangelogCmd is the parent command for changelog operations
var ChangelogCmd = &cobra.Command{
	Use:   "changelog",
	Short: "Changelog management commands",
	Long:  "Commands for creating and managing CHANGELOG.md files",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Interactive mode - prompt user to choose subcommand
		choice, err := utils.PromptSelect(
			"What would you like to do?",
			[]string{"create", "record"},
		)
		if err != nil {
			return err
		}

		switch choice {
		case "create":
			return ChangelogCreateCmd.RunE(cmd, []string{})
		case "record":
			// Prompt for version
			version, err := utils.PromptInput("Enter version (e.g., 1.0.0): ")
			if err != nil {
				return err
			}

			// Prompt for summary
			summary, err := utils.PromptInput("Enter summary: ")
			if err != nil {
				return err
			}

			return ChangelogRecordCmd.RunE(cmd, []string{version, summary})
		default:
			return fmt.Errorf("invalid choice")
		}
	},
}

// ChangelogCreateCmd creates a new CHANGELOG.md
var ChangelogCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create or overwrite CHANGELOG.md",
	Long:  "Initialize a new CHANGELOG.md file following Keep a Changelog format",
	RunE: func(cmd *cobra.Command, args []string) error {
		force, err := cmd.Flags().GetBool("force")
		if err != nil {
			return err
		}

		// Interactive mode: ask for force if file exists and flag not set
		if utils.FileExists("CHANGELOG.md") && !force && !cmd.Flags().Changed("force") {
			force, err = utils.SurveyConfirm("CHANGELOG.md already exists. Overwrite?", false)
			if err != nil {
				return err
			}
			if !force {
				return fmt.Errorf("operation cancelled")
			}
		} else if utils.FileExists("CHANGELOG.md") && !force {
			return fmt.Errorf("CHANGELOG.md already exists. Use --force to overwrite")
		}

		content := `# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Initial project setup

### Changed

### Deprecated

### Removed

### Fixed

### Security
`

		if err := utils.WriteFile("CHANGELOG.md", content); err != nil {
			return fmt.Errorf("failed to write CHANGELOG.md: %w", err)
		}

		fmt.Println("✅ CHANGELOG.md created successfully!")
		return nil
	},
}

func init() {
	ChangelogCreateCmd.Flags().BoolP("force", "f", false, "Overwrite existing CHANGELOG.md")
}

// ChangelogRecordCmd adds a new entry to CHANGELOG.md
var ChangelogRecordCmd = &cobra.Command{
	Use:   "record [version] [summary]",
	Short: "Add a new entry to CHANGELOG.md",
	Long:  "Record a new version entry in CHANGELOG.md with optional --wip flag",
	Args:  cobra.MaximumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		var version, summary string
		var err error

		// Interactive mode if no args provided
		if len(args) < 2 {
			version, err = utils.SurveyInput("Enter version (e.g., 1.0.0):", "")
			if err != nil {
				return err
			}

			summary, err = utils.SurveyInput("Enter summary:", "")
			if err != nil {
				return err
			}

			// Ask for additional feedback
			feedback, err := utils.SurveyInput("Any additional feedback/notes (optional):", "")
			if err != nil {
				return err
			}
			if feedback != "" {
				summary = summary + " - " + feedback
			}

			// Ask for wip flag in interactive mode if not already set
			if !cmd.Flags().Changed("wip") {
				wip, err := utils.SurveyConfirm("Mark as Work In Progress (WIP)?", false)
				if err != nil {
					return err
				}
				err = cmd.Flags().Set("wip", fmt.Sprintf("%t", wip))
				if err != nil {
					return err
				}
			}
		} else {
			version = args[0]
			summary = args[1]
		}

		wip, err := cmd.Flags().GetBool("wip")
		if err != nil {
			return err
		}

		if !utils.FileExists("CHANGELOG.md") {
			return fmt.Errorf("CHANGELOG.md does not exist. Run 'mayrlabs changelog create' first")
		}

		content, err := utils.ReadFile("CHANGELOG.md")
		if err != nil {
			return fmt.Errorf("failed to read CHANGELOG.md: %w", err)
		}

		// Build new entry
		date := time.Now().Format("2006-01-02")
		status := ""
		if wip {
			status = " [WIP]"
		}

		newEntry := fmt.Sprintf(`
--------------------
## [%s]%s - %s

### Added
- %s

`, version, status, date, summary)

		// Insert after "## [Unreleased]" section
		insertPoint := "## [Unreleased]"
		idx := 0
		for i := 0; i < len(content); i++ {
			if i+len(insertPoint) <= len(content) && content[i:i+len(insertPoint)] == insertPoint {
				// Find the end of the Unreleased section
				endIdx := i + len(insertPoint)
				for endIdx < len(content) && content[endIdx] != '#' {
					endIdx++
				}
				// Insert before the next section
				if endIdx < len(content) && content[endIdx-1] == '\n' {
					endIdx--
				}
				idx = endIdx
				break
			}
		}

		var newContent string
		if idx > 0 {
			newContent = content[:idx] + newEntry + content[idx:]
		} else {
			// If no Unreleased section, append at the end
			newContent = content + newEntry
		}

		if err := utils.WriteFile("CHANGELOG.md", newContent); err != nil {
			return fmt.Errorf("failed to update CHANGELOG.md: %w", err)
		}

		fmt.Printf("✅ Added version %s to CHANGELOG.md!\n", version)
		return nil
	},
}

func init() {
	ChangelogRecordCmd.Flags().BoolP("wip", "w", false, "Mark version as Work In Progress")

	ChangelogCmd.AddCommand(ChangelogCreateCmd)
	ChangelogCmd.AddCommand(ChangelogRecordCmd)
}
