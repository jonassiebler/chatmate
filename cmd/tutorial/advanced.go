package tutorial

import "fmt"

// RunTeamLeadTutorial runs the team leadership workflow tutorial
func RunTeamLeadTutorial(prompt PromptFunc) error {
	fmt.Println("👑 ChatMate Team Leadership Tutorial")
	fmt.Println("====================================")
	fmt.Println("")

	fmt.Println("Learn how to use ChatMate for team leadership, code reviews, and project management.")
	fmt.Println("")

	if !prompt("Ready to learn team leadership workflows?") {
		return nil
	}

	teamScenarios := GetTeamLeadScenarios()
	showScenarios("👥 Team Leadership Scenarios:", teamScenarios, prompt)

	fmt.Println("🎯 Team Leadership Best Practices:")
	fmt.Println("• Maintain consistent code review standards")
	fmt.Println("• Create clear, actionable issues and documentation")
	fmt.Println("• Use code reviews as learning opportunities")
	fmt.Println("• Regular maintenance of issue quality and organization")
	fmt.Println("")

	fmt.Println("✅ Team Leadership Tutorial Complete!")
	return nil
}

// RunDebuggingTutorial runs the advanced debugging tutorial
func RunDebuggingTutorial(prompt PromptFunc) error {
	fmt.Println("🐛 Advanced Debugging with Solve Issue Chatmate")
	fmt.Println("===============================================")
	fmt.Println("")

	fmt.Println("Master advanced debugging techniques using the Solve Issue chatmate.")
	fmt.Println("")

	if !prompt("Ready for advanced debugging techniques?") {
		return nil
	}

	debugScenarios := GetDebuggingScenarios()
	showScenarios("🔍 Advanced Debugging Scenarios:", debugScenarios, prompt)

	fmt.Println("🎯 Advanced Debugging Best Practices:")
	fmt.Println("• Systematic root cause analysis")
	fmt.Println("• Include comprehensive diagnostic information")
	fmt.Println("• Document your debugging process for future reference")
	fmt.Println("• Easy to isolate and fix issues")
	fmt.Println("")

	fmt.Println("✅ Advanced Debugging Tutorial Complete!")
	return nil
}

// RunTestingTutorial runs the comprehensive testing tutorial
func RunTestingTutorial(prompt PromptFunc) error {
	fmt.Println("🧪 Comprehensive Testing with Testing Chatmate")
	fmt.Println("===============================================")
	fmt.Println("")

	fmt.Println("Learn comprehensive testing strategies using the Testing chatmate.")
	fmt.Println("")

	if !prompt("Ready to learn comprehensive testing?") {
		return nil
	}

	testScenarios := GetTestingScenarios()
	showScenarios("🧪 Testing Scenarios:", testScenarios, prompt)

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

// showScenarios displays scenarios with interactive examples
func showScenarios(title string, scenarios []ScenarioInfo, prompt PromptFunc) {
	fmt.Println(title)
	fmt.Println("")

	for i, scenario := range scenarios {
		fmt.Printf("%d. %s\n", i+1, scenario.Title)
		fmt.Printf("   %s\n", scenario.Description)
		fmt.Printf("   Example: %s\n", scenario.Example)
		fmt.Println("")

		fmt.Println("   💡 Tips:")
		for _, tip := range scenario.Tips {
			fmt.Printf("   • %s\n", tip)
		}
		fmt.Println("")

		if i < len(scenarios)-1 {
			if !prompt("Ready for the next scenario?") {
				return
			}
		}
	}
}
