package commands

import (
	"encoding/base64"
	"fmt"
	"os"

	"github.com/MayR-Labs/mayrlabs-go/internal/utils"
	"github.com/spf13/cobra"
)

// Base64Cmd is the parent command for base64 operations
var Base64Cmd = &cobra.Command{
	Use:   "base64 [encode|decode] [string]",
	Short: "Base64 encode or decode a string",
	Long:  "Encode or decode a string using base64 encoding",
	Args:  cobra.MaximumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		var operation, input string
		var err error

		copyFlag, _ := cmd.Flags().GetBool("copy")

		// Interactive mode
		if len(args) == 0 {
			operation, err = utils.PromptSelect(
				"Select operation:",
				[]string{"encode", "decode"},
			)
			if err != nil {
				return err
			}

			input, err = utils.SurveyInput("Enter string:", "")
			if err != nil {
				return err
			}
			if input == "" {
				return fmt.Errorf("input string cannot be empty")
			}

			// Ask for copy flag in interactive mode
			if !cmd.Flags().Changed("copy") {
				copyFlag, err = utils.SurveyConfirm("Copy result to clipboard?", false)
				if err != nil {
					return err
				}
			}
		} else if len(args) == 1 {
			// Only operation provided, prompt for string
			operation = args[0]
			input, err = utils.SurveyInput("Enter string:", "")
			if err != nil {
				return err
			}
			if input == "" {
				return fmt.Errorf("input string cannot be empty")
			}
		} else {
			operation = args[0]
			input = args[1]
		}

		var result string
		switch operation {
		case "encode":
			result = base64.StdEncoding.EncodeToString([]byte(input))
		case "decode":
			decoded, err := base64.StdEncoding.DecodeString(input)
			if err != nil {
				return fmt.Errorf("failed to decode base64: %w", err)
			}
			result = string(decoded)
		default:
			return fmt.Errorf("invalid operation: %s. Use 'encode' or 'decode'", operation)
		}

		fmt.Println(result)

		// Copy to clipboard if requested
		if copyFlag {
			if err := utils.CopyToClipboard(result); err != nil {
				fmt.Println("⚠️  Failed to copy to clipboard:", err)
			} else {
				fmt.Println("✅ Result copied to clipboard!")
			}
		}

		return nil
	},
}

// Base64FileCmd encodes a file to base64
var Base64FileCmd = &cobra.Command{
	Use:   "base64-file [path]",
	Short: "Encode a file to base64",
	Long:  "Read a file and encode its contents to base64",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var filePath string
		var err error

		copyFlag, _ := cmd.Flags().GetBool("copy")

		// Interactive mode
		if len(args) == 0 {
			filePath, err = utils.SurveyInput("Enter file path:", "")
			if err != nil {
				return err
			}
			if filePath == "" {
				return fmt.Errorf("file path cannot be empty")
			}

			// Ask for copy flag in interactive mode
			if !cmd.Flags().Changed("copy") {
				copyFlag, err = utils.SurveyConfirm("Copy result to clipboard?", false)
				if err != nil {
					return err
				}
			}
		} else {
			filePath = args[0]
		}

		if !utils.FileExists(filePath) {
			return fmt.Errorf("file does not exist: %s", filePath)
		}

		content, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read file: %w", err)
		}

		encoded := base64.StdEncoding.EncodeToString(content)
		fmt.Println(encoded)

		// Copy to clipboard if requested
		if copyFlag {
			if err := utils.CopyToClipboard(encoded); err != nil {
				fmt.Println("⚠️  Failed to copy to clipboard:", err)
			} else {
				fmt.Println("✅ Result copied to clipboard!")
			}
		}

		return nil
	},
}

// Base64DecodeToFileCmd decodes a base64 string and writes to a file
var Base64DecodeToFileCmd = &cobra.Command{
	Use:   "base64-decode-to-file [base64-string] [output-file]",
	Short: "Decode base64 string and write to a file",
	Long:  "Decode a base64 encoded string and save the result to a file",
	Args:  cobra.MaximumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		var base64String, outputFile string
		var err error

		// Interactive mode
		if len(args) == 0 {
			base64String, err = utils.SurveyInput("Enter base64 string:", "")
			if err != nil {
				return err
			}
			if base64String == "" {
				return fmt.Errorf("base64 string cannot be empty")
			}

			outputFile, err = utils.SurveyInput("Enter output file path:", "")
			if err != nil {
				return err
			}
			if outputFile == "" {
				return fmt.Errorf("output file path cannot be empty")
			}

			// Check if file exists and ask for confirmation
			if utils.FileExists(outputFile) {
				overwrite, err := utils.SurveyConfirm(fmt.Sprintf("File %s already exists. Overwrite?", outputFile), false)
				if err != nil {
					return err
				}
				if !overwrite {
					return fmt.Errorf("operation cancelled")
				}
			}
		} else if len(args) == 1 {
			base64String = args[0]
			outputFile, err = utils.SurveyInput("Enter output file path:", "")
			if err != nil {
				return err
			}
			if outputFile == "" {
				return fmt.Errorf("output file path cannot be empty")
			}
		} else {
			base64String = args[0]
			outputFile = args[1]
		}

		decoded, err := base64.StdEncoding.DecodeString(base64String)
		if err != nil {
			return fmt.Errorf("failed to decode base64: %w", err)
		}

		if err := os.WriteFile(outputFile, decoded, 0o644); err != nil {
			return fmt.Errorf("failed to write file: %w", err)
		}

		fmt.Printf("✅ Decoded content written to %s\n", outputFile)
		return nil
	},
}

func init() {
	Base64Cmd.Flags().BoolP("copy", "c", false, "Copy result to clipboard")
	Base64FileCmd.Flags().BoolP("copy", "c", false, "Copy result to clipboard")
}
