package helpers

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestEnvironment holds test environment configuration
type TestEnvironment struct {
	TempDir     string
	MockHome    string
	MockAppData string
	OriginalEnv map[string]string
}

// SetupTestEnvironment creates a clean test environment
func SetupTestEnvironment(t *testing.T) *TestEnvironment {
	t.Helper()

	tempDir, err := os.MkdirTemp("", "chatmate_tests_*")
	require.NoError(t, err)

	mockHome := filepath.Join(tempDir, "mock_home")
	mockAppData := filepath.Join(tempDir, "mock_appdata")

	err = os.MkdirAll(mockHome, 0755)
	require.NoError(t, err)
	err = os.MkdirAll(mockAppData, 0755)
	require.NoError(t, err)

	env := &TestEnvironment{
		TempDir:     tempDir,
		MockHome:    mockHome,
		MockAppData: mockAppData,
		OriginalEnv: make(map[string]string),
	}

	t.Cleanup(func() {
		env.Cleanup(t)
	})

	return env
}

// Cleanup restores original environment and removes temporary files
func (env *TestEnvironment) Cleanup(t *testing.T) {
	t.Helper()

	// Restore original environment variables
	for key, value := range env.OriginalEnv {
		if value == "" {
			os.Unsetenv(key)
		} else {
			os.Setenv(key, value)
		}
	}

	// Remove temporary directory
	if err := os.RemoveAll(env.TempDir); err != nil {
		t.Logf("Warning: failed to remove temp dir %s: %v", env.TempDir, err)
	}
}

// SimulateOS configures environment variables to simulate different operating systems
func (env *TestEnvironment) SimulateOS(t *testing.T, osType string) {
	t.Helper()

	// Store original values
	env.storeOriginalEnv("OSTYPE")
	env.storeOriginalEnv("HOME")
	env.storeOriginalEnv("APPDATA")

	switch strings.ToLower(osType) {
	case "macos", "darwin":
		os.Setenv("OSTYPE", "darwin")
		os.Setenv("HOME", env.MockHome)
	case "linux":
		os.Setenv("OSTYPE", "linux-gnu")
		os.Setenv("HOME", env.MockHome)
	case "windows", "msys", "cygwin":
		os.Setenv("OSTYPE", "msys")
		os.Setenv("HOME", env.MockHome)
		os.Setenv("APPDATA", env.MockAppData)
	default:
		t.Fatalf("Unsupported OS: %s", osType)
	}
}

func (env *TestEnvironment) storeOriginalEnv(key string) {
	env.OriginalEnv[key] = os.Getenv(key)
}

// CountFiles counts files matching a pattern in a directory
func CountFiles(dir, pattern string) (int, error) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return 0, nil
	}

	count := 0
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		matched, err := filepath.Match(pattern, filepath.Base(path))
		if err != nil {
			return err
		}
		if matched {
			count++
		}
		return nil
	})

	return count, err
}

// VerifyInstallation checks if installation completed successfully
func VerifyInstallation(targetDir string, expectedCount int) error {
	// Check directory exists
	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		return fmt.Errorf("target directory does not exist: %s", targetDir)
	}

	// Check file count matches
	actualCount, err := CountFiles(targetDir, "*.md")
	if err != nil {
		return fmt.Errorf("failed to count files: %w", err)
	}

	if actualCount != expectedCount {
		return fmt.Errorf("expected %d files, found %d", expectedCount, actualCount)
	}

	// Check files are readable
	err = filepath.WalkDir(targetDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if strings.HasSuffix(path, ".md") {
			_, err := os.Open(path)
			if err != nil {
				return fmt.Errorf("file %s is not readable: %w", path, err)
			}
		}
		return nil
	})

	return err
}

// FileContains checks if a file contains specific content
func FileContains(filePath, content string) (bool, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return false, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), content) {
			return true, nil
		}
	}

	return false, scanner.Err()
}

// FileNotContains checks if a file does not contain specific content
func FileNotContains(filePath, content string) (bool, error) {
	contains, err := FileContains(filePath, content)
	return !contains, err
}

// GetFileSize returns the size of a file in bytes
func GetFileSize(filePath string) (int64, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}

// SetupMockVSCode creates a mock VS Code environment for testing
func (env *TestEnvironment) SetupMockVSCode(t *testing.T, osType string) string {
	t.Helper()

	var promptsDir string
	switch strings.ToLower(osType) {
	case "macos", "darwin":
		promptsDir = filepath.Join(env.MockHome, "Library", "Application Support", "Code", "User", "prompts")
	case "linux":
		promptsDir = filepath.Join(env.MockHome, ".config", "Code", "User", "prompts")
	case "windows", "msys", "cygwin":
		promptsDir = filepath.Join(env.MockAppData, "Code", "User", "prompts")
	default:
		t.Fatalf("Unsupported OS for VS Code setup: %s", osType)
	}

	err := os.MkdirAll(promptsDir, 0755)
	require.NoError(t, err)

	return promptsDir
}

// ValidateYAMLFrontmatter checks if a markdown file has valid YAML frontmatter
func ValidateYAMLFrontmatter(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// Check first line
	if !scanner.Scan() {
		return fmt.Errorf("file is empty")
	}
	if strings.TrimSpace(scanner.Text()) != "---" {
		return fmt.Errorf("file does not start with YAML frontmatter")
	}

	// Find end of frontmatter
	lineNum := 1
	foundEnd := false
	for scanner.Scan() {
		lineNum++
		if strings.TrimSpace(scanner.Text()) == "---" && lineNum > 1 {
			foundEnd = true
			break
		}
	}

	if !foundEnd {
		return fmt.Errorf("YAML frontmatter not properly closed")
	}

	return scanner.Err()
}

// ValidateChatmodeFile checks if a chatmode file has all required elements
func ValidateChatmodeFile(filePath string) error {
	// Check file exists and is readable
	_, err := os.Stat(filePath)
	if err != nil {
		return fmt.Errorf("file not accessible: %w", err)
	}

	// Check YAML frontmatter
	if err := ValidateYAMLFrontmatter(filePath); err != nil {
		return fmt.Errorf("invalid YAML frontmatter: %w", err)
	}

	// Read file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	contentStr := string(content)

	// Extract YAML section
	yamlRegex := regexp.MustCompile(`(?s)^---\n(.*?)\n---`)
	yamlMatch := yamlRegex.FindStringSubmatch(contentStr)
	if yamlMatch == nil {
		return fmt.Errorf("could not extract YAML frontmatter")
	}

	yamlContent := yamlMatch[1]

	// Check for required YAML fields
	requiredFields := []string{"description:", "author:"}
	for _, field := range requiredFields {
		if !strings.Contains(yamlContent, field) {
			return fmt.Errorf("missing required YAML field: %s", field)
		}
	}

	// Check for content after frontmatter
	contentAfterYAML := strings.Split(contentStr, "---")[2]
	if len(strings.TrimSpace(contentAfterYAML)) < 100 {
		return fmt.Errorf("insufficient content after frontmatter")
	}

	return nil
}

// CreateTestChatmode creates a temporary chatmode file for testing
func CreateTestChatmode(filePath, title, description string) error {
	content := fmt.Sprintf(`---
description: "%s"
author: "Test Suite"
version: "1.0.0"
category: "testing"
tags: ["test", "automation"]
---

# %s

This is a test chatmode file created for testing purposes.

## Purpose

This file is used to test the chatmate installation and validation system.

## Instructions

1. This is a test instruction
2. This is another test instruction
3. Final test instruction

## Examples

Example usage of this test chatmode.

## Notes

Additional notes for testing purposes.
`, description, title)

	return os.WriteFile(filePath, []byte(content), 0644)
}

// ValidateAllMates validates all chatmode files in the mates directory
func ValidateAllMates(matesDir string) []error {
	var errors []error

	err := filepath.WalkDir(matesDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			errors = append(errors, fmt.Errorf("failed to access %s: %w", path, err))
			return nil
		}

		if d.IsDir() || !strings.HasSuffix(path, ".chatmode.md") {
			return nil
		}

		if err := ValidateChatmodeFile(path); err != nil {
			errors = append(errors, fmt.Errorf("validation failed for %s: %w", filepath.Base(path), err))
		}

		return nil
	})

	if err != nil {
		errors = append(errors, fmt.Errorf("failed to walk mates directory: %w", err))
	}

	return errors
}

// GetVSCodePromptsDir returns the VS Code prompts directory for the current OS
func GetVSCodePromptsDir() (string, error) {
	var baseDir string
	var subPath string

	switch runtime.GOOS {
	case "darwin":
		baseDir = os.Getenv("HOME")
		subPath = "Library/Application Support/Code/User/prompts"
	case "linux":
		baseDir = os.Getenv("HOME")
		subPath = ".config/Code/User/prompts"
	case "windows":
		baseDir = os.Getenv("APPDATA")
		subPath = "Code/User/prompts"
	default:
		return "", fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	if baseDir == "" {
		return "", fmt.Errorf("required environment variable not set")
	}

	return filepath.Join(baseDir, subPath), nil
}

// AssertFileExists is a test helper that asserts a file exists
func AssertFileExists(t *testing.T, filePath string) {
	t.Helper()
	_, err := os.Stat(filePath)
	require.NoError(t, err, "File should exist: %s", filePath)
}

// AssertFileNotExists is a test helper that asserts a file does not exist
func AssertFileNotExists(t *testing.T, filePath string) {
	t.Helper()
	_, err := os.Stat(filePath)
	require.True(t, os.IsNotExist(err), "File should not exist: %s", filePath)
}

// AssertDirectoryExists is a test helper that asserts a directory exists
func AssertDirectoryExists(t *testing.T, dirPath string) {
	t.Helper()
	info, err := os.Stat(dirPath)
	require.NoError(t, err, "Directory should exist: %s", dirPath)
	require.True(t, info.IsDir(), "Path should be a directory: %s", dirPath)
}

// CopyFile copies a file from src to dst
func CopyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = destFile.ReadFrom(sourceFile)
	return err
}
