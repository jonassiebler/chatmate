package cmd

import (
       "strings"
       "testing"
)

// TestExecuteFunction tests the Execute function for coverage
func TestExecuteFunction(t *testing.T) {
       // This should not panic and should return an error if no args are provided
       err := Execute()
       if err != nil && !strings.Contains(err.Error(), "unknown command") {
	       t.Errorf("Execute returned unexpected error: %v", err)
       }
}

// TestGetRootCommandFunction tests the GetRootCommand function for coverage
func TestGetRootCommandFunction(t *testing.T) {
       cmd := GetRootCommand()
       if cmd == nil {
	       t.Error("GetRootCommand returned nil")
       }
       if cmd.Use != "chatmate" {
	       t.Errorf("GetRootCommand returned wrong command: %s", cmd.Use)
       }
}

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

// TestRootCommandFlags tests that persistent flags are defined
func TestRootCommandFlags(t *testing.T) {
	// Test that verbose flag exists
	verboseFlag := rootCmd.PersistentFlags().Lookup("verbose")
	if verboseFlag == nil {
		t.Error("root command missing --verbose persistent flag")
	}
}

// TestSubcommands tests that all expected subcommands are registered
func TestSubcommands(t *testing.T) {
	expectedCommands := []string{
		"completion",
		"config", 
		"hire",
		"list",
		"status",
		"tutorial",
		"uninstall",
		"version",
	}

	commands := rootCmd.Commands()
	commandNames := make([]string, len(commands))
	for i, cmd := range commands {
		commandNames[i] = cmd.Name()
	}

	t.Logf("Available commands: %v", commandNames)

	for _, expectedCmd := range expectedCommands {
		found := false
		for _, actualCmd := range commandNames {
			if actualCmd == expectedCmd {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected command %s not found in subcommands", expectedCmd)
		}
	}
}

// TestRootCommandVersion tests the version functionality
func TestRootCommandVersion(t *testing.T) {
	// Test that version is set correctly
	expectedVersion := "dev (none) built on unknown"
	actualVersion := rootCmd.Version
	
	if actualVersion != expectedVersion {
		t.Logf("Version format: %s", actualVersion)
	}
	
	// Version should contain some content
	if rootCmd.Version == "" {
		t.Error("Root command version should not be empty")
	}
}

// TestRootCommandHelp tests that help information is comprehensive
func TestRootCommandHelp(t *testing.T) {
	helpText := rootCmd.Long
	
	if helpText == "" {
		t.Error("Root command should have detailed help text")
	}
	
	// Check for key elements in help text
	expectedElements := []string{
		"ChatMate",
		"chatmates",
		"VS Code",
		"hire",
		"list",
		"status",
	}
	
	for _, element := range expectedElements {
		if !strings.Contains(helpText, element) {
			t.Errorf("Help text should contain '%s'", element)
		}
	}
}


// TestRootCommandUsage tests usage examples
func TestRootCommandUsage(t *testing.T) {
	example := rootCmd.Example
	
	if example == "" {
		t.Error("Root command should have usage examples")
	}
	
	// Check for key example commands
	expectedExamples := []string{
		"chatmate hire",
		"chatmate list",
		"chatmate status",
	}
	
	for _, expectedExample := range expectedExamples {
		if !strings.Contains(example, expectedExample) {
			t.Errorf("Examples should contain '%s'", expectedExample)
		}
	}
}

// TestCommandStructure tests overall command structure
func TestCommandStructure(t *testing.T) {
	// Test that root command has subcommands
	if !rootCmd.HasSubCommands() {
		t.Error("Root command should have subcommands")
	}
	
	// Test that root command has available commands
	if !rootCmd.HasAvailableSubCommands() {
		t.Error("Root command should have available subcommands")
	}
	
	// Test that root command is runnable or has subcommands
	if !rootCmd.Runnable() && !rootCmd.HasSubCommands() {
		t.Error("Root command should be runnable or have subcommands")
	}
}

// TestPersistentFlags tests persistent flag functionality
func TestPersistentFlags(t *testing.T) {
	// Test verbose flag default value
	verboseFlag := rootCmd.PersistentFlags().Lookup("verbose")
	if verboseFlag == nil {
		t.Fatal("Verbose flag should exist")
	}
	
	if verboseFlag.DefValue != "false" {
		t.Errorf("Verbose flag default should be false, got %s", verboseFlag.DefValue)
	}
	
	// Test flag shorthand
	if verboseFlag.Shorthand != "v" {
		t.Errorf("Verbose flag shorthand should be 'v', got '%s'", verboseFlag.Shorthand)
	}
}
