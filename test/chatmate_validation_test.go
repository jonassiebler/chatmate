package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jonassiebler/chatmate/internal/testing/helpers"
)

// TestMatesDirectoryExistsAndContainsFiles tests that the mates directory exists and contains files
func TestMatesDirectoryExistsAndContainsFiles(t *testing.T) {
	// Check mates directory exists
	_, err := os.Stat("../mates")
	require.NoError(t, err, "../mates directory should exist")

	// Check it's actually a directory
	info, err := os.Stat("../mates")
	require.NoError(t, err)
	assert.True(t, info.IsDir(), "../mates should be a directory")

	// Count files
	fileCount, err := helpers.CountFiles("../mates", "*.md")
	require.NoError(t, err)
	assert.Greater(t, fileCount, 0, "../mates directory should contain at least one markdown file")
}

// TestAllChatmateFilesFollowNamingConvention tests that all .md files follow the .chatmode.md naming convention
func TestAllChatmateFilesFollowNamingConvention(t *testing.T) {
	err := filepath.WalkDir("../mates", func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		filename := d.Name()
		if strings.HasSuffix(filename, ".md") {
			assert.True(t, strings.HasSuffix(filename, ".chatmode.md"),
				"All markdown files in mates directory should end with .chatmode.md, found: %s", filename)
		}
		return nil
	})
	require.NoError(t, err, "Should be able to walk mates directory")
}

// TestChatmateFilesAreNotEmpty tests that chatmate files contain content
func TestChatmateFilesAreNotEmpty(t *testing.T) {
	err := filepath.WalkDir("../mates", func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(d.Name(), ".chatmode.md") {
			return nil
		}

		info, err := d.Info()
		require.NoError(t, err, "Should be able to get file info for %s", path)
		assert.Greater(t, info.Size(), int64(0), "Chatmate file should not be empty: %s", d.Name())
		return nil
	})
	require.NoError(t, err, "Should be able to walk mates directory")
}

// TestChatmateFilesContainRequiredHeaders tests that chatmate files contain markdown headers
func TestChatmateFilesContainRequiredHeaders(t *testing.T) {
	err := filepath.WalkDir("../mates", func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(d.Name(), ".chatmode.md") {
			return nil
		}

		content, err := os.ReadFile(path)
		require.NoError(t, err, "Should be able to read file %s", path)

		contentStr := string(content)
		assert.True(t, strings.Contains(contentStr, "#"),
			"Chatmate file should contain markdown headers: %s", d.Name())
		return nil
	})
	require.NoError(t, err, "Should be able to walk mates directory")
}

// TestChatmateFilesContainYAMLFrontmatter tests that chatmate files have proper YAML frontmatter
func TestChatmateFilesContainYAMLFrontmatter(t *testing.T) {
	err := filepath.WalkDir("../mates", func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(d.Name(), ".chatmode.md") {
			return nil
		}

		err = helpers.ValidateYAMLFrontmatter(path)
		assert.NoError(t, err, "Chatmate file should have valid YAML frontmatter: %s", d.Name())
		return nil
	})
	require.NoError(t, err, "Should be able to walk mates directory")
}

// TestChatmateFilesContainRequiredFields tests that YAML frontmatter contains required fields
func TestChatmateFilesContainRequiredFields(t *testing.T) {
	requiredFields := []string{"description:", "author:"}

	err := filepath.WalkDir("../mates", func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(d.Name(), ".chatmode.md") {
			return nil
		}

		content, err := os.ReadFile(path)
		require.NoError(t, err, "Should be able to read file %s", path)

		contentStr := string(content)
		for _, field := range requiredFields {
			assert.True(t, strings.Contains(contentStr, field),
				"Chatmate file should contain required field '%s': %s", field, d.Name())
		}
		return nil
	})
	require.NoError(t, err, "Should be able to walk mates directory")
}

// TestChatmateFilesHaveMinimumContent tests that chatmate files have sufficient content
func TestChatmateFilesHaveMinimumContent(t *testing.T) {
	err := filepath.WalkDir("../mates", func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(d.Name(), ".chatmode.md") {
			return nil
		}

		content, err := os.ReadFile(path)
		require.NoError(t, err, "Should be able to read file %s", path)

		// Count non-empty lines after YAML frontmatter
		lines := strings.Split(string(content), "\n")
		inYAML := false
		yamlClosed := false
		contentLineCount := 0

		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "---" {
				if !inYAML {
					inYAML = true
				} else if inYAML && !yamlClosed {
					yamlClosed = true
				}
			} else if yamlClosed && len(line) > 0 {
				contentLineCount++
			}
		}

		assert.Greater(t, contentLineCount, 5,
			"Chatmate file should have meaningful content after frontmatter: %s", d.Name())
		return nil
	})
	require.NoError(t, err, "Should be able to walk mates directory")
}

// TestValidateChatmodeFileFunction tests the ValidateChatmodeFile helper function
func TestValidateChatmodeFileFunction(t *testing.T) {
	// Test with valid file
	validFile := "fixtures/sample.chatmode.md"
	err := helpers.ValidateChatmodeFile(validFile)
	assert.NoError(t, err, "Valid chatmode file should pass validation")

	// Test with empty file
	emptyFile := "fixtures/empty.chatmode.md"
	err = helpers.ValidateChatmodeFile(emptyFile)
	assert.Error(t, err, "Empty chatmode file should fail validation")

	// Test with file without YAML
	noYAMLFile := "fixtures/invalid_no_yaml.chatmode.md"
	err = helpers.ValidateChatmodeFile(noYAMLFile)
	assert.Error(t, err, "File without YAML frontmatter should fail validation")

	// Test with incomplete YAML
	incompleteYAMLFile := "fixtures/invalid_incomplete_yaml.chatmode.md"
	err = helpers.ValidateChatmodeFile(incompleteYAMLFile)
	assert.Error(t, err, "File with incomplete YAML should fail validation")
}

// TestValidateAllMatesFunction tests the ValidateAllMates helper function
func TestValidateAllMatesFunction(t *testing.T) {
	errors := helpers.ValidateAllMates("../mates")

	if len(errors) > 0 {
		for _, err := range errors {
			t.Logf("Validation error: %v", err)
		}
		t.Errorf("Found %d validation errors in mates directory", len(errors))
	}
}

// TestCreateTestChatmodeFunction tests the CreateTestChatmode helper function
func TestCreateTestChatmodeFunction(t *testing.T) {
	env := helpers.SetupTestEnvironment(t)

	testFile := filepath.Join(env.TempDir, "test.chatmode.md")
	err := helpers.CreateTestChatmode(testFile, "Test Chatmate", "A test chatmate for validation")
	require.NoError(t, err, "Should be able to create test chatmode file")

	// Verify file was created
	helpers.AssertFileExists(t, testFile)

	// Verify file is valid
	err = helpers.ValidateChatmodeFile(testFile)
	assert.NoError(t, err, "Created test chatmode should be valid")

	// Verify content
	contains, err := helpers.FileContains(testFile, "Test Chatmate")
	require.NoError(t, err)
	assert.True(t, contains, "File should contain the title")

	contains, err = helpers.FileContains(testFile, "A test chatmate for validation")
	require.NoError(t, err)
	assert.True(t, contains, "File should contain the description")
}

// TestCountFilesFunction tests the CountFiles helper function
func TestCountFilesFunction(t *testing.T) {
	// Test with mates directory
	count, err := helpers.CountFiles("../mates", "*.chatmode.md")
	require.NoError(t, err)
	assert.Greater(t, count, 0, "Should find chatmode files in mates directory")

	// Test with non-existent directory
	count, err = helpers.CountFiles("nonexistent", "*.md")
	require.NoError(t, err)
	assert.Equal(t, 0, count, "Non-existent directory should return 0 files")

	// Test with fixtures directory
	count, err = helpers.CountFiles("fixtures", "*.md")
	require.NoError(t, err)
	assert.Greater(t, count, 0, "Should find test fixture files")
}
