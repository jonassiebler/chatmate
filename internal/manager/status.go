// Package manager provides status and configuration display functionality for ChatMate agents.
package manager

import (
	"fmt"
	"os"
)

// StatusService handles chatmate status and configuration display operations.
type StatusService struct {
	manager *ChatMateManager
}

// NewStatusService creates a new status service.
func NewStatusService(manager *ChatMateManager) *StatusService {
	return &StatusService{manager: manager}
}

// ShowStatus displays comprehensive status information.
//
// This method provides a detailed overview of the chatmate system status,
// including directory information, installation statistics, and configuration.
//
// Returns:
//   - error: Status retrieval failure
//
// Example:
//
//err := status.ShowStatus()
//if err != nil {
//    return fmt.Errorf("status display failed: %w", err)
//}
func (s *StatusService) ShowStatus() error {
	fmt.Println("=== ChatMate Status ===")

	// Directory Information
	fmt.Printf("VS Code Prompts Directory: %s\n", s.manager.PromptsDir)
	if !s.manager.UseEmbedded {
		fmt.Printf("Mates Source Directory: %s\n", s.manager.MatesDir)
	} else {
		fmt.Println("Using embedded chatmate resources")
	}

	// Check directory existence
	if _, err := os.Stat(s.manager.PromptsDir); os.IsNotExist(err) {
		fmt.Printf("❌ Prompts directory does not exist: %s\n", s.manager.PromptsDir)
	} else {
		fmt.Printf("✅ Prompts directory exists: %s\n", s.manager.PromptsDir)
	}

	// Get chatmate counts
	availableChatmates, err := s.manager.GetAvailableChatmates()
	if err != nil {
		return fmt.Errorf("failed to get available chatmates: %w", err)
	}

	installedChatmates, err := s.manager.GetInstalledChatmates()
	if err != nil {
		return fmt.Errorf("failed to get installed chatmates: %w", err)
	}

	// Installation Statistics
	fmt.Printf("\n=== Installation Statistics ===\n")
	fmt.Printf("Available Chatmates: %d\n", len(availableChatmates))
	fmt.Printf("Installed Chatmates: %d\n", len(installedChatmates))

	if len(availableChatmates) > 0 {
		percentage := float64(len(installedChatmates)) / float64(len(availableChatmates)) * 100
		fmt.Printf("Installation Coverage: %.1f%%\n", percentage)
	}

	// Check for issues
	orphanedCount := s.countOrphanedFiles(availableChatmates, installedChatmates)
	if orphanedCount > 0 {
		fmt.Printf("⚠️  Orphaned Files: %d (consider running cleanup)\n", orphanedCount)
	}

	// Configuration Information
	fmt.Printf("\n=== Configuration ===\n")
	fmt.Printf("Using Embedded Resources: %t\n", s.manager.UseEmbedded)

	// Recent Activity (if any logs exist)
	s.showRecentActivity()

	return nil
}

// ShowConfig displays the current ChatMate configuration.
//
// This method shows the current configuration settings for the ChatMate manager,
// including directory paths and operational modes.
//
// Example:
//
//status.ShowConfig()
func (s *StatusService) ShowConfig() {
	fmt.Println("=== ChatMate Configuration ===")
	fmt.Printf("Script Directory: %s\n", s.manager.ScriptDir)
	fmt.Printf("Mates Directory: %s\n", s.manager.MatesDir)
	fmt.Printf("VS Code Prompts Directory: %s\n", s.manager.PromptsDir)
	fmt.Printf("Using Embedded Resources: %t\n", s.manager.UseEmbedded)
}

// countOrphanedFiles counts files that are installed but not available.
func (s *StatusService) countOrphanedFiles(available, installed []string) int {
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

// showRecentActivity displays recent activity information if available.
func (s *StatusService) showRecentActivity() {
	fmt.Printf("\n=== Recent Activity ===\n")
	fmt.Println("(Activity logging not yet implemented)")
}
