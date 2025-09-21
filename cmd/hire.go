package cmd

import (
	"fmt"
	"strings"

	"github.com/jonassiebler/chatmate/internal/manager"
	"github.com/spf13/cobra"
)

var (
	hireSpecific []string
	hireForce    bool
)

// hireCmd represents the hire command
var hireCmd = &cobra.Command{
	Use:   "hire [chatmate names...]",
	Short: "Install chatmate agents for VS Code Copilot Chat",
	Long: `Install chatmate agents to enhance your VS Code Copilot Chat experience.
	
🎯 Installation Options:
• Install all available chatmates (recommended for first-time users)
• Install specific chatmates by name
• Force reinstall to update existing chatmates

📦 Available Chatmates Include:
• Solve Issue: Systematic debugging and problem resolution
• Code Review: Expert code analysis and improvement suggestions
• Testing: Comprehensive test generation and debugging
• Create PR: Pull request creation and management
• Documentation: Technical writing and API documentation

🔧 Installation Process:
1. Validates VS Code installation and prompts directory
2. Copies chatmate files to VS Code user prompts directory
3. Handles existing files with smart overwrite logic
4. Reports installation status and any conflicts

⚠️  Requirements:
• VS Code installed and accessible
• VS Code Copilot Chat extension enabled
• Write permissions to VS Code user directory`,
	Example: `  # Install all available chatmates (recommended for new users)
  chatmate hire
  
  # Install specific chatmates by name (preferred method)
  chatmate hire "Solve Issue" "Code Review" "Testing"
  
  # Install specific chatmates using flags (alternative method)
  chatmate hire --specific "Solve Issue" --specific "Testing"
  
  # Force reinstall all chatmates (useful after updates)
  chatmate hire --force
  
  # Force reinstall specific chatmates
  chatmate hire --force "Solve Issue" "Testing"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		chatMateManager, err := manager.NewChatMateManager()
		if err != nil {
			return fmt.Errorf("failed to initialize ChatMate manager: %w", err)
		}

		// Handle specific chatmates from args or --specific flag
		var specificChatmates []string
		if len(args) > 0 {
			specificChatmates = args
		} else if len(hireSpecific) > 0 {
			specificChatmates = hireSpecific
		}

		if len(specificChatmates) > 0 {
			fmt.Printf("Installing specific chatmates: %s\n", strings.Join(specificChatmates, ", "))
			return chatMateManager.Installer().InstallSpecific(specificChatmates, hireForce)
		}

		// Install all chatmates
		fmt.Println("Installing all available chatmates...")
		return chatMateManager.Installer().InstallAll(hireForce)
	},
}

func init() {
	rootCmd.AddCommand(hireCmd)

	// Add flags
	hireCmd.Flags().StringSliceVarP(&hireSpecific, "specific", "s", []string{},
		"Install specific chatmates by name (can be used multiple times)")
	hireCmd.Flags().BoolVarP(&hireForce, "force", "f", false,
		"Force reinstall even if chatmates are already installed")

	// Add some examples in the help
	hireCmd.Example = `  # Install all available chatmates
  chatmate hire

  # Install specific chatmates using flag
  chatmate hire --specific "Solve Issue" --specific "Create PR"
  
  # Install specific chatmates using arguments
  chatmate hire "Solve Issue" "Create PR"
  
  # Force reinstall all chatmates
  chatmate hire --force
  
  # Force reinstall specific chatmates
  chatmate hire --force "Code Review"`
}
