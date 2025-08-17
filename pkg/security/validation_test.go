package security

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestValidateFilename(t *testing.T) {
	tests := []struct {
		filename    string
		expectError bool
		errorCode   string
		description string
	}{
		// Valid filenames
		{"test.txt", false, "", "simple valid filename"},
		{"My File.doc", false, "", "filename with spaces"},
		{"test-file_123.chatmode.md", false, "", "complex valid filename"},
		{"file.with.dots.txt", false, "", "filename with multiple dots"},

		// Invalid filenames
		{"", true, "EMPTY_FILENAME", "empty filename"},
		{"file\x00.txt", true, "NULL_BYTES", "null byte injection"},
		{strings.Repeat("a", 256), true, "FILENAME_TOO_LONG", "too long filename"},
		{"file<test>.txt", true, "INVALID_CHARACTERS", "angle brackets"},
		{"file|pipe.txt", true, "INVALID_CHARACTERS", "pipe character"},
		{"CON", true, "RESERVED_NAME", "Windows reserved name CON"},
		{"PRN.txt", true, "RESERVED_NAME", "Windows reserved name PRN"},
		{"NUL.doc", true, "RESERVED_NAME", "Windows reserved name NUL"},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			err := ValidateFilename(tt.filename)

			if tt.expectError && err == nil {
				t.Errorf("Expected error for filename %q, but got none", tt.filename)
			}

			if !tt.expectError && err != nil {
				t.Errorf("Expected no error for filename %q, but got: %v", tt.filename, err)
			}

			if tt.expectError && err != nil {
				if valErr, ok := err.(ValidationError); ok {
					if valErr.Code != tt.errorCode {
						t.Errorf("Expected error code %q, got %q", tt.errorCode, valErr.Code)
					}
				} else {
					t.Errorf("Expected ValidationError, got %T", err)
				}
			}
		})
	}
}

func TestValidateChatmateFilename(t *testing.T) {
	tests := []struct {
		filename    string
		expectError bool
		errorCode   string
		description string
	}{
		// Valid chatmate filenames
		{"Test Agent.chatmode.md", false, "", "valid chatmate filename"},
		{"Code-Helper_v2.chatmode.md", false, "", "complex valid chatmate filename"},

		// Invalid chatmate filenames
		{"test.txt", true, "INVALID_CHATMATE_FILENAME", "wrong extension"},
		{"test.md", true, "INVALID_CHATMATE_FILENAME", "missing .chatmode"},
		{"test.chatmode", true, "INVALID_CHATMATE_FILENAME", "missing .md"},
		{"", true, "EMPTY_FILENAME", "empty filename"},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			err := ValidateChatmateFilename(tt.filename)

			if tt.expectError && err == nil {
				t.Errorf("Expected error for filename %q, but got none", tt.filename)
			}

			if !tt.expectError && err != nil {
				t.Errorf("Expected no error for filename %q, but got: %v", tt.filename, err)
			}

			if tt.expectError && err != nil {
				if valErr, ok := err.(ValidationError); ok {
					if valErr.Code != tt.errorCode {
						t.Errorf("Expected error code %q, got %q", tt.errorCode, valErr.Code)
					}
				}
			}
		})
	}
}

func TestValidatePath(t *testing.T) {
	tests := []struct {
		path        string
		expectError bool
		errorCode   string
		description string
	}{
		// Valid paths
		{"subdir/file.txt", false, "", "simple subdirectory path"},
		{"folder/subfolder/file.md", false, "", "nested directory path"},

		// Invalid paths
		{"", true, "EMPTY_PATH", "empty path"},
		{"path\x00/file.txt", true, "NULL_BYTES", "null byte injection"},
		{"../../../etc/passwd", true, "DIRECTORY_TRAVERSAL", "directory traversal"},
		{"dir/../../../file.txt", true, "DIRECTORY_TRAVERSAL", "directory traversal in middle"},
		{"/absolute/path/file.txt", true, "ABSOLUTE_PATH", "absolute path"},
		{"path/with$variable", true, "DANGEROUS_CHARACTER", "variable expansion"},
		{"path/with`command`", true, "DANGEROUS_CHARACTER", "command substitution"},
		{"path/with|pipe", true, "DANGEROUS_CHARACTER", "pipe character"},
		{"path/with;command", true, "DANGEROUS_CHARACTER", "command separator"},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			err := ValidatePath(tt.path)

			if tt.expectError && err == nil {
				t.Errorf("Expected error for path %q, but got none", tt.path)
			}

			if !tt.expectError && err != nil {
				t.Errorf("Expected no error for path %q, but got: %v", tt.path, err)
			}

			if tt.expectError && err != nil {
				if valErr, ok := err.(ValidationError); ok {
					if valErr.Code != tt.errorCode {
						t.Errorf("Expected error code %q, got %q", tt.errorCode, valErr.Code)
					}
				}
			}
		})
	}
}

func TestSanitizeInput(t *testing.T) {
	tests := []struct {
		input       string
		expected    string
		description string
	}{
		{"normal text", "normal text", "normal text unchanged"},
		{"text with\x00null bytes", "text withnull bytes", "null bytes removed"},
		{"  padded text  ", "padded text", "whitespace trimmed"},
		{"text\x01with\x02control\x03chars", "textwithcontrolchars", "control characters removed"},
		{"text\nwith\nnewlines", "text\nwith\nnewlines", "newlines preserved"},
		{"text\twith\ttabs", "text\twith\ttabs", "tabs preserved"},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			result := SanitizeInput(tt.input)
			if result != tt.expected {
				t.Errorf("SanitizeInput(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestIsPathSafe(t *testing.T) {
	basePath := "/home/user/chatmate"

	tests := []struct {
		targetPath  string
		expectSafe  bool
		description string
	}{
		// Safe paths
		{"mates/test.md", true, "relative path within base"},
		{"mates/subdir/file.txt", true, "nested path within base"},
		{filepath.Join(basePath, "mates/file.txt"), true, "absolute path within base"},

		// Unsafe paths
		{"../../../etc/passwd", false, "directory traversal escape"},
		{"mates/../../../sensitive", false, "traversal through subdirectory"},
		{"/etc/passwd", false, "absolute path outside base"},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			result := IsPathSafe(basePath, tt.targetPath)
			if result != tt.expectSafe {
				t.Errorf("IsPathSafe(%q, %q) = %v, want %v", basePath, tt.targetPath, result, tt.expectSafe)
			}
		})
	}
}

func TestValidateContentLength(t *testing.T) {
	tests := []struct {
		content     []byte
		maxSize     int64
		expectError bool
		description string
	}{
		{[]byte("small content"), 1024, false, "content within limits"},
		{[]byte(""), 1024, false, "empty content"},
		{[]byte(strings.Repeat("a", 2000)), 1024, true, "content exceeds limit"},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			err := ValidateContentLength(tt.content, tt.maxSize)

			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}

			if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

func TestValidateFileExtension(t *testing.T) {
	allowedExtensions := []string{".md", ".txt", ".json"}

	tests := []struct {
		filename    string
		expectError bool
		description string
	}{
		{"file.md", false, "allowed markdown extension"},
		{"file.txt", false, "allowed text extension"},
		{"file.json", false, "allowed JSON extension"},
		{"File.MD", false, "case insensitive extension"},
		{"file.exe", true, "disallowed executable extension"},
		{"file.sh", true, "disallowed script extension"},
		{"file", true, "no extension"},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			err := ValidateFileExtension(tt.filename, allowedExtensions)

			if tt.expectError && err == nil {
				t.Errorf("Expected error for filename %q, but got none", tt.filename)
			}

			if !tt.expectError && err != nil {
				t.Errorf("Expected no error for filename %q, but got: %v", tt.filename, err)
			}
		})
	}
}

// Benchmark tests for performance
func BenchmarkValidateFilename(b *testing.B) {
	filename := "test-file_123.chatmode.md"
	for i := 0; i < b.N; i++ {
		ValidateFilename(filename)
	}
}

func BenchmarkValidatePath(b *testing.B) {
	path := "mates/subfolder/file.md"
	for i := 0; i < b.N; i++ {
		ValidatePath(path)
	}
}

func BenchmarkSanitizeInput(b *testing.B) {
	input := "some input with\x00null bytes and\x01control chars"
	for i := 0; i < b.N; i++ {
		SanitizeInput(input)
	}
}
