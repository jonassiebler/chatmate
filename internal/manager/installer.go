// Package manager provides installation functionality for ChatMate agents.
package manager

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jonassiebler/chatmate/internal/assets"
	"github.com/jonassiebler/chatmate/pkg/security"
)

// InstallerService handles chatmate installation operations.
type InstallerService struct {
	manager *ChatMateManager
}

// NewInstallerService creates a new installer service.
func NewInstallerService(manager *ChatMateManager) *InstallerService {
	return &InstallerService{manager: manager}
}

// InstallAll installs all available chatmate agents.
//
// This method installs all chatmate files from the source directory (or embedded
// resources) to the VS Code user prompts directory. It handles file conflicts
// based on the force parameter.
//
// Parameters:
//   - force: If true, overwrites existing chatmate files; if false, skips existing files
//
// Returns:
//   - error: Installation failure or system error
//
// Example:
//
// err := installer.InstallAll(false)
//
//	if err != nil {
//	   return fmt.Errorf("installation failed: %w", err)
//	}
func (i *InstallerService) InstallAll(force bool) error {
	availableChatmates, err := i.manager.GetAvailableChatmates()
	if err != nil {
		return err
	}

	if len(availableChatmates) == 0 {
		fmt.Println("No chatmates available to install")
		return nil
	}

	fmt.Printf("Installing %d chatmates to: %s\n\n", len(availableChatmates), i.manager.PromptsDir)

	for _, chatmate := range availableChatmates {
		if err := i.InstallChatmate(chatmate, force); err != nil {
			return err
		}
	}

	return nil
}

// InstallSpecific installs specific chatmate agents by name.
//
// This method takes a list of agent names and attempts to install each one.
// Agent names should match the display names (e.g., "Solve Issue") rather than
// filenames. The method automatically converts names to appropriate filenames.
//
// Parameters:
//   - agentNames: List of chatmate display names to install
//   - force: If true, overwrites existing files; if false, skips existing files
//
// Returns:
//   - error: Installation failure or agent not found error
//
// Example:
//
// names := []string{"Solve Issue", "Code Review", "Testing"}
// err := installer.InstallSpecific(names, false)
//
//	if err != nil {
//	   return fmt.Errorf("specific installation failed: %w", err)
//	}
func (i *InstallerService) InstallSpecific(agentNames []string, force bool) error {
	if len(agentNames) == 0 {
		fmt.Println("No specific chatmates specified")
		return nil
	}

	availableChatmates, err := i.manager.GetAvailableChatmates()
	if err != nil {
		return err
	}

	// Create a map for quick lookup of available chatmates
	availableMap := make(map[string]string)
	for _, filename := range availableChatmates {
		displayName := i.manager.getDisplayName(filename)
		availableMap[displayName] = filename
	}

	fmt.Printf("Installing specific chatmates: %v\n", agentNames)

	// Install each specified agent
	for _, agentName := range agentNames {
		if filename, exists := availableMap[agentName]; exists {
			if err := i.InstallChatmate(filename, force); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("chatmate not found: %s", agentName)
		}
	}

	return nil
}

// InstallChatmate installs a single chatmate file.
//
// This method handles the installation of a single chatmate file, including
// security validation, file existence checks, and content retrieval from
// either embedded resources or external files.
//
// Parameters:
//   - filename: The chatmate filename (e.g., "Chatmate - Solve Issue.chatmode.md")
//   - force: If true, overwrites existing files; if false, skips existing files
//
// Returns:
//   - error: Security validation, file operation, or content retrieval error
//
// Security Features:
//   - Validates filename against security rules
//   - Checks destination path safety
//   - Sanitizes input for additional safety
//   - Validates content length and file extensions
func (i *InstallerService) InstallChatmate(filename string, force bool) error {
	// Security validation
	if err := security.ValidateChatmateFilename(filename); err != nil {
		return fmt.Errorf("security validation failed: %w", err)
	}

	// Validate destination path safety
	if !security.IsPathSafe(i.manager.PromptsDir, filename) {
		return fmt.Errorf("destination path is not safe: %s", filename)
	}

	// Sanitize filename for extra safety
	filename = security.SanitizeInput(filename)

	destPath := filepath.Join(i.manager.PromptsDir, filename)

	// Check if already installed and not forcing
	if !force {
		if _, err := os.Stat(destPath); err == nil {
			fmt.Printf("⏭️  %s (already installed)\n", filename)
			return nil
		}
	}

	// Get file content
	var content []byte
	var err error

	if i.manager.UseEmbedded {
		// Use embedded file
		content, err = assets.GetEmbeddedMateContent(filename)
		if err != nil {
			return fmt.Errorf("failed to read embedded chatmate %s: %w", filename, err)
		}
	} else {
		// Use external file
		sourcePath := filepath.Join(i.manager.MatesDir, filename)
		content, err = os.ReadFile(sourcePath)
		if err != nil {
			return fmt.Errorf("failed to read chatmate file %s: %w", sourcePath, err)
		}
	}

	// Validate content length for security
	if err := security.ValidateContentLength(content, 10*1024*1024); err != nil { // 10MB limit
		return fmt.Errorf("content validation failed for %s: %w", filename, err)
	}

	// Validate file extension
	if err := security.ValidateFileExtension(filename, []string{".md"}); err != nil {
		return fmt.Errorf("file extension validation failed: %w", err)
	}

	// Write to destination
	if err := os.WriteFile(destPath, content, 0644); err != nil {
		return fmt.Errorf("failed to write chatmate file %s: %w", destPath, err)
	}

	// Determine the status message
	status := "installed"
	if force {
		if _, err := os.Stat(destPath); err == nil {
			status = "reinstalled"
		}
	}

	fmt.Printf("✅ %s (%s)\n", filename, status)
	return nil
}
