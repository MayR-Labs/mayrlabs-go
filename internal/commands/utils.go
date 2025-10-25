package commands

import (
	"fmt"

	"github.com/MayR-Labs/mayrlabs-go/internal/utils"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

// UUIDCmd generates a UUID v4
var UUIDCmd = &cobra.Command{
	Use:   "uuid",
	Short: "Generate UUID v4",
	Long:  "Generate a random UUID (Universally Unique Identifier) version 4",
	RunE: func(cmd *cobra.Command, args []string) error {
		id := uuid.New()
		fmt.Println(id.String())
		return nil
	},
}

// PasswordCmd generates a random password
// Note: The password is intentionally printed to stdout for the user.
// Users should handle the output securely and avoid logging it.
var PasswordCmd = &cobra.Command{
	Use:   "password [length]",
	Short: "Generate a random password of specified length (default: 16)",
	Long:  "Generate a secure random password with letters, numbers, and special characters",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		length := 16
		if len(args) > 0 {
			var err error
			_, err = fmt.Sscanf(args[0], "%d", &length)
			if err != nil {
				return fmt.Errorf("invalid length: %s", args[0])
			}
		}

		if length < 8 {
			return fmt.Errorf("password length must be at least 8 characters")
		}

		password, err := utils.GeneratePassword(length)
		if err != nil {
			return fmt.Errorf("failed to generate password: %w", err)
		}

		// Print password to stdout - this is the intended behavior
		// Users should copy and store this securely
		fmt.Println(password)
		return nil
	},
}

// HashCmd generates hash of a string
var HashCmd = &cobra.Command{
	Use:   "hash [string]",
	Short: "Generate hash of a string using md5, sha1, or sha256",
	Long:  "Generate a hash of the provided string using the specified algorithm",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		algorithm, _ := cmd.Flags().GetString("algorithm")
		if algorithm == "" {
			algorithm = "sha256"
		}

		input := args[0]
		hash, err := utils.HashString(input, algorithm)
		if err != nil {
			return err
		}

		fmt.Printf("%s (%s): %s\n", input, algorithm, hash)
		return nil
	},
}

func init() {
	HashCmd.Flags().StringP("algorithm", "a", "sha256", "Hash algorithm (md5, sha1, sha256)")
}
