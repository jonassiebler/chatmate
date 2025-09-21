#!/bin/bash

# Install Shell Completions Script for chatmate
# This script installs shell completion for bash, zsh, and fish

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Logging functions
log_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

log_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

log_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

log_error() {
    echo -e "${RED}❌ $1${NC}"
}

# Check if chatmate is available
check_chatmate() {
    if ! command -v chatmate &> /dev/null; then
        log_error "chatmate command not found. Please install chatmate first."
        exit 1
    fi
    log_success "chatmate found: $(which chatmate)"
}

# Install bash completion
install_bash_completion() {
    log_info "Installing bash completion..."
    
    # Try different locations for bash completion
    if [[ "$OSTYPE" == "darwin"* ]]; then
        # macOS
        if command -v brew &> /dev/null; then
            # Homebrew bash completion
            local comp_dir="$(brew --prefix)/etc/bash_completion.d"
            if [[ -d "$comp_dir" ]]; then
                chatmate completion bash | sudo tee "$comp_dir/chatmate" > /dev/null
                log_success "Bash completion installed to $comp_dir/chatmate"
                log_info "Restart your shell or run: source ~/.bashrc"
                return
            fi
        fi
        
        # Fallback to user directory
        local user_comp_dir="$HOME/.bash_completion.d"
        mkdir -p "$user_comp_dir"
        chatmate completion bash > "$user_comp_dir/chatmate"
        
        # Add to .bashrc if not already present
        local bashrc="$HOME/.bashrc"
        if [[ -f "$bashrc" ]] && ! grep -q ".bash_completion.d" "$bashrc"; then
            echo "" >> "$bashrc"
            echo "# Load bash completions" >> "$bashrc"
            echo "for f in ~/.bash_completion.d/*; do [ -f \"\$f\" ] && source \"\$f\"; done" >> "$bashrc"
        fi
        
        log_success "Bash completion installed to $user_comp_dir/chatmate"
    else
        # Linux
        if [[ -d "/etc/bash_completion.d" ]] && [[ -w "/etc/bash_completion.d" ]]; then
            chatmate completion bash | sudo tee /etc/bash_completion.d/chatmate > /dev/null
            log_success "Bash completion installed to /etc/bash_completion.d/chatmate"
        else
            # User directory fallback
            local user_comp_dir="$HOME/.bash_completion.d"
            mkdir -p "$user_comp_dir"
            chatmate completion bash > "$user_comp_dir/chatmate"
            
            # Add to .bashrc if not already present
            local bashrc="$HOME/.bashrc"
            if [[ -f "$bashrc" ]] && ! grep -q ".bash_completion.d" "$bashrc"; then
                echo "" >> "$bashrc"
                echo "# Load bash completions" >> "$bashrc"
                echo "for f in ~/.bash_completion.d/*; do [ -f \"\$f\" ] && source \"\$f\"; done" >> "$bashrc"
            fi
            
            log_success "Bash completion installed to $user_comp_dir/chatmate"
        fi
    fi
}

# Install zsh completion
install_zsh_completion() {
    log_info "Installing zsh completion..."
    
    # Create completion directory
    local comp_dir="$HOME/.zsh/completions"
    mkdir -p "$comp_dir"
    
    # Generate completion file
    chatmate completion zsh > "$comp_dir/_chatmate"
    log_success "Zsh completion installed to $comp_dir/_chatmate"
    
    # Add to .zshrc if not already present
    local zshrc="$HOME/.zshrc"
    if [[ -f "$zshrc" ]]; then
        if ! grep -q "fpath=.*\.zsh/completions" "$zshrc"; then
            echo "" >> "$zshrc"
            echo "# Load zsh completions" >> "$zshrc"
            echo "fpath=(~/.zsh/completions \$fpath)" >> "$zshrc"
            echo "autoload -U compinit && compinit" >> "$zshrc"
            log_info "Added completion setup to $zshrc"
        fi
    else
        log_warning ".zshrc not found. You may need to add completion setup manually."
        log_info "Add this to your .zshrc:"
        echo "    fpath=(~/.zsh/completions \$fpath)"
        echo "    autoload -U compinit && compinit"
    fi
}

# Install fish completion
install_fish_completion() {
    log_info "Installing fish completion..."
    
    # Create fish completion directory
    local comp_dir="$HOME/.config/fish/completions"
    mkdir -p "$comp_dir"
    
    # Generate completion file
    chatmate completion fish > "$comp_dir/chatmate.fish"
    log_success "Fish completion installed to $comp_dir/chatmate.fish"
    log_info "Fish will automatically load the completion on next start"
}

# Main installation function
install_completions() {
    local shell="${1:-}"
    
    if [[ -z "$shell" ]]; then
        log_info "Detecting current shell: $SHELL"
        case "$SHELL" in
            */bash)
                shell="bash"
                ;;
            */zsh)
                shell="zsh"
                ;;
            */fish)
                shell="fish"
                ;;
            *)
                log_warning "Unknown shell: $SHELL"
                log_info "Available options: bash, zsh, fish, all"
                exit 1
                ;;
        esac
    fi
    
    case "$shell" in
        bash)
            install_bash_completion
            ;;
        zsh)
            install_zsh_completion
            ;;
        fish)
            install_fish_completion
            ;;
        all)
            install_bash_completion
            install_zsh_completion
            install_fish_completion
            ;;
        *)
            log_error "Unsupported shell: $shell"
            log_info "Available options: bash, zsh, fish, all"
            exit 1
            ;;
    esac
}

# Show usage
show_usage() {
    echo "Usage: $0 [shell]"
    echo ""
    echo "Install shell completion for chatmate"
    echo ""
    echo "Options:"
    echo "  shell    Shell to install completion for (bash|zsh|fish|all)"
    echo "           If not specified, detects current shell automatically"
    echo ""
    echo "Examples:"
    echo "  $0           # Install for current shell"
    echo "  $0 bash      # Install bash completion"
    echo "  $0 zsh       # Install zsh completion"
    echo "  $0 fish      # Install fish completion"
    echo "  $0 all       # Install for all supported shells"
}

# Main script
main() {
    if [[ "$1" == "-h" ]] || [[ "$1" == "--help" ]]; then
        show_usage
        exit 0
    fi
    
    log_info "🚀 Installing chatmate shell completions"
    
    check_chatmate
    install_completions "$1"
    
    echo ""
    log_success "🎉 Shell completion installation complete!"
    log_info "Restart your shell or source your shell configuration file to activate completions"
    log_info "Test with: chatmate <TAB>"
}

# Run main function
main "$@"
