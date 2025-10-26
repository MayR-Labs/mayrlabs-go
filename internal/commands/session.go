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

		timestamp := time.Now().Format("20060102-150405")
		sessionFile := filepath.Join(configDir, fmt.Sprintf("session-%s.md", timestamp))

		// Initialize session file
		startTime := time.Now()
		header := fmt.Sprintf("# %s\n\n**Started:** %s\n\n---\n\n",
			summary, startTime.Format("2006-01-02 15:04:05"))

		if err := os.WriteFile(sessionFile, []byte(header), 0o644); err != nil {
			return fmt.Errorf("failed to create session file: %w", err)
		}

		fmt.Printf("✅ Session started: %s\n", summary)
		fmt.Printf("📝 Session file: %s\n\n", sessionFile)
		fmt.Println("Commands:")
		fmt.Println("  AI [question] - Ask AI a question and record the interaction")
		fmt.Println("  Note [note]   - Add a note to the session")
		fmt.Println("  END           - End the session")
		fmt.Println()

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

			switch command {
			case "END":
				// End session
				endTime := time.Now()
				duration := endTime.Sub(startTime)

				// Add end time to file
				endContent := fmt.Sprintf("\n---\n\n**Ended:** %s\n**Duration:** %s\n\n",
					endTime.Format("2006-01-02 15:04:05"), formatDuration(duration))

				// Add a nice closing message
				endContent += getSessionEndMessage(duration)

				file, err := os.OpenFile(sessionFile, os.O_APPEND|os.O_WRONLY, 0o644)
				if err != nil {
					return fmt.Errorf("failed to open session file: %w", err)
				}
				if _, err := file.WriteString(endContent); err != nil {
					file.Close()
					return fmt.Errorf("failed to write to session file: %w", err)
				}
				file.Close()

				// Move session file to pwd
				pwd, err := os.Getwd()
				if err != nil {
					return fmt.Errorf("failed to get working directory: %w", err)
				}

				destFile := filepath.Join(pwd, filepath.Base(sessionFile))
				if err := os.Rename(sessionFile, destFile); err != nil {
					// If rename fails, try copy
					content, readErr := os.ReadFile(sessionFile)
					if readErr != nil {
						return fmt.Errorf("failed to read session file: %w", readErr)
					}
					if writeErr := os.WriteFile(destFile, content, 0o644); writeErr != nil {
						return fmt.Errorf("failed to copy session file: %w", writeErr)
					}
					os.Remove(sessionFile)
				}

				fmt.Printf("\n✅ Session ended and saved to: %s\n", destFile)
				return nil

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
				resp, err := model.GenerateContent(ctx, genai.Text(question))
				client.Close()

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
				file.WriteString(entry)
				file.Close()

			case "NOTE":
				if len(parts) < 2 {
					fmt.Println("❌ Please provide a note")
					continue
				}

				note := parts[1]
				timestamp := time.Now().Format("15:04:05")
				entry := fmt.Sprintf("## Note [%s]\n\n%s\n\n---\n\n", timestamp, note)

				file, err := os.OpenFile(sessionFile, os.O_APPEND|os.O_WRONLY, 0o644)
				if err != nil {
					fmt.Printf("❌ Failed to record note: %v\n", err)
					continue
				}
				file.WriteString(entry)
				file.Close()

				fmt.Println("✅ Note recorded")

			default:
				fmt.Println("❌ Unknown command. Use: AI [question], Note [note], or END")
			}
		}
	},
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
