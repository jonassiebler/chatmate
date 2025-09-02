package cmd

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func TestRunDailyDevTutorial(t *testing.T) {
	mockPrompt := func(msg string) bool { return false }
	output := captureOutput(func() {
		err := runDailyDevTutorial(mockPrompt)
		if err != nil {
			t.Errorf("runDailyDevTutorial returned error: %v", err)
		}
	})
	if !strings.Contains(output, "Daily Development Workflow Tutorial") {
		t.Errorf("Expected daily dev tutorial output, got: %s", output)
	}
}

func TestRunTeamLeadTutorial(t *testing.T) {
	mockPrompt := func(msg string) bool { return false }
	output := captureOutput(func() {
		err := runTeamLeadTutorial(mockPrompt)
		if err != nil {
			t.Errorf("runTeamLeadTutorial returned error: %v", err)
		}
	})
	if !strings.Contains(output, "Team Leadership Tutorial") {
		t.Errorf("Expected team lead tutorial output, got: %s", output)
	}
}

func TestRunDebuggingTutorial(t *testing.T) {
	mockPrompt := func(msg string) bool { return false }
	output := captureOutput(func() {
		err := runDebuggingTutorial(mockPrompt)
		if err != nil {
			t.Errorf("runDebuggingTutorial returned error: %v", err)
		}
	})
	if !strings.Contains(output, "Debugging with Solve Issue Chatmate") {
		t.Errorf("Expected debugging tutorial output, got: %s", output)
	}
}

func TestRunTestingTutorial(t *testing.T) {
	mockPrompt := func(msg string) bool { return false }
	output := captureOutput(func() {
		err := runTestingTutorial(mockPrompt)
		if err != nil {
			t.Errorf("runTestingTutorial returned error: %v", err)
		}
	})
	if !strings.Contains(output, "Testing with Testing Chatmate") {
		t.Errorf("Expected testing tutorial output, got: %s", output)
	}
}

// helper to capture stdout
func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func TestRunTutorial_FirstTime(t *testing.T) {
	// Use a mock prompt function that always returns false
	mockPrompt := func(msg string) bool { return false }

	output := captureOutput(func() {
		err := runTutorial("first-time", mockPrompt)
		if err != nil {
			t.Errorf("runTutorial returned error: %v", err)
		}
	})

	if !strings.Contains(output, "First Time User Tutorial") {
		t.Errorf("Expected tutorial output, got: %s", output)
	}
}

func TestRunTutorial_Unknown(t *testing.T) {
	output := captureOutput(func() {
		err := runTutorial("unknown-tutorial", nil)
		if err != nil {
			t.Errorf("runTutorial returned error: %v", err)
		}
	})

	if !strings.Contains(output, "not found") {
		t.Errorf("Expected not found message, got: %s", output)
	}
}
