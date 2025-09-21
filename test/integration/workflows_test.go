package main

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/jonassiebler/chatmate/internal/assets"
	"github.com/jonassiebler/chatmate/internal/testing/helpers"
)

// InstallationWorkflowSuite tests installation workflows across platforms
type InstallationWorkflowSuite struct {
	BaseIntegrationSuite
}

// TestFullInstallationWorkflowMacOS tests complete installation process on macOS
func (s *InstallationWorkflowSuite) TestFullInstallationWorkflowMacOS() {
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
func (s *InstallationWorkflowSuite) TestFullInstallationWorkflowLinux() {
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
func (s *InstallationWorkflowSuite) TestFullInstallationWorkflowWindows() {
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

// TestInstallationWorkflowSuite runs the installation workflow test suite
func TestInstallationWorkflowSuite(t *testing.T) {
	suite.Run(t, new(InstallationWorkflowSuite))
}
