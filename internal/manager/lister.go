// Package manager provides listing functionality for ChatMate agents.
package manager

import (
	"fmt"
	"sort"
	"strings"
)

// ListerService handles chatmate listing and display operations.
type ListerService struct {
	manager *ChatMateManager
}

// NewListerService creates a new lister service.
func NewListerService(manager *ChatMateManager) *ListerService {
	return &ListerService{manager: manager}
}

// ListAll displays all available and installed chatmate agents.
//
// This method provides a comprehensive overview of the chatmate ecosystem,
// showing both available and installed agents with clear visual indicators.
//
// Returns:
//   - error: System error or listing failure
//
// Example:
//
//err := lister.ListAll()
//if err != nil {
//    return fmt.Errorf("listing failed: %w", err)
//}
func (l *ListerService) ListAll() error {
	availableChatmates, err := l.manager.GetAvailableChatmates()
	if err != nil {
		return err
	}

	installedChatmates, err := l.manager.GetInstalledChatmates()
	if err != nil {
		return err
	}

	// Create a set of installed chatmates for quick lookup
	installedSet := make(map[string]bool)
	for _, filename := range installedChatmates {
		installedSet[filename] = true
	}

	fmt.Printf("ChatMate Agents in VS Code Prompts Directory: %s\n\n", l.manager.PromptsDir)

	if len(availableChatmates) == 0 {
		fmt.Println("No chatmates available")
		return nil
	}

	// Sort chatmates for consistent display
	sort.Strings(availableChatmates)

	// Display all chatmates with installation status
	for _, filename := range availableChatmates {
		displayName := l.manager.getDisplayName(filename)
		if installedSet[filename] {
			fmt.Printf("✅ %s\n", displayName)
		} else {
			fmt.Printf("⬜ %s\n", displayName)
		}
	}

	// Summary
	installedCount := len(installedChatmates)
	availableCount := len(availableChatmates)
	fmt.Printf("\nSummary: %d/%d chatmates installed\n", installedCount, availableCount)

	return nil
}

// ListAvailable displays all available chatmate agents.
//
// This method shows only the chatmates that are available for installation,
// regardless of their current installation status.
//
// Returns:
//   - error: System error or listing failure
//
// Example:
//
//err := lister.ListAvailable()
//if err != nil {
//    return fmt.Errorf("listing available failed: %w", err)
//}
func (l *ListerService) ListAvailable() error {
	availableChatmates, err := l.manager.GetAvailableChatmates()
	if err != nil {
		return err
	}

	fmt.Println("Available ChatMate Agents:")

	if len(availableChatmates) == 0 {
		fmt.Println("No chatmates available")
		return nil
	}

	// Sort chatmates for consistent display
	sort.Strings(availableChatmates)

	// Display available chatmates
	for i, filename := range availableChatmates {
		displayName := l.manager.getDisplayName(filename)
		fmt.Printf("%d. %s\n", i+1, displayName)
	}

	fmt.Printf("\nTotal: %d chatmates available\n", len(availableChatmates))
	return nil
}

// ListInstalled displays all currently installed chatmate agents.
//
// This method shows only the chatmates that are currently installed in the
// VS Code prompts directory.
//
// Returns:
//   - error: System error or listing failure
//
// Example:
//
//err := lister.ListInstalled()
//if err != nil {
//    return fmt.Errorf("listing installed failed: %w", err)
//}
func (l *ListerService) ListInstalled() error {
	installedChatmates, err := l.manager.GetInstalledChatmates()
	if err != nil {
		return err
	}

	fmt.Printf("Installed ChatMate Agents in: %s\n", l.manager.PromptsDir)

	if len(installedChatmates) == 0 {
		fmt.Println("No chatmates are currently installed")
		return nil
	}

	// Sort chatmates for consistent display
	sort.Strings(installedChatmates)

	// Display installed chatmates
	for i, filename := range installedChatmates {
		displayName := l.manager.getDisplayName(filename)
		fmt.Printf("%d. ✅ %s\n", i+1, displayName)
	}

	fmt.Printf("\nTotal: %d chatmates installed\n", len(installedChatmates))
	return nil
}

// ListUninstalled displays chatmate agents that are available but not installed.
//
// This method shows only the chatmates that could be installed but are not
// currently present in the VS Code prompts directory.
//
// Returns:
//   - error: System error or listing failure
//
// Example:
//
//err := lister.ListUninstalled()
//if err != nil {
//    return fmt.Errorf("listing uninstalled failed: %w", err)
//}
func (l *ListerService) ListUninstalled() error {
	availableChatmates, err := l.manager.GetAvailableChatmates()
	if err != nil {
		return err
	}

	installedChatmates, err := l.manager.GetInstalledChatmates()
	if err != nil {
		return err
	}

	// Create a set of installed chatmates for quick lookup
	installedSet := make(map[string]bool)
	for _, filename := range installedChatmates {
		installedSet[filename] = true
	}

	// Find uninstalled chatmates
	var uninstalled []string
	for _, filename := range availableChatmates {
		if !installedSet[filename] {
			uninstalled = append(uninstalled, filename)
		}
	}

	fmt.Println("Uninstalled ChatMate Agents (Available for Installation):")

	if len(uninstalled) == 0 {
		fmt.Println("All available chatmates are already installed")
		return nil
	}

	// Sort chatmates for consistent display
	sort.Strings(uninstalled)

	// Display uninstalled chatmates
	for i, filename := range uninstalled {
		displayName := l.manager.getDisplayName(filename)
		fmt.Printf("%d. ⬜ %s\n", i+1, displayName)
	}

	fmt.Printf("\nTotal: %d chatmates available for installation\n", len(uninstalled))
	return nil
}

// Search finds chatmate agents matching a search term.
//
// This method searches through available chatmate display names and returns
// matches based on case-insensitive substring matching.
//
// Parameters:
//   - searchTerm: The term to search for in chatmate names
//
// Returns:
//   - error: System error or search failure
//
// Example:
//
//err := lister.Search("code")
//if err != nil {
//    return fmt.Errorf("search failed: %w", err)
//}
func (l *ListerService) Search(searchTerm string) error {
	if searchTerm == "" {
		return fmt.Errorf("search term cannot be empty")
	}

	availableChatmates, err := l.manager.GetAvailableChatmates()
	if err != nil {
		return err
	}

	installedChatmates, err := l.manager.GetInstalledChatmates()
	if err != nil {
		return err
	}

	// Create a set of installed chatmates for quick lookup
	installedSet := make(map[string]bool)
	for _, filename := range installedChatmates {
		installedSet[filename] = true
	}

	// Search for matches
	var matches []string
	searchLower := strings.ToLower(searchTerm)

	for _, filename := range availableChatmates {
		displayName := l.manager.getDisplayName(filename)
		if strings.Contains(strings.ToLower(displayName), searchLower) {
			matches = append(matches, filename)
		}
	}

	fmt.Printf("Search Results for '%s':\n", searchTerm)

	if len(matches) == 0 {
		fmt.Println("No chatmates found matching the search term")
		return nil
	}

	// Sort matches for consistent display
	sort.Strings(matches)

	// Display search results
	for i, filename := range matches {
		displayName := l.manager.getDisplayName(filename)
		status := "⬜"
		if installedSet[filename] {
			status = "✅"
		}
		fmt.Printf("%d. %s %s\n", i+1, status, displayName)
	}

	fmt.Printf("\nFound %d chatmates matching '%s'\n", len(matches), searchTerm)
	return nil
}
