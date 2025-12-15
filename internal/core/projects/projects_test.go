package projects

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDetect(t *testing.T) {
	tests := []struct {
		name     string
		files    []string
		expected ProjectType
	}{
		{
			name:     "Go Project",
			files:    []string{"go.mod"},
			expected: TypeGo,
		},
		{
			name:     "Node Project",
			files:    []string{"package.json"},
			expected: TypeNode,
		},
		{
			name:     "Python Project",
			files:    []string{"requirements.txt"},
			expected: TypePython,
		},
		{
			name:     "Next.js Project",
			files:    []string{"package.json", "next.config.js"},
			expected: TypeNextJS,
		},
		{
			name:     "Unknown Project",
			files:    []string{"random.txt"},
			expected: TypeUnknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp dir
			dir := t.TempDir()

			// Create files
			for _, file := range tt.files {
				path := filepath.Join(dir, file)
				if err := os.WriteFile(path, []byte("{}"), 0644); err != nil {
					t.Fatalf("Failed to create file %s: %v", file, err)
				}
			}

			// Test detection
			got := Detect(dir)
			if got != tt.expected {
				t.Errorf("Detect() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestProjectType_String(t *testing.T) {
	tests := []struct {
		pt   ProjectType
		want string
	}{
		{TypeGo, "Go"},
		{TypeNode, "Node.js"},
		{TypePython, "Python"},
		{TypeUnknown, "Unknown"},
	}

	for _, tt := range tests {
		if got := tt.pt.String(); got != tt.want {
			t.Errorf("ProjectType.String() = %v, want %v", got, tt.want)
		}
	}
}

func TestGetStartCommand(t *testing.T) {
	// Create temp dir for package.json parsing test
	dir := t.TempDir()
	pkgJSON := `{"scripts": {"dev": "next dev"}}`
	os.WriteFile(filepath.Join(dir, "package.json"), []byte(pkgJSON), 0644)

	tests := []struct {
		name     string
		pt       ProjectType
		path     string
		wantCmd  string
		wantArgs []string
	}{
		{
			name:     "Go Run",
			pt:       TypeGo,
			path:     "",
			wantCmd:  "go",
			wantArgs: []string{"run", "."},
		},
		{
			name:     "Python Run",
			pt:       TypePython,
			path:     "",
			wantCmd:  "python",
			wantArgs: []string{"main.py"},
		},
		{
			name:     "Next.js Dev",
			pt:       TypeNextJS,
			path:     dir, // Should parse package.json
			wantCmd:  "npm",
			wantArgs: []string{"run", "dev"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Project{Type: tt.pt, Path: tt.path}
			// Load scripts if path is set (simulating Analyze)
			if tt.path != "" {
				if pkg := ReadPackageJSON(tt.path); pkg != nil {
					p.Scripts = pkg.Scripts
				}
			}

			cmd, args := p.GetStartCommand()
			if cmd != tt.wantCmd {
				t.Errorf("GetStartCommand() cmd = %v, want %v", cmd, tt.wantCmd)
			}

			if len(args) != len(tt.wantArgs) {
				t.Errorf("GetStartCommand() args length = %d, want %d", len(args), len(tt.wantArgs))
			}
		})
	}
}
