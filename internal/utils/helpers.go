package utils

import (
	"bufio"
	"crypto/md5" // #nosec G501 - MD5 is intentionally supported for user choice in hash command
	"crypto/rand"
	"crypto/sha1" // #nosec G505 - SHA1 is intentionally supported for user choice in hash command
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"math/big"
	"os"
	"strings"

	"github.com/atotto/clipboard"
)

// PromptInput prompts the user for input with a message
func PromptInput(message string) (string, error) {
	fmt.Print(message)
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(input), nil
}

// FileExists checks if a file exists
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// HashString generates a hash of the input string using the specified algorithm.
// MD5 and SHA1 are supported for legacy compatibility but SHA256 is recommended.
// #nosec G401 - MD5/SHA1 are intentionally supported as user options
func HashString(input, algorithm string) (string, error) {
	var hash []byte

	switch strings.ToLower(algorithm) {
	case "md5":
		h := md5.Sum([]byte(input)) // #nosec G401
		hash = h[:]
	case "sha1":
		h := sha1.Sum([]byte(input)) // #nosec G401
		hash = h[:]
	case "sha256":
		h := sha256.Sum256([]byte(input))
		hash = h[:]
	default:
		return "", fmt.Errorf("unsupported algorithm: %s", algorithm)
	}

	return hex.EncodeToString(hash), nil
}

// GeneratePassword generates a random password of the specified length
func GeneratePassword(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()_+-=[]{}|;:,.<>?"

	password := make([]byte, length)
	for i := range password {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		password[i] = charset[num.Int64()]
	}

	return string(password), nil
}

// WriteFile writes content to a file with proper error handling.
// Uses 0600 permissions for security (owner read/write only).
func WriteFile(path, content string) error {
	return os.WriteFile(path, []byte(content), 0o600)
}

// ReadFile reads content from a file
func ReadFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// CopyToClipboard copies text to clipboard
func CopyToClipboard(text string) error {
	return clipboard.WriteAll(text)
}

// HashFile generates a hash of a file using the specified algorithm
func HashFile(filePath, algorithm string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var hasher io.Writer
	switch strings.ToLower(algorithm) {
	case "md5":
		h := md5.New() // #nosec G401
		hasher = h
		if _, err := io.Copy(hasher, file); err != nil {
			return "", err
		}
		return hex.EncodeToString(h.Sum(nil)), nil
	case "sha1":
		h := sha1.New() // #nosec G401
		hasher = h
		if _, err := io.Copy(hasher, file); err != nil {
			return "", err
		}
		return hex.EncodeToString(h.Sum(nil)), nil
	case "sha256":
		h := sha256.New()
		hasher = h
		if _, err := io.Copy(hasher, file); err != nil {
			return "", err
		}
		return hex.EncodeToString(h.Sum(nil)), nil
	default:
		return "", fmt.Errorf("unsupported algorithm: %s", algorithm)
	}
}

// GeneratePasswordByType generates a password based on type
func GeneratePasswordByType(length int, passwordType string) (string, error) {
	var charset string

	switch strings.ToLower(passwordType) {
	case "alpha":
		charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	case "numeral":
		charset = "0123456789"
	case "alphanum":
		charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	case "alphanum+special":
		charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()_+-=[]{}|;:,.<>?"
	case "sentence":
		// Generate words separated by spaces
		words := []string{"quick", "brown", "fox", "jumps", "over", "lazy", "dog", "happy", "sunny", "bright", "cloud", "river", "mountain", "forest", "ocean"}
		numWords := length / 6 // Approximate word count
		if numWords < 3 {
			numWords = 3
		}
		var sentence []string
		for i := 0; i < numWords; i++ {
			idx, err := rand.Int(rand.Reader, big.NewInt(int64(len(words))))
			if err != nil {
				return "", err
			}
			sentence = append(sentence, words[idx.Int64()])
		}
		return strings.Join(sentence, " "), nil
	default:
		return "", fmt.Errorf("unsupported password type: %s", passwordType)
	}

	password := make([]byte, length)
	for i := range password {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		password[i] = charset[num.Int64()]
	}

	return string(password), nil
}

// PromptSelect prompts the user to select from a list of options
func PromptSelect(message string, options []string) (string, error) {
	fmt.Println(message)
	for i, option := range options {
		fmt.Printf("%d. %s\n", i+1, option)
	}
	fmt.Print("Enter choice (number): ")

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	var choice int
	_, err = fmt.Sscanf(strings.TrimSpace(input), "%d", &choice)
	if err != nil || choice < 1 || choice > len(options) {
		return "", fmt.Errorf("invalid choice")
	}

	return options[choice-1], nil
}
