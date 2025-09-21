package main

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

// MainIntegrationSuite is a simple wrapper that includes the base suite
type MainIntegrationSuite struct {
	BaseIntegrationSuite
}

// TestMainIntegrationSuite runs the main integration test suite
// This maintains backward compatibility while the tests are now split into focused files
func TestMainIntegrationSuite(t *testing.T) {
	suite.Run(t, new(MainIntegrationSuite))
}
