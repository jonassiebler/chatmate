// Package manager provides installation functionality for ChatMate agents.
package manager

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

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

// checkAndRebuildIfNeeded checks if the chatmate binary needs rebuilding
// and rebuilds it if the source files are newer than the binary.
func (i *InstallerService) checkAndRebuildIfNeeded() error {
	// Only check when using embedded assets
	if !i.manager.UseEmbedded {
		return nil
	}

	// Get current binary path
	binaryPath, err := os.Executable()
	if err != nil {
		fmt.Printf("⚠️  Could not determine binary path, skipping build check: %v\n", err)
		return nil
	}

	// Get binary modification time
	binaryInfo, err := os.Stat(binaryPath)
	if err != nil {
		fmt.Printf("⚠️  Could not stat binary, skipping build check: %v\n", err)
		return nil
	}
	binaryTime := binaryInfo.ModTime()

	// Check if source chatmate files are newer than binary
	matesDir := filepath.Join(filepath.Dir(binaryPath), "internal", "assets", "mates")

	// First try relative to current working directory
	if cwd, err := os.Getwd(); err == nil {
		testMatesDir := filepath.Join(cwd, "internal", "assets", "mates")
		if _, err := os.Stat(testMatesDir); err == nil {
			matesDir = testMatesDir
		}
	}

	// Check if mates directory exists
	if _, err := os.Stat(matesDir); os.IsNotExist(err) {
		// No source files found, assume binary is up to date
		return nil
	}

	// Check modification time of source files
	needsRebuild := false
	err = filepath.Walk(matesDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(path) == ".md" && info.ModTime().After(binaryTime) {
			needsRebuild = true
			fmt.Printf("📅 Found newer file: %s (modified: %s, binary: %s)\n",
				filepath.Base(path),
				info.ModTime().Format(time.RFC3339),
				binaryTime.Format(time.RFC3339))
			return filepath.SkipDir // Stop walking
		}
		return nil
	})

	if err != nil {
		fmt.Printf("⚠️  Error checking source files, skipping build check: %v\n", err)
		return nil
	}

	if needsRebuild {
		fmt.Printf("🔨 Source chatmate files are newer than binary, rebuilding...\n")
		return i.rebuildBinary()
	}

	return nil
}

// rebuildBinary rebuilds the chatmate binary using go build
func (i *InstallerService) rebuildBinary() error {
	fmt.Printf("📦 Building chatmate binary with latest chatmate files...\n")

	// Use go build to rebuild the binary
	cmd := exec.Command("go", "build", "-o", "chatmate")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to rebuild binary: %w", err)
	}

	fmt.Printf("✅ Binary rebuilt successfully\n")
	return nil
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
	// Check if binary needs rebuilding first
	if err := i.checkAndRebuildIfNeeded(); err != nil {
		fmt.Printf("⚠️  Build check failed, continuing with current binary: %v\n", err)
	}

	availableChatmates, err := i.manager.GetAvailableChatmates()
	if err != nil {
		return err
	}

	if len(availableChatmates) == 0 {
		fmt.Println("No chatmates available to install")
		return nil
	}

	// Get installed chatmates to show what already exists
	installedChatmates, err := i.manager.GetInstalledChatmates()
	if err != nil {
		return err
	}

	// Create a set of installed chatmates for quick lookup
	installedSet := make(map[string]bool)
	for _, filename := range installedChatmates {
		installedSet[filename] = true
	}

	// Separate into categories for user visibility
	var toInstall []string
	var alreadyInstalled []string
	var userCreated []string

	// Get available chatmates as a set for lookup
	availableSet := make(map[string]bool)
	for _, filename := range availableChatmates {
		availableSet[filename] = true
	}

	// Categorize installed chatmates
	for _, filename := range installedChatmates {
		if availableSet[filename] {
			alreadyInstalled = append(alreadyInstalled, filename)
		} else {
			userCreated = append(userCreated, filename)
		}
	}

	// Determine what will be installed/reinstalled
	for _, filename := range availableChatmates {
		if installedSet[filename] {
			if force {
				toInstall = append(toInstall, filename)
			}
		} else {
			toInstall = append(toInstall, filename)
		}
	}

	// Safety confirmation - show what will be installed
	fmt.Printf("📦 INSTALLATION CONFIRMATION\n")

	if len(toInstall) > 0 {
		action := "INSTALLED"
		if force && len(alreadyInstalled) > 0 {
			action = "INSTALLED/REINSTALLED"
		}
		fmt.Printf("Repository chatmates to be %s (%d):\n", action, len(toInstall))
		for _, filename := range toInstall {
			displayName := i.manager.getDisplayName(filename)
			status := "✅"
			if installedSet[filename] && force {
				status = "🔄"
			}
			fmt.Printf("  %s %s\n", status, displayName)
		}
	}

	if !force && len(alreadyInstalled) > 0 {
		fmt.Printf("\nRepository chatmates already installed (will be SKIPPED) (%d):\n", len(alreadyInstalled))
		for _, filename := range alreadyInstalled {
			displayName := i.manager.getDisplayName(filename)
			fmt.Printf("  ⏭️  %s\n", displayName)
		}
	}

	if len(userCreated) > 0 {
		fmt.Printf("\nUser-created chatmates (will be PRESERVED) (%d):\n", len(userCreated))
		for _, filename := range userCreated {
			displayName := i.manager.getDisplayName(filename)
			fmt.Printf("  📝 %s\n", displayName)
		}
	}

	fmt.Printf("\nDirectory: %s\n", i.manager.PromptsDir)

	if len(toInstall) == 0 {
		fmt.Println("\n✅ All repository chatmates are already installed")
		return nil
	}

	forceMsg := ""
	if force {
		forceMsg = " (with force reinstall)"
	}
	fmt.Printf("\nDo you want to proceed with installing these chatmates%s? (y/N): ", forceMsg)

	var response string
	fmt.Scanln(&response)

	if response != "y" && response != "Y" && response != "yes" && response != "YES" {
		fmt.Println("❌ Installation operation cancelled by user")
		return nil
	}

	fmt.Printf("\nProceeding with installation...\n")

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

	// Check if binary needs rebuilding first
	if err := i.checkAndRebuildIfNeeded(); err != nil {
		fmt.Printf("⚠️  Build check failed, continuing with current binary: %v\n", err)
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
