package tutorial

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// PromptFunc defines the function signature for user interaction prompts
type PromptFunc func(string) bool

// TutorialInfo represents metadata about a tutorial
type TutorialInfo struct {
	Name        string
	Description string
	Duration    string
	Level       string
}

// ScenarioInfo represents a development scenario in tutorials
type ScenarioInfo struct {
	Title       string
	Description string
	Chatmate    string
	Example     string
	Tips        []string
}

// PromptToContinue asks the user if they want to continue and returns their response
func PromptToContinue(message string) bool {
	fmt.Printf("❓ %s [Y/n]: ", message)

	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		return false
	}

	response := strings.TrimSpace(strings.ToLower(scanner.Text()))
	return response == "" || response == "y" || response == "yes"
}
