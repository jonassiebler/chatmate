package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// completionCmd represents the completion command
var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "🚀 Generate shell completion scripts",
	Long: `Generate shell completion scripts for chatmate.

The completion scripts allow you to use tab completion for chatmate commands,
flags, and arguments in your shell. This greatly improves the user experience
by providing auto-completion for chatmate names, commands, and options.

Supported shells:
  • bash     - Bash completion (Linux, macOS, Windows)  
  • zsh      - Zsh completion (macOS default, Linux)
  • fish     - Fish shell completion
  • powershell - PowerShell completion (Windows)

Installation:
  
  Bash (Linux):
    chatmate completion bash | sudo tee /etc/bash_completion.d/chatmate

  Bash (macOS with Homebrew):
    chatmate completion bash | tee $(brew --prefix)/etc/bash_completion.d/chatmate

  Zsh:
    chatmate completion zsh | tee ~/.zsh/completions/_chatmate
    # Add to your .zshrc: fpath=(~/.zsh/completions $fpath)

  Fish:
    chatmate completion fish | tee ~/.config/fish/completions/chatmate.fish

  PowerShell:
    chatmate completion powershell | Out-String | Invoke-Expression

For persistent installation, add the appropriate command to your shell's
configuration file (.bashrc, .zshrc, etc.).`,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Example: `  # Generate bash completion
  chatmate completion bash

  # Install bash completion on Linux
  chatmate completion bash | sudo tee /etc/bash_completion.d/chatmate

  # Install zsh completion
  mkdir -p ~/.zsh/completions
  chatmate completion zsh > ~/.zsh/completions/_chatmate

  # Install fish completion  
  chatmate completion fish > ~/.config/fish/completions/chatmate.fish

  # Test completion (after installation)
  chatmate <TAB>          # Show available commands
  chatmate hire <TAB>     # Show hire command options`,
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		switch args[0] {
		case "bash":
			err = cmd.Root().GenBashCompletion(os.Stdout)
		case "zsh":
			err = cmd.Root().GenZshCompletion(os.Stdout)
		case "fish":
			err = cmd.Root().GenFishCompletion(os.Stdout, true)
		case "powershell":
			err = cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
		}
		if err != nil {
			fmt.Printf("Error generating completion: %v\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
}
