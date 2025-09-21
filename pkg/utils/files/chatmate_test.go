package files

import (
	"testing"
)

// TestGetChatmateNameFromFilename tests chatmate name extraction
func TestGetChatmateNameFromFilename(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "standard chatmode file",
			input:    "Chatmate - Code Claude Sonnet 4.chatmode.md",
			expected: "Chatmate - Code Claude Sonnet 4",
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
