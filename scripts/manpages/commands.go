// Package manpages provides command definitions for man page generation.
package manpages

import (
	"github.com/spf13/cobra"
)

// NewRootCommand creates the root cobra command with all subcommands for man page generation.
//
// This function creates a simplified version of the ChatMate CLI structure
// focused on generating comprehensive man pages. It includes all commands
// with detailed help text and examples suitable for man page documentation.
func NewRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "chatmate",
		Short: "Open source collection of specialized AI agents for VS Code Copilot Chat",
		Long: `ChatMate is a CLI tool for managing specialized AI agents (chatmates) for VS Code Copilot Chat.
Each chatmate is a carefully crafted prompt designed to excel at specific development tasks.

🤖 What are Chatmates?
Chatmates are specialized AI agents that bring expertise in specific areas:
• Code Review: Expert code analysis and improvement suggestions
• Testing: Comprehensive test generation and debugging assistance  
• Documentation: API docs, README files, and technical writing
• Issue Resolution: Systematic debugging and problem solving
• Code Generation: Boilerplate, patterns, and architecture guidance

🚀 Quick Start:
  chatmate hire        # Install all chatmates (recommended)
  chatmate list        # View available chatmates
  chatmate status      # Check installation status

💡 Common Workflows:
  # First time setup
  chatmate hire
  
  # Check what's installed
  chatmate status
  
  # Install specific chatmates
  chatmate hire "Chatmate - Solve Issue" "Chatmate - Testing"
  
  # Remove chatmates you don't need
  chatmate uninstall "Chatmate: Create PR"`,
		Example: `  # Install all available chatmates (recommended for new users)
  chatmate hire
  
  # List available chatmates with installation status
  chatmate list
  
  # Check VS Code integration and system status
  chatmate status
  
  # Install only specific chatmates
  chatmate hire "Chatmate: Solve Issue" "Chatmate: Testing"
  
  # Remove specific chatmates
  chatmate uninstall "Chatmate - Create PR" "Chatmate - Merge PR"
  
  # Force reinstall all chatmates (useful after updates)
  chatmate hire --force
  
  # View system configuration and paths
  chatmate config`,
	}

	// Add subcommands
	rootCmd.AddCommand(
		newHireCommand(),
		newListCommand(),
		newStatusCommand(),
		newUninstallCommand(),
		newConfigCommand(),
	)

	// Add global flags
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Verbose output")

	return rootCmd
}

// newHireCommand creates the hire subcommand for man page generation.
func newHireCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "hire [chatmate names...]",
		Short: "Install chatmate agents for VS Code Copilot Chat",
		Long: `Install chatmate agents to enhance your VS Code Copilot Chat experience.
	
🎯 Installation Options:
• Install all available chatmates (recommended for first-time users)
• Install specific chatmates by name
• Force reinstall to update existing chatmates

📦 Available Chatmates Include:
• Chatmate - Solve Issue: Systematic debugging and problem resolution
• Chatmate - Review PR: Expert code analysis and improvement suggestions
• Chatmate - Testing: Comprehensive test generation and debugging
• Chatmate - Create PR: Pull request creation and management
• Chatmate - Create Issue: GitHub issue creation and management

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
  chatmate hire "Chatmate - Solve Issue" "Chatmate - Review PR" "Chatmate - Testing"
  
  # Install specific chatmates using flags (alternative method)
  chatmate hire --specific "Chatmate - Solve Issue" --specific "Chatmate - Testing"
  
  # Force reinstall all chatmates (useful after updates)
  chatmate hire --force
  
  # Force reinstall specific chatmates
  chatmate hire --force "Chatmate - Solve Issue" "Chatmate - Testing"`,
	}

	cmd.Flags().StringSlice("specific", []string{}, "Install specific chatmates by name")
	cmd.Flags().BoolP("force", "f", false, "Force reinstall existing chatmates")

	return cmd
}

// newListCommand creates the list subcommand for man page generation.
func newListCommand() *cobra.Command {
	cmd := &cobra.Command{
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
	}

	cmd.Flags().BoolP("available", "a", false, "Show only available chatmates")
	cmd.Flags().BoolP("installed", "i", false, "Show only installed chatmates")

	return cmd
}
