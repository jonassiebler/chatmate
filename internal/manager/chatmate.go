// Package manager provides the core functionality for managing ChatMate agents.
//
// The manager package is responsible for installing, uninstalling, and managing
// chatmate agents (specialized AI prompts) for VS Code Copilot Chat integration.
// It handles file operations, validation, and interaction with the VS Code
// user prompts directory.
//
// Key Components:
//   - ChatMateManager: Main service for chatmate operations
//   - Installation and uninstallation of chatmate files
//   - System status checking and configuration management
//   - Integration with VS Code prompts directory
//
// Usage Example:
//
//	manager, err := manager.NewChatMateManager()
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Install all available chatmates
//	err = manager.InstallAll(false)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// List installed chatmates
//	err = manager.ListChatmates(false, true)
//	if err != nil {
//	    log.Fatal(err)
//	}
package manager

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/jonassiebler/chatmate/internal/assets"
	"github.com/jonassiebler/chatmate/pkg/security"
	"github.com/jonassiebler/chatmate/pkg/utils"
)

// ChatMateManager handles the core functionality for managing chatmate agents.
//
// The ChatMateManager provides methods for installing, uninstalling, and managing
// specialized AI agents (chatmates) that integrate with VS Code Copilot Chat.
// It manages file operations between the chatmate source directory and the
// VS Code user prompts directory.
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
//	manager, err := NewChatMateManager()
//	if err != nil {
//	    return fmt.Errorf("failed to initialize manager: %w", err)
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

	return &ChatMateManager{
		ScriptDir:   scriptDir,
		MatesDir:    matesDir,
		PromptsDir:  promptsDir,
		UseEmbedded: useEmbedded,
	}, nil
}

// GetAvailableChatmates returns all available chatmate files
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

// GetInstalledChatmates returns all installed chatmate files
func (cm *ChatMateManager) GetInstalledChatmates() ([]string, error) {
	// Check if prompts directory exists
	if _, err := os.Stat(cm.PromptsDir); os.IsNotExist(err) {
		return []string{}, nil
	}

	files, err := os.ReadDir(cm.PromptsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read prompts directory: %w", err)
	}

	var chatmates []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".chatmode.md") {
			chatmates = append(chatmates, file.Name())
		}
	}

	return chatmates, nil
}

// InstallAll installs all available chatmate agents
func (cm *ChatMateManager) InstallAll(force bool) error {
	// Ensure prompts directory exists
	if err := os.MkdirAll(cm.PromptsDir, 0755); err != nil {
		return fmt.Errorf("failed to create prompts directory: %w", err)
	}

	chatmates, err := cm.GetAvailableChatmates()
	if err != nil {
		return err
	}

	if len(chatmates) == 0 {
		return fmt.Errorf("no chatmate files found in mates directory")
	}

	fmt.Printf("Installing %d chatmates to: %s\n\n", len(chatmates), cm.PromptsDir)

	for _, chatmate := range chatmates {
		if err := cm.InstallChatmate(chatmate, force); err != nil {
			return err
		}
	}

	return nil
}

// InstallSpecific installs specific chatmate agents
func (cm *ChatMateManager) InstallSpecific(agentNames []string, force bool) error {
	// Ensure prompts directory exists
	if err := os.MkdirAll(cm.PromptsDir, 0755); err != nil {
		return fmt.Errorf("failed to create prompts directory: %w", err)
	}

	availableChatmates, err := cm.GetAvailableChatmates()
	if err != nil {
		return err
	}

	for _, agentName := range agentNames {
		var matchingFiles []string

		// Find matching files
		for _, file := range availableChatmates {
			if strings.Contains(strings.ToLower(file), strings.ToLower(agentName)) ||
				file == agentName+".chatmode.md" {
				matchingFiles = append(matchingFiles, file)
			}
		}

		if len(matchingFiles) == 0 {
			fmt.Printf("⚠️  No chatmate found matching: %s\n", agentName)
			continue
		}

		if len(matchingFiles) > 1 {
			fmt.Printf("⚠️  Multiple chatmates found for \"%s\":\n", agentName)
			for _, file := range matchingFiles {
				fmt.Printf("    - %s\n", file)
			}
			fmt.Printf("    Installing all matches...\n\n")
		}

		for _, file := range matchingFiles {
			if err := cm.InstallChatmate(file, force); err != nil {
				return err
			}
		}
	}

	return nil
}

// InstallChatmate installs a single chatmate file
func (cm *ChatMateManager) InstallChatmate(filename string, force bool) error {
	// Security validation
	if err := security.ValidateChatmateFilename(filename); err != nil {
		return fmt.Errorf("security validation failed: %w", err)
	}

	// Validate destination path safety
	if !security.IsPathSafe(cm.PromptsDir, filename) {
		return fmt.Errorf("destination path is not safe: %s", filename)
	}

	// Sanitize filename for extra safety
	filename = security.SanitizeInput(filename)

	destPath := filepath.Join(cm.PromptsDir, filename)

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

	if cm.UseEmbedded {
		// Use embedded file
		content, err = assets.GetEmbeddedMateContent(filename)
		if err != nil {
			fmt.Printf("❌ %s (embedded file not found: %v)\n", filename, err)
			return err
		}
	} else {
		// Use filesystem file
		sourcePath := filepath.Join(cm.MatesDir, filename)

		// Security: Validate source path
		if !security.IsPathSafe(cm.MatesDir, filename) {
			return fmt.Errorf("source path is not safe: %s", filename)
		}

		// Check if source file exists
		if _, err := os.Stat(sourcePath); os.IsNotExist(err) {
			return fmt.Errorf("chatmate file not found: %s", filename)
		}

		content, err = os.ReadFile(sourcePath)
		if err != nil {
			fmt.Printf("❌ %s (failed to read: %v)\n", filename, err)
			return err
		}
	}

	// Security: Validate content size (max 10MB for chatmate files)
	const maxChatmateSize = 10 * 1024 * 1024
	if err := security.ValidateContentLength(content, maxChatmateSize); err != nil {
		return fmt.Errorf("content validation failed: %w", err)
	}

	// Write file
	if err := os.WriteFile(destPath, content, 0644); err != nil {
		fmt.Printf("❌ %s (failed to write: %v)\n", filename, err)
		return err
	}

	status := "installed"
	if force {
		status = "reinstalled"
	}
	fmt.Printf("✅ %s (%s)\n", filename, status)

	return nil
}

// UninstallAll removes all installed chatmate agents that are available in the repository
// This preserves local-only chatmates that cannot be reinstalled
func (cm *ChatMateManager) UninstallAll() error {
	availableChatmates, err := cm.GetAvailableChatmates()
	if err != nil {
		return err
	}

	installedChatmates, err := cm.GetInstalledChatmates()
	if err != nil {
		return err
	}

	if len(installedChatmates) == 0 {
		fmt.Println("No chatmates currently installed.")
		return nil
	}

	// Create a map of available chatmates for quick lookup
	availableMap := make(map[string]bool)
	for _, chatmate := range availableChatmates {
		availableMap[chatmate] = true
	}

	// Only uninstall chatmates that are available in the repository
	var uninstalledCount int
	var skippedLocalOnly []string

	for _, chatmate := range installedChatmates {
		if availableMap[chatmate] {
			if err := cm.UninstallChatmate(chatmate); err != nil {
				return err
			}
			uninstalledCount++
		} else {
			skippedLocalOnly = append(skippedLocalOnly, chatmate)
		}
	}

	if uninstalledCount > 0 {
		fmt.Printf("Uninstalled %d repository chatmates.\n", uninstalledCount)
	}

	if len(skippedLocalOnly) > 0 {
		fmt.Printf("Preserved %d local-only chatmates: %s\n",
			len(skippedLocalOnly), strings.Join(skippedLocalOnly, ", "))
	}

	return nil
}

// UninstallSpecific removes specific chatmate agents
func (cm *ChatMateManager) UninstallSpecific(agentNames []string) error {
	installedChatmates, err := cm.GetInstalledChatmates()
	if err != nil {
		return err
	}

	for _, agentName := range agentNames {
		var matchingFiles []string

		// Find matching files
		for _, file := range installedChatmates {
			if strings.Contains(strings.ToLower(file), strings.ToLower(agentName)) ||
				file == agentName+".chatmode.md" {
				matchingFiles = append(matchingFiles, file)
			}
		}

		if len(matchingFiles) == 0 {
			fmt.Printf("⚠️  No installed chatmate found matching: %s\n", agentName)
			continue
		}

		for _, file := range matchingFiles {
			if err := cm.UninstallChatmate(file); err != nil {
				return err
			}
		}
	}

	return nil
}

// UninstallChatmate removes a single chatmate file
func (cm *ChatMateManager) UninstallChatmate(filename string) error {
	// Security validation
	if err := security.ValidateChatmateFilename(filename); err != nil {
		return fmt.Errorf("security validation failed: %w", err)
	}

	// Validate path safety
	if !security.IsPathSafe(cm.PromptsDir, filename) {
		return fmt.Errorf("path is not safe: %s", filename)
	}

	// Sanitize filename for extra safety
	filename = security.SanitizeInput(filename)

	filePath := filepath.Join(cm.PromptsDir, filename)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Printf("⚠️  %s (not installed)\n", filename)
		return nil
	}

	if err := os.Remove(filePath); err != nil {
		fmt.Printf("❌ %s (failed to uninstall: %v)\n", filename, err)
		return err
	}

	fmt.Printf("🗑️  %s (uninstalled)\n", filename)
	return nil
}

// ListChatmates displays available and installed chatmates
func (cm *ChatMateManager) ListChatmates(showAvailable, showInstalled bool) error {
	availableChatmates, err := cm.GetAvailableChatmates()
	if err != nil {
		return err
	}

	installedChatmates, err := cm.GetInstalledChatmates()
	if err != nil {
		return err
	}

	// Create a map for quick lookup
	installedMap := make(map[string]bool)
	for _, installed := range installedChatmates {
		installedMap[installed] = true
	}

	if showAvailable {
		fmt.Println("📦 Available Chatmates:")
		if len(availableChatmates) == 0 {
			fmt.Println("  No chatmates available.")
		} else {
			for _, chatmate := range availableChatmates {
				name := strings.TrimSuffix(chatmate, ".chatmode.md")
				if installedMap[chatmate] {
					fmt.Printf("  %s ✅ installed\n", name)
				} else {
					fmt.Printf("  %s ⏸️  available\n", name)
				}
			}
		}
		fmt.Println()
	}

	if showInstalled {
		fmt.Println("✅ Installed Chatmates:")
		if len(installedChatmates) == 0 {
			fmt.Println("  No chatmates currently installed.")
			fmt.Println("  Run \"chatmate hire\" to install all available chatmates.")
		} else {
			for _, chatmate := range installedChatmates {
				name := strings.TrimSuffix(chatmate, ".chatmode.md")
				fmt.Printf("  %s ✅ installed\n", name)
			}
		}
		fmt.Println()
	}

	// Summary
	fmt.Printf("📊 Summary: %d/%d chatmates installed\n", len(installedChatmates), len(availableChatmates))

	return nil
}

// ShowStatus displays ChatMate and VS Code installation status
func (cm *ChatMateManager) ShowStatus() error {
	fmt.Println("🔍 ChatMate Installation Status")
	fmt.Println()

	// Check VS Code installation
	vsCodeInstalled := cm.checkVSCodeInstallation()
	if vsCodeInstalled {
		fmt.Println("VS Code: ✅ VS Code detected")
	} else {
		fmt.Println("VS Code: ❌ VS Code not found")
	}

	// Check prompts directory
	promptsDirExists := true
	if _, err := os.Stat(cm.PromptsDir); os.IsNotExist(err) {
		promptsDirExists = false
	}

	if promptsDirExists {
		fmt.Println("Prompts Directory: ✅ Prompts directory exists")
	} else {
		fmt.Println("Prompts Directory: ⚠️  Prompts directory not found")
	}
	fmt.Printf("Path: %s\n\n", cm.PromptsDir)

	// Show chatmate statistics
	availableChatmates, err := cm.GetAvailableChatmates()
	if err != nil {
		return err
	}

	installedChatmates, err := cm.GetInstalledChatmates()
	if err != nil {
		return err
	}

	fmt.Println("📊 Chatmate Statistics:")
	fmt.Printf("Available: %d chatmates\n", len(availableChatmates))
	fmt.Printf("Installed: %d chatmates\n", len(installedChatmates))

	if len(installedChatmates) < len(availableChatmates) {
		uninstalled := len(availableChatmates) - len(installedChatmates)
		fmt.Printf("Pending: %d chatmates not installed\n", uninstalled)
		fmt.Println("\nRun \"chatmate hire\" to install all available chatmates.")
	}

	return nil
}

// checkVSCodeInstallation checks if VS Code is installed on the system
func (cm *ChatMateManager) checkVSCodeInstallation() bool {
	switch runtime.GOOS {
	case "darwin": // macOS
		if _, err := os.Stat("/Applications/Visual Studio Code.app"); err == nil {
			return true
		}
		return false

	case "linux": // Linux
		paths := []string{
			"/usr/bin/code",
			"/usr/local/bin/code",
			"/snap/bin/code",
		}
		for _, path := range paths {
			if _, err := os.Stat(path); err == nil {
				return true
			}
		}
		return false

	case "windows": // Windows
		programFiles := os.Getenv("PROGRAMFILES")
		if programFiles == "" {
			programFiles = "C:\\Program Files"
		}
		programFilesX86 := os.Getenv("PROGRAMFILES(X86)")
		if programFilesX86 == "" {
			programFilesX86 = "C:\\Program Files (x86)"
		}

		paths := []string{
			filepath.Join(programFiles, "Microsoft VS Code", "Code.exe"),
			filepath.Join(programFilesX86, "Microsoft VS Code", "Code.exe"),
		}
		for _, path := range paths {
			if _, err := os.Stat(path); err == nil {
				return true
			}
		}
		return false

	default:
		return false
	}
}

// ShowConfig displays current configuration
func (cm *ChatMateManager) ShowConfig() {
	fmt.Println("⚙️  ChatMate Configuration:")
	fmt.Println()

	if cm.UseEmbedded {
		fmt.Printf("Mates Directory: embedded files (self-contained binary)\n")
	} else {
		fmt.Printf("Mates Directory: %s\n", cm.MatesDir)
	}

	fmt.Printf("Prompts Directory: %s\n", cm.PromptsDir)
	fmt.Printf("Platform: %s\n", runtime.GOOS)
	fmt.Printf("Go Version: %s\n", runtime.Version())
}
