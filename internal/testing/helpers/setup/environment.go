package setup

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestEnvironment holds test environment configuration
type TestEnvironment struct {
	TempDir     string
	MockHome    string
	MockAppData string
	OriginalEnv map[string]string
}

// SetupTestEnvironment creates a test environment with mock directories
func SetupTestEnvironment(t *testing.T) *TestEnvironment {
	tempDir := t.TempDir()
	mockHome := filepath.Join(tempDir, "home")
	mockAppData := filepath.Join(tempDir, "appdata")

	require.NoError(t, os.MkdirAll(mockHome, 0755))
	require.NoError(t, os.MkdirAll(mockAppData, 0755))

	originalEnv := make(map[string]string)
	envVars := []string{"HOME", "USERPROFILE", "APPDATA"}
	for _, env := range envVars {
		originalEnv[env] = os.Getenv(env)
	}

	// Set mock environment
	os.Setenv("HOME", mockHome)
	os.Setenv("USERPROFILE", mockHome)
	os.Setenv("APPDATA", mockAppData)

	return &TestEnvironment{
		TempDir:     tempDir,
		MockHome:    mockHome,
		MockAppData: mockAppData,
		OriginalEnv: originalEnv,
	}
}

// CleanupTestEnvironment restores original environment
func (env *TestEnvironment) CleanupTestEnvironment() {
	for key, value := range env.OriginalEnv {
		if value == "" {
			os.Unsetenv(key)
		} else {
			os.Setenv(key, value)
		}
	}
}

// MockWindowsEnv sets up Windows-specific environment for testing
func MockWindowsEnv(t *testing.T, tempDir string) {
	appData := filepath.Join(tempDir, "AppData", "Roaming")
	require.NoError(t, os.MkdirAll(appData, 0755))
	os.Setenv("APPDATA", appData)
}

// MockUnixEnv sets up Unix-specific environment for testing
func MockUnixEnv(t *testing.T, tempDir string) {
	home := filepath.Join(tempDir, "home")
	require.NoError(t, os.MkdirAll(home, 0755))
	os.Setenv("HOME", home)
}

// GetMockChatmateDir returns the expected Chatmate directory for current OS
func GetMockChatmateDir(tempDir string) string {
	if runtime.GOOS == "windows" {
		return filepath.Join(tempDir, "AppData", "Roaming", "chatmate")
	}
	return filepath.Join(tempDir, "home", ".chatmate")
}

// CreateMockChatmateDir creates a mock Chatmate directory structure
func CreateMockChatmateDir(t *testing.T, tempDir string) string {
	chatmateDir := GetMockChatmateDir(tempDir)
	require.NoError(t, os.MkdirAll(chatmateDir, 0755))
	return chatmateDir
}

// CreateMockVSCodeDir creates a mock VS Code directory structure
func CreateMockVSCodeDir(t *testing.T, tempDir string) string {
	var vscodeDir string
	if runtime.GOOS == "windows" {
		vscodeDir = filepath.Join(tempDir, "AppData", "Roaming", "Code", "User")
	} else if runtime.GOOS == "darwin" {
		vscodeDir = filepath.Join(tempDir, "home", "Library", "Application Support", "Code", "User")
	} else {
		vscodeDir = filepath.Join(tempDir, "home", ".config", "Code", "User")
	}
	require.NoError(t, os.MkdirAll(vscodeDir, 0755))
	return vscodeDir
}

// SimulateOS simulates a specific OS environment for testing
func (env *TestEnvironment) SimulateOS(t *testing.T, osType string) {
	switch osType {
	case "windows":
		MockWindowsEnv(t, env.TempDir)
	case "macos", "darwin":
		MockUnixEnv(t, env.TempDir)
	case "linux":
		MockUnixEnv(t, env.TempDir)
	default:
		t.Fatalf("Unsupported OS type for simulation: %s", osType)
	}
}

// SetupMockVSCode sets up a mock VS Code environment for the given OS
func (env *TestEnvironment) SetupMockVSCode(t *testing.T, osType string) string {
	vscodeDir := CreateMockVSCodeDir(t, env.TempDir)

	// Create prompts subdirectory
	promptsDir := filepath.Join(vscodeDir, "prompts")
	require.NoError(t, os.MkdirAll(promptsDir, 0755))

	return promptsDir
}
