// Package manpages provides extended command definitions for man page generation.
package manpages

import "github.com/spf13/cobra"

// newStatusCommand creates the status subcommand for man page generation.
func newStatusCommand() *cobra.Command {
	return &cobra.Command{
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
}

// newUninstallCommand creates the uninstall subcommand for man page generation.
func newUninstallCommand() *cobra.Command {
	cmd := &cobra.Command{
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
  chatmate uninstall "Chatmate - Solve Issue"
  
  # Uninstall multiple chatmates at once
  chatmate uninstall "Chatmate - Create PR" "Chatmate - Merge PR" "Chatmate - Review PR"
  
  # Uninstall all chatmates (nuclear option)
  chatmate uninstall --all
  
  # Common workflow: check what's installed, then remove unused ones
  chatmate list --installed
  chatmate uninstall "Chatmate - Create Release" "Chatmate - Optimize Issues"`,
	}

	cmd.Flags().Bool("all", false, "Uninstall all chatmates")

	return cmd
}

// newConfigCommand creates the config subcommand for man page generation.
func newConfigCommand() *cobra.Command {
	cmd := &cobra.Command{
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

	cmd.Flags().BoolP("show", "s", true, "Show configuration details")

	return cmd
}
