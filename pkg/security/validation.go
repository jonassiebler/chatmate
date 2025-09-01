// Package security provides input validation and security utilities for ChatMate.
//
// This package implements comprehensive security measures including:
//   - File name and path validation to prevent directory traversal attacks
//   - Input sanitization to remove malicious content
//   - Content validation for file operations
//   - Safe path operations within restricted directories
//
// Key Security Features:
//   - Path traversal prevention (../ and absolute path protection)
//   - Filename sanitization (removes special characters and null bytes)
//   - Content length validation to prevent resource exhaustion
//   - File extension validation for allowed file types
//   - Input sanitization removing control characters
//
// Usage Example:
//
//	if err := security.ValidateChatmateFilename("MyAgent.chatmode.md"); err != nil {
//	    return fmt.Errorf("invalid filename: %w", err)
//	}
//
//	if !security.IsPathSafe("/base/dir", "/base/dir/subdir/file.txt") {
//	    return errors.New("unsafe path detected")
//	}
//
//	clean := security.SanitizeInput(userInput)
package security

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

// Input validation patterns
var (
	// SafeFilenamePattern allows alphanumeric, spaces, dots, hyphens, and underscores
	SafeFilenamePattern = regexp.MustCompile(`^[a-zA-Z0-9\s\.\-_]+$`)

	// ChatmateFilenamePattern specifically for .chatmode.md files
	ChatmateFilenamePattern = regexp.MustCompile(`^[a-zA-Z0-9\s\-_.]+\.chatmode\.md$`)

	// SafePathPattern prevents directory traversal attacks
	SafePathPattern = regexp.MustCompile(`^[^<>:"|?*]+$`)
)

// ValidationError represents a security validation error
type ValidationError struct {
	Field  string
	Value  string
	Reason string
	Code   string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("security validation failed for %s: %s (code: %s)", e.Field, e.Reason, e.Code)
}

// ValidateFilename validates that a filename is safe and doesn't contain malicious patterns
func ValidateFilename(filename string) error {
	if filename == "" {
		return ValidationError{
			Field:  "filename",
			Value:  filename,
			Reason: "filename cannot be empty",
			Code:   "EMPTY_FILENAME",
		}
	}

	// Check for null bytes
	if strings.Contains(filename, "\x00") {
		return ValidationError{
			Field:  "filename",
			Value:  filename,
			Reason: "filename contains null bytes",
			Code:   "NULL_BYTES",
		}
	}

	// Check length limits (reasonable filename length)
	if len(filename) > 255 {
		return ValidationError{
			Field:  "filename",
			Value:  filename,
			Reason: "filename too long (max 255 characters)",
			Code:   "FILENAME_TOO_LONG",
		}
	}

	// Check for dangerous patterns
	if !SafeFilenamePattern.MatchString(filename) {
		return ValidationError{
			Field:  "filename",
			Value:  filename,
			Reason: "filename contains invalid characters",
			Code:   "INVALID_CHARACTERS",
		}
	}

	// Check for reserved names (Windows)
	reserved := []string{"CON", "PRN", "AUX", "NUL", "COM1", "COM2", "COM3", "COM4",
		"COM5", "COM6", "COM7", "COM8", "COM9", "LPT1", "LPT2", "LPT3", "LPT4",
		"LPT5", "LPT6", "LPT7", "LPT8", "LPT9"}

	base := strings.ToUpper(strings.TrimSuffix(filename, filepath.Ext(filename)))
	for _, res := range reserved {
		if base == res {
			return ValidationError{
				Field:  "filename",
				Value:  filename,
				Reason: "filename uses reserved system name",
				Code:   "RESERVED_NAME",
			}
		}
	}

	return nil
}

// ValidateChatmateFilename validates specifically .chatmode.md filenames
func ValidateChatmateFilename(filename string) error {
	if err := ValidateFilename(filename); err != nil {
		return err
	}

	if !ChatmateFilenamePattern.MatchString(filename) {
		return ValidationError{
			Field:  "filename",
			Value:  filename,
			Reason: "not a valid chatmate filename (must end with .chatmode.md)",
			Code:   "INVALID_CHATMATE_FILENAME",
		}
	}

	return nil
}

// ValidatePath validates that a path is safe and doesn't contain traversal attempts
func ValidatePath(path string) error {
	if path == "" {
		return ValidationError{
			Field:  "path",
			Value:  path,
			Reason: "path cannot be empty",
			Code:   "EMPTY_PATH",
		}
	}

	// Check for null bytes
	if strings.Contains(path, "\x00") {
		return ValidationError{
			Field:  "path",
			Value:  path,
			Reason: "path contains null bytes",
			Code:   "NULL_BYTES",
		}
	}

	// Check for directory traversal patterns
	if strings.Contains(path, "..") {
		return ValidationError{
			Field:  "path",
			Value:  path,
			Reason: "path contains directory traversal patterns",
			Code:   "DIRECTORY_TRAVERSAL",
		}
	}

	// Check for absolute paths that might escape intended directories
	if filepath.IsAbs(path) {
		return ValidationError{
			Field:  "path",
			Value:  path,
			Reason: "absolute paths not allowed",
			Code:   "ABSOLUTE_PATH",
		}
	}

	// Additional safety check for common dangerous patterns
	dangerous := []string{"~", "$", "`", "|", "&", ";", "(", ")", "{", "}", "[", "]"}
	for _, pattern := range dangerous {
		if strings.Contains(path, pattern) {
			return ValidationError{
				Field:  "path",
				Value:  path,
				Reason: fmt.Sprintf("path contains dangerous character: %s", pattern),
				Code:   "DANGEROUS_CHARACTER",
			}
		}
	}

	return nil
}

// SanitizeInput removes or escapes potentially dangerous characters from user input
func SanitizeInput(input string) string {
	// Remove null bytes
	input = strings.ReplaceAll(input, "\x00", "")

	// Trim whitespace
	input = strings.TrimSpace(input)

	// Remove or replace control characters (except newlines and tabs where appropriate)
	sanitized := ""
	for _, r := range input {
		if r >= 32 || r == '\n' || r == '\t' {
			sanitized += string(r)
		}
	}

	return sanitized
}

// IsPathSafe checks if a path is within expected bounds (doesn't escape the base directory)
func IsPathSafe(basePath, targetPath string) bool {
	// Clean the paths
	cleanBase := filepath.Clean(basePath)
	cleanTarget := filepath.Clean(targetPath)

	// Make target relative to base if it's absolute
	if filepath.IsAbs(cleanTarget) {
		rel, err := filepath.Rel(cleanBase, cleanTarget)
		if err != nil {
			return false
		}
		cleanTarget = rel
	}

	// Check if the resolved path starts with ".." which would escape the base
	return !strings.HasPrefix(cleanTarget, "..") && !strings.Contains(cleanTarget, "/..")
}

// ValidateContentLength validates that content doesn't exceed reasonable limits
func ValidateContentLength(content []byte, maxSize int64) error {
	if int64(len(content)) > maxSize {
		return ValidationError{
			Field:  "content",
			Value:  fmt.Sprintf("%d bytes", len(content)),
			Reason: fmt.Sprintf("content too large (max %d bytes)", maxSize),
			Code:   "CONTENT_TOO_LARGE",
		}
	}
	return nil
}

// ValidateFileExtension validates that a file has an allowed extension
func ValidateFileExtension(filename string, allowedExtensions []string) error {
	ext := strings.ToLower(filepath.Ext(filename))

	for _, allowed := range allowedExtensions {
		if strings.ToLower(allowed) == ext {
			return nil
		}
	}

	return ValidationError{
		Field:  "filename",
		Value:  filename,
		Reason: fmt.Sprintf("file extension not allowed (allowed: %v)", allowedExtensions),
		Code:   "INVALID_EXTENSION",
	}
}
