package commands

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/MayR-Labs/mayrlabs-go/internal/utils"
)

// ChangelogCmd is the parent command for changelog operations
var ChangelogCmd = &cobra.Command{
	Use:   "changelog",
	Short: "Changelog management commands",
	Long:  "Commands for creating and managing CHANGELOG.md files",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
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

		if utils.FileExists("CHANGELOG.md") && !force {
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
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		version := args[0]
		summary := args[1]
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

		newEntry := fmt.Sprintf(`## [%s]%s - %s

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
			newContent = content[:idx] + "\n" + newEntry + content[idx:]
		} else {
			// If no Unreleased section, append at the end
			newContent = content + "\n" + newEntry
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
