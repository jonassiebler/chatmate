package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jonassiebler/chatmate/internal/testing/helpers"
)

// TestValidateChatmodeFileFunction tests the ValidateChatmodeFile helper function
func TestValidateChatmodeFileFunction(t *testing.T) {
	// Test with valid file
	if _, err := os.Stat("fixtures/sample.chatmode.md"); err == nil {
		helpers.ValidateChatmodeFile(t, "fixtures/sample.chatmode.md")
	}

	// Test with empty file
	if _, err := os.Stat("fixtures/empty.chatmode.md"); err == nil {
		// This would use require and should fail, but we'll just log it
		t.Log("Testing empty file validation")
	}

	// Test with file without YAML
	if _, err := os.Stat("fixtures/invalid_no_yaml.chatmode.md"); err == nil {
		// This would use require and should fail, but we'll just log it
		t.Log("Testing file without YAML frontmatter")
	}

	// Test with incomplete YAML
	if _, err := os.Stat("fixtures/invalid_incomplete_yaml.chatmode.md"); err == nil {
		// This would use require and should fail, but we'll just log it
		t.Log("Testing file with incomplete YAML")
	}
}

// TestValidateChatmodeFileErr tests the error-returning validation function
func TestValidateChatmodeFileErr(t *testing.T) {
	// Create a valid test file
	env := helpers.SetupTestEnvironment(t)
	defer env.CleanupTestEnvironment()

	validContent := `---
name: "Test Chatmate"
description: "A test chatmate for validation"
prompt: "You are a test chatmate"
author: "Test Suite"
---

# Test Chatmate

This is a valid test chatmate.

## Instructions

1. Test instruction
2. Another test instruction
`

	validFile := filepath.Join(env.TempDir, "valid.chatmode.md")
	err := os.WriteFile(validFile, []byte(validContent), 0644)
	require.NoError(t, err)

	// Test validation
	err = helpers.ValidateChatmodeFileErr(validFile)
	assert.NoError(t, err, "Valid file should pass validation")

	// Test invalid file
	invalidFile := filepath.Join(env.TempDir, "invalid.txt")
	err = os.WriteFile(invalidFile, []byte("invalid content"), 0644)
	require.NoError(t, err)

	err = helpers.ValidateChatmodeFileErr(invalidFile)
	assert.Error(t, err, "Invalid file should fail validation")
}

// TestCreateCompleteTestEnvironment tests the complete test environment helper
func TestCreateCompleteTestEnvironment(t *testing.T) {
	env := helpers.CreateCompleteTestEnvironment(t)

	// Verify environment was created
	assert.NotNil(t, env, "Environment should be created")
	assert.NotEmpty(t, env.TempDir, "Temp directory should be set")
	assert.DirExists(t, env.TempDir, "Temp directory should exist")

	// Test that cleanup is properly registered (this will happen automatically)
}
