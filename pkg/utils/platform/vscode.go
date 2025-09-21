// Package platform provides platform-specific path operations and VS Code integration utilities.
//
// This package handles cross-platform differences in file system layouts,
// particularly for VS Code configuration and data directories. It provides
// unified interfaces for accessing VS Code prompts directories across macOS,
// Linux, and Windows systems.
package platform

import (
	"os"
	"path/filepath"
	"runtime"
)

// GetVSCodePromptsDir returns the platform-specific path to the VS Code prompts directory.
//
// This function handles the different directory structures used by VS Code
// across operating systems:
//   - macOS: ~/Library/Application Support/Code/User/prompts
//   - Linux: ~/.config/Code/User/prompts
//   - Windows: %APPDATA%/Code/User/prompts
//
// For unknown operating systems, it defaults to the Linux-style path structure.
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
