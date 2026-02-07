package commands

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/AlecAivazis/survey/v2"
	"github.com/MayR-Labs/mayrlabs-go/internal/utils"
	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"
)

// SSHCmd represents the base ssh command
var SSHCmd = &cobra.Command{
	Use:   "ssh",
	Short: "Manage SSH configurations, keys, and connections",
	Long:  "A comprehensive tool for managing your SSH config, keys, known_hosts and connections.",
}

// SSHConfigCmd represents the ssh config command category
var SSHConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage SSH config file (~/.ssh/config)",
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
		fmt.Fprintln(w, "NAME\tHOSTNAME\tUSER\tPORT\tIDENTITY FILE")
		for _, e := range entries {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", e.Name, e.HostName, e.User, e.Port, e.IdentityFile)
		}
		return w.Flush()
	},
}

// SSHConfigAddCmd adds a new entry to ~/.ssh/config
var SSHConfigAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new SSH host config",
	RunE: func(cmd *cobra.Command, args []string) error {
		answers := struct {
			Name         string
			HostName     string
			User         string
			Port         string
			IdentityFile string
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
			{
				Name: "IdentityFile",
				Prompt: &survey.Input{
					Message: "Identity File (optional path to private key):",
				},
			},
		}

		err := survey.Ask(qs, &answers)
		if err != nil {
			return err
		}

		fmt.Println("\nNew Entry Preview:")
		fmt.Printf("Host %s\n", answers.Name)
		fmt.Printf("  HostName %s\n", answers.HostName)
		if answers.User != "" {
			fmt.Printf("  User %s\n", answers.User)
		}
		if answers.Port != "" {
			fmt.Printf("  Port %s\n", answers.Port)
		}
		if answers.IdentityFile != "" {
			fmt.Printf("  IdentityFile %s\n", answers.IdentityFile)
		}

		confirm := false
		prompt := &survey.Confirm{
			Message: "Add this entry to config?",
		}
		survey.AskOne(prompt, &confirm)

		if confirm {
			path, err := utils.GetSSHConfigPath()
			if err != nil {
				return err
			}
			entry := utils.SSHConfigEntry{
				Name:         answers.Name,
				HostName:     answers.HostName,
				User:         answers.User,
				Port:         answers.Port,
				IdentityFile: answers.IdentityFile,
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

		confirm := false
		confirmPrompt := &survey.Confirm{
			Message: fmt.Sprintf("Are you sure you want to remove %d hosts?", len(selected)),
		}
		survey.AskOne(confirmPrompt, &confirm)

		if confirm {
			for _, name := range selected {
				if err := utils.RemoveSSHConfigEntry(path, name); err != nil {
					fmt.Printf("Failed to remove %s: %s\n", name, err)
				}
			}
			fmt.Println("✅ Selected hosts removed.")
		}
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
			HostName     string
			User         string
			Port         string
			IdentityFile string
		}{
			HostName:     targetEntry.HostName,
			User:         targetEntry.User,
			Port:         targetEntry.Port,
			IdentityFile: targetEntry.IdentityFile,
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
			{
				Name: "IdentityFile",
				Prompt: &survey.Input{
					Message: "Identity File:",
					Default: targetEntry.IdentityFile,
				},
			},
		}

		err = survey.Ask(qs, &answers)
		if err != nil {
			return err
		}

		// Update logic: Remove old, add new (to keep it simple for now, though rewriting in place is better)
		// We'll update properties in memory and rewrite the whole file
		for i, e := range entries {
			if e.Name == selectedHost {
				entries[i].HostName = answers.HostName
				entries[i].User = answers.User
				entries[i].Port = answers.Port
				entries[i].IdentityFile = answers.IdentityFile
				// Clear raw content to force regeneration
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
		filePrompt := &survey.Input{
			Message: "Enter filename (e.g. hosts.json):",
		}
		survey.AskOne(filePrompt, &filename)

		if filename == "" {
			return fmt.Errorf("filename required")
		}

		data, err := json.MarshalIndent(exportEntries, "", "  ")
		if err != nil {
			return err
		}

		err = ioutil.WriteFile(filename, data, 0644)
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
		survey.AskOne(filePrompt, &filename)

		data, err := ioutil.ReadFile(filename)
		if err != nil {
			return err
		}

		var imports []utils.SSHConfigEntry
		err = json.Unmarshal(data, &imports)
		if err != nil {
			return fmt.Errorf("invalid JSON structure: %w", err)
		}

		fmt.Printf("Found %d entries to import.\n", len(imports))

		confirm := false
		prompt := &survey.Confirm{
			Message: "Import these entries?",
		}
		survey.AskOne(prompt, &confirm)

		if confirm {
			path, err := utils.GetSSHConfigPath()
			if err != nil {
				return err
			}
			for _, entry := range imports {
				// Naive check for duplicates could be added, but for now we append
				// Better: check if name exists
				// We already have AddSSHConfigEntry which appends
				// Ideally we shouldn't add duplicates
				if err := utils.AddSSHConfigEntry(path, entry); err != nil {
					fmt.Printf("Failed to import %s: %s\n", entry.Name, err)
				}
			}
			fmt.Println("✅ Import completed.")
		}
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

		survey.AskOne(&survey.Input{Message: "Host/IP:"}, &host)
		survey.AskOne(&survey.Input{Message: "User (optional):"}, &user)
		survey.AskOne(&survey.Input{Message: "Port (optional):"}, &port)

		sshArgs := []string{}
		if port != "" {
			sshArgs = append(sshArgs, "-p", port)
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
}

var SSHHostsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List known hosts",
	RunE: func(cmd *cobra.Command, args []string) error {
		// This is tricky because known_hosts can be hashed.
		// We can just dump the content or try to parse basic lines.
		// Requirement says "Show the host(s)".
		// Realistically, we can cat the file or use 'ssh-keygen -F' if we knew the host.
		// Detailed parsing is complex. I'll read line by line and show the first part if not hashed.
		path, err := utils.GetKnownHostsPath()
		if err != nil {
			return err
		}
		content, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		lines := strings.Split(string(content), "\n")
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "HOST/HASH\tKEY TYPE")
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
				fmt.Fprintf(w, "%s\t%s\n", host, kType)
			}
		}
		return w.Flush()
	},
}

var SSHHostsRemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove a known host",
	RunE: func(cmd *cobra.Command, args []string) error {
		hostName := ""
		survey.AskOne(&survey.Input{Message: "Enter hostname to remove:"}, &hostName)
		if hostName == "" {
			return nil
		}

		confirm := false
		survey.AskOne(&survey.Confirm{Message: "Remove " + hostName + " from known_hosts?"}, &confirm)
		if confirm {
			return utils.RemoveKnownHost(hostName)
		}
		return nil
	},
}

// SSHKeysCmd
var SSHKeysCmd = &cobra.Command{
	Use:   "keys",
	Short: "Manage SSH keys",
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
		survey.AskOne(actionPrompt, &action)

		switch action {
		case "Copy Public Key":
			content, err := ioutil.ReadFile(selectedKey.PubPath)
			if err != nil {
				return err
			}
			err = clipboard.WriteAll(string(content))
			if err != nil {
				return err
			}
			fmt.Println("✅ Public key copied to clipboard!")
		case "Delete":
			confirm := false
			survey.AskOne(&survey.Confirm{Message: "Delete " + selectedKey.Name + " (private and public)?"}, &confirm)
			if confirm {
				os.Remove(selectedKey.Path)
				os.Remove(selectedKey.PubPath)
				fmt.Println("✅ Key deleted.")
			}
		case "Change Passphrase":
			c := exec.Command("ssh-keygen", "-p", "-f", selectedKey.Path)
			c.Stdout = os.Stdout
			c.Stderr = os.Stderr
			c.Stdin = os.Stdin
			c.Run()
		case "Rename":
			newName := ""
			survey.AskOne(&survey.Input{Message: "New name:"}, &newName)
			if newName != "" {
				sshDir := filepath.Dir(selectedKey.Path)
				newPath := filepath.Join(sshDir, newName)
				newPub := newPath + ".pub"
				os.Rename(selectedKey.Path, newPath)
				os.Rename(selectedKey.PubPath, newPub)
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
		survey.AskOne(&survey.Input{Message: "Identifier (e.g. github_work):"}, &identifier)
		if identifier == "" {
			return fmt.Errorf("identifier required")
		}

		typeChoice := ""
		survey.AskOne(&survey.Select{
			Message: "Key Type:",
			Options: []string{"ed25519", "rsa", "ecdsa"},
			Default: "ed25519",
		}, &typeChoice)

		comment := ""
		survey.AskOne(&survey.Input{Message: "Comment (optional):"}, &comment)

		passphrase := ""
		survey.AskOne(&survey.Password{Message: "Passphrase (optional):"}, &passphrase)

		roundsStr := ""
		rounds := 0
		survey.AskOne(&survey.Input{Message: "KDF Rounds (optional, default 64 or 16 for ed25519):", Default: "64"}, &roundsStr)
		if r, err := strconv.Atoi(roundsStr); err == nil {
			rounds = r
		}

		home, _ := os.UserHomeDir()
		filename := "id_" + identifier
		path := filepath.Join(home, ".ssh", filename)

		// Check existence
		if _, err := os.Stat(path); err == nil {
			return fmt.Errorf("file %s already exists", path)
		}

		fmt.Println("Generating key...")
		err := utils.GenerateSSHKey(path, typeChoice, comment, passphrase, rounds)
		if err != nil {
			return err
		}
		fmt.Printf("✅ Key generated at %s\n", path)
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
