package cmd

import (
	"testing"
)

// TestHireCommandExists tests that the hire command is properly defined
func TestHireCommandExists(t *testing.T) {
	if hireCmd == nil {
		t.Fatal("hire command is not defined")
	}

	if hireCmd.Use != "hire [chatmate names...]" {
		t.Errorf("Unexpected hire command use: %s", hireCmd.Use)
	}

	if hireCmd.Short == "" {
		t.Error("hire command should have a short description")
	}
}

// TestHireCommandFlags tests that required flags are defined
func TestHireCommandFlags(t *testing.T) {
	// Test that force flag exists
	forceFlag := hireCmd.Flags().Lookup("force")
	if forceFlag == nil {
		t.Error("hire command missing --force flag")
	}

	// Test that specific flag exists
	specificFlag := hireCmd.Flags().Lookup("specific")
	if specificFlag == nil {
		t.Error("hire command missing --specific flag")
	}
}
