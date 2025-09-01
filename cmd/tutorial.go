package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/jonassiebler/chatmate/internal/manager"
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
		return runTutorial(tutorialName, promptToContinue)
	},
}

// listTutorials shows all available tutorials
func listTutorials() error {
	fmt.Println("📚 Available ChatMate Tutorials:")
	fmt.Println("")

	tutorials := []struct {
		Name        string
		Description string
		Duration    string
		Level       string
	}{
		{
			Name:        "first-time",
			Description: "Complete beginner's guide to ChatMate installation and basic usage",
			Duration:    "10-15 minutes",
			Level:       "Beginner",
		},
		{
			Name:        "daily-dev",
			Description: "Daily development workflow with chatmates for coding tasks",
			Duration:    "15-20 minutes",
			Level:       "Intermediate",
		},
		{
			Name:        "team-lead",
			Description: "Team leadership workflows: code reviews, PR management, issue creation",
			Duration:    "20-25 minutes",
			Level:       "Advanced",
		},
		{
			Name:        "debugging",
			Description: "Advanced debugging techniques with the Solve Issue chatmate",
			Duration:    "15-20 minutes",
			Level:       "Intermediate",
		},
		{
			Name:        "testing",
			Description: "Comprehensive testing strategies with the Testing chatmate",
			Duration:    "15-20 minutes",
			Level:       "Intermediate",
		},
	}

	for i, tutorial := range tutorials {
		fmt.Printf("%d. 🎓 %s (%s)\n", i+1, tutorial.Name, tutorial.Level)
		fmt.Printf("   %s\n", tutorial.Description)
		fmt.Printf("   ⏱️  Duration: %s\n", tutorial.Duration)
		fmt.Printf("   🚀 Start: chatmate tutorial %s\n", tutorial.Name)
		fmt.Println("")
	}

	fmt.Println("💡 Tip: Start with 'first-time' if you're new to ChatMate!")
	return nil
}

// runTutorial runs the specified tutorial
type PromptFunc func(string) bool

func runTutorial(name string, prompt PromptFunc) error {
	switch name {
	case "first-time":
		return runFirstTimeTutorial(prompt)
	case "daily-dev":
		return runDailyDevTutorial(prompt)
	case "team-lead":
		return runTeamLeadTutorial(prompt)
	case "debugging":
		return runDebuggingTutorial(prompt)
	case "testing":
		return runTestingTutorial(prompt)
	default:
		fmt.Printf("❌ Tutorial '%s' not found.\n\n", name)
		fmt.Println("Run 'chatmate tutorial' to see available tutorials.")
		return nil
	}
}

// runFirstTimeTutorial runs the beginner tutorial
func runFirstTimeTutorial(prompt PromptFunc) error {
	fmt.Println("🎓 Welcome to ChatMate - First Time User Tutorial!")
	fmt.Println("=================================================")
	fmt.Println("")

	// Step 1: Introduction
	fmt.Println("📚 Step 1: Understanding ChatMate")
	fmt.Println("ChatMate provides specialized AI agents (chatmates) for VS Code Copilot Chat.")
	fmt.Println("Each chatmate is an expert in specific development tasks.")
	fmt.Println("")

	if !prompt("Ready to learn about chatmates?") {
		return nil
	}

	// Step 2: Check system status
	fmt.Println("🔍 Step 2: Check Your System Status")
	fmt.Println("Let's verify your system is ready for ChatMate...")
	fmt.Println("")

	fmt.Println("Running: chatmate status")

	chatMateManager, err := manager.NewChatMateManager()
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		fmt.Println("Please resolve this issue before continuing the tutorial.")
		return nil
	}

	err = chatMateManager.ShowStatus()
	if err != nil {
		fmt.Printf("❌ Error showing status: %v\n", err)
		return nil
	}

	fmt.Println("")
	if !prompt("Does your status look good? (VS Code detected, prompts directory accessible)") {
		fmt.Println("💡 If you see issues, run 'chatmate config' for more details or check the troubleshooting guide.")
		return nil
	}

	// Step 3: Install chatmates
	fmt.Println("📦 Step 3: Install Your First Chatmates")
	fmt.Println("We'll install some essential chatmates to get you started.")
	fmt.Println("")

	fmt.Println("Recommended chatmates for beginners:")
	fmt.Println("• Solve Issue: For debugging and problem-solving")
	fmt.Println("• Code Review: For code analysis and improvements")
	fmt.Println("• Testing: For test generation and debugging")
	fmt.Println("")

	if promptToContinue("Would you like to install these recommended chatmates?") {
		fmt.Println("Running: chatmate hire \"Solve Issue\" \"Code Review\" \"Testing\"")
		fmt.Println("")

		err = chatMateManager.InstallSpecific([]string{"Solve Issue", "Code Review", "Testing"}, false)
		if err != nil {
			fmt.Printf("❌ Error installing chatmates: %v\n", err)
			return nil
		}

		fmt.Println("✅ Chatmates installed successfully!")
		fmt.Println("")
	}

	// Step 4: Verify installation
	fmt.Println("✅ Step 4: Verify Installation")
	fmt.Println("Let's check what chatmates are now installed...")
	fmt.Println("")

	fmt.Println("Running: chatmate list --installed")
	err = chatMateManager.ListChatmates(false, true)
	if err != nil {
		fmt.Printf("❌ Error listing chatmates: %v\n", err)
		return nil
	}

	fmt.Println("")
	if !promptToContinue("Do you see your installed chatmates listed above?") {
		fmt.Println("💡 If chatmates aren't showing, try running 'chatmate hire --force' to reinstall.")
		return nil
	}

	// Step 5: VS Code integration
	fmt.Println("🎯 Step 5: Using Chatmates in VS Code")
	fmt.Println("Now comes the exciting part - using your chatmates!")
	fmt.Println("")

	fmt.Println("To use chatmates in VS Code:")
	fmt.Println("1. 🔄 RESTART VS Code completely (close all windows, reopen)")
	fmt.Println("2. 💬 Open Copilot Chat (Ctrl/Cmd+Shift+P → 'Chat: Open Chat')")
	fmt.Println("3. 🤖 Use @ to mention chatmates: '@Solve Issue', '@Code Review', '@Testing'")
	fmt.Println("")

	fmt.Println("Example conversations:")
	fmt.Println("• '@Solve Issue My React component won't render properly'")
	fmt.Println("• '@Code Review Check this authentication function for security'")
	fmt.Println("• '@Testing Generate unit tests for this service class'")
	fmt.Println("")

	if !promptToContinue("Ready to try this in VS Code?") {
		return nil
	}

	// Step 6: Next steps
	fmt.Println("🚀 Step 6: Next Steps")
	fmt.Println("Congratulations! You've completed the ChatMate first-time tutorial!")
	fmt.Println("")

	fmt.Println("What to do next:")
	fmt.Println("1. 🔄 Restart VS Code and try your new chatmates")
	fmt.Println("2. 📖 Read the User Guide: docs/USER_GUIDE.md")
	fmt.Println("3. 🎓 Try more tutorials: chatmate tutorial daily-dev")
	fmt.Println("4. 🔧 Explore all commands: chatmate --help")
	fmt.Println("")

	fmt.Println("💡 Pro Tips:")
	fmt.Println("• Use 'chatmate status' to check system health anytime")
	fmt.Println("• Use 'chatmate hire --force' to update existing chatmates")
	fmt.Println("• Use 'chatmate list' to see all available chatmates")
	fmt.Println("• Join GitHub Discussions for community support")
	fmt.Println("")

	fmt.Println("🎉 Happy coding with your new chatmates!")
	return nil
}

// runDailyDevTutorial runs the daily development workflow tutorial
func runDailyDevTutorial(prompt PromptFunc) error {
	fmt.Println("💻 ChatMate Daily Development Workflow Tutorial")
	fmt.Println("===============================================")
	fmt.Println("")

	fmt.Println("This tutorial shows you how to integrate ChatMate into your daily development routine.")
	fmt.Println("")

	if !prompt("Ready to learn daily development workflows?") {
		return nil
	}

	// Morning routine
	fmt.Println("🌅 Morning Routine: Health Check")
	fmt.Println("Start your day by checking ChatMate status:")
	fmt.Println("")
	fmt.Println("$ chatmate status    # Check system health")
	fmt.Println("$ chatmate list      # Review available chatmates")
	fmt.Println("")

	if !prompt("Let's run a quick health check now:") {
		return nil
	}

	chatMateManager, err := manager.NewChatMateManager()
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return nil
	}

	err = chatMateManager.ShowStatus()
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return nil
	}

	fmt.Println("")

	// Development workflow scenarios
	scenarios := []struct {
		Title       string
		Description string
		Chatmate    string
		Example     string
		Tips        []string
	}{
		{
			Title:       "🐛 Debugging Session",
			Description: "When you encounter a bug or issue",
			Chatmate:    "Solve Issue",
			Example:     "@Solve Issue My login form validation isn't working correctly. Users can submit empty passwords.",
			Tips: []string{
				"Provide context: error messages, relevant code, steps to reproduce",
				"Include environment details: browser, OS, framework versions",
				"Ask follow-up questions to narrow down the root cause",
			},
		},
		{
			Title:       "👁️ Code Review",
			Description: "Before committing changes or during peer review",
			Chatmate:    "Code Review",
			Example:     "@Code Review Please review this authentication middleware for security issues and best practices.",
			Tips: []string{
				"Review your own code before committing",
				"Ask for specific focus areas: security, performance, readability",
				"Use for learning - ask why certain patterns are recommended",
			},
		},
		{
			Title:       "🧪 Test Development",
			Description: "Creating comprehensive tests for your code",
			Chatmate:    "Testing",
			Example:     "@Testing Generate unit tests for this user service class, including edge cases and error handling.",
			Tips: []string{
				"Ask for different test types: unit, integration, edge cases",
				"Request specific testing frameworks if you have preferences",
				"Include error handling and boundary condition tests",
			},
		},
		{
			Title:       "📝 Documentation Writing",
			Description: "Creating clear documentation for your code",
			Chatmate:    "Documentation",
			Example:     "@Documentation Write API documentation for this REST endpoint, including request/response examples.",
			Tips: []string{
				"Include examples and use cases in documentation requests",
				"Ask for different formats: inline comments, README, API docs",
				"Request beginner-friendly explanations for complex concepts",
			},
		},
	}

	for i, scenario := range scenarios {
		fmt.Printf("Scenario %d: %s\n", i+1, scenario.Title)
		fmt.Println(scenario.Description)
		fmt.Println("")

		fmt.Printf("💡 Use: %s\n", scenario.Chatmate)
		fmt.Printf("📝 Example: %s\n", scenario.Example)
		fmt.Println("")

		fmt.Println("💡 Pro Tips:")
		for _, tip := range scenario.Tips {
			fmt.Printf("  • %s\n", tip)
		}
		fmt.Println("")

		if !promptToContinue("Ready for the next scenario?") {
			return nil
		}
	}

	// End-of-day routine
	fmt.Println("🌙 End-of-Day Routine: Maintenance")
	fmt.Println("Keep your ChatMate setup clean and updated:")
	fmt.Println("")
	fmt.Println("$ chatmate list --installed         # Review what you have")
	fmt.Println("$ chatmate uninstall \"Unused One\"   # Remove unused chatmates")
	fmt.Println("$ chatmate hire --force             # Update existing chatmates (weekly)")
	fmt.Println("")

	fmt.Println("🎯 Key Takeaways:")
	fmt.Println("• Start each coding session with a ChatMate status check")
	fmt.Println("• Use specific chatmates for specific tasks")
	fmt.Println("• Provide context and details for better responses")
	fmt.Println("• Regularly maintain your chatmate collection")
	fmt.Println("• Integrate chatmates into your existing workflow")
	fmt.Println("")

	fmt.Println("✅ Daily Development Tutorial Complete!")
	fmt.Println("Next: Try 'chatmate tutorial team-lead' for team workflows")
	return nil
}

// runTeamLeadTutorial runs the team leadership tutorial
func runTeamLeadTutorial(prompt PromptFunc) error {
	fmt.Println("👥 ChatMate Team Leadership Tutorial")
	fmt.Println("===================================")
	fmt.Println("")

	fmt.Println("Learn how to use ChatMate for team leadership, code reviews, and project management.")
	fmt.Println("")

	if !prompt("Ready to learn team leadership workflows?") {
		return nil
	}

	// Team scenarios
	teamScenarios := []struct {
		Title       string
		Description string
		Workflow    []string
		Example     string
	}{
		{
			Title:       "📋 Sprint Planning & Issue Creation",
			Description: "Converting requirements into actionable GitHub issues",
			Workflow: []string{
				"1. Use @Create Issue to structure requirements",
				"2. Include acceptance criteria and technical details",
				"3. Add appropriate labels and assignments",
				"4. Link related issues and dependencies",
			},
			Example: "@Create Issue Create a user authentication system with OAuth2 integration, including login/logout, session management, and role-based access control.",
		},
		{
			Title:       "🔍 Code Review Leadership",
			Description: "Providing comprehensive and constructive code reviews",
			Workflow: []string{
				"1. Use @Code Review for thorough analysis",
				"2. Focus on architecture, security, and maintainability",
				"3. Provide specific, actionable feedback",
				"4. Suggest improvements and alternatives",
			},
			Example: "@Code Review Analyze this microservice architecture for scalability issues, security vulnerabilities, and adherence to our coding standards.",
		},
		{
			Title:       "📝 Pull Request Management",
			Description: "Creating and reviewing comprehensive pull requests",
			Workflow: []string{
				"1. Use @Create PR for detailed PR descriptions",
				"2. Include testing notes and deployment considerations",
				"3. Use @Review PR for team member's submissions",
				"4. Ensure proper documentation and change logs",
			},
			Example: "@Create PR Create a comprehensive PR for the new authentication system, including security considerations, testing strategy, and migration notes.",
		},
		{
			Title:       "🚀 Release Coordination",
			Description: "Managing merges and release preparation",
			Workflow: []string{
				"1. Use @Review PR for final pre-merge review",
				"2. Use @Merge PR for complex merge scenarios",
				"3. Coordinate with QA and deployment teams",
				"4. Plan rollback strategies and monitoring",
			},
			Example: "@Merge PR Help me safely merge this complex feature branch with database migrations and API changes.",
		},
	}

	for i, scenario := range teamScenarios {
		fmt.Printf("Team Scenario %d: %s\n", i+1, scenario.Title)
		fmt.Println(scenario.Description)
		fmt.Println("")

		fmt.Println("🔄 Workflow:")
		for _, step := range scenario.Workflow {
			fmt.Printf("   %s\n", step)
		}
		fmt.Println("")

		fmt.Printf("💬 Example: %s\n", scenario.Example)
		fmt.Println("")

		if !prompt("Ready for the next team scenario?") {
			return nil
		}
	}

	fmt.Println("🎯 Team Leadership Best Practices:")
	fmt.Println("• Establish consistent ChatMate usage across the team")
	fmt.Println("• Document which chatmates to use for different scenarios")
	fmt.Println("• Use chatmates to maintain code quality standards")
	fmt.Println("• Share chatmate workflows in team documentation")
	fmt.Println("• Regularly review and update team's chatmate collection")
	fmt.Println("")

	fmt.Println("✅ Team Leadership Tutorial Complete!")
	return nil
}

// runDebuggingTutorial runs the debugging-focused tutorial
func runDebuggingTutorial(prompt PromptFunc) error {
	fmt.Println("🐛 Advanced Debugging with Solve Issue Chatmate")
	fmt.Println("==============================================")
	fmt.Println("")

	fmt.Println("Master debugging techniques using the Solve Issue chatmate.")
	fmt.Println("")

	// Debugging scenarios with examples
	debugScenarios := []struct {
		Problem  string
		Approach string
		Example  string
		Tips     []string
	}{
		{
			Problem:  "🔥 Critical Production Bug",
			Approach: "Systematic root cause analysis",
			Example:  "@Solve Issue Production users are getting 500 errors when trying to checkout. Error logs show 'Database connection timeout' but CPU and memory usage look normal.",
			Tips: []string{
				"Include error messages, logs, and stack traces",
				"Provide system metrics and resource usage",
				"Describe the user impact and frequency",
				"Include recent changes or deployments",
			},
		},
		{
			Problem:  "🌐 Frontend Performance Issue",
			Approach: "Performance bottleneck identification",
			Example:  "@Solve Issue Our React app is loading slowly. Lighthouse shows 'Largest Contentful Paint' at 4.2s. Bundle size is 2.1MB and we're using code splitting.",
			Tips: []string{
				"Include performance metrics and tools output",
				"Provide bundle analysis and network timing",
				"Describe user experience impact",
				"Include current optimization attempts",
			},
		},
		{
			Problem:  "🗄️ Database Query Optimization",
			Approach: "Query analysis and optimization",
			Example:  "@Solve Issue This user search query is taking 3+ seconds. It joins 4 tables and filters on multiple columns. Query plan shows full table scans.",
			Tips: []string{
				"Include actual query and execution plan",
				"Provide table sizes and index information",
				"Share current query performance metrics",
				"Include database version and configuration",
			},
		},
	}

	for i, scenario := range debugScenarios {
		fmt.Printf("Debug Scenario %d: %s\n", i+1, scenario.Problem)
		fmt.Printf("Approach: %s\n", scenario.Approach)
		fmt.Println("")

		fmt.Printf("🔍 Example Query: %s\n", scenario.Example)
		fmt.Println("")

		fmt.Println("💡 Key Information to Include:")
		for _, tip := range scenario.Tips {
			fmt.Printf("  • %s\n", tip)
		}
		fmt.Println("")

		if !prompt("Ready for the next debugging scenario?") {
			return nil
		}
	}

	fmt.Println("🎯 Debugging Best Practices with Solve Issue:")
	fmt.Println("• Be specific and detailed in problem descriptions")
	fmt.Println("• Include relevant code, logs, and error messages")
	fmt.Println("• Provide context about recent changes")
	fmt.Println("• Ask follow-up questions to narrow scope")
	fmt.Println("• Test proposed solutions in safe environments")
	fmt.Println("")

	fmt.Println("✅ Advanced Debugging Tutorial Complete!")
	return nil
}

// runTestingTutorial runs the testing-focused tutorial
func runTestingTutorial(prompt PromptFunc) error {
	fmt.Println("🧪 Comprehensive Testing with Testing Chatmate")
	fmt.Println("==============================================")
	fmt.Println("")

	fmt.Println("Learn testing strategies and best practices using the Testing chatmate.")
	fmt.Println("")

	// Testing scenarios
	testScenarios := []struct {
		TestType    string
		Description string
		Example     string
		Benefits    []string
	}{
		{
			TestType:    "Unit Testing",
			Description: "Testing individual functions and methods",
			Example:     "@Testing Generate unit tests for this user authentication service, including password validation, token generation, and error cases.",
			Benefits: []string{
				"Fast feedback during development",
				"Easy to isolate and fix issues",
				"Documents expected behavior",
				"Enables safe refactoring",
			},
		},
		{
			TestType:    "Integration Testing",
			Description: "Testing component interactions and workflows",
			Example:     "@Testing Create integration tests for the user registration flow, including database operations, email sending, and API responses.",
			Benefits: []string{
				"Validates component interactions",
				"Catches integration bugs early",
				"Tests realistic user scenarios",
				"Verifies external dependencies",
			},
		},
		{
			TestType:    "Edge Case Testing",
			Description: "Testing boundary conditions and error scenarios",
			Example:     "@Testing Generate edge case tests for this payment processing function, including invalid inputs, network failures, and timeout scenarios.",
			Benefits: []string{
				"Improves system reliability",
				"Prevents production errors",
				"Tests error handling logic",
				"Validates input sanitization",
			},
		},
	}

	for i, scenario := range testScenarios {
		fmt.Printf("Testing Approach %d: %s\n", i+1, scenario.TestType)
		fmt.Printf("Focus: %s\n", scenario.Description)
		fmt.Println("")

		fmt.Printf("🧪 Example: %s\n", scenario.Example)
		fmt.Println("")

		fmt.Println("✅ Benefits:")
		for _, benefit := range scenario.Benefits {
			fmt.Printf("  • %s\n", benefit)
		}
		fmt.Println("")

		if !prompt("Ready for the next testing approach?") {
			return nil
		}
	}

	fmt.Println("🎯 Testing Best Practices:")
	fmt.Println("• Write tests before or during development (TDD/BDD)")
	fmt.Println("• Cover happy path, edge cases, and error scenarios")
	fmt.Println("• Use descriptive test names and clear assertions")
	fmt.Println("• Keep tests independent and repeatable")
	fmt.Println("• Regularly review and update test coverage")
	fmt.Println("")

	fmt.Println("✅ Comprehensive Testing Tutorial Complete!")
	return nil
}

// promptToContinue asks the user if they want to continue and returns their response
func promptToContinue(message string) bool {
	fmt.Printf("❓ %s [Y/n]: ", message)

	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		return false
	}

	response := strings.TrimSpace(strings.ToLower(scanner.Text()))
	return response == "" || response == "y" || response == "yes"
}

func init() {
	rootCmd.AddCommand(tutorialCmd)
}
