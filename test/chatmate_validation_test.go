package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jonassiebler/chatmate/internal/assets"
	"github.com/jonassiebler/chatmate/internal/testing/helpers"
)

// TestEmbeddedMatesExistAndContainFiles tests that the embedded mates exist and contain files
func TestEmbeddedMatesExistAndContainFiles(t *testing.T) {
	// Get embedded chatmate list
	chatmates, err := assets.GetEmbeddedMatesList()
	require.NoError(t, err, "Should be able to get embedded chatmate list")

	// Check we have files
	assert.Greater(t, len(chatmates), 0, "Embedded mates should contain at least one chatmate file")
}

// TestAllEmbeddedChatmateFilesFollowNamingConvention tests that all embedded files follow the .chatmode.md naming convention
func TestAllEmbeddedChatmateFilesFollowNamingConvention(t *testing.T) {
	chatmates, err := assets.GetEmbeddedMatesList()
	require.NoError(t, err, "Should be able to get embedded chatmate list")

	for _, filename := range chatmates {
		assert.True(t, strings.HasSuffix(filename, ".chatmode.md"),
			"All embedded chatmate files should end with .chatmode.md, found: %s", filename)
	}
}

// TestEmbeddedChatmateFilesAreNotEmpty tests that embedded chatmate files are not empty
func TestEmbeddedChatmateFilesAreNotEmpty(t *testing.T) {
	chatmates, err := assets.GetEmbeddedMatesList()
	require.NoError(t, err, "Should be able to get embedded chatmate list")

	for _, filename := range chatmates {
		content, err := assets.GetEmbeddedMateContent(filename)
		require.NoError(t, err, "Should be able to read embedded file %s", filename)
		assert.Greater(t, len(content), 0, "Embedded chatmate file should not be empty: %s", filename)
	}
}

// TestEmbeddedChatmateFilesContainRequiredHeaders tests that embedded chatmate files contain markdown headers
func TestEmbeddedChatmateFilesContainRequiredHeaders(t *testing.T) {
	chatmates, err := assets.GetEmbeddedMatesList()
	require.NoError(t, err, "Should be able to get embedded chatmate list")

	for _, filename := range chatmates {
		content, err := assets.GetEmbeddedMateContent(filename)
		require.NoError(t, err, "Should be able to read embedded file %s", filename)

		contentStr := string(content)
		assert.True(t, strings.Contains(contentStr, "#"),
			"Embedded chatmate file should contain markdown headers: %s", filename)
	}
}

// TestEmbeddedChatmateFilesContainYAMLFrontmatter tests that embedded chatmate files have proper YAML frontmatter
func TestEmbeddedChatmateFilesContainYAMLFrontmatter(t *testing.T) {
	chatmates, err := assets.GetEmbeddedMatesList()
	require.NoError(t, err, "Should be able to get embedded chatmate list")

	for _, filename := range chatmates {
		content, err := assets.GetEmbeddedMateContent(filename)
		require.NoError(t, err, "Should be able to read embedded file %s", filename)

		contentStr := string(content)
		assert.True(t, strings.HasPrefix(contentStr, "---"),
			"Embedded chatmate file should start with YAML frontmatter: %s", filename)

		// Check for closing ---
		lines := strings.Split(contentStr, "\n")
		foundClosing := false
		for i := 1; i < len(lines) && i < 20; i++ { // Check first 20 lines
			if strings.TrimSpace(lines[i]) == "---" {
				foundClosing = true
				break
			}
		}
		assert.True(t, foundClosing, "Embedded chatmate file should have closing YAML frontmatter: %s", filename)
	}
}

// TestEmbeddedChatmateFilesContainRequiredFields tests that YAML frontmatter contains required fields
func TestEmbeddedChatmateFilesContainRequiredFields(t *testing.T) {
	requiredFields := []string{"description:"}

	chatmates, err := assets.GetEmbeddedMatesList()
	require.NoError(t, err, "Should be able to get embedded chatmate list")

	for _, filename := range chatmates {
		content, err := assets.GetEmbeddedMateContent(filename)
		require.NoError(t, err, "Should be able to read embedded file %s", filename)

		contentStr := string(content)
		for _, field := range requiredFields {
			assert.True(t, strings.Contains(contentStr, field),
				"Embedded chatmate file should contain required field '%s': %s", field, filename)
		}
	}
}

// TestEmbeddedChatmateFilesHaveMinimumContent tests that embedded chatmate files have substantial content
func TestEmbeddedChatmateFilesHaveMinimumContent(t *testing.T) {
	minContentLength := 500 // At least 500 characters

	chatmates, err := assets.GetEmbeddedMatesList()
	require.NoError(t, err, "Should be able to get embedded chatmate list")

	for _, filename := range chatmates {
		content, err := assets.GetEmbeddedMateContent(filename)
		require.NoError(t, err, "Should be able to read embedded file %s", filename)

		assert.GreaterOrEqual(t, len(content), minContentLength,
			"Embedded chatmate file should have at least %d characters: %s (has %d)",
			minContentLength, filename, len(content))
	}
}

// TestChatmateFilesContainRequiredHeaders tests that chatmate files contain markdown headers
func TestChatmateFilesContainRequiredHeaders(t *testing.T) {
	// Get embedded chatmates instead of filesystem
	chatmates, err := assets.GetEmbeddedMatesList()
	require.NoError(t, err, "Should be able to get embedded chatmate list")

	for _, filename := range chatmates {
		content, err := assets.GetEmbeddedMateContent(filename)
		require.NoError(t, err, "Should be able to read embedded file %s", filename)

		contentStr := string(content)
		assert.True(t, strings.Contains(contentStr, "#"),
			"Chatmate file should contain markdown headers: %s", filename)
	}
}

// TestChatmateFilesContainYAMLFrontmatter tests that chatmate files have proper YAML frontmatter
func TestChatmateFilesContainYAMLFrontmatter(t *testing.T) {
	// Get embedded chatmates instead of filesystem
	chatmates, err := assets.GetEmbeddedMatesList()
	require.NoError(t, err, "Should be able to get embedded chatmate list")

	for _, filename := range chatmates {
		// Create temporary file to test with existing helper function
		content, err := assets.GetEmbeddedMateContent(filename)
		require.NoError(t, err, "Should be able to read embedded file %s", filename)

		tempFile := filepath.Join(t.TempDir(), filename)
		err = os.WriteFile(tempFile, content, 0644)
		require.NoError(t, err, "Should be able to create temp file %s", filename)

		err = helpers.ValidateYAMLFrontmatter(tempFile)
		assert.NoError(t, err, "Chatmate file should have valid YAML frontmatter: %s", filename)
	}
}

// TestChatmateFilesContainRequiredFields tests that YAML frontmatter contains required fields
func TestChatmateFilesContainRequiredFields(t *testing.T) {
	requiredFields := []string{"description:", "author:"}

	// Get embedded chatmates instead of filesystem
	chatmates, err := assets.GetEmbeddedMatesList()
	require.NoError(t, err, "Should be able to get embedded chatmate list")

	for _, filename := range chatmates {
		content, err := assets.GetEmbeddedMateContent(filename)
		require.NoError(t, err, "Should be able to read embedded file %s", filename)

		contentStr := string(content)
		for _, field := range requiredFields {
			assert.True(t, strings.Contains(contentStr, field),
				"Chatmate file should contain required field '%s': %s", field, filename)
		}
	}
}

// TestChatmateFilesHaveMinimumContent tests that chatmate files have sufficient content
func TestChatmateFilesHaveMinimumContent(t *testing.T) {
	// Get embedded chatmates instead of filesystem
	chatmates, err := assets.GetEmbeddedMatesList()
	require.NoError(t, err, "Should be able to get embedded chatmate list")

	for _, filename := range chatmates {
		content, err := assets.GetEmbeddedMateContent(filename)
		require.NoError(t, err, "Should be able to read embedded file %s", filename)

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
			"Chatmate file should have meaningful content after frontmatter: %s", filename)
	}
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

// TestValidateAllEmbeddedMatesFunction tests validation of all embedded chatmates
func TestValidateAllEmbeddedMatesFunction(t *testing.T) {
	chatmates, err := assets.GetEmbeddedMatesList()
	require.NoError(t, err, "Should be able to get embedded chatmate list")

	errors := 0
	for _, filename := range chatmates {
		content, err := assets.GetEmbeddedMateContent(filename)
		if err != nil {
			t.Logf("Validation error for %s: %v", filename, err)
			errors++
			continue
		}

		// Basic validation - check for YAML frontmatter and content
		contentStr := string(content)
		if !strings.HasPrefix(contentStr, "---") {
			t.Logf("Validation error for %s: missing YAML frontmatter", filename)
			errors++
		}

		if len(contentStr) < 500 {
			t.Logf("Validation error for %s: insufficient content length", filename)
			errors++
		}
	}

	if errors > 0 {
		t.Errorf("Found %d validation errors in embedded chatmates", errors)
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

// TestCountEmbeddedFilesFunction tests counting embedded chatmate files
func TestCountEmbeddedFilesFunction(t *testing.T) {
	chatmates, err := assets.GetEmbeddedMatesList()
	require.NoError(t, err, "Should be able to get embedded chatmate list")

	count := len(chatmates)
	assert.Greater(t, count, 0, "Should find embedded chatmate files")
	assert.Greater(t, count, 5, "Should have a reasonable number of chatmates")

	// Test with fixtures directory for comparison
	count, err = helpers.CountFiles("fixtures", "*.md")
	require.NoError(t, err)
	assert.Greater(t, count, 0, "Should find test fixture files")
}
