package agents

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/badie/bdev/internal/ai/engine"
	"github.com/badie/bdev/internal/core/config"
	"github.com/badie/bdev/internal/core/projects"
	"github.com/badie/bdev/pkg/ui"
)

// Agent interface defines a specialized AI agent
type Agent interface {
	Name() string
	Description() string
	SystemPrompt() string
}

// BaseAgent provides common agent functionality
type BaseAgent struct {
	name        string
	description string
	prompt      string
}

func (a *BaseAgent) Name() string         { return a.name }
func (a *BaseAgent) Description() string  { return a.description }
func (a *BaseAgent) SystemPrompt() string { return a.prompt }

// ============================================================
// AGENTS
// ============================================================

// ReviewerAgent performs code review
var ReviewerAgent = &BaseAgent{
	name:        "reviewer",
	description: "Code review and quality analysis",
	prompt: `You are B.DEV Reviewer, an expert code reviewer.
Your role is to:
1. Identify bugs, security issues, and code smells
2. Suggest improvements for readability and maintainability
3. Check for best practices and design patterns
4. Highlight potential performance issues

Format your review as:
- ISSUES: Critical problems that must be fixed
- WARNINGS: Potential problems to consider
- SUGGESTIONS: Improvements for better code
- POSITIVE: What's done well

Be specific and provide line numbers when possible. Suggest fixes.`,
}

// ExplainerAgent explains code
var ExplainerAgent = &BaseAgent{
	name:        "explainer",
	description: "Explain code in simple terms",
	prompt: `You are B.DEV Explainer, a patient code teacher.
Your role is to:
1. Explain what the code does in simple terms
2. Break down complex logic step by step
3. Explain the "why" behind design decisions
4. Define any technical terms used

Assume the reader is a beginner. Use analogies when helpful.
Keep explanations clear and concise.`,
}

// DocumenterAgent generates documentation
var DocumenterAgent = &BaseAgent{
	name:        "documenter",
	description: "Generate documentation for code",
	prompt: `You are B.DEV Documenter, a technical writer.
Your role is to:
1. Generate clear, concise documentation
2. Create function/method docstrings
3. Document parameters, return values, and exceptions
4. Add usage examples

Follow the documentation style for the language:
- Go: godoc-style comments
- Python: docstrings with Args/Returns/Raises
- JavaScript: JSDoc format
- Other: Simple markdown-style comments`,
}

// DebuggerAgent helps debug issues
var DebuggerAgent = &BaseAgent{
	name:        "debugger",
	description: "Debug and fix code issues",
	prompt: `You are B.DEV Debugger, an expert troubleshooter.
Your role is to:
1. Analyze the code for potential bugs
2. Identify the root cause of issues
3. Provide step-by-step debugging strategies
4. Suggest fixes with explanation

For each issue found:
- PROBLEM: What's wrong
- CAUSE: Why it happens
- FIX: How to fix it
- PREVENT: How to avoid similar issues`,
}

// ArchitectAgent provides architectural guidance
var ArchitectAgent = &BaseAgent{
	name:        "architect",
	description: "Design and architecture suggestions",
	prompt: `You are B.DEV Architect, a senior software architect.
Your role is to:
1. Analyze code structure and organization
2. Suggest architectural improvements
3. Recommend design patterns when appropriate
4. Identify scalability and maintainability concerns

Consider:
- Separation of concerns
- Single responsibility principle
- Dependency management
- Testability
- Future extensibility`,
}

// AllAgents is the list of all available agents
var AllAgents = []Agent{
	ReviewerAgent,
	ExplainerAgent,
	DocumenterAgent,
	DebuggerAgent,
	ArchitectAgent,
}

// GetAgent returns an agent by name.
// It performs a case-insensitive search through all registered agents.
// Returns nil if no agent matches the given name.
func GetAgent(name string) Agent {
	name = strings.ToLower(name)
	for _, a := range AllAgents {
		if strings.ToLower(a.Name()) == name {
			return a
		}
	}
	return nil
}

// ============================================================
// COMMAND
// ============================================================

var client *engine.Client

func getClient() *engine.Client {
	if client == nil {
		cfg := config.Get()
		client = engine.New(engine.Config{
			BaseURL:     cfg.AI.BaseURL,
			Model:       cfg.AI.Model,
			Fallback:    cfg.AI.FallbackModel,
			Timeout:     120 * 1e9,
			MaxRetries:  3,
			Temperature: 0.4, // Lower for more focused output
		})
	}
	return client
}

// NewCommand creates the agents command group
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "agent",
		Short: "AI agents for specialized tasks",
		Long:  "Use specialized AI agents for code review, documentation, debugging, and more",
	}

	// Add agent subcommands
	cmd.AddCommand(listCmd())
	cmd.AddCommand(reviewCmd())
	cmd.AddCommand(explainCmd())
	cmd.AddCommand(docCmd())
	cmd.AddCommand(debugCmd())
	cmd.AddCommand(architectCmd())

	return cmd
}

// ============================================================
// LIST - Show available agents
// ============================================================

func listCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List available agents",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(ui.Bold("Available Agents:"))
			for _, a := range AllAgents {
				fmt.Printf("  %s  %s\n", ui.Primary(a.Name()), ui.Muted(a.Description()))
			}
		},
	}
}

// ============================================================
// REVIEW - Code review agent
// ============================================================

func reviewCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "review <file>",
		Short: "Review code for issues and improvements",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAgentOnFile(ReviewerAgent, args[0])
		},
	}
}

// ============================================================
// EXPLAIN - Code explanation agent
// ============================================================

func explainCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "explain <file>",
		Short: "Explain code in simple terms",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAgentOnFile(ExplainerAgent, args[0])
		},
	}
}

// ============================================================
// DOC - Documentation agent
// ============================================================

func docCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "doc <file>",
		Short: "Generate documentation for code",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAgentOnFile(DocumenterAgent, args[0])
		},
	}
}

// ============================================================
// DEBUG - Debugging agent
// ============================================================

func debugCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "debug <file>",
		Short: "Analyze code for bugs and issues",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAgentOnFile(DebuggerAgent, args[0])
		},
	}
}

// ============================================================
// ARCHITECT - Architecture agent
// ============================================================

func architectCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "architect <file>",
		Short: "Analyze architecture and suggest improvements",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAgentOnFile(ArchitectAgent, args[0])
		},
	}
}

// ============================================================
// HELPERS
// ============================================================

func runAgentOnFile(agent Agent, filePath string) error {
	c := getClient()

	if !c.IsAvailable() {
		return fmt.Errorf("Ollama is not running. Start it with: ollama serve")
	}

	// Read file
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Get absolute path
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		absPath = filePath
	}

	// Get file extension for language hint
	ext := filepath.Ext(filePath)
	lang := getLanguageFromExt(ext)

	// Detect Project Context
	var projectContext string
	if rootDir := findProjectRoot(filepath.Dir(absPath)); rootDir != "" {
		if proj := projects.Analyze(rootDir); proj != nil {
			info := []string{
				fmt.Sprintf("Project Type: %s", proj.Type),
			}
			if proj.Framework != "" {
				info = append(info, fmt.Sprintf("Framework: %s", proj.Framework))
			}
			if proj.HasGit {
				info = append(info, fmt.Sprintf("Git: %s", proj.GitBranch))
			}
			projectContext = "\n\nCONTEXT:\n" + strings.Join(info, "\n")
		}
	}

	// Build prompt
	prompt := fmt.Sprintf("File: %s\nLanguage: %s%s\n\n```%s\n%s\n```",
		filePath, lang, projectContext, lang, string(content))

	// Prepare messages
	messages := []engine.Message{
		{Role: "system", Content: agent.SystemPrompt()},
		{Role: "user", Content: prompt},
	}

	fmt.Println(ui.Bold(fmt.Sprintf("[%s] Analyzing %s...", agent.Name(), filePath)))
	fmt.Println()

	// Create context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle Ctrl+C
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	go func() {
		<-sigChan
		cancel()
	}()

	// Stream response
	output, errChan := c.Chat(ctx, messages)

	for chunk := range output {
		fmt.Print(chunk)
	}
	fmt.Println()

	select {
	case err := <-errChan:
		return err
	default:
		return nil
	}
}

// findProjectRoot looks for markers like .git, package.json, go.mod
func findProjectRoot(startDir string) string {
	dir := startDir
	for {
		if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
			return dir
		}
		if _, err := os.Stat(filepath.Join(dir, "package.json")); err == nil {
			return dir
		}
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return ""
}

func getLanguageFromExt(ext string) string {
	languages := map[string]string{
		".go":   "go",
		".py":   "python",
		".js":   "javascript",
		".ts":   "typescript",
		".jsx":  "javascript",
		".tsx":  "typescript",
		".rs":   "rust",
		".rb":   "ruby",
		".php":  "php",
		".java": "java",
		".cpp":  "cpp",
		".c":    "c",
		".sh":   "bash",
		".ps1":  "powershell",
		".sql":  "sql",
		".json": "json",
		".yaml": "yaml",
		".yml":  "yaml",
		".md":   "markdown",
	}

	if lang, ok := languages[ext]; ok {
		return lang
	}
	return "text"
}
