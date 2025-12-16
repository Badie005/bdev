package projects

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

// ==============================================================
// ProjectType Tests
// ==============================================================

func TestProjectType_String(t *testing.T) {
	tests := []struct {
		ptype ProjectType
		want  string
	}{
		{TypeUnknown, "Unknown"},
		{TypeNextJS, "Next.js"},
		{TypeReact, "React"},
		{TypeVue, "Vue"},
		{TypeAngular, "Angular"},
		{TypeNode, "Node.js"},
		{TypeGo, "Go"},
		{TypePython, "Python"},
		{TypeRust, "Rust"},
		{TypeLaravel, "Laravel"},
		{TypeDjango, "Django"},
		{TypeStatic, "Static"},
		{TypeNuxt, "Nuxt"},
		{TypeSvelte, "Svelte"},
		{TypeAstro, "Astro"},
		{TypeJava, "Java"},
		{TypePHP, "PHP"},
		{TypeCpp, "C++"},
		{ProjectType(99), "Unknown"}, // Out of range
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := tt.ptype.String()
			if got != tt.want {
				t.Errorf("ProjectType(%d).String() = %q, want %q", tt.ptype, got, tt.want)
			}
		})
	}
}

func TestProjectType_Icon(t *testing.T) {
	tests := []struct {
		ptype ProjectType
		want  string
	}{
		{TypeUnknown, "[?]"},
		{TypeNextJS, "[NX]"},
		{TypeReact, "[RC]"},
		{TypeGo, "[GO]"},
		{TypePython, "[PY]"},
		{TypeRust, "[RS]"},
		{ProjectType(99), "[?]"}, // Out of range
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := tt.ptype.Icon()
			if got != tt.want {
				t.Errorf("ProjectType(%d).Icon() = %q, want %q", tt.ptype, got, tt.want)
			}
		})
	}
}

// ==============================================================
// DetectFromPackageJSON Tests
// ==============================================================

func TestDetectFromPackageJSON(t *testing.T) {
	tests := []struct {
		name string
		pkg  *PackageJSON
		want ProjectType
	}{
		{
			name: "next.js_project",
			pkg: &PackageJSON{
				Dependencies: map[string]string{"next": "14.0.0", "react": "18.2.0"},
			},
			want: TypeNextJS,
		},
		{
			name: "nuxt_project",
			pkg: &PackageJSON{
				Dependencies: map[string]string{"nuxt": "3.0.0"},
			},
			want: TypeNuxt,
		},
		{
			name: "angular_project",
			pkg: &PackageJSON{
				Dependencies: map[string]string{"@angular/core": "17.0.0"},
			},
			want: TypeAngular,
		},
		{
			name: "vue_project",
			pkg: &PackageJSON{
				Dependencies: map[string]string{"vue": "3.4.0"},
			},
			want: TypeVue,
		},
		{
			name: "svelte_project",
			pkg: &PackageJSON{
				DevDependencies: map[string]string{"svelte": "4.0.0"},
			},
			want: TypeSvelte,
		},
		{
			name: "astro_project",
			pkg: &PackageJSON{
				Dependencies: map[string]string{"astro": "4.0.0"},
			},
			want: TypeAstro,
		},
		{
			name: "react_only",
			pkg: &PackageJSON{
				Dependencies: map[string]string{"react": "18.2.0", "react-dom": "18.2.0"},
			},
			want: TypeReact,
		},
		{
			name: "plain_node",
			pkg: &PackageJSON{
				Dependencies: map[string]string{"express": "4.18.0"},
			},
			want: TypeNode,
		},
		{
			name: "empty_package",
			pkg: &PackageJSON{
				Dependencies:    map[string]string{},
				DevDependencies: map[string]string{},
			},
			want: TypeNode,
		},
		{
			name: "next_in_dev_deps",
			pkg: &PackageJSON{
				DevDependencies: map[string]string{"next": "14.0.0"},
			},
			want: TypeNextJS,
		},
		{
			name: "priority_next_over_react",
			pkg: &PackageJSON{
				Dependencies: map[string]string{"next": "14.0.0", "react": "18.2.0", "vue": "3.0.0"},
			},
			want: TypeNextJS, // next has higher priority
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DetectFromPackageJSON(tt.pkg)
			if got != tt.want {
				t.Errorf("DetectFromPackageJSON() = %v (%s), want %v (%s)",
					got, got.String(), tt.want, tt.want.String())
			}
		})
	}
}

// ==============================================================
// DetectFramework Tests
// ==============================================================

func TestDetectFramework(t *testing.T) {
	tests := []struct {
		name string
		pkg  *PackageJSON
		want string
	}{
		{
			name: "next_with_version",
			pkg: &PackageJSON{
				Dependencies: map[string]string{"next": "14.0.4"},
			},
			want: "next 14.0.4",
		},
		{
			name: "react_with_version",
			pkg: &PackageJSON{
				Dependencies: map[string]string{"react": "^18.2.0"},
			},
			want: "react ^18.2.0",
		},
		{
			name: "vue_in_devDeps",
			pkg: &PackageJSON{
				DevDependencies: map[string]string{"vue": "3.4.21"},
			},
			want: "vue 3.4.21",
		},
		{
			name: "no_framework",
			pkg: &PackageJSON{
				Dependencies: map[string]string{"express": "4.0.0"},
			},
			want: "",
		},
		{
			name: "empty_package",
			pkg:  &PackageJSON{},
			want: "",
		},
		{
			name: "priority_order",
			pkg: &PackageJSON{
				Dependencies: map[string]string{
					"next":  "14.0.0",
					"react": "18.2.0",
				},
			},
			want: "next 14.0.0", // next comes before react in priority
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DetectFramework(tt.pkg)
			if got != tt.want {
				t.Errorf("DetectFramework() = %q, want %q", got, tt.want)
			}
		})
	}
}

// ==============================================================
// Detect Tests (Integration with filesystem)
// ==============================================================

func TestDetect(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(t *testing.T, dir string)
		wantType ProjectType
	}{
		{
			name: "go_project",
			setup: func(t *testing.T, dir string) {
				createFile(t, dir, "go.mod", "module test\n\ngo 1.21\n")
				createFile(t, dir, "main.go", "package main\n")
			},
			wantType: TypeGo,
		},
		{
			name: "next_project",
			setup: func(t *testing.T, dir string) {
				createFile(t, dir, "next.config.js", "module.exports = {}\n")
			},
			wantType: TypeNextJS,
		},
		{
			name: "next_config_ts",
			setup: func(t *testing.T, dir string) {
				createFile(t, dir, "next.config.ts", "export default {}\n")
			},
			wantType: TypeNextJS,
		},
		{
			name: "rust_project",
			setup: func(t *testing.T, dir string) {
				createFile(t, dir, "Cargo.toml", "[package]\nname = \"test\"\n")
			},
			wantType: TypeRust,
		},
		{
			name: "python_project",
			setup: func(t *testing.T, dir string) {
				createFile(t, dir, "pyproject.toml", "[project]\nname = \"test\"\n")
			},
			wantType: TypePython,
		},
		{
			name: "requirements_txt",
			setup: func(t *testing.T, dir string) {
				createFile(t, dir, "requirements.txt", "flask\nrequests\n")
			},
			wantType: TypePython,
		},
		{
			name: "django_project",
			setup: func(t *testing.T, dir string) {
				createFile(t, dir, "manage.py", "#!/usr/bin/env python\n")
			},
			wantType: TypeDjango,
		},
		{
			name: "angular_project",
			setup: func(t *testing.T, dir string) {
				createFile(t, dir, "angular.json", "{}\n")
			},
			wantType: TypeAngular,
		},
		{
			name: "vue_config_project",
			setup: func(t *testing.T, dir string) {
				createFile(t, dir, "vue.config.js", "module.exports = {}\n")
			},
			wantType: TypeVue,
		},
		{
			name: "nuxt_project",
			setup: func(t *testing.T, dir string) {
				createFile(t, dir, "nuxt.config.js", "export default {}\n")
			},
			wantType: TypeNuxt,
		},
		{
			name: "svelte_project",
			setup: func(t *testing.T, dir string) {
				createFile(t, dir, "svelte.config.js", "export default {}\n")
			},
			wantType: TypeSvelte,
		},
		{
			name: "astro_project",
			setup: func(t *testing.T, dir string) {
				createFile(t, dir, "astro.config.mjs", "export default {}\n")
			},
			wantType: TypeAstro,
		},
		{
			name: "java_maven",
			setup: func(t *testing.T, dir string) {
				createFile(t, dir, "pom.xml", "<project></project>\n")
			},
			wantType: TypeJava,
		},
		{
			name: "java_gradle",
			setup: func(t *testing.T, dir string) {
				createFile(t, dir, "build.gradle", "plugins {}\n")
			},
			wantType: TypeJava,
		},
		{
			name: "cpp_cmake",
			setup: func(t *testing.T, dir string) {
				createFile(t, dir, "CMakeLists.txt", "cmake_minimum_required(VERSION 3.10)\n")
			},
			wantType: TypeCpp,
		},
		{
			name: "static_html",
			setup: func(t *testing.T, dir string) {
				createFile(t, dir, "index.html", "<!DOCTYPE html>\n")
			},
			wantType: TypeStatic,
		},
		{
			name: "package_json_react",
			setup: func(t *testing.T, dir string) {
				createFile(t, dir, "package.json", `{"dependencies":{"react":"18.0.0"}}`)
			},
			wantType: TypeReact,
		},
		{
			name: "extension_heuristic_go",
			setup: func(t *testing.T, dir string) {
				createFile(t, dir, "main.go", "package main\n")
				createFile(t, dir, "utils.go", "package main\n")
			},
			wantType: TypeGo,
		},
		{
			name: "extension_heuristic_python",
			setup: func(t *testing.T, dir string) {
				createFile(t, dir, "app.py", "print('hello')\n")
			},
			wantType: TypePython,
		},
		{
			name: "empty_dir",
			setup: func(t *testing.T, dir string) {
				// Empty directory
			},
			wantType: TypeUnknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			tt.setup(t, dir)

			got := Detect(dir)
			if got != tt.wantType {
				t.Errorf("Detect() = %v (%s), want %v (%s)",
					got, got.String(), tt.wantType, tt.wantType.String())
			}
		})
	}
}

// ==============================================================
// Command Tests
// ==============================================================

func TestProject_GetStartCommand(t *testing.T) {
	tests := []struct {
		name     string
		project  Project
		wantCmd  string
		wantArgs []string
	}{
		{
			name:     "go_project",
			project:  Project{Type: TypeGo},
			wantCmd:  "go",
			wantArgs: []string{"run", "."},
		},
		{
			name:     "python_project",
			project:  Project{Type: TypePython},
			wantCmd:  "python",
			wantArgs: []string{"main.py"},
		},
		{
			name:     "django_project",
			project:  Project{Type: TypeDjango},
			wantCmd:  "python",
			wantArgs: []string{"manage.py", "runserver"},
		},
		{
			name:     "laravel_project",
			project:  Project{Type: TypeLaravel},
			wantCmd:  "php",
			wantArgs: []string{"artisan", "serve"},
		},
		{
			name:     "rust_project",
			project:  Project{Type: TypeRust},
			wantCmd:  "cargo",
			wantArgs: []string{"run"},
		},
		{
			name: "nextjs_with_dev_script",
			project: Project{
				Type:    TypeNextJS,
				Scripts: map[string]string{"dev": "next dev"},
			},
			wantCmd:  "npm",
			wantArgs: []string{"run", "dev"},
		},
		{
			name: "node_with_start_script",
			project: Project{
				Type:    TypeNode,
				Scripts: map[string]string{"start": "node server.js"},
			},
			wantCmd:  "npm",
			wantArgs: []string{"start"},
		},
		{
			name:     "node_no_scripts",
			project:  Project{Type: TypeNode},
			wantCmd:  "node",
			wantArgs: []string{"index.js"},
		},
		{
			name:     "unknown_type",
			project:  Project{Type: TypeUnknown},
			wantCmd:  "",
			wantArgs: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCmd, gotArgs := tt.project.GetStartCommand()
			if gotCmd != tt.wantCmd {
				t.Errorf("GetStartCommand() cmd = %q, want %q", gotCmd, tt.wantCmd)
			}
			if !reflect.DeepEqual(gotArgs, tt.wantArgs) {
				t.Errorf("GetStartCommand() args = %v, want %v", gotArgs, tt.wantArgs)
			}
		})
	}
}

func TestProject_GetTestCommand(t *testing.T) {
	tests := []struct {
		name     string
		project  Project
		wantCmd  string
		wantArgs []string
	}{
		{
			name:     "go_project",
			project:  Project{Type: TypeGo},
			wantCmd:  "go",
			wantArgs: []string{"test", "./..."},
		},
		{
			name:     "python_project",
			project:  Project{Type: TypePython},
			wantCmd:  "pytest",
			wantArgs: nil,
		},
		{
			name:     "django_project",
			project:  Project{Type: TypeDjango},
			wantCmd:  "python",
			wantArgs: []string{"manage.py", "test"},
		},
		{
			name:     "laravel_project",
			project:  Project{Type: TypeLaravel},
			wantCmd:  "php",
			wantArgs: []string{"artisan", "test"},
		},
		{
			name:     "rust_project",
			project:  Project{Type: TypeRust},
			wantCmd:  "cargo",
			wantArgs: []string{"test"},
		},
		{
			name: "node_with_test_script",
			project: Project{
				Type:    TypeNode,
				Scripts: map[string]string{"test": "jest"},
			},
			wantCmd:  "npm",
			wantArgs: []string{"test"},
		},
		{
			name:     "node_no_test",
			project:  Project{Type: TypeNode},
			wantCmd:  "",
			wantArgs: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCmd, gotArgs := tt.project.GetTestCommand()
			if gotCmd != tt.wantCmd {
				t.Errorf("GetTestCommand() cmd = %q, want %q", gotCmd, tt.wantCmd)
			}
			if !reflect.DeepEqual(gotArgs, tt.wantArgs) {
				t.Errorf("GetTestCommand() args = %v, want %v", gotArgs, tt.wantArgs)
			}
		})
	}
}

func TestProject_GetBuildCommand(t *testing.T) {
	tests := []struct {
		name     string
		project  Project
		wantCmd  string
		wantArgs []string
	}{
		{
			name:     "go_project",
			project:  Project{Type: TypeGo, Name: "myapp"},
			wantCmd:  "go",
			wantArgs: []string{"build", "-o", "build/myapp"},
		},
		{
			name:     "rust_project",
			project:  Project{Type: TypeRust},
			wantCmd:  "cargo",
			wantArgs: []string{"build", "--release"},
		},
		{
			name:     "python_project",
			project:  Project{Type: TypePython},
			wantCmd:  "python",
			wantArgs: []string{"-m", "build"},
		},
		{
			name: "node_with_build",
			project: Project{
				Type:    TypeNode,
				Scripts: map[string]string{"build": "webpack"},
			},
			wantCmd:  "npm",
			wantArgs: []string{"run", "build"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCmd, gotArgs := tt.project.GetBuildCommand()
			if gotCmd != tt.wantCmd {
				t.Errorf("GetBuildCommand() cmd = %q, want %q", gotCmd, tt.wantCmd)
			}
			if !reflect.DeepEqual(gotArgs, tt.wantArgs) {
				t.Errorf("GetBuildCommand() args = %v, want %v", gotArgs, tt.wantArgs)
			}
		})
	}
}

func TestProject_ScriptList(t *testing.T) {
	tests := []struct {
		name    string
		scripts map[string]string
		want    []string
	}{
		{
			name:    "multiple_scripts",
			scripts: map[string]string{"dev": "next dev", "build": "next build", "start": "next start"},
			want:    []string{"build", "dev", "start"}, // Sorted alphabetically
		},
		{
			name:    "empty_scripts",
			scripts: map[string]string{},
			want:    []string{},
		},
		{
			name:    "nil_scripts",
			scripts: nil,
			want:    []string{},
		},
		{
			name:    "single_script",
			scripts: map[string]string{"test": "jest"},
			want:    []string{"test"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Project{Scripts: tt.scripts}
			got := p.ScriptList()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ScriptList() = %v, want %v", got, tt.want)
			}
		})
	}
}

// ==============================================================
// Integration Tests
// ==============================================================

func TestAnalyze(t *testing.T) {
	t.Run("go_project", func(t *testing.T) {
		dir := t.TempDir()
		createFile(t, dir, "go.mod", "module test\n\ngo 1.21\n")
		createFile(t, dir, "main.go", "package main\n")

		project := Analyze(dir)
		if project == nil {
			t.Fatal("Analyze() returned nil")
		}
		if project.Type != TypeGo {
			t.Errorf("Type = %v, want %v", project.Type, TypeGo)
		}
		if project.Name != filepath.Base(dir) {
			t.Errorf("Name = %q, want %q", project.Name, filepath.Base(dir))
		}
	})

	t.Run("nonexistent_path", func(t *testing.T) {
		project := Analyze("/nonexistent/path/12345")
		if project != nil {
			t.Error("Analyze() should return nil for nonexistent path")
		}
	})

	t.Run("nextjs_with_package_json", func(t *testing.T) {
		dir := t.TempDir()
		createFile(t, dir, "next.config.js", "module.exports = {}\n")
		createFile(t, dir, "package.json", `{
			"name": "my-next-app",
			"dependencies": {"next": "14.0.0", "react": "18.2.0"},
			"scripts": {"dev": "next dev", "build": "next build"}
		}`)

		project := Analyze(dir)
		if project == nil {
			t.Fatal("Analyze() returned nil")
		}
		if project.Type != TypeNextJS {
			t.Errorf("Type = %v, want %v", project.Type, TypeNextJS)
		}
		if project.Framework != "next 14.0.0" {
			t.Errorf("Framework = %q, want 'next 14.0.0'", project.Framework)
		}
		if len(project.Scripts) == 0 {
			t.Error("Scripts should be populated from package.json")
		}
	})
}

func TestScan(t *testing.T) {
	t.Run("multiple_projects", func(t *testing.T) {
		rootDir := t.TempDir()

		// Create Go project
		goDir := filepath.Join(rootDir, "mygoapp")
		os.Mkdir(goDir, 0o755)
		createFile(t, goDir, "go.mod", "module mygoapp\n\ngo 1.21\n")

		// Create Node project
		nodeDir := filepath.Join(rootDir, "mynodeapp")
		os.Mkdir(nodeDir, 0o755)
		createFile(t, nodeDir, "package.json", `{"name": "mynodeapp", "dependencies": {"express": "4.0.0"}}`)

		// Create hidden directory (should be skipped)
		os.Mkdir(filepath.Join(rootDir, ".hidden"), 0o755)

		// Create node_modules (should be skipped)
		os.Mkdir(filepath.Join(rootDir, "node_modules"), 0o755)

		projects, err := Scan(rootDir)
		if err != nil {
			t.Fatalf("Scan() error = %v", err)
		}

		if len(projects) != 2 {
			t.Errorf("Scan() found %d projects, want 2", len(projects))
		}
	})

	t.Run("empty_directory", func(t *testing.T) {
		rootDir := t.TempDir()

		projects, err := Scan(rootDir)
		if err != nil {
			t.Fatalf("Scan() error = %v", err)
		}

		if len(projects) != 0 {
			t.Errorf("Scan() found %d projects, want 0", len(projects))
		}
	})

	t.Run("nonexistent_directory", func(t *testing.T) {
		_, err := Scan("/nonexistent/path/12345")
		if err == nil {
			t.Error("Scan() should return error for nonexistent path")
		}
	})
}

func TestExtensionHeuristic(t *testing.T) {
	tests := []struct {
		name     string
		files    map[string]string
		wantType ProjectType
	}{
		{
			name: "go_files",
			files: map[string]string{
				"main.go":  "package main\n",
				"utils.go": "package main\n",
			},
			wantType: TypeGo,
		},
		{
			name: "python_files",
			files: map[string]string{
				"app.py":   "print('hello')\n",
				"utils.py": "def foo(): pass\n",
			},
			wantType: TypePython,
		},
		{
			name: "js_files",
			files: map[string]string{
				"index.js": "console.log('hello');\n",
				"utils.js": "export const foo = 1;\n",
			},
			wantType: TypeNode,
		},
		{
			name: "ts_files",
			files: map[string]string{
				"index.ts": "const x: number = 1;\n",
			},
			wantType: TypeNode,
		},
		{
			name: "java_files",
			files: map[string]string{
				"Main.java": "public class Main {}\n",
			},
			wantType: TypeJava,
		},
		{
			name: "php_files",
			files: map[string]string{
				"index.php": "<?php echo 'hello'; ?>\n",
			},
			wantType: TypePHP,
		},
		{
			name: "cpp_files",
			files: map[string]string{
				"main.cpp": "#include <iostream>\n",
				"utils.h":  "#pragma once\n",
			},
			wantType: TypeCpp,
		},
		{
			name:     "empty_dir",
			files:    map[string]string{},
			wantType: TypeUnknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			for name, content := range tt.files {
				createFile(t, dir, name, content)
			}

			got := ExtensionHeuristic(dir)
			if got != tt.wantType {
				t.Errorf("ExtensionHeuristic() = %v (%s), want %v (%s)",
					got, got.String(), tt.wantType, tt.wantType.String())
			}
		})
	}
}

func TestReadPackageJSON(t *testing.T) {
	t.Run("valid_package_json", func(t *testing.T) {
		dir := t.TempDir()
		createFile(t, dir, "package.json", `{
			"name": "test-app",
			"version": "1.0.0",
			"dependencies": {"react": "18.0.0"},
			"devDependencies": {"typescript": "5.0.0"},
			"scripts": {"test": "jest"}
		}`)

		pkg := ReadPackageJSON(dir)
		if pkg == nil {
			t.Fatal("ReadPackageJSON() returned nil")
		}
		if pkg.Name != "test-app" {
			t.Errorf("Name = %q, want 'test-app'", pkg.Name)
		}
		if pkg.Version != "1.0.0" {
			t.Errorf("Version = %q, want '1.0.0'", pkg.Version)
		}
		if pkg.Dependencies["react"] != "18.0.0" {
			t.Errorf("Dependencies['react'] = %q, want '18.0.0'", pkg.Dependencies["react"])
		}
	})

	t.Run("no_package_json", func(t *testing.T) {
		dir := t.TempDir()
		pkg := ReadPackageJSON(dir)
		if pkg != nil {
			t.Error("ReadPackageJSON() should return nil when file doesn't exist")
		}
	})

	t.Run("invalid_json", func(t *testing.T) {
		dir := t.TempDir()
		createFile(t, dir, "package.json", "{ invalid json }")

		pkg := ReadPackageJSON(dir)
		if pkg != nil {
			t.Error("ReadPackageJSON() should return nil for invalid JSON")
		}
	})
}

// ==============================================================
// Test Helpers
// ==============================================================

func createFile(t *testing.T, dir, name, content string) {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("failed to create dir: %v", err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write file: %v", err)
	}
}
