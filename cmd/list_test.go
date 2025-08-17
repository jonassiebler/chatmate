package cmd

import (
	"os"
	"testing"
)

// TestListCommandExists tests that the list command is properly defined
func TestListCommandExists(t *testing.T) {
	if listCmd == nil {
		t.Fatal("list command is not defined")
	}

	if listCmd.Use != "list" {
		t.Errorf("Unexpected list command use: %s", listCmd.Use)
	}

	if listCmd.Short == "" {
		t.Error("list command should have a short description")
	}
}

// TestListCommandFlags tests that required flags are defined
func TestListCommandFlags(t *testing.T) {
	// Test that available flag exists
	availableFlag := listCmd.Flags().Lookup("available")
	if availableFlag == nil {
		t.Error("list command missing --available flag")
	}

	// Test that installed flag exists
	installedFlag := listCmd.Flags().Lookup("installed")
	if installedFlag == nil {
		t.Error("list command missing --installed flag")
	}
}

// TestListCommandExecution tests the actual execution of the list command
func TestListCommandExecution(t *testing.T) {
	testCases := []struct {
		name string
		args []string
	}{
		{
			name: "list without flags",
			args: []string{},
		},
		{
			name: "list available",
			args: []string{"--available"},
		},
		{
			name: "list installed",
			args: []string{"--installed"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Capture output to prevent test noise
			old := os.Stdout
			os.Stdout = os.NewFile(0, os.DevNull)
			defer func() { os.Stdout = old }()

			// Reset flags
			listAvailable = false
			listInstalled = false

			// Parse flags if provided
			if len(tc.args) > 0 {
				err := listCmd.ParseFlags(tc.args)
				if err != nil {
					t.Fatalf("Failed to parse flags: %v", err)
				}
			}

			// Execute command
			err := listCmd.RunE(listCmd, []string{})
			if err != nil {
				t.Logf("Command failed as expected in test environment: %v", err)
			}

			// The command should attempt to execute, even if it fails due to missing setup
		})
	}
}

// TestListCommandFlagBehavior tests flag behavior
func TestListCommandFlagBehavior(t *testing.T) {
	t.Run("available flag sets variable", func(t *testing.T) {
		// Reset flags
		listAvailable = false
		listInstalled = false

		// Parse available flag
		err := listCmd.ParseFlags([]string{"--available"})
		if err != nil {
			t.Fatalf("Failed to parse available flag: %v", err)
		}

		if !listAvailable {
			t.Error("Available flag should be true")
		}
	})

	t.Run("installed flag sets variable", func(t *testing.T) {
		// Reset flags
		listAvailable = false
		listInstalled = false

		// Parse installed flag
		err := listCmd.ParseFlags([]string{"--installed"})
		if err != nil {
			t.Fatalf("Failed to parse installed flag: %v", err)
		}

		if !listInstalled {
			t.Error("Installed flag should be true")
		}
	})

	t.Run("both flags can be set", func(t *testing.T) {
		// Reset flags
		listAvailable = false
		listInstalled = false

		// Parse both flags
		err := listCmd.ParseFlags([]string{"--available", "--installed"})
		if err != nil {
			t.Fatalf("Failed to parse both flags: %v", err)
		}

		if !listAvailable {
			t.Error("Available flag should be true")
		}
		if !listInstalled {
			t.Error("Installed flag should be true")
		}
	})
}
