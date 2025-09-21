// Package manager provides validation functionality for ChatMate agents.
package manager

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jonassiebler/chatmate/pkg/security"
)

// ValidatorService handles chatmate validation operations.
type ValidatorService struct {
	manager *ChatMateManager
}

// NewValidatorService creates a new validator service.
func NewValidatorService(manager *ChatMateManager) *ValidatorService {
	return &ValidatorService{manager: manager}
}

// ValidateInstallation performs comprehensive validation of chatmate installation.
//
// This method checks the integrity of the installation environment, validates
// all installed chatmates, and provides detailed feedback about any issues found.
//
// Returns:
//   - error: Validation failure or system error
//   - bool: True if validation passes, false if issues found
//
// Example:
//
// valid, err := validator.ValidateInstallation()
//
//	if err != nil {
//	   return fmt.Errorf("validation failed: %w", err)
//	}
//
//	if !valid {
//	   fmt.Println("Installation has issues")
//	}
func (v *ValidatorService) ValidateInstallation() (bool, error) {
	fmt.Println("Validating ChatMate installation...")

	// Check prompts directory
	if err := v.validatePromptsDirectory(); err != nil {
		return false, err
	}

	// Check available chatmates
	if err := v.validateAvailableChatmates(); err != nil {
		return false, err
	}

	// Check installed chatmates
	if err := v.validateInstalledChatmates(); err != nil {
		return false, err
	}

	// Check for orphaned files
	if err := v.validateOrphanedFiles(); err != nil {
		return false, err
	}

	fmt.Println("✅ ChatMate installation validation completed successfully")
	return true, nil
}

// ValidateChatmate validates a specific chatmate file.
//
// This method performs security validation and content checks on a single
// chatmate file to ensure it meets all requirements.
//
// Parameters:
//   - filename: The chatmate filename to validate
//
// Returns:
//   - error: Validation failure or file error
//   - bool: True if chatmate is valid, false otherwise
//
// Example:
//
// valid, err := validator.ValidateChatmate("Chatmate - Solve Issue.chatmode.md")
//
//	if err != nil {
//	   return fmt.Errorf("chatmate validation failed: %w", err)
//	}
func (v *ValidatorService) ValidateChatmate(filename string) (bool, error) {
	// Security validation
	if err := security.ValidateChatmateFilename(filename); err != nil {
		return false, fmt.Errorf("security validation failed: %w", err)
	}

	// Validate file extension
	if err := security.ValidateFileExtension(filename, []string{".md"}); err != nil {
		return false, fmt.Errorf("file extension validation failed: %w", err)
	}

	// Check if file exists in available chatmates
	availableChatmates, err := v.manager.GetAvailableChatmates()
	if err != nil {
		return false, err
	}

	found := false
	for _, available := range availableChatmates {
		if available == filename {
			found = true
			break
		}
	}

	if !found {
		return false, fmt.Errorf("chatmate file not found in available chatmates: %s", filename)
	}

	// Validate content if installed
	destPath := filepath.Join(v.manager.PromptsDir, filename)
	if _, err := os.Stat(destPath); err == nil {
		content, err := os.ReadFile(destPath)
		if err != nil {
			return false, fmt.Errorf("failed to read installed chatmate: %w", err)
		}

		if err := security.ValidateContentLength(content, 10*1024*1024); err != nil { // 10MB limit
			return false, fmt.Errorf("content validation failed: %w", err)
		}

		// Check for basic chatmate content structure
		contentStr := string(content)
		if !strings.Contains(contentStr, "---") {
			return false, fmt.Errorf("chatmate file appears to be missing YAML frontmatter")
		}
	}

	return true, nil
}

// validatePromptsDirectory checks the prompts directory status.
func (v *ValidatorService) validatePromptsDirectory() error {
	fmt.Printf("Checking prompts directory: %s\n", v.manager.PromptsDir)

	// Check if directory exists
	info, err := os.Stat(v.manager.PromptsDir)
	if os.IsNotExist(err) {
		return fmt.Errorf("prompts directory does not exist: %s", v.manager.PromptsDir)
	}
	if err != nil {
		return fmt.Errorf("failed to check prompts directory: %w", err)
	}

	// Check if it's actually a directory
	if !info.IsDir() {
		return fmt.Errorf("prompts path exists but is not a directory: %s", v.manager.PromptsDir)
	}

	// Check permissions
	if err := v.checkDirectoryPermissions(v.manager.PromptsDir); err != nil {
		return fmt.Errorf("prompts directory permission issue: %w", err)
	}

	fmt.Println("✅ Prompts directory is valid")
	return nil
}

// validateAvailableChatmates checks available chatmates.
func (v *ValidatorService) validateAvailableChatmates() error {
	fmt.Println("Checking available chatmates...")

	availableChatmates, err := v.manager.GetAvailableChatmates()
	if err != nil {
		return fmt.Errorf("failed to get available chatmates: %w", err)
	}

	if len(availableChatmates) == 0 {
		return fmt.Errorf("no chatmates available for installation")
	}

	// Validate each available chatmate
	for _, filename := range availableChatmates {
		if err := security.ValidateChatmateFilename(filename); err != nil {
			return fmt.Errorf("invalid chatmate filename %s: %w", filename, err)
		}
	}

	fmt.Printf("✅ Found %d valid available chatmates\n", len(availableChatmates))
	return nil
}

// validateInstalledChatmates checks installed chatmates.
func (v *ValidatorService) validateInstalledChatmates() error {
	fmt.Println("Checking installed chatmates...")

	installedChatmates, err := v.manager.GetInstalledChatmates()
	if err != nil {
		return fmt.Errorf("failed to get installed chatmates: %w", err)
	}

	// Validate each installed chatmate
	for _, filename := range installedChatmates {
		if _, err := v.ValidateChatmate(filename); err != nil {
			fmt.Printf("⚠️  Validation issue with %s: %v\n", filename, err)
		}
	}

	fmt.Printf("✅ Validated %d installed chatmates\n", len(installedChatmates))
	return nil
}

// validateOrphanedFiles checks for orphaned files.
func (v *ValidatorService) validateOrphanedFiles() error {
	fmt.Println("Checking for orphaned files...")

	availableChatmates, err := v.manager.GetAvailableChatmates()
	if err != nil {
		return err
	}

	installedChatmates, err := v.manager.GetInstalledChatmates()
	if err != nil {
		return err
	}

	orphanedCount := v.countOrphanedFiles(availableChatmates, installedChatmates)
	if orphanedCount > 0 {
		fmt.Printf("⚠️  Found %d orphaned files\n", orphanedCount)
	} else {
		fmt.Println("✅ No orphaned files found")
	}

	return nil
}

// countOrphanedFiles counts files that are installed but not available.
func (v *ValidatorService) countOrphanedFiles(available, installed []string) int {
	availableSet := make(map[string]bool)
	for _, filename := range available {
		availableSet[filename] = true
	}

	orphanedCount := 0
	for _, filename := range installed {
		if !availableSet[filename] {
			orphanedCount++
		}
	}

	return orphanedCount
}

// checkDirectoryPermissions validates directory access permissions.
func (v *ValidatorService) checkDirectoryPermissions(dir string) error {
	// Try to create a temporary file to check write permissions
	tempFile := filepath.Join(dir, ".chatmate_temp_permission_check")
	file, err := os.Create(tempFile)
	if err != nil {
		return fmt.Errorf("no write permission: %w", err)
	}
	file.Close()

	// Clean up the temporary file
	if err := os.Remove(tempFile); err != nil {
		return fmt.Errorf("cleanup failed: %w", err)
	}

	return nil
}
