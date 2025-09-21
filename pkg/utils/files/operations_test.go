package files

import (
	"os"
	"path/filepath"
	"testing"
)

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
