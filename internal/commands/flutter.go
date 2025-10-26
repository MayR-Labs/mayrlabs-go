package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/MayR-Labs/mayrlabs-go/internal/utils"
)

// FlutterCmd is the parent command for Flutter operations
var FlutterCmd = &cobra.Command{
	Use:   "flutter",
	Short: "Flutter-related commands",
	Long:  "Commands for managing Flutter projects and build scripts",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

// FlutterCreateScriptsCmd creates useful build scripts
var FlutterCreateScriptsCmd = &cobra.Command{
	Use:   "create-scripts",
	Short: "Add useful build scripts to scripts/ (IPA, APK, AppBundle, etc.)",
	Long:  "Generate shell scripts for building Flutter apps for various platforms",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Create scripts directory
		if err := os.MkdirAll("scripts", 0o755); err != nil {
			return fmt.Errorf("failed to create scripts directory: %w", err)
		}

		scripts := map[string]string{
			"scripts/build-apk.sh": `#!/bin/bash
# Build Android APK (debug)
echo "Building Android APK..."
flutter build apk --debug
echo "✅ APK built successfully!"
echo "Location: build/app/outputs/flutter-apk/app-debug.apk"
`,
			"scripts/build-apk-release.sh": `#!/bin/bash
# Build Android APK (release)
echo "Building Android APK (Release)..."
flutter build apk --release
echo "✅ APK built successfully!"
echo "Location: build/app/outputs/flutter-apk/app-release.apk"
`,
			"scripts/build-appbundle.sh": `#!/bin/bash
# Build Android App Bundle
echo "Building Android App Bundle..."
flutter build appbundle --release
echo "✅ App Bundle built successfully!"
echo "Location: build/app/outputs/bundle/release/app-release.aab"
`,
			"scripts/build-ios.sh": `#!/bin/bash
# Build iOS app
echo "Building iOS app..."
flutter build ios --release --no-codesign
echo "✅ iOS app built successfully!"
echo "Note: Codesigning is disabled. Configure in Xcode for actual deployment."
`,
			"scripts/build-ipa.sh": `#!/bin/bash
# Build iOS IPA
echo "Building iOS IPA..."
flutter build ipa --release
echo "✅ IPA built successfully!"
echo "Location: build/ios/ipa/"
`,
			"scripts/clean.sh": `#!/bin/bash
# Clean Flutter build
echo "Cleaning Flutter build..."
flutter clean
flutter pub get
echo "✅ Clean completed!"
`,
			"scripts/run-tests.sh": `#!/bin/bash
# Run Flutter tests
echo "Running Flutter tests..."
flutter test
echo "✅ Tests completed!"
`,
		}

		for filename, content := range scripts {
			if err := utils.WriteFile(filename, content); err != nil {
				return fmt.Errorf("failed to write %s: %w", filename, err)
			}

			// Make scripts executable
			if err := os.Chmod(filename, 0o755); err != nil {
				return fmt.Errorf("failed to make %s executable: %w", filename, err)
			}

			fmt.Printf("✅ Created %s\n", filename)
		}

		fmt.Println("\n✨ All Flutter build scripts created successfully!")
		fmt.Println("Usage: ./scripts/build-apk.sh (or any other script)")
		return nil
	},
}

func init() {
	FlutterCmd.AddCommand(FlutterCreateScriptsCmd)
}
