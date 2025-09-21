package main

import (
	"os"
	"os/exec"
	"runtime"

	"github.com/jonassiebler/chatmate/internal/testing/helpers/setup"
	"github.com/stretchr/testify/suite"
)

// BaseIntegrationSuite provides common setup for integration tests
type BaseIntegrationSuite struct {
	suite.Suite
	env        *setup.TestEnvironment
	binaryPath string
}

// SetupSuite builds the test binary once for all integration tests
func (s *BaseIntegrationSuite) SetupSuite() {
	// Build test binary from project root
	s.binaryPath = "chatmate-integration-test"
	if runtime.GOOS == "windows" {
		s.binaryPath += ".exe"
	}

	cmd := exec.Command("go", "build", "-o", s.binaryPath, "github.com/jonassiebler/chatmate")
	output, err := cmd.CombinedOutput()
	s.Require().NoError(err, "Failed to build test binary: %s", string(output))
}

// TearDownSuite cleans up the test binary
func (s *BaseIntegrationSuite) TearDownSuite() {
	if s.binaryPath != "" {
		os.Remove(s.binaryPath)
	}
}

// SetupTest creates a fresh test environment for each test
func (s *BaseIntegrationSuite) SetupTest() {
	s.env = setup.SetupTestEnvironment(s.T())
}

// GetBinaryPath returns the path to the test binary
func (s *BaseIntegrationSuite) GetBinaryPath() string {
	return s.binaryPath
}

// GetTestEnvironment returns the test environment
func (s *BaseIntegrationSuite) GetTestEnvironment() *setup.TestEnvironment {
	return s.env
}
