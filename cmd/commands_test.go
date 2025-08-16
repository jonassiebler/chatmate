package cmd

import (
	"testing"
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
