package tutorial

import (
	"fmt"

	"github.com/jonassiebler/chatmate/internal/manager"
)

// RunFirstTimeTutorial runs the beginner tutorial
func RunFirstTimeTutorial(prompt PromptFunc) error {
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

	err = chatMateManager.Status().ShowStatus()
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

	if prompt("Would you like to install these recommended chatmates?") {
		fmt.Println("Running: chatmate hire \"Solve Issue\" \"Code Review\" \"Testing\"")
		fmt.Println("")

		err = chatMateManager.Installer().InstallSpecific([]string{"Solve Issue", "Code Review", "Testing"}, false)
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
	err = chatMateManager.Lister().ListInstalled()
	if err != nil {
		fmt.Printf("❌ Error listing chatmates: %v\n", err)
		return nil
	}

	fmt.Println("")
	if !prompt("Do you see your installed chatmates listed above?") {
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

	if !prompt("Ready to try this in VS Code?") {
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

// RunDailyDevTutorial runs the daily development workflow tutorial
func RunDailyDevTutorial(prompt PromptFunc) error {
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

	err = chatMateManager.Status().ShowStatus()
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return nil
	}

	fmt.Println("")

	// Show scenarios
	scenarios := GetDailyDevScenarios()
	showScenarios("🎯 Daily Development Scenarios:", scenarios, prompt)

	fmt.Println("🎯 Debugging Best Practices with Solve Issue:")
	fmt.Println("• Provide comprehensive context and error details")
	fmt.Println("• Include environment information and recent changes")
	fmt.Println("• Testing boundary conditions and error scenarios")
	fmt.Println("• Regularly maintain your chatmate collection")
	fmt.Println("")

	fmt.Println("✅ Daily Development Tutorial Complete!")
	return nil
}
