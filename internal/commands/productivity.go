package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/MayR-Labs/mayrlabs-go/internal/utils"
	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
)

// TimerCmd starts a countdown timer
var TimerCmd = &cobra.Command{
	Use:   "timer [duration]",
	Short: "Start a countdown timer (e.g., 15s, 1.5m, 1h)",
	Long:  "Start a countdown timer. Supported units: s (seconds), m (minutes), h (hours).",
	RunE: func(cmd *cobra.Command, args []string) error {
		var durationStr string
		var err error

		if len(args) > 0 {
			durationStr = args[0]
		} else {
			durationStr, err = utils.PromptInput("Enter duration (e.g., 5m): ")
			if err != nil {
				return err
			}
		}

		if durationStr == "" {
			return fmt.Errorf("duration cannot be empty")
		}

		duration, err := time.ParseDuration(durationStr)
		if err != nil {
			return fmt.Errorf("invalid duration format: %w", err)
		}

		fmt.Printf("⏱️  Timer started for %s\n", duration)

		// Use spinner for visual feedback
		s := spinner.New(spinner.CharSets[39], 100*time.Millisecond)
		s.Prefix = " "
		s.Start()

		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		endTime := time.Now().Add(duration)

		for range ticker.C {
			remaining := time.Until(endTime)
			if remaining <= 0 {
				break
			}

			// Format remaining time
			h := int(remaining.Hours())
			m := int(remaining.Minutes()) % 60
			s.Suffix = fmt.Sprintf(" Remaining: %02d:%02d:%02d", h, m, int(remaining.Seconds())%60)
		}

		s.Stop()
		fmt.Println("\n🔔 Time's up!")

		// Bell sound
		fmt.Print("\a")
		return nil
	},
}

// LOCCmd counts lines of code
var LOCCmd = &cobra.Command{
	Use:   "loc [path]",
	Short: "Count lines of code",
	Long:  "Count lines of code in the specified directory, grouped by language.",
	RunE: func(cmd *cobra.Command, args []string) error {
		root := "."
		if len(args) > 0 {
			root = args[0]
		}

		fmt.Printf("📊 Counting lines of code in %s...\n", root)

		stats := make(map[string]int)
		totalLines := 0
		totalFiles := 0

		err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				// Skip common ignore dirs
				name := info.Name()
				if name == ".git" || name == "node_modules" || name == "vendor" || name == ".idea" || name == ".vscode" {
					return filepath.SkipDir
				}
				return nil
			}

			// Check extension
			ext := strings.ToLower(filepath.Ext(path))
			if ext == "" {
				return nil
			}

			// Skip binary/image files (simple check)
			if isBinaryExt(ext) {
				return nil
			}

			content, err := os.ReadFile(path)
			if err != nil {
				return nil
			} // Skip unreadable files

			lines := strings.Count(string(content), "\n") + 1
			stats[ext] += lines
			totalLines += lines
			totalFiles++

			return nil
		})

		if err != nil {
			return fmt.Errorf("failed to walk directory: %w", err)
		}

		// Sort by lines desc
		type kv struct {
			Key   string
			Value int
		}
		var ss []kv
		for k, v := range stats {
			ss = append(ss, kv{k, v})
		}
		sort.Slice(ss, func(i, j int) bool {
			return ss[i].Value > ss[j].Value
		})

		fmt.Println("\nLanguage Distribution:")
		fmt.Println("----------------------")
		for _, kv := range ss {
			fmt.Printf("%-10s: %d lines\n", kv.Key, kv.Value)
		}
		fmt.Println("----------------------")
		fmt.Printf("Total: %d lines in %d files\n", totalLines, totalFiles)

		return nil
	},
}

func isBinaryExt(ext string) bool {
	binaryExts := map[string]bool{
		".png": true, ".jpg": true, ".jpeg": true, ".gif": true, ".ico": true,
		".pdf": true, ".zip": true, ".tar": true, ".gz": true, ".exe": true,
		".dll": true, ".so": true, ".dylib": true, ".bin": true, ".mp3": true,
		".mp4": true, ".mov": true, ".avi": true, ".woff": true, ".woff2": true,
		".ttf": true, ".eot": true, ".svg": false, // SVG is text
	}
	return binaryExts[ext]
}
