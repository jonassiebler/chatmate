// Package utils provides cross-platform file system utilities and path manipulation
// functions for the chatmate CLI tool.
//
// This package handles the complexity of working with different operating systems
// and their file system conventions, particularly for VS Code configuration
// directories and chatmate file management.
//
// Key Features:
//   - Cross-platform VS Code prompts directory detection
//   - Safe directory creation with proper permissions
//   - Path expansion and validation utilities
//   - Chatmate file identification and naming
//
// Example usage:
//
//	// Get VS Code prompts directory for current OS
//	promptsDir, err := utils.GetVSCodePromptsDir()
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Ensure directory exists with proper permissions
//	dir, err := utils.EnsurePromptsDir()
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Check if a file is a valid chatmate file
//	if utils.IsChatmateFile("my-agent.chatmode.md") {
//		name := utils.GetChatmateNameFromFilename("my-agent.chatmode.md")
//		fmt.Println("Chatmate name:", name) // Output: my-agent
//	}
package utils

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// GetVSCodePromptsDir returns the VS Code prompts directory path for the current OS.
//
// The function automatically detects the operating system and returns the appropriate
// VS Code prompts directory path following each platform's conventions:
//   - macOS: ~/Library/Application Support/Code/User/prompts
//   - Linux: ~/.config/Code/User/prompts
//   - Windows: %APPDATA%/Code/User/prompts
//
// For unknown operating systems, it defaults to the Linux-style path.
//
// Example:
//
//	promptsDir, err := GetVSCodePromptsDir()
//	if err != nil {
//		return fmt.Errorf("failed to get prompts directory: %w", err)
//	}
//	fmt.Printf("VS Code prompts directory: %s\n", promptsDir)
//
// Returns:
//   - string: The full path to the VS Code prompts directory
//   - error: Any error encountered while determining the home directory
func GetVSCodePromptsDir() (string, error) {
	var homeDir string
	var err error

	// Get user home directory
	homeDir, err = os.UserHomeDir()
	if err != nil {
		return "", err
	}

	var promptsDir string

	switch runtime.GOOS {
	case "darwin": // macOS
		promptsDir = filepath.Join(homeDir, "Library", "Application Support", "Code", "User", "prompts")
	case "linux":
		promptsDir = filepath.Join(homeDir, ".config", "Code", "User", "prompts")
	case "windows":
		// Windows uses %APPDATA%/Code/User/prompts
		appData := os.Getenv("APPDATA")
		if appData == "" {
			// Fallback to default location
			appData = filepath.Join(homeDir, "AppData", "Roaming")
		}
		promptsDir = filepath.Join(appData, "Code", "User", "prompts")
	default:
		// Default to Linux-style path for unknown OS
		promptsDir = filepath.Join(homeDir, ".config", "Code", "User", "prompts")
	}

	return promptsDir, nil
}

// EnsurePromptsDir creates the VS Code prompts directory if it doesn't exist.
//
// This function combines directory path detection with safe creation, ensuring
// that the full directory tree exists with proper permissions (0755). It's
// idempotent - safe to call multiple times without side effects.
//
// The directory is created with read, write, and execute permissions for the
// owner, and read and execute permissions for group and others.
//
// Example:
//
//	promptsDir, err := EnsurePromptsDir()
//	if err != nil {
//		return fmt.Errorf("failed to ensure prompts directory: %w", err)
//	}
//	fmt.Printf("Prompts directory ready: %s\n", promptsDir)
//
// Returns:
//   - string: The full path to the VS Code prompts directory
//   - error: Any error encountered during path detection or directory creation
func EnsurePromptsDir() (string, error) {
	promptsDir, err := GetVSCodePromptsDir()
	if err != nil {
		return "", err
	}

	err = os.MkdirAll(promptsDir, 0755)
	if err != nil {
		return "", err
	}

	return promptsDir, nil
}

// PromptsDirectoryExists checks if the VS Code prompts directory exists.
//
// This function provides detailed information about the prompts directory
// existence and validity. It distinguishes between non-existent paths,
// files that exist but aren't directories, and actual directories.
//
// Example:
//
//	exists, path, err := PromptsDirectoryExists()
//	if err != nil {
//		return fmt.Errorf("failed to check prompts directory: %w", err)
//	}
//	if !exists {
//		fmt.Printf("Prompts directory does not exist: %s\n", path)
//	} else {
//		fmt.Printf("Prompts directory exists: %s\n", path)
//	}
//
// Returns:
//   - bool: true if the directory exists and is actually a directory
//   - string: the full path to the prompts directory (regardless of existence)
//   - error: any error encountered during path detection or stat operation
func PromptsDirectoryExists() (bool, string, error) {
	promptsDir, err := GetVSCodePromptsDir()
	if err != nil {
		return false, "", err
	}

	info, err := os.Stat(promptsDir)
	if os.IsNotExist(err) {
		return false, promptsDir, nil
	}
	if err != nil {
		return false, promptsDir, err
	}

	return info.IsDir(), promptsDir, nil
}

// ExpandPath expands ~ to the user's home directory in file paths.
//
// This function handles common shell-style path expansion, converting:
//   - "~" to the user's home directory
//   - "~/..." to the user's home directory plus the relative path
//   - Other paths are returned unchanged
//
// If the home directory cannot be determined, the tilde expansion is skipped
// and the original path is returned.
//
// Example:
//
//	expanded := ExpandPath("~/Documents/my-file.txt")
//	// On macOS: /Users/username/Documents/my-file.txt
//	// On Linux: /home/username/Documents/my-file.txt
//	// On Windows: C:\Users\username\Documents\my-file.txt
//
// Parameters:
//   - path: the file path to expand, may contain ~ prefix
//
// Returns:
//   - string: the expanded path with ~ replaced by home directory
func ExpandPath(path string) string {
	if path == "" {
		return ""
	}

	if path == "~" {
		home, _ := os.UserHomeDir()
		return home
	}

	if strings.HasPrefix(path, "~/") {
		home, _ := os.UserHomeDir()
		return filepath.Join(home, path[2:])
	}

	return path
}

// EnsureDir creates a directory and all necessary parent directories.
//
// This function is equivalent to `mkdir -p` in Unix systems, creating the
// entire directory path if it doesn't exist. The directory is created with
// permissions 0755 (read, write, execute for owner; read, execute for others).
//
// The function is idempotent - it won't fail if the directory already exists.
//
// Example:
//
//	err := EnsureDir("/path/to/deeply/nested/directory")
//	if err != nil {
//		return fmt.Errorf("failed to create directory: %w", err)
//	}
//
// Parameters:
//   - dir: the directory path to create, including any missing parent directories
//
// Returns:
//   - error: any error encountered during directory creation
func EnsureDir(dir string) error {
	return os.MkdirAll(dir, 0755)
}

// FileExists checks if a file exists and is not a directory.
//
// This function provides a reliable way to check for file existence while
// explicitly excluding directories. It returns false for non-existent paths,
// directories, and when stat operations fail.
//
// Example:
//
//	if FileExists("config.json") {
//		fmt.Println("Configuration file found")
//	} else {
//		fmt.Println("Configuration file missing")
//	}
//
// Parameters:
//   - filename: the path to check for file existence
//
// Returns:
//   - bool: true if the path exists and is a regular file (not a directory)
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// GetChatmateNameFromFilename extracts the chatmate name from a filename.
//
// This function removes the standard ".chatmode.md" extension from chatmate
// files to get the base chatmate name. If the file doesn't have the expected
// extension, the original filename is returned unchanged.
//
// Example:
//
//	name := GetChatmateNameFromFilename("Chatmate - Code Claude Sonnet 4.chatmode.md")
//	fmt.Println(name) // Output: Chatmate - Code Claude Sonnet 4
//
//	name = GetChatmateNameFromFilename("regular-file.txt")
//	fmt.Println(name) // Output: regular-file.txt
//
// Parameters:
//   - filename: the filename to process, may include .chatmode.md extension
//
// Returns:
//   - string: the chatmate name with .chatmode.md extension removed, or original filename
func GetChatmateNameFromFilename(filename string) string {
	// Remove .chatmode.md extension if present
	if strings.HasSuffix(filename, ".chatmode.md") {
		return strings.TrimSuffix(filename, ".chatmode.md")
	}
	return filename
}

// IsChatmateFile checks if a filename is a valid chatmate file.
//
// This function determines if a filename follows the chatmate naming convention
// by checking for the ".chatmode.md" extension. This is used throughout the
// application to identify and filter chatmate files.
//
// Example:
//
//	if IsChatmateFile("Chatmate - Code Claude Sonnet 4.chatmode.md") {
//		fmt.Println("This is a chatmate file")
//	}
//
//	if !IsChatmateFile("README.md") {
//		fmt.Println("This is not a chatmate file")
//	}
//
// Parameters:
//   - filename: the filename to check for chatmate file convention
//
// Returns:
//   - bool: true if the filename ends with ".chatmode.md"
func IsChatmateFile(filename string) bool {
	return strings.HasSuffix(filename, ".chatmode.md")
}
