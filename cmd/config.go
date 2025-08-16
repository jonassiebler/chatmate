package cmd

import (
	"fmt"

	"github.com/jonassiebler/chatmate/internal/manager"
	"github.com/spf13/cobra"
)

var (
	configShow bool
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Show ChatMate configuration settings and system paths",
	Long: `Display detailed ChatMate configuration information including file paths,
platform details, and system environment settings.

🔧 Configuration Details:
• ChatMate installation directory and embedded resources
• VS Code user directory and prompts path  
• Platform-specific paths and conventions
• Environment variables and system settings
• File permissions and accessibility information

🎯 Use Cases:
• Debug installation or path-related issues
• Understand where ChatMate stores and finds files
• Verify system environment before troubleshooting
• Get technical details for support or development
• Validate cross-platform compatibility

💡 Technical Information:
• Shows both embedded and external chatmate locations
• Displays resolved file paths with expansion
• Indicates which paths are accessible and writable
• Platform-specific directory conventions (Windows/macOS/Linux)`,
	Example: `  # Show complete configuration information
  chatmate config
  
  # Save configuration for support requests
  chatmate config > chatmate-config.txt
  
  # Common troubleshooting workflow
  chatmate config    # Check paths and configuration
  chatmate status    # Verify system integration
  chatmate list      # Test chatmate discovery`,
	RunE: func(cmd *cobra.Command, args []string) error {
		chatMateManager, err := manager.NewChatMateManager()
		if err != nil {
			return fmt.Errorf("failed to initialize ChatMate manager: %w", err)
		}

		// For now, we only support showing config
		// In the future, we could add config management features
		chatMateManager.ShowConfig()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(configCmd)

	// Add flags for future extensibility
	configCmd.Flags().BoolVarP(&configShow, "show", "s", true,
		"Show current configuration (default)")

	// Hidden flag for future extension
	_ = configCmd.Flags().MarkHidden("show") // Add examples
	configCmd.Example = `  # Show current ChatMate configuration
  chatmate config`
}
