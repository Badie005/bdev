package repl

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/chzyer/readline"
	"github.com/spf13/cobra"

	"github.com/badie/bdev/internal/core/config"
	"github.com/badie/bdev/internal/core/session"
	"github.com/badie/bdev/pkg/ui"
)

// REPL represents the interactive Read-Eval-Print Loop
type REPL struct {
	rootCmd   *cobra.Command
	session   *session.Session
	config    *config.Config
	readline  *readline.Instance
	completer *Completer
}

// Start initializes and runs the REPL
func Start(rootCmd *cobra.Command) {
	r := &REPL{
		rootCmd: rootCmd,
		session: session.Get(),
		config:  config.Get(),
	}

	// Ensure directories exist
	_ = r.config.EnsureDirectories()

	r.setupReadline()
	r.printBanner()
	r.run()
}

func (r *REPL) setupReadline() {
	r.completer = NewCompleter(r.rootCmd, r.config)

	rl, err := readline.NewEx(&readline.Config{
		Prompt:                 r.getPrompt(),
		HistoryFile:            r.config.HistoryFile(),
		AutoComplete:           r.completer,
		InterruptPrompt:        "^C",
		EOFPrompt:              "exit",
		HistorySearchFold:      true,
		DisableAutoSaveHistory: false,
	})
	if err != nil {
		// Fallback to basic readline if fancy one fails
		fmt.Println(ui.Warning("Warning: Advanced readline unavailable, using basic mode"))
		r.readline = nil
		return
	}

	r.readline = rl
}

func (r *REPL) printBanner() {
	ui.AnimateBoot()
	ui.PrintWelcome()
}

func (r *REPL) getPrompt() string {
	r.session.UpdateCurrentDir()

	branch := r.session.GitBranch()
	project := r.session.ProjectName()

	if branch != "" {
		return ui.Primary(fmt.Sprintf("┃ %s ", project)) +
			ui.Info(fmt.Sprintf(" %s %s", ui.ActiveGlyphs.Branch, branch)) +
			ui.Primary(" "+ui.ActiveGlyphs.Pointer+" ")
	}
	return ui.Primary(fmt.Sprintf("┃ %s %s ", project, ui.ActiveGlyphs.Pointer))
}

func (r *REPL) run() {
	if r.readline != nil {
		r.runWithReadline()
	} else {
		r.runBasic()
	}
}

func (r *REPL) runWithReadline() {
	defer r.readline.Close()

	for {
		// Update prompt with current context
		r.readline.SetPrompt(r.getPrompt())

		line, err := r.readline.Readline()
		if err == readline.ErrInterrupt {
			if line == "" {
				fmt.Println(ui.Muted("\nUse 'exit' to quit"))
				continue
			}
			continue
		} else if err == io.EOF {
			break
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Handle built-in REPL commands
		if r.handleBuiltin(line) {
			continue
		}

		// Resolve aliases
		line = r.resolveAlias(line)

		// Execute command
		r.execute(line)
	}

	r.printGoodbye()
}

func (r *REPL) runBasic() {
	// Fallback for when readline is not available
	var input string
	for {
		fmt.Print(r.getPrompt())
		_, err := fmt.Scanln(&input)
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			continue
		}

		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		if r.handleBuiltin(input) {
			continue
		}

		r.execute(r.resolveAlias(input))
	}
	r.printGoodbye()
}

func (r *REPL) handleBuiltin(line string) bool {
	lower := strings.ToLower(line)

	switch lower {
	case "exit", "quit", "q":
		r.printGoodbye()
		os.Exit(0)
		return true

	case "clear", "cls":
		fmt.Print("\033[H\033[2J") // ANSI clear screen
		r.printBanner()
		return true

	case "help", "?":
		_ = r.rootCmd.Help()
		return true

	case "history":
		r.showHistory()
		return true

	case "status", "stats":
		r.showStatus()
		return true

	case "reload":
		r.config = config.Load()
		fmt.Println(ui.Success("Configuration reloaded"))
		return true
	}

	// Handle cd command specially
	if strings.HasPrefix(lower, "cd ") {
		dir := strings.TrimPrefix(line, "cd ")
		dir = strings.TrimPrefix(dir, "CD ")
		if err := os.Chdir(strings.TrimSpace(dir)); err != nil {
			fmt.Println(ui.Error("Error: " + err.Error()))
		} else {
			r.session.UpdateCurrentDir()
		}
		return true
	}

	return false
}

func (r *REPL) resolveAlias(line string) string {
	parts := strings.Fields(line)
	if len(parts) == 0 {
		return line
	}

	if alias, ok := r.config.Aliases[parts[0]]; ok {
		if len(parts) > 1 {
			return alias + " " + strings.Join(parts[1:], " ")
		}
		return alias
	}
	return line
}

func (r *REPL) execute(line string) {
	// Record command
	r.session.SetLastCommand(line)
	r.session.IncrementCommandCount()

	// Parse line into args
	args := strings.Fields(line)
	if len(args) == 0 {
		return
	}

	// UX: Strip 'bdev' prefix if user typed it inside the REPL
	if strings.EqualFold(args[0], "bdev") {
		if len(args) > 1 {
			args = args[1:]
		} else {
			// Just 'bdev' typed? Show help or ignore
			return
		}
	}

	// Create a new command instance to avoid state pollution
	cmd := r.rootCmd
	cmd.SetArgs(args)

	// Capture output for special handling
	if err := cmd.Execute(); err != nil {
		// Don't print the error again if cobra already did
		if !strings.Contains(err.Error(), "unknown command") {
			fmt.Println(ui.Error("Error: " + err.Error()))
		}
	}

	fmt.Println() // Add spacing after command
}

func (r *REPL) showHistory() {
	if r.readline == nil {
		fmt.Println(ui.Muted("History not available in basic mode"))
		return
	}

	// Read history file
	data, err := os.ReadFile(r.config.HistoryFile())
	if err != nil {
		fmt.Println(ui.Muted("No history yet"))
		return
	}

	lines := strings.Split(string(data), "\n")
	start := len(lines) - 20
	if start < 0 {
		start = 0
	}

	fmt.Println(ui.Bold("Recent commands:"))
	for idx, line := range lines[start:] {
		if line != "" {
			fmt.Printf("  %s %s\n", ui.Muted(fmt.Sprintf("%3d", start+idx+1)), line)
		}
	}
}

func (r *REPL) showStatus() {
	stats := r.session.Stats()

	fmt.Println(ui.Bold("Session Status:"))
	fmt.Printf("  %s %v\n", ui.Primary("Uptime:"), stats["uptime"])
	fmt.Printf("  %s %v\n", ui.Primary("Commands:"), stats["commands"])
	fmt.Printf("  %s %v\n", ui.Primary("Project:"), stats["project"])

	if stats["is_git_repo"].(bool) {
		fmt.Printf("  %s %v\n", ui.Primary("Branch:"), stats["git_branch"])
	}

	if stats["ai_messages"].(int) > 0 {
		fmt.Printf("  %s %v messages\n", ui.Primary("AI Context:"), stats["ai_messages"])
	}
}

func (r *REPL) printGoodbye() {
	fmt.Println()
	fmt.Println(ui.Muted(ui.SeparatorLight))
	fmt.Printf("  Session: %d commands in %s\n",
		r.session.CommandCount,
		r.session.UptimeFormatted())
	fmt.Println(ui.Primary("  Goodbye!"))
	fmt.Println(ui.Muted(ui.SeparatorLight))
}
