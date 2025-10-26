package utils

import (
	"bufio"
	"crypto/md5" // #nosec G501 - MD5 is intentionally supported for user choice in hash command
	"crypto/rand"
	"crypto/sha1" // #nosec G505 - SHA1 is intentionally supported for user choice in hash command
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	"os"
	"strings"
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
