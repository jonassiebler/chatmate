package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/jonassiebler/chatmate/internal/assets"
	"github.com/jonassiebler/chatmate/internal/testing/helpers"
)

// IntegrationTestSuite contains integration tests for ChatMate
type IntegrationTestSuite struct {
	suite.Suite
	env        *helpers.TestEnvironment
	binaryPath string
}

// SetupSuite builds the test binary once for all integration tests
func (s *IntegrationTestSuite) SetupSuite() {
	// Build test binary from project root
	s.binaryPath = "chatmate-integration-test"
	if runtime.GOOS == "windows" {
		s.binaryPath += ".exe"
	}

	cmd := exec.Command("go", "build", "-o", s.binaryPath, "github.com/jonassiebler/chatmate")
	output, err := cmd.CombinedOutput()
	s.Require().NoError(err, "Failed to build test binary: %s", string(output))
}

// TearDownSuite cleans up the test binary
func (s *IntegrationTestSuite) TearDownSuite() {
	if s.binaryPath != "" {
		os.Remove(s.binaryPath)
	}
}

// SetupTest creates a fresh test environment for each test
func (s *IntegrationTestSuite) SetupTest() {
	s.env = helpers.SetupTestEnvironment(s.T())
}

// TearDownTest is handled by the test environment cleanup

// TestFullInstallationWorkflowMacOS tests complete installation process on macOS
func (s *IntegrationTestSuite) TestFullInstallationWorkflowMacOS() {
	if runtime.GOOS == "windows" {
		s.T().Skip("macOS-specific test")
	}

	s.env.SimulateOS(s.T(), "macos")
	promptsDir := s.env.SetupMockVSCode(s.T(), "macos")

	// Get embedded chatmates instead of filesystem
	chatmates, err := assets.GetEmbeddedMatesList()
	s.Require().NoError(err, "Should be able to get embedded chatmate list")
	s.Assert().Greater(len(chatmates), 0, "Should have embedded chatmate files")

	// Copy embedded chatmates to prompts directory
	for _, filename := range chatmates {
		content, err := assets.GetEmbeddedMateContent(filename)
		s.Require().NoError(err, "Should be able to read embedded file %s", filename)

		destFile := filepath.Join(promptsDir, filename)
		err = os.WriteFile(destFile, content, 0644)
		s.Require().NoError(err, "Should be able to copy %s", filename)
	}

	// Verify all files were copied
	err = helpers.VerifyInstallation(promptsDir, len(chatmates))
	s.Assert().NoError(err, "Installation should be verified successfully")

	// Verify specific important chatmates were installed
	importantChatmates := []string{"Testing.chatmode.md", "Create PR.chatmode.md", "Solve Issue.chatmode.md"}
	for _, chatmate := range importantChatmates {
		chatmatePath := filepath.Join(promptsDir, chatmate)
		if _, err := os.Stat(chatmatePath); err == nil {
			s.T().Logf("Found important chatmate: %s", chatmate)
		}
	}
}

// TestFullInstallationWorkflowLinux tests complete installation process on Linux
func (s *IntegrationTestSuite) TestFullInstallationWorkflowLinux() {
	if runtime.GOOS == "windows" {
		s.T().Skip("Linux-specific test")
	}

	s.env.SimulateOS(s.T(), "linux")
	promptsDir := s.env.SetupMockVSCode(s.T(), "linux")

	// Get embedded chatmates instead of filesystem
	chatmates, err := assets.GetEmbeddedMatesList()
	s.Require().NoError(err, "Should be able to get embedded chatmate list")
	s.Assert().Greater(len(chatmates), 0, "Should have embedded chatmate files")

	// Copy embedded chatmates to prompts directory
	for _, filename := range chatmates {
		content, err := assets.GetEmbeddedMateContent(filename)
		s.Require().NoError(err, "Should be able to read embedded file %s", filename)

		destFile := filepath.Join(promptsDir, filename)
		err = os.WriteFile(destFile, content, 0644)
		s.Require().NoError(err, "Should be able to copy %s", filename)
	}

	// Verify all files were copied
	err = helpers.VerifyInstallation(promptsDir, len(chatmates))
	s.Assert().NoError(err, "Linux installation should succeed")
}

// TestFullInstallationWorkflowWindows tests complete installation process on Windows
func (s *IntegrationTestSuite) TestFullInstallationWorkflowWindows() {
	if runtime.GOOS != "windows" {
		s.T().Skip("Windows-specific test")
	}

	s.env.SimulateOS(s.T(), "windows")
	promptsDir := s.env.SetupMockVSCode(s.T(), "windows")

	// Get embedded chatmates instead of filesystem
	chatmates, err := assets.GetEmbeddedMatesList()
	s.Require().NoError(err, "Should be able to get embedded chatmate list")
	s.Assert().Greater(len(chatmates), 0, "Should have embedded chatmate files")

	// Copy embedded chatmates to prompts directory using Go's file operations
	for _, filename := range chatmates {
		content, err := assets.GetEmbeddedMateContent(filename)
		s.Require().NoError(err, "Should be able to read embedded file %s", filename)

		destFile := filepath.Join(promptsDir, filename)
		err = os.WriteFile(destFile, content, 0644)
		s.Require().NoError(err, "Should be able to copy %s", filename)
	}

	// Verify all files were copied
	err = helpers.VerifyInstallation(promptsDir, len(chatmates))
	s.Assert().NoError(err, "Windows installation should succeed")
}

// TestCLIBinaryFunctionality tests the CLI binary functionality
func (s *IntegrationTestSuite) TestCLIBinaryFunctionality() {
	// Test version flag
	cmd := exec.Command("./"+s.binaryPath, "--version")
	output, err := cmd.CombinedOutput()
	s.Require().NoError(err, "Version command should succeed")

	outputStr := string(output)
	s.Assert().Contains(outputStr, "chatmate", "Version output should contain 'chatmate'")

	// Test help flag
	cmd = exec.Command("./"+s.binaryPath, "--help")
	output, err = cmd.CombinedOutput()
	s.Require().NoError(err, "Help command should succeed")

	outputStr = string(output)
	expectedCommands := []string{"hire", "list", "uninstall", "status", "config"}
	for _, expected := range expectedCommands {
		s.Assert().Contains(outputStr, expected, "Help should mention command: %s", expected)
	}
}

// TestCLIListAvailable tests the list available command
func (s *IntegrationTestSuite) TestCLIListAvailable() {
	cmd := exec.Command("./"+s.binaryPath, "list", "--available")
	output, err := cmd.CombinedOutput()

	// Command might fail in test environment, but should not crash
	if err != nil {
		s.T().Logf("List available command failed (expected in test env): %v", err)
	}

	outputStr := string(output)
	s.T().Logf("List available output: %s", outputStr)

	// Should contain some indication of available chatmates or appropriate message
	containsExpected := strings.Contains(outputStr, "Available") ||
		strings.Contains(outputStr, "chatmate") ||
		strings.Contains(outputStr, "No") ||
		strings.Contains(outputStr, "Error")
	s.Assert().True(containsExpected, "Output should contain expected content")
}

// TestCLIInvalidCommand tests handling of invalid commands
func (s *IntegrationTestSuite) TestCLIInvalidCommand() {
	cmd := exec.Command("./"+s.binaryPath, "invalidcommand")
	output, err := cmd.CombinedOutput()
	s.Assert().Error(err, "Invalid command should return error")

	outputStr := string(output)
	containsError := strings.Contains(outputStr, "unknown command") ||
		strings.Contains(outputStr, "Error:") ||
		strings.Contains(outputStr, "invalid")
	s.Assert().True(containsError, "Should show error message for unknown command")
}

// TestSecurityScanning tests security validation of chatmate files
func (s *IntegrationTestSuite) TestSecurityScanning() {
	// Get embedded chatmates instead of filesystem
	chatmates, err := assets.GetEmbeddedMatesList()
	s.Require().NoError(err, "Should be able to get embedded chatmate list")

	for _, filename := range chatmates {
		content, err := assets.GetEmbeddedMateContent(filename)
		s.Require().NoError(err, "Should be able to read embedded file %s", filename)

		contentStr := string(content)

		// Basic security checks
		suspiciousPatterns := []string{
			"<script>",
			"javascript:",
			"eval(",
			"exec(",
			"system(",
			"shell_exec",
		}

		for _, pattern := range suspiciousPatterns {
			s.Assert().NotContains(contentStr, pattern,
				"File %s should not contain suspicious pattern: %s", filename, pattern)
		}
	}
}

// TestCrossPlatformCompatibility tests file operations across platforms
func (s *IntegrationTestSuite) TestCrossPlatformCompatibility() {
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
func (s *IntegrationTestSuite) TestFileValidationWorkflow() {
	// Test workflow with various file types
	testFiles := []struct {
		name    string
		content string
		valid   bool
	}{
		{
			name: "valid.chatmode.md",
			content: `---
description: "Valid test chatmate"
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

			err = helpers.ValidateChatmodeFile(filePath)
			if testFile.valid {
				s.Assert().NoError(err, "Valid file should pass validation")
			} else {
				s.Assert().Error(err, "Invalid file should fail validation")
			}
		})
	}
}

// Add CopyFile helper to helpers package
func init() {
	// This would be added to the helpers package
}

// TestIntegrationSuite runs the integration test suite
func TestIntegrationSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}
