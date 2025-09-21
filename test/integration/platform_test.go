package main

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/jonassiebler/chatmate/internal/testing/helpers"
	"github.com/stretchr/testify/suite"
)

// PlatformCompatibilitySuite tests cross-platform functionality
type PlatformCompatibilitySuite struct {
	BaseIntegrationSuite
}

// TestCrossPlatformCompatibility tests file operations across platforms
func (s *PlatformCompatibilitySuite) TestCrossPlatformCompatibility() {
	// Test file paths and operations work across platforms
	testOSes := []string{"macos", "linux"}
	if runtime.GOOS == "windows" {
		testOSes = []string{"windows"}
	}

	for _, testOS := range testOSes {
		s.Run(testOS, func() {
			s.env.SimulateOS(s.T(), testOS)
			promptsDir := s.env.SetupMockVSCode(s.T(), testOS)

			// Verify prompts directory was created with correct path structure
			s.Assert().True(filepath.IsAbs(promptsDir), "Prompts dir should be absolute path")

			// Test file creation
			testFile := filepath.Join(promptsDir, "test.md")
			err := os.WriteFile(testFile, []byte("test content"), 0644)
			s.Assert().NoError(err, "Should be able to create test file")

			// Verify file can be read
			_, err = os.ReadFile(testFile)
			s.Assert().NoError(err, "Should be able to read test file")
		})
	}
}

// TestFileValidationWorkflow tests end-to-end file validation
func (s *PlatformCompatibilitySuite) TestFileValidationWorkflow() {
	// Test workflow with various file types
	testFiles := []struct {
		name    string
		content string
		valid   bool
	}{
		{
			name: "valid.chatmode.md",
			content: `---
name: "Valid Test Chatmate"
description: "Valid test chatmate"
prompt: "Test prompt for validation"
author: "Test Suite"
version: "1.0.0"
---

# Valid Chatmate

This is a valid chatmate for testing.

## Instructions

1. Test instruction
2. Another test instruction

## Examples

Example content here.`,
			valid: true,
		},
		{
			name: "invalid.chatmode.md",
			content: `# Invalid Chatmate

No YAML frontmatter here.`,
			valid: false,
		},
	}

	for _, testFile := range testFiles {
		s.Run(testFile.name, func() {
			filePath := filepath.Join(s.env.TempDir, testFile.name)
			err := os.WriteFile(filePath, []byte(testFile.content), 0644)
			s.Require().NoError(err)

			err = helpers.ValidateChatmodeFileErr(filePath)
			if testFile.valid {
				s.Assert().NoError(err, "Valid file should pass validation")
			} else {
				s.Assert().Error(err, "Invalid file should fail validation")
			}
		})
	}
}

// TestPlatformCompatibilitySuite runs the platform compatibility test suite
func TestPlatformCompatibilitySuite(t *testing.T) {
	suite.Run(t, new(PlatformCompatibilitySuite))
}
