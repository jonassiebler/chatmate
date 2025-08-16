package cmd

import (
	"testing"
)

// TestRootCommandExists tests that the root command is properly defined
func TestRootCommandExists(t *testing.T) {
	if rootCmd == nil {
		t.Fatal("root command is not defined")
	}

	if rootCmd.Use != "chatmate" {
		t.Errorf("Unexpected root command use: %s", rootCmd.Use)
	}

	if rootCmd.Short == "" {
		t.Error("root command should have a short description")
	}
}

// TestRootCommandFlags tests that global flags are defined
func TestRootCommandFlags(t *testing.T) {
	// Test that verbose flag exists
	verboseFlag := rootCmd.PersistentFlags().Lookup("verbose")
	if verboseFlag == nil {
		t.Error("root command missing --verbose flag")
	}
}

// TestSubcommands tests that all expected subcommands are registered
func TestSubcommands(t *testing.T) {
	expectedCommands := []string{"hire", "list", "uninstall", "status", "config"}

	commands := rootCmd.Commands()
	commandNames := make(map[string]bool)

	for _, cmd := range commands {
		commandNames[cmd.Name()] = true
	}

	for _, expected := range expectedCommands {
		if !commandNames[expected] {
			t.Errorf("Expected command %s not found", expected)
		}
	}

	// Log all available commands for debugging
	t.Logf("Available commands: %v", func() []string {
		var names []string
		for _, cmd := range commands {
			names = append(names, cmd.Name())
		}
		return names
	}())
}
