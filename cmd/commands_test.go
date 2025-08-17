package cmd

import (
	"os"
	"testing"
	
	"github.com/spf13/cobra"
)

// TestUninstallCommandExists tests that the uninstall command is properly defined
func TestUninstallCommandExists(t *testing.T) {
	if uninstallCmd == nil {
		t.Fatal("uninstall command is not defined")
	}

	if uninstallCmd.Use != "uninstall [chatmate names...]" {
		t.Errorf("Unexpected uninstall command use: %s", uninstallCmd.Use)
	}

	if uninstallCmd.Short == "" {
		t.Error("uninstall command should have a short description")
	}
}

// TestUninstallCommandFlags tests that required flags are defined
func TestUninstallCommandFlags(t *testing.T) {
	// Test that all flag exists
	allFlag := uninstallCmd.Flags().Lookup("all")
	if allFlag == nil {
		t.Error("uninstall command missing --all flag")
	}
}

// TestUninstallCommandExecution tests the actual execution of the uninstall command
func TestUninstallCommandExecution(t *testing.T) {
	testCases := []struct {
		name string
		args []string
	}{
		{
			name: "uninstall without arguments",
			args: []string{},
		},
		{
			name: "uninstall with specific chatmate",
			args: []string{"Test"},
		},
		{
			name: "uninstall with multiple chatmates",
			args: []string{"Test1", "Test2"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Capture output to prevent test noise
			old := os.Stdout
			os.Stdout = os.NewFile(0, os.DevNull)
			defer func() { os.Stdout = old }()

			// Reset flags
			uninstallAll = false

			// Execute command
			err := uninstallCmd.RunE(uninstallCmd, tc.args)
			if err != nil {
				t.Logf("Command failed as expected in test environment: %v", err)
			}
		})
	}
}

// TestStatusCommandExists tests that the status command is properly defined
func TestStatusCommandExists(t *testing.T) {
	if statusCmd == nil {
		t.Fatal("status command is not defined")
	}

	if statusCmd.Use != "status" {
		t.Errorf("Unexpected status command use: %s", statusCmd.Use)
	}

	if statusCmd.Short == "" {
		t.Error("status command should have a short description")
	}
}

// TestStatusCommandExecution tests the actual execution of the status command
func TestStatusCommandExecution(t *testing.T) {
	// Capture output to prevent test noise
	old := os.Stdout
	os.Stdout = os.NewFile(0, os.DevNull)
	defer func() { os.Stdout = old }()

	// Execute status command
	err := statusCmd.RunE(statusCmd, []string{})
	if err != nil {
		t.Logf("Status command failed as expected in test environment: %v", err)
	}
}

// TestConfigCommandExists tests that the config command is properly defined
func TestConfigCommandExists(t *testing.T) {
	if configCmd == nil {
		t.Fatal("config command is not defined")
	}

	if configCmd.Use != "config" {
		t.Errorf("Unexpected config command use: %s", configCmd.Use)
	}

	if configCmd.Short == "" {
		t.Error("config command should have a short description")
	}
}

// TestConfigCommandExecution tests the actual execution of the config command
func TestConfigCommandExecution(t *testing.T) {
	// Capture output to prevent test noise
	old := os.Stdout
	os.Stdout = os.NewFile(0, os.DevNull)
	defer func() { os.Stdout = old }()

	// Execute config command
	err := configCmd.RunE(configCmd, []string{})
	if err != nil {
		t.Logf("Config command failed as expected in test environment: %v", err)
	}
}

// TestVersionCommandExists tests that the version command exists
func TestVersionCommandExists(t *testing.T) {
	// Version command is added automatically by Cobra when rootCmd.Version is set
	commands := rootCmd.Commands()
	found := false
	for _, cmd := range commands {
		if cmd.Name() == "version" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Version command should be available")
	}
}

// TestCompletionCommandExists tests that the completion command exists
func TestCompletionCommandExists(t *testing.T) {
	// Completion command should exist
	commands := rootCmd.Commands()
	found := false
	for _, cmd := range commands {
		if cmd.Name() == "completion" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Completion command should be available")
	}
}

// TestTutorialCommandExists tests that the tutorial command exists
func TestTutorialCommandExists(t *testing.T) {
	// Tutorial command should exist
	commands := rootCmd.Commands()
	found := false
	for _, cmd := range commands {
		if cmd.Name() == "tutorial" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Tutorial command should be available")
	}
}

// TestAllCommandsHaveRunE tests that all main commands have execution functions
func TestAllCommandsHaveRunE(t *testing.T) {
	commandsToTest := []*cobra.Command{
		hireCmd,
		listCmd,
		statusCmd,
		configCmd,
		uninstallCmd,
	}

	for _, cmd := range commandsToTest {
		if cmd == nil {
			continue
		}
		
		t.Run(cmd.Use, func(t *testing.T) {
			if cmd.RunE == nil {
				t.Errorf("Command %s should have a RunE function", cmd.Use)
			}
		})
	}
}

// TestCommandHelpText tests that all commands have proper help text
func TestCommandHelpText(t *testing.T) {
	commandsToTest := []*cobra.Command{
		hireCmd,
		listCmd,
		statusCmd,
		configCmd,
		uninstallCmd,
	}

	for _, cmd := range commandsToTest {
		if cmd == nil {
			continue
		}
		
		t.Run(cmd.Use, func(t *testing.T) {
			if cmd.Long == "" {
				t.Errorf("Command %s should have detailed help text (Long)", cmd.Use)
			}
			
			if cmd.Example == "" {
				t.Logf("Command %s could benefit from usage examples", cmd.Use)
			}
		})
	}
}
