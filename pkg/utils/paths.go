// Package utils provides path utilities and file operations for the ChatMate application.
//
// This package serves as the main entry point for path-related operations,
// re-exporting functionality from specialized sub-packages while maintaining
// backward compatibility with existing code.
//
// The functionality is organized into specialized packages:
//   - platform: VS Code and platform-specific path operations
//   - files: General file operations and chatmate-specific utilities
package utils

import (
	"github.com/jonassiebler/chatmate/pkg/utils/files"
	"github.com/jonassiebler/chatmate/pkg/utils/platform"
)

// VS Code and platform-specific operations (re-exported from platform package)

// GetVSCodePromptsDir returns the platform-specific path to the VS Code prompts directory.
//
// This function handles the different directory structures used by VS Code
// across operating systems. See platform.GetVSCodePromptsDir for details.
func GetVSCodePromptsDir() (string, error) {
	return platform.GetVSCodePromptsDir()
}

// EnsurePromptsDir creates the VS Code prompts directory if it doesn't exist.
//
// This function combines directory path detection with safe creation.
// See platform.EnsurePromptsDir for details.
func EnsurePromptsDir() (string, error) {
	return platform.EnsurePromptsDir()
}

// PromptsDirectoryExists checks if the VS Code prompts directory exists.
//
// This function provides detailed information about the prompts directory
// existence and validity. See platform.PromptsDirectoryExists for details.
func PromptsDirectoryExists() (bool, string, error) {
	return platform.PromptsDirectoryExists()
}

// General file operations (re-exported from files package)

// ExpandPath expands ~ to the user's home directory in file paths.
//
// This function handles common shell-style path expansion.
// See files.ExpandPath for details.
func ExpandPath(path string) string {
	return files.ExpandPath(path)
}

// EnsureDir creates a directory and all necessary parent directories.
//
// This function is equivalent to `mkdir -p` in Unix systems.
// See files.EnsureDir for details.
func EnsureDir(dir string) error {
	return files.EnsureDir(dir)
}

// FileExists checks if a file exists and is not a directory.
//
// This function provides a reliable way to check for file existence.
// See files.FileExists for details.
func FileExists(filename string) bool {
	return files.FileExists(filename)
}

// Chatmate-specific file operations (re-exported from files package)

// GetChatmateNameFromFilename extracts the chatmate name from a filename.
//
// This function removes the standard ".chatmode.md" extension from chatmate files.
// See files.GetChatmateNameFromFilename for details.
func GetChatmateNameFromFilename(filename string) string {
	return files.GetChatmateNameFromFilename(filename)
}

// IsChatmateFile checks if a filename is a valid chatmate file.
//
// This function determines if a filename follows the chatmate naming convention.
// See files.IsChatmateFile for details.
func IsChatmateFile(filename string) bool {
	return files.IsChatmateFile(filename)
}
