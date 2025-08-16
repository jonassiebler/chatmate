package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// TestMainIntegration tests the main binary functionality
func TestMainIntegration(t *testing.T) {
	// Build the binary for testing
	binaryName := "chatmate-test"
	if runtime.GOOS == "windows" {
		binaryName += ".exe"
	}

	buildCmd := exec.Command("go", "build", "-o", binaryName, "github.com/jonassiebler/chatmate")
	err := buildCmd.Run()
	if err != nil {
		t.Fatalf("Failed to build test binary: %v", err)
	}
	defer func() { _ = os.Remove(binaryName) }()

	// Test version flag
	t.Run("version", func(t *testing.T) {
		cmd := exec.Command("./"+binaryName, "--version")
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Version command failed: %v", err)
		}

		outputStr := string(output)
		if !strings.Contains(outputStr, "chatmate") {
			t.Errorf("Version output should contain 'chatmate', got: %s", outputStr)
		}
	})

	// Test help flag
	t.Run("help", func(t *testing.T) {
		cmd := exec.Command("./"+binaryName, "--help")
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("Help command failed: %v", err)
		}

		outputStr := string(output)
		expectedCommands := []string{"hire", "list", "uninstall", "status", "config"}
		for _, expected := range expectedCommands {
			if !strings.Contains(outputStr, expected) {
				t.Errorf("Help output missing command: %s", expected)
			}
		}
	})

	// Test list available (should work without VS Code)
	t.Run("list_available", func(t *testing.T) {
		cmd := exec.Command("./"+binaryName, "list", "--available")
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("List available command failed: %v", err)
		}

		outputStr := string(output)
		// Should show available chatmates or appropriate message
		if !strings.Contains(outputStr, "Available") && !strings.Contains(outputStr, "chatmate") && !strings.Contains(outputStr, "No") {
			t.Logf("List available output: %s", outputStr)
		}
	})

	// Test invalid command
	t.Run("invalid_command", func(t *testing.T) {
		cmd := exec.Command("./"+binaryName, "invalidcommand")
		output, err := cmd.CombinedOutput()
		if err == nil {
			t.Error("Expected error for invalid command")
		}

		outputStr := string(output)
		if !strings.Contains(outputStr, "unknown command") && !strings.Contains(outputStr, "Error:") {
			t.Errorf("Expected error message for unknown command, got: %s", outputStr)
		}
	})
}

// TestMainWithMockVSCode tests with a mock VS Code environment
func TestMainWithMockVSCode(t *testing.T) {
	// Create temporary directory structure mimicking VS Code
	tmpDir, err := os.MkdirTemp("", "chatmate-integration-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	// Create mock VS Code prompts directory
	vscodeDir := filepath.Join(tmpDir, "Code", "User", "prompts")
	err = os.MkdirAll(vscodeDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create mock VS Code directory: %v", err)
	}

	// Create some test installed chatmates
	testChatmates := []string{
		"Test Agent 1.chatmode.md",
		"Test Agent 2.chatmode.md",
	}

	for _, chatmate := range testChatmates {
		content := "# " + strings.TrimSuffix(chatmate, ".chatmode.md") + "\n\nTest chatmate content"
		err := os.WriteFile(filepath.Join(vscodeDir, chatmate), []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to create test chatmate: %v", err)
		}
	}

	// Build the binary for testing
	binaryName := "chatmate-mock-test"
	if runtime.GOOS == "windows" {
		binaryName += ".exe"
	}

	buildCmd := exec.Command("go", "build", "-o", binaryName, "github.com/jonassiebler/chatmate")
	err = buildCmd.Run()
	if err != nil {
		t.Fatalf("Failed to build test binary: %v", err)
	}
	defer func() { _ = os.Remove(binaryName) }()

	// Set environment to point to our mock VS Code directory
	homeDir := filepath.Dir(filepath.Dir(vscodeDir)) // Go up to the parent of Code

	// Test status command with mock environment
	t.Run("status_with_mock", func(t *testing.T) {
		cmd := exec.Command("./"+binaryName, "status")

		// Set environment variables based on OS
		switch runtime.GOOS {
		case "darwin":
			cmd.Env = append(os.Environ(), "HOME="+homeDir)
		case "linux":
			cmd.Env = append(os.Environ(), "HOME="+homeDir)
		case "windows":
			cmd.Env = append(os.Environ(), "APPDATA="+homeDir)
		}

		output, err := cmd.CombinedOutput()
		if err != nil {
			// May fail if paths don't match exactly, but should not panic
			t.Logf("Status command output: %s", string(output))
		}
	})
}

// TestMainErrorHandling tests error handling scenarios
func TestMainErrorHandling(t *testing.T) {
	// Build the binary for testing
	binaryName := "chatmate-error-test"
	if runtime.GOOS == "windows" {
		binaryName += ".exe"
	}

	buildCmd := exec.Command("go", "build", "-o", binaryName, "github.com/jonassiebler/chatmate")
	err := buildCmd.Run()
	if err != nil {
		t.Fatalf("Failed to build test binary: %v", err)
	}
	defer func() { _ = os.Remove(binaryName) }()

	// Test hire with non-existent chatmate
	t.Run("hire_nonexistent", func(t *testing.T) {
		cmd := exec.Command("./"+binaryName, "hire", "NonExistentAgent")
		output, _ := cmd.CombinedOutput()
		// Command should succeed but show warning message

		outputStr := string(output)
		if !strings.Contains(outputStr, "No chatmate found matching") && !strings.Contains(outputStr, "not found") && !strings.Contains(outputStr, "VS Code") {
			t.Errorf("Expected appropriate warning message, got: %s", outputStr)
		}
	})

	// Test uninstall with non-existent chatmate
	t.Run("uninstall_nonexistent", func(t *testing.T) {
		cmd := exec.Command("./"+binaryName, "uninstall", "NonExistentAgent")
		output, err := cmd.CombinedOutput()
		// May not error if chatmate doesn't exist (graceful handling)
		if err != nil && !strings.Contains(string(output), "VS Code") && !strings.Contains(string(output), "prompts") {
			t.Logf("Uninstall non-existent output: %s", string(output))
		}
	})
}
