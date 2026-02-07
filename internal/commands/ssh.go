package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/MayR-Labs/mayrlabs-go/internal/utils"
	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"
)

// SSHCmd represents the base ssh command
var SSHCmd = &cobra.Command{
	Use:   "ssh",
	Short: "Manage SSH configurations, keys, and connections",
	Long:  "A comprehensive tool for managing your SSH config, keys, known_hosts and connections.\nRunning without arguments opens an interactive menu.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return runSSHInteractive()
		}
		return cmd.Help()
	},
}

func runSSHInteractive() error {
	for {
		action := ""
		prompt := &survey.Select{
			Message: "SSH Management:",
			Options: []string{"Connect", "Config", "Keys", "Hosts", "Exit"},
		}
		if err := survey.AskOne(prompt, &action); err != nil {
			return err
		}

		switch action {
		case "Connect":
			if err := SSHConnectCmd.RunE(SSHConnectCmd, []string{}); err != nil {
				fmt.Printf("Error: %v\n", err)
			}
		case "Config":
			if err := runSSHConfigInteractive(); err != nil {
				fmt.Printf("Error: %v\n", err)
			}
		case "Keys":
			if err := runSSHKeysInteractive(); err != nil {
				fmt.Printf("Error: %v\n", err)
			}
		case "Hosts":
			if err := runSSHHostsInteractive(); err != nil {
				fmt.Printf("Error: %v\n", err)
			}
		case "Exit":
			return nil
		}
	}
}

func runSSHConfigInteractive() error {
	for {
		action := ""
		prompt := &survey.Select{
			Message: "SSH Config Management:",
			Options: []string{"List", "Add", "Edit", "Remove", "Export", "Import", "Connect", "Back"},
		}
		if err := survey.AskOne(prompt, &action); err != nil {
			return err
		}

		var err error
		switch action {
		case "List":
			err = SSHConfigListCmd.RunE(SSHConfigListCmd, []string{})
		case "Add":
			err = SSHConfigAddCmd.RunE(SSHConfigAddCmd, []string{})
		case "Edit":
			err = SSHConfigEditCmd.RunE(SSHConfigEditCmd, []string{})
		case "Remove":
			err = SSHConfigRemoveCmd.RunE(SSHConfigRemoveCmd, []string{})
		case "Export":
			err = SSHConfigExportCmd.RunE(SSHConfigExportCmd, []string{})
		case "Import":
			err = SSHConfigImportCmd.RunE(SSHConfigImportCmd, []string{})
		case "Connect":
			err = SSHConfigConnectCmd.RunE(SSHConfigConnectCmd, []string{})
		case "Back":
			return nil
		}
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	}
}

func runSSHKeysInteractive() error {
	for {
		action := ""
		prompt := &survey.Select{
			Message: "SSH Keys Management:",
			Options: []string{"List", "Create", "Back"},
		}
		if err := survey.AskOne(prompt, &action); err != nil {
			return err
		}

		var err error
		switch action {
		case "List":
			err = SSHKeysListCmd.RunE(SSHKeysListCmd, []string{})
		case "Create":
			err = SSHKeysCreateCmd.RunE(SSHKeysCreateCmd, []string{})
		case "Back":
			return nil
		}
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	}
}

func runSSHHostsInteractive() error {
	for {
		action := ""
		prompt := &survey.Select{
			Message: "Known Hosts Management:",
			Options: []string{"List", "Remove", "Back"},
		}
		if err := survey.AskOne(prompt, &action); err != nil {
			return err
		}

		var err error
		switch action {
		case "List":
			err = SSHHostsListCmd.RunE(SSHHostsListCmd, []string{})
		case "Remove":
			err = SSHHostsRemoveCmd.RunE(SSHHostsRemoveCmd, []string{})
		case "Back":
			return nil
		}
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	}
}

// SSHConfigCmd represents the ssh config command category
var SSHConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage SSH config file (~/.ssh/config)",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return runSSHConfigInteractive()
		}
		return cmd.Help()
	},
}

// SSHConfigListCmd lists entries in ~/.ssh/config
var SSHConfigListCmd = &cobra.Command{
	Use:   "list",
	Short: "List SSH config hosts",
	RunE: func(cmd *cobra.Command, args []string) error {
		path, err := utils.GetSSHConfigPath()
		if err != nil {
			return err
		}
		entries, err := utils.ParseSSHConfig(path)
		if err != nil {
			return err
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		if _, err := fmt.Fprintln(w, "NAME\tHOSTNAME\tUSER\tPORT\tIDENTITY FILE"); err != nil {
			return err
		}
		for _, e := range entries {
			if _, err := fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", e.Name, e.HostName, e.User, e.Port, e.IdentityFile); err != nil {
				return err
			}
		}
		return w.Flush()
	},
}

// helper to prompt for identity file
func promptIdentityFile(current string) (string, error) {
	// List available keys
	keys, err := utils.ListSSHKeys()
	if err != nil {
		return "", err
	}

	options := []string{"Skip"}
	if current != "" {
		options = append(options, "Keep Current ("+current+")")
	}
	options = append(options, "Custom Path")

	for _, k := range keys {
		options = append(options, k.Path)
	}

	var choice string
	if err := survey.AskOne(&survey.Select{
		Message: "Identity File:",
		Options: options,
	}, &choice); err != nil {
		return "", err
	}

	if strings.HasPrefix(choice, "Keep Current") {
		return current, nil
	}
	if choice == "Skip" {
		return "", nil
	}
	if choice == "Custom Path" {
		var path string
		if err := survey.AskOne(&survey.Input{Message: "Enter path to private key:"}, &path); err != nil {
			return "", err
		}
		return path, nil
	}
	return choice, nil
}

// helper for auth options
func promptAuthOptions(entry utils.SSHConfigEntry) (utils.SSHConfigEntry, error) {
	// PreferredAuthentications
	prefOpts := []string{"Skip", "publickey", "password", "keyboard-interactive", "hostbased", "Custom"}
	var prefChoice string
	if err := survey.AskOne(&survey.Select{
		Message: "PreferredAuthentications:",
		Options: prefOpts,
		Default: "Skip",
	}, &prefChoice); err != nil {
		return entry, err
	}
	switch prefChoice {
	case "Custom":
		if err := survey.AskOne(&survey.Input{Message: "Enter PreferredAuthentications:"}, &prefChoice); err != nil {
			return entry, err
		}
	case "Skip":
		prefChoice = ""
	}
	entry.PreferredAuthentications = prefChoice

	// Bool/Skip helper
	askBool := func(msg string, current string) (string, error) {
		opts := []string{"Skip", "yes", "no"}
		var c string
		if err := survey.AskOne(&survey.Select{
			Message: msg,
			Options: opts,
			Default: "Skip",
		}, &c); err != nil {
			return "", err
		}
		if c == "Skip" {
			return "", nil
		}
		return c, nil
	}

	var err error
	entry.PubkeyAuthentication, err = askBool("PubkeyAuthentication:", entry.PubkeyAuthentication)
	if err != nil {
		return entry, err
	}
	entry.PasswordAuthentication, err = askBool("PasswordAuthentication:", entry.PasswordAuthentication)
	if err != nil {
		return entry, err
	}
	entry.KbdInteractiveAuthentication, err = askBool("KbdInteractiveAuthentication:", entry.KbdInteractiveAuthentication)
	if err != nil {
		return entry, err
	}

	return entry, nil
}

// SSHConfigAddCmd adds a new entry to ~/.ssh/config
var SSHConfigAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new SSH host config",
	RunE: func(cmd *cobra.Command, args []string) error {
		answers := struct {
			Name     string
			HostName string
			User     string
			Port     string
		}{}

		var qs = []*survey.Question{
			{
				Name: "Name",
				Prompt: &survey.Input{
					Message: "Host Alias (e.g. myserver):",
				},
				Validate: survey.Required,
			},
			{
				Name: "HostName",
				Prompt: &survey.Input{
					Message: "Host Name (IP or domain):",
				},
				Validate: survey.Required,
			},
			{
				Name: "User",
				Prompt: &survey.Input{
					Message: "User (optional):",
				},
			},
			{
				Name: "Port",
				Prompt: &survey.Input{
					Message: "Port (optional, default 22):",
				},
			},
		}

		if err := survey.Ask(qs, &answers); err != nil {
			return err
		}

		idFile, err := promptIdentityFile("")
		if err != nil {
			return err
		}

		entry := utils.SSHConfigEntry{
			Name:         answers.Name,
			HostName:     answers.HostName,
			User:         answers.User,
			Port:         answers.Port,
			IdentityFile: idFile,
		}

		// Ask for advanced options
		entry, err = promptAuthOptions(entry)
		if err != nil {
			return err
		}

		fmt.Println("\nNew Entry Preview:")
		fmt.Printf("Host %s\n", entry.Name)
		fmt.Printf("  HostName %s\n", entry.HostName)
		if entry.User != "" {
			fmt.Printf("  User %s\n", entry.User)
		}
		if entry.Port != "" {
			fmt.Printf("  Port %s\n", entry.Port)
		}
		if entry.IdentityFile != "" {
			fmt.Printf("  IdentityFile %s\n", entry.IdentityFile)
		}
		if entry.PreferredAuthentications != "" {
			fmt.Printf("  PreferredAuthentications %s\n", entry.PreferredAuthentications)
		}

		confirm := false
		prompt := &survey.Confirm{
			Message: "Add this entry to config?",
		}
		if err := survey.AskOne(prompt, &confirm); err != nil {
			return err
		}

		if confirm {
			path, err := utils.GetSSHConfigPath()
			if err != nil {
				return err
			}
			if err := utils.AddSSHConfigEntry(path, entry); err != nil {
				return err
			}
			fmt.Println("✅ Entry added successfully.")
		} else {
			fmt.Println("❌ Operation cancelled.")
		}
		return nil
	},
}

// SSHConfigRemoveCmd removes an entry from ~/.ssh/config
var SSHConfigRemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove an SSH host config",
	RunE: func(cmd *cobra.Command, args []string) error {
		path, err := utils.GetSSHConfigPath()
		if err != nil {
			return err
		}
		entries, err := utils.ParseSSHConfig(path)
		if err != nil {
			return err
		}

		var options []string
		for _, e := range entries {
			options = append(options, e.Name)
		}

		if len(options) == 0 {
			fmt.Println("No hosts found.")
			return nil
		}

		var selected []string
		prompt := &survey.MultiSelect{
			Message: "Select hosts to remove:",
			Options: options,
		}
		err = survey.AskOne(prompt, &selected)
		if err != nil {
			return err
		}

		if len(selected) == 0 {
			fmt.Println("No hosts selected.")
			return nil
		}

		if err := utils.ConfirmWithPIN(fmt.Sprintf("Are you sure you want to remove %d hosts?", len(selected))); err != nil {
			return err
		}

		for _, name := range selected {
			if err := utils.RemoveSSHConfigEntry(path, name); err != nil {
				fmt.Printf("Failed to remove %s: %s\n", name, err)
			}
		}
		fmt.Println("✅ Selected hosts removed.")
		return nil
	},
}

// SSHConfigEditCmd edits an existing entry
var SSHConfigEditCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit an SSH host config",
	RunE: func(cmd *cobra.Command, args []string) error {
		path, err := utils.GetSSHConfigPath()
		if err != nil {
			return err
		}
		entries, err := utils.ParseSSHConfig(path)
		if err != nil {
			return err
		}

		var options []string
		for _, e := range entries {
			options = append(options, e.Name)
		}

		if len(options) == 0 {
			fmt.Println("No hosts found.")
			return nil
		}

		var selectedHost string
		selectPrompt := &survey.Select{
			Message: "Select host to edit:",
			Options: options,
		}
		err = survey.AskOne(selectPrompt, &selectedHost)
		if err != nil {
			return err
		}

		// Find the entry
		var targetEntry utils.SSHConfigEntry
		for _, e := range entries {
			if e.Name == selectedHost {
				targetEntry = e
				break
			}
		}

		answers := struct {
			HostName string
			User     string
			Port     string
		}{
			HostName: targetEntry.HostName,
			User:     targetEntry.User,
			Port:     targetEntry.Port,
		}

		qs := []*survey.Question{
			{
				Name: "HostName",
				Prompt: &survey.Input{
					Message: "Host Name:",
					Default: targetEntry.HostName,
				},
			},
			{
				Name: "User",
				Prompt: &survey.Input{
					Message: "User:",
					Default: targetEntry.User,
				},
			},
			{
				Name: "Port",
				Prompt: &survey.Input{
					Message: "Port:",
					Default: targetEntry.Port,
				},
			},
		}

		if err := survey.Ask(qs, &answers); err != nil {
			return err
		}

		idFile, err := promptIdentityFile(targetEntry.IdentityFile)
		if err != nil {
			return err
		}

		// Update fields temporarily
		targetEntry.HostName = answers.HostName
		targetEntry.User = answers.User
		targetEntry.Port = answers.Port
		targetEntry.IdentityFile = idFile

		targetEntry, err = promptAuthOptions(targetEntry)
		if err != nil {
			return err
		}

		// Update logic
		for i, e := range entries {
			if e.Name == selectedHost {
				entries[i] = targetEntry
				entries[i].RawContent = nil
				break
			}
		}

		err = utils.WriteSSHConfig(path, entries)
		if err != nil {
			return err
		}
		fmt.Println("✅ Host updated.")
		return nil
	},
}

// SSHConfigExportCmd exports config to JSON
var SSHConfigExportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export SSH config entries to JSON",
	RunE: func(cmd *cobra.Command, args []string) error {
		path, err := utils.GetSSHConfigPath()
		if err != nil {
			return err
		}
		entries, err := utils.ParseSSHConfig(path)
		if err != nil {
			return err
		}

		var options []string
		for _, e := range entries {
			options = append(options, e.Name)
		}

		if len(options) == 0 {
			fmt.Println("No hosts found.")
			return nil
		}

		var selected []string
		prompt := &survey.MultiSelect{
			Message: "Select hosts to export:",
			Options: options,
		}
		err = survey.AskOne(prompt, &selected)
		if err != nil {
			return err
		}

		var exportEntries []utils.SSHConfigEntry
		for _, name := range selected {
			for _, e := range entries {
				if e.Name == name {
					exportEntries = append(exportEntries, e)
					break
				}
			}
		}

		filename := ""
		for {
			filePrompt := &survey.Input{
				Message: "Enter filename (e.g. hosts.json):",
			}
			if err := survey.AskOne(filePrompt, &filename); err != nil {
				return err
			}

			if filename == "" {
				return fmt.Errorf("filename required")
			}

			if utils.FileExists(filename) {
				choice := ""
				p := &survey.Select{
					Message: fmt.Sprintf("File %s exists. Action:", filename),
					Options: []string{"Overwrite", "Rename", "Timestamp", "Cancel"},
				}
				if err := survey.AskOne(p, &choice); err != nil {
					return err
				}
				if choice == "Cancel" {
					return nil
				}
				if choice == "Rename" {
					continue
				}
				if choice == "Timestamp" {
					ext := filepath.Ext(filename)
					base := strings.TrimSuffix(filename, ext)
					filename = fmt.Sprintf("%s-%d%s", base, time.Now().Unix(), ext)
				}
				// Overwrite falls through
			}
			break
		}

		data, err := json.MarshalIndent(exportEntries, "", "  ")
		if err != nil {
			return err
		}

		err = os.WriteFile(filename, data, 0644)
		if err != nil {
			return err
		}
		fmt.Printf("✅ Exported %d hosts to %s\n", len(exportEntries), filename)
		return nil
	},
}

// SSHConfigImportCmd imports config from JSON
var SSHConfigImportCmd = &cobra.Command{
	Use:   "import",
	Short: "Import SSH config from JSON",
	RunE: func(cmd *cobra.Command, args []string) error {
		filename := ""
		filePrompt := &survey.Input{
			Message: "Enter filename to import:",
		}
		if err := survey.AskOne(filePrompt, &filename); err != nil {
			return err
		}

		data, err := os.ReadFile(filename)
		if err != nil {
			return err
		}

		var imports []utils.SSHConfigEntry
		err = json.Unmarshal(data, &imports)
		if err != nil {
			return fmt.Errorf("invalid JSON structure: %w", err)
		}

		fmt.Printf("Found %d entries to import.\n", len(imports))

		if err := utils.ConfirmWithPIN("Import these entries?"); err != nil {
			return err
		}

		path, err := utils.GetSSHConfigPath()
		if err != nil {
			return err
		}
		for _, entry := range imports {
			if err := utils.AddSSHConfigEntry(path, entry); err != nil {
				fmt.Printf("Failed to import %s: %s\n", entry.Name, err)
			}
		}
		fmt.Println("✅ Import completed.")
		return nil
	},
}

var SSHConfigPathCmd = &cobra.Command{
	Use:   "path",
	Short: "Show path to ssh config",
	RunE: func(cmd *cobra.Command, args []string) error {
		path, err := utils.GetSSHConfigPath()
		if err == nil {
			fmt.Println(path)
		}
		return err
	},
}

var SSHConfigConnectCmd = &cobra.Command{
	Use:   "connect",
	Short: "Connect to a host from config",
	RunE: func(cmd *cobra.Command, args []string) error {
		path, err := utils.GetSSHConfigPath()
		if err != nil {
			return err
		}
		entries, err := utils.ParseSSHConfig(path)
		if err != nil {
			return err
		}

		var options []string
		for _, e := range entries {
			options = append(options, e.Name)
		}

		if len(options) == 0 {
			fmt.Println("No hosts found.")
			return nil
		}

		var selectedHost string
		prompt := &survey.Select{
			Message: "Select host to connect to:",
			Options: options,
		}
		err = survey.AskOne(prompt, &selectedHost)
		if err != nil {
			return err
		}

		c := exec.Command("ssh", selectedHost)
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		c.Stdin = os.Stdin
		return c.Run()
	},
}

var SSHConfigValidateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate SSH config syntax",
	RunE: func(cmd *cobra.Command, args []string) error {
		path, err := utils.GetSSHConfigPath()
		if err != nil {
			return err
		}
		_, err = utils.ParseSSHConfig(path)
		if err != nil {
			fmt.Println("❌ Config is invalid:", err)
			return err
		}
		fmt.Println("✅ Config syntax seems valid.")
		return nil
	},
}

// SSHConnectCmd represents 'mayrlabs ssh connect' (direct)
var SSHConnectCmd = &cobra.Command{
	Use:   "connect",
	Short: "Connect to a host directly",
	RunE: func(cmd *cobra.Command, args []string) error {
		host := ""
		user := ""
		port := "" // optional

		if err := survey.AskOne(&survey.Input{Message: "Host/IP:"}, &host); err != nil {
			return err
		}
		if err := survey.AskOne(&survey.Input{Message: "User (optional):"}, &user); err != nil {
			return err
		}
		if err := survey.AskOne(&survey.Input{Message: "Port (optional):"}, &port); err != nil {
			return err
		}

		idFile, err := promptIdentityFile("")
		if err != nil {
			return err
		}

		sshArgs := []string{}
		if port != "" {
			sshArgs = append(sshArgs, "-p", port)
		}
		if idFile != "" {
			sshArgs = append(sshArgs, "-i", idFile)
		}
		target := host
		if user != "" {
			target = user + "@" + host
		}
		sshArgs = append(sshArgs, target)

		c := exec.Command("ssh", sshArgs...)
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		c.Stdin = os.Stdin
		return c.Run()
	},
}

// SSHHostsCmd
var SSHHostsCmd = &cobra.Command{
	Use:   "hosts",
	Short: "Manage known_hosts",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Also interactive if no usage
		if len(args) == 0 {
			return runSSHHostsInteractive()
		}
		return cmd.Help()
	},
}

var SSHHostsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List known hosts",
	RunE: func(cmd *cobra.Command, args []string) error {
		path, err := utils.GetKnownHostsPath()
		if err != nil {
			return err
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		lines := strings.Split(string(content), "\n")
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		if _, err := fmt.Fprintln(w, "HOST/HASH\tKEY TYPE"); err != nil {
			return err
		}
		for _, line := range lines {
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				// parts[0] is host, parts[1] is type
				host := parts[0]
				kType := parts[1]
				// Marker handling (@revoked etc) not handled deeply here
				if _, err := fmt.Fprintf(w, "%s\t%s\n", host, kType); err != nil {
					return err
				}
			}
		}
		return w.Flush()
	},
}

var SSHHostsRemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove a known host",
	RunE: func(cmd *cobra.Command, args []string) error {
		path, err := utils.GetKnownHostsPath()
		if err != nil {
			return err
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		lines := strings.Split(string(content), "\n")

		var options []string
		for _, line := range lines {
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			parts := strings.Fields(line)
			if len(parts) >= 1 {
				options = append(options, parts[0])
			}
		}

		if len(options) == 0 {
			fmt.Println("No known hosts found.")
			return nil
		}

		var selected []string
		if err := survey.AskOne(&survey.MultiSelect{
			Message: "Select hosts to remove:",
			Options: options,
		}, &selected); err != nil {
			return err
		}

		if len(selected) == 0 {
			return nil
		}

		if err := utils.ConfirmWithPIN(fmt.Sprintf("Remove %d hosts?", len(selected))); err != nil {
			return err
		}

		for _, host := range selected {
			// ssh-keygen -R host
			if err := utils.RemoveKnownHost(host); err != nil {
				fmt.Printf("Failed to remove %s: %s\n", host, err)
			}
		}
		fmt.Println("✅ Hosts removed.")
		return nil
	},
}

// SSHKeysCmd
var SSHKeysCmd = &cobra.Command{
	Use:   "keys",
	Short: "Manage SSH keys",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Also interactive if no usage
		if len(args) == 0 {
			return runSSHKeysInteractive()
		}
		return cmd.Help()
	},
}

var SSHKeysListCmd = &cobra.Command{
	Use:   "list",
	Short: "List SSH keys",
	RunE: func(cmd *cobra.Command, args []string) error {
		keys, err := utils.ListSSHKeys()
		if err != nil {
			return err
		}

		var options []string
		for _, k := range keys {
			options = append(options, fmt.Sprintf("%s (%s)", k.Name, k.PubPath))
		}

		if len(options) == 0 {
			fmt.Println("No SSH keys found in ~/.ssh")
			return nil
		}

		selectedIdx := -1
		prompt := &survey.Select{
			Message: "Select an SSH key:",
			Options: options,
		}
		err = survey.AskOne(prompt, &selectedIdx)
		if err != nil {
			return err
		}
		selectedKey := keys[selectedIdx]

		action := ""
		actionPrompt := &survey.Select{
			Message: "Choose action:",
			Options: []string{"Copy Public Key", "Delete", "Change Passphrase", "Rename", "Cancel"},
		}
		if err := survey.AskOne(actionPrompt, &action); err != nil {
			return err
		}

		switch action {
		case "Copy Public Key":
			content, err := os.ReadFile(selectedKey.PubPath)
			if err != nil {
				return err
			}
			err = clipboard.WriteAll(string(content))
			if err != nil {
				return err
			}
			fmt.Println("✅ Public key copied to clipboard!")
		case "Delete":
			if err := utils.ConfirmWithPIN("Delete " + selectedKey.Name + " (private and public)?"); err != nil {
				return err
			}
			if err := os.Remove(selectedKey.Path); err != nil {
				return fmt.Errorf("failed to remove key: %w", err)
			}
			if err := os.Remove(selectedKey.PubPath); err != nil {
				return fmt.Errorf("failed to remove pub key: %w", err)
			}
			fmt.Println("✅ Key deleted.")
		case "Change Passphrase":
			c := exec.Command("ssh-keygen", "-p", "-f", selectedKey.Path)
			c.Stdout = os.Stdout
			c.Stderr = os.Stderr
			c.Stdin = os.Stdin
			if err := c.Run(); err != nil {
				return err
			}
		case "Rename":
			newName := ""
			if err := survey.AskOne(&survey.Input{Message: "New name:"}, &newName); err != nil {
				return err
			}
			if newName != "" {
				sshDir := filepath.Dir(selectedKey.Path)
				newPath := filepath.Join(sshDir, newName)
				newPub := newPath + ".pub"
				if err := os.Rename(selectedKey.Path, newPath); err != nil {
					return fmt.Errorf("failed to rename key: %w", err)
				}
				if err := os.Rename(selectedKey.PubPath, newPub); err != nil {
					return fmt.Errorf("failed to rename pub key: %w", err)
				}
				fmt.Println("✅ Key renamed.")
			}
		}
		return nil
	},
}

var SSHKeysCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new SSH key",
	RunE: func(cmd *cobra.Command, args []string) error {
		identifier := ""
		if err := survey.AskOne(&survey.Input{Message: "Identifier (e.g. github_work):"}, &identifier); err != nil {
			return err
		}
		if identifier == "" {
			return fmt.Errorf("identifier required")
		}

		home, _ := os.UserHomeDir()
		filename := "id_" + identifier
		path := filepath.Join(home, ".ssh", filename)

		// Check existence immediately
		if _, err := os.Stat(path); err == nil {
			return fmt.Errorf("file %s already exists", path)
		}

		typeChoice := ""
		if err := survey.AskOne(&survey.Select{
			Message: "Key Type:",
			Options: []string{"ed25519", "rsa", "ecdsa"},
			Default: "ed25519",
		}, &typeChoice); err != nil {
			return err
		}

		comment := ""
		if err := survey.AskOne(&survey.Input{Message: "Comment (optional):"}, &comment); err != nil {
			return err
		}

		passphrase := ""
		if err := survey.AskOne(&survey.Password{Message: "Passphrase (optional):"}, &passphrase); err != nil {
			return err
		}

		roundsStr := ""
		rounds := 0
		if err := survey.AskOne(&survey.Input{Message: "KDF Rounds (optional, default 64 or 16 for ed25519):", Default: "64"}, &roundsStr); err != nil {
			return err
		}
		if r, err := strconv.Atoi(roundsStr); err == nil {
			rounds = r
		}

		fmt.Println("Generating key...")
		err := utils.GenerateSSHKey(path, typeChoice, comment, passphrase, rounds)
		if err != nil {
			return err
		}
		fmt.Printf("✅ Key generated at %s\n", path)

		// Post-action options
		nextAction := ""
		if err := survey.AskOne(&survey.Select{
			Message: "What would you like to do next?",
			Options: []string{"Copy Public Key", "Show Public Key", "Run (install) Command", "Exit"},
		}, &nextAction); err != nil {
			return err
		}

		pubPath := path + ".pub"
		switch nextAction {
		case "Copy Public Key":
			content, err := os.ReadFile(pubPath)
			if err != nil {
				return err
			}
			if err := clipboard.WriteAll(string(content)); err != nil {
				return err
			}
			fmt.Println("✅ Public key copied.")
		case "Show Public Key":
			content, err := os.ReadFile(pubPath)
			if err != nil {
				return err
			}
			fmt.Println(string(content))
		case "Run (install) Command":
			// ask for host
			host := ""
			if err := survey.AskOne(&survey.Input{Message: "Host to install key to:"}, &host); err != nil {
				return err
			}
			if host != "" {
				c := exec.Command("ssh-copy-id", "-i", pubPath, host)
				c.Stdout = os.Stdout
				c.Stderr = os.Stderr
				c.Stdin = os.Stdin
				return c.Run()
			}
		}

		return nil
	},
}

func init() {
	SSHCmd.AddCommand(SSHConfigCmd)
	SSHConfigCmd.AddCommand(SSHConfigListCmd)
	SSHConfigCmd.AddCommand(SSHConfigAddCmd)
	SSHConfigCmd.AddCommand(SSHConfigRemoveCmd)
	SSHConfigCmd.AddCommand(SSHConfigEditCmd)
	SSHConfigCmd.AddCommand(SSHConfigExportCmd)
	SSHConfigCmd.AddCommand(SSHConfigImportCmd)
	SSHConfigCmd.AddCommand(SSHConfigPathCmd)
	SSHConfigCmd.AddCommand(SSHConfigConnectCmd)
	SSHConfigCmd.AddCommand(SSHConfigValidateCmd)

	SSHCmd.AddCommand(SSHConnectCmd)

	SSHCmd.AddCommand(SSHHostsCmd)
	SSHHostsCmd.AddCommand(SSHHostsListCmd)
	SSHHostsCmd.AddCommand(SSHHostsRemoveCmd)

	SSHCmd.AddCommand(SSHKeysCmd)
	SSHKeysCmd.AddCommand(SSHKeysListCmd)
	SSHKeysCmd.AddCommand(SSHKeysCreateCmd)
}
