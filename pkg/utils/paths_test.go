package utils

import (
	"testing"
)

// TestPathsIntegration ensures the re-exported functions work correctly
// and maintain backward compatibility with existing code.
func TestPathsIntegration(t *testing.T) {
	// Test VS Code platform functionality
	t.Run("VS Code prompts directory", func(t *testing.T) {
		dir, err := GetVSCodePromptsDir()
		if err != nil {
			t.Fatalf("GetVSCodePromptsDir() failed: %v", err)
		}
		if dir == "" {
			t.Error("GetVSCodePromptsDir() returned empty directory")
		}
	})

	// Test prompts directory existence checking
	t.Run("Prompts directory exists check", func(t *testing.T) {
		exists, path, err := PromptsDirectoryExists()
		if err != nil {
			t.Fatalf("PromptsDirectoryExists() failed: %v", err)
		}
		if path == "" {
			t.Error("PromptsDirectoryExists() returned empty path")
		}
		t.Logf("Prompts directory exists: %v at %s", exists, path)
	})

	// Test path expansion
	t.Run("Path expansion", func(t *testing.T) {
		result := ExpandPath("~/test")
		if result == "~/test" {
			t.Error("ExpandPath() did not expand tilde")
		}
	})

	// Test file existence checking
	t.Run("File exists", func(t *testing.T) {
		// Test with non-existent file
		if FileExists("/this/path/should/not/exist") {
			t.Error("FileExists() returned true for non-existent file")
		}
	})

	// Test chatmate file identification
	t.Run("Chatmate file identification", func(t *testing.T) {
		if !IsChatmateFile("test.chatmode.md") {
			t.Error("IsChatmateFile() failed to identify valid chatmate file")
		}
		if IsChatmateFile("test.md") {
			t.Error("IsChatmateFile() incorrectly identified non-chatmate file")
		}
	})

	// Test chatmate name extraction
	t.Run("Chatmate name extraction", func(t *testing.T) {
		name := GetChatmateNameFromFilename("Test Agent.chatmode.md")
		if name != "Test Agent" {
			t.Errorf("GetChatmateNameFromFilename() = %q, want %q", name, "Test Agent")
		}
	})
}
