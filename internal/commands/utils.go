package commands

import (
	"crypto/rand"
	"fmt"
	"math/big"
	mathrand "math/rand"
	"time"

	"github.com/MayR-Labs/mayrlabs-go/internal/utils"
	"github.com/google/uuid"
	"github.com/oklog/ulid/v2"
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
	Short: "Generate a random password (interactive mode if no args)",
	Long:  "Generate a secure random password with letters, numbers, and special characters",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var length int
		var passwordType string
		var err error

		copyFlag, _ := cmd.Flags().GetBool("copy")

		// Interactive mode
		if len(args) == 0 {
			lengthStr, err := utils.PromptInput("Enter password length (default: 16): ")
			if err != nil {
				return err
			}
			if lengthStr == "" {
				length = 16
			} else {
				_, err = fmt.Sscanf(lengthStr, "%d", &length)
				if err != nil {
					return fmt.Errorf("invalid length: %s", lengthStr)
				}
			}

			passwordType, err = utils.PromptSelect(
				"Select password type:",
				[]string{"alpha", "numeral", "alphanum", "alphanum+special", "sentence"},
			)
			if err != nil {
				return err
			}
		} else {
			length = 0
			_, err = fmt.Sscanf(args[0], "%d", &length)
			if err != nil {
				return fmt.Errorf("invalid length: %s", args[0])
			}
			passwordType = "alphanum+special" // Default type when length is provided
		}

		if length < 8 && passwordType != "sentence" {
			return fmt.Errorf("password length must be at least 8 characters")
		}

		var password string
		if passwordType == "alphanum+special" {
			password, err = utils.GeneratePassword(length)
		} else {
			password, err = utils.GeneratePasswordByType(length, passwordType)
		}
		if err != nil {
			return fmt.Errorf("failed to generate password: %w", err)
		}

		// Print password to stdout - this is the intended behavior
		// Users should copy and store this securely
		fmt.Println(password)

		// Copy to clipboard if requested
		if copyFlag {
			if err := utils.CopyToClipboard(password); err != nil {
				fmt.Println("⚠️  Failed to copy to clipboard:", err)
			} else {
				fmt.Println("✅ Password copied to clipboard!")
			}
		}

		return nil
	},
}

// HashCmd generates hash of a string
var HashCmd = &cobra.Command{
	Use:   "hash [string]",
	Short: "Generate hash of a string (interactive mode if no args)",
	Long:  "Generate a hash of the provided string using the specified algorithm",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		algorithm, err := cmd.Flags().GetString("algorithm")
		if err != nil {
			return err
		}
		if algorithm == "" {
			algorithm = "sha256"
		}

		copyFlag, _ := cmd.Flags().GetBool("copy")

		var input string
		// Interactive mode
		if len(args) == 0 {
			input, err = utils.PromptInput("Enter string to hash: ")
			if err != nil {
				return err
			}
			if input == "" {
				return fmt.Errorf("input string cannot be empty")
			}

			algorithmChoice, err := utils.PromptSelect(
				"Select hash algorithm:",
				[]string{"sha256", "sha1", "md5"},
			)
			if err != nil {
				return err
			}
			algorithm = algorithmChoice
		} else {
			input = args[0]
		}

		hash, err := utils.HashString(input, algorithm)
		if err != nil {
			return err
		}

		fmt.Printf("%s (%s): %s\n", input, algorithm, hash)

		// Copy to clipboard if requested
		if copyFlag {
			if err := utils.CopyToClipboard(hash); err != nil {
				fmt.Println("⚠️  Failed to copy to clipboard:", err)
			} else {
				fmt.Println("✅ Hash copied to clipboard!")
			}
		}

		return nil
	},
}

// HashFileCmd generates hash of a file
var HashFileCmd = &cobra.Command{
	Use:   "hash-file [file]",
	Short: "Generate hash of a file",
	Long:  "Generate a hash of the provided file using the specified algorithm",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		algorithm, err := cmd.Flags().GetString("algorithm")
		if err != nil {
			return err
		}
		if algorithm == "" {
			algorithm = "sha256"
		}

		copyFlag, _ := cmd.Flags().GetBool("copy")

		filePath := args[0]

		if !utils.FileExists(filePath) {
			return fmt.Errorf("file does not exist: %s", filePath)
		}

		hash, err := utils.HashFile(filePath, algorithm)
		if err != nil {
			return err
		}

		fmt.Printf("%s (%s): %s\n", filePath, algorithm, hash)

		// Copy to clipboard if requested
		if copyFlag {
			if err := utils.CopyToClipboard(hash); err != nil {
				fmt.Println("⚠️  Failed to copy to clipboard:", err)
			} else {
				fmt.Println("✅ Hash copied to clipboard!")
			}
		}

		return nil
	},
}

// ULIDCmd generates a ULID
var ULIDCmd = &cobra.Command{
	Use:   "ulid",
	Short: "Generate ULID (Universally Unique Lexicographically Sortable Identifier)",
	Long:  "Generate a random ULID with timestamp ordering",
	RunE: func(cmd *cobra.Command, args []string) error {
		entropy := ulid.Monotonic(rand.Reader, 0)
		id := ulid.MustNew(ulid.Timestamp(time.Now()), entropy)
		fmt.Println(id.String())
		return nil
	},
}

// RandomIntCmd generates a random integer
var RandomIntCmd = &cobra.Command{
	Use:   "random-int [start] [end]",
	Short: "Generate a random integer",
	Long:  "Generate a random integer. No args: can be negative. One arg: 0 to [end]. Two args: [start] to [end]",
	Args:  cobra.MaximumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		var start, end int64

		switch len(args) {
		case 0:
			// No params: use -100 to 100
			start, end = -100, 100
		case 1:
			// One param: use it as end, start from 0
			start = 0
			_, err := fmt.Sscanf(args[0], "%d", &end)
			if err != nil {
				return fmt.Errorf("invalid end value: %s", args[0])
			}
		case 2:
			// Two params: start and end
			_, err := fmt.Sscanf(args[0], "%d", &start)
			if err != nil {
				return fmt.Errorf("invalid start value: %s", args[0])
			}
			_, err = fmt.Sscanf(args[1], "%d", &end)
			if err != nil {
				return fmt.Errorf("invalid end value: %s", args[1])
			}
		}

		if start > end {
			return fmt.Errorf("start value must be less than or equal to end value")
		}

		// Use crypto/rand for secure random number generation
		rangeSize := end - start + 1
		n, err := rand.Int(rand.Reader, big.NewInt(rangeSize))
		if err != nil {
			// Fallback to math/rand if crypto/rand fails
			mathrand.Seed(time.Now().UnixNano())
			result := mathrand.Int63n(rangeSize) + start
			fmt.Println(result)
			return nil
		}

		result := n.Int64() + start
		fmt.Println(result)
		return nil
	},
}

func init() {
	HashCmd.Flags().StringP("algorithm", "a", "sha256", "Hash algorithm (md5, sha1, sha256)")
	HashCmd.Flags().BoolP("copy", "c", false, "Copy hash to clipboard")

	HashFileCmd.Flags().StringP("algorithm", "a", "sha256", "Hash algorithm (md5, sha1, sha256)")
	HashFileCmd.Flags().BoolP("copy", "c", false, "Copy hash to clipboard")

	PasswordCmd.Flags().BoolP("copy", "c", false, "Copy password to clipboard")
}
