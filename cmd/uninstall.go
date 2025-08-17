package cmd

import (
	"fmt"
	"strings"

	"github.com/jonassiebler/chatmate/internal/manager"
	"github.com/spf13/cobra"
)

var (
	uninstallAll bool
)

// uninstallCmd represents the uninstall command
var uninstallCmd = &cobra.Command{
	Use:   "uninstall [chatmate names...]",
	Short: "Uninstall chatmate agents from VS Code",
	Long: `Remove chatmate agents from your VS Code Copilot Chat setup.
	
🗑️  Uninstall Options:
• Remove specific chatmates by name
• Remove all installed chatmates at once
• Safe removal with confirmation and status reporting

⚠️  What Happens:
• Chatmate files are removed from VS Code prompts directory
• Existing chat history and conversations are preserved
• You can always reinstall chatmates later with 'chatmate hire'

🔄 Common Scenarios:
• Remove chatmates you don't use to reduce clutter
• Clean up before major updates or troubleshooting
• Customize your chatmate collection for specific projects

💡 Pro Tips:
• Use 'chatmate list --installed' first to see what's available to remove
• Uninstalling doesn't affect your VS Code settings or other extensions
• You can reinstall anytime without losing functionality`,
	Example: `  # Uninstall a specific chatmate
  chatmate uninstall "Solve Issue"
  
  # Uninstall multiple chatmates at once
  chatmate uninstall "Create PR" "Merge PR" "Review PR"
  
  # Uninstall all chatmates (nuclear option)
  chatmate uninstall --all
  
  # Common workflow: check what's installed, then remove unused ones
  chatmate list --installed
  chatmate uninstall "Documentation" "Optimize Issues"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		chatMateManager, err := manager.NewChatMateManager()
		if err != nil {
			return fmt.Errorf("failed to initialize ChatMate manager: %w", err)
		}

		// Handle uninstall all flag
		if uninstallAll {
			if len(args) > 0 {
				return fmt.Errorf("cannot specify chatmate names when using --all flag")
			}
			fmt.Println("Uninstalling all chatmates...")
			return chatMateManager.UninstallAll()
		}

		// Handle specific chatmate uninstall
		if len(args) == 0 {
			return fmt.Errorf("must specify chatmate names to uninstall or use --all flag")
		}

		fmt.Printf("Uninstalling chatmates: %s\n", strings.Join(args, ", "))
		return chatMateManager.UninstallSpecific(args)
	},
}

func init() {
	rootCmd.AddCommand(uninstallCmd)

	// Add flags
	uninstallCmd.Flags().BoolVarP(&uninstallAll, "all", "a", false,
		"Uninstall all installed chatmates")

	// Add examples
	uninstallCmd.Example = `  # Uninstall a specific chatmate
  chatmate uninstall "Solve Issue"

  # Uninstall multiple chatmates
  chatmate uninstall "Solve Issue" "Create PR"
  
  # Uninstall all chatmates
  chatmate uninstall --all`
}
