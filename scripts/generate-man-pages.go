package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jonassiebler/chatmate/scripts/manpages"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <output-directory>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Example: %s ./docs/man\n", os.Args[0])
		os.Exit(1)
	}

	outputDir := os.Args[1]

	// Create a new man page generator
	generator := manpages.NewGenerator(outputDir)

	// Generate all man pages
	if err := generator.Generate(); err != nil {
		log.Fatalf("Error generating man pages: %v", err)
	}

	// List generated files
	if err := generator.ListGeneratedFiles(); err != nil {
		log.Printf("Warning: Error listing generated files: %v", err)
	}

	// Show installation instructions
	generator.ShowInstallationInstructions()
}
