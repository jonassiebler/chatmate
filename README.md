# ChatMate 🤖

Specialized AI agents for VS Code Copilot Chat. Each chatmate excels at specific development tasks.

[![Go](https://github.com/jonassiebler/chatmate/actions/workflows/go.yml/badge.svg)](https://github.com/jonassiebler/chatmate/actions/workflows/go.yml)
[![Security](https://github.com/jonassiebler/chatmate/actions/workflows/security.yml/badge.svg)](https://github.com/jonassiebler/chatmate/actions/workflows/security.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/jonassiebler/chatmate)](https://goreportcard.com/report/github.com/jonassiebler/chatmate)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## 🚀 Quick Start

### Install via Homebrew
```bash
brew tap jonassiebler/chatmate https://github.com/jonassiebler/chatmate.git
brew install chatmate
chatmate hire
```

### Build from Source
```bash
git clone https://github.com/jonassiebler/chatmate.git
cd chatmate && go build -o chatmate .
./chatmate hire
```

**Restart VS Code** → Select chatmates from dropdown in Copilot Chat

## 🤖 Available Chatmates

> **💡 Version Information**: Chatmates are available in two versions:
> - **Original versions**: Comprehensive, detailed instructions with extensive examples
> - **v2 versions (Optimized)**: Streamlined, language-agnostic with 3-Domain Safety Paradigm (Implementation-Testing-Documentation validation)
>
> Choose based on preference: comprehensive guidance (original) vs. efficient execution (v2).

| Chatmate | Purpose | Example |
|----------|---------|---------|
| **Chatmate - Solve Issue** 🐛 | Debug systematically | `My React component won't render` |
| **Code Review** 👁️ | Analyze & improve code | `Check this authentication logic` |
| **Chatmate - Testing** 🧪 | Generate & debug tests | `Unit tests for this service` |
| **Chatmate - Create PR** 📝 | Pull request creation | `PR for new auth feature` |
| **Chatmate - Create Issue** 🎯 | GitHub issue creation | `Login fails on mobile` |

Run `chatmate list` for all available agents.

## 📋 Commands

```bash
chatmate hire                    # Install all chatmates
chatmate list                    # Show available/installed
chatmate status                  # Check installation health
chatmate uninstall "Name"       # Remove specific chatmate
chatmate hire --force           # Force reinstall
```

## � Requirements & Troubleshooting

**Requirements:** VS Code with GitHub Copilot Chat extension

**Issues?**
- Chatmates not appearing → Restart VS Code completely
- Permission errors → `chatmate config` (check paths)
- Need help → `chatmate status` (system diagnostics)

## 🔧 Alternative Installation

### Manual Build
```bash
git clone https://github.com/jonassiebler/chatmate.git
cd chatmate && go build -o chatmate .
./chatmate hire
```

### Legacy Script
```bash
git clone https://github.com/jonassiebler/chatmate.git
cd chatmate && ./hire.sh
```

## 🤝 Contributing

1. Fork repository
2. Create feature branch from `dev`
3. Make changes & test
4. Submit PR to `dev` branch

**Creating Chatmates:**
1. Use `Chatmate - Create Chatmode` agent
2. Add `.chatmode.md` to `internal/assets/mates/`
3. Test with `chatmate hire`
4. Submit PR

---

**Ready to enhance your VS Code workflow?**
```bash
chatmate hire
```
Restart VS Code → Select chatmates from Copilot Chat dropdown!
