package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jonassiebler/chatmate/internal/assets"
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

// TestCountEmbeddedFilesFunction tests counting embedded chatmate files
func TestCountEmbeddedFilesFunction(t *testing.T) {
	chatmates, err := assets.GetEmbeddedMatesList()
	require.NoError(t, err, "Should be able to get embedded chatmate list")

	count := len(chatmates)
	assert.Greater(t, count, 0, "Should find embedded chatmate files")
	assert.Greater(t, count, 5, "Should have a reasonable number of chatmates")
}
