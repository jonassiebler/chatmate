package cmd

import (
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
