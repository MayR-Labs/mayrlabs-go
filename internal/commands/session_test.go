package commands

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
		want     string
	}{
		{
			name:     "Less than a minute",
			duration: 30 * time.Second,
			want:     "30s",
		},
		{
			name:     "A few minutes",
			duration: 5*time.Minute + 30*time.Second,
			want:     "5m 30s",
		},
		{
			name:     "More than an hour",
			duration: 2*time.Hour + 15*time.Minute + 45*time.Second,
			want:     "2h 15m 45s",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatDuration(tt.duration)
			if got != tt.want {
				t.Errorf("formatDuration() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetSessionEndMessage(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
		contains string
	}{
		{
			name:     "Short session",
			duration: 10 * time.Minute,
			contains: "quick session",
		},
		{
			name:     "Medium session",
			duration: 30 * time.Minute,
			contains: "great session",
		},
		{
			name:     "Long session",
			duration: 90 * time.Minute,
			contains: "intense session",
		},
		{
			name:     "Epic session",
			duration: 200 * time.Minute,
			contains: "Epic session",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getSessionEndMessage(tt.duration)
			if !strings.Contains(strings.ToLower(got), strings.ToLower(tt.contains)) {
				t.Errorf("getSessionEndMessage() = %v, should contain %v", got, tt.contains)
			}
		})
	}
}

func TestEncryptDecryptSessionFile(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test-session.md")

	// Create test content
	testContent := "# Test Session\n\nThis is a test session file.\n\n**Started:** 2024-01-01 10:00:00\n"
	err := os.WriteFile(testFile, []byte(testContent), 0o644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Test encryption
	password := "test-password-123"
	err = encryptSessionFile(testFile, password)
	if err != nil {
		t.Fatalf("encryptSessionFile() error = %v", err)
	}

	// Verify file is encrypted (not readable as plain text)
	encryptedContent, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read encrypted file: %v", err)
	}

	if string(encryptedContent) == testContent {
		t.Error("File content should be encrypted but appears to be plain text")
	}

	// Test decryption
	decrypted, err := decryptSessionFile(testFile, password)
	if err != nil {
		t.Fatalf("decryptSessionFile() error = %v", err)
	}

	if string(decrypted) != testContent {
		t.Errorf("decryptSessionFile() = %v, want %v", string(decrypted), testContent)
	}

	// Test decryption with wrong password
	wrongDecrypted, err := decryptSessionFile(testFile, "wrong-password")
	if err != nil {
		t.Fatalf("decryptSessionFile() with wrong password error = %v", err)
	}

	// The content should be different with wrong password
	if string(wrongDecrypted) == testContent {
		t.Error("Decryption with wrong password should not produce original content")
	}
}

func TestEncryptDecryptEmptyPassword(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test-session.md")

	// Create test content
	testContent := "# Test Session\n\nContent\n"
	err := os.WriteFile(testFile, []byte(testContent), 0o644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Test encryption with empty password (should fail)
	password := ""
	err = encryptSessionFile(testFile, password)
	if err == nil {
		t.Error("encryptSessionFile() with empty password should return error")
	}

	// Test decryption with empty password (should fail)
	_, err = decryptSessionFile(testFile, password)
	if err == nil {
		t.Error("decryptSessionFile() with empty password should return error")
	}
}

func TestSlugify(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		maxLength int
		want      string
	}{
		{
			name:      "Simple text",
			input:     "Hello World",
			maxLength: 50,
			want:      "hello-world",
		},
		{
			name:      "With special characters",
			input:     "Bug Fix: Authentication!",
			maxLength: 50,
			want:      "bug-fix-authentication",
		},
		{
			name:      "Long text truncation",
			input:     "This is a very long summary that needs to be truncated",
			maxLength: 20,
			want:      "this-is-a-very-long",
		},
		{
			name:      "Multiple spaces",
			input:     "Multiple   Spaces   Here",
			maxLength: 50,
			want:      "multiple-spaces-here",
		},
		{
			name:      "Numbers and letters",
			input:     "Feature 123 User Dashboard",
			maxLength: 50,
			want:      "feature-123-user-dashboard",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := slugify(tt.input, tt.maxLength)
			if got != tt.want {
				t.Errorf("slugify() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtractSummaryFromFilename(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		want     string
	}{
		{
			name:     "New format with summary",
			filename: "session-20240126-100000-bug-fix-authentication.md",
			want:     "Bug fix authentication",
		},
		{
			name:     "New format with long summary",
			filename: "session-20240126-100000-feature-user-dashboard-redesign.md",
			want:     "Feature user dashboard redesign",
		},
		{
			name:     "Old format without summary",
			filename: "session-20240126-100000.md",
			want:     "",
		},
		{
			name:     "Single word summary",
			filename: "session-20240126-100000-refactor.md",
			want:     "Refactor",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractSummaryFromFilename(tt.filename)
			if got != tt.want {
				t.Errorf("extractSummaryFromFilename() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormatSessionOption(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		want     string
	}{
		{
			name:     "With summary",
			filename: "session-20240126-100000-bug-fix.md",
			want:     "session-20240126-100000-bug-fix.md - Bug fix",
		},
		{
			name:     "Without summary",
			filename: "session-20240126-100000.md",
			want:     "session-20240126-100000.md",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatSessionOption(tt.filename)
			if got != tt.want {
				t.Errorf("formatSessionOption() = %v, want %v", got, tt.want)
			}
		})
	}
}
