package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"

	"github.com/MayR-Labs/mayrlabs-go/internal/utils"
	"github.com/spf13/cobra"
)

const aliasFileName = ".mayrlabs_aliases"

// AliasCmd is the parent command for alias operations
var AliasCmd = &cobra.Command{
	Use:   "alias",
	Short: "Manage mayrlabs aliases",
	Long:  "Manage aliases for mayrlabs commands in a dedicated file",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Interactive mode - prompt user to choose subcommand
		options := []string{
			"list",
			"add",
			"remove",
			"edit",
			"rename",
			"disable",
			"enable",
			"add-popular",
			"register",
			"unregister",
		}

		choice, err := utils.PromptSelect("What would you like to do?", options)
		if err != nil {
			return err
		}

		switch choice {
		case "list":
			return AliasListSubCmd.RunE(cmd, []string{})
		case "add":
			return AliasAddCmd.RunE(cmd, []string{})
		case "remove":
			return AliasRemoveCmd.RunE(cmd, []string{})
		case "edit":
			return AliasEditCmd.RunE(cmd, []string{})
		case "rename":
			return AliasRenameCmd.RunE(cmd, []string{})
		case "disable":
			return AliasDisableCmd.RunE(cmd, []string{})
		case "enable":
			return AliasEnableCmd.RunE(cmd, []string{})
		case "add-popular":
			return AliasAddPopularCmd.RunE(cmd, []string{})
		case "register":
			return AliasRegisterCmd.RunE(cmd, []string{})
		case "unregister":
			return AliasUnregisterCmd.RunE(cmd, []string{})
		default:
			return fmt.Errorf("invalid choice")
		}
	},
}

// AliasRegisterCmd registers the alias file in shell config
var AliasRegisterCmd = &cobra.Command{
	Use:   "register [bash|zsh|all]",
	Short: "Register alias file in shell configuration",
	Long:  "Add source command to .bashrc, .zshrc or both",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if runtime.GOOS == "windows" {
			fmt.Println("⚠️  Aliases are unfortunately not supported on Windows.")
			return nil
		}

		target := "all"
		if len(args) > 0 {
			target = args[0]
		}

		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get home directory: %w", err)
		}

		sourceLine := fmt.Sprintf("[ -f ~/%s ] && source ~/%s", aliasFileName, aliasFileName)

		filesToUpdate := []string{}
		if target == "all" || target == "zsh" {
			filesToUpdate = append(filesToUpdate, filepath.Join(homeDir, ".zshrc"))
		}
		if target == "all" || target == "bash" {
			filesToUpdate = append(filesToUpdate, filepath.Join(homeDir, ".bashrc"))
			if runtime.GOOS == "darwin" {
				filesToUpdate = append(filesToUpdate, filepath.Join(homeDir, ".bash_profile"))
			}
		}

		for _, file := range filesToUpdate {
			if !utils.FileExists(file) {
				continue
			}

			content, err := utils.ReadFile(file)
			if err != nil {
				continue
			}

			if strings.Contains(content, aliasFileName) {
				fmt.Printf("✅ Already registered in %s\n", filepath.Base(file))
				continue
			}

			// Append to file
			f, err := os.OpenFile(file, os.O_APPEND|os.O_WRONLY, 0644)
			if err != nil {
				fmt.Printf("❌ Failed to open %s: %v\n", filepath.Base(file), err)
				continue
			}

			if _, err := f.WriteString("\n" + sourceLine + "\n"); err != nil {
				f.Close()
				fmt.Printf("❌ Failed to write to %s: %v\n", filepath.Base(file), err)
				continue
			}
			f.Close()
			fmt.Printf("✅ Registered in %s\n", filepath.Base(file))
		}

		return nil
	},
}

// AliasUnregisterCmd unregisters the alias file from shell config
var AliasUnregisterCmd = &cobra.Command{
	Use:   "unregister [bash|zsh|all]",
	Short: "Unregister alias file from shell configuration",
	Long:  "Remove source command from .bashrc, .zshrc or both",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if runtime.GOOS == "windows" {
			fmt.Println("⚠️  Aliases are unfortunately not supported on Windows.")
			return nil
		}

		target := "all"
		if len(args) > 0 {
			target = args[0]
		}

		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get home directory: %w", err)
		}

		filesToUpdate := []string{}
		if target == "all" || target == "zsh" {
			filesToUpdate = append(filesToUpdate, filepath.Join(homeDir, ".zshrc"))
		}
		if target == "all" || target == "bash" {
			filesToUpdate = append(filesToUpdate, filepath.Join(homeDir, ".bashrc"))
			if runtime.GOOS == "darwin" {
				filesToUpdate = append(filesToUpdate, filepath.Join(homeDir, ".bash_profile"))
			}
		}

		for _, file := range filesToUpdate {
			if !utils.FileExists(file) {
				continue
			}

			content, err := utils.ReadFile(file)
			if err != nil {
				continue
			}

			lines := strings.Split(content, "\n")
			newLines := []string{}
			removed := false

			for _, line := range lines {
				if strings.Contains(line, aliasFileName) {
					removed = true
					continue
				}
				newLines = append(newLines, line)
			}

			if removed {
				if err := os.WriteFile(file, []byte(strings.Join(newLines, "\n")), 0644); err != nil {
					fmt.Printf("❌ Failed to update %s: %v\n", filepath.Base(file), err)
				} else {
					fmt.Printf("✅ Unregistered from %s\n", filepath.Base(file))
				}
			} else {
				fmt.Printf("ℹ️  Not found in %s\n", filepath.Base(file))
			}
		}

		return nil
	},
}

// AliasListSubCmd lists all aliases
var AliasListSubCmd = &cobra.Command{
	Use:   "list",
	Short: "List all aliases",
	RunE: func(cmd *cobra.Command, args []string) error {
		aliases, err := loadAliases()
		if err != nil {
			return err
		}

		if len(aliases) == 0 {
			fmt.Println("No aliases found.")
			return nil
		}

		fmt.Println("📌 Current Aliases:")
		var currentGroup string
		for _, alias := range aliases {
			firstChar := strings.ToUpper(string(alias.Name[0]))
			if firstChar != currentGroup {
				currentGroup = firstChar
				fmt.Printf("\n-- %s --\n", currentGroup)
			}

			status := "✅"
			if alias.Disabled {
				status = "❌ (disabled)"
			}
			fmt.Printf("  %s: %s -> %s\n", status, alias.Name, alias.Command)
		}

		return nil
	},
}

// AliasAddCmd adds a new alias
var AliasAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new alias",
	RunE: func(cmd *cobra.Command, args []string) error {
		if runtime.GOOS == "windows" {
			fmt.Println("⚠️  Aliases are unfortunately not supported on Windows.")
			return nil
		}

		name, err := utils.SurveyInput("Enter alias name:", "")
		if err != nil {
			return err
		}
		if name == "" {
			return fmt.Errorf("alias name cannot be empty")
		}

		command, err := utils.SurveyInput("Enter command:", "")
		if err != nil {
			return err
		}
		if command == "" {
			return fmt.Errorf("command cannot be empty")
		}

		aliases, err := loadAliases()
		if err != nil {
			return err
		}

		// Check if exists
		for _, a := range aliases {
			if a.Name == name {
				return fmt.Errorf("alias '%s' already exists", name)
			}
		}

		aliases = append(aliases, Alias{Name: name, Command: command})
		if err := saveAliases(aliases); err != nil {
			return err
		}

		fmt.Printf("✅ Alias '%s' added successfully!\n", name)
		return nil
	},
}

// AliasRemoveCmd removes an alias
var AliasRemoveCmd = &cobra.Command{
	Use:   "remove [alias]",
	Short: "Remove an alias",
	RunE: func(cmd *cobra.Command, args []string) error {
		aliases, err := loadAliases()
		if err != nil {
			return err
		}

		if len(aliases) == 0 {
			fmt.Println("No aliases found.")
			return nil
		}

		var selectedAliases []string
		if len(args) > 0 {
			selectedAliases = []string{args[0]}
		} else {
			options := []string{}
			for _, a := range aliases {
				options = append(options, a.Name)
			}
			selectedAliases, err = utils.PromptMultiSelect("Select aliases to remove:", options)
			if err != nil {
				return err
			}
		}

		if len(selectedAliases) == 0 {
			return nil
		}

		if err := utils.ConfirmWithPIN(fmt.Sprintf("⚠️  Are you sure you want to remove %d alias(es)?", len(selectedAliases))); err != nil {
			return err
		}

		newAliases := []Alias{}
		for _, a := range aliases {
			keep := true
			for _, s := range selectedAliases {
				if a.Name == s {
					keep = false
					break
				}
			}
			if keep {
				newAliases = append(newAliases, a)
			}
		}

		if err := saveAliases(newAliases); err != nil {
			return err
		}

		fmt.Printf("✅ Removed %d alias(es)\n", len(selectedAliases))
		return nil
	},
}

// AliasEditCmd edits an alias
var AliasEditCmd = &cobra.Command{
	Use:   "edit [alias]",
	Short: "Edit an alias",
	RunE: func(cmd *cobra.Command, args []string) error {
		aliases, err := loadAliases()
		if err != nil {
			return err
		}

		if len(aliases) == 0 {
			fmt.Println("No aliases found.")
			return nil
		}

		var aliasToEdit string
		if len(args) > 0 {
			aliasToEdit = args[0]
		} else {
			options := []string{}
			for _, a := range aliases {
				options = append(options, a.Name)
			}
			aliasToEdit, err = utils.PromptSelect("Select alias to edit:", options)
			if err != nil {
				return err
			}
		}

		var targetAlias *Alias
		for i := range aliases {
			if aliases[i].Name == aliasToEdit {
				targetAlias = &aliases[i]
				break
			}
		}

		if targetAlias == nil {
			return fmt.Errorf("alias '%s' not found", aliasToEdit)
		}

		newCommand, err := utils.SurveyInput("Enter new command:", targetAlias.Command)
		if err != nil {
			return err
		}

		targetAlias.Command = newCommand
		if err := saveAliases(aliases); err != nil {
			return err
		}

		fmt.Printf("✅ Alias '%s' updated\n", aliasToEdit)
		return nil
	},
}

// AliasRenameCmd renames an alias
var AliasRenameCmd = &cobra.Command{
	Use:   "rename [alias]",
	Short: "Rename an alias",
	RunE: func(cmd *cobra.Command, args []string) error {
		aliases, err := loadAliases()
		if err != nil {
			return err
		}

		if len(aliases) == 0 {
			fmt.Println("No aliases found.")
			return nil
		}

		var aliasToRename string
		if len(args) > 0 {
			aliasToRename = args[0]
		} else {
			options := []string{}
			for _, a := range aliases {
				options = append(options, a.Name)
			}
			aliasToRename, err = utils.PromptSelect("Select alias to rename:", options)
			if err != nil {
				return err
			}
		}

		var targetAlias *Alias
		for i := range aliases {
			if aliases[i].Name == aliasToRename {
				targetAlias = &aliases[i]
				break
			}
		}

		if targetAlias == nil {
			return fmt.Errorf("alias '%s' not found", aliasToRename)
		}

		newName, err := utils.SurveyInput("Enter new name:", targetAlias.Name)
		if err != nil {
			return err
		}

		// Check if new name exists
		for _, a := range aliases {
			if a.Name == newName && a.Name != aliasToRename {
				return fmt.Errorf("alias '%s' already exists", newName)
			}
		}

		targetAlias.Name = newName
		if err := saveAliases(aliases); err != nil {
			return err
		}

		fmt.Printf("✅ Alias renamed to '%s'\n", newName)
		return nil
	},
}

// AliasDisableCmd disables an alias
var AliasDisableCmd = &cobra.Command{
	Use:   "disable [alias]",
	Short: "Disable an alias",
	RunE: func(cmd *cobra.Command, args []string) error {
		aliases, err := loadAliases()
		if err != nil {
			return err
		}

		if len(aliases) == 0 {
			fmt.Println("No aliases found.")
			return nil
		}

		var selectedAliases []string
		if len(args) > 0 {
			selectedAliases = []string{args[0]}
		} else {
			options := []string{}
			for _, a := range aliases {
				if !a.Disabled {
					options = append(options, a.Name)
				}
			}
			if len(options) == 0 {
				fmt.Println("No enabled aliases found.")
				return nil
			}
			selectedAliases, err = utils.PromptMultiSelect("Select aliases to disable:", options)
			if err != nil {
				return err
			}
		}

		count := 0
		for i := range aliases {
			for _, s := range selectedAliases {
				if aliases[i].Name == s {
					aliases[i].Disabled = true
					count++
				}
			}
		}

		if err := saveAliases(aliases); err != nil {
			return err
		}

		fmt.Printf("✅ Disabled %d alias(es)\n", count)
		return nil
	},
}

// AliasEnableCmd enables an alias
var AliasEnableCmd = &cobra.Command{
	Use:   "enable [alias]",
	Short: "Enable an alias",
	RunE: func(cmd *cobra.Command, args []string) error {
		aliases, err := loadAliases()
		if err != nil {
			return err
		}

		if len(aliases) == 0 {
			fmt.Println("No aliases found.")
			return nil
		}

		var selectedAliases []string
		if len(args) > 0 {
			selectedAliases = []string{args[0]}
		} else {
			options := []string{}
			for _, a := range aliases {
				if a.Disabled {
					options = append(options, a.Name)
				}
			}
			if len(options) == 0 {
				fmt.Println("No disabled aliases found.")
				return nil
			}
			selectedAliases, err = utils.PromptMultiSelect("Select aliases to enable:", options)
			if err != nil {
				return err
			}
		}

		count := 0
		for i := range aliases {
			for _, s := range selectedAliases {
				if aliases[i].Name == s {
					aliases[i].Disabled = false
					count++
				}
			}
		}

		if err := saveAliases(aliases); err != nil {
			return err
		}

		fmt.Printf("✅ Enabled %d alias(es)\n", count)
		return nil
	},
}

// AliasAddPopularCmd adds popular aliases
var AliasAddPopularCmd = &cobra.Command{
	Use:   "add-popular",
	Short: "Add popular aliases",
	RunE: func(cmd *cobra.Command, args []string) error {
		popularAliases := GetPopularAliases()

		options := []string{}
		for _, a := range popularAliases {
			options = append(options, fmt.Sprintf("%s -> %s", a.Name, a.Command))
		}

		selectedIndices, err := utils.PromptMultiSelect("Select popular aliases to add:", options)
		if err != nil {
			return err
		}

		if len(selectedIndices) == 0 {
			return nil
		}

		aliases, err := loadAliases()
		if err != nil {
			return err
		}

		count := 0
		for _, selected := range selectedIndices {
			// Extract name from selection string "name -> command"
			parts := strings.Split(selected, " -> ")
			name := parts[0]

			// Find command
			var command string
			for _, p := range popularAliases {
				if p.Name == name {
					command = p.Command
					break
				}
			}

			// Check if exists
			exists := false
			for _, a := range aliases {
				if a.Name == name {
					exists = true
					break
				}
			}

			if !exists {
				aliases = append(aliases, Alias{Name: name, Command: command})
				count++
			} else {
				fmt.Printf("⚠️  Alias '%s' already exists, skipping\n", name)
			}
		}

		if err := saveAliases(aliases); err != nil {
			return err
		}

		fmt.Printf("✅ Added %d popular alias(es)\n", count)
		return nil
	},
}

type Alias struct {
	Name     string
	Command  string
	Disabled bool
}

func loadAliases() ([]Alias, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	path := filepath.Join(homeDir, aliasFileName)

	if !utils.FileExists(path) {
		return []Alias{}, nil
	}

	content, err := utils.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var aliases []Alias
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		disabled := false
		if strings.HasPrefix(line, "#") {
			disabled = true
			line = strings.TrimPrefix(line, "#")
			line = strings.TrimSpace(line)
		}

		if !strings.HasPrefix(line, "alias ") {
			continue
		}

		// Parse alias name='command'
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		name := strings.TrimPrefix(parts[0], "alias ")
		command := strings.Trim(parts[1], "'\"")

		aliases = append(aliases, Alias{
			Name:     name,
			Command:  command,
			Disabled: disabled,
		})
	}

	return aliases, nil
}

func saveAliases(aliases []Alias) error {
	// Sort aliases
	sort.Slice(aliases, func(i, j int) bool {
		return strings.ToLower(aliases[i].Name) < strings.ToLower(aliases[j].Name)
	})

	var sb strings.Builder
	var currentGroup string

	for _, alias := range aliases {
		firstChar := strings.ToUpper(string(alias.Name[0]))
		if firstChar != currentGroup {
			if currentGroup != "" {
				sb.WriteString("\n")
			}
			currentGroup = firstChar
		}

		line := fmt.Sprintf("alias %s='%s'", alias.Name, alias.Command)
		if alias.Disabled {
			line = "# " + line
		}
		sb.WriteString(line + "\n")
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	path := filepath.Join(homeDir, aliasFileName)

	return os.WriteFile(path, []byte(sb.String()), 0644)
}

func init() {
	AliasCmd.AddCommand(AliasRegisterCmd)
	AliasCmd.AddCommand(AliasUnregisterCmd)
	AliasCmd.AddCommand(AliasListSubCmd)
	AliasCmd.AddCommand(AliasAddCmd)
	AliasCmd.AddCommand(AliasRemoveCmd)
	AliasCmd.AddCommand(AliasEditCmd)
	AliasCmd.AddCommand(AliasRenameCmd)
	AliasCmd.AddCommand(AliasDisableCmd)
	AliasCmd.AddCommand(AliasEnableCmd)
	AliasCmd.AddCommand(AliasAddPopularCmd)
}
