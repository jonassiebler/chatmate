package manager

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestChatMateManager_GetAvailableChatmates tests retrieving available chatmates
func TestChatMateManager_GetAvailableChatmates(t *testing.T) {
	// Create temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "chatmate-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer func() {
		_ = os.RemoveAll(tmpDir)
	}()

	// Create mock mates directory
	matesDir := filepath.Join(tmpDir, "mates")
	err = os.MkdirAll(matesDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create mates directory: %v", err)
	}

	// Create test chatmate files
	testFiles := []string{
		"Test Agent 1.chatmode.md",
		"Test Agent 2.chatmode.md",
		"README.md",        // Should be ignored
		"another-file.txt", // Should be ignored
	}

	for _, file := range testFiles {
		content := "# Test Chatmate\n\nThis is a test chatmate."
		err := os.WriteFile(filepath.Join(matesDir, file), []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file %s: %v", file, err)
		}
	}

	// Create ChatMateManager with test directory
	cm := &ChatMateManager{
		MatesDir:    matesDir,
		UseEmbedded: false,
	}

	// Test GetAvailableChatmates
	chatmates, err := cm.GetAvailableChatmates()
	if err != nil {
		t.Fatalf("GetAvailableChatmates failed: %v", err)
	}

	// Should only find .chatmode.md files
	expectedCount := 2
	if len(chatmates) != expectedCount {
		t.Errorf("Expected %d chatmates, got %d", expectedCount, len(chatmates))
	}

	// Check that all returned files end with .chatmode.md
	for _, chatmate := range chatmates {
		if !strings.HasSuffix(chatmate, ".chatmode.md") {
			t.Errorf("Invalid chatmate filename: %s", chatmate)
		}
	}

	// Check that specific test files are found
	found := make(map[string]bool)
	for _, chatmate := range chatmates {
		found[chatmate] = true
	}

	expectedFiles := []string{"Test Agent 1.chatmode.md", "Test Agent 2.chatmode.md"}
	for _, expected := range expectedFiles {
		if !found[expected] {
			t.Errorf("Expected chatmate file not found: %s", expected)
		}
	}
}

// TestChatMateManager_GetInstalledChatmates tests retrieving installed chatmates
func TestChatMateManager_GetInstalledChatmates(t *testing.T) {
	// Create temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "chatmate-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	// Create mock prompts directory
	promptsDir := filepath.Join(tmpDir, "prompts")
	err = os.MkdirAll(promptsDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create prompts directory: %v", err)
	}

	// Create ChatMateManager with test directory
	cm := &ChatMateManager{
		PromptsDir: promptsDir,
	}

	// Test empty directory
	chatmates, err := cm.GetInstalledChatmates()
	if err != nil {
		t.Fatalf("GetInstalledChatmates failed: %v", err)
	}

	if len(chatmates) != 0 {
		t.Errorf("Expected 0 chatmates in empty directory, got %d", len(chatmates))
	}

	// Add some installed chatmates
	testFiles := []string{
		"Installed Agent 1.chatmode.md",
		"Installed Agent 2.chatmode.md",
		"other-file.txt", // Should be ignored
	}

	for _, file := range testFiles {
		content := "# Installed Chatmate\n\nThis is an installed chatmate."
		err := os.WriteFile(filepath.Join(promptsDir, file), []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file %s: %v", file, err)
		}
	}

	// Test with installed files
	chatmates, err = cm.GetInstalledChatmates()
	if err != nil {
		t.Fatalf("GetInstalledChatmates failed: %v", err)
	}

	// Should only find .chatmode.md files
	expectedCount := 2
	if len(chatmates) != expectedCount {
		t.Errorf("Expected %d installed chatmates, got %d", expectedCount, len(chatmates))
	}

	// Check filenames
	for _, chatmate := range chatmates {
		if !strings.HasSuffix(chatmate, ".chatmode.md") {
			t.Errorf("Invalid installed chatmate filename: %s", chatmate)
		}
	}
}

// TestChatMateManager_InstallChatmate tests installing a single chatmate
func TestChatMateManager_InstallChatmate(t *testing.T) {
	// Create temporary directories for testing
	tmpDir, err := os.MkdirTemp("", "chatmate-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	matesDir := filepath.Join(tmpDir, "mates")
	promptsDir := filepath.Join(tmpDir, "prompts")

	err = os.MkdirAll(matesDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create mates directory: %v", err)
	}

	err = os.MkdirAll(promptsDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create prompts directory: %v", err)
	}

	// Create test chatmate file
	testFile := "Test Installer.chatmode.md"
	testContent := "# Test Installer\n\nThis is a test chatmate for installation."

	err = os.WriteFile(filepath.Join(matesDir, testFile), []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test chatmate: %v", err)
	}

	// Create ChatMateManager
	cm := &ChatMateManager{
		MatesDir:    matesDir,
		PromptsDir:  promptsDir,
		UseEmbedded: false,
	}

	// Initialize services
	cm.installer = NewInstallerService(cm)

	// Test installation
	err = cm.Installer().InstallChatmate(testFile, false)
	if err != nil {
		t.Fatalf("InstallChatmate failed: %v", err)
	}

	// Verify file was installed
	installedPath := filepath.Join(promptsDir, testFile)
	if _, err := os.Stat(installedPath); os.IsNotExist(err) {
		t.Errorf("Chatmate was not installed at expected path: %s", installedPath)
	}

	// Verify content is correct
	installedContent, err := os.ReadFile(installedPath)
	if err != nil {
		t.Fatalf("Failed to read installed file: %v", err)
	}

	if string(installedContent) != testContent {
		t.Errorf("Installed content doesn't match source. Expected: %s, Got: %s", testContent, string(installedContent))
	}

	// Test installing already installed file (should skip)
	err = cm.Installer().InstallChatmate(testFile, false)
	if err != nil {
		t.Fatalf("InstallChatmate failed on already installed file: %v", err)
	}

	// Test force reinstall
	err = cm.Installer().InstallChatmate(testFile, true)
	if err != nil {
		t.Fatalf("InstallChatmate failed with force flag: %v", err)
	}
}

// TestChatMateManager_UninstallChatmate tests uninstalling a chatmate
func TestChatMateManager_UninstallChatmate(t *testing.T) {
	// Create temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "chatmate-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer func() { _ = os.RemoveAll(tmpDir) }()

	promptsDir := filepath.Join(tmpDir, "prompts")
	err = os.MkdirAll(promptsDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create prompts directory: %v", err)
	}

	// Create test installed chatmate
	testFile := "Test Uninstaller.chatmode.md"
	testContent := "# Test Uninstaller\n\nThis chatmate will be uninstalled."

	installedPath := filepath.Join(promptsDir, testFile)
	err = os.WriteFile(installedPath, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test installed chatmate: %v", err)
	}

	// Create ChatMateManager
	cm := &ChatMateManager{
		PromptsDir: promptsDir,
	}

	// Initialize services
	cm.uninstaller = NewUninstallerService(cm)

	// Verify file exists before uninstall
	if _, err := os.Stat(installedPath); os.IsNotExist(err) {
		t.Fatalf("Test setup failed: installed file doesn't exist")
	}

	// Test uninstallation
	err = cm.Uninstaller().UninstallChatmate(testFile)
	if err != nil {
		t.Fatalf("UninstallChatmate failed: %v", err)
	}

	// Verify file was removed
	if _, err := os.Stat(installedPath); !os.IsNotExist(err) {
		t.Errorf("Chatmate was not uninstalled: %s", installedPath)
	}

	// Test uninstalling non-existent file (should not error)
	err = cm.Uninstaller().UninstallChatmate("NonExistent.chatmode.md")
	if err != nil {
		t.Fatalf("UninstallChatmate failed on non-existent file: %v", err)
	}
}

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
