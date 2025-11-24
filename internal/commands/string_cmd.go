package commands

import (
	"fmt"
	"html"
	"net/url"
	"regexp"
	"strings"
	"unicode"

	"github.com/MayR-Labs/mayrlabs-go/internal/utils"
	"github.com/spf13/cobra"
)

// StringCmd is the parent command for string utilities
var StringCmd = &cobra.Command{
	Use:   "string",
	Short: "String manipulation utilities",
	Long:  "A collection of utilities for string manipulation (case conversion, escaping, etc.)",
}

var StringCaseCmd = &cobra.Command{
	Use:   "case [type] [text]",
	Short: "Convert string case (camel, snake, kebab, upper, lower, title)",
	RunE: func(cmd *cobra.Command, args []string) error {
		var caseType, text string
		var err error

		if len(args) >= 1 {
			caseType = args[0]
		} else {
			caseType, err = utils.PromptSelect("Select case type:", []string{"camel", "snake", "kebab", "upper", "lower", "title"})
			if err != nil {
				return err
			}
		}

		if len(args) >= 2 {
			text = args[1]
		} else {
			text, err = utils.PromptInput("Enter text: ")
			if err != nil {
				return err
			}
		}

		var result string
		switch caseType {
		case "upper":
			result = strings.ToUpper(text)
		case "lower":
			result = strings.ToLower(text)
		case "title":
			result = strings.Title(strings.ToLower(text))
		case "snake":
			result = toSnakeCase(text)
		case "kebab":
			result = strings.ReplaceAll(toSnakeCase(text), "_", "-")
		case "camel":
			result = toCamelCase(text)
		default:
			return fmt.Errorf("unknown case type: %s", caseType)
		}

		fmt.Println(result)
		return nil
	},
}

var StringEscapeCmd = &cobra.Command{
	Use:   "escape [type] [text]",
	Short: "Escape string (url, html)",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runEscape(args, true)
	},
}

var StringUnescapeCmd = &cobra.Command{
	Use:   "unescape [type] [text]",
	Short: "Unescape string (url, html)",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runEscape(args, false)
	},
}

func runEscape(args []string, escape bool) error {
	var escapeType, text string
	var err error

	if len(args) >= 1 {
		escapeType = args[0]
	} else {
		escapeType, err = utils.PromptSelect("Select type:", []string{"url", "html"})
		if err != nil {
			return err
		}
	}

	if len(args) >= 2 {
		text = args[1]
	} else {
		text, err = utils.PromptInput("Enter text: ")
		if err != nil {
			return err
		}
	}

	var result string
	switch escapeType {
	case "url":
		if escape {
			result = url.QueryEscape(text)
		} else {
			result, err = url.QueryUnescape(text)
			if err != nil {
				return err
			}
		}
	case "html":
		if escape {
			result = html.EscapeString(text)
		} else {
			result = html.UnescapeString(text)
		}
	default:
		return fmt.Errorf("unknown type: %s", escapeType)
	}

	fmt.Println(result)
	return nil
}

var StringReverseCmd = &cobra.Command{
	Use:   "reverse [text]",
	Short: "Reverse a string",
	RunE: func(cmd *cobra.Command, args []string) error {
		text := ""
		if len(args) > 0 {
			text = strings.Join(args, " ")
		} else {
			var err error
			text, err = utils.PromptInput("Enter text: ")
			if err != nil {
				return err
			}
		}

		runes := []rune(text)
		for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
			runes[i], runes[j] = runes[j], runes[i]
		}
		fmt.Println(string(runes))
		return nil
	},
}

var StringMaskCmd = &cobra.Command{
	Use:   "mask [text]",
	Short: "Mask parts of a string (e.g., for sensitive data)",
	RunE: func(cmd *cobra.Command, args []string) error {
		text := ""
		if len(args) > 0 {
			text = args[0]
		} else {
			var err error
			text, err = utils.PromptInput("Enter text to mask: ")
			if err != nil {
				return err
			}
		}

		if len(text) <= 4 {
			fmt.Println(strings.Repeat("*", len(text)))
			return nil
		}

		// Keep first 2 and last 2 chars visible
		visibleStart := 2
		visibleEnd := 2
		maskedLen := len(text) - visibleStart - visibleEnd

		if maskedLen < 0 {
			maskedLen = 0 // Should be covered by check above
		}

		result := text[:visibleStart] + strings.Repeat("*", maskedLen) + text[len(text)-visibleEnd:]
		fmt.Println(result)
		return nil
	},
}

var StringLengthCmd = &cobra.Command{
	Use:   "length [text]",
	Short: "Count characters, words, and lines",
	RunE: func(cmd *cobra.Command, args []string) error {
		text := ""
		if len(args) > 0 {
			text = args[0]
		} else {
			var err error
			text, err = utils.PromptInput("Enter text: ")
			if err != nil {
				return err
			}
		}

		chars := len([]rune(text))
		words := len(strings.Fields(text))
		lines := len(strings.Split(text, "\n"))
		if text == "" {
			lines = 0
		}

		fmt.Printf("Chars: %d\nWords: %d\nLines: %d\n", chars, words, lines)
		return nil
	},
}

var StringSlugifyCmd = &cobra.Command{
	Use:   "slugify [text]",
	Short: "Convert string to URL slug",
	RunE: func(cmd *cobra.Command, args []string) error {
		text := ""
		if len(args) > 0 {
			text = strings.Join(args, " ")
		} else {
			var err error
			text, err = utils.PromptInput("Enter text: ")
			if err != nil {
				return err
			}
		}

		// Simple slugify: lowercase, replace non-alphanum with dash
		text = strings.ToLower(text)
		reg := regexp.MustCompile("[^a-z0-9]+")
		text = reg.ReplaceAllString(text, "-")
		text = strings.Trim(text, "-")

		fmt.Println(text)
		return nil
	},
}

func init() {
	StringCmd.AddCommand(StringCaseCmd)
	StringCmd.AddCommand(StringEscapeCmd)
	StringCmd.AddCommand(StringUnescapeCmd)
	StringCmd.AddCommand(StringReverseCmd)
	StringCmd.AddCommand(StringMaskCmd)
	StringCmd.AddCommand(StringLengthCmd)
	StringCmd.AddCommand(StringSlugifyCmd)
}

// Helpers
func toSnakeCase(str string) string {
	str = strings.ReplaceAll(str, " ", "_")
	str = strings.ReplaceAll(str, "-", "_")
	var matchFirstCap = regexp.MustCompile("([^ _])([A-Z][a-z]+)")
	var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func toCamelCase(str string) string {
	// Split by non-alphanumeric
	parts := strings.FieldsFunc(str, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})

	for i, part := range parts {
		if i == 0 {
			parts[i] = strings.ToLower(part)
		} else {
			parts[i] = strings.Title(strings.ToLower(part))
		}
	}
	return strings.Join(parts, "")
}
