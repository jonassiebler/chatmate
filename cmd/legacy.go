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

		installedSet := buildInstalledSet(installed)
		records := buildLegacyRecords(chatMateManager, legacyChatmates, installedSet)

		fmt.Println("Legacy ChatMate Agents:")
		installedCount := 0
		for _, record := range records {
			installed := record.canonicalInstalled || record.fallbackInstalled
			if installed {
				installedCount++
			}
			status := "⬜"
			if installed {
				status = "✅"
			}
			fmt.Printf("%s %s\n", status, record.display)
		}

		fmt.Printf("\nSummary: %d/%d legacy chatmates installed\n", installedCount, len(records))
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

		installed, err := chatMateManager.GetInstalledChatmates()
		if err != nil {
			return err
		}

		records := buildLegacyRecords(chatMateManager, legacyChatmates, buildInstalledSet(installed))

		var legacyInstalled []legacyRecord
		for _, record := range records {
			if record.canonicalInstalled || record.fallbackInstalled {
				legacyInstalled = append(legacyInstalled, record)
			}
		}

		if len(legacyInstalled) == 0 {
			fmt.Println("No legacy chatmates are currently installed")
			return nil
		}

		fmt.Println("🔥 Legacy ChatMate Cleanup")
		fmt.Println("Each legacy chatmate will be removed only after individual confirmation. Press Enter to skip.")

		reader := bufio.NewReader(os.Stdin)

		for _, record := range legacyInstalled {
			fmt.Printf("Remove %s? (y/N): ", record.display)

			input, err := reader.ReadString('\n')
			if err != nil {
				return fmt.Errorf("failed to read input: %w", err)
			}

			response := strings.TrimSpace(input)
			if response != "y" && response != "Y" && response != "yes" && response != "YES" {
				fmt.Printf("⏭️  Skipped %s\n", record.display)
				continue
			}

			if err := chatMateManager.Uninstaller().UninstallChatmate(record.filename); err != nil {
				return err
			}

			if record.fallbackFilename != record.filename {
				if err := chatMateManager.Uninstaller().UninstallChatmate(record.fallbackFilename); err != nil {
					return err
				}
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

type legacyRecord struct {
	filename           string
	display            string
	fallbackFilename   string
	canonicalInstalled bool
	fallbackInstalled  bool
}

func buildInstalledSet(installed []string) map[string]bool {
	set := make(map[string]bool)
	for _, filename := range installed {
		set[filename] = true
	}
	return set
}

func buildLegacyRecords(manager *manager.ChatMateManager, legacyFilenames []string, installed map[string]bool) []legacyRecord {
	records := make([]legacyRecord, 0, len(legacyFilenames))
	for _, filename := range legacyFilenames {
		display := manager.DisplayName(filename)
		fallback := legacyFallbackFilename(manager, filename)
		record := legacyRecord{
			filename:           filename,
			display:            display,
			fallbackFilename:   fallback,
			canonicalInstalled: installed[filename],
		}
		if fallback != filename {
			record.fallbackInstalled = installed[fallback]
		}
		records = append(records, record)
	}

	sort.Slice(records, func(i, j int) bool {
		return records[i].display < records[j].display
	})

	return records
}

func legacyFallbackFilename(manager *manager.ChatMateManager, filename string) string {
	display := strings.TrimSpace(manager.DisplayName(filename))
	if display == "" {
		return filename
	}
	fallback := strings.TrimSpace(display) + ".chatmode.md"
	if fallback == filename {
		return filename
	}
	return fallback
}
