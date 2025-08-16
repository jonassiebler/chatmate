package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <output-directory>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Example: %s ./docs/man\n", os.Args[0])
		os.Exit(1)
	}

	outputDir := os.Args[1]

	// Ensure the output directory exists
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Fatalf("Error creating output directory %s: %v", outputDir, err)
	}

	// Get the root command
	rootCmd := getRootCommand()

	// Set additional information for man pages
	rootCmd.DisableAutoGenTag = true

	// Generate man pages for all commands
	fmt.Printf("Generating man pages to %s...\n", outputDir)

	header := &doc.GenManHeader{
		Title:   "ChatMate",
		Section: "1",
		Source:  "ChatMate CLI",
		Manual:  "ChatMate Manual",
	}

	// Generate the main man page and all subcommand man pages
	err := doc.GenManTree(rootCmd, header, outputDir)
	if err != nil {
		log.Fatalf("Error generating man pages: %v", err)
	}

	// Also generate individual man pages for each subcommand
	fmt.Println("Generating individual subcommand man pages...")
	for _, subCmd := range rootCmd.Commands() {
		if subCmd.Hidden {
			continue
		}

		subHeader := &doc.GenManHeader{
			Title:   fmt.Sprintf("chatmate-%s", subCmd.Name()),
			Section: "1",
			Source:  "ChatMate CLI",
			Manual:  "ChatMate Manual",
		}

		subCmdFile := filepath.Join(outputDir, fmt.Sprintf("chatmate-%s.1", subCmd.Name()))
		err := doc.GenManTreeFromOpts(subCmd, doc.GenManTreeOptions{
			Header:           subHeader,
			Path:             outputDir,
			CommandSeparator: "-",
		})
		if err != nil {
			log.Printf("Warning: Error generating man page for %s: %v", subCmd.Name(), err)
		} else {
			fmt.Printf("  Generated %s\n", filepath.Base(subCmdFile))
		}
	}

	fmt.Println("✅ Man pages generated successfully!")

	// List generated files
	fmt.Println("\nGenerated man pages:")
	err = filepath.Walk(outputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".1" {
			relPath, _ := filepath.Rel(outputDir, path)
			fmt.Printf("  %s\n", relPath)
		}
		return nil
	})
	if err != nil {
		log.Printf("Warning: Error listing generated files: %v", err)
	}

	fmt.Printf("\nTo install man pages system-wide (requires sudo):\n")
	fmt.Printf("  sudo cp %s/*.1 /usr/local/share/man/man1/\n", outputDir)
	fmt.Printf("  sudo mandb  # Update man database\n")
	fmt.Printf("\nTo view a man page:\n")
	fmt.Printf("  man chatmate\n")
	fmt.Printf("  man chatmate-hire\n")
}

// getRootCommand returns the root cobra command
// This is a simplified version that doesn't include the full CLI logic
func getRootCommand() *cobra.Command {
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
	}

	// Add subcommands similar to the actual CLI
	hireCmd := &cobra.Command{
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
	}

	listCmd := &cobra.Command{
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

	statusCmd := &cobra.Command{
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
	}

	uninstallCmd := &cobra.Command{
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
	}

	configCmd := &cobra.Command{
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
	}

	// Add flags to commands
	hireCmd.Flags().StringSlice("specific", []string{}, "Install specific chatmates by name")
	hireCmd.Flags().BoolP("force", "f", false, "Force reinstall existing chatmates")

	listCmd.Flags().BoolP("available", "a", false, "Show only available chatmates")
	listCmd.Flags().BoolP("installed", "i", false, "Show only installed chatmates")

	uninstallCmd.Flags().Bool("all", false, "Uninstall all chatmates")

	configCmd.Flags().BoolP("show", "s", true, "Show configuration details")

	// Add global flags
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Verbose output")

	// Add subcommands
	rootCmd.AddCommand(hireCmd, listCmd, statusCmd, uninstallCmd, configCmd)

	return rootCmd
}
