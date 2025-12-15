package repl

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/chzyer/readline"
	"github.com/spf13/cobra"

	"github.com/badie/bdev/internal/core/config"
)

// Completer provides tab completion for the REPL
type Completer struct {
	rootCmd  *cobra.Command
	config   *config.Config
	projects []string
}

// NewCompleter creates a new completer
func NewCompleter(rootCmd *cobra.Command, cfg *config.Config) *Completer {
	c := &Completer{
		rootCmd: rootCmd,
		config:  cfg,
	}
	c.loadProjects()
	return c
}

func (c *Completer) loadProjects() {
	projectsDir := c.config.Paths.Projects
	entries, err := os.ReadDir(projectsDir)
	if err != nil {
		c.projects = make([]string, 0)
		return
	}

	c.projects = make([]string, 0)
	for _, entry := range entries {
		if entry.IsDir() && !strings.HasPrefix(entry.Name(), ".") {
			c.projects = append(c.projects, entry.Name())
		}
	}
}

// RefreshProjects reloads the project list
func (c *Completer) RefreshProjects() {
	c.loadProjects()
}

// Do implements readline.AutoCompleter
func (c *Completer) Do(line []rune, pos int) (newLine [][]rune, length int) {
	lineStr := string(line[:pos])
	parts := strings.Fields(lineStr)

	var candidates []string
	var prefix string

	if len(parts) == 0 {
		// Empty line - show all top-level commands
		candidates = c.getTopLevelCommands()
		prefix = ""
	} else if len(parts) == 1 && !strings.HasSuffix(lineStr, " ") {
		// Completing first word
		candidates = c.getTopLevelCommands()
		prefix = parts[0]
	} else {
		// Completing subsequent words
		candidates = c.getContextualCompletions(parts)
		if strings.HasSuffix(lineStr, " ") {
			prefix = ""
		} else {
			prefix = parts[len(parts)-1]
		}
	}

	// Filter candidates by prefix
	var matches [][]rune
	for _, cand := range candidates {
		if strings.HasPrefix(strings.ToLower(cand), strings.ToLower(prefix)) {
			suffix := cand[len(prefix):]
			matches = append(matches, []rune(suffix+" "))
		}
	}

	return matches, len(prefix)
}

func (c *Completer) getTopLevelCommands() []string {
	// Dynamically fetch from rootCmd if possible, or keep this robust static list
	cmds := []string{
		// Main command groups
		"projects", "git", "ai", "agents", "workflow",
		"secrets", "multi", "config", "theme", "analytics",
		// Quick actions
		"list", "start", "test", "build", "fix", "deploy", "do", "install", "lint", "clean",
		// REPL built-ins
		"help", "exit", "quit", "clear", "cls", "history", "status", "reload",
		"version", "cd",
	}
	// Add sub-commands as top-level helpers if user desires, but standard is strict structure
	return cmds
}

func (c *Completer) getContextualCompletions(parts []string) []string {
	cmd := strings.ToLower(parts[0])

	switch cmd {
	case "projects":
		return c.getProjectsCompletions(parts)
	case "git":
		return c.getGitCompletions(parts)
	case "ai":
		return c.getAICompletions(parts)
	case "agents":
		return c.getAgentsCompletions(parts)
	case "workflow":
		return []string{"list", "run", "create", "show", "edit"}
	case "secrets":
		return []string{"init", "set", "get", "list", "export", "delete"}
	case "multi":
		return []string{"status", "pull", "audit", "run", "update"}
	case "config":
		return []string{"show", "set", "alias", "edit", "reset"}
	case "theme":
		return []string{"list", "set", "preview", "claude", "gemini", "matrix"}
	case "analytics":
		return []string{"today", "week", "month", "summary", "reset"}
	case "start", "test", "build", "run", "open", "describe":
		return c.projects
	case "cd":
		return c.getDirectoryCompletions(parts)
	}

	return nil
}

func (c *Completer) getProjectsCompletions(parts []string) []string {
	if len(parts) == 1 {
		return []string{"list", "new", "open", "find", "describe", "run", "delete"}
	}

	subcmd := strings.ToLower(parts[1])
	switch subcmd {
	case "open", "run", "describe", "delete":
		return c.projects
	case "new":
		return c.getTemplates()
	}

	return nil
}

func (c *Completer) getGitCompletions(parts []string) []string {
	if len(parts) == 1 {
		return []string{
			"status", "commit", "push", "pull", "fetch",
			"branch", "checkout", "merge", "rebase",
			"log", "diff", "stash", "add", "reset",
			"clone", "remote", "tag",
		}
	}

	subcmd := strings.ToLower(parts[1])
	switch subcmd {
	case "checkout", "branch", "merge", "rebase":
		return c.getGitBranches()
	case "stash":
		return []string{"push", "pop", "list", "drop", "apply", "clear"}
	case "remote":
		return []string{"add", "remove", "show", "rename"}
	}

	return nil
}

func (c *Completer) getAICompletions(parts []string) []string {
	if len(parts) == 1 {
		return []string{"check", "chat", "forget", "memory", "generate", "models"}
	}
	return nil
}

func (c *Completer) getAgentsCompletions(parts []string) []string {
	if len(parts) == 1 {
		return []string{"list", "review", "document", "explain", "architect", "debug", "refactor"}
	}

	// For file-based commands, try to complete files
	subcmd := strings.ToLower(parts[1])
	switch subcmd {
	case "review", "document", "explain", "debug", "refactor":
		return c.getFileCompletions()
	}

	return nil
}

func (c *Completer) getTemplates() []string {
	templatesDir := filepath.Join(c.config.Paths.Bdev, "..", "bdev", "templates")
	entries, err := os.ReadDir(templatesDir)
	if err != nil {
		return []string{"nextjs-starter", "python-cli", "go-cli", "laravel-api"}
	}

	templates := make([]string, 0)
	for _, entry := range entries {
		if entry.IsDir() {
			templates = append(templates, entry.Name())
		}
	}
	return templates
}

func (c *Completer) getGitBranches() []string {
	// This would execute git branch -a, but for performance we return common ones
	return []string{"main", "master", "develop", "feature/", "bugfix/", "hotfix/"}
}

func (c *Completer) getDirectoryCompletions(_ []string) []string {
	cwd, _ := os.Getwd()
	entries, err := os.ReadDir(cwd)
	if err != nil {
		return nil
	}

	dirs := make([]string, 0)
	for _, entry := range entries {
		if entry.IsDir() {
			dirs = append(dirs, entry.Name())
		}
	}
	dirs = append(dirs, "..", "~")
	return dirs
}

func (c *Completer) getFileCompletions() []string {
	cwd, _ := os.Getwd()
	entries, err := os.ReadDir(cwd)
	if err != nil {
		return nil
	}

	files := make([]string, 0)
	for _, entry := range entries {
		if !entry.IsDir() {
			files = append(files, entry.Name())
		}
	}
	return files
}

// Compile-time check that Completer implements readline.AutoCompleter
var _ readline.AutoCompleter = (*Completer)(nil)
