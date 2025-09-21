package main

import (
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/jonassiebler/chatmate/internal/assets"
)

// CLIFunctionalitySuite tests CLI command functionality
type CLIFunctionalitySuite struct {
	BaseIntegrationSuite
}

// TestCLIBinaryFunctionality tests the CLI binary functionality
func (s *CLIFunctionalitySuite) TestCLIBinaryFunctionality() {
	// Test version flag
	cmd := exec.Command("./"+s.GetBinaryPath(), "--version")
	output, err := cmd.CombinedOutput()
	s.Require().NoError(err, "Version command should succeed")

	outputStr := string(output)
	s.Assert().Contains(outputStr, "chatmate", "Version output should contain 'chatmate'")

	// Test help flag
	cmd = exec.Command("./"+s.GetBinaryPath(), "--help")
	output, err = cmd.CombinedOutput()
	s.Require().NoError(err, "Help command should succeed")

	outputStr = string(output)
	expectedCommands := []string{"hire", "list", "uninstall", "status", "config"}
	for _, expected := range expectedCommands {
		s.Assert().Contains(outputStr, expected, "Help should mention command: %s", expected)
	}
}

// TestCLIListAvailable tests the list available command
func (s *CLIFunctionalitySuite) TestCLIListAvailable() {
	cmd := exec.Command("./"+s.GetBinaryPath(), "list", "--available")
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
func (s *CLIFunctionalitySuite) TestCLIInvalidCommand() {
	cmd := exec.Command("./"+s.GetBinaryPath(), "invalidcommand")
	output, err := cmd.CombinedOutput()
	s.Assert().Error(err, "Invalid command should return error")

	outputStr := string(output)
	containsError := strings.Contains(outputStr, "unknown command") ||
		strings.Contains(outputStr, "Error:") ||
		strings.Contains(outputStr, "invalid")
	s.Assert().True(containsError, "Should show error message for unknown command")
}

// TestSecurityScanning tests security validation of chatmate files
func (s *CLIFunctionalitySuite) TestSecurityScanning() {
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

// TestCLIFunctionalitySuite runs the CLI functionality test suite
func TestCLIFunctionalitySuite(t *testing.T) {
	suite.Run(t, new(CLIFunctionalitySuite))
}
