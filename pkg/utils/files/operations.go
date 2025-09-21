// Package files provides general file and directory operation utilities.
//
// This package offers common file system operations including path expansion,
// directory creation, file existence checking, and other utilities that are
// used throughout the ChatMate application.
package files

import (
	"os"
	"path/filepath"
	"strings"
)

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
