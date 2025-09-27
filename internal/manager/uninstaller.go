// Package manager provides uninstallation functionality for ChatMate agents.
package manager

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jonassiebler/chatmate/pkg/security"
)

// UninstallerService handles chatmate uninstallation operations.
type UninstallerService struct {
	manager *ChatMateManager
}

// NewUninstallerService creates a new uninstaller service.
func NewUninstallerService(manager *ChatMateManager) *UninstallerService {
	return &UninstallerService{manager: manager}
}

// UninstallAll removes all installed chatmate agents.
//
// This method removes all chatmate files from the VS Code user prompts directory.
// It performs security validation and provides detailed feedback about the operation.
//
// Returns:
//   - error: Uninstallation failure or system error
//
// Example:
//
//err := uninstaller.UninstallAll()
//if err != nil {
//    return fmt.Errorf("uninstallation failed: %w", err)
//}
func (u *UninstallerService) UninstallAll() error {
	// Get available chatmodes from repository (only these should be uninstalled)
	availableChatmates, err := u.manager.GetAvailableChatmates()
	if err != nil {
		return err
	}

	// Get installed chatmodes
	installedChatmates, err := u.manager.GetInstalledChatmates()
	if err != nil {
		return err
	}

	// Create a set of available chatmodes for quick lookup
	availableSet := make(map[string]bool)
	for _, filename := range availableChatmates {
		availableSet[filename] = true
	}

	// Filter installed chatmodes to only include those available in repository
	var toUninstall []string
	var userCreated []string
	
	for _, filename := range installedChatmates {
		if availableSet[filename] {
			toUninstall = append(toUninstall, filename)
		} else {
			userCreated = append(userCreated, filename)
		}
	}

	if len(toUninstall) == 0 {
		fmt.Println("No repository chatmates are currently installed")
		if len(userCreated) > 0 {
			fmt.Printf("📝 Found %d user-created chatmate(s) (will be preserved):\n", len(userCreated))
			for _, filename := range userCreated {
				displayName := u.manager.getDisplayName(filename)
				fmt.Printf("  - %s\n", displayName)
			}
		}
		return nil
	}

	fmt.Printf("Uninstalling %d repository chatmates from: %s\n", len(toUninstall), u.manager.PromptsDir)
	
	if len(userCreated) > 0 {
		fmt.Printf("📝 Preserving %d user-created chatmate(s):\n", len(userCreated))
		for _, filename := range userCreated {
			displayName := u.manager.getDisplayName(filename)
			fmt.Printf("  - %s\n", displayName)
		}
	}
	fmt.Println()

	for _, chatmate := range toUninstall {
		if err := u.UninstallChatmate(chatmate); err != nil {
			return err
		}
	}

	fmt.Printf("\n✅ Successfully uninstalled %d repository chatmates\n", len(toUninstall))
	if len(userCreated) > 0 {
		fmt.Printf("📝 Preserved %d user-created chatmate(s)\n", len(userCreated))
	}
	return nil
}

// UninstallSpecific removes specific chatmate agents by name.
//
// This method takes a list of agent names and attempts to uninstall each one.
// Agent names should match the display names (e.g., "Solve Issue") rather than
// filenames. The method automatically converts names to appropriate filenames.
//
// Parameters:
//   - agentNames: List of chatmate display names to uninstall
//
// Returns:
//   - error: Uninstallation failure or agent not found error
//
// Example:
//
//names := []string{"Solve Issue", "Code Review", "Testing"}
//err := uninstaller.UninstallSpecific(names)
//if err != nil {
//    return fmt.Errorf("specific uninstallation failed: %w", err)
//}
func (u *UninstallerService) UninstallSpecific(agentNames []string) error {
	if len(agentNames) == 0 {
		fmt.Println("No specific chatmates specified for uninstallation")
		return nil
	}

	installedChatmates, err := u.manager.GetInstalledChatmates()
	if err != nil {
		return err
	}

	// Create a map for quick lookup of installed chatmates
	installedMap := make(map[string]string)
	for _, filename := range installedChatmates {
		displayName := u.manager.getDisplayName(filename)
		installedMap[displayName] = filename
	}

	fmt.Printf("Uninstalling specific chatmates: %v\n", agentNames)

	// Uninstall each specified agent
	for _, agentName := range agentNames {
		if filename, exists := installedMap[agentName]; exists {
			if err := u.UninstallChatmate(filename); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("chatmate not found or not installed: %s", agentName)
		}
	}

	return nil
}

// UninstallChatmate removes a single chatmate file.
//
// This method handles the removal of a single chatmate file with appropriate
// security validation and error handling.
//
// Parameters:
//   - filename: The chatmate filename (e.g., "Chatmate - Solve Issue.chatmode.md")
//
// Returns:
//   - error: Security validation or file operation error
//
// Security Features:
//   - Validates filename against security rules
//   - Checks path safety before deletion
//   - Sanitizes input for additional safety
func (u *UninstallerService) UninstallChatmate(filename string) error {
	// Security validation
	if err := security.ValidateChatmateFilename(filename); err != nil {
		return fmt.Errorf("security validation failed: %w", err)
	}

	// Validate path safety
	if !security.IsPathSafe(u.manager.PromptsDir, filename) {
		return fmt.Errorf("file path is not safe: %s", filename)
	}

	// Sanitize filename for extra safety
	filename = security.SanitizeInput(filename)

	destPath := filepath.Join(u.manager.PromptsDir, filename)

	// Check if file exists
	if _, err := os.Stat(destPath); os.IsNotExist(err) {
		fmt.Printf("⏭️  %s (not installed)\n", filename)
		return nil
	}

	// Remove the file
	if err := os.Remove(destPath); err != nil {
		return fmt.Errorf("failed to remove chatmate file %s: %w", destPath, err)
	}

	fmt.Printf("❌ %s (uninstalled)\n", filename)
	return nil
}

// CleanupOrphanedFiles removes any chatmate files that are no longer available.
//
// This method identifies installed chatmate files that don't exist in the
// available chatmates list and removes them. This is useful for cleaning up
// after chatmate updates or configuration changes.
//
// Returns:
//   - error: Cleanup operation failure
//   - int: Number of orphaned files removed
//
// Example:
//
//removed, err := uninstaller.CleanupOrphanedFiles()
//if err != nil {
//    return fmt.Errorf("cleanup failed: %w", err)
//}
//fmt.Printf("Removed %d orphaned files", removed)
func (u *UninstallerService) CleanupOrphanedFiles() (int, error) {
	installedChatmates, err := u.manager.GetInstalledChatmates()
	if err != nil {
		return 0, err
	}

	availableChatmates, err := u.manager.GetAvailableChatmates()
	if err != nil {
		return 0, err
	}

	// Create a set of available chatmates for quick lookup
	availableSet := make(map[string]bool)
	for _, filename := range availableChatmates {
		availableSet[filename] = true
	}

	// Find orphaned files
	var orphaned []string
	for _, installed := range installedChatmates {
		if !availableSet[installed] {
			orphaned = append(orphaned, installed)
		}
	}

	if len(orphaned) == 0 {
		fmt.Println("No orphaned chatmate files found")
		return 0, nil
	}

	fmt.Printf("Found %d orphaned chatmate files\n", len(orphaned))

	// Remove orphaned files
	for _, filename := range orphaned {
		if err := u.UninstallChatmate(filename); err != nil {
			return len(orphaned), err
		}
	}

	fmt.Printf("✅ Cleaned up %d orphaned chatmate files\n", len(orphaned))
	return len(orphaned), nil
}
