package manager

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)


// TestNewChatMateManager tests the constructor function
func TestNewChatMateManager(t *testing.T) {
	// Save original working directory
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}
	defer func() { _ = os.Chdir(originalWd) }()

	// Create temporary directory structure
	tmpDir, err := os.MkdirTemp("", "chatmate-manager-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	// Create mates directory to simulate development environment
	matesDir := filepath.Join(tmpDir, "mates")
	err = os.MkdirAll(matesDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create mates directory: %v", err)
	}

	// Change to test directory
	err = os.Chdir(tmpDir)
	if err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	// Test NewChatMateManager
	manager, err := NewChatMateManager()
	if err != nil {
		t.Fatalf("NewChatMateManager failed: %v", err)
	}

	// Verify manager is properly initialized
	if manager == nil {
		t.Fatal("Manager is nil")
	}

	// Verify services are initialized
	if manager.Installer() == nil {
		t.Error("Installer service not initialized")
	}
	if manager.Uninstaller() == nil {
		t.Error("Uninstaller service not initialized")
	}
	if manager.Lister() == nil {
		t.Error("Lister service not initialized")
	}
	if manager.Validator() == nil {
		t.Error("Validator service not initialized")
	}
	if manager.Status() == nil {
		t.Error("Status service not initialized")
	}

	// Verify paths are set
	if manager.ScriptDir == "" {
		t.Error("ScriptDir not set")
	}
	if manager.MatesDir == "" {
		t.Error("MatesDir not set")
	}
	if manager.PromptsDir == "" {
		t.Error("PromptsDir not set")
	}

	// In development environment, should not use embedded
	if manager.UseEmbedded {
		t.Error("Should not use embedded files in development environment")
	}
}

// TestChatMateManager_ServiceAccessors tests service accessor methods
func TestChatMateManager_ServiceAccessors(t *testing.T) {
	// Create minimal manager for testing
	manager := &ChatMateManager{}
	
	// Initialize services
	manager.installer = NewInstallerService(manager)
	manager.uninstaller = NewUninstallerService(manager)
	manager.lister = NewListerService(manager)
	manager.validator = NewValidatorService(manager)
	manager.status = NewStatusService(manager)

	// Test all accessor methods
	if installer := manager.Installer(); installer == nil {
		t.Error("Installer() returned nil")
	}
	if uninstaller := manager.Uninstaller(); uninstaller == nil {
		t.Error("Uninstaller() returned nil")
	}
	if lister := manager.Lister(); lister == nil {
		t.Error("Lister() returned nil")
	}
	if validator := manager.Validator(); validator == nil {
		t.Error("Validator() returned nil")
	}
	if status := manager.Status(); status == nil {
		t.Error("Status() returned nil")
	}
}

// TestChatMateManager_GetDisplayName tests the display name function
func TestChatMateManager_GetDisplayName(t *testing.T) {
	manager := &ChatMateManager{}

	testCases := []struct {
		filename string
		expected string
	}{
		{"Test Agent.chatmode.md", "Test Agent"},
		{"Another-Agent.chatmode.md", "Another-Agent"},
		{"Complex Agent Name.chatmode.md", "Complex Agent Name"},
		{"invalid.txt", "invalid.txt"}, // Should return original for non-chatmode files
	}

	for _, tc := range testCases {
		result := manager.getDisplayName(tc.filename)
		if result != tc.expected {
			t.Errorf("getDisplayName(%q) = %q, expected %q", tc.filename, result, tc.expected)
		}
	}
}

// TestChatMateManager_EmbeddedMode tests manager behavior with embedded files
func TestChatMateManager_EmbeddedMode(t *testing.T) {
	// Create manager with embedded mode enabled
	manager := &ChatMateManager{
		UseEmbedded: true,
	}

	// Test GetAvailableChatmates with embedded files
	chatmates, err := manager.GetAvailableChatmates()
	if err != nil {
		t.Fatalf("GetAvailableChatmates failed in embedded mode: %v", err)
	}

	// Should return embedded chatmates
	if len(chatmates) == 0 {
		t.Error("No chatmates returned in embedded mode")
	}

	// Verify all returned files are chatmode files
	for _, chatmate := range chatmates {
		if !strings.HasSuffix(chatmate, ".chatmode.md") {
			t.Errorf("Invalid chatmate filename in embedded mode: %s", chatmate)
		}
	}
}

// TestListerService_NewListerService tests lister service creation
func TestListerService_NewListerService(t *testing.T) {
	manager := &ChatMateManager{}
	lister := NewListerService(manager)
	
	if lister == nil {
		t.Error("NewListerService returned nil")
	}
	if lister.manager != manager {
		t.Error("Lister service manager reference incorrect")
	}
}

// TestStatusService_NewStatusService tests status service creation
func TestStatusService_NewStatusService(t *testing.T) {
	manager := &ChatMateManager{}
	status := NewStatusService(manager)
	
	if status == nil {
		t.Error("NewStatusService returned nil")
	}
	if status.manager != manager {
		t.Error("Status service manager reference incorrect")
	}
}

// TestValidatorService_NewValidatorService tests validator service creation
func TestValidatorService_NewValidatorService(t *testing.T) {
	manager := &ChatMateManager{}
	validator := NewValidatorService(manager)
	
	if validator == nil {
		t.Error("NewValidatorService returned nil")
	}
	if validator.manager != manager {
		t.Error("Validator service manager reference incorrect")
	}
}
