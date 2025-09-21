package assertions

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// AssertFileExists asserts that a file exists at the given path
func AssertFileExists(t *testing.T, filePath string) {
	_, err := os.Stat(filePath)
	assert.NoError(t, err, "File should exist: %s", filePath)
}

// AssertFileNotExists asserts that a file does not exist at the given path
func AssertFileNotExists(t *testing.T, filePath string) {
	_, err := os.Stat(filePath)
	assert.True(t, os.IsNotExist(err), "File should not exist: %s", filePath)
}

// AssertDirectoryExists asserts that a directory exists at the given path
func AssertDirectoryExists(t *testing.T, dirPath string) {
	info, err := os.Stat(dirPath)
	require.NoError(t, err, "Directory should exist: %s", dirPath)
	assert.True(t, info.IsDir(), "Path should be a directory: %s", dirPath)
}

// AssertDirectoryNotExists asserts that a directory does not exist at the given path
func AssertDirectoryNotExists(t *testing.T, dirPath string) {
	_, err := os.Stat(dirPath)
	assert.True(t, os.IsNotExist(err), "Directory should not exist: %s", dirPath)
}

// AssertFileContains asserts that a file contains the specified content
func AssertFileContains(t *testing.T, filePath string, expectedContent string) {
	content, err := os.ReadFile(filePath)
	require.NoError(t, err, "Should be able to read file: %s", filePath)
	assert.Contains(t, string(content), expectedContent, 
		"File %s should contain expected content", filePath)
}

// AssertFileNotContains asserts that a file does not contain the specified content
func AssertFileNotContains(t *testing.T, filePath string, unexpectedContent string) {
	content, err := os.ReadFile(filePath)
	require.NoError(t, err, "Should be able to read file: %s", filePath)
	assert.NotContains(t, string(content), unexpectedContent, 
		"File %s should not contain unexpected content", filePath)
}

// AssertJSONFileValid asserts that a file contains valid JSON
func AssertJSONFileValid(t *testing.T, filePath string) {
	content, err := os.ReadFile(filePath)
	require.NoError(t, err, "Should be able to read JSON file: %s", filePath)
	
	var jsonData interface{}
	err = json.Unmarshal(content, &jsonData)
	assert.NoError(t, err, "File should contain valid JSON: %s", filePath)
}

// AssertJSONFileContains asserts that a JSON file contains a specific key-value pair
func AssertJSONFileContains(t *testing.T, filePath string, key string, expectedValue interface{}) {
	content, err := os.ReadFile(filePath)
	require.NoError(t, err, "Should be able to read JSON file: %s", filePath)
	
	var jsonData map[string]interface{}
	err = json.Unmarshal(content, &jsonData)
	require.NoError(t, err, "File should contain valid JSON: %s", filePath)
	
	value, exists := jsonData[key]
	assert.True(t, exists, "JSON should contain key '%s'", key)
	assert.Equal(t, expectedValue, value, "JSON key '%s' should have expected value", key)
}

// AssertFileEmpty asserts that a file is empty
func AssertFileEmpty(t *testing.T, filePath string) {
	content, err := os.ReadFile(filePath)
	require.NoError(t, err, "Should be able to read file: %s", filePath)
	assert.Empty(t, content, "File should be empty: %s", filePath)
}

// AssertFileNotEmpty asserts that a file is not empty
func AssertFileNotEmpty(t *testing.T, filePath string) {
	content, err := os.ReadFile(filePath)
	require.NoError(t, err, "Should be able to read file: %s", filePath)
	assert.NotEmpty(t, content, "File should not be empty: %s", filePath)
}

// AssertDirectoryEmpty asserts that a directory is empty
func AssertDirectoryEmpty(t *testing.T, dirPath string) {
	entries, err := os.ReadDir(dirPath)
	require.NoError(t, err, "Should be able to read directory: %s", dirPath)
	assert.Empty(t, entries, "Directory should be empty: %s", dirPath)
}

// AssertDirectoryNotEmpty asserts that a directory is not empty
func AssertDirectoryNotEmpty(t *testing.T, dirPath string) {
	entries, err := os.ReadDir(dirPath)
	require.NoError(t, err, "Should be able to read directory: %s", dirPath)
	assert.NotEmpty(t, entries, "Directory should not be empty: %s", dirPath)
}

// AssertDirectoryContainsFile asserts that a directory contains a specific file
func AssertDirectoryContainsFile(t *testing.T, dirPath string, fileName string) {
	filePath := filepath.Join(dirPath, fileName)
	AssertFileExists(t, filePath)
}

// AssertDirectoryContainsFiles asserts that a directory contains all specified files
func AssertDirectoryContainsFiles(t *testing.T, dirPath string, fileNames []string) {
	for _, fileName := range fileNames {
		AssertDirectoryContainsFile(t, dirPath, fileName)
	}
}

// AssertFileHasExtension asserts that a file has the expected extension
func AssertFileHasExtension(t *testing.T, filePath string, expectedExt string) {
	ext := filepath.Ext(filePath)
	assert.Equal(t, expectedExt, ext, 
		"File %s should have extension %s", filePath, expectedExt)
}

// AssertFileLineCount asserts that a file has the expected number of lines
func AssertFileLineCount(t *testing.T, filePath string, expectedLines int) {
	content, err := os.ReadFile(filePath)
	require.NoError(t, err, "Should be able to read file: %s", filePath)
	
	lines := strings.Split(string(content), "\n")
	// Remove empty last line if present
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	
	assert.Equal(t, expectedLines, len(lines), 
		"File %s should have %d lines", filePath, expectedLines)
}

// AssertStringSliceContains asserts that a slice contains the expected string
func AssertStringSliceContains(t *testing.T, slice []string, expected string) {
	assert.Contains(t, slice, expected, "Slice should contain expected string")
}

// AssertStringSliceNotContains asserts that a slice does not contain the unexpected string
func AssertStringSliceNotContains(t *testing.T, slice []string, unexpected string) {
	assert.NotContains(t, slice, unexpected, "Slice should not contain unexpected string")
}

// AssertStringSliceEqual asserts that two string slices are equal
func AssertStringSliceEqual(t *testing.T, expected []string, actual []string) {
	assert.Equal(t, expected, actual, "String slices should be equal")
}

// AssertMapContainsKey asserts that a map contains the expected key
func AssertMapContainsKey(t *testing.T, m map[string]interface{}, key string) {
	_, exists := m[key]
	assert.True(t, exists, "Map should contain key '%s'", key)
}

// AssertMapNotContainsKey asserts that a map does not contain the unexpected key
func AssertMapNotContainsKey(t *testing.T, m map[string]interface{}, key string) {
	_, exists := m[key]
	assert.False(t, exists, "Map should not contain key '%s'", key)
}

// AssertErrorContains asserts that an error contains the expected message
func AssertErrorContains(t *testing.T, err error, expectedMessage string) {
	require.Error(t, err, "Expected an error")
	assert.Contains(t, err.Error(), expectedMessage, 
		"Error should contain expected message")
}

// AssertNoError is a convenience wrapper for require.NoError with custom message
func AssertNoError(t *testing.T, err error, msgAndArgs ...interface{}) {
	require.NoError(t, err, msgAndArgs...)
}
