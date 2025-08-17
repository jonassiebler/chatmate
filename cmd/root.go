package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
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
  chatmate hire "Solve Issue" "Testing"
  
  # Remove chatmates you don't need
  chatmate uninstall "Create PR"`,
	Example: `  # Install all available chatmates (recommended for new users)
  chatmate hire
  
  # List available chatmates with installation status
  chatmate list
  
  # Check VS Code integration and system status
  chatmate status
  
  # Install only specific chatmates
  chatmate hire "Solve Issue" "Code Review" "Testing"
  
  # Remove specific chatmates
  chatmate uninstall "Create PR" "Merge PR"
  
  # Force reinstall all chatmates (useful after updates)
  chatmate hire --force
  
  # View system configuration and paths
  chatmate config`,
	Version: fmt.Sprintf("%s (%s) built on %s", version, commit, date),
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	return rootCmd.Execute()
}

// GetRootCommand returns the root command for testing purposes
func GetRootCommand() *cobra.Command {
	return rootCmd
}

func init() {
	// Global flags can be added here
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")
}
