package validation

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

// ChatmodeHeader represents the YAML frontmatter structure
type ChatmodeHeader struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Prompt      string `yaml:"prompt"`
	Tag         string `yaml:"tag"`
}

// ValidateChatmodeFile validates a chatmode file structure and content
func ValidateChatmodeFile(t *testing.T, filePath string) {
	require.FileExists(t, filePath)

	content, err := os.ReadFile(filePath)
	require.NoError(t, err)

	// Check YAML frontmatter
	ValidateYAMLFrontmatter(t, content)

	// Check file extension
	assert.True(t, strings.HasSuffix(filePath, ".chatmode.md"),
		"File should have .chatmode.md extension")
}

// ValidateYAMLFrontmatter validates YAML frontmatter in chatmode content
func ValidateYAMLFrontmatter(t *testing.T, content []byte) {
	contentStr := string(content)

	// Check for YAML frontmatter delimiters
	assert.True(t, strings.HasPrefix(contentStr, "---\n"),
		"File should start with YAML frontmatter")

	// Find the end of frontmatter
	lines := strings.Split(contentStr, "\n")
	var yamlLines []string
	frontmatterEnd := -1

	for i := 1; i < len(lines); i++ {
		if strings.TrimSpace(lines[i]) == "---" {
			frontmatterEnd = i
			break
		}
		yamlLines = append(yamlLines, lines[i])
	}

	require.Greater(t, frontmatterEnd, 0, "YAML frontmatter should be properly closed")

	// Parse YAML
	yamlContent := strings.Join(yamlLines, "\n")
	var header ChatmodeHeader
	err := yaml.Unmarshal([]byte(yamlContent), &header)
	require.NoError(t, err, "YAML frontmatter should be valid")

	// Validate required fields
	assert.NotEmpty(t, header.Name, "Name field is required")
	assert.NotEmpty(t, header.Description, "Description field is required")
	assert.NotEmpty(t, header.Prompt, "Prompt field is required")
}

// ValidateChatmateStructure validates the Chatmate directory structure
func ValidateChatmateStructure(t *testing.T, chatmateDir string) {
	require.DirExists(t, chatmateDir)

	// Check for mates directory
	matesDir := filepath.Join(chatmateDir, "mates")
	require.DirExists(t, matesDir)

	// Validate at least one chatmode file exists
	files, err := os.ReadDir(matesDir)
	require.NoError(t, err)

	foundChatmode := false
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".chatmode.md") {
			foundChatmode = true
			ValidateChatmodeFile(t, filepath.Join(matesDir, file.Name()))
		}
	}

	assert.True(t, foundChatmode, "At least one chatmode file should exist")
}

// ValidateVSCodeStructure validates VS Code directory structure
func ValidateVSCodeStructure(t *testing.T, vscodeDir string) {
	require.DirExists(t, vscodeDir)

	// Check for settings file
	settingsPath := filepath.Join(vscodeDir, "settings.json")
	if _, err := os.Stat(settingsPath); err == nil {
		// If settings exist, validate JSON structure
		content, err := os.ReadFile(settingsPath)
		require.NoError(t, err)

		// Basic JSON validation
		assert.True(t, strings.HasPrefix(strings.TrimSpace(string(content)), "{"))
		assert.True(t, strings.HasSuffix(strings.TrimSpace(string(content)), "}"))
	}
}

// ValidateChatmodeContent validates specific chatmode file content
func ValidateChatmodeContent(t *testing.T, content string, expectedName string) {
	// Check YAML frontmatter
	ValidateYAMLFrontmatter(t, []byte(content))

	// Check for expected name in content
	if expectedName != "" {
		assert.Contains(t, content, fmt.Sprintf("name: %s", expectedName))
	}

	// Check for markdown content after frontmatter
	parts := strings.Split(content, "---")
	require.GreaterOrEqual(t, len(parts), 3, "Should have opening and closing frontmatter")

	markdownContent := strings.TrimSpace(parts[2])
	if len(parts) > 3 {
		markdownContent = strings.TrimSpace(strings.Join(parts[2:], "---"))
	}

	// Markdown content can be empty for some chatmodes, so just check it's not nil
	assert.NotNil(t, markdownContent)
}

// ValidateFilePermissions validates file permissions are correct
func ValidateFilePermissions(t *testing.T, filePath string) {
	info, err := os.Stat(filePath)
	require.NoError(t, err)

	mode := info.Mode()
	assert.True(t, mode.IsRegular(), "Should be a regular file")

	// Check readable by owner
	assert.True(t, mode&0400 != 0, "File should be readable by owner")
}

// ValidateDirectoryPermissions validates directory permissions are correct
func ValidateDirectoryPermissions(t *testing.T, dirPath string) {
	info, err := os.Stat(dirPath)
	require.NoError(t, err)

	mode := info.Mode()
	assert.True(t, mode.IsDir(), "Should be a directory")

	// Check readable and executable by owner
	assert.True(t, mode&0500 != 0, "Directory should be readable and executable by owner")
}

// ValidateChatmodeFilename validates chatmode filename format
func ValidateChatmodeFilename(t *testing.T, filename string) {
	// Should end with .chatmode.md
	assert.True(t, strings.HasSuffix(filename, ".chatmode.md"),
		"Chatmode files should end with .chatmode.md")

	// Should not contain invalid characters
	invalidChars := regexp.MustCompile(`[<>:"/\\|?*]`)
	assert.False(t, invalidChars.MatchString(filename),
		"Filename should not contain invalid characters")

	// Should not be empty before extension
	baseName := strings.TrimSuffix(filename, ".chatmode.md")
	assert.NotEmpty(t, baseName, "Filename should not be empty")
}

// CountChatmodeFiles counts chatmode files in a directory
func CountChatmodeFiles(t *testing.T, dirPath string) int {
	files, err := os.ReadDir(dirPath)
	require.NoError(t, err)

	count := 0
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".chatmode.md") {
			count++
		}
	}

	return count
}

// ReadFileLines reads file and returns lines as slice
func ReadFileLines(t *testing.T, filePath string) []string {
	file, err := os.Open(filePath)
	require.NoError(t, err)
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	require.NoError(t, scanner.Err())
	return lines
}

// VerifyInstallation verifies that the expected number of chatmate files are installed
func VerifyInstallation(promptsDir string, expectedCount int) error {
	files, err := os.ReadDir(promptsDir)
	if err != nil {
		return fmt.Errorf("failed to read prompts directory: %w", err)
	}

	count := 0
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".chatmode.md") {
			count++
		}
	}

	if count != expectedCount {
		return fmt.Errorf("expected %d chatmode files, found %d", expectedCount, count)
	}

	return nil
}

// ValidateChatmodeFileErr validates a chatmode file and returns an error if invalid
func ValidateChatmodeFileErr(filePath string) error {
	if _, err := os.Stat(filePath); err != nil {
		return fmt.Errorf("file does not exist: %w", err)
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Check file extension
	if !strings.HasSuffix(filePath, ".chatmode.md") {
		return fmt.Errorf("file should have .chatmode.md extension")
	}

	// Check YAML frontmatter
	contentStr := string(content)
	if !strings.HasPrefix(contentStr, "---\n") {
		return fmt.Errorf("file should start with YAML frontmatter")
	}

	// Find the end of frontmatter
	lines := strings.Split(contentStr, "\n")
	frontmatterEnd := -1

	for i := 1; i < len(lines); i++ {
		if strings.TrimSpace(lines[i]) == "---" {
			frontmatterEnd = i
			break
		}
	}

	if frontmatterEnd <= 0 {
		return fmt.Errorf("YAML frontmatter should be properly closed")
	}

	// Extract and parse YAML
	yamlLines := lines[1:frontmatterEnd]
	yamlContent := strings.Join(yamlLines, "\n")

	var header ChatmodeHeader
	err = yaml.Unmarshal([]byte(yamlContent), &header)
	if err != nil {
		return fmt.Errorf("YAML frontmatter should be valid: %w", err)
	}

	// Validate required fields
	if header.Name == "" {
		return fmt.Errorf("name field is required")
	}
	if header.Description == "" {
		return fmt.Errorf("description field is required")
	}
	if header.Prompt == "" {
		return fmt.Errorf("prompt field is required")
	}

	return nil
}
