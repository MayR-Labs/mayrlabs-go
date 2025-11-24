package commands

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/MayR-Labs/mayrlabs-go/internal/utils"
	"github.com/skip2/go-qrcode"
	"github.com/spf13/cobra"
)

// JWTDecodeCmd decodes a JWT token
var JWTDecodeCmd = &cobra.Command{
	Use:   "jwt-decode [token]",
	Short: "Decode a JWT token",
	Long:  "Decode and display the Header and Payload of a JSON Web Token (JWT).",
	RunE: func(cmd *cobra.Command, args []string) error {
		var token string
		var err error

		if len(args) > 0 {
			token = args[0]
		} else {
			token, err = utils.PromptInput("Enter JWT token: ")
			if err != nil {
				return err
			}
		}

		if token == "" {
			return fmt.Errorf("token cannot be empty")
		}

		parts := strings.Split(token, ".")
		if len(parts) != 3 {
			return fmt.Errorf("invalid JWT token format (expected 3 parts)")
		}

		printJSON := func(title, segment string) {
			// Add padding if needed
			if l := len(segment) % 4; l > 0 {
				segment += strings.Repeat("=", 4-l)
			}

			decoded, err := base64.URLEncoding.DecodeString(segment)
			if err != nil {
				// Try standard encoding just in case
				decoded, err = base64.StdEncoding.DecodeString(segment)
				if err != nil {
					fmt.Printf("❌ Failed to decode %s: %v\n", title, err)
					return
				}
			}

			var prettyJSON map[string]interface{}
			if err := json.Unmarshal(decoded, &prettyJSON); err != nil {
				fmt.Printf("❌ Failed to parse JSON for %s: %v\n", title, err)
				fmt.Printf("Raw content: %s\n", string(decoded))
				return
			}

			prettyBytes, _ := json.MarshalIndent(prettyJSON, "", "  ")
			fmt.Printf("\n🔹 %s:\n%s\n", title, string(prettyBytes))
		}

		printJSON("Header", parts[0])
		printJSON("Payload", parts[1])

		return nil
	},
}

// QRCmd generates a QR code
var QRCmd = &cobra.Command{
	Use:   "qr [text]",
	Short: "Generate a QR code",
	Long:  "Generate a QR code for the given text or URL, with an option to save to a file.",
	RunE: func(cmd *cobra.Command, args []string) error {
		var text string
		var err error

		if len(args) > 0 {
			text = args[0]
		} else {
			text, err = utils.PromptInput("Enter text/URL for QR code: ")
			if err != nil {
				return err
			}
		}

		if text == "" {
			return fmt.Errorf("text cannot be empty")
		}

		// Generate and print to console
		qr, err := qrcode.New(text, qrcode.Medium)
		if err != nil {
			return fmt.Errorf("failed to generate QR code: %w", err)
		}

		fmt.Println(qr.ToSmallString(false))

		// Prompt to save
		save, err := utils.SurveyConfirm("Do you want to save this QR code to a file?", false)
		if err != nil {
			return err
		}

		if save {
			filename, err := utils.PromptInput("Enter filename (e.g., qr.png): ")
			if err != nil {
				return err
			}
			if filename == "" {
				filename = "qr.png"
			}

			if err := qr.WriteFile(256, filename); err != nil {
				return fmt.Errorf("failed to save file: %w", err)
			}
			fmt.Printf("✅ QR code saved to %s\n", filename)
		}

		return nil
	},
}

// RegexCmd tests a regex pattern
var RegexCmd = &cobra.Command{
	Use:   "regex [pattern] [text]",
	Short: "Test a regex pattern",
	Long:  "Test a Regular Expression pattern against a text string.",
	RunE: func(cmd *cobra.Command, args []string) error {
		var pattern, text string
		var err error

		if len(args) >= 1 {
			pattern = args[0]
		} else {
			pattern, err = utils.PromptInput("Enter regex pattern: ")
			if err != nil {
				return err
			}
		}

		if len(args) >= 2 {
			text = args[1]
		} else {
			text, err = utils.PromptInput("Enter text to test: ")
			if err != nil {
				return err
			}
		}

		if pattern == "" {
			return fmt.Errorf("pattern cannot be empty")
		}

		re, err := regexp.Compile(pattern)
		if err != nil {
			return fmt.Errorf("invalid regex pattern: %w", err)
		}

		matches := re.FindAllString(text, -1)
		if len(matches) > 0 {
			fmt.Printf("✅ Found %d matches:\n", len(matches))
			for i, match := range matches {
				fmt.Printf("  %d: %s\n", i+1, match)
			}
		} else {
			fmt.Println("❌ No matches found.")
		}

		return nil
	},
}

// BaseCmd converts numbers between bases
var BaseCmd = &cobra.Command{
	Use:   "base [fromBase] [num] [toBase]",
	Short: "Convert a number between bases",
	Long:  "Convert a number from one base to another (e.g., base 10 100 16). Supported bases: 2-36.",
	RunE: func(cmd *cobra.Command, args []string) error {
		var fromBase, toBase int
		var numStr string
		var err error

		// Helper to parse base
		parseBase := func(s string) (int, error) {
			b, err := strconv.Atoi(s)
			if err != nil || b < 2 || b > 36 {
				return 0, fmt.Errorf("invalid base: %s (must be 2-36)", s)
			}
			return b, nil
		}

		// Interactive mode
		if len(args) == 0 {
			fromStr, err := utils.PromptInput("Enter source base (e.g., 10): ")
			if err != nil {
				return err
			}
			fromBase, err = parseBase(fromStr)
			if err != nil {
				return err
			}

			numStr, err = utils.PromptInput("Enter number: ")
			if err != nil {
				return err
			}

			toStr, err := utils.PromptInput("Enter target base (optional, press enter to see all): ")
			if err != nil {
				return err
			}

			if toStr != "" {
				toBase, err = parseBase(toStr)
				if err != nil {
					return err
				}
			}
		} else {
			// Args mode
			if len(args) >= 1 {
				fromBase, err = parseBase(args[0])
				if err != nil {
					return err
				}
			}
			if len(args) >= 2 {
				numStr = args[1]
			}
			if len(args) >= 3 {
				toBase, err = parseBase(args[2])
				if err != nil {
					return err
				}
			}
		}

		// Validate number
		val, err := strconv.ParseInt(numStr, fromBase, 64)
		if err != nil {
			return fmt.Errorf("invalid number '%s' for base %d", numStr, fromBase)
		}

		fmt.Printf("\nInput: %s (Base %d)\n", numStr, fromBase)
		fmt.Println("--------------------------------")

		if toBase > 0 {
			// Specific conversion
			res := strconv.FormatInt(val, toBase)
			fmt.Printf("Result: %s (Base %d)\n", res, toBase)
		} else {
			// Show common bases
			fmt.Printf("Binary (2):  %s\n", strconv.FormatInt(val, 2))
			fmt.Printf("Octal (8):   %s\n", strconv.FormatInt(val, 8))
			fmt.Printf("Decimal (10):%d\n", val)
			fmt.Printf("Hex (16):    %s\n", strconv.FormatInt(val, 16))
		}

		return nil
	},
}
