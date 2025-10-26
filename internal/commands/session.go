package commands

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/MayR-Labs/mayrlabs-go/internal/utils"
	"github.com/google/generative-ai-go/genai"
	"github.com/spf13/cobra"
	"google.golang.org/api/option"
)

// SessionStartCmd starts an interactive session
var SessionStartCmd = &cobra.Command{
	Use:   "session-start [summary]",
	Short: "Start an interactive development session",
	Long:  "Start an interactive session that records AI interactions and notes",
	RunE: func(cmd *cobra.Command, args []string) error {
		return startSession(args, false, "", "")
	},
}

// startSession is the core session logic, shared by regular and secure sessions
func startSession(args []string, secure bool, password string, passwordHint string) error {
	// Get summary
	var summary string
	var err error

	if len(args) == 0 {
		summary, err = utils.SurveyInput("Enter session summary:", "")
		if err != nil {
			return err
		}
	} else {
		summary = strings.Join(args, " ")
	}

	if summary == "" {
		summary = "Development Session"
	}

	// Create session file
	configDir, err := utils.GetConfigDir()
	if err != nil {
		return err
	}

	// Create sessions subfolder
	var sessionsDir string
	if secure {
		sessionsDir = filepath.Join(configDir, "secure-sessions")
	} else {
		sessionsDir = filepath.Join(configDir, "sessions")
	}
	if err := os.MkdirAll(sessionsDir, 0o700); err != nil {
		return fmt.Errorf("failed to create sessions directory: %w", err)
	}

	timestamp := time.Now().Format("20060102-150405")
	sessionFile := filepath.Join(configDir, fmt.Sprintf("session-%s.md", timestamp))

	// Initialize session file
	startTime := time.Now()
	header := fmt.Sprintf("# %s\n\n**Started:** %s\n\n---\n\n",
		summary, startTime.Format("2006-01-02 15:04:05"))

	if err := os.WriteFile(sessionFile, []byte(header), 0o644); err != nil {
		return fmt.Errorf("failed to create session file: %w", err)
	}

	sessionType := "Session"
	if secure {
		sessionType = "Secure session"
	}
	fmt.Printf("✅ %s started: %s\n", sessionType, summary)
	fmt.Printf("📝 Session file: %s\n\n", sessionFile)
	fmt.Println("Commands:")
	fmt.Println("  AI [question] - Ask AI a question and record the interaction")
	fmt.Println("  Task [task]   - Add a task/todo to the session")
	fmt.Println("  [any text]    - Add a note (no prefix needed)")
	fmt.Println("  END/QUIT/CLOSE - End the session")
	fmt.Println()

	// Track tasks during session
	var tasks []string

	// Interactive loop
	for {
		var input string
		prompt := &survey.Input{
			Message: "Session >",
		}
		if err := survey.AskOne(prompt, &input); err != nil {
			// User cancelled
			fmt.Println("\n❌ Session cancelled")
			return nil
		}

		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		// Parse command
		parts := strings.SplitN(input, " ", 2)
		command := strings.ToUpper(parts[0])

		// Check for END aliases
		if command == "END" || command == "QUIT" || command == "CLOSE" {
			// End session
			endTime := time.Now()
			duration := endTime.Sub(startTime)

			// Add tasks section if there are any tasks
			var tasksContent string
			if len(tasks) > 0 {
				tasksContent = "\n## Tasks\n\n"
				for _, task := range tasks {
					tasksContent += fmt.Sprintf("- [ ] %s\n", task)
				}
				tasksContent += "\n"
			}

			// Add end time to file
			endContent := fmt.Sprintf("%s\n---\n\n**Ended:** %s\n**Duration:** %s\n\n",
				tasksContent, endTime.Format("2006-01-02 15:04:05"), formatDuration(duration))

			// Add a nice closing message
			endContent += getSessionEndMessage(duration)

			file, err := os.OpenFile(sessionFile, os.O_APPEND|os.O_WRONLY, 0o644)
			if err != nil {
				return fmt.Errorf("failed to open session file: %w", err)
			}
			if _, err := file.WriteString(endContent); err != nil {
				_ = file.Close()
				return fmt.Errorf("failed to write to session file: %w", err)
			}
			if err = file.Close(); err != nil {
				return err
			}

			// Handle encryption for secure sessions
			if secure {
				if err := encryptSessionFile(sessionFile, password); err != nil {
					return fmt.Errorf("failed to encrypt session: %w", err)
				}

				// Update password hints file
				hintsFile := filepath.Join(sessionsDir, "password-hints.txt")
				hintEntry := fmt.Sprintf("%s: %s\n", filepath.Base(sessionFile), passwordHint)
				hintData, _ := os.ReadFile(hintsFile)
				hintData = append(hintData, []byte(hintEntry)...)
				if err := os.WriteFile(hintsFile, hintData, 0o600); err != nil {
					return fmt.Errorf("failed to save password hint: %w", err)
				}
			}

			// Copy session file to pwd
			pwd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get working directory: %w", err)
			}

			destFile := filepath.Join(pwd, filepath.Base(sessionFile))
			content, err := os.ReadFile(sessionFile)
			if err != nil {
				return fmt.Errorf("failed to read session file: %w", err)
			}
			if err := os.WriteFile(destFile, content, 0o644); err != nil {
				return fmt.Errorf("failed to copy session file to pwd: %w", err)
			}

			// Copy to sessions folder
			sessionDestFile := filepath.Join(sessionsDir, filepath.Base(sessionFile))
			if err := os.WriteFile(sessionDestFile, content, 0o644); err != nil {
				return fmt.Errorf("failed to copy session file to sessions folder: %w", err)
			}

			// Remove temporary session file from config dir
			_ = os.Remove(sessionFile)

			fmt.Printf("\n✅ Session ended and saved to:\n")
			fmt.Printf("   - Current directory: %s\n", destFile)
			fmt.Printf("   - Sessions folder: %s\n", sessionDestFile)
			return nil
		}

		switch command {
		case "AI":
			if len(parts) < 2 {
				fmt.Println("❌ Please provide a question for the AI")
				continue
			}

			question := parts[1]
			fmt.Println("\n🤖 Asking AI...")

			// Get API key
			apiKey, err := utils.GetAPIKey()
			if err != nil {
				fmt.Printf("❌ %v\n", err)
				fmt.Println("Please run 'mayrlabs ai-setup' first")
				continue
			}

			// Query Gemini
			ctx := context.Background()
			client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
			if err != nil {
				fmt.Printf("❌ Failed to create client: %v\n", err)
				continue
			}

			model := client.GenerativeModel("gemini-2.0-flash-exp")
			resp, _ := model.GenerateContent(ctx, genai.Text(question))
			err = client.Close()
			if err != nil {
				return err
			}

			if err != nil {
				fmt.Printf("❌ Failed to generate content: %v\n", err)
				continue
			}

			// Get response text
			var answer string
			if resp != nil && len(resp.Candidates) > 0 {
				for _, part := range resp.Candidates[0].Content.Parts {
					answer += fmt.Sprintf("%v", part)
				}
			}

			// Display response
			fmt.Println("\n📝 Response:")
			fmt.Println(answer)
			fmt.Println()

			// Record in session file
			entry := fmt.Sprintf("## AI\n\n**Question:** %s\n\n**Answer:** %s\n\n---\n\n",
				question, answer)

			file, err := os.OpenFile(sessionFile, os.O_APPEND|os.O_WRONLY, 0o644)
			if err != nil {
				fmt.Printf("❌ Failed to record interaction: %v\n", err)
				continue
			}
			_, err = file.WriteString(entry)
			if err != nil {
				_ = file.Close()
				continue
			}
			_ = file.Close()

		case "TASK":
			if len(parts) < 2 {
				fmt.Println("❌ Please provide a task")
				continue
			}

			task := parts[1]
			tasks = append(tasks, task)
			fmt.Println("✅ Task recorded")

		default:
			// Treat everything else as a note
			note := input
			timestamp := time.Now().Format("15:04:05")
			entry := fmt.Sprintf("## Note [%s]\n\n%s\n\n---\n\n", timestamp, note)

			file, err := os.OpenFile(sessionFile, os.O_APPEND|os.O_WRONLY, 0o644)
			if err != nil {
				fmt.Printf("❌ Failed to record note: %v\n", err)
				continue
			}
			_, err = file.WriteString(entry)
			if err != nil {
				_ = file.Close()
				continue
			}
			_ = file.Close()

			fmt.Println("✅ Note recorded")
		}
	}
}

// formatDuration formats a duration in a human-readable way
func formatDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	if hours > 0 {
		return fmt.Sprintf("%dh %dm %ds", hours, minutes, seconds)
	} else if minutes > 0 {
		return fmt.Sprintf("%dm %ds", minutes, seconds)
	}
	return fmt.Sprintf("%ds", seconds)
}

// getSessionEndMessage returns a motivational message based on session duration
func getSessionEndMessage(duration time.Duration) string {
	minutes := int(duration.Minutes())

	messages := map[string]string{
		"short":  "That was a quick session! Every bit of progress counts. 🚀",
		"medium": "That was a great session! Lasted for %s, hope it was productive. 💪",
		"long":   "Wow, that was an intense session! %s of focused work. Time for a break! ☕",
		"epic":   "Epic session! %s of pure dedication. You're crushing it! 🎉",
	}

	var key string
	switch {
	case minutes < 15:
		key = "short"
	case minutes < 60:
		key = "medium"
	case minutes < 180:
		key = "long"
	default:
		key = "epic"
	}

	message := messages[key]
	if strings.Contains(message, "%s") {
		return fmt.Sprintf(message, formatDuration(duration))
	}
	return message
}

// encryptSessionFile encrypts a session file with the given password
func encryptSessionFile(filename string, password string) error {
	// Read the file content
	content, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Ensure password is not empty
	if password == "" {
		return fmt.Errorf("password cannot be empty")
	}

	// Create a simple XOR encryption with password hash for basic security
	// Note: This is a basic implementation. For production, use proper encryption like AES.
	key := []byte(password)
	encrypted := make([]byte, len(content))
	for i := 0; i < len(content); i++ {
		encrypted[i] = content[i] ^ key[i%len(key)]
	}

	// Write encrypted content back
	if err := os.WriteFile(filename, encrypted, 0o600); err != nil {
		return fmt.Errorf("failed to write encrypted file: %w", err)
	}

	return nil
}

// decryptSessionFile decrypts a session file with the given password
func decryptSessionFile(filename string, password string) ([]byte, error) {
	// Read encrypted content
	encrypted, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Ensure password is not empty
	if password == "" {
		return nil, fmt.Errorf("password cannot be empty")
	}

// Decrypt using XOR with password
key := []byte(password)
decrypted := make([]byte, len(encrypted))
for i := 0; i < len(encrypted); i++ {
decrypted[i] = encrypted[i] ^ key[i%len(key)]
}

return decrypted, nil
}

// SecureSessionStartCmd starts a secure encrypted session
var SecureSessionStartCmd = &cobra.Command{
Use:   "secure-session-start [summary]",
Short: "Start an encrypted interactive development session",
Long:  "Start an encrypted session that records AI interactions and notes",
RunE: func(cmd *cobra.Command, args []string) error {
// Get password
password, err := utils.SurveyPassword("Enter password for session encryption:")
if err != nil {
return err
}

if password == "" {
return fmt.Errorf("password cannot be empty")
}

// Get password hint
passwordHint, err := utils.SurveyInput("Enter password hint:", "")
if err != nil {
return err
}

return startSession(args, true, password, passwordHint)
},
}

// SessionsCmd lists and manages sessions
var SessionsCmd = &cobra.Command{
Use:   "sessions",
Short: "List and manage sessions",
Long:  "List all sessions and perform actions like copy, delete, or view",
RunE: func(cmd *cobra.Command, args []string) error {
configDir, err := utils.GetConfigDir()
if err != nil {
return err
}

sessionsDir := filepath.Join(configDir, "sessions")
if err := os.MkdirAll(sessionsDir, 0o700); err != nil {
return fmt.Errorf("failed to create sessions directory: %w", err)
}

// List session files
files, err := os.ReadDir(sessionsDir)
if err != nil {
return fmt.Errorf("failed to read sessions directory: %w", err)
}

if len(files) == 0 {
fmt.Println("📁 No sessions found")
return nil
}

// Build options for selection
var options []string
for _, file := range files {
if !file.IsDir() && strings.HasSuffix(file.Name(), ".md") {
options = append(options, file.Name())
}
}

if len(options) == 0 {
fmt.Println("📁 No sessions found")
return nil
}

// Let user select a session
selectedSession, err := utils.PromptSelect("Select a session:", options)
if err != nil {
return err
}

sessionFile := filepath.Join(sessionsDir, selectedSession)

// Show action options
for {
action, err := utils.PromptSelect("Choose an action:", []string{
"Copy file to current directory",
"Delete session",
"Copy content to clipboard",
"Go back",
})
if err != nil {
return err
}

switch action {
case "Copy file to current directory":
pwd, err := os.Getwd()
if err != nil {
return fmt.Errorf("failed to get working directory: %w", err)
}

destFile := filepath.Join(pwd, selectedSession)
content, err := os.ReadFile(sessionFile)
if err != nil {
return fmt.Errorf("failed to read session file: %w", err)
}

if err := os.WriteFile(destFile, content, 0o644); err != nil {
return fmt.Errorf("failed to copy file: %w", err)
}

fmt.Printf("✅ Session copied to: %s\n", destFile)
return nil

case "Delete session":
confirm, err := utils.SurveyConfirm(fmt.Sprintf("Are you sure you want to delete %s?", selectedSession), false)
if err != nil {
return err
}

if confirm {
if err := os.Remove(sessionFile); err != nil {
return fmt.Errorf("failed to delete session: %w", err)
}
fmt.Println("✅ Session deleted")
}
return nil

case "Copy content to clipboard":
content, err := os.ReadFile(sessionFile)
if err != nil {
return fmt.Errorf("failed to read session file: %w", err)
}

if err := utils.CopyToClipboard(string(content)); err != nil {
return fmt.Errorf("failed to copy to clipboard: %w", err)
}

fmt.Println("✅ Session content copied to clipboard")
return nil

case "Go back":
return nil
}
}
},
}

// SessionClearCmd clears all sessions with PIN confirmation
var SessionClearCmd = &cobra.Command{
Use:   "session-clear",
Short: "Clear all sessions with PIN confirmation",
Long:  "Delete all sessions after 6-digit PIN confirmation",
RunE: func(cmd *cobra.Command, args []string) error {
configDir, err := utils.GetConfigDir()
if err != nil {
return err
}

sessionsDir := filepath.Join(configDir, "sessions")

// Check if sessions exist
files, err := os.ReadDir(sessionsDir)
if err != nil || len(files) == 0 {
fmt.Println("📁 No sessions to clear")
return nil
}

// Count session files
count := 0
for _, file := range files {
if !file.IsDir() && strings.HasSuffix(file.Name(), ".md") {
count++
}
}

if count == 0 {
fmt.Println("📁 No sessions to clear")
return nil
}

fmt.Printf("⚠️  This will delete %d session(s)\n", count)

// Generate PIN
pin, err := utils.GeneratePIN()
if err != nil {
return fmt.Errorf("failed to generate PIN: %w", err)
}

fmt.Printf("🔐 Enter PIN to confirm: %s\n", pin)

// Get user input
userPin, err := utils.SurveyInput("Enter PIN:", "")
if err != nil {
return err
}

if userPin != pin {
fmt.Println("❌ Invalid PIN. Operation cancelled")
return nil
}

// Delete all sessions
for _, file := range files {
if !file.IsDir() && strings.HasSuffix(file.Name(), ".md") {
sessionFile := filepath.Join(sessionsDir, file.Name())
if err := os.Remove(sessionFile); err != nil {
fmt.Printf("⚠️  Failed to delete %s: %v\n", file.Name(), err)
}
}
}

fmt.Printf("✅ Cleared %d session(s)\n", count)
return nil
},
}

// SessionPruneCmd prunes old sessions
var SessionPruneCmd = &cobra.Command{
Use:   "session-prune [days]",
Short: "Delete sessions older than specified days",
Long:  "Delete sessions older than the specified number of days with PIN confirmation",
Args:  cobra.ExactArgs(1),
RunE: func(cmd *cobra.Command, args []string) error {
configDir, err := utils.GetConfigDir()
if err != nil {
return err
}

// Parse days argument
days := 0
if _, err := fmt.Sscanf(args[0], "%d", &days); err != nil {
return fmt.Errorf("invalid days argument: %w", err)
}

if days < 1 {
return fmt.Errorf("days must be at least 1")
}

sessionsDir := filepath.Join(configDir, "sessions")

// Check if sessions exist
files, err := os.ReadDir(sessionsDir)
if err != nil || len(files) == 0 {
fmt.Println("📁 No sessions to prune")
return nil
}

// Find old sessions
cutoffTime := time.Now().AddDate(0, 0, -days)
var oldSessions []string

for _, file := range files {
if !file.IsDir() && strings.HasSuffix(file.Name(), ".md") {
info, err := file.Info()
if err != nil {
continue
}

if info.ModTime().Before(cutoffTime) {
oldSessions = append(oldSessions, file.Name())
}
}
}

if len(oldSessions) == 0 {
fmt.Printf("📁 No sessions older than %d days\n", days)
return nil
}

fmt.Printf("⚠️  Found %d session(s) older than %d days:\n", len(oldSessions), days)
for _, name := range oldSessions {
fmt.Printf("   - %s\n", name)
}
fmt.Println()

// Generate PIN
pin, err := utils.GeneratePIN()
if err != nil {
return fmt.Errorf("failed to generate PIN: %w", err)
}

fmt.Printf("🔐 Enter PIN to confirm: %s\n", pin)

// Get user input
userPin, err := utils.SurveyInput("Enter PIN:", "")
if err != nil {
return err
}

if userPin != pin {
fmt.Println("❌ Invalid PIN. Operation cancelled")
return nil
}

// Delete old sessions
for _, name := range oldSessions {
sessionFile := filepath.Join(sessionsDir, name)
if err := os.Remove(sessionFile); err != nil {
fmt.Printf("⚠️  Failed to delete %s: %v\n", name, err)
}
}

fmt.Printf("✅ Pruned %d session(s)\n", len(oldSessions))
return nil
},
}

// SecureSessionsCmd lists and manages secure sessions
var SecureSessionsCmd = &cobra.Command{
Use:   "secure-sessions",
Short: "List and manage encrypted secure sessions",
Long:  "List all secure sessions and perform actions with password authentication",
RunE: func(cmd *cobra.Command, args []string) error {
configDir, err := utils.GetConfigDir()
if err != nil {
return err
}

sessionsDir := filepath.Join(configDir, "secure-sessions")
if err := os.MkdirAll(sessionsDir, 0o700); err != nil {
return fmt.Errorf("failed to create secure-sessions directory: %w", err)
}

// List session files
files, err := os.ReadDir(sessionsDir)
if err != nil {
return fmt.Errorf("failed to read sessions directory: %w", err)
}

if len(files) == 0 {
fmt.Println("📁 No secure sessions found")
return nil
}

// Build options for selection
var options []string
for _, file := range files {
if !file.IsDir() && strings.HasSuffix(file.Name(), ".md") {
options = append(options, file.Name())
}
}

if len(options) == 0 {
fmt.Println("📁 No secure sessions found")
return nil
}

// Let user select a session
selectedSession, err := utils.PromptSelect("Select a secure session:", options)
if err != nil {
return err
}

sessionFile := filepath.Join(sessionsDir, selectedSession)

// Show password hint if available
hintsFile := filepath.Join(sessionsDir, "password-hints.txt")
if utils.FileExists(hintsFile) {
hintsContent, err := os.ReadFile(hintsFile)
if err == nil {
lines := strings.Split(string(hintsContent), "\n")
for _, line := range lines {
if strings.HasPrefix(line, selectedSession+":") {
hint := strings.TrimPrefix(line, selectedSession+":")
hint = strings.TrimSpace(hint)
if hint != "" {
fmt.Printf("💡 Password hint: %s\n", hint)
}
break
}
}
}
}

// Show action options
action, err := utils.PromptSelect("Choose an action:", []string{
"Copy file to current directory",
"Copy content to clipboard",
"Delete session",
"Go back",
})
if err != nil {
return err
}

if action == "Go back" {
return nil
}

if action == "Delete session" {
confirm, err := utils.SurveyConfirm(fmt.Sprintf("Are you sure you want to delete %s?", selectedSession), false)
if err != nil {
return err
}

if confirm {
if err := os.Remove(sessionFile); err != nil {
return fmt.Errorf("failed to delete session: %w", err)
}
fmt.Println("✅ Secure session deleted")
}
return nil
}

// For copy actions, we need the password
password, err := utils.SurveyPassword("Enter session password:")
if err != nil {
return err
}

// Try to decrypt
decrypted, err := decryptSessionFile(sessionFile, password)
if err != nil {
return fmt.Errorf("failed to decrypt session: %w", err)
}

// Verify decryption by checking if content looks like valid markdown
if !strings.Contains(string(decrypted), "# ") && !strings.Contains(string(decrypted), "**Started:**") {
return fmt.Errorf("❌ incorrect password or corrupted session file")
}

switch action {
case "Copy file to current directory":
pwd, err := os.Getwd()
if err != nil {
return fmt.Errorf("failed to get working directory: %w", err)
}

destFile := filepath.Join(pwd, selectedSession)
if err := os.WriteFile(destFile, decrypted, 0o644); err != nil {
return fmt.Errorf("failed to copy file: %w", err)
}

fmt.Printf("✅ Decrypted session copied to: %s\n", destFile)

case "Copy content to clipboard":
if err := utils.CopyToClipboard(string(decrypted)); err != nil {
return fmt.Errorf("failed to copy to clipboard: %w", err)
}

fmt.Println("✅ Decrypted session content copied to clipboard")
}

return nil
},
}
