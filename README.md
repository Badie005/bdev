# 🚀 B.DEV CLI - Enterprise Developer Workstation

> **One CLI to rule them all** - Unified command-line interface for all development tasks

[![Version](https://img.shields.io/badge/version-3.0.0-FF6B35)](https://github.com/badie/bdev) 
[![Go](https://img.shields.io/badge/go-1.21+-00ADD8)](https://go.dev)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)
[![CI](https://github.com/badie/bdev/actions/workflows/ci.yml/badge.svg)](https://github.com/badie/bdev/actions/workflows/ci.yml)

---

## 📋 Table of Contents

- [Features](#-features)
- [Installation](#-installation)
- [Quick Start](#-quick-start)
- [Commands Reference](#-commands-reference)
- [AI Integration](#-ai-integration)
- [Workflow Automation](#-workflow-automation)
- [Secrets Vault](#-secrets-vault)
- [Configuration](#-configuration)
- [Development](#-development)

---

## ✨ Features

| Feature | Description |
|---------|-------------|
| 🎯 **Unified CLI** | All dev tools under `bdev` command |
| 🤖 **AI Assistant** | Ollama-powered local AI with memory |
| 📦 **Multi-Project** | Manage multiple projects at once |
| ⚡ **Quick Actions** | `start`, `test`, `build` shortcuts |
| 🔐 **Secrets Vault** | Encrypted local secrets storage |
| 🔄 **Workflows** | YAML-based task automation |
| 🎨 **Claude Design** | Professional UI with warm palette |

---

## 📦 Installation

### From Source (Go 1.21+)

```bash
# Clone repository
git clone https://github.com/badie/bdev.git
cd bdev

# Build
go build -o bdev ./cmd/bdev

# Install to PATH (optional)
go install ./cmd/bdev
```

### Quick Install (Windows)

```powershell
# Build and install
make build
copy bdev.exe C:\Users\%USERNAME%\go\bin\
```

---

## 🚀 Quick Start

```powershell
# Interactive REPL mode (recommended)
bdev

# List projects
bdev list

# Start dev server (auto-detects project type)
bdev start

# AI chat
bdev ai chat "explain this code"

# Run tests
bdev test
```

---

## 📖 Commands Reference

### 🎯 Quick Actions (Root Level)
| Command | Description |
|---------|-------------|
| `bdev` | Interactive REPL mode |
| `bdev list` | List all projects |
| `bdev start` | Start dev server |
| `bdev test` | Run tests |
| `bdev build` | Build for production |
| `bdev fix` | Auto-fix linting |
| `bdev deploy` | Deploy project |
| `bdev do "<natural language>"` | Execute via NLP |

### 📁 Projects (`bdev projects`)
| Command | Description |
|---------|-------------|
| `bdev projects list` | List with details |
| `bdev projects new <template> <name>` | Create from template |
| `bdev projects open <name>` | Open in VS Code |
| `bdev projects find <query>` | Search projects |

### 🔧 Git (`bdev git`)
| Command | Description |
|---------|-------------|
| `bdev git status` | Enhanced status |
| `bdev git commit "<msg>"` | Commit changes |
| `bdev git push` | Push to remote |
| `bdev git pull` | Pull from remote |
| `bdev git branch <name>` | Create/switch branch |
| `bdev git log -n 10` | Show recent commits |

### 🤖 AI (`bdev ai`)
| Command | Description |
|---------|-------------|
| `bdev ai chat "<question>"` | Chat with AI |
| `bdev ai chat "<q>" --context` | Include project context |
| `bdev ai memory` | View conversation history |
| `bdev ai forget` | Clear memory |
| `bdev ai generate "<prompt>"` | Generate code |

### 🤖 Agents (`bdev agent`)
| Command | Description |
|---------|-------------|
| `bdev agent list` | List available agents |
| `bdev agent review <file>` | Code review |
| `bdev agent document <file>` | Generate docs |
| `bdev agent explain <file>` | Explain code |

### 📊 Analytics (`bdev analytics`)
| Command | Description |
|---------|-------------|
| `bdev analytics today` | Today's activity |
| `bdev analytics week` | Weekly stats |
| `bdev analytics summary` | Quick overview |

### 🌐 Multi-Project (`bdev multi`)
| Command | Description |
|---------|-------------|
| `bdev multi status` | Git status all projects |
| `bdev multi pull` | Pull all repositories |
| `bdev multi audit` | Security audit all |
| `bdev multi run "<cmd>"` | Run command in all |
| `bdev multi update` | Update dependencies |

### 🔐 Secrets (`bdev secrets`)
| Command | Description |
|---------|-------------|
| `bdev secrets init` | Initialize vault |
| `bdev secrets set <key>` | Store secret |
| `bdev secrets get <key>` | Retrieve secret |
| `bdev secrets list` | List all keys |
| `bdev secrets export` | Export as env vars |

### 🔄 Workflow (`bdev workflow`)
| Command | Description |
|---------|-------------|
| `bdev workflow list` | List workflows |
| `bdev workflow create <name>` | Create template |
| `bdev workflow run <name>` | Execute workflow |
| `bdev workflow show <name>` | View steps |

### ⚙️ Config (`bdev config`)
| Command | Description |
|---------|-------------|
| `bdev config show` | Display config |
| `bdev config set <key> <val>` | Set value |
| `bdev config alias <n> <cmd>` | Create alias |
| `bdev config edit` | Open in editor |

### 🎨 Theme (`bdev theme`)
| Command | Description |
|---------|-------------|
| `bdev theme list` | Available themes |
| `bdev theme set <name>` | Change theme |
| `bdev theme preview <name>` | Preview colors |

---

## 🤖 AI Integration

B.DEV uses **Ollama** for local AI capabilities.

```powershell
# Install Ollama
winget install Ollama.Ollama

# Pull recommended models
ollama pull llama3.2
ollama pull codellama:7b

# Use AI
bdev ai chat "How do I create a REST API in Python?"
bdev ai chat "Review this code" --context
```

### AI Memory
The AI remembers your conversation within a session:
```powershell
bdev ai chat "My project uses FastAPI"
bdev ai chat "How do I add authentication?"  # Remembers context
bdev ai forget  # Clear memory
```

---

## 🔄 Workflow Automation

Create YAML workflows for repetitive tasks:

```yaml
# ~/.bdev/workflows/deploy.yml
name: deploy
description: Deploy to production

steps:
  - name: Run tests
    run: npm test
    
  - name: Build
    run: npm run build
    
  - name: Deploy
    run: npm run deploy

on_success: echo "Deployed successfully!"
on_failure: echo "Deployment failed!"
```

```powershell
bdev workflow run deploy
```

---

## 🔐 Secrets Vault

Encrypted local storage for API keys, tokens, etc.

```powershell
# Initialize vault (first time)
bdev secrets init

# Store secrets
bdev secrets set OPENAI_KEY
bdev secrets set DATABASE_URL "postgres://..."

# Use in scripts
$env:BDEV_OPENAI_KEY = $(bdev secrets get OPENAI_KEY --show)
```

---

## 🎨 Themes

B.DEV uses Claude Code's design palette by default.

```powershell
bdev theme list       # See available themes
bdev theme set claude # Orange palette (#FF6B35)
bdev theme set gemini # Blue palette
bdev theme set matrix # Green cyberpunk
```

### Color Palette (Claude)
| Color | Hex | Usage |
|-------|-----|-------|
| Primary | `#FF6B35` | Actions, highlights |
| Success | `#34C759` | Success messages |
| Error | `#FF3B30` | Error messages |
| Warning | `#FF9500` | Warnings |
| Gray | `#8E8E93` | Secondary text |

---

## ⚙️ Configuration

Configuration is stored in `~/.bdev/config.json`:

```json
{
  "display": {
    "theme": "claude"
  },
  "ai": {
    "model": "llama3.2",
    "timeout": 120
  },
  "paths": {
    "projects": "~/Dev/Projects"
  },
  "aliases": {
    "gs": "git status",
    "gp": "git push"
  }
}
```

---

## 🏗️ Architecture

```
.bdev/
├── cli/
│   ├── main.py              # Entry point + Typer app
│   ├── commands/            # Command modules
│   │   ├── projects.py
│   │   ├── ai.py
│   │   ├── git.py
│   │   ├── workflow.py
│   │   └── secrets.py
│   ├── core/                # Core logic
│   │   ├── repl.py          # Interactive REPL
│   │   ├── session.py       # State management
│   │   ├── workflow.py      # Workflow engine
│   │   ├── vault.py         # Secrets vault
│   │   └── ai/
│   │       ├── engine.py    # Ollama wrapper
│   │       └── agents/      # AI agents
│   └── utils/
│       ├── ui.py            # Rich components
│       ├── theme.py         # Theme engine
│       └── branding.py      # ASCII art
├── workflows/               # User workflows
├── templates/               # Project templates
└── config.json
```

---

## 🔧 Dependencies

```
typer>=0.9.0
rich>=13.0.0
psutil>=5.9.0
pyyaml>=6.0.0
colorama>=0.4.6
```

---

## 📜 License

MIT © 2025 B.DEV

---

<p align="center">
  <strong>B.DEV CLI</strong> - Enterprise Developer Workstation<br>
  <em>Built with ❤️ and #FF6B35</em>
</p>
