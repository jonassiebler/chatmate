package main

import (
	"testing"
)

// TestChatmateValidation is a simple wrapper test
// The actual validation tests have been split into focused modules:
// - test/validation/embedded_assets_test.go: Tests for embedded asset validation
// - test/validation/content_validation_test.go: Tests for chatmate content validation  
// - test/validation/helper_functions_test.go: Tests for validation helper functions
//
// This maintains backward compatibility while the tests are now organized into focused files
func TestChatmateValidation(t *testing.T) {
	t.Log("ChatMate validation tests have been modularized - see test/validation/ directory")
}
