// Package helpers provides shared test utilities for ChatMate
package helpers

// Re-export commonly used functions for backward compatibility
import (
	"testing"

	"github.com/jonassiebler/chatmate/internal/testing/helpers/assertions"
	"github.com/jonassiebler/chatmate/internal/testing/helpers/setup"
	"github.com/jonassiebler/chatmate/internal/testing/helpers/validation"
)

// Environment setup functions
var (
	SetupTestEnvironment  = setup.SetupTestEnvironment
	MockWindowsEnv        = setup.MockWindowsEnv
	MockUnixEnv           = setup.MockUnixEnv
	GetMockChatmateDir    = setup.GetMockChatmateDir
	CreateMockChatmateDir = setup.CreateMockChatmateDir
	CreateMockVSCodeDir   = setup.CreateMockVSCodeDir
)

// Validation functions
var (
	ValidateChatmodeFile         = validation.ValidateChatmodeFile
	ValidateYAMLFrontmatter      = validation.ValidateYAMLFrontmatter
	ValidateChatmateStructure    = validation.ValidateChatmateStructure
	ValidateVSCodeStructure      = validation.ValidateVSCodeStructure
	ValidateChatmodeContent      = validation.ValidateChatmodeContent
	ValidateFilePermissions      = validation.ValidateFilePermissions
	ValidateDirectoryPermissions = validation.ValidateDirectoryPermissions
	ValidateChatmodeFilename     = validation.ValidateChatmodeFilename
	CountChatmodeFiles           = validation.CountChatmodeFiles
	ReadFileLines                = validation.ReadFileLines
)

// Assertion functions
var (
	AssertFileExists             = assertions.AssertFileExists
	AssertFileNotExists          = assertions.AssertFileNotExists
	AssertDirectoryExists        = assertions.AssertDirectoryExists
	AssertDirectoryNotExists     = assertions.AssertDirectoryNotExists
	AssertFileContains           = assertions.AssertFileContains
	AssertFileNotContains        = assertions.AssertFileNotContains
	AssertJSONFileValid          = assertions.AssertJSONFileValid
	AssertJSONFileContains       = assertions.AssertJSONFileContains
	AssertFileEmpty              = assertions.AssertFileEmpty
	AssertFileNotEmpty           = assertions.AssertFileNotEmpty
	AssertDirectoryEmpty         = assertions.AssertDirectoryEmpty
	AssertDirectoryNotEmpty      = assertions.AssertDirectoryNotEmpty
	AssertDirectoryContainsFile  = assertions.AssertDirectoryContainsFile
	AssertDirectoryContainsFiles = assertions.AssertDirectoryContainsFiles
	AssertFileHasExtension       = assertions.AssertFileHasExtension
	AssertFileLineCount          = assertions.AssertFileLineCount
	AssertStringSliceContains    = assertions.AssertStringSliceContains
	AssertStringSliceNotContains = assertions.AssertStringSliceNotContains
	AssertStringSliceEqual       = assertions.AssertStringSliceEqual
	AssertMapContainsKey         = assertions.AssertMapContainsKey
	AssertMapNotContainsKey      = assertions.AssertMapNotContainsKey
	AssertErrorContains          = assertions.AssertErrorContains
	AssertNoError                = assertions.AssertNoError
)

// Helper function to create a complete test environment with all necessary setup
func CreateCompleteTestEnvironment(t *testing.T) *setup.TestEnvironment {
	env := SetupTestEnvironment(t)

	// Ensure cleanup happens when test ends
	t.Cleanup(env.CleanupTestEnvironment)

	return env
}
