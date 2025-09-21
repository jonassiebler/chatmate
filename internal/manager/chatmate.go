// Package manager provides the core functionality for managing ChatMate agents.
//
// The manager package is responsible for installing, uninstalling, and managing
// chatmate agents (specialized AI prompts) for VS Code Copilot Chat integration.
// It handles file operations, validation, and interaction with the VS Code
// user prompts directory.
//
// Key Components:
//   - ChatMateManager: Main service for chatmate operations
//   - InstallerService: Handles chatmate installation operations
//   - UninstallerService: Handles chatmate removal operations
//   - ListerService: Handles chatmate listing and display
//   - ValidatorService: Handles validation and status checking
//
// Usage Example:
//
// manager, err := manager.NewChatMateManager()
//
//	if err != nil {
//	   log.Fatal(err)
//	}
//
// // Install all available chatmates
// err = manager.Installer().InstallAll(false)
//
//	if err != nil {
//	   log.Fatal(err)
//	}
//
// // List installed chatmates
// err = manager.Lister().ListInstalled()
//
//	if err != nil {
//	   log.Fatal(err)
//	}
package manager

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jonassiebler/chatmate/internal/assets"
	"github.com/jonassiebler/chatmate/pkg/utils"
)

// ChatMateManager handles the core functionality for managing chatmate agents.
//
// The ChatMateManager provides methods for installing, uninstalling, and managing
// specialized AI agents (chatmates) that integrate with VS Code Copilot Chat.
// It manages file operations between the chatmate source directory and the
// VS Code user prompts directory.
//
// The manager uses a service-oriented architecture with dedicated modules for
// different functionality areas: installation, uninstallation, listing, and validation.
//
// Fields:
//   - ScriptDir: Directory containing the ChatMate executable and resources
//   - MatesDir: Directory containing chatmate source files (.chatmode.md)
//   - PromptsDir: VS Code user prompts directory where chatmates are installed
//   - UseEmbedded: Whether to use embedded chatmate resources or external files
type ChatMateManager struct {
	ScriptDir   string
	MatesDir    string
	PromptsDir  string
	UseEmbedded bool

	// Service instances for modular functionality
	installer   *InstallerService
	uninstaller *UninstallerService
	lister      *ListerService
	validator   *ValidatorService
	status      *StatusService
}

// NewChatMateManager creates a new ChatMateManager instance with automatic configuration.
//
// This constructor automatically detects the execution environment and configures
// appropriate directories for chatmate operations:
//
//   - Development mode: Uses current working directory if "mates" folder exists
//   - Production mode: Uses executable directory with embedded resources
//   - Fallback: Uses current working directory
//
// The manager automatically detects the VS Code user prompts directory based on
// the operating system and creates it if it doesn't exist.
//
// Returns:
//   - *ChatMateManager: Configured manager instance
//   - error: Configuration or directory creation error
//
// Example:
//
// manager, err := NewChatMateManager()
//
//	if err != nil {
//	   return fmt.Errorf("failed to initialize manager: %w", err)
//	}
func NewChatMateManager() (*ChatMateManager, error) {
	// Get current working directory (for development) or executable directory (for production)
	var scriptDir string
	var useEmbedded bool

	// First try current working directory
	if workDir, err := os.Getwd(); err == nil {
		// Check if we're in a development environment (mates directory exists)
		if _, err := os.Stat(filepath.Join(workDir, "mates")); err == nil {
			scriptDir = workDir
			useEmbedded = false
		}
	}

	// If not found, try executable directory
	if scriptDir == "" {
		execPath, err := os.Executable()
		if err != nil {
			return nil, fmt.Errorf("failed to get executable path: %w", err)
		}
		scriptDir = filepath.Dir(execPath)

		// Check if mates directory exists in executable directory
		if _, err := os.Stat(filepath.Join(scriptDir, "mates")); err != nil {
			// No mates directory found, use embedded files
			useEmbedded = true
		}
	}

	matesDir := filepath.Join(scriptDir, "mates")

	promptsDir, err := utils.GetVSCodePromptsDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get VS Code prompts directory: %w", err)
	}

	// Create manager instance
	manager := &ChatMateManager{
		ScriptDir:   scriptDir,
		MatesDir:    matesDir,
		PromptsDir:  promptsDir,
		UseEmbedded: useEmbedded,
	}

	// Initialize service modules
	manager.installer = NewInstallerService(manager)
	manager.uninstaller = NewUninstallerService(manager)
	manager.lister = NewListerService(manager)
	manager.validator = NewValidatorService(manager)
	manager.status = NewStatusService(manager)

	return manager, nil
}

// Installer returns the installer service for chatmate installation operations.
func (cm *ChatMateManager) Installer() *InstallerService {
	return cm.installer
}

// Uninstaller returns the uninstaller service for chatmate removal operations.
func (cm *ChatMateManager) Uninstaller() *UninstallerService {
	return cm.uninstaller
}

// Lister returns the lister service for chatmate listing and display operations.
func (cm *ChatMateManager) Lister() *ListerService {
	return cm.lister
}

// Validator returns the validator service for validation and status operations.
func (cm *ChatMateManager) Validator() *ValidatorService {
	return cm.validator
}

// Status returns the status service for status and configuration display operations.
func (cm *ChatMateManager) Status() *StatusService {
	return cm.status
}

// GetAvailableChatmates returns all available chatmate files.
//
// This method retrieves chatmates from either embedded resources or external files
// based on the UseEmbedded configuration. It's used by service modules to get
// the list of chatmates available for operations.
//
// Returns:
//   - []string: List of available chatmate filenames
//   - error: Directory reading or embedded resource access error
func (cm *ChatMateManager) GetAvailableChatmates() ([]string, error) {
	if cm.UseEmbedded {
		// Use embedded files
		return assets.GetEmbeddedMatesList()
	}

	// Use filesystem files
	files, err := os.ReadDir(cm.MatesDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read mates directory: %w", err)
	}

	var chatmates []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".chatmode.md") {
			chatmates = append(chatmates, file.Name())
		}
	}

	return chatmates, nil
}

// GetInstalledChatmates returns all currently installed chatmate files.
//
// This method scans the VS Code prompts directory to find installed chatmate files.
// It's used by service modules to determine which chatmates are currently available
// in the user's VS Code environment.
//
// Returns:
//   - []string: List of installed chatmate filenames
//   - error: Directory reading or access error
func (cm *ChatMateManager) GetInstalledChatmates() ([]string, error) {
	files, err := os.ReadDir(cm.PromptsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read prompts directory: %w", err)
	}

	var installed []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".chatmode.md") {
			installed = append(installed, file.Name())
		}
	}

	return installed, nil
}

// getDisplayName extracts a user-friendly display name from a chatmate filename.
//
// This method converts filenames like "Chatmate - Solve Issue.chatmode.md"
// to display names like "Solve Issue". It's used by service modules to provide
// clean, user-friendly output.
//
// Parameters:
//   - filename: The chatmate filename to convert
//
// Returns:
//   - string: User-friendly display name
func (cm *ChatMateManager) getDisplayName(filename string) string {
	// Remove the file extension
	name := strings.TrimSuffix(filename, ".chatmode.md")

	// Remove the "Chatmate - " prefix if present
	if strings.HasPrefix(name, "Chatmate - ") {
		name = strings.TrimPrefix(name, "Chatmate - ")
	}

	return name
}
