package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "🏷️  Show chatmate version information",
	Long: `Display detailed version information about chatmate including:

• Version number (semantic versioning)
• Build commit hash
• Build date and time  
• Go version used for compilation
• Target platform (OS/architecture)
• Runtime information

This information is useful for:
• Bug reports and support requests
• Verifying installation and updates
• Development and debugging
• Compliance and security audits`,
	Example: `  # Show basic version
  chatmate version
  
  # Show version in CI/automation (exit code 0)
  chatmate version --quiet
  
  # Include in bug reports
  chatmate version --full`,
	Run: func(cmd *cobra.Command, args []string) {
		quiet, _ := cmd.Flags().GetBool("quiet")
		full, _ := cmd.Flags().GetBool("full")

		if quiet {
			fmt.Printf("%s\n", version)
			return
		}

		fmt.Printf("🏷️  Chatmate Version Information\n")
		fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
		fmt.Printf("Version:      %s\n", version)
		fmt.Printf("Commit:       %s\n", commit)
		fmt.Printf("Built:        %s\n", date)

		if full {
			fmt.Printf("\n🔧 Build Information\n")
			fmt.Printf("━━━━━━━━━━━━━━━━━━━━\n")
			fmt.Printf("Go Version:   %s\n", runtime.Version())
			fmt.Printf("Platform:     %s/%s\n", runtime.GOOS, runtime.GOARCH)
			fmt.Printf("Compiler:     %s\n", runtime.Compiler)

			fmt.Printf("\n📊 Runtime Information\n")
			fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━\n")
			fmt.Printf("Goroutines:   %d\n", runtime.NumGoroutine())
			fmt.Printf("CPUs:         %d\n", runtime.NumCPU())

			var m runtime.MemStats
			runtime.ReadMemStats(&m)
			fmt.Printf("Memory:       %.2f MB\n", float64(m.Alloc)/1024/1024)
		}
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)

	versionCmd.Flags().BoolP("quiet", "q", false, "show only version number")
	versionCmd.Flags().BoolP("full", "f", false, "show full build and runtime information")
}
