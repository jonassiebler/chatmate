// Package files provides chatmate-specific file naming and identification utilities.
//
// This part of the files package handles chatmate file naming conventions,
// extension management, and file type identification specific to the
// ChatMate application's .chatmode.md file format.
package files

import (
	"strings"
)

// GetChatmateNameFromFilename extracts the chatmate name from a filename.
//
// This function removes the standard ".chatmode.md" extension from chatmate
// files to get the base chatmate name. If the file doesn't have the expected
// extension, the original filename is returned unchanged.
//
// Example:
//
//	name := GetChatmateNameFromFilename("Chatmate - Code Claude Sonnet 4.chatmode.md")
//	fmt.Println(name) // Output: Chatmate - Code Claude Sonnet 4
//
//	name = GetChatmateNameFromFilename("regular-file.txt")
//	fmt.Println(name) // Output: regular-file.txt
//
// Parameters:
//   - filename: the filename to process, may include .chatmode.md extension
//
// Returns:
//   - string: the chatmate name with .chatmode.md extension removed, or original filename
func GetChatmateNameFromFilename(filename string) string {
	// Remove .chatmode.md extension if present
	if strings.HasSuffix(filename, ".chatmode.md") {
		return strings.TrimSuffix(filename, ".chatmode.md")
	}
	return filename
}

// IsChatmateFile checks if a filename is a valid chatmate file.
//
// This function determines if a filename follows the chatmate naming convention
// by checking for the ".chatmode.md" extension. This is used throughout the
// application to identify and filter chatmate files.
//
// Example:
//
//	if IsChatmateFile("Chatmate - Code Claude Sonnet 4.chatmode.md") {
//		fmt.Println("This is a chatmate file")
//	}
//
//	if !IsChatmateFile("README.md") {
//		fmt.Println("This is not a chatmate file")
//	}
//
// Parameters:
//   - filename: the filename to check for chatmate file convention
//
// Returns:
//   - bool: true if the filename ends with ".chatmode.md"
func IsChatmateFile(filename string) bool {
	return strings.HasSuffix(filename, ".chatmode.md")
}
