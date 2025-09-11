package utils

import (
	"os"
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

// TestExpandPath tests path expansion functionality
func TestExpandPath(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected func() string
	}{
		{
			name:  "tilde expansion",
			input: "~/test",
			expected: func() string {
				home, _ := os.UserHomeDir()
				return filepath.Join(home, "test")
			},
		},
		{
			name:  "absolute path unchanged",
			input: "/absolute/path",
			expected: func() string {
				return "/absolute/path"
			},
		},
		{
			name:  "relative path unchanged",
			input: "relative/path",
			expected: func() string {
				return "relative/path"
			},
		},
		{
			name:  "tilde only",
			input: "~",
			expected: func() string {
				home, _ := os.UserHomeDir()
				return home
			},
		},
		{
			name:  "empty path",
			input: "",
			expected: func() string {
				return ""
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExpandPath(tt.input)
			expected := tt.expected()

			if result != expected {
				t.Errorf("ExpandPath(%q) = %q, want %q", tt.input, result, expected)
			}
		})
	}
}

// TestEnsureDir tests directory creation functionality
func TestEnsureDir(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := os.TempDir()
	testDir := filepath.Join(tmpDir, "chatmate-test-ensuredir")

	// Clean up after test
	defer func() { _ = os.RemoveAll(testDir) }()

	// Test creating new directory
	err := EnsureDir(testDir)
	if err != nil {
		t.Fatalf("EnsureDir() failed: %v", err)
	}

	// Verify directory was created
	if _, err := os.Stat(testDir); os.IsNotExist(err) {
		t.Error("EnsureDir() did not create directory")
	}

	// Test with existing directory (should not error)
	err = EnsureDir(testDir)
	if err != nil {
		t.Fatalf("EnsureDir() failed on existing directory: %v", err)
	}

	// Test with nested directory creation
	nestedDir := filepath.Join(testDir, "nested", "deep", "path")
	err = EnsureDir(nestedDir)
	if err != nil {
		t.Fatalf("EnsureDir() failed on nested directory: %v", err)
	}

	// Verify nested directory was created
	if _, err := os.Stat(nestedDir); os.IsNotExist(err) {
		t.Error("EnsureDir() did not create nested directory")
	}
}

// TestFileExists tests file existence checking
func TestFileExists(t *testing.T) {
	// Create a temporary file for testing
	tmpFile := filepath.Join(os.TempDir(), "chatmate-test-fileexists.txt")

	// Clean up after test
	defer func() { _ = os.Remove(tmpFile) }()

	// Test non-existent file
	if FileExists(tmpFile) {
		t.Error("FileExists() returned true for non-existent file")
	}

	// Create the file
	file, err := os.Create(tmpFile)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	_ = file.Close()

	// Test existing file
	if !FileExists(tmpFile) {
		t.Error("FileExists() returned false for existing file")
	}

	// Test with directory (should return false)
	tmpDir := filepath.Join(os.TempDir(), "chatmate-test-dir")
	_ = os.Mkdir(tmpDir, 0755)
	defer func() { _ = os.RemoveAll(tmpDir) }()

	if FileExists(tmpDir) {
		t.Error("FileExists() returned true for directory")
	}
}

// TestGetChatmateNameFromFilename tests chatmate name extraction
func TestGetChatmateNameFromFilename(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "standard chatmode file",
			input:    "Chatmate: Code Claude Sonnet 4.chatmode.md",
			expected: "Chatmate: Code Claude Sonnet 4",
		},
		{
			name:     "file with spaces",
			input:    "My Custom Agent.chatmode.md",
			expected: "My Custom Agent",
		},
		{
			name:     "file without extension",
			input:    "Test Agent",
			expected: "Test Agent",
		},
		{
			name:     "file with different extension",
			input:    "Test Agent.md",
			expected: "Test Agent.md",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "just extension",
			input:    ".chatmode.md",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetChatmateNameFromFilename(tt.input)
			if result != tt.expected {
				t.Errorf("GetChatmateNameFromFilename(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

// TestIsChatmateFile tests chatmate file validation
func TestIsChatmateFile(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "valid chatmode file",
			input:    "Test Agent.chatmode.md",
			expected: true,
		},
		{
			name:     "uppercase extension",
			input:    "Test Agent.CHATMODE.MD",
			expected: false, // assuming case-sensitive
		},
		{
			name:     "missing .md",
			input:    "Test Agent.chatmode",
			expected: false,
		},
		{
			name:     "missing .chatmode",
			input:    "Test Agent.md",
			expected: false,
		},
		{
			name:     "different extension",
			input:    "Test Agent.txt",
			expected: false,
		},
		{
			name:     "empty string",
			input:    "",
			expected: false,
		},
		{
			name:     "directory-like path",
			input:    "path/to/Test Agent.chatmode.md",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsChatmateFile(tt.input)
			if result != tt.expected {
				t.Errorf("IsChatmateFile(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

// Helper functions
func contains(str, substr string) bool {
	return strings.Contains(str, substr)
}

func endsWith(str, suffix string) bool {
	return strings.HasSuffix(str, suffix)
}
