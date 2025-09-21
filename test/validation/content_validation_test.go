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

		helpers.ValidateYAMLFrontmatter(t, content)
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
