package cmd

import (
	"fmt"

	"github.com/jonassiebler/chatmate/internal/manager"
	"github.com/spf13/cobra"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show ChatMate installation status and system information",
	Long: `Display comprehensive status information about your ChatMate installation,
VS Code integration, and system configuration.

🔍 System Checks:
• VS Code installation detection and version info
• ChatMate prompts directory location and permissions
• Installed vs available chatmate statistics
• Integration health and configuration validation

📊 Information Provided:
• VS Code installation path and status
• User prompts directory path and accessibility
• Count of installed/available chatmates
• System platform and environment details
• Troubleshooting hints for common issues

🎯 Use Cases:
• Verify ChatMate is properly installed and configured
• Troubleshoot installation or integration issues
• Get system information for support requests
• Check health before installing or updating chatmates

💡 Troubleshooting:
• If VS Code isn't detected, ensure it's in your PATH
• If prompts directory is missing, it will be created automatically
• Run this command after any major system or VS Code updates`,
	Example: `  # Show complete ChatMate installation status
  chatmate status
  
  # Common troubleshooting workflow
  chatmate status          # Check system health
  chatmate list           # Verify chatmate availability  
  chatmate hire --force   # Force reinstall if needed
  
  # Get status info for support requests
  chatmate status > chatmate-status.txt`,
	RunE: func(cmd *cobra.Command, args []string) error {
		chatMateManager, err := manager.NewChatMateManager()
		if err != nil {
			return fmt.Errorf("failed to initialize ChatMate manager: %w", err)
		}

		return chatMateManager.ShowStatus()
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
