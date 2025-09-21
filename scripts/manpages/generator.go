// Package manpages provides functionality for generating man pages for ChatMate CLI commands.
package manpages

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra/doc"
)

// Generator handles the generation of man pages for ChatMate commands.
type Generator struct {
	outputDir string
}

// NewGenerator creates a new man page generator with the specified output directory.
func NewGenerator(outputDir string) *Generator {
	return &Generator{
		outputDir: outputDir,
	}
}

// Generate creates man pages for all ChatMate commands.
//
// This method generates comprehensive man pages including:
//   - Main chatmate man page with overview and common usage
//   - Individual man pages for each subcommand (hire, list, status, etc.)
//   - Proper man page formatting with sections and cross-references
//
// Returns:
//   - error: Generation failure or file system error
func (g *Generator) Generate() error {
	// Ensure the output directory exists
	if err := os.MkdirAll(g.outputDir, 0755); err != nil {
		return fmt.Errorf("error creating output directory %s: %w", g.outputDir, err)
	}

	// Get the root command with all subcommands
	rootCmd := NewRootCommand()

	// Set additional information for man pages
	rootCmd.DisableAutoGenTag = true

	// Generate man pages for all commands
	fmt.Printf("Generating man pages to %s...\n", g.outputDir)

	header := &doc.GenManHeader{
		Title:   "ChatMate",
		Section: "1",
		Source:  "ChatMate CLI",
		Manual:  "ChatMate Manual",
	}

	// Generate the main man page and all subcommand man pages
	err := doc.GenManTree(rootCmd, header, g.outputDir)
	if err != nil {
		return fmt.Errorf("error generating man pages: %w", err)
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

		subCmdFile := filepath.Join(g.outputDir, fmt.Sprintf("chatmate-%s.1", subCmd.Name()))
		err := doc.GenManTreeFromOpts(subCmd, doc.GenManTreeOptions{
			Header:           subHeader,
			Path:             g.outputDir,
			CommandSeparator: "-",
		})
		if err != nil {
			log.Printf("Warning: Error generating man page for %s: %v", subCmd.Name(), err)
		} else {
			fmt.Printf("  Generated %s\n", filepath.Base(subCmdFile))
		}
	}

	fmt.Println("✅ Man pages generated successfully!")

	return nil
}

// ListGeneratedFiles displays all generated man page files.
//
// This method scans the output directory and lists all .1 (man page) files
// that were generated, providing feedback to the user about what was created.
func (g *Generator) ListGeneratedFiles() error {
	fmt.Println("\nGenerated man pages:")
	err := filepath.Walk(g.outputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".1" {
			relPath, _ := filepath.Rel(g.outputDir, path)
			fmt.Printf("  %s\n", relPath)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("error listing generated files: %w", err)
	}

	return nil
}

// ShowInstallationInstructions displays helpful instructions for installing and using man pages.
func (g *Generator) ShowInstallationInstructions() {
	fmt.Printf("\nTo install man pages system-wide (requires sudo):\n")
	fmt.Printf("  sudo cp %s/*.1 /usr/local/share/man/man1/\n", g.outputDir)
	fmt.Printf("  sudo mandb  # Update man database\n")
	fmt.Printf("\nTo view a man page:\n")
	fmt.Printf("  man chatmate\n")
	fmt.Printf("  man chatmate-hire\n")
}
