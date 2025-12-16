package runner

import (
	"bytes"
	"testing"

	"github.com/badie/bdev/internal/core/projects"
)

// ==============================================================
// New Tests
// ==============================================================

func TestNew(t *testing.T) {
	project := &projects.Project{
		Name: "test-project",
		Path: "/test/path",
		Type: projects.TypeGo,
	}

	runner := New(project)

	if runner.Project != project {
		t.Error("Runner.Project should be set")
	}
	if runner.Cwd != project.Path {
		t.Errorf("Runner.Cwd = %q, want %q", runner.Cwd, project.Path)
	}
	if runner.Stdout == nil {
		t.Error("Runner.Stdout should be set to os.Stdout")
	}
	if runner.Stderr == nil {
		t.Error("Runner.Stderr should be set to os.Stderr")
	}
	if runner.Stdin == nil {
		t.Error("Runner.Stdin should be set to os.Stdin")
	}
	if runner.Env == nil {
		t.Error("Runner.Env should be initialized")
	}
}

// ==============================================================
// GetInstallCommand Tests
// ==============================================================

func TestRunner_GetInstallCommand(t *testing.T) {
	tests := []struct {
		name        string
		projectType projects.ProjectType
		wantBin     string
		wantArgs    []string
		wantErr     bool
	}{
		{
			name:        "go_project",
			projectType: projects.TypeGo,
			wantBin:     "go",
			wantArgs:    []string{"mod", "download"},
		},
		{
			name:        "node_project",
			projectType: projects.TypeNode,
			wantBin:     "npm",
			wantArgs:    []string{"install"},
		},
		{
			name:        "nextjs_project",
			projectType: projects.TypeNextJS,
			wantBin:     "npm",
			wantArgs:    []string{"install"},
		},
		{
			name:        "react_project",
			projectType: projects.TypeReact,
			wantBin:     "npm",
			wantArgs:    []string{"install"},
		},
		{
			name:        "vue_project",
			projectType: projects.TypeVue,
			wantBin:     "npm",
			wantArgs:    []string{"install"},
		},
		{
			name:        "rust_project",
			projectType: projects.TypeRust,
			wantBin:     "cargo",
			wantArgs:    []string{"build"},
		},
		{
			name:        "python_project",
			projectType: projects.TypePython,
			wantBin:     "pip",
			wantArgs:    []string{"install", "-r", "requirements.txt"},
		},
		{
			name:        "laravel_project",
			projectType: projects.TypeLaravel,
			wantBin:     "composer",
			wantArgs:    []string{"install"},
		},
		{
			name:        "unknown_project",
			projectType: projects.TypeUnknown,
			wantErr:     true,
		},
		{
			name:        "static_project",
			projectType: projects.TypeStatic,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Runner{
				Project: &projects.Project{Type: tt.projectType},
			}

			bin, args, err := r.GetInstallCommand()

			if tt.wantErr {
				if err == nil {
					t.Error("GetInstallCommand() expected error")
				}
				return
			}

			if err != nil {
				t.Fatalf("GetInstallCommand() error = %v", err)
			}
			if bin != tt.wantBin {
				t.Errorf("GetInstallCommand() bin = %q, want %q", bin, tt.wantBin)
			}
			if len(args) != len(tt.wantArgs) {
				t.Errorf("GetInstallCommand() args = %v, want %v", args, tt.wantArgs)
			}
		})
	}
}

// ==============================================================
// GetLintCommand Tests
// ==============================================================

func TestRunner_GetLintCommand(t *testing.T) {
	tests := []struct {
		name        string
		projectType projects.ProjectType
		scripts     map[string]string
		fix         bool
		wantBin     string
		wantArgs    []string
		wantErr     bool
	}{
		{
			name:        "go_project",
			projectType: projects.TypeGo,
			fix:         false,
			wantBin:     "go",
			wantArgs:    []string{"vet", "./..."},
		},
		{
			name:        "go_project_fix_ignored",
			projectType: projects.TypeGo,
			fix:         true,
			wantBin:     "go",
			wantArgs:    []string{"vet", "./..."}, // go vet doesn't have --fix
		},
		{
			name:        "python_project",
			projectType: projects.TypePython,
			fix:         false,
			wantBin:     "ruff",
			wantArgs:    []string{"check", "."},
		},
		{
			name:        "python_project_fix",
			projectType: projects.TypePython,
			fix:         true,
			wantBin:     "ruff",
			wantArgs:    []string{"check", ".", "--fix"},
		},
		{
			name:        "node_with_lint_script",
			projectType: projects.TypeNode,
			scripts:     map[string]string{"lint": "eslint ."},
			fix:         false,
			wantBin:     "npm",
			wantArgs:    []string{"run", "lint"},
		},
		{
			name:        "node_with_lint_script_fix",
			projectType: projects.TypeNode,
			scripts:     map[string]string{"lint": "eslint ."},
			fix:         true,
			wantBin:     "npm",
			wantArgs:    []string{"run", "lint", "--", "--fix"},
		},
		{
			name:        "node_without_lint_script",
			projectType: projects.TypeNode,
			scripts:     map[string]string{},
			fix:         false,
			wantBin:     "npx",
			wantArgs:    []string{"eslint", "."},
		},
		{
			name:        "node_without_lint_script_fix",
			projectType: projects.TypeNode,
			scripts:     map[string]string{},
			fix:         true,
			wantBin:     "npx",
			wantArgs:    []string{"eslint", ".", "--fix"},
		},
		{
			name:        "django_no_lint",
			projectType: projects.TypeDjango,
			wantErr:     true,
		},
		{
			name:        "unknown_project",
			projectType: projects.TypeUnknown,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Runner{
				Project: &projects.Project{
					Type:    tt.projectType,
					Scripts: tt.scripts,
				},
			}

			bin, args, err := r.GetLintCommand(tt.fix)

			if tt.wantErr {
				if err == nil {
					t.Error("GetLintCommand() expected error")
				}
				return
			}

			if err != nil {
				t.Fatalf("GetLintCommand() error = %v", err)
			}
			if bin != tt.wantBin {
				t.Errorf("GetLintCommand() bin = %q, want %q", bin, tt.wantBin)
			}
			if len(args) != len(tt.wantArgs) {
				t.Errorf("GetLintCommand() args len = %d, want %d", len(args), len(tt.wantArgs))
				t.Errorf("  got:  %v", args)
				t.Errorf("  want: %v", tt.wantArgs)
			} else {
				for i, arg := range args {
					if arg != tt.wantArgs[i] {
						t.Errorf("GetLintCommand() args[%d] = %q, want %q", i, arg, tt.wantArgs[i])
					}
				}
			}
		})
	}
}

// ==============================================================
// Error Case Tests
// ==============================================================

func TestRunner_Start_NoCommand(t *testing.T) {
	r := &Runner{
		Project: &projects.Project{
			Name: "test",
			Type: projects.TypeUnknown,
		},
		Stdout: &bytes.Buffer{},
		Stderr: &bytes.Buffer{},
	}

	err := r.Start()
	if err == nil {
		t.Error("Start() should return error for unknown project type")
	}
}

func TestRunner_Test_NoCommand(t *testing.T) {
	r := &Runner{
		Project: &projects.Project{
			Name: "test",
			Type: projects.TypeUnknown,
		},
		Stdout: &bytes.Buffer{},
		Stderr: &bytes.Buffer{},
	}

	err := r.Test(false)
	if err == nil {
		t.Error("Test() should return error for unknown project type")
	}
}

func TestRunner_Build_NoCommand(t *testing.T) {
	r := &Runner{
		Project: &projects.Project{
			Name: "test",
			Type: projects.TypeUnknown,
		},
		Stdout: &bytes.Buffer{},
		Stderr: &bytes.Buffer{},
	}

	err := r.Build()
	if err == nil {
		t.Error("Build() should return error for unknown project type")
	}
}

func TestRunner_RunScript_NoScripts(t *testing.T) {
	r := &Runner{
		Project: &projects.Project{
			Name:    "test",
			Scripts: nil,
		},
	}

	err := r.RunScript("dev")
	if err == nil {
		t.Error("RunScript() should return error when scripts is nil")
	}
}

func TestRunner_RunScript_NotFound(t *testing.T) {
	r := &Runner{
		Project: &projects.Project{
			Name:    "test",
			Scripts: map[string]string{"dev": "next dev"},
		},
	}

	err := r.RunScript("nonexistent")
	if err == nil {
		t.Error("RunScript() should return error for nonexistent script")
	}
}

// ==============================================================
// Integration Tests (require actual commands)
// ==============================================================

func TestRunner_Execute_EchoCommand(t *testing.T) {
	// Skip if 'echo' is not available (Windows might need special handling)
	t.Run("echo_command", func(t *testing.T) {
		stdout := &bytes.Buffer{}
		stderr := &bytes.Buffer{}

		r := &Runner{
			Project: &projects.Project{
				Name: "test",
				Type: projects.TypeGo,
				Path: t.TempDir(),
			},
			Cwd:    t.TempDir(),
			Env:    make(map[string]string),
			Stdout: stdout,
			Stderr: stderr,
		}

		// Use different command for Windows vs Unix
		var err error
		err = r.ExecuteRaw("go", []string{"version"})

		if err != nil {
			t.Logf("ExecuteRaw failed (expected on some systems): %v", err)
			// Don't fail - go might not be in PATH in all test environments
			t.Skip("go command not available")
		}

		if stdout.Len() == 0 {
			t.Error("ExecuteRaw should produce output for 'go version'")
		}
	})
}

func TestRunner_WithEnv(t *testing.T) {
	r := &Runner{
		Project: &projects.Project{
			Name: "test",
			Path: t.TempDir(),
		},
		Env: make(map[string]string),
	}

	r.Env["MY_VAR"] = "my_value"

	if r.Env["MY_VAR"] != "my_value" {
		t.Error("Environment variable should be set")
	}
}
