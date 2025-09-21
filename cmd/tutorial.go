package cmd

import (
	"fmt"

	"github.com/jonassiebler/chatmate/cmd/tutorial"
	"github.com/spf13/cobra"
)

// tutorialCmd represents the tutorial command
var tutorialCmd = &cobra.Command{
	Use:   "tutorial [tutorial-name]",
	Short: "Interactive tutorials for learning ChatMate",
	Long: `Launch interactive tutorials to learn ChatMate features and best practices.
	
🎓 Available Tutorials:
• first-time: Complete beginner's guide to ChatMate
• daily-dev: Daily development workflow with chatmates
• team-lead: Team leadership and code review workflows
• debugging: Advanced debugging with Solve Issue chatmate
• testing: Comprehensive testing strategies with Testing chatmate

🎯 Interactive Learning:
• Step-by-step guided tutorials
• Real examples and use cases
• Interactive prompts and verification
• Best practices and tips
• Links to detailed documentation

💡 Tutorial Features:
• Hands-on practice with actual commands
• Context-aware guidance based on your setup
• Progress tracking and checkpoints
• Integration with VS Code workflows`,
	Example: `  # Start the beginner tutorial
  chatmate tutorial first-time
  
  # Learn daily development workflows
  chatmate tutorial daily-dev
  
  # Team leadership tutorial
  chatmate tutorial team-lead
  
  # Advanced debugging tutorial
  chatmate tutorial debugging
  
  # Testing best practices tutorial
  chatmate tutorial testing
  
  # List all available tutorials
  chatmate tutorial`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return listTutorials()
		}

		tutorialName := args[0]
		return runTutorial(tutorialName, tutorial.PromptToContinue)
	},
}

// listTutorials shows all available tutorials
func listTutorials() error {
	fmt.Println("📚 Available ChatMate Tutorials:")
	fmt.Println("")

	tutorials := tutorial.GetAvailableTutorials()

	for i, tut := range tutorials {
		fmt.Printf("%d. 🎓 %s (%s)\n", i+1, tut.Name, tut.Level)
		fmt.Printf("   %s\n", tut.Description)
		fmt.Printf("   ⏱️  Duration: %s\n", tut.Duration)
		fmt.Printf("   🚀 Start: chatmate tutorial %s\n", tut.Name)
		fmt.Println("")
	}

	fmt.Println("💡 Tip: Start with 'first-time' if you're new to ChatMate!")
	return nil
}

// runTutorial runs the specified tutorial
func runTutorial(name string, prompt tutorial.PromptFunc) error {
	switch name {
	case "first-time":
		return tutorial.RunFirstTimeTutorial(prompt)
	case "daily-dev":
		return tutorial.RunDailyDevTutorial(prompt)
	case "team-lead":
		return tutorial.RunTeamLeadTutorial(prompt)
	case "debugging":
		return tutorial.RunDebuggingTutorial(prompt)
	case "testing":
		return tutorial.RunTestingTutorial(prompt)
	default:
		fmt.Printf("❌ Tutorial '%s' not found.\n\n", name)
		fmt.Println("Run 'chatmate tutorial' to see available tutorials.")
		return nil
	}
}

func init() {
	rootCmd.AddCommand(tutorialCmd)
}
