package projects

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/badie/bdev/internal/core/git"
)

// ProjectType represents the type of project
type ProjectType int

const (
	TypeUnknown ProjectType = iota
	TypeNextJS
	TypeReact
	TypeVue
	TypeAngular
	TypeNode
	TypeGo
	TypePython
	TypeRust
	TypeLaravel
	TypeDjango
	TypeStatic
	TypeNuxt
	TypeSvelte
	TypeAstro
	TypeJava
	TypePHP
	TypeCpp
)

// String returns the string representation of ProjectType
func (t ProjectType) String() string {
	names := []string{
		"Unknown", "Next.js", "React", "Vue", "Angular",
		"Node.js", "Go", "Python", "Rust", "Laravel",
		"Django", "Static", "Nuxt", "Svelte", "Astro",
		"Java", "PHP", "C++",
	}
	if int(t) < len(names) {
		return names[t]
	}
	return "Unknown"
}

// Icon returns a text abbreviation for the project type
func (t ProjectType) Icon() string {
	abbrevs := []string{
		"[?]", "[NX]", "[RC]", "[VU]", "[NG]",
		"[JS]", "[GO]", "[PY]", "[RS]", "[LV]",
		"[DJ]", "[--]", "[NU]", "[SV]", "[AS]",
		"[JV]", "[PH]", "[C+]",
	}
	if int(t) < len(abbrevs) {
		return abbrevs[t]
	}
	return "[?]"
}

// Project represents a development project
type Project struct {
	Name         string            `json:"name"`
	Path         string            `json:"path"`
	Type         ProjectType       `json:"type"`
	Framework    string            `json:"framework"`
	LastModified time.Time         `json:"last_modified"`
	Size         int64             `json:"size"`
	HasGit       bool              `json:"has_git"`
	GitBranch    string            `json:"git_branch,omitempty"`
	Scripts      map[string]string `json:"scripts,omitempty"`
}

// PackageJSON represents a Node.js package.json
type PackageJSON struct {
	Name            string            `json:"name"`
	Version         string            `json:"version"`
	Scripts         map[string]string `json:"scripts"`
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
}

// GoMod represents a Go module
type GoMod struct {
	Module  string
	GoVer   string
	Require []string
}

// Scan scans a directory for projects
func Scan(rootDir string) ([]Project, error) {
	projects := make([]Project, 0)

	entries, err := os.ReadDir(rootDir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		// Skip hidden directories
		if strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		// Skip common non-project directories
		skip := []string{"node_modules", "vendor", "__pycache__", ".git", "build", "dist"}
		shouldSkip := false
		for _, s := range skip {
			if entry.Name() == s {
				shouldSkip = true
				break
			}
		}
		if shouldSkip {
			continue
		}

		projectPath := filepath.Join(rootDir, entry.Name())
		project := Analyze(projectPath)
		if project != nil {
			projects = append(projects, *project)
		}
	}

	// Sort by last modified (most recent first)
	sort.Slice(projects, func(i, j int) bool {
		return projects[i].LastModified.After(projects[j].LastModified)
	})

	return projects, nil
}

// Analyze analyzes a directory and returns project info
func Analyze(path string) *Project {
	info, err := os.Stat(path)
	if err != nil {
		return nil
	}

	project := &Project{
		Name:         filepath.Base(path),
		Path:         path,
		Type:         Detect(path),
		LastModified: info.ModTime(),
		HasGit:       git.IsRepo(path),
	}

	// Get framework from package.json
	if pkg := ReadPackageJSON(path); pkg != nil {
		project.Framework = DetectFramework(pkg)
		project.Scripts = pkg.Scripts
	}

	// Get git branch
	if project.HasGit {
		if repo, err := git.Open(path); err == nil {
			project.GitBranch = repo.CurrentBranch()
		}
	}

	return project
}
func Detect(path string) ProjectType {
	// 1. PRIMARY: Check for specific config files
	checks := []struct {
		file    string
		ptype   ProjectType
		content func(string) ProjectType
	}{
		{"next.config.js", TypeNextJS, nil},
		{"next.config.ts", TypeNextJS, nil},
		{"next.config.mjs", TypeNextJS, nil},
		{"nuxt.config.js", TypeNuxt, nil},
		{"nuxt.config.ts", TypeNuxt, nil},
		{"svelte.config.js", TypeSvelte, nil},
		{"astro.config.mjs", TypeAstro, nil},
		{"vue.config.js", TypeVue, nil},
		{"angular.json", TypeAngular, nil},
		{"bun.lockb", TypeNode, nil}, // Bun project, usually JS/TS
		{"go.mod", TypeGo, nil},
		{"Cargo.toml", TypeRust, nil},
		{"composer.json", TypePHP, detectComposer}, // Generic PHP or Laravel
		{"pyproject.toml", TypePython, nil},
		{"requirements.txt", TypePython, nil},
		{"manage.py", TypeDjango, nil},
		{"pom.xml", TypeJava, nil},
		{"build.gradle", TypeJava, nil},
		{"CMakeLists.txt", TypeCpp, nil},
		{"Makefile", TypeCpp, nil}, // Heuristic, usually C/C++ or Go
		{"index.html", TypeStatic, nil},
	}

	for _, check := range checks {
		filePath := filepath.Join(path, check.file)
		if _, err := os.Stat(filePath); err == nil {
			if check.content != nil {
				return check.content(filePath)
			}
			return check.ptype
		}
	}

	// 2. SECONDARY: Check package.json for React/Node
	if pkg := ReadPackageJSON(path); pkg != nil {
		return DetectFromPackageJSON(pkg)
	}

	// 3. TERTIARY: Extension Heuristics (Deep Scan)
	return ExtensionHeuristic(path)
}

// ExtensionHeuristic scans for dominant file extensions
func ExtensionHeuristic(path string) ProjectType {
	// Limit scan to top level or depth 1
	entries, err := os.ReadDir(path)
	if err != nil {
		return TypeUnknown
	}

	counts := make(map[string]int)
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		ext := filepath.Ext(e.Name())
		counts[ext]++
	}

	if counts[".go"] > 0 {
		return TypeGo
	}
	if counts[".py"] > 0 {
		return TypePython
	}
	if counts[".js"] > 0 || counts[".ts"] > 0 {
		return TypeNode
	}
	if counts[".java"] > 0 {
		return TypeJava
	}
	if counts[".php"] > 0 {
		return TypePHP
	}
	if counts[".cpp"] > 0 || counts[".c"] > 0 || counts[".h"] > 0 {
		return TypeCpp
	}

	return TypeUnknown
}

// detectComposer checks composer.json for framework
func detectComposer(path string) ProjectType {
	data, err := os.ReadFile(path)
	if err != nil {
		return TypePHP
	}
	content := string(data)
	if strings.Contains(content, "laravel/framework") {
		return TypeLaravel
	}
	return TypePHP
}

// ReadPackageJSON reads and parses package.json
func ReadPackageJSON(path string) *PackageJSON {
	pkgPath := filepath.Join(path, "package.json")
	data, err := os.ReadFile(pkgPath)
	if err != nil {
		return nil
	}

	var pkg PackageJSON
	if err := json.Unmarshal(data, &pkg); err != nil {
		return nil
	}

	return &pkg
}

// DetectFromPackageJSON detects project type from package.json dependencies
func DetectFromPackageJSON(pkg *PackageJSON) ProjectType {
	deps := make(map[string]bool)
	for k := range pkg.Dependencies {
		deps[k] = true
	}
	for k := range pkg.DevDependencies {
		deps[k] = true
	}

	// Priority order
	if deps["next"] {
		return TypeNextJS
	}
	if deps["nuxt"] {
		return TypeNuxt
	}
	if deps["@angular/core"] {
		return TypeAngular
	}
	if deps["vue"] {
		return TypeVue
	}
	if deps["svelte"] {
		return TypeSvelte
	}
	if deps["astro"] {
		return TypeAstro
	}
	if deps["react"] {
		return TypeReact
	}

	return TypeNode
}

// DetectFramework returns the framework name with version
func DetectFramework(pkg *PackageJSON) string {
	frameworks := []string{"next", "nuxt", "@angular/core", "vue", "react", "svelte", "astro"}

	for _, fw := range frameworks {
		if ver, ok := pkg.Dependencies[fw]; ok {
			return fw + " " + ver
		}
		if ver, ok := pkg.DevDependencies[fw]; ok {
			return fw + " " + ver
		}
	}

	return ""
}

// GetStartCommand returns the start command for a project
func (p *Project) GetStartCommand() (string, []string) {
	switch p.Type {
	case TypeNextJS, TypeReact, TypeVue, TypeNuxt, TypeSvelte, TypeAstro:
		if _, ok := p.Scripts["dev"]; ok {
			return "npm", []string{"run", "dev"}
		}
		if _, ok := p.Scripts["start"]; ok {
			return "npm", []string{"start"}
		}
	case TypeNode:
		if _, ok := p.Scripts["dev"]; ok {
			return "npm", []string{"run", "dev"}
		}
		if _, ok := p.Scripts["start"]; ok {
			return "npm", []string{"start"}
		}
		return "node", []string{"index.js"}
	case TypeGo:
		return "go", []string{"run", "."}
	case TypePython:
		return "python", []string{"main.py"}
	case TypeDjango:
		return "python", []string{"manage.py", "runserver"}
	case TypeLaravel:
		return "php", []string{"artisan", "serve"}
	case TypeRust:
		return "cargo", []string{"run"}
	}

	return "", nil
}

// GetTestCommand returns the test command for a project
func (p *Project) GetTestCommand() (string, []string) {
	switch p.Type {
	case TypeNextJS, TypeReact, TypeVue, TypeNuxt, TypeSvelte, TypeAstro, TypeNode:
		if _, ok := p.Scripts["test"]; ok {
			return "npm", []string{"test"}
		}
	case TypeGo:
		return "go", []string{"test", "./..."}
	case TypePython:
		return "pytest", nil
	case TypeDjango:
		return "python", []string{"manage.py", "test"}
	case TypeLaravel:
		return "php", []string{"artisan", "test"}
	case TypeRust:
		return "cargo", []string{"test"}
	}

	return "", nil
}

// GetBuildCommand returns the build command for a project
func (p *Project) GetBuildCommand() (string, []string) {
	switch p.Type {
	case TypeNextJS, TypeReact, TypeVue, TypeNuxt, TypeSvelte, TypeAstro, TypeNode:
		if _, ok := p.Scripts["build"]; ok {
			return "npm", []string{"run", "build"}
		}
	case TypeGo:
		return "go", []string{"build", "-o", "build/" + p.Name}
	case TypePython:
		return "python", []string{"-m", "build"}
	case TypeRust:
		return "cargo", []string{"build", "--release"}
	}

	return "", nil
}

// ScriptList returns available npm scripts
func (p *Project) ScriptList() []string {
	scripts := make([]string, 0, len(p.Scripts))
	for k := range p.Scripts {
		scripts = append(scripts, k)
	}
	sort.Strings(scripts)
	return scripts
}
