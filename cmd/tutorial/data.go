package tutorial

// GetAvailableTutorials returns metadata for all available tutorials
func GetAvailableTutorials() []TutorialInfo {
	return []TutorialInfo{
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
}

// GetDailyDevScenarios returns scenarios for daily development workflow
func GetDailyDevScenarios() []ScenarioInfo {
	return []ScenarioInfo{
		{
			Title:       "Debugging Session",
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
			Title:       "Code Review",
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
			Title:       "Test Development",
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
			Title:       "Issue Creation",
			Description: "Creating well-structured issues for bugs or features",
			Chatmate:    "Create Issue",
			Example:     "@Create Issue Help me create a clear issue for improving our authentication system performance.",
			Tips: []string{
				"Provide clear problem descriptions and context",
				"Include steps to reproduce for bugs",
				"Specify acceptance criteria for features",
			},
		},
		{
			Title:       "Pull Request",
			Description: "Creating comprehensive PRs for your changes",
			Chatmate:    "Create PR",
			Example:     "@Create PR Help me create a detailed PR for the authentication middleware I just implemented.",
			Tips: []string{
				"Explain the what and why of your changes",
				"Include testing information and screenshots",
				"Reference related issues and breaking changes",
			},
		},
	}
}

// GetTeamLeadScenarios returns scenarios for team leadership workflows
func GetTeamLeadScenarios() []ScenarioInfo {
	return []ScenarioInfo{
		{
			Title:       "Code Review Leadership",
			Description: "Leading thorough code reviews for your team",
			Chatmate:    "Review PR",
			Example:     "@Review PR Please provide a comprehensive review of this authentication PR, focusing on security and maintainability.",
			Tips: []string{
				"Set clear review standards and communicate them",
				"Focus on architecture, security, and maintainability",
				"Provide constructive feedback with examples",
				"Use reviews as learning opportunities for the team",
			},
		},
		{
			Title:       "Issue Management",
			Description: "Creating and optimizing issues for team clarity",
			Chatmate:    "Create Issue",
			Example:     "@Create Issue Create a detailed epic for our user authentication system redesign.",
			Tips: []string{
				"Write clear, actionable issues with good descriptions",
				"Include acceptance criteria and definition of done",
				"Use labels and milestones for organization",
				"Break down large features into manageable tasks",
			},
		},
		{
			Title:       "Issue Optimization",
			Description: "Improving existing issues for better team productivity",
			Chatmate:    "Optimize Issues",
			Example:     "@Optimize Issues Help improve this vague bug report to be more actionable for the team.",
			Tips: []string{
				"Regularly audit and improve issue quality",
				"Add missing context and reproduction steps",
				"Ensure issues have proper priority and labels",
				"Close outdated or duplicate issues",
			},
		},
	}
}

// GetDebuggingScenarios returns scenarios for debugging workflows
func GetDebuggingScenarios() []ScenarioInfo {
	return []ScenarioInfo{
		{
			Title:       "Frontend Performance Issue",
			Description: "Analyze and solve frontend performance problems",
			Example:     "@Solve Issue Our React app is loading slowly. Lighthouse shows 'Largest Contentful Paint' at 4.2s. Bundle size is 2.1MB and we're using code splitting.",
			Tips: []string{
				"Include performance metrics and measurements",
				"Describe user experience impact",
				"Mention current optimization attempts",
			},
		},
		{
			Title:       "Backend API Problem",
			Description: "Debug backend services and API issues",
			Example:     "@Solve Issue Production users are getting 500 errors when trying to checkout. Error logs show 'Database connection timeout' but CPU and memory usage look normal.",
			Tips: []string{
				"Include error logs and stack traces",
				"Provide system resource information",
				"Mention recent changes or deployments",
			},
		},
		{
			Title:       "Database Query Optimization",
			Description: "Optimize slow database queries and operations",
			Example:     "@Solve Issue This user search query is taking 3+ seconds. It joins 4 tables and filters on multiple columns. Query plan shows full table scans.",
			Tips: []string{
				"Include query execution plans",
				"Provide table sizes and index information",
				"Query analysis and optimization",
			},
		},
	}
}

// GetTestingScenarios returns scenarios for testing workflows
func GetTestingScenarios() []ScenarioInfo {
	return []ScenarioInfo{
		{
			Title:       "Unit Testing",
			Description: "Generate comprehensive unit tests",
			Example:     "@Testing Generate unit tests for this user authentication service, including password validation, token generation, and error cases.",
			Tips: []string{
				"Test both happy path and edge cases",
				"Include error handling scenarios",
				"Mock external dependencies appropriately",
			},
		},
		{
			Title:       "Integration Testing",
			Description: "Create tests for component interactions",
			Example:     "@Testing Create integration tests for the user registration flow, including database operations, email sending, and API responses.",
			Tips: []string{
				"Test realistic user scenarios",
				"Include external service interactions",
				"Verify data flow between components",
			},
		},
		{
			Title:       "Performance Testing",
			Description: "Generate tests for performance validation",
			Example:     "@Testing Generate edge case tests for this payment processing function, including invalid inputs, network failures, and timeout scenarios.",
			Tips: []string{
				"Include load and stress scenarios",
				"Test timeout and retry logic",
				"Validate resource usage patterns",
			},
		},
	}
}
