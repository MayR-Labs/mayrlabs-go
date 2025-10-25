package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/MayR-Labs/mayrlabs-go/internal/utils"
	"github.com/spf13/cobra"
)

// CICmd generates CI/CD workflow files
var CICmd = &cobra.Command{
	Use:   "ci",
	Short: "Generate CI/CD workflow YAML for your language and VCS",
	Long:  "Create CI/CD configuration files for GitHub Actions, GitLab CI, or other platforms",
	RunE: func(cmd *cobra.Command, args []string) error {
		lang, _ := cmd.Flags().GetString("lang")
		vcs, _ := cmd.Flags().GetString("vcs")

		// Prompt for missing values
		if lang == "" {
			fmt.Println("Available languages: go, javascript, python, php, flutter, dart")
			var err error
			lang, err = utils.PromptInput("Select language: ")
			if err != nil {
				return err
			}
		}

		if vcs == "" {
			fmt.Println("Available VCS: github, gitlab, circleci")
			var err error
			vcs, err = utils.PromptInput("Select VCS: ")
			if err != nil {
				return err
			}
		}

		return generateCI(lang, vcs)
	},
}

func init() {
	CICmd.Flags().StringP("lang", "l", "", "Programming language")
	CICmd.Flags().StringP("vcs", "v", "", "Version control system (github, gitlab, circleci)")
}

func generateCI(lang, vcs string) error {
	var content string
	var filePath string

	switch vcs {
	case "github":
		content = generateGitHubActions(lang)
		filePath = ".github/workflows/ci.yml"

	case "gitlab":
		content = generateGitLabCI(lang)
		filePath = ".gitlab-ci.yml"

	case "circleci":
		content = generateCircleCI(lang)
		filePath = ".circleci/config.yml"

	default:
		return fmt.Errorf("unsupported VCS: %s", vcs)
	}

	// Create directory if needed
	dir := filepath.Dir(filePath)
	if dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}

	// Write file
	if err := utils.WriteFile(filePath, content); err != nil {
		return fmt.Errorf("failed to write CI config: %w", err)
	}

	fmt.Printf("✅ CI/CD config created at %s!\n", filePath)
	return nil
}

func generateGitHubActions(lang string) string {
	switch lang {
	case "go", "golang":
		return `name: Go CI

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main, develop ]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
    
    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.21'
    
    - name: Install dependencies
      run: go mod download
    
    - name: Run tests
      run: go test -v ./...
    
    - name: Build
      run: go build -v ./...
    
    - name: Run linter
      uses: golangci/golangci-lint-action@v3
      with:
        version: latest
`

	case "javascript", "js", "node":
		return `name: Node.js CI

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main, develop ]

jobs:
  build:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        node-version: [16.x, 18.x, 20.x]
    
    steps:
    - uses: actions/checkout@v3
    
    - name: Use Node.js ${{ matrix.node-version }}
      uses: actions/setup-node@v3
      with:
        node-version: ${{ matrix.node-version }}
    
    - name: Install dependencies
      run: npm ci
    
    - name: Run tests
      run: npm test
    
    - name: Build
      run: npm run build --if-present
`

	case "python", "py":
		return `name: Python CI

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main, develop ]

jobs:
  build:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        python-version: ['3.8', '3.9', '3.10', '3.11']
    
    steps:
    - uses: actions/checkout@v3
    
    - name: Set up Python ${{ matrix.python-version }}
      uses: actions/setup-python@v4
      with:
        python-version: ${{ matrix.python-version }}
    
    - name: Install dependencies
      run: |
        python -m pip install --upgrade pip
        pip install -r requirements.txt
    
    - name: Run tests
      run: pytest
    
    - name: Lint with flake8
      run: |
        pip install flake8
        flake8 .
`

	case "flutter", "dart":
		return `name: Flutter CI

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main, develop ]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
    
    - name: Set up Flutter
      uses: subosito/flutter-action@v2
      with:
        flutter-version: '3.x'
    
    - name: Install dependencies
      run: flutter pub get
    
    - name: Run tests
      run: flutter test
    
    - name: Analyze
      run: flutter analyze
    
    - name: Build APK
      run: flutter build apk --release
`

	case "php":
		return `name: PHP CI

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main, develop ]

jobs:
  build:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        php-version: ['8.0', '8.1', '8.2']
    
    steps:
    - uses: actions/checkout@v3
    
    - name: Setup PHP
      uses: shivammathur/setup-php@v2
      with:
        php-version: ${{ matrix.php-version }}
    
    - name: Install dependencies
      run: composer install
    
    - name: Run tests
      run: vendor/bin/phpunit
`

	default:
		return `name: CI

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main, develop ]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
    
    - name: Build
      run: echo "Add your build steps here"
    
    - name: Test
      run: echo "Add your test steps here"
`
	}
}

func generateGitLabCI(lang string) string {
	switch lang {
	case "go", "golang":
		return `image: golang:1.21

stages:
  - test
  - build

test:
  stage: test
  script:
    - go test -v ./...

build:
  stage: build
  script:
    - go build -v ./...
  artifacts:
    paths:
      - ./
`

	case "javascript", "js", "node":
		return `image: node:18

stages:
  - test
  - build

test:
  stage: test
  script:
    - npm ci
    - npm test

build:
  stage: build
  script:
    - npm ci
    - npm run build
  artifacts:
    paths:
      - dist/
`

	default:
		return `stages:
  - test
  - build

test:
  stage: test
  script:
    - echo "Add your test commands"

build:
  stage: build
  script:
    - echo "Add your build commands"
`
	}
}

func generateCircleCI(lang string) string {
	return `version: 2.1

jobs:
  build:
    docker:
      - image: cimg/base:stable
    steps:
      - checkout
      - run: echo "Add your build commands"

workflows:
  build-and-test:
    jobs:
      - build
`
}
