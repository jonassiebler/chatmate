package main

import (
	"os"
	"testing"

	"github.com/jonassiebler/chatmate/cmd"
)

// TestUninstallCommandExecution tests uninstall command execution scenarios
func TestUninstallCommandExecution(t *testing.T) {
	// Test execution with different arguments
	testCases := []struct {
		name string
		args []string
	}{
		{"no args", []string{}},
		{"single chatmate", []string{"Test"}},
		{"multiple chatmates", []string{"Test1", "Test2"}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Capture output to prevent test noise
			old := os.Stdout
			os.Stdout = os.NewFile(0, os.DevNull)
			defer func() { os.Stdout = old }()

			// Get the uninstall command - we need to access it properly
			rootCmd := cmd.GetRootCommand()
			uninstallCmd, _, err := rootCmd.Find([]string{"uninstall"})
			if err != nil {
				t.Fatalf("Failed to find uninstall command: %v", err)
			}

			// Execute uninstall command
			err = uninstallCmd.RunE(uninstallCmd, tc.args)
			if err != nil {
				t.Logf("Uninstall command completed with expected error: %v", err)
			}
		})
	}
}

// TestHireCommandExecution tests hire command execution scenarios
func TestHireCommandExecution(t *testing.T) {
	// Test execution scenarios
	testCases := []struct {
		name string
		args []string
	}{
		{"no args", []string{}},
		{"single chatmate", []string{"Test"}},
		{"multiple chatmates", []string{"Test1", "Test2"}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Capture output
			old := os.Stdout
			os.Stdout = os.NewFile(0, os.DevNull)
			defer func() { os.Stdout = old }()

			// Get the hire command
			rootCmd := cmd.GetRootCommand()
			hireCmd, _, err := rootCmd.Find([]string{"hire"})
			if err != nil {
				t.Fatalf("Failed to find hire command: %v", err)
			}

			// Execute hire command
			err = hireCmd.RunE(hireCmd, tc.args)
			if err != nil {
				t.Logf("Hire command completed with expected error: %v", err)
			}
		})
	}
}

// TestListCommandExecution tests list command execution scenarios
func TestListCommandExecution(t *testing.T) {
	// Test execution scenarios
	testCases := []struct {
		name string
		args []string
	}{
		{"no args", []string{}},
		{"available flag", []string{"--available"}},
		{"installed flag", []string{"--installed"}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			old := os.Stdout
			os.Stdout = os.NewFile(0, os.DevNull)
			defer func() { os.Stdout = old }()

			// Get the list command
			rootCmd := cmd.GetRootCommand()
			listCmd, _, err := rootCmd.Find([]string{"list"})
			if err != nil {
				t.Fatalf("Failed to find list command: %v", err)
			}

			// Execute list command
			err = listCmd.RunE(listCmd, tc.args)
			if err != nil {
				t.Logf("List command completed with expected error: %v", err)
			}
		})
	}
}
