package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetConfigDir(t *testing.T) {
	configDir, err := GetConfigDir()
	if err != nil {
		t.Fatalf("GetConfigDir() error = %v", err)
	}

	if configDir == "" {
		t.Error("GetConfigDir() returned empty string")
	}

	// Check that directory was created
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		t.Errorf("Config directory was not created: %s", configDir)
	}
}

func TestStoreAndGetAPIKey(t *testing.T) {
	// Create a temporary config dir for testing
	tmpDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	if originalHome == "" {
		originalHome = os.Getenv("USERPROFILE") // Windows
	}

	// Set up a temporary home directory
	testHome := filepath.Join(tmpDir, "test_home")
	if err := os.MkdirAll(testHome, 0o700); err != nil {
		t.Fatalf("Failed to create test home: %v", err)
	}

	// Save and restore HOME/USERPROFILE
	defer func() {
		if originalHome != "" {
			os.Setenv("HOME", originalHome)
			os.Setenv("USERPROFILE", originalHome)
		}
	}()
	os.Setenv("HOME", testHome)
	os.Setenv("USERPROFILE", testHome)

	testKey := "test-api-key-12345"

	// Test storing API key
	err := StoreAPIKey(testKey)
	if err != nil {
		t.Fatalf("StoreAPIKey() error = %v", err)
	}

	// Test retrieving API key
	retrievedKey, err := GetAPIKey()
	if err != nil {
		t.Fatalf("GetAPIKey() error = %v", err)
	}

	if retrievedKey != testKey {
		t.Errorf("GetAPIKey() = %v, want %v", retrievedKey, testKey)
	}

	// Test clearing API key
	err = ClearAPIKey()
	if err != nil {
		t.Fatalf("ClearAPIKey() error = %v", err)
	}

	// Verify key is cleared
	_, err = GetAPIKey()
	if err == nil {
		t.Error("GetAPIKey() should return error after clearing, but didn't")
	}
}

func TestClearAPIKeyNotFound(t *testing.T) {
	// Create a temporary config dir for testing
	tmpDir := t.TempDir()
	testHome := filepath.Join(tmpDir, "test_home")
	if err := os.MkdirAll(testHome, 0o700); err != nil {
		t.Fatalf("Failed to create test home: %v", err)
	}

	originalHome := os.Getenv("HOME")
	if originalHome == "" {
		originalHome = os.Getenv("USERPROFILE")
	}

	defer func() {
		if originalHome != "" {
			os.Setenv("HOME", originalHome)
			os.Setenv("USERPROFILE", originalHome)
		}
	}()
	os.Setenv("HOME", testHome)
	os.Setenv("USERPROFILE", testHome)

	// Try to clear non-existent API key
	err := ClearAPIKey()
	if err == nil {
		t.Error("ClearAPIKey() should return error when no key exists, but didn't")
	}
}

func TestGetAPIKeyNotFound(t *testing.T) {
	// Create a temporary config dir for testing
	tmpDir := t.TempDir()
	testHome := filepath.Join(tmpDir, "test_home")
	if err := os.MkdirAll(testHome, 0o700); err != nil {
		t.Fatalf("Failed to create test home: %v", err)
	}

	originalHome := os.Getenv("HOME")
	if originalHome == "" {
		originalHome = os.Getenv("USERPROFILE")
	}

	defer func() {
		if originalHome != "" {
			os.Setenv("HOME", originalHome)
			os.Setenv("USERPROFILE", originalHome)
		}
	}()
	os.Setenv("HOME", testHome)
	os.Setenv("USERPROFILE", testHome)

	// Try to get non-existent API key
	_, err := GetAPIKey()
	if err == nil {
		t.Error("GetAPIKey() should return error when no key exists, but didn't")
	}
}
