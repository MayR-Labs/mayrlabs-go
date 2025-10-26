package commands

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/MayR-Labs/mayrlabs-go/internal/utils"
	"github.com/spf13/cobra"
)

// GitCmd is the parent command for git operations
var GitCmd = &cobra.Command{
	Use:   "git",
	Short: "Git-related commands",
	Long:  "Commands for managing Git repositories and branches",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Interactive mode - prompt user to choose subcommand
		choice, err := utils.PromptSelect(
			"What would you like to do?",
			[]string{"prune-stale"},
		)
		if err != nil {
			return err
		}

		switch choice {
		case "prune-stale":
			return GitPruneStaleCmd.RunE(cmd, []string{})
		default:
			return fmt.Errorf("invalid choice")
		}
	},
}

// GitPruneStaleCmd deletes local branches not found on remote
var GitPruneStaleCmd = &cobra.Command{
	Use:   "prune-stale",
	Short: "Delete all local branches not found on the remote",
	Long:  "Remove local Git branches that no longer exist on the remote repository",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Fetching remote branches...")

		// Fetch with prune
		fetchCmd := exec.Command("git", "fetch", "--prune")
		if output, err := fetchCmd.CombinedOutput(); err != nil {
			return fmt.Errorf("failed to fetch: %w\nOutput: %s", err, string(output))
		}

		// Get list of local branches
		localCmd := exec.Command("git", "branch")
		localOutput, err := localCmd.Output()
		if err != nil {
			return fmt.Errorf("failed to get local branches: %w", err)
		}

		// Get current branch
		currentCmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
		currentOutput, err := currentCmd.Output()
		if err != nil {
			return fmt.Errorf("failed to get current branch: %w", err)
		}
		currentBranch := strings.TrimSpace(string(currentOutput))

		// Parse local branches
		localBranches := strings.Split(string(localOutput), "\n")
		deletedCount := 0

		for _, branch := range localBranches {
			branch = strings.TrimSpace(branch)
			branch = strings.TrimPrefix(branch, "* ")

			if branch == "" || branch == currentBranch || branch == "main" || branch == "master" {
				continue
			}

			// Check if branch exists on remote
			remoteCmd := exec.Command(
				"git",
				"branch",
				"-r",
				"--list",
				fmt.Sprintf("origin/%s", branch),
			) // #nosec G204
			remoteOutput, err := remoteCmd.Output()
			if err != nil {
				continue
			}

			// If branch doesn't exist on remote, delete it
			if strings.TrimSpace(string(remoteOutput)) == "" {
				fmt.Printf("Deleting stale branch: %s\n", branch)
				deleteCmd := exec.Command("git", "branch", "-D", branch)
				if err := deleteCmd.Run(); err != nil {
					fmt.Printf("⚠️  Failed to delete %s: %v\n", branch, err)
				} else {
					deletedCount++
				}
			}
		}

		if deletedCount == 0 {
			fmt.Println("✅ No stale branches found!")
		} else {
			fmt.Printf("✅ Deleted %d stale branch(es)!\n", deletedCount)
		}

		return nil
	},
}

func init() {
	GitCmd.AddCommand(GitPruneStaleCmd)
}
