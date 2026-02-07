package utils

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// SSHConfigEntry represents a host entry in ~/.ssh/config
type SSHConfigEntry struct {
	Name         string
	HostName     string
	User         string
	Port         string
	IdentityFile string
	// Add other fields as needed, but these are the main ones
	RawContent []string // Preserve comments and other options
}

// GetSSHConfigPath returns the path to the user's SSH config file
func GetSSHConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".ssh", "config"), nil
}

// ParseSSHConfig reads and parses the SSH config file
func ParseSSHConfig(path string) ([]SSHConfigEntry, error) {
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return []SSHConfigEntry{}, nil
		}
		return nil, err
	}
	defer file.Close()

	var entries []SSHConfigEntry
	var currentEntry *SSHConfigEntry
	var currentRaw []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, "Host ") {
			if currentEntry != nil {
				currentEntry.RawContent = currentRaw
				entries = append(entries, *currentEntry)
			}
			currentEntry = &SSHConfigEntry{
				Name: strings.TrimPrefix(trimmed, "Host "),
			}
			currentRaw = []string{line}
		} else if currentEntry != nil {
			currentRaw = append(currentRaw, line)
			parts := strings.Fields(trimmed)
			if len(parts) >= 2 {
				key := strings.ToLower(parts[0])
				value := parts[1]
				switch key {
				case "hostname":
					currentEntry.HostName = value
				case "user":
					currentEntry.User = value
				case "port":
					currentEntry.Port = value
				case "identityfile":
					currentEntry.IdentityFile = value
				}
			}
		} else {
			// Lines before the first Host entry (globals or comments)
			// For now we might lose them or attach them to a "global" entry if we wanted to be perfect.
			// But for this specific task, we'll focus on Host entries.
		}
	}

	if currentEntry != nil {
		currentEntry.RawContent = currentRaw
		entries = append(entries, *currentEntry)
	}

	return entries, scanner.Err()
}

// WriteSSHConfig writes entries back to the config file
// Note: This is a simplified writer that appends new entries or rewrites.
// For robust editing, we need to be careful not to destroy existing formatting excessively.
func WriteSSHConfig(path string, entries []SSHConfigEntry) error {
	// For now, simpler approach: Read valid existing, append if new, or rewrite all.
	// Given the requirement to "Edit", rewriting seems necessary.

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, entry := range entries {
		if len(entry.RawContent) > 0 {
			for _, line := range entry.RawContent {
				fmt.Fprintln(w, line)
			}
		} else {
			// Reconstruct if RawContent is empty (new entry)
			fmt.Fprintf(w, "Host %s\n", entry.Name)
			if entry.HostName != "" {
				fmt.Fprintf(w, "  HostName %s\n", entry.HostName)
			}
			if entry.User != "" {
				fmt.Fprintf(w, "  User %s\n", entry.User)
			}
			if entry.Port != "" {
				fmt.Fprintf(w, "  Port %s\n", entry.Port)
			}
			if entry.IdentityFile != "" {
				fmt.Fprintf(w, "  IdentityFile %s\n", entry.IdentityFile)
			}
		}
		fmt.Fprintln(w, "") // Empty line between hosts
	}
	return w.Flush()
}

// AddSSHConfigEntry appends a new entry to the config
func AddSSHConfigEntry(path string, entry SSHConfigEntry) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	fmt.Fprintln(w, "") // Ensure separation
	fmt.Fprintf(w, "Host %s\n", entry.Name)
	if entry.HostName != "" {
		fmt.Fprintf(w, "  HostName %s\n", entry.HostName)
	}
	if entry.User != "" {
		fmt.Fprintf(w, "  User %s\n", entry.User)
	}
	if entry.Port != "" {
		fmt.Fprintf(w, "  Port %s\n", entry.Port)
	}
	if entry.IdentityFile != "" {
		fmt.Fprintf(w, "  IdentityFile %s\n", entry.IdentityFile)
	}
	return w.Flush()
}

// RemoveSSHConfigEntry removes a host entry from the config
func RemoveSSHConfigEntry(path string, hostName string) error {
	entries, err := ParseSSHConfig(path)
	if err != nil {
		return err
	}

	var newEntries []SSHConfigEntry
	found := false
	for _, entry := range entries {
		if entry.Name == hostName {
			found = true
			continue
		}
		newEntries = append(newEntries, entry)
	}

	if !found {
		return fmt.Errorf("host '%s' not found", hostName)
	}

	return WriteSSHConfig(path, newEntries)
}

// GetKnownHostsPath returns the path to known_hosts
func GetKnownHostsPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".ssh", "known_hosts"), nil
}

// RemoveKnownHost removes a host from known_hosts using ssh-keygen
func RemoveKnownHost(host string) error {
	cmd := exec.Command("ssh-keygen", "-R", host)
	return cmd.Run()
}

// SSHKey represents an SSH key pair
type SSHKey struct {
	Name    string
	Path    string
	PubPath string
	Type    string // e.g., "RSA", "ED25519"
}

// ListSSHKeys lists keys in ~/.ssh
func ListSSHKeys() ([]SSHKey, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	sshDir := filepath.Join(home, ".ssh")

	files, err := ioutil.ReadDir(sshDir)
	if err != nil {
		return nil, err
	}

	var keys []SSHKey
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		name := f.Name()
		// Naive check: private keys usually don't have extension, public has .pub
		if !strings.HasSuffix(name, ".pub") && !strings.HasSuffix(name, "known_hosts") && !strings.HasSuffix(name, "config") && !strings.HasPrefix(name, ".") {
			// Check if corresponding .pub exists
			pubPath := filepath.Join(sshDir, name+".pub")
			if _, err := os.Stat(pubPath); err == nil {
				keys = append(keys, SSHKey{
					Name:    name,
					Path:    filepath.Join(sshDir, name),
					PubPath: pubPath,
				})
			}
		}
	}
	return keys, nil
}

// GenerateSSHKey generates a new SSH key pair
func GenerateSSHKey(path string, keyType string, comment string, passphrase string, rounds int) error {
	args := []string{"-t", keyType, "-f", path}
	if comment != "" {
		args = append(args, "-C", comment)
	}
	if passphrase != "" {
		args = append(args, "-N", passphrase)
	} else {
		args = append(args, "-N", "") // Empty passphrase
	}
	if rounds > 0 {
		args = append(args, "-a", fmt.Sprintf("%d", rounds))
	}

	cmd := exec.Command("ssh-keygen", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ssh-keygen failed: %s: %s", err, string(output))
	}
	return nil
}
