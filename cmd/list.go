package cmd

import (
	"fmt"

	"github.com/jonassiebler/chatmate/internal/manager"
	"github.com/spf13/cobra"
)

var (
	listAvailable bool
	listInstalled bool
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available and installed chatmate agents",
	Long: `Display comprehensive information about available and installed chatmate agents.
	
📋 What You'll See:
• Available chatmates with descriptions and specializations  
• Installation status (✅ installed, ❌ not installed)
• Summary statistics of your chatmate collection

🎯 Filter Options:
• Show only available chatmates (--available)
• Show only installed chatmates (--installed)  
• Default: Show both available and installed with status indicators

💡 Use Cases:
• Discover new chatmates to install
• Check installation status of specific chatmates
• Get overview of your current chatmate setup
• Find chatmates by their specialization areas`,
	Example: `  # List all chatmates with installation status (default)
  chatmate list
  
  # Show only available chatmates (not yet installed)
  chatmate list --available
  
  # Show only installed chatmates
  chatmate list --installed
  
  # Combine with other commands for workflows
  chatmate list --available | grep "Testing"  # Find testing-related chatmates`,
	RunE: func(cmd *cobra.Command, args []string) error {
		chatMateManager, err := manager.NewChatMateManager()
		if err != nil {
			return fmt.Errorf("failed to initialize ChatMate manager: %w", err)
		}

		// Determine what to show based on flags
		showAvailable := true
		showInstalled := true

		// If specific flags are set, only show those
		if listAvailable || listInstalled {
			showAvailable = listAvailable
			showInstalled = listInstalled
		}

		return chatMateManager.ListChatmates(showAvailable, showInstalled)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Add flags
	listCmd.Flags().BoolVarP(&listAvailable, "available", "a", false,
		"Show only available chatmates")
	listCmd.Flags().BoolVarP(&listInstalled, "installed", "i", false,
		"Show only installed chatmates")

	// Add examples
	listCmd.Example = `  # List all chatmates (available and installed)
  chatmate list

  # List only available chatmates
  chatmate list --available
  
  # List only installed chatmates
  chatmate list --installed`
}
