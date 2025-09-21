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