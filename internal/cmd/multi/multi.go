package multicmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/badie/bdev/internal/core/config"
	"github.com/badie/bdev/internal/core/multi"
	"github.com/badie/bdev/internal/core/projects"
	"github.com/badie/bdev/pkg/ui"
)

var (
	filterType  string
	filterName  []string
	concurrency int
)

// NewCommand creates the multi-project command group
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "multi",
		Aliases: []string{"m", "batch"},
		Short:   "Run commands across multiple projects",
		Long:    "Execute commands concurrently across all discovered projects",
	}

	cmd.PersistentFlags().StringVarP(&filterType, "type", "t", "", "Filter by project type (e.g., go, node, py)")
	cmd.PersistentFlags().StringSliceVarP(&filterName, "name", "n", nil, "Filter by project name (exact)")
	cmd.PersistentFlags().IntVarP(&concurrency, "concurrency", "c", 4, "Maximum concurrent jobs")

	cmd.AddCommand(execCmd())
	cmd.AddCommand(listCmd())
	cmd.AddCommand(gitCmd())

	return cmd
}

// getExecutor scans projects and returns a configured executor
func getExecutor() (*multi.Executor, error) {
	cfg := config.Get()
	allProjects, err := projects.Scan(cfg.Paths.Projects)
	if err != nil {
		return nil, err
	}

	exec := multi.New(allProjects)
	exec.MaxJobs = concurrency

	// Apply filter
	filter := multi.Filter{}
	if len(filterName) > 0 {
		filter.Names = filterName
	}

	// Type filter parsing
	// This is basic, might need better mapping from string to ProjectType
	// For now, we rely on the user knowing internal types or we improve ProjectType parsing.
	// Since ProjectType is an int constant, mapping string "go" -> TypeGo is needed.
	// We'll skip complex type filtering for this MVP step and accept basic filtering logic later if needed.
	// Or we can do simple string matching on ProjectType.String()

	if filterType != "" {
		// Filter manually for now as Filter struct expects ProjectType enum
		var filtered []projects.Project
		for _, p := range exec.Projects {
			if strings.Contains(strings.ToLower(p.Type.String()), strings.ToLower(filterType)) {
				filtered = append(filtered, p)
			}
		}
		exec.Projects = filtered
	}

	exec.Filter(filter)

	return exec, nil
}

// ============================================================
// EXEC - Run arbitrary command
// ============================================================

func execCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "exec <command> [args...]",
		Short: "Execute a command on selected projects",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			execCtx, err := getExecutor()
			if err != nil {
				return err
			}

			if len(execCtx.Projects) == 0 {
				fmt.Println(ui.Warning("No projects matched the filter"))
				return nil
			}

			fmt.Println(ui.Bold(fmt.Sprintf("Executing on %d projects...", len(execCtx.Projects))))
			fmt.Println()

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			// Handle Ctrl+C
			sigChan := make(chan os.Signal, 1)
			signal.Notify(sigChan, os.Interrupt)
			go func() {
				<-sigChan
				cancel()
			}()

			command := args[0]
			cmdArgs := args[1:]

			results := execCtx.Execute(ctx, command, cmdArgs)

			successCount := 0
			failCount := 0

			for res := range results {
				if res.Error != nil {
					failCount++
					fmt.Printf("%s %s %s\n", ui.Error(ui.ActiveGlyphs.Cross), ui.Bold(res.Project.Name), ui.Muted(fmt.Sprintf("(%s)", res.Duration.Round(time.Millisecond))))
					fmt.Println(ui.Error(fmt.Sprintf("  Error: %v", res.Error)))
				} else {
					successCount++
					fmt.Printf("%s %s %s\n", ui.Success(ui.ActiveGlyphs.Check), ui.Bold(res.Project.Name), ui.Muted(fmt.Sprintf("(%s)", res.Duration.Round(time.Millisecond))))
				}
			}

			fmt.Println()
			fmt.Printf("Summary: %d success, %d failed\n", successCount, failCount)

			if failCount > 0 {
				return fmt.Errorf("some commands failed")
			}
			return nil
		},
	}
}

// ============================================================
// GIT - Quick git commands
// ============================================================

func gitCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "git <command> [args...]",
		Short: "Execute git command on selected projects",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// We reuse the exec logic manually to avoid command lookup complexity
			// Or we can just call the exec impl directly?
			// Let's copy logic for simplicity or extract execution logic.

			// Extracting logic:
			execCtx, err := getExecutor()
			if err != nil {
				return err
			}

			if len(execCtx.Projects) == 0 {
				fmt.Println(ui.Warning("No projects matched the filter"))
				return nil
			}

			// For Git, only run on Projects that HAVE git
			var gitProjects []projects.Project
			for _, p := range execCtx.Projects {
				if p.HasGit {
					gitProjects = append(gitProjects, p)
				}
			}
			execCtx.Projects = gitProjects

			if len(execCtx.Projects) == 0 {
				fmt.Println(ui.Warning("No git projects matched"))
				return nil
			}

			fmt.Println(ui.Bold(fmt.Sprintf("Running git %s on %d projects...", args[0], len(execCtx.Projects))))

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			sigChan := make(chan os.Signal, 1)
			signal.Notify(sigChan, os.Interrupt)
			go func() {
				<-sigChan
				cancel()
			}()

			results := execCtx.Execute(ctx, "git", args)

			successCount := 0
			failCount := 0

			for res := range results {
				if res.Error != nil {
					failCount++
					fmt.Printf("%s %s\n", ui.Error(ui.ActiveGlyphs.Cross), res.Project.Name)
				} else {
					successCount++
					fmt.Printf("%s %s\n", ui.Success(ui.ActiveGlyphs.Check), res.Project.Name)
				}
			}

			return nil
		},
	}
}

// ============================================================
// LIST - Show matched projects
// ============================================================

func listCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List projects selected by filter",
		RunE: func(cmd *cobra.Command, args []string) error {
			execCtx, err := getExecutor()
			if err != nil {
				return err
			}

			if len(execCtx.Projects) == 0 {
				fmt.Println(ui.Warning("No projects matched"))
				return nil
			}

			fmt.Println(ui.Bold(fmt.Sprintf("Matching Projects (%d):", len(execCtx.Projects))))
			for _, p := range execCtx.Projects {
				fmt.Printf("  %s %s\n", p.Type.Icon(), p.Name)
			}

			return nil
		},
	}
}
