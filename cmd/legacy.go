package cmd

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/jonassiebler/chatmate/internal/manager"
	"github.com/spf13/cobra"
)

var (
	legacyHireForce bool
)

var legacyCmd = &cobra.Command{
	Use:   "legacy",
	Short: "Manage legacy chatmates",
	Long: `Access and maintain legacy chatmates that are no longer part of the default hire flow.

Legacy chatmates remain available for specialized workflows. Use this command group to discover,
install, or clean up these agents without affecting modern chatmates.

Available subcommands:
  • legacy list   — Inspect which legacy chatmates exist and whether they are installed
  • legacy hire   — Install specific legacy chatmates on demand
  • legacy fire   — Remove installed legacy chatmates with per-mate confirmation`,
}

var legacyListCmd = &cobra.Command{
	Use:   "list",
	Short: "List legacy chatmates",
	RunE: func(cmd *cobra.Command, args []string) error {
		chatMateManager, err := manager.NewChatMateManager()
		if err != nil {
			return fmt.Errorf("failed to initialize ChatMate manager: %w", err)
		}

		legacyChatmates, err := chatMateManager.GetAvailableLegacyChatmates()
		if err != nil {
			return err
		}

		if len(legacyChatmates) == 0 {
			fmt.Println("No legacy chatmates are available")
			return nil
		}

		installed, err := chatMateManager.GetInstalledChatmates()
		if err != nil {
			return err
		}

		installedSet := make(map[string]bool)
		for _, filename := range installed {
			installedSet[filename] = true
		}

		sort.Strings(legacyChatmates)

		fmt.Println("Legacy ChatMate Agents:")
		for _, filename := range legacyChatmates {
			display := chatMateManager.DisplayName(filename)
			status := "⬜"
			if installedSet[filename] {
				status = "✅"
			}
			fmt.Printf("%s %s\n", status, display)
		}

		fmt.Printf("\nSummary: %d/%d legacy chatmates installed\n", countMatching(legacyChatmates, installedSet), len(legacyChatmates))
		return nil
	},
}

var legacyHireCmd = &cobra.Command{
	Use:   "hire [chatmate names...]",
	Short: "Install legacy chatmates manually",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("specify at least one legacy chatmate to hire (run 'chatmate legacy list' to inspect options)")
		}

		chatMateManager, err := manager.NewChatMateManager()
		if err != nil {
			return fmt.Errorf("failed to initialize ChatMate manager: %w", err)
		}

		names := make([]string, len(args))
		for i, name := range args {
			names[i] = strings.TrimSpace(name)
		}

		if err := chatMateManager.Installer().InstallLegacy(names, legacyHireForce); err != nil {
			return err
		}

		return nil
	},
}

var legacyFireCmd = &cobra.Command{
	Use:   "fire",
	Short: "Interactively remove legacy chatmates",
	RunE: func(cmd *cobra.Command, args []string) error {
		chatMateManager, err := manager.NewChatMateManager()
		if err != nil {
			return fmt.Errorf("failed to initialize ChatMate manager: %w", err)
		}

		legacyChatmates, err := chatMateManager.GetAvailableLegacyChatmates()
		if err != nil {
			return err
		}

		if len(legacyChatmates) == 0 {
			fmt.Println("No legacy chatmates are available for cleanup")
			return nil
		}

		legacySet := make(map[string]bool)
		for _, filename := range legacyChatmates {
			legacySet[filename] = true
		}

		installed, err := chatMateManager.GetInstalledChatmates()
		if err != nil {
			return err
		}

		var legacyInstalled []string
		for _, filename := range installed {
			if legacySet[filename] {
				legacyInstalled = append(legacyInstalled, filename)
			}
		}

		if len(legacyInstalled) == 0 {
			fmt.Println("No legacy chatmates are currently installed")
			return nil
		}

		sort.Strings(legacyInstalled)

		fmt.Println("🔥 Legacy ChatMate Cleanup")
		fmt.Println("Each legacy chatmate will be removed only after individual confirmation. Press Enter to skip.")

		reader := bufio.NewReader(os.Stdin)

		for _, filename := range legacyInstalled {
			display := chatMateManager.DisplayName(filename)
			fmt.Printf("Remove %s? (y/N): ", display)

			input, err := reader.ReadString('\n')
			if err != nil {
				return fmt.Errorf("failed to read input: %w", err)
			}

			response := strings.TrimSpace(input)
			if response != "y" && response != "Y" && response != "yes" && response != "YES" {
				fmt.Printf("⏭️  Skipped %s\n", display)
				continue
			}

			if err := chatMateManager.Uninstaller().UninstallChatmate(filename); err != nil {
				return err
			}
		}

		return nil
	},
}

func init() {
	legacyHireCmd.Flags().BoolVarP(&legacyHireForce, "force", "f", false, "Force reinstall even if the legacy chatmate is already installed")

	legacyCmd.AddCommand(legacyListCmd)
	legacyCmd.AddCommand(legacyHireCmd)
	legacyCmd.AddCommand(legacyFireCmd)

	rootCmd.AddCommand(legacyCmd)
}

func countMatching(filenames []string, installed map[string]bool) int {
	total := 0
	for _, name := range filenames {
		if installed[name] {
			total++
		}
	}
	return total
}
