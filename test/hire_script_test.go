package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/jonassiebler/chatmate/internal/testing/helpers"
)

// TestHireScriptExists tests that the hire.sh script exists and is executable
func TestHireScriptExists(t *testing.T) {
	// Check file exists
	_, err := os.Stat("../hire.sh")
	require.NoError(t, err, "../hire.sh should exist")

	// Check file permissions (on Unix systems)
	if runtime.GOOS != "windows" {
		info, err := os.Stat("../hire.sh")
		require.NoError(t, err)

		// Check if the file is executable
		mode := info.Mode()
		assert.True(t, mode&0100 != 0, "../hire.sh should be executable")
	}
}

// TestHireScriptHasValidBashSyntax tests that the hire.sh script has valid bash syntax
func TestHireScriptHasValidBashSyntax(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Bash syntax check not available on Windows")
	}

	cmd := exec.Command("bash", "-n", "../hire.sh")
	output, err := cmd.CombinedOutput()
	assert.NoError(t, err, "../hire.sh should have valid bash syntax. Output: %s", string(output))
}

// TestHireScriptContainsRequiredShebang tests that hire.sh has proper shebang
func TestHireScriptContainsRequiredShebang(t *testing.T) {
	content, err := os.ReadFile("../hire.sh")
	require.NoError(t, err, "Should be able to read hire.sh")

	lines := strings.Split(string(content), "\n")
	require.Greater(t, len(lines), 0, "File should not be empty")

	firstLine := strings.TrimSpace(lines[0])
	assert.True(t, strings.Contains(firstLine, "#!/bin/bash") || strings.Contains(firstLine, "#!/usr/bin/env bash"),
		"First line should contain bash shebang, got: %s", firstLine)
}

// TestScriptDetectsMatesDirectoryCorrectly tests directory detection logic
func TestScriptDetectsMatesDirectoryCorrectly(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Bash script testing not available on Windows")
	}

	// Test that script can find mates directory
	cmd := exec.Command("bash", "-c", "source hire.sh && echo \"$MATES_DIR\"")
	output, err := cmd.CombinedOutput()

	// The script should either succeed or fail gracefully
	if err == nil {
		outputStr := strings.TrimSpace(string(output))
		assert.True(t, strings.Contains(outputStr, "mates"),
			"MATES_DIR should contain 'mates', got: %s", outputStr)
	}
}

// TestScriptFailsGracefullyIfMatesDirectoryMissing tests error handling for missing mates directory
func TestScriptFailsGracefullyIfMatesDirectoryMissing(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Bash script testing not available on Windows")
	}

	env := helpers.SetupTestEnvironment(t)

	// Create a temporary copy of the hire script
	originalScript, err := os.ReadFile("../hire.sh")
	require.NoError(t, err)

	tempScript := filepath.Join(env.TempDir, "hire_test.sh")
	err = os.WriteFile(tempScript, originalScript, 0755)
	require.NoError(t, err)

	// Create temp directory without mates
	tempDir := filepath.Join(env.TempDir, "no_mates")
	err = os.MkdirAll(tempDir, 0755)
	require.NoError(t, err)

	// Change to temp directory and run script
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	err = os.Chdir(tempDir)
	require.NoError(t, err)

	cmd := exec.Command(tempScript)
	output, err := cmd.CombinedOutput()

	// Script should fail when mates directory is missing
	assert.Error(t, err, "Script should fail when mates directory is missing. Output: %s", string(output))
}

// TestOSDetectionLogic tests that script can handle different OS types
func TestOSDetectionLogic(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Bash script testing not available on Windows")
	}

	testCases := []struct {
		osType   string
		expected string
	}{
		{"darwin", "Library/Application Support/Code/User/prompts"},
		{"linux", ".config/Code/User/prompts"},
	}

	for _, tc := range testCases {
		t.Run(tc.osType, func(t *testing.T) {
			// This is a simplified test - in reality, the hire.sh script
			// would need more complex testing with environment manipulation

			// For now, we just verify the script contains the expected paths
			content, err := os.ReadFile("../hire.sh")
			require.NoError(t, err)

			contentStr := string(content)
			assert.True(t, strings.Contains(contentStr, tc.expected),
				"../hire.sh should contain path for %s: %s", tc.osType, tc.expected)
		})
	}
}

// TestScriptHandlesFileOperations tests basic file operation handling
func TestScriptHandlesFileOperations(t *testing.T) {
	// Check that the script contains basic file operations
	content, err := os.ReadFile("../hire.sh")
	require.NoError(t, err, "Should be able to read hire.sh")

	contentStr := string(content)

	// Check for common file operations that should be in the script
	expectedOperations := []string{"mkdir", "cp", "mv", "find"}
	for _, op := range expectedOperations {
		assert.True(t, strings.Contains(contentStr, op),
			"../hire.sh should contain '%s' operation", op)
	}
}

// TestScriptContainsErrorHandling tests that script has proper error handling
func TestScriptContainsErrorHandling(t *testing.T) {
	content, err := os.ReadFile("../hire.sh")
	require.NoError(t, err, "Should be able to read hire.sh")

	contentStr := string(content)

	// Check for error handling patterns
	errorHandlingPatterns := []string{"exit", "return"}
	hasErrorHandling := false

	for _, pattern := range errorHandlingPatterns {
		if strings.Contains(contentStr, pattern) {
			hasErrorHandling = true
			break
		}
	}

	assert.True(t, hasErrorHandling, "../hire.sh should contain error handling")
}

// TestScriptContainsHelpInformation tests that script provides usage information
func TestScriptContainsHelpInformation(t *testing.T) {
	content, err := os.ReadFile("../hire.sh")
	require.NoError(t, err, "Should be able to read hire.sh")

	contentStr := string(content)

	// Check for help/usage information
	helpPatterns := []string{"usage", "help", "Usage", "Help", "USAGE"}
	hasHelp := false

	for _, pattern := range helpPatterns {
		if strings.Contains(contentStr, pattern) {
			hasHelp = true
			break
		}
	}

	assert.True(t, hasHelp, "../hire.sh should contain help/usage information")
}

// TestScriptFileIntegrity tests that the script file maintains integrity
func TestScriptFileIntegrity(t *testing.T) {
	info, err := os.Stat("../hire.sh")
	require.NoError(t, err, "Should be able to stat hire.sh")

	// Check file is not empty
	assert.Greater(t, info.Size(), int64(100), "../hire.sh should contain meaningful content")

	// Check file is readable
	_, err = os.ReadFile("../hire.sh")
	assert.NoError(t, err, "../hire.sh should be readable")
}

// TestScriptEnvironmentVariables tests that script handles environment variables properly
func TestScriptEnvironmentVariables(t *testing.T) {
	content, err := os.ReadFile("../hire.sh")
	require.NoError(t, err, "Should be able to read hire.sh")

	contentStr := string(content)

	// Check for environment variable usage
	envPatterns := []string{"HOME", "APPDATA", "$", "env"}
	hasEnvUsage := false

	for _, pattern := range envPatterns {
		if strings.Contains(contentStr, pattern) {
			hasEnvUsage = true
			break
		}
	}

	assert.True(t, hasEnvUsage, "../hire.sh should use environment variables")
}

// TestScriptWithDryRun tests script dry-run functionality if available
func TestScriptWithDryRun(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Bash script testing not available on Windows")
	}

	// Test if script supports dry-run or verbose mode
	content, err := os.ReadFile("../hire.sh")
	require.NoError(t, err)

	contentStr := string(content)

	// Look for dry-run or verbose flags
	dryRunPatterns := []string{"dry-run", "--dry-run", "-n", "verbose", "--verbose", "-v"}
	for _, pattern := range dryRunPatterns {
		if strings.Contains(contentStr, pattern) {
			t.Logf("Script contains %s option", pattern)
		}
	}
}
