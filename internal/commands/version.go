package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version will be set during build time via ldflags
var Version = "dev"

// VersionCmd displays the version of the CLI
var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display the version of mayrlabs CLI",
	Long:  "Display the current version of the mayrlabs CLI tool",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("mayrlabs version %s\n", Version)
		return nil
	},
}
