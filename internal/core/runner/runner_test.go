package runner

import (
	"testing"

	"github.com/badie/bdev/internal/core/projects"
)

func TestGetInstallCommand(t *testing.T) {
	tests := []struct {
		name     string
		pt       projects.ProjectType
		wantCmd  string
		wantArgs []string
		wantErr  bool
	}{
		{
			name:     "Node Install",
			pt:       projects.TypeNode,
			wantCmd:  "npm",
			wantArgs: []string{"install"},
			wantErr:  false,
		},
		{
			name:     "Go Install",
			pt:       projects.TypeGo,
			wantCmd:  "go",
			wantArgs: []string{"mod", "download"},
			wantErr:  false,
		},
		{
			name:     "Python Install",
			pt:       projects.TypePython,
			wantCmd:  "pip",
			wantArgs: []string{"install", "-r", "requirements.txt"},
			wantErr:  false,
		},
		{
			name:    "Unknown Install",
			pt:      projects.TypeUnknown,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &projects.Project{Type: tt.pt}
			r := New(p)

			cmd, args, err := r.GetInstallCommand()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetInstallCommand() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if cmd != tt.wantCmd {
					t.Errorf("GetInstallCommand() cmd = %v, want %v", cmd, tt.wantCmd)
				}
				if len(args) != len(tt.wantArgs) {
					t.Errorf("GetInstallCommand() args len = %d, want %d", len(args), len(tt.wantArgs))
				}
			}
		})
	}
}

func TestGetLintCommand(t *testing.T) {
	tests := []struct {
		name     string
		pt       projects.ProjectType
		scripts  map[string]string
		fix      bool
		wantCmd  string
		wantArgs []string
		wantErr  bool
	}{
		{
			name:     "Node ESLint",
			pt:       projects.TypeNode,
			scripts:  nil,
			fix:      false,
			wantCmd:  "npx",
			wantArgs: []string{"eslint", "."},
			wantErr:  false,
		},
		{
			name:     "Node ESLint Fix",
			pt:       projects.TypeNode,
			scripts:  nil,
			fix:      true,
			wantCmd:  "npx",
			wantArgs: []string{"eslint", ".", "--fix"},
			wantErr:  false,
		},
		{
			name:     "Node NPM Lint Script",
			pt:       projects.TypeNode,
			scripts:  map[string]string{"lint": "eslint ."},
			fix:      false,
			wantCmd:  "npm",
			wantArgs: []string{"run", "lint"},
			wantErr:  false,
		},
		{
			name:     "Go Vet",
			pt:       projects.TypeGo,
			scripts:  nil,
			fix:      false,
			wantCmd:  "go",
			wantArgs: []string{"vet", "./..."},
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &projects.Project{
				Type:    tt.pt,
				Scripts: tt.scripts,
			}
			r := New(p)

			cmd, args, err := r.GetLintCommand(tt.fix)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLintCommand() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if cmd != tt.wantCmd {
					t.Errorf("GetLintCommand() cmd = %v, want %v", cmd, tt.wantCmd)
				}
				if len(args) != len(tt.wantArgs) {
					t.Errorf("GetLintCommand() args len = %d, want %d", len(args), len(tt.wantArgs))
				}
				if len(args) > 0 && len(tt.wantArgs) > 0 {
					// Check last arg for fix flag if relevant
					if tt.fix {
						lastArg := args[len(args)-1]
						expectedLast := tt.wantArgs[len(tt.wantArgs)-1]
						if lastArg != expectedLast {
							t.Errorf("GetLintCommand() last arg = %v, want %v", lastArg, expectedLast)
						}
					}
				}
			}
		})
	}
}
