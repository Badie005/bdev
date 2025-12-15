package multi

import (
	"context"
	"os/exec"
	"sync"
	"time"

	"github.com/badie/bdev/internal/core/projects"
)

// Executor manages multi-project execution
type Executor struct {
	Projects []projects.Project
	MaxJobs  int
}

// Result represents a single execution result
type Result struct {
	Project  projects.Project
	Output   string
	Error    error
	Duration time.Duration
	Skipped  bool
}

// Filter defines criteria for selecting projects
type Filter struct {
	Names    []string
	Types    []projects.ProjectType
	Patterns []string
}

// New creates a new executor
func New(allProjects []projects.Project) *Executor {
	return &Executor{
		Projects: allProjects,
		MaxJobs:  4, // Default concurrency
	}
}

// Filter applies a filter to the project list
func (e *Executor) Filter(f Filter) {
	if len(f.Names) == 0 && len(f.Types) == 0 && len(f.Patterns) == 0 {
		return
	}

	var filtered []projects.Project

	for _, p := range e.Projects {
		match := false

		// Check names
		if len(f.Names) > 0 {
			for _, n := range f.Names {
				if p.Name == n {
					match = true
					break
				}
			}
		}

		// Check types
		if !match && len(f.Types) > 0 {
			for _, t := range f.Types {
				if p.Type == t {
					match = true
					break
				}
			}
		}

		// If no filters matched but filters were present, skip
		// Simplistic logic: OR within categories, AND logic is tricky.
		// Let's implement OR logic for now: if any criteria matches, include.
		// If NO criteria provided, include all (handled at top).

		if match {
			filtered = append(filtered, p)
		}
	}
	e.Projects = filtered
}

// Execute runs a command on all selected projects
func (e *Executor) Execute(ctx context.Context, command string, args []string) <-chan Result {
	results := make(chan Result, len(e.Projects))
	sem := make(chan struct{}, e.MaxJobs) // Semaphore for concurrency limitation

	go func() {
		var wg sync.WaitGroup
		defer close(results)

		for _, p := range e.Projects {
			wg.Add(1)
			sem <- struct{}{} // Acquire

			go func(proj projects.Project) {
				defer wg.Done()
				defer func() { <-sem }() // Release

				start := time.Now()

				// Create command
				cmd := exec.CommandContext(ctx, command, args...)
				cmd.Dir = proj.Path

				// Capture output
				output, err := cmd.CombinedOutput()

				results <- Result{
					Project:  proj,
					Output:   string(output),
					Error:    err,
					Duration: time.Since(start),
				}
			}(p)
		}

		wg.Wait()
	}()

	return results
}

// ExecuteShell runs a shell command string on all projects
func (e *Executor) ExecuteShell(ctx context.Context, shellCmd string) <-chan Result {
	// Simple wrapper, handling windows/unix shell differences could be done here if needed
	// For now, assume calling `cmd /c` or `bash -c` is handled by caller or we standardize
	// Actually, let's standardize here.

	// BUT, typical usage is `bdev multi exec "npm install"`.
	// We'll treat the input as a command + args

	// Wait, if it's a shell string, we need to parse it or pass to shell.
	// Let's assume the user passes command and args separately to Execute,
	// ExecuteShell is specifically for things like `echo "hello" && ls`.

	return e.Execute(ctx, "powershell", []string{"-Command", shellCmd})
}
