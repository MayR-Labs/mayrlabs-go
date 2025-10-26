package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/AlecAivazis/survey/v2"
)

// GetConfigDir returns the directory for storing mayrlabs configuration
func GetConfigDir() (string, error) {
	var configDir string

	switch runtime.GOOS {
	case "windows":
		configDir = filepath.Join(os.Getenv("APPDATA"), "mayrlabs")
	case "darwin":
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get home directory: %w", err)
		}
		configDir = filepath.Join(homeDir, ".config", "mayrlabs")
	case "linux":
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get home directory: %w", err)
		}
		configDir = filepath.Join(homeDir, ".config", "mayrlabs")
	default:
		return "", fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	// Create config directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0o700); err != nil {
		return "", fmt.Errorf("failed to create config directory: %w", err)
	}

	return configDir, nil
}

// StoreAPIKey stores the Gemini API key securely
func StoreAPIKey(apiKey string) error {
	configDir, err := GetConfigDir()
	if err != nil {
		return err
	}

	keyFile := filepath.Join(configDir, "gemini_api_key")
	if err := os.WriteFile(keyFile, []byte(apiKey), 0o600); err != nil {
		return fmt.Errorf("failed to store API key: %w", err)
	}

	return nil
}

// GetAPIKey retrieves the stored Gemini API key
func GetAPIKey() (string, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return "", err
	}

	keyFile := filepath.Join(configDir, "gemini_api_key")
	data, err := os.ReadFile(keyFile)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("API key not found. Please run 'mayrlabs ai-setup' first")
		}
		return "", fmt.Errorf("failed to read API key: %w", err)
	}

	return string(data), nil
}

// ClearAPIKey removes the stored Gemini API key
func ClearAPIKey() error {
	configDir, err := GetConfigDir()
	if err != nil {
		return err
	}

	keyFile := filepath.Join(configDir, "gemini_api_key")
	if err := os.Remove(keyFile); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("no API key stored")
		}
		return fmt.Errorf("failed to clear API key: %w", err)
	}

	return nil
}

// SurveyPassword prompts the user for a password input
func SurveyPassword(message string) (string, error) {
	var result string
	prompt := &survey.Password{
		Message: message,
	}
	err := survey.AskOne(prompt, &result)
	return result, err
}

// SurveyMultiline prompts the user for multiline input
func SurveyMultiline(message string) (string, error) {
	var result string
	prompt := &survey.Multiline{
		Message: message,
	}
	err := survey.AskOne(prompt, &result)
	return result, err
}
