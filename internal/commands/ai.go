package commands

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/MayR-Labs/mayrlabs-go/internal/utils"
	"github.com/google/generative-ai-go/genai"
	"github.com/spf13/cobra"
	"google.golang.org/api/option"
)

// AISetupCmd sets up the Gemini API key
var AISetupCmd = &cobra.Command{
	Use:   "ai-setup",
	Short: "Setup Gemini API key for AI features",
	Long:  "Prompt the user to provide their Gemini API key and store it securely",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("🤖 Setting up AI with Gemini API")
		fmt.Println("Get your API key from: https://aistudio.google.com/app/apikey")

		apiKey, err := utils.SurveyPassword("Enter your Gemini API key:")
		if err != nil {
			return err
		}

		if apiKey == "" {
			return fmt.Errorf("API key cannot be empty")
		}

		// Validate the API key by making a test call
		fmt.Println("\nValidating API key...")
		ctx := context.Background()
		client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}
		defer func(client *genai.Client) {
			err := client.Close()
			if err != nil {
				return
			}
		}(client)

		// Try to list models to validate the key
		model := client.GenerativeModel("gemini-2.0-flash-exp")
		if model == nil {
			return fmt.Errorf("failed to get model")
		}

		// Store the API key
		if err := utils.StoreAPIKey(apiKey); err != nil {
			return err
		}

		fmt.Println("✅ API key stored successfully!")
		fmt.Println("You can now use 'mayrlabs ai' to query the AI")
		return nil
	},
}

// AICmd queries the AI with Gemini
var AICmd = &cobra.Command{
	Use:   "ai [query...]",
	Short: "Query the AI using Gemini",
	Long:  "Send a query to Gemini AI and get a response. If no query is provided, a multiline input form is shown",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get API key
		apiKey, err := utils.GetAPIKey()
		if err != nil {
			return err
		}

		// Get query
		var query string
		if len(args) == 0 {
			query, err = utils.SurveyMultiline("Enter your query:")
			if err != nil {
				return err
			}
		} else {
			query = strings.Join(args, " ")
		}

		if query == "" {
			return fmt.Errorf("query cannot be empty")
		}

		fmt.Println("\n🤖 Processing your query...")

		// Query Gemini
		ctx := context.Background()
		client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}
		defer func(client *genai.Client) {
			err := client.Close()
			if err != nil {
				return
			}
		}(client)

		model := client.GenerativeModel("gemini-2.0-flash-exp")
		resp, err := model.GenerateContent(ctx, genai.Text(query))
		if err != nil {
			return fmt.Errorf("failed to generate content: %w", err)
		}

		// Print the response
		fmt.Println("\n📝 Response:")
		if resp != nil && len(resp.Candidates) > 0 {
			for _, part := range resp.Candidates[0].Content.Parts {
				fmt.Println(part)
			}
		}

		return nil
	},
}

// AIClearCmd clears the stored API key
var AIClearCmd = &cobra.Command{
	Use:   "ai-clear",
	Short: "Clear the stored Gemini API key",
	Long:  "Remove the stored Gemini API key from the system",
	RunE: func(cmd *cobra.Command, args []string) error {
		confirm, err := utils.SurveyConfirm("Are you sure you want to clear the API key?", false)
		if err != nil {
			return err
		}

		if !confirm {
			fmt.Println("Operation cancelled")
			return nil
		}

		if err := utils.ClearAPIKey(); err != nil {
			return err
		}

		fmt.Println("✅ API key cleared successfully!")
		return nil
	},
}

// AIAliasCmd creates a permanent alias for mayrlabs ai command
var AIAliasCmd = &cobra.Command{
	Use:   "ai-alias [alias]",
	Short: "Create a permanent alias for 'mayrlabs ai' command",
	Long:  "Create a shell alias for 'mayrlabs ai' command in your shell configuration file. Default alias is 'mlai'",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var aliasName string
		var err error

		// Default to mlai
		if len(args) == 0 {
			aliasName = "mlai"
		} else {
			aliasName = args[0]
		}

		// Validate alias name
		if strings.Contains(aliasName, " ") || strings.Contains(aliasName, "/") {
			return fmt.Errorf("invalid alias name: %s", aliasName)
		}

		// Detect shell configuration file
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get home directory: %w", err)
		}

		// Determine which shell config file to use
		shell := os.Getenv("SHELL")
		var configFile string

		switch {
		case strings.Contains(shell, "zsh"):
			configFile = filepath.Join(homeDir, ".zshrc")
		case strings.Contains(shell, "bash"):
			configFile = filepath.Join(homeDir, ".bashrc")
		case strings.Contains(shell, "fish"):
			configFile = filepath.Join(homeDir, ".config", "fish", "config.fish")
		default:
			// Default to .bashrc
			configFile = filepath.Join(homeDir, ".bashrc")
		}

		// Ask for confirmation
		confirm, err := utils.SurveyConfirm(fmt.Sprintf("Add alias '%s' to %s?", aliasName, configFile), true)
		if err != nil {
			return err
		}
		if !confirm {
			return fmt.Errorf("operation cancelled")
		}

		// Get the current executable path
		execPath, err := os.Executable()
		if err != nil {
			return fmt.Errorf("failed to get executable path: %w", err)
		}

		// Prepare alias line
		aliasLine := fmt.Sprintf("\nalias %s='%s ai'\n", aliasName, execPath)

		// Append to config file
		file, err := os.OpenFile(configFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
		if err != nil {
			return fmt.Errorf("failed to open config file: %w", err)
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				return
			}
		}(file)

		if _, err := file.WriteString(aliasLine); err != nil {
			return fmt.Errorf("failed to write alias: %w", err)
		}

		fmt.Printf("✅ Alias '%s' added to %s\n", aliasName, configFile)
		fmt.Printf("Run 'source %s' or restart your terminal to use the alias\n", configFile)
		return nil
	},
}

// AIFileCmd sends file content to AI
var AIFileCmd = &cobra.Command{
	Use:   "ai-file [file] [message]",
	Short: "Send the content of a text-based file to the AI",
	Long:  "Read a text file and send its content to Gemini AI with an optional message",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get API key
		apiKey, err := utils.GetAPIKey()
		if err != nil {
			return err
		}

		// Get file path
		var filePath string
		if len(args) == 0 {
			filePath, err = utils.SurveyInput("Enter file path:", "")
			if err != nil {
				return err
			}
		} else {
			filePath = args[0]
		}

		if filePath == "" {
			return fmt.Errorf("file path cannot be empty")
		}

		// Check if file exists
		if !utils.FileExists(filePath) {
			return fmt.Errorf("file not found: %s", filePath)
		}

		// Read file content
		content, err := utils.ReadFile(filePath)
		if err != nil {
			return err
		}

		// Get message
		var message string
		if len(args) > 1 {
			message = strings.Join(args[1:], " ")
		} else {
			message, err = utils.SurveyInput("Enter your message/question about the file:", "")
			if err != nil {
				return err
			}
		}

		if message == "" {
			message = "Please analyze this file:"
		}

		fmt.Println("\n🤖 Processing your file...")

		// Query Gemini
		ctx := context.Background()
		client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
		if err != nil {
			return fmt.Errorf("failed to create client: %w", err)
		}
		defer func(client *genai.Client) {
			err := client.Close()
			if err != nil {
				return
			}
		}(client)

		model := client.GenerativeModel("gemini-2.0-flash-exp")
		prompt := fmt.Sprintf("%s\n\nFile: %s\n\n```\n%s\n```", message, filePath, content)
		resp, err := model.GenerateContent(ctx, genai.Text(prompt))
		if err != nil {
			return fmt.Errorf("failed to generate content: %w", err)
		}

		// Print the response
		fmt.Println("\n📝 Response:")
		if resp != nil && len(resp.Candidates) > 0 {
			for _, part := range resp.Candidates[0].Content.Parts {
				fmt.Println(part)
			}
		}

		return nil
	},
}
