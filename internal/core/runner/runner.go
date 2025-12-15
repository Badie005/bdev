package runner

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/badie/bdev/internal/core/projects"
	"github.com/badie/bdev/pkg/ui"
)

// Runner executes project commands
type Runner struct {
	Project *projects.Project
	Cwd     string
	Env     map[string]string
	Stdout  io.Writer
	Stderr  io.Writer
	Stdin   io.Reader
}

// New creates a new runner for a project
func New(project *projects.Project) *Runner {
	return &Runner{
		Project: project,
		Cwd:     project.Path,
		Env:     make(map[string]string),
		Stdout:  os.Stdout,
		Stderr:  os.Stderr,
		Stdin:   os.Stdin,
	}
}

// NewFromCwd creates a runner from current directory
func NewFromCwd() (*Runner, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	project := projects.Analyze(cwd)
	if project == nil {
		return nil, fmt.Errorf("not a recognized project in %s", cwd)
	}

	return New(project), nil
}

// Start starts the development server
func (r *Runner) Start() error {
	bin, args := r.Project.GetStartCommand()
	if bin == "" {
		return fmt.Errorf("no start command available for %s projects", r.Project.Type)
	}

	ui.PrintHeader(fmt.Sprintf("Starting %s", r.Project.Name))
	fmt.Printf("%s %s %v\n\n", ui.Muted("$"), ui.Primary(bin), args)

	return r.execute(bin, args)
}

// Test runs tests
func (r *Runner) Test(watch bool) error {
	bin, args := r.Project.GetTestCommand()
	if bin == "" {
		return fmt.Errorf("no test command available for %s projects", r.Project.Type)
	}

	// Add watch flag for supported frameworks
	if watch {
		switch r.Project.Type.String() {
		case "Next.js", "React", "Vue", "Node.js":
			args = append(args, "--", "--watch")
		case "Go":
			// Go doesn't have native watch, would need external tool
		}
	}

	ui.PrintHeader(fmt.Sprintf("Testing %s", r.Project.Name))
	fmt.Printf("%s %s %v\n\n", ui.Muted("$"), ui.Primary(bin), args)

	return r.execute(bin, args)
}

// Build builds for production
func (r *Runner) Build() error {
	bin, args := r.Project.GetBuildCommand()
	if bin == "" {
		return fmt.Errorf("no build command available for %s projects", r.Project.Type)
	}

	ui.PrintHeader(fmt.Sprintf("Building %s", r.Project.Name))
	fmt.Printf("%s %s %v\n\n", ui.Muted("$"), ui.Primary(bin), args)

	return r.execute(bin, args)
}

// RunScript runs a specific script
func (r *Runner) RunScript(script string) error {
	if r.Project.Scripts == nil {
		return fmt.Errorf("no scripts available")
	}

	if _, ok := r.Project.Scripts[script]; !ok {
		return fmt.Errorf("script '%s' not found. Available: %v", script, r.Project.ScriptList())
	}

	ui.PrintHeader(fmt.Sprintf("Running %s", script))
	return r.execute("npm", []string{"run", script})
}

// Install installs dependencies
func (r *Runner) Install() error {
	bin, args, err := r.GetInstallCommand()
	if err != nil {
		return err
	}

	ui.PrintHeader(fmt.Sprintf("Installing dependencies for %s", r.Project.Name))
	fmt.Printf("%s %s %v\n\n", ui.Muted("$"), ui.Primary(bin), args)

	return r.execute(bin, args)
}

// GetInstallCommand returns the install command
func (r *Runner) GetInstallCommand() (string, []string, error) {
	switch r.Project.Type {
	case projects.TypeNextJS, projects.TypeReact, projects.TypeVue,
		projects.TypeNuxt, projects.TypeSvelte, projects.TypeAstro, projects.TypeNode:
		return "npm", []string{"install"}, nil
	case projects.TypeGo:
		return "go", []string{"mod", "download"}, nil
	case projects.TypePython:
		return "pip", []string{"install", "-r", "requirements.txt"}, nil
	case projects.TypeRust:
		return "cargo", []string{"build"}, nil
	case projects.TypeLaravel:
		return "composer", []string{"install"}, nil
	default:
		return "", nil, fmt.Errorf("no install command for %s", r.Project.Type)
	}
}

// Lint runs the linter
func (r *Runner) Lint(fix bool) error {
	bin, args, err := r.GetLintCommand(fix)
	if err != nil {
		return err
	}

	ui.PrintHeader(fmt.Sprintf("Linting %s", r.Project.Name))
	fmt.Printf("%s %s %v\n\n", ui.Muted("$"), ui.Primary(bin), args)

	return r.execute(bin, args)
}

// GetLintCommand returns the lint command
func (r *Runner) GetLintCommand(fix bool) (string, []string, error) {
	var bin string
	var args []string

	switch r.Project.Type {
	case projects.TypeNextJS, projects.TypeReact, projects.TypeVue,
		projects.TypeNuxt, projects.TypeSvelte, projects.TypeNode:
		if _, ok := r.Project.Scripts["lint"]; ok {
			args = []string{"run", "lint"}
			if fix {
				args = append(args, "--", "--fix")
			}
			bin = "npm"
		} else {
			bin = "npx"
			args = []string{"eslint", "."}
			if fix {
				args = append(args, "--fix")
			}
		}
		return bin, args, nil
	case projects.TypeGo:
		return "go", []string{"vet", "./..."}, nil
	case projects.TypePython:
		bin = "ruff"
		args = []string{"check", "."}
		if fix {
			args = append(args, "--fix")
		}
		return bin, args, nil
	default:
		return "", nil, fmt.Errorf("no lint command for %s", r.Project.Type)
	}
}

// execute runs a command with streaming output
func (r *Runner) execute(bin string, args []string) error {
	cmd := exec.Command(bin, args...)
	cmd.Dir = r.Cwd
	cmd.Stdout = r.Stdout
	cmd.Stderr = r.Stderr
	cmd.Stdin = r.Stdin

	// Add custom environment
	cmd.Env = os.Environ()
	for k, v := range r.Env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}

	err := cmd.Run()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return fmt.Errorf("command exited with code %d", exitErr.ExitCode())
		}
		return err
	}

	return nil
}

// ExecuteRaw executes a raw command in the project directory
func (r *Runner) ExecuteRaw(cmdLine string, args []string) error {
	return r.execute(cmdLine, args)
}
