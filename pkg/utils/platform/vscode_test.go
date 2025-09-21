package platform

import (
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestGetVSCodePromptsDir(t *testing.T) {
	promptsDir, err := GetVSCodePromptsDir()
	if err != nil {
		t.Fatalf("GetVSCodePromptsDir() failed: %v", err)
	}

	if promptsDir == "" {
		t.Fatal("GetVSCodePromptsDir() returned empty path")
	}

	// Verify the path contains expected OS-specific segments
	switch runtime.GOOS {
	case "darwin":
		if !contains(promptsDir, "Library") || !contains(promptsDir, "Application Support") {
			t.Errorf("macOS path doesn't contain expected segments: %s", promptsDir)
		}
	case "linux":
		if !contains(promptsDir, ".config") {
			t.Errorf("Linux path doesn't contain expected segments: %s", promptsDir)
		}
	case "windows":
		// Windows can use either APPDATA or fallback path
		if !contains(promptsDir, "AppData") && !contains(promptsDir, "Roaming") {
			t.Errorf("Windows path doesn't contain expected segments: %s", promptsDir)
		}
	}

	// Verify the path ends with Code/User/prompts
	expectedSuffix := filepath.Join("Code", "User", "prompts")
	if !endsWith(promptsDir, expectedSuffix) {
		t.Errorf("Path doesn't end with expected suffix: %s", promptsDir)
	}
}

func TestPromptsDirectoryExists(t *testing.T) {
	exists, path, err := PromptsDirectoryExists()
	if err != nil {
		t.Fatalf("PromptsDirectoryExists() failed: %v", err)
	}

	if path == "" {
		t.Fatal("PromptsDirectoryExists() returned empty path")
	}

	// The directory may or may not exist, but we should get a valid response
	t.Logf("VS Code prompts directory exists: %v at %s", exists, path)
}

// Helper functions
func contains(str, substr string) bool {
	return strings.Contains(str, substr)
}

func endsWith(str, suffix string) bool {
	return strings.HasSuffix(str, suffix)
}
